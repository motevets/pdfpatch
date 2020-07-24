package pdfpatch_test

import (
	"bytes"
	"path"

	"github.com/ledongthuc/pdf"
	"github.com/motevets/pdfpatch/pkg/pdfpatch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const finalOutput = `
Goodbye from page 1.
Auf wiedersehen von Seite 2.
`
const computedPatch = "@@ -1,9 +1,12 @@\n-H\n+%0AGoodby\n e\n-llo\n  fro\n@@ -19,13 +19,23 @@\n  1.%0A\n-Hallo\n+Auf wiedersehen\n  von\n"

var _ = Describe("pdfpatch", func() {
	Describe("GeneratePatch", func() {
		It("generates a patch from the PDF files", func() {
			patch, err := pdfpatch.GeneratePatch("../../test/fixtures/", []string{"hello_from_page_1.pdf", "hallo_von_seite_2.pdf"}, finalOutput)
			Expect(err).ToNot(HaveOccurred())
			Expect(patch).To(Equal(computedPatch))
		})
	})

	Describe("ApplyPatch", func() {
		It("applies the path to the PDF to make the desired output", func() {
			output, err := pdfpatch.ApplyPatch("../../test/fixtures/", []string{"hello_from_page_1.pdf", "hallo_von_seite_2.pdf"}, computedPatch)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(finalOutput))
		})
	})

	Describe("PatchPDF", func() {
		const outputPDFFile = "../../test/output/out.pdf"
		const fixturesPath = "../../test/fixtures"
		var pdfFiles = []string{"hello_from_page_1.pdf", "hallo_von_seite_2.pdf"}
		var cssFile = path.Join(fixturesPath, "patch_bundle/book.css")
		var patchFile = path.Join(fixturesPath, "patch_bundle/book.patch")

		It("uses the patch bundle and input PDFs to make a patched PDF", func() {
			var err error
			err = pdfpatch.PatchPDF(fixturesPath, pdfFiles, patchFile, cssFile, outputPDFFile)
			Expect(err).NotTo(HaveOccurred())
			numPages, text, err := statPDF(outputPDFFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(numPages).To(Equal(2))
			Expect(text).To(Equal("Goodbye from page 1.1Auf wiedersehen von Seite 2.2"))
		})
	})
})

func statPDF(path string) (numPages int, text string, err error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return
	}
	buf.ReadFrom(b)
	return r.NumPage(), buf.String(), nil
}
