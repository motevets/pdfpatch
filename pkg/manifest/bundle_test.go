package manifest_test

import (
	"path"

	"github.com/motevets/pdfpatch/pkg/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bundle", func() {
	const bundlePath = "../../test/fixtures/patch_bundle.zip"
	var (
		bundle manifest.Bundle
		err    error
	)

	BeforeEach(func() {
		bundle, err = manifest.UnpackBundle(bundlePath)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe(".UnpackBundle", func() {
		It("returns a struct with the manifest", func() {
			Expect(bundle.Manifest).To(Equal(manifest.Manifest{
				Sources: []manifest.Source{
					{
						URL: "http://example.com/title_pages.pdf",
						PatchedFiles: []string{
							"title.md",
							"dedication.md",
						},
					},
					{
						URL: "http://example.com/chapter_1.pdf",
						PatchedFiles: []string{
							"chapter_1.md",
						},
					},
				},
				Styles: []manifest.Style{
					{
						Name:        "Traditional",
						Description: "Classical rendering of this iconic text.",
						StyleSheet:  "book.css",
					},
					{
						Name:        "Large Print",
						Description: "Large print format.",
						StyleSheet:  "large_print.css",
					},
				},
			}))
		})

		It("returns a Bundle with the directory with CSS files", func() {
			expectedCSSFile := path.Join(bundle.CSSDir, "book.css")
			Expect(expectedCSSFile).To(BeARegularFile())
		})

		It("returns a Bundle with the directory with patch files", func() {
			Expect(path.Join(bundle.PatchesDir, "title_pages.pdf.patch")).To(BeARegularFile())
			Expect(path.Join(bundle.PatchesDir, "chapter_1.pdf.patch")).To(BeARegularFile())
		})

		It("has the path to the manifest file", func() {
			Expect(bundle.ManifestPath).To(BeARegularFile())
		})
	})

	Describe("Bundle#StyleSheetPath", func() {
		It("returns the path to the stylesheet by it's name", func() {
			Expect(bundle.CSSFilePath("large_print.css")).To(BeARegularFile())
		})

		When("the stylesheet does not exist", func() {
			It("returns and error", func() {
				_, err := bundle.CSSFilePath("nope.css")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("nope.css is not a style sheet in the bundle"))
			})
		})
	})
})
