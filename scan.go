package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

func fileExists(filenames ...string) scanFn {
	return func(dir string) bool {
		for _, filename := range filenames {
			info, err := os.Stat(filepath.Join(dir, filename))
			if err != nil {
				continue
			}
			if !info.IsDir() {
				return true
			}
		}
		return false
	}
}

func fileContains(path string, pattern string) bool {
	file, err := os.Open(path)

	if err != nil {
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		re := regexp.MustCompile(pattern)
		if re.MatchString(scanner.Text()) {
			return true
		}
	}

	return false
}

func dirContains(glob string, patterns ...string) scanFn {
	return func(dir string) bool {
		for _, pattern := range patterns {
			filenames, _ := filepath.Glob(filepath.Join(dir, glob))
			for _, filename := range filenames {
				if fileContains(filename, pattern) {
					return true
				}
			}
		}
		return false
	}
}

type scanFn func(dir string) bool

func sourceTriggers(sourceDir string, scanners ...scanFn) bool {
	for _, check := range scanners {
		if check(sourceDir) {
			return true
		}
	}
	return false
}
