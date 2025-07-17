package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println(`Please provide a pdf file.
			Ex= tpdf example.pdf`)
		os.Exit(1)
	}
	pdf := os.Args[2]
	tpdf := "./output.txt"

	err := utils.ConvertPdfToText(pdf, tpdf)

	if err != nil {
		fmt.Println("[!] Error converting PDF:", err)
		return
	}

	pages := utils.SplitIntoPages(tpdf)
	utils.RenderPage(pages, 0)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nCommand (gotoPage N / gotoChapter Title / exit): ")

		if !reader.Scan() {
			break
		}

		input := reader.Text()

		if strings.HasPrefix(input, "gotoPage") {
			var pageNum int
			fmt.Sscanf(input, "gotoPage %d", &pageNum)
			utils.RenderPage(pages, pageNum-1)
		} else if strings.HasPrefix(input, "gotoChapter") {
			title := strings.TrimPrefix(input, "gotoChapter ")
			index := utils.FindChapterPage(pages, title)
			if index >= 0 {
				utils.RenderPage(pages, index)
			} else {
				fmt.Println("Chapter not found.")
			}
		} else if input == "exit" {
			break
		}
	}

}
