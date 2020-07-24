package extractor

import (
	"bytes"
	"path"

	"github.com/ledongthuc/pdf"
)

func TextFromPDFs(directory string, files []string) (extractedText string, err error) {
	for _, file := range files {
		pdfPath := path.Join(directory, file)
		fileText, err := textFromPDF(pdfPath)
		if err != nil {
			return "", err
		}
		extractedText += fileText + "\n"
	}
	return
}

func textFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
