package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prem0x01/tpdf/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Please provide a pdf file.
			Ex= tpdf example.pdf`)
		os.Exit(1)
	}
	pdf := os.Args[1]
	tpdf := "./output.txt"

	err := utils.ConvertPDFToText(pdf, tpdf)

	if err != nil {
		fmt.Println("[!] Error converting PDF:", err)
		return
	}

	pages := utils.SplitIntoPages(tpdf)
	currentPage := 0

	utils.RenderPage(pages, 0)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nCommand (gotoPage N / gotoChapter Title / next / previous / exit): ")
		if !reader.Scan() {
			break
		}
		input := reader.Text()

		switch {
		case strings.HasPrefix(input, "gotoPage"):
			var pageNum int
			fmt.Sscanf(input, "gotoPage %d", &pageNum)
			if pageNum >= 1 && pageNum <= len(pages) {
				currentPage = pageNum - 1
				utils.RenderPage(pages, currentPage)
			} else {
				fmt.Println("Invalid page number.")
			}

		case strings.HasPrefix(input, "gotoChapter"):
			title := strings.TrimPrefix(input, "gotoChapter ")
			index := utils.FindChapterPage(pages, title)
			if index >= 0 {
				currentPage = index
				utils.RenderPage(pages, currentPage)
			} else {
				fmt.Println("Chapter not found.")
			}

		case input == "next":
			if currentPage < len(pages)-1 {
				currentPage++
				utils.RenderPage(pages, currentPage)
			} else {
				fmt.Println("You are at the last page.")
			}

		case input == "previous":
			if currentPage > 0 {
				currentPage--
				utils.RenderPage(pages, currentPage)
			} else {
				fmt.Println("You are at the first page.")
			}

		case input == "exit":
			return

		default:
			fmt.Println("Unknown command.")
		}
	}

}
