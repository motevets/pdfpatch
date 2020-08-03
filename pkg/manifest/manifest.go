package manifest

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Manifest represents the manifest for a pdfpatch bundle
// Example:
//   sources:
//   - file_name: foo
//     url: http://example.com/foo.md
//     md5sum: a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a
//   styles:
//   - name: Regular
//     description: This is the regular formatting of the book.
//     style_sheet: regular.css
type Manifest struct {
	Sources []Source
	Styles  []Style
}

// Source represent a source file for patching
// FileName (required) an is the filename (without an path) for the pdf file to be patched
// Md5Sum (optional) is the check md5sum for the file
// URL (optional) is the URL from which the PDF can be obtained
// PatchedFiles (required) are the PDFs from which file names of the patches in order that the PDF text should patch to
type Source struct {
	URL          string
	FileName     string `yaml:"file_name"`
	Md5Sum       string
	PatchedFiles []string `yaml:"patched_files"`
}

// Style are a list of stylesheets that can be used to style the patched text
// Name (required) is the human readable name for the style
// Description (optional) is the human readable description for the style
// StyleSheet (required) is the file name (no path) for the style_sheet used for the style
type Style struct {
	Name        string
	Description string
	StyleSheet  string `yaml:"style_sheet"`
}

// ParseFile reads/parses a manifest from path
func ParseFile(path string) (theManifest Manifest, err error) {
	manifestData, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(manifestData, &theManifest)
	if err != nil {
		return
	}

	return
}

// SourceFileNames returns an array of the source file name in the manifest
func (m Manifest) SourceFileNames() (fileNames []string) {
	fileNames = make([]string, len(m.Sources))
	for i, source := range m.Sources {
		fileNames[i] = source.FileName
	}
	return
}
