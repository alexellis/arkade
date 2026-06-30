// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg/get"
)

type scoreRank struct {
	Tool  get.Tool
	Score float64
}

// aliases maps common shorthand to full term. Each alias expands to one or more
// space-separated terms so that multi-word expansions are scored correctly.
var aliasMap = map[string]string{
	"k8s":    "kubernetes",
	"kube":   "kubernetes",
	"eksctl": "amazon eks kubernetes cluster management",
	"gke":    "google kubernetes engine",
	"aks":    "azure kubernetes service",
}

func MakeSearch() *cobra.Command {
	tools := get.MakeTools()

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: `Search for a tool available in arkade get`,
		Long:  `Search for tools by name or description using relevance ranking. Tools that share keywords with your query are ranked first. Common aliases like k8s are expanded to kubernetes, and fuzzy matching finds similar names (e.g., "openfaas" matches faas-cli). Multi-word queries match and rank tools containing multiple terms higher.`,
		Example: `  arkade search helm

   # Expand "k8s" to Kubernetes and rank by relevance
   arkade search k8s

   # Fuzzy name matching — finds faas-cli even though the user types "openfaas"
   arkade search openfaas

   # Multi-word query (tools with both words ranked higher)
   arkade search container runtime

   # Show as a list instead of table
   arkade search helm --format list`,
	}

	cmd.Flags().String("format", "table", "Output format: list, markdown or table")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			return errors.New("please provide a search query")
		}

		format, _ := cmd.Flags().GetString("format")

		ranked := rankByTFIDF(tools, query)

		sort.SliceStable(ranked, func(i, j int) bool {
			return ranked[i].Score > ranked[j].Score
		})

		var matches []scoreRank
		for _, r := range ranked {
			if r.Score > 0 {
				matches = append(matches, r)
			}
		}

		// Last resort: substring fallback on Name, Owner and Repo when TF-IDF found nothing.
		if len(matches) == 0 {
			queryTerms := tokenize(expandAliases(query))
			matches = fuzzySubstringFallback(tools, queryTerms)
		}

		if len(matches) == 0 {
			cmd.Printf("No tools found matching \"%s\"\n", query)
			return nil
		}

		switch format {
		case "list":
			for i, r := range matches {
				fmt.Printf("%d. %s (%.3f)\t%s\n", i+1, r.Tool.Name, r.Score, r.Tool.Description)
			}
		case "markdown":
			fmt.Println("| Rank | Name | Score | Description |")
			fmt.Println("|------|------|-------|-------------|")
			for i, r := range matches {
				fmt.Printf("| %d | %s | %.3f | %s |\n", i+1, r.Tool.Name, r.Score, r.Tool.Description)
			}
		default:
			matchesOnly := make([]get.Tool, len(matches))
			for i, r := range matches {
				matchesOnly[i] = r.Tool
			}
			cmd.Printf("Found %d tool(s) matching \"%s\":\n\n", len(matches), query)
			get.CreateToolsTable(matchesOnly, get.TableStyle)
		}

		return nil
	}

	return cmd
}

