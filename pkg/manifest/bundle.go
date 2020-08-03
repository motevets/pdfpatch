package manifest

import (
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
