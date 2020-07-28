package pdfpatch

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/motevets/pdfpatch/pkg/pdfbinder"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func GeneratePatch(inputPDFFile string, markdownFiles []string) (patch string, err error) {
	extractedText, err := extractor.TextFromPDF(inputPDFFile)
	if err != nil {
		return
	}
	markdownFilesText, err := concatFilesToString(markdownFiles)
	if err != nil {
		return
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(extractedText, markdownFilesText, false)
	patches := dmp.PatchMake(diffs)
	return dmp.PatchToText(patches), nil
}

func concatFilesToString(files []string) (output string, err error) {
	for _, file := range files {
		var fileText []byte
		fileText, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}
		output += "\n" + string(fileText) + "\n"
	}
	return
}

func ApplyPatch(inputFilesDir string, inputFiles []string, patch string) (newText string, err error) {
	extractedText, err := extractor.TextFromPDFs(inputFilesDir, inputFiles)
	if err != nil {
		return
	}
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(patch)
	if err != nil {
		return
	}
	newText, _ = dmp.PatchApply(patches, extractedText)
	return
}

func PatchPDF(inputPDFsDir string, inputPDFs []string, patchFile string, cssFile string, outputPDFPath string) (err error) {
	patch, err := ioutil.ReadFile(patchFile)
	patchedMarkdownDir, err := ioutil.TempDir("", "patched_markdowns")
	patchedMarkdownFilePath := path.Join(patchedMarkdownDir, "patched.md")
	if err != nil {
		return
	}
	patchedText, err := ApplyPatch(inputPDFsDir, inputPDFs, string(patch))
	if err != nil {
		return
	}
	err = ioutil.WriteFile(patchedMarkdownFilePath, []byte(patchedText), 0755)
	if err != nil {
		return
	}
	log.Println("patched mardown written:", patchedMarkdownFilePath)
	err = pdfbinder.BindPdf(patchedMarkdownDir, cssFile, outputPDFPath)
	return
}