// fuzzySubstringFallback does a simple case-insensitive substring match across
// Name, Owner, Repo (not Description) when the TF-IDF index returned no results.
func fuzzySubstringFallback(tools []get.Tool, queryTerms []string) []scoreRank {
	results := make([]scoreRank, 0)
	for _, t := range tools {
		var score float64
		for _, q := range queryTerms {
			if strings.Contains(strings.ToLower(t.Name), q) {
				score += 3.0
			} else if strings.Contains(strings.ToLower(t.Owner), q) {
				score += 2.0
			} else if strings.Contains(strings.ToLower(t.Repo), q) {
				score += 1.5
			}
		}
		if score > 0 {
			results = append(results, scoreRank{Tool: t, Score: score})
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

// splitOnSeparators returns a slice of substrings obtained by splitting on - and _.
func splitOnSeparators(s string) []string {
	parts := strings.Split(strings.ToLower(s), "-")
	result := make([]string, 0)
	for _, p := range parts {
		subparts := strings.Split(p, "_")
		for _, sp := range subparts {
			if sp != "" {
				result = append(result, sp)
			}
		}
	}
	return result
}

func rankByTFIDF(tools []get.Tool, rawQuery string) []scoreRank {
	docs := make([][]string, len(tools))
	for i, t := range tools {
		// Tokenize name with separator splitting to create individual IDF entries
		// for parts like "faas" in "faas-cli", then append the expanded description.
		nameTokens := splitOnSeparators(t.Name)
		docs[i] = tokenize(
			strings.Join(nameTokens, " ") + " " +
				t.Owner + " " +
				t.Repo + " " +
				expandAliases(t.Description),
		)
	}

	df := make(map[string]int)
	for _, doc := range docs {
		seen := map[string]bool{}
		for _, w := range doc {
			if !seen[w] {
				df[w]++
				seen[w] = true
			}
		}
	}

	nDocs := float64(len(tools))
	idf := make(map[string]float64)
	for t, freq := range df {
		idf[t] = math.Log(nDocs/float64(freq)) + 1.0
	}

	queryTerms := tokenize(expandAliases(rawQuery))

	scores := make([]scoreRank, len(tools))
	for i, t := range tools {
		tf := termFreq(docs[i])

		var score float64
		termsMatched := 0

		// TF-IDF contribution from name + description.
		for _, q := range queryTerms {
			if tf[q] > 0 {
				score += tf[q] * idf[q]
				termsMatched++
			}
		}

		// Exact name match bonus: if the tool name equals any query term (or vice versa),
		// give it a very high score so it ranks first.
		lowerName := strings.ToLower(t.Name)
		for _, q := range queryTerms {
			if lowerName == q {
				score += 5.0
			}
		}

		// Substring name bonus: if a query term is a substring of the tool name but not
		// matched by TF-IDF (because it's not an independent token), score using the
		// best available IDF weight from that tool's own tokens.
		for _, q := range queryTerms {
			if tf[q] == 0 && strings.Contains(lowerName, q) {
				// Use a fallback weight: log(N/1)+1 which is the maximum possible IDF,
				// giving partial name matches strong relevance.
				maxIDF := math.Log(nDocs) + 1.0
				score += maxIDF * 2.0
				termsMatched++
			} else if tf[q] > 0 {
				// Name and description both matched — slight extra boost.
				score += idf[q] * 0.5
			}
		}

		// Levenshtein fuzzy matching on tool name parts only.
		fuzzyBonus := levenshteinFuzzyScore(queryTerms, t.Name)
		score += fuzzyBonus
		if fuzzyBonus > 0 {
			termsMatched++
		}

		// Multi-word query boost: tools that match multiple distinct query terms are
		// ranked higher. The boost is proportional to the fraction of query terms matched.
		if len(queryTerms) > 1 {
			frac := float64(termsMatched) / float64(len(queryTerms))
			score += score * frac * 0.5
		}

		scores[i] = scoreRank{Tool: t, Score: score}
	}

	return scores
}

// levenshteinFuzzyScore computes a bonus for tools whose name contains
// words within edit distance of any query token that aren't already matched by exact text.
func levenshteinFuzzyScore(queryTerms []string, toolName string) float64 {
	if len(queryTerms) == 0 {
		return 0
	}

	nameLower := strings.ToLower(toolName)
	// Use separator-split name parts for comparison rather than whitespace tokens.
	nameWords := splitOnSeparators(nameLower)

	var bonus float64
	distMax := 2

	for _, q := range queryTerms {
		if len(q) < 5 {
			continue
		}

		// Quick exact substring check — if already matched, no need to fuzzy.
		if strings.Contains(nameLower, q) {
			continue
		}

		bestDist := distMax + 1
		for _, w := range nameWords {
			if len(w) < 4 {
				continue
			}
			d := levenshteinDistance(q, w)
			if d < bestDist && d <= distMax {
				bestDist = d
			}
		}

		if bestDist > 0 && bestDist <= distMax {
			bonus += float64(distMax-bestDist+1) * 0.8
		}
	}

	return bonus
}

func levenshteinDistance(a, b string) int {
	aLen := len(a)
	bLen := len(b)

	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}

	dp := make([]int, bLen+1)
	for j := 0; j <= bLen; j++ {
		dp[j] = j
	}

	for i := 1; i <= aLen; i++ {
		prevDiag := dp[0]
		dp[0] = i
		for j := 1; j <= bLen; j++ {
			temp := dp[j]
			if a[i-1] == b[j-1] {
				dp[j] = prevDiag
			} else {
				m := dp[j-1] + 1
				if d := dp[j] + 1; d < m {
					m = d
				}
				if r := prevDiag + 1; r < m {
					m = r
				}
				dp[j] = m
			}
			prevDiag = temp
		}
	}

	return dp[bLen]
}

func expandAliases(s string) string {
	s = strings.ToLower(s)
	for alias, expansion := range aliasMap {
		padded := " " + s + " "
		s = strings.ReplaceAll(padded, " "+alias+" ", " "+expansion+" ")
	}
	return strings.Trim(s, " ")
}

func tokenize(s string) []string {
	lower := strings.ToLower(s)
	fields := strings.Fields(lower)
	cleaned := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.Trim(f, ".,:;()'\"")
		if f != "" {
			cleaned = append(cleaned, f)
		}
	}
	return cleaned
}

func termFreq(words []string) map[string]float64 {
	counts := make(map[string]int)
	for _, w := range words {
		counts[w]++
	}
	tf := make(map[string]float64)
	n := float64(len(words))
	if n == 0 {
		return tf
	}
	for w, c := range counts {
		tf[w] = float64(c) / n
	}
	return tf
}
