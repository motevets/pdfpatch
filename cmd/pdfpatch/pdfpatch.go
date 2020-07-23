package main

import (
	"fmt"
	"os"

	"github.com/motevets/pdfpatch/pkg/extractor"
	"github.com/motevets/pdfpatch/pkg/manifest"
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
