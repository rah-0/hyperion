package util

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// DirectoryGetHighestVersion finds the directory with the highest numeration (e.g., v1, v2, v100, etc.)
func DirectoryGetHighestVersion(path string) (string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`^v(\d+)$`)
	var versions []int
	dirMap := make(map[int]string)

	for _, entry := range entries {
		if entry.IsDir() {
			matches := re.FindStringSubmatch(entry.Name())
			if matches != nil {
				num, err := strconv.Atoi(matches[1])
				if err == nil {
					versions = append(versions, num)
					dirMap[num] = entry.Name()
				}
			}
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versioned directories found")
	}

	sort.Sort(sort.Reverse(sort.IntSlice(versions))) // Sort descending
	return dirMap[versions[0]], nil
}

// DirectoryIsEmpty checks if a directory is empty
func DirectoryIsEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// Read at most one entry from the directory
	_, err = dir.Readdirnames(1)
	if err == nil {
		return false, nil // Directory is not empty
	} else if errors.Is(err, os.ErrNotExist) {
		return true, nil // Directory is empty
	}
	return true, nil // If EOF, directory is empty
}

// DirectoryCreate creates a new directory with the given path
func DirectoryCreate(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// DirectoryRemove deletes the directory and its contents
func DirectoryRemove(path string) error {
	return os.RemoveAll(path)
}

// DirectoriesInPath returns a list of directories in the given path (first-level only).
func DirectoriesInPath(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}

	return directories, nil
}
