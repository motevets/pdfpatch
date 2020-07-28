package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/motevets/pdfpatch/pkg/manifest"
	"github.com/motevets/pdfpatch/pkg/pdfbinder"
	"github.com/motevets/pdfpatch/pkg/pdfpatch"
)

const usage = `
pdfpatch SUBCOMMAND ARGS

  SUBCOMMAND: must be extract-text, make-patch, apply-patch, bind-pdf, or patch-pdfs
`

const extractTextUsage = `
pdfpatch extract-text MANIFEST_PATH PDF_DIR

  MANIFEST_PATH: file page to manifest file
  PDF_DIR:       path to director with source PDF files
`

const makePatchUsage = `
pdfpatch make-patch PDF_FILE MARKDOWN_FILE [ADDITIONAL_MARKDOWN_FILES ...]

  PDF_FILE:                  original source PDF file
  MARKDOWN_FILE:             file to diff against to make the patch
  ADDITIONAL_MARKDOWN_FILES: (optional) additional files appended to the first file with which to make the patch
`

const applyPatchUsage = `
pdfpatch apply-patch MANIFEST_PATH PDF_DIR PATCH_FILE

  MANIFEST_PATH: file page to manifest file
  PDF_DIR:       path to director with source PDF files
  PATCH_FILE:    path to the patch file (optional, default: /dev/stdin)
`

const bindPdfUsage = `
pdfpatch bind-pdf INPUT_MARKDOWNS_DIR INPUT_CSS_FILE OUTPUT_FILE_PATH

  INPUT_MARKDOWNS_DIR:    directory containing markdown file
  INPUT_CSS_FILE:         path to file used to style the book
  OUTPUT_FILE_PATH:       path where printable output HTML file is to be written
`

const patchPDFsUsage = `
pdfpatch patch-pdfs PATCH_BUNDLE_ZIP INPUT_PDF_DIR OUTPUT_PDF_PATH

  PATH_BUNDLE_ZIP:        bundle file containing assets needed to patch PDFs
  INPUT_PDF_DIR:          put the directory containing PDFs to patch
  OUTPUT_PDF_PATH:        path where output PDF should be written
`

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(2)
	}

	subcommand := os.Args[1]

	if subcommand == "extract-text" {
		if len(os.Args) != 4 {
			fmt.Println(extractTextUsage)
			os.Exit(2)
		}
		manifest := parseManifest(os.Args[2])
		fileNames, err := manifest.SourceFileNames()
		exitOnError(err, "Could not get file names from sources")
		text, err := extractor.TextFromPDFs(os.Args[3], fileNames)
		exitOnError(err, "Could not extract text")
		fmt.Println(text)
		os.Exit(0)
	} else if subcommand == "make-patch" {
		if len(os.Args) < 4 {
			fmt.Println(makePatchUsage)
			os.Exit(2)
		}
		patch, err := pdfpatch.GeneratePatch(os.Args[2], os.Args[3:])
		exitOnError(err, "Could not generate patch")
		fmt.Println(patch)
		os.Exit(0)
	} else if subcommand == "apply-patch" {
		var fileName string
		if len(os.Args) == 4 {
			fileName = "/dev/stdin"
		} else if len(os.Args) == 5 {
			fileName = os.Args[4]
		} else {
			fmt.Println(applyPatchUsage)
			os.Exit(2)
		}
		patchFileBytes, err := ioutil.ReadFile(fileName)
		exitOnError(err, "Could not read PATCH_FILE")
		manifest := parseManifest(os.Args[2])
		fileNames, err := manifest.SourceFileNames()
		exitOnError(err, "Could not get file names from sources")
		result, err := pdfpatch.ApplyPatch(os.Args[3], fileNames, string(patchFileBytes))
		exitOnError(err, "Could not generate patch")
		fmt.Println(result)
		os.Exit(0)
	} else if subcommand == "bind-pdf" {
		if len(os.Args) != 5 {
			fmt.Println(bindPdfUsage)
			os.Exit(2)
		}
		err := pdfbinder.BindPdf(os.Args[2], os.Args[3], os.Args[4])
		exitOnError(err, "Unable to bind PDF")
		os.Exit(0)
	} else {
		fmt.Println(usage)
		os.Exit(2)
	}

}

func parseManifest(manifestPath string) manifest.Manifest {
	theManifest, err := manifest.ParseFile(manifestPath)
	exitOnError(err, "Could not parse manifest")
	return theManifest
}

func exitOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		os.Exit(1)
	}
}
