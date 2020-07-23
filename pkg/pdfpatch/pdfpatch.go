package pdfpatch

import (
	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func GeneratePatch(inputFilesDir string, inputFiles []string, diffWith string) (patch string, err error) {
	extractedText, err := extractor.TextFromPdf(inputFilesDir, inputFiles)
	if err != nil {
		return
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(extractedText, diffWith, false)
	patches := dmp.PatchMake(diffs)
	return dmp.PatchToText(patches), nil
}

func ApplyPatch(inputFilesDir string, inputFiles []string, patch string) (newText string, err error) {
	extractedText, err := extractor.TextFromPdf(inputFilesDir, inputFiles)
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
