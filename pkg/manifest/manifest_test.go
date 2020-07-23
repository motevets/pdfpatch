package manifest_test

import (
	"io/ioutil"
	"log"

	"github.com/motevets/pdfpatch/pkg/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const validManifest = `
sources:
    - url: http://example.com/foo.pdf
      md5sum: 2b00042f7481c7b056c4b410d28f33cf
    - url: http://example.com/bar.pdf
      md5sum: b9eb9d6228842aeb05d64f30d56b361e`

var _ = Describe("manifest", func() {
	Describe("ParseFile", func() {
		When("a valid manifest exists at path", func() {
			var manifestFilePath string
			var theManifest manifest.Manifest
			BeforeEach(func() {
				var err error
				manifestFilePath = writeTmpFile(validManifest)
				theManifest, err = manifest.ParseFile(manifestFilePath)
				Expect(err).NotTo(HaveOccurred())
			})

			It("parses the manifest", func() {
				Expect(len(theManifest.Sources)).To(Equal(2))
				Expect(theManifest.Sources[0]).To(Equal(manifest.Source{
					URL:    "http://example.com/foo.pdf",
					Md5Sum: "2b00042f7481c7b056c4b410d28f33cf",
				}))
				Expect(theManifest.Sources[1]).To(Equal(manifest.Source{
					URL:    "http://example.com/bar.pdf",
					Md5Sum: "b9eb9d6228842aeb05d64f30d56b361e",
				}))
			})
		})

	})

	Describe("Manifest", func() {
		Describe("#SourceFileNames", func() {
			It("returns the file names in the order specified", func() {
				theManifest := manifest.Manifest{
					Sources: []manifest.Source{
						{
							URL: "http://example.com/foo.pdf",
						},
						{
							URL: "http://example.com/bar.pdf",
						},
					},
				}
				Expect(theManifest.SourceFileNames()).To(Equal([]string{"foo.pdf", "bar.pdf"}))
			})
		})
	})

	Describe("Source", func() {
		Describe("#FileName", func() {
			It("returns the filename based off of the URL", func() {
				source := manifest.Source{
					URL: "http://example.com/foo.pdf",
				}
				Expect(source.FileName()).To(Equal("foo.pdf"))
			})

			When("a URL is invalid", func() {
				It("returns an error", func() {
					source := manifest.Source{
						URL: "!@#$%^&*()",
					}
					_, err := source.FileName()
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})

func writeTmpFile(content string) string {
	tmpfile, err := ioutil.TempFile("", "manifest.yml")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}
