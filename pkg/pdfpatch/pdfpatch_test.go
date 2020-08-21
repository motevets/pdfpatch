package pdfpatch_test

import (
	"bytes"
	"path"
	"time"

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
-Hell
+%0APAGE 1%0A%0AG
 o
+odbye
  fro
@@ -29,15 +29,33 @@
 r 1.
+%0A%0APAGE
  
+2%0A%0AAuf
  
-Hallo
+wiedersehen
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

	Describe("GeneratePatches", func() {
		const fixturesPath = "../../test/fixtures/multiple_patches"
		var pdfsDir = path.Join(fixturesPath, "pdfs")
		var markdownsDir = path.Join(fixturesPath, "markdowns")
		var pdfMarkdowns = []pdfpatch.PDFMarkdowns{
			{
				PDFFileName: "title_pages.pdf",
				MarkdownFileNames: []string{
					"title.md",
					"dedication.md",
				},
			},
			{
				PDFFileName: "chapter_1.pdf",
				MarkdownFileNames: []string{
					"chapter_1.md",
				},
			},
		}

		var pdfPatches = []pdfpatch.PDFPatch{
			{
				PDFFileName: "title_pages.pdf",
				Patch:       "@@ -1,8 +1,21 @@\n+%0APAGE 1%0A%0ANEW \n TITLE PA\n@@ -20,10 +20,18 @@\n PAGE\n- \n+%0A%0APAGE\n  \n+2%0A%0A\n Dedi\n@@ -52,9 +52,7 @@\n llow\n- men\n+s\n .\n+%0A\n",
			},
			{
				PDFFileName: "chapter_1.pdf",
				Patch:       "@@ -1,8 +1,17 @@\n+%0APAGE 3%0A%0A\n This is \n@@ -20,8 +20,30 @@\n apter 1.\n+%0A%0AIt's pretty great.%0A%0A\n",
			},
		}

		It("generates a patch from the PDF files", func() {
			patches, err := pdfpatch.GeneratePatches(pdfMarkdowns, pdfsDir, markdownsDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(patches)).To(Equal(2))
			Expect(patches).To(Equal(pdfPatches))
		})
	})

	Describe("ApplyPatch", func() {
		const fixturesPath = "../../test/fixtures/one_pdf_two_markdowns"
		var pdfPath = path.Join(fixturesPath, "original.pdf")
		var patchPath = path.Join(fixturesPath, "original.pdf.patch")
		It("applies the path to the PDF to make the desired output", func() {
			output, err := pdfpatch.ApplyPatch(pdfPath, patchPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(finalOutput))
		})
	})

	Describe("PatchPDF", func() {
		var outputPDFFile = "../../test/output/" + time.Now().Format(time.RFC3339) + "-patch-pdf-out.pdf"
		const fixturesPath = "../../test/fixtures/pdfs_patches_and_csses"
		var pdfFiles = []string{"title_pages.pdf", "chapter_1.pdf"}
		var patchesDir = path.Join(fixturesPath, "patches")
		var pdfsDir = path.Join(fixturesPath, "pdfs")
		var cssFile = path.Join(fixturesPath, "css/book.css")

		It("uses the patch bundle and input PDFs to make a patched PDF", func() {
			var err error
			err = pdfpatch.PatchPDF(pdfFiles, pdfsDir, patchesDir, cssFile, outputPDFFile)
			Expect(err).NotTo(HaveOccurred())
			numPages, text, err := statPDF(outputPDFFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(numPages).To(Equal(3))
			Expect(text).To(Equal("NEW TITLE PAGE1Dedicated to my fellows.2This is chapter 1.It's pretty great.3"))
		})
	})

	Describe("PatchBundle", func() {
		var outputPDFFile = "../../test/output/" + time.Now().Format(time.RFC3339) + "-patch-bundle-out.pdf"
		const pdfsDir = "../../test/fixtures/patch_bundle_pdfs"
		const bundlePath = "../../test/fixtures/patch_bundle.zip"
		const styleSheet = "large_print.css"

		It("uses the patch bundle and input PDFs to make a patched PDF", func() {
			var err error
			err = pdfpatch.PatchBundle(bundlePath, pdfsDir, styleSheet, outputPDFFile)
			Expect(err).NotTo(HaveOccurred())
			numPages, text, err := statPDF(outputPDFFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(numPages).To(Equal(3))
			Expect(text).To(Equal("NEW TITLE PAGE1Dedicated to myfellows.2This is chapter 1.It's pretty great.3"))
		})
	})
})

func statPDF(path string) (numPages int, text string, err error) {
	f, r, err := pdf.Open(path)
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
