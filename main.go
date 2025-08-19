package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

// getPDFToTextCommand returns the appropriate pdftotext command for the platform
func getPDFToTextCommand() (string, error) {
	// Try common command names
	commands := []string{"pdftotext"}
	
	// Add Windows-specific paths
	if runtime.GOOS == "windows" {
		commands = append([]string{
			"pdftotext.exe",
			"C:\\poppler\\bin\\pdftotext.exe",
			"C:\\Program Files\\poppler\\bin\\pdftotext.exe",
		}, commands...)
	}
	
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err == nil {
			return cmd, nil
		}
	}
	
	// Return installation instructions based on OS
	var installMsg string
	switch runtime.GOOS {
	case "windows":
		installMsg = "Download Poppler from: https://github.com/oschwartz10612/poppler-windows/releases/"
	case "darwin":
		installMsg = "Install with: brew install poppler"
	default:
		installMsg = "Install with: sudo apt install poppler-utils (Ubuntu/Debian) or sudo pacman -S poppler (Arch)"
	}
	
	return "", fmt.Errorf("pdftotext not found. %s", installMsg)
}

// ConvertPDFToText converts a PDF file to text using pdftotext
func (v *PDFViewer) ConvertPDFToText() error {
	cmd, err := getPDFToTextCommand()
	if err != nil {
		return err
	}
	
	// Check if PDF file exists
	if _, err := os.Stat(v.PDFPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF file not found: %s", v.PDFPath)
	}
	
	execCmd := exec.Command(cmd, "-layout", v.PDFPath, v.TextPath)
	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("failed to convert PDF: %v", err)
	}
	
	return nil
}

// LoadPages loads and splits the text file into pages
func (v *PDFViewer) LoadPages() error {
	content, err := os.ReadFile(v.TextPath)
	if err != nil {
		return fmt.Errorf("failed to read converted text: %v", err)
	}
	
	if len(content) == 0 {
		return fmt.Errorf("PDF appears to be empty or conversion failed")
	}

	lines := strings.Split(string(content), "\n")
	const maxLinesPerPage = 40
	
	var pages []Page
	var currentPageLines []string

	for i, line := range lines {
		currentPageLines = append(currentPageLines, line)
		
		if (i+1)%maxLinesPerPage == 0 {
			pages = append(pages, Page{Lines: currentPageLines})
			currentPageLines = []string{}
		}
	}
	
	if len(currentPageLines) > 0 {
		pages = append(pages, Page{Lines: currentPageLines})
	}
	
	if len(pages) == 0 {
		return fmt.Errorf("no content found in PDF")
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
Commands:
  n, next         Next page
  p, prev         Previous page  
  g N, gotoPage N Jump to page N
  gc T            Jump to chapter T
  q, exit         Quit
  h, help         Show this help
`)
}

// ProcessCommand processes user commands
func (v *PDFViewer) ProcessCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	
	switch parts[0] {
	case "gotoPage", "g":
		if len(parts) < 2 {
			fmt.Println("Usage: gotoPage N")
			return
		}
		var pageNum int
		if _, err := fmt.Sscanf(parts[1], "%d", &pageNum); err != nil {
			fmt.Println("Invalid page number")
			return
		}
		v.GotoPage(pageNum)

	case "gotoChapter", "gc":
		if len(parts) < 2 {
			fmt.Println("Usage: gotoChapter Title")
			return
		}
		title := strings.Join(parts[1:], " ")
		v.GotoChapter(title)

	case "next", "n":
		v.NextPage()

	case "previous", "prev", "p":
		v.PreviousPage()

	case "exit", "quit", "q":
		v.cleanup()
		os.Exit(0)

	case "help", "h":
		v.ShowHelp()

	default:
		fmt.Printf("Unknown command: %s. Type 'help' for commands.\n", parts[0])
	}
}

// cleanup removes temporary files
func (v *PDFViewer) cleanup() {
	os.Remove(v.TextPath)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tpdf <file.pdf>")
		os.Exit(1)
	}

	viewer := &PDFViewer{
		PDFPath:  os.Args[1],
		TextPath: ".tpdf_temp.txt",
	}
	defer viewer.cleanup()

	fmt.Println("Converting PDF...")
	if err := viewer.ConvertPDFToText(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := viewer.LoadPages(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	viewer.RenderPage()

	// Command loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nCommand (n/p/g N/gc Title/q/h): ")
		if !scanner.Scan() {
			break
		}
		if input := strings.TrimSpace(scanner.Text()); input != "" {
			viewer.ProcessCommand(input)
		}
	}
}
