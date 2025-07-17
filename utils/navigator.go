package utils

import "strings"

func FindChapterPage(pages [][]string, title string) int {
	title = strings.ToLower(title)
	for i, page := range pages {
		for _, line := range page {
			if strings.Contains(strings.ToLower(line), title) {
				return i
			}
		}
	}
	return -1
}

func clearScreen() {
	print("\033[h\033[2j")
}
