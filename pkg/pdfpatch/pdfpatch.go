package pdfpatch

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/motevets/pdfpatch/pkg/pdfbinder"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type PDFMarkdowns struct {
	PDFFileName       string
	MarkdownFileNames []string
}

type PDFPatch struct {
	PDFFileName string
	Patch       string
}

func GeneratePatch(inputPDFFile string, markdownFiles []string) (patch string, err error) {
	if len(markdownFiles) == 0 {
		log.Println("WARNING: empty list of markdown files to diff against", inputPDFFile)
	}
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

func GeneratePatches(pdfMarkdownsList []PDFMarkdowns, pdfsDir string, markdownsDir string) (patches []PDFPatch, err error) {
	patches = make([]PDFPatch, len(pdfMarkdownsList))
	for i, pdfMarkdowns := range pdfMarkdownsList {
		var patch string
		patches[i] = PDFPatch{PDFFileName: pdfMarkdowns.PDFFileName}
		pdfFile := path.Join(pdfsDir, pdfMarkdowns.PDFFileName)
		markdownFiles := make([]string, len(pdfMarkdowns.MarkdownFileNames))
		for j, markdownFileName := range pdfMarkdowns.MarkdownFileNames {
			markdownFiles[j] = path.Join(markdownsDir, markdownFileName)
		}
		patch, err = GeneratePatch(pdfFile, markdownFiles)
		if err != nil {
			return
		}
		patches[i].Patch = patch
	}
	return
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

func ApplyPatch(inputPDFFilePath string, patchFilePath string) (newText string, err error) {
	extractedText, err := extractor.TextFromPDF(inputPDFFilePath)
	if err != nil {
		return
	}
	patch, err := ioutil.ReadFile(patchFilePath)
	if err != nil {
		return
	}

	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(string(patch))
	if err != nil {
		return
	}
	newText, _ = dmp.PatchApply(patches, extractedText)
	return
}

func PatchPDF(inputPDFs []string, inputPDFsDir string, patchFilesDir string, cssFile string, outputPDFPath string) (err error) {
	patchedMarkdownDir, err := ioutil.TempDir("", "patched_markdowns")
	for i, pdfFileName := range inputPDFs {
		var patchedText string
		pdfFilePath := path.Join(inputPDFsDir, pdfFileName)
		patchFilePath := path.Join(patchFilesDir, pdfFileName+".patch")
		patchedMarkdownFileName := fmt.Sprintf("%04d_%s.md", i, pdfFileName)
		patchedMarkdownPath := path.Join(patchedMarkdownDir, patchedMarkdownFileName)

		patchedText, err = ApplyPatch(pdfFilePath, patchFilePath)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(patchedMarkdownPath, []byte(patchedText), 0644)
		if err != nil {
			return
		}
	}
	log.Println("patched mardowns written:", patchedMarkdownDir)
	err = pdfbinder.BindPdf(patchedMarkdownDir, cssFile, outputPDFPath)
	return
}
