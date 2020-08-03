package manifest

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/mholt/archiver"
)

// Bundle represents the contents of a packaged (compressed) bundle
type Bundle struct {
	Manifest     Manifest
	ManifestPath string
	CSSDir       string
	PatchesDir   string
}

// UnpackBundle unzips a bundle and returns its content
//
// Bundle file can be a tar, gzip(tar), or zip
//
// Bundle must have the following directory structure:
//
//   bundle.zip
//   ├── css
//   │   ├── CSS_FILE_1.css
//   │   ├── OPTIONAL_CSS_FILE_2.css
//   │   └── ...
//   ├── manifest.yml
//   └── patches
//       ├── PATCH_FILE_1.css
//       ├── OPTIONAL_PATCH_FILE_2.css
//       └── ...
func UnpackBundle(bundleFilePath string) (bundle Bundle, err error) {
	var (
		tempDir     string
		theManifest Manifest
	)

	tempDir, err = ioutil.TempDir("", "manifest_bundle")
	err = archiver.Unarchive(bundleFilePath, tempDir)
	if err != nil {
		return
	}
	bundle.ManifestPath = path.Join(tempDir, "manifest.yml")
	theManifest, err = ParseFile(bundle.ManifestPath)
	bundle.Manifest = theManifest
	bundle.CSSDir = path.Join(tempDir, "css")
	bundle.PatchesDir = path.Join(tempDir, "patches")
	return
}

// CSSFilePath returns the path to the style sheet of an extracted bundle
// note: currently there is no validation that the stylesheet exists and err will always be nil
func (bundle Bundle) CSSFilePath(styleSheet string) (styleSheetPath string, err error) {
	var foundStyleSheet = false
	for _, style := range bundle.Manifest.Styles {
		if style.StyleSheet == styleSheet {
			foundStyleSheet = true
			break
		}
	}
	if foundStyleSheet {
		styleSheetPath = path.Join(bundle.CSSDir, styleSheet)
	} else {
		err = fmt.Errorf("%s is not a style sheet in the bundle", styleSheet)
	}
	return
}
