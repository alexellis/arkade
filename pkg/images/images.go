package images

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-containerregistry/pkg/crane"
)

// TagAttributes describes the format of a container image tag
type TagAttributes struct {
	HasSuffix bool
	HasMajor  bool
	HasMinor  bool
	HasPatch  bool
	Original  string
}

// AttributesMatch returns true if the tag attributes match the same version format
func (c *TagAttributes) AttributesMatch(n TagAttributes) bool {
	return c.HasMajor == n.HasMajor &&
		c.HasMinor == n.HasMinor &&
		c.HasPatch == n.HasPatch &&
		c.HasSuffix == n.HasSuffix
}

// SplitImageName splits an image reference into name and tag parts
func SplitImageName(reposName string) (string, string) {
	nameParts := strings.SplitN(reposName, ":", 2)
	return nameParts[0], nameParts[1]
}

// UpdateImage checks if a newer version of an image exists and returns the updated reference
func UpdateImage(iName string, verbose bool) (bool, string, error) {
	imageName, tag := SplitImageName(iName)
	ref, err := crane.ListTags(imageName)
	if err != nil {
		return false, iName, errors.New("unable to list tags for " + imageName)
	}

	candidateTag, hasSemVerTag := GetCandidateTag(ref, tag)
	if !hasSemVerTag {
		return false, iName, fmt.Errorf("no valid semver tags of current format found for %s", imageName)
	}

	updated := false

	// Don't upgrade to an RC tag, even if it's newer.
	if TagIsUpgradeable(tag, candidateTag) {
		updated = true
		iName = fmt.Sprintf("%s:%s", imageName, candidateTag)
		if verbose {
			log.Printf("[%s] %s => %s", imageName, tag, candidateTag)
		}
	}

	return updated, iName, nil
}

// UpdateImagePinned checks for a newer patch version within the same major.minor range.
// For example, golang:1.24 would find 1.24.4 but not 1.25.
func UpdateImagePinned(iName string, verbose bool) (bool, string, error) {
	imageName, tag := SplitImageName(iName)

	currentVer, err := semver.NewVersion(tag)
	if err != nil {
		return false, iName, fmt.Errorf("unable to parse tag %s as semver: %s", tag, err)
	}

	remoteTags, err := crane.ListTags(imageName)
	if err != nil {
		return false, iName, errors.New("unable to list tags for " + imageName)
	}

	prefix := fmt.Sprintf("%d.%d.", currentVer.Major(), currentVer.Minor())
	unpatchedTag := fmt.Sprintf("%d.%d", currentVer.Major(), currentVer.Minor())

	var candidates []*semver.Version
	for _, t := range remoteTags {
		if !strings.HasPrefix(t, prefix) && t != unpatchedTag {
			continue
		}

		v, err := semver.NewVersion(t)
		if err != nil {
			continue
		}

		if v.Prerelease() != currentVer.Prerelease() {
			continue
		}

		candidates = append(candidates, v)
	}

	if len(candidates) == 0 {
		return false, iName, fmt.Errorf("no valid tags found for %s within %s.x", imageName, unpatchedTag)
	}

	sort.Sort(sort.Reverse(semver.Collection(candidates)))
	latest := candidates[0]

	if TagIsUpgradeable(tag, latest.Original()) {
		newRef := fmt.Sprintf("%s:%s", imageName, latest.Original())
		if verbose {
			log.Printf("[%s] %s => %s (pinned to %s.x)", imageName, tag, latest.Original(), unpatchedTag)
		}
		return true, newRef, nil
	}

	return false, iName, nil
}

// TagIsUpgradeable returns true if the candidate tag is a newer version than current
func TagIsUpgradeable(current, candidate string) bool {
	if strings.EqualFold(current, "latest") {
		return false
	}

	currentSemVer, _ := semver.NewVersion(current)
	candidateSemVer, _ := semver.NewVersion(candidate)

	return candidateSemVer.Compare(currentSemVer) == 1 && candidateSemVer.Prerelease() == currentSemVer.Prerelease()
}

// GetCandidateTag finds the best candidate tag from discovered tags that matches the current tag format
func GetCandidateTag(discoveredTags []string, currentTag string) (string, bool) {
	var candidateTags []*semver.Version
	for _, tag := range discoveredTags {
		v, err := semver.NewVersion(tag)
		if err == nil {
			candidateTags = append(candidateTags, v)
		}
	}

	if len(candidateTags) > 0 {
		currentTagAttr := GetTagAttributes(currentTag)
		sort.Sort(sort.Reverse(semver.Collection(candidateTags)))

		for _, candidate := range candidateTags {
			candidateTagAttr := GetTagAttributes(candidate.Original())
			if currentTagAttr.AttributesMatch(candidateTagAttr) {
				return candidate.Original(), true
			}
		}
	}

	return "", false
}

// GetTagAttributes returns the attributes of a tag string
func GetTagAttributes(t string) TagAttributes {
	tagParts := strings.Split(t, "-")
	tagLevels := strings.Split(tagParts[0], ".")

	return TagAttributes{
		HasSuffix: len(tagParts) > 1,
		HasMajor:  len(tagLevels) >= 1 && tagLevels[0] != "",
		HasMinor:  len(tagLevels) >= 2,
		HasPatch:  len(tagLevels) == 3,
		Original:  t,
	}
}
