package extractor

import (
	"path"

	"code.sajari.com/docconv"
)

func TextFromPDFs(directory string, files []string) (extractedText string, err error) {
	for _, file := range files {
		pdfPath := path.Join(directory, file)
		fileText, err := TextFromPDF(pdfPath)
		if err != nil {
			return "", err
		}
		extractedText += fileText + "\n"
	}
	return
}

func TextFromPDF(path string) (string, error) {
	res, err := docconv.ConvertPath(path)
	if err != nil {
		return "", err
	}
	return res.Body, err
}
