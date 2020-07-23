package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/motevets/pdfpatch/pkg/manifest"
	"github.com/motevets/pdfpatch/pkg/pdfpatch"
)

const usage = `
pdfpatch SUBCOMMAND ARGS

  SUBCOMMAND: must be create-patch or extract-text
`

const extractTextUsage = `
pdfpatch extract-text MANIFEST_PATH PDF_DIR

  MANIFEST_PATH: file page to manifest file
  PDF_DIR:       path to director with source PDF files
`

const makePatchUsage = `
pdfpatch make-patch MANIFEST_PATH PDF_DIR TO_FILE

  MANIFEST_PATH: file page to manifest file
  PDF_DIR:       path to director with source PDF files
  TO_FILE:       the file to diff against (optional, default: /dev/stdin)
`

const applyPatchUsage = `
pdfpatch apply-patch MANIFEST_PATH PDF_DIR PATCH_FILE

  MANIFEST_PATH: file page to manifest file
  PDF_DIR:       path to director with source PDF files
  PATCH_FILE:    path to the patch file (optional, default: /dev/stdin)
`

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(2)
	}

	subcommand := os.Args[1]

	if subcommand == "create-patch" {
		fmt.Println("Yippee")
	} else if subcommand == "extract-text" {
		if len(os.Args) != 4 {
			fmt.Println(extractTextUsage)
			os.Exit(2)
		}
		manifest := parseManifest(os.Args[2])
		fileNames, err := manifest.SourceFileNames()
		exitOnError(err, "Could not get file names from sources")
		text, err := extractor.TextFromPdf(os.Args[3], fileNames)
		exitOnError(err, "Could not extract text")
		fmt.Println(text)
		os.Exit(0)
	} else if subcommand == "make-patch" {
		var fileName string
		if len(os.Args) == 4 {
			fileName = "/dev/stdin"
		} else if len(os.Args) == 5 {
			fileName = os.Args[4]
		} else {
			fmt.Println(makePatchUsage)
			os.Exit(2)
		}
		toFileBytes, err := ioutil.ReadFile(fileName)
		exitOnError(err, "Could not read TO_FILE")
		manifest := parseManifest(os.Args[2])
		fileNames, err := manifest.SourceFileNames()
		exitOnError(err, "Could not get file names from sources")
		patch, err := pdfpatch.GeneratePatch(os.Args[3], fileNames, string(toFileBytes))
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
