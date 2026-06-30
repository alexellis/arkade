package cmd

import (
	"testing"

	"github.com/alexellis/arkade/pkg/get"
)

func makeTestTools() []get.Tool {
	return []get.Tool{
		{Name: "helm", Owner: "helm", Repo: "helm", Description: "The Kubernetes Package Manager"},
		{Name: "faas-cli", Owner: "openfaas", Repo: "faas-cli", Description: "CLI for OpenFaaS"},
		{Name: "kubectl", Owner: "kubernetes", Repo: "kubernetes", Description: "Control plane CLI"},
	}
}

func Test_ExactNameMatchRanksFirst(t *testing.T) {
	ranked := rankByTFIDF(makeTestTools(), "helm")

	if ranked[0].Tool.Name != "helm" {
		t.Errorf("expected helm to rank first, got %s", ranked[0].Tool.Name)
	}
}

func Test_MultiWordQueryBoost(t *testing.T) {
	ranked := rankByTFIDF(makeTestTools(), "kubernetes package")

	var found int
	for _, r := range ranked {
		if r.Tool.Name == "helm" && r.Score > 0 {
			found++
		}
	}
	if found != 1 {
		t.Error("expected helm to match 'kubernetes package' query")
	}
}

func Test_OwnerRepoInTFIDF(t *testing.T) {
	ranked := rankByTFIDF(makeTestTools(), "openfaas")

	var found bool
	for _, r := range ranked {
		if r.Tool.Name == "faas-cli" && r.Score > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected faas-cli to match 'openfaas' via Owner field")
	}
}

func Test_FallbackFiresWhenTFIDFFails(t *testing.T) {
	matches := fuzzySubstringFallback(makeTestTools(), []string{"kube"})

	var found bool
	for _, r := range matches {
		if r.Tool.Name == "kubectl" && r.Score > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected kubectl to match 'kube' via substring fallback on Owner")
	}
}

func Test_LevenshteinFuzzyNearMiss(t *testing.T) {
	score := levenshteinFuzzyScore([]string{"faasd"}, "faas-cli")
	if score <= 0 {
		t.Error("expected positive fuzzy score for 'faasd' vs 'faas-cli'")
	}
}

func Test_FallbackReturnsEmptyForNoMatch(t *testing.T) {
	matches := fuzzySubstringFallback(makeTestTools(), []string{"xyznonexistent"})
	if len(matches) != 0 {
		t.Errorf("expected no fallback matches, got %d", len(matches))
	}
}

func Test_FallbackSortedByScore(t *testing.T) {
	matches := fuzzySubstringFallback(makeTestTools(), []string{"cli"})
	for i := 1; i < len(matches); i++ {
		if matches[i].Score > matches[i-1].Score {
			t.Errorf("results not sorted by score: %.2f > %.2f", matches[i].Score, matches[i-1].Score)
		}
	}
}
