package dockerfile

import (
	"strings"
)

// ImageRef represents a container image reference found in a Dockerfile
type ImageRef struct {
	Image string
	Tag   string
}

// Ref returns the full image reference as image:tag
func (i ImageRef) Ref() string {
	return i.Image + ":" + i.Tag
}

// FindImages parses Dockerfile content and returns image references from FROM lines.
// Images with variable references in tags (containing $) are skipped.
// Images without an explicit tag are skipped.
func FindImages(content string) []ImageRef {
	var images []ImageRef
	seen := map[string]bool{}

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToUpper(line), "FROM ") {
			continue
		}

		rest := line[len("FROM "):]
		rest = strings.TrimSpace(rest)

		// Skip --platform or other flags
		if strings.HasPrefix(rest, "--") {
			idx := strings.Index(rest, " ")
			if idx < 0 {
				continue
			}
			rest = strings.TrimSpace(rest[idx+1:])
		}

		// Take the image reference (first field, before AS or end of line)
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			continue
		}
		ref := fields[0]

		image, tag, ok := splitRef(ref)
		if !ok {
			continue
		}

		if strings.Contains(tag, "$") {
			continue
		}

		if seen[ref] {
			continue
		}
		seen[ref] = true

		images = append(images, ImageRef{
			Image: image,
			Tag:   tag,
		})
	}

	return images
}

// splitRef splits an image reference into image name and tag,
// handling registry:port/image:tag format correctly.
func splitRef(ref string) (string, string, bool) {
	slashIdx := strings.LastIndex(ref, "/")
	searchFrom := 0
	if slashIdx >= 0 {
		searchFrom = slashIdx
	}
	colonIdx := strings.Index(ref[searchFrom:], ":")
	if colonIdx < 0 {
		return "", "", false
	}
	colonIdx += searchFrom
	return ref[:colonIdx], ref[colonIdx+1:], true
}

// ReplaceImage replaces an old image:tag reference with a new one
// in Dockerfile content, only within FROM lines.
func ReplaceImage(content, oldRef, newRef string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(trimmed), "FROM ") {
			lines[i] = strings.Replace(line, oldRef, newRef, 1)
		}
	}
	return strings.Join(lines, "\n")
}
