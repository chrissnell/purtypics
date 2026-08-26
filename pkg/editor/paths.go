package editor

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
)

// errPathEscape is returned when a user-supplied path resolves outside the root.
var errPathEscape = errors.New("path escapes root")

// resolveWithinRoot joins subPath onto root and verifies the cleaned result
// stays inside root. It rejects any path containing ".." segments — cheap
// early rejection — and re-checks via filepath.Rel to catch symlink-adjacent
// tricks and absolute inputs.
func resolveWithinRoot(root, subPath string) (string, error) {
	// Reject any ".." segment before joining. This is the fast path; the
	// filepath.Rel check below catches absolute inputs and other tricks.
	if slices.Contains(strings.Split(filepath.ToSlash(subPath), "/"), "..") {
		return "", errPathEscape
	}

	joined := filepath.Join(root, subPath)
	rel, err := filepath.Rel(root, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errPathEscape
	}
	return joined, nil
}
