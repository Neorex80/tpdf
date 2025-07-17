package utils

import (
	"fmt"
	"os"
	"strings"
)

func SplitIntoPages(textFile string) [][]string {
	content, err := os.ReadFile(textFile)
	if err != nil {
		fmt.Println("[!]Error reading file:", err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	var pages [][]string
	var page []string
	lineCount := 0
	const maxLinePerPage = 40

	for _, line := range lines {
		page = append(page, line)
		lineCount++
		if lineCount >= maxLinePerPage {
			pages = append(pages, page)
			page = []string{}
			lineCount = 0
		}
	}
	if len(page) > 0 {
		pages = append(pages, page)
	}
	return pages
}

func RenderPage(pages [][]string, index int) {
	if index < 0 || index >= len(pages) {
		fmt.Println("Invalid page number")
		return
	}

	clearScreen()

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(" Page %d \n", index+1)
	fmt.Println(strings.Repeat("-", 80))

	for _, line := range pages[index] {
		fmt.Println(line)
	}
	fmt.Println(strings.Repeat("-", 80))
}
