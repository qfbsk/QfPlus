package app

import (
	"os"
	"path/filepath"
)

// scanTreeSize returns the number of items and total bytes under root without
// following symlinks or junctions. A link/junction counts as a single item so
// sdks/<name> never double-counts the runtime root it points at — the copy step
// also treats links as shallow, so the preview must agree with the real copy.
func scanTreeSize(root string) (int, int64, error) {
	if _, lerr := os.Readlink(root); lerr == nil {
		return 1, 0, nil
	}
	info, err := os.Lstat(root)
	if err != nil {
		return 0, 0, err
	}
	if !info.IsDir() {
		return 1, info.Size(), nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0, err
	}
	total := 1
	var size int64
	for _, entry := range entries {
		count, childSize, err := scanTreeSize(filepath.Join(root, entry.Name()))
		if err != nil {
			return 0, 0, err
		}
		total += count
		size += childSize
	}
	return total, size, nil
}
