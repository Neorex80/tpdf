package utils

import (
	"fmt"
	"os/exec"
)

func isUtilityInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// using pdftotext utility to convert pdf file to text.
func ConvertPDFToText(pdfFile, outputFile string) error {
	utility := "pdftotext"
	if !isUtilityInstalled(utility) {
		fmt.Printf("'%s' is NOT installed, Please install it to use this tool.\n", utility)
	}
	cmd := exec.Command("pdftotext", "-layout", pdfFile, outputFile)
	return cmd.Run()

}
