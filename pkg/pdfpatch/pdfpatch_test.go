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
PAGE 1

Goodbye from chapter 1.

PAGE 2

Auf wiedersehen von Kapitel 2.

`

const computedPatch = `@@ -1,9 +1,20 @@
-H
+%0APAGE 1%0A%0AGoodby
 e
-llo
  fro
@@ -31,13 +31,31 @@
 1.%0A%0A
-Hallo
+PAGE 2%0A%0AAuf wiedersehen
  von
@@ -65,8 +65,9 @@
 pitel 2.
+%0A
`

var _ = Describe("pdfpatch", func() {
	Describe("GeneratePatch", func() {
		const fixturesPath = "../../test/fixtures/one_pdf_two_markdowns"
		var pdfPath = path.Join(fixturesPath, "original.pdf")
		var markdownPaths = []string{path.Join(fixturesPath, "chapter_1.md"), path.Join(fixturesPath, "chapter_2.md")}

		It("generates a patch from the PDF files", func() {
			patch, err := pdfpatch.GeneratePatch(pdfPath, markdownPaths)
			Expect(err).ToNot(HaveOccurred())
			Expect(patch).To(Equal(computedPatch))
		})
	})

	Describe("ApplyPatch", func() {
		It("applies the path to the PDF to make the desired output", func() {
			output, err := pdfpatch.ApplyPatch("../../test/fixtures/one_pdf_two_markdowns", []string{"original.pdf"}, computedPatch)
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
