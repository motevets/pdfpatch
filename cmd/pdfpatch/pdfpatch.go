package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/motevets/pdfpatch/pkg/manifest"
	"github.com/motevets/pdfpatch/pkg/pdfbinder"
	"github.com/motevets/pdfpatch/pkg/pdfpatch"
)

const usage = `
pdfpatch SUBCOMMAND ARGS

  SUBCOMMAND: must be extract-text, make-patch, make-patches, apply-patch, bind-pdf, or patch-pdfs
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

const makePatchesUsage = `
pdfpatch make-patches MANIFEST_PATH PDF_DIR MARKDOWN_DIR OUTPUT_DIR

  MANIFEST_PATH: file page to manifest file
  PDF_DIR:       path to directory with source PDF files
  MARKDOWN_FILE: path to directory with files to diff against to make the patch
  OUTPUT_DIR:    path where patches should be written
`

const applyPatchUsage = `
pdfpatch apply-patch PDF_FILE [PATCH_FILE]

  PDF_FILE:      path to source PDF file with which to patch
  PATCH_FILE:    path to the patch file (optional, default: /dev/stdin)
`

const bindPdfUsage = `
pdfpatch bind-pdf INPUT_MARKDOWNS_DIR INPUT_CSS_FILE OUTPUT_FILE_PATH

  INPUT_MARKDOWNS_DIR:    directory containing markdown file
  INPUT_CSS_FILE:         path to file used to style the book
  OUTPUT_FILE_PATH:       path where printable output HTML file is to be written
`

const patchPDFsUsage = `
pdfpatch patch-pdfs MANIFEST_PATH INPUT_PDF_DIR PATCHES_DIR CSS_PATH OUTPUT_PDF_PATH

  MANIFEST_PATH:   file page to manifest file
  INPUT_PDF_DIR:   put the directory containing PDFs to patch
  PATCHES_DIR:	   directory containing patches with filenames like "input_pdf_file.pdf.patch" for each PDF file
  CSS_PATH:        path to the CSS file used to style the output PDF
  OUTPUT_PDF_PATH: path where output PDF should be written
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
		fileNames := manifest.SourceFileNames()
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
	} else if subcommand == "make-patches" {
		if len(os.Args) != 6 {
			fmt.Println(makePatchesUsage)
			os.Exit(2)
		}
		manifest := parseManifest(os.Args[2])
		pdfMarkdowns := make([]pdfpatch.PDFMarkdowns, len(manifest.Sources))
		for i, source := range manifest.Sources {
			pdfFilename := source.FileName
			pdfMarkdowns[i] = pdfpatch.PDFMarkdowns{
				PDFFileName:       pdfFilename,
				MarkdownFileNames: source.PatchedFiles,
			}
		}
		patches, err := pdfpatch.GeneratePatches(pdfMarkdowns, os.Args[3], os.Args[4])
		exitOnError(err, "Could not generate patches")
		for _, patch := range patches {
			outputPath := path.Join(os.Args[5], patch.PDFFileName+".patch")
			err := ioutil.WriteFile(outputPath, []byte(patch.Patch), 0755)
			exitOnError(err, "Could not write patch file")
		}
		os.Exit(0)
	} else if subcommand == "apply-patch" {
		var patchFileName string
		if len(os.Args) == 3 {
			patchFileName = "/dev/stdin"
		} else if len(os.Args) == 4 {
			patchFileName = os.Args[3]
		} else {
			fmt.Println(applyPatchUsage)
			os.Exit(2)
		}
		result, err := pdfpatch.ApplyPatch(os.Args[2], patchFileName)
		exitOnError(err, "Could not apply patch")
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
	} else if subcommand == "patch-pdfs" {
		if len(os.Args) != 7 {
			fmt.Println(patchPDFsUsage)
			os.Exit(2)
		}
		manifest := parseManifest(os.Args[2])
		pdfFiles := manifest.SourceFileNames()
		err := pdfpatch.PatchPDF(pdfFiles, os.Args[3], os.Args[4], os.Args[5], os.Args[6])
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
