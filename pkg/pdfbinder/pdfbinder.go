package pdfbinder

import (
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/their-sober-press/alcobinder/pkg/alcobinder"
)

func BindPdf(inputFolder string, inputCSSFile string, outputPDFPath string) (err error) {
	htmlFilePath, err := makeHTMLFile(inputFolder, inputCSSFile)
	log.Println("HTML file written:", htmlFilePath)
	err = renderPDF(htmlFilePath, outputPDFPath)
	return
}

func makeHTMLFile(inputFolder string, inputCSSFile string) (htmlFilePath string, err error) {
	tempFile, err := ioutil.TempFile("", "bound-*.html")
	if err != nil {
		return
	}
	htmlFilePath = tempFile.Name()
	err = alcobinder.BindMarkdownsToFile(inputFolder, inputCSSFile, htmlFilePath)
	return
}

func renderPDF(pathToHTML string, outputPDFPath string) (err error) {
	cmd := exec.Command("weasyprint", "--presentational-hints", pathToHTML, outputPDFPath)
	err = cmd.Run()
	return
}
