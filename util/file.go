package util

import (
	"bytes"
	"os"
	"path/filepath"
)

func FileCreate(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func FileDelete(path string) error {
	return os.Remove(path)
}

func FileRead(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func FileIsEditable(filename string) bool {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

type FileExpanderTags struct {
	StartTag   []byte
	EndTag     []byte
	ExpandWith []byte
	Count      int
}

func FileExpand(filePath string, tags []FileExpanderTags) error {
	content, err := FileRead(filePath)
	if err != nil {
		return err
	}

	lines := bytes.Split(content, []byte("\n"))

	for _, tag := range tags {
		var newLines [][]byte
		var appendOriginalLine bool = true
		replacementCount := 0
		insideBlock := false

		for _, line := range lines {
			if bytes.Contains(line, tag.StartTag) {
				insideBlock = true
				if tag.Count == 0 || replacementCount < tag.Count {
					newLines = append(newLines, line)
					newLines = append(newLines, tag.ExpandWith)
					replacementCount++
					appendOriginalLine = false
				} else {
					appendOriginalLine = true
				}
			}

			if bytes.Contains(line, tag.EndTag) {
				insideBlock = false
				appendOriginalLine = true
			}

			if appendOriginalLine || !insideBlock {
				newLines = append(newLines, line)
			}
		}
		lines = newLines
	}

	newContent := bytes.Join(lines, []byte("\n"))
	return FileCreate(filePath, newContent)
}
