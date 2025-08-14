package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Page represents a single page of text
type Page struct {
	Lines []string
}

// PDFViewer holds the state of our PDF viewer
type PDFViewer struct {
	Pages       []Page
	CurrentPage int
	PDFPath     string
	TextPath    string
}

// isUtilityInstalled checks if a command-line utility is available
func isUtilityInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// ConvertPDFToText converts a PDF file to text using pdftotext
func (v *PDFViewer) ConvertPDFToText() error {
	utility := "pdftotext"
	if !isUtilityInstalled(utility) {
		return fmt.Errorf("'%s' is not installed. Please install poppler-utils to use this tool", utility)
	}
	
	cmd := exec.Command("pdftotext", "-layout", v.PDFPath, v.TextPath)
	return cmd.Run()
}

// LoadPages loads and splits the text file into pages
func (v *PDFViewer) LoadPages() error {
	content, err := os.ReadFile(v.TextPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	const maxLinesPerPage = 40
	
	var pages []Page
	var currentPageLines []string
	lineCount := 0

	for _, line := range lines {
		currentPageLines = append(currentPageLines, line)
		lineCount++
		
		if lineCount >= maxLinesPerPage {
			pages = append(pages, Page{Lines: currentPageLines})
			currentPageLines = []string{}
			lineCount = 0
		}
	}
	
	if len(currentPageLines) > 0 {
		pages = append(pages, Page{Lines: currentPageLines})
	}
	
	v.Pages = pages
	return nil
}

// RenderPage displays the current page
func (v *PDFViewer) RenderPage() {
	if v.CurrentPage < 0 || v.CurrentPage >= len(v.Pages) {
		fmt.Println("Invalid page number")
		return
	}

	// Clear screen
	print("\033[H\033[2J")

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(" Page %d of %d \n", v.CurrentPage+1, len(v.Pages))
	fmt.Println(strings.Repeat("-", 80))

	for _, line := range v.Pages[v.CurrentPage].Lines {
		fmt.Println(line)
	}
	fmt.Println(strings.Repeat("-", 80))
}

// FindChapterPage finds a page containing the specified chapter title
func (v *PDFViewer) FindChapterPage(title string) int {
	title = strings.ToLower(title)
	for i, page := range v.Pages {
		for _, line := range page.Lines {
			if strings.Contains(strings.ToLower(line), title) {
				return i
			}
		}
	}
	return -1
}

// NextPage moves to the next page
func (v *PDFViewer) NextPage() {
	if v.CurrentPage < len(v.Pages)-1 {
		v.CurrentPage++
		v.RenderPage()
	} else {
		fmt.Println("You are at the last page.")
	}
}

// PreviousPage moves to the previous page
func (v *PDFViewer) PreviousPage() {
	if v.CurrentPage > 0 {
		v.CurrentPage--
		v.RenderPage()
	} else {
		fmt.Println("You are at the first page.")
	}
}

// GotoPage jumps to a specific page
func (v *PDFViewer) GotoPage(pageNum int) {
	if pageNum >= 1 && pageNum <= len(v.Pages) {
		v.CurrentPage = pageNum - 1
		v.RenderPage()
	} else {
		fmt.Println("Invalid page number.")
	}
}

// GotoChapter jumps to a page containing the specified chapter
func (v *PDFViewer) GotoChapter(title string) {
	index := v.FindChapterPage(title)
	if index >= 0 {
		v.CurrentPage = index
		v.RenderPage()
	} else {
		fmt.Println("Chapter not found.")
	}
}

// ShowHelp displays the help screen
func (v *PDFViewer) ShowHelp() {
	fmt.Println(`
tpdf - Terminal PDF Viewer

Usage:
  tpdf <file.pdf>

Commands:
  next            Go to next page
  previous        Go to previous page
  gotoPage N      Jump to page N
  gotoChapter T   Jump to chapter titled T
  exit            Exit viewer

Press 'q' to quit.
`)
}

// ProcessCommand processes user commands
func (v *PDFViewer) ProcessCommand(input string) {
	switch {
	case strings.HasPrefix(input, "gotoPage"):
		var pageNum int
		_, err := fmt.Sscanf(input, "gotoPage %d", &pageNum)
		if err != nil {
			fmt.Println("Invalid command format. Use: gotoPage N")
			return
		}
		v.GotoPage(pageNum)

	case strings.HasPrefix(input, "gotoChapter"):
		title := strings.TrimPrefix(input, "gotoChapter ")
		if title == "" {
			fmt.Println("Invalid command format. Use: gotoChapter Title")
			return
		}
		v.GotoChapter(title)

	case input == "next":
		v.NextPage()

	case input == "previous":
		v.PreviousPage()

	case input == "exit":
		os.Remove(v.TextPath) // Clean up temporary file
		os.Exit(0)

	case input == "q":
		os.Remove(v.TextPath) // Clean up temporary file
		os.Exit(0)

	default:
		fmt.Println("Unknown command. Type 'exit' to quit.")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tpdf <file.pdf>")
		fmt.Println("Use 'tpdf' without arguments to show help.")
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	textPath := "./.tpdf_temp.txt" // Hidden temporary file

	viewer := &PDFViewer{
		PDFPath:  pdfPath,
		TextPath: textPath,
	}

	// Convert PDF to text
	fmt.Println("Loading PDF...")
	if err := viewer.ConvertPDFToText(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Load pages
	if err := viewer.LoadPages(); err != nil {
		fmt.Printf("Error loading pages: %v\n", err)
		os.Exit(1)
	}

	// Clean up temporary file when exiting
	defer os.Remove(textPath)

	// Display first page
	viewer.RenderPage()

	// Command loop
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nCommand (next/previous/gotoPage/gotoChapter/exit): ")
		if !reader.Scan() {
			break
		}
		input := strings.TrimSpace(reader.Text())
		if input != "" {
			viewer.ProcessCommand(input)
		}
	}
}
