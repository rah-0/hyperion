package util

import (
	"bytes"
	"fmt"
	"io"
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

func FileSizeHuman(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	size := info.Size()
	units := []string{"B", "KB", "MB", "GB", "TB"}
	var i int
	floatSize := float64(size)

	for floatSize >= 1024 && i < len(units)-1 {
		floatSize /= 1024
		i++
	}

	return fmt.Sprintf("%.2f %s", floatSize, units[i]), nil
}

func FileCOpy(pathSource, pathDestiny string) error {
	// Open the source file
	sourceFile, err := os.Open(pathSource)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(pathDestiny)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure all data is written to disk
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
