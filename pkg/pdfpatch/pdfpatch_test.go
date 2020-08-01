package pdfpatch_test

import (
	"bytes"
	"fmt"
	"path"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
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
				Patch: heredoc.Doc(`
					@@ -1,8 +1,21 @@
					+%0APAGE 1%0A%0ANEW 
					 TITLE PA
					@@ -18,16 +18,24 @@
					 E PAGE%0A%0A
					+PAGE 2%0A%0A
					 Dedicate
					@@ -52,9 +52,7 @@
					 llow
					- men
					+s
					 .
					+%0A
				`),
			},
			{
				PDFFileName: "chapter_1.pdf",
				Patch: heredoc.Doc(`
					@@ -1,8 +1,17 @@
					+%0APAGE 3%0A%0A
					 This is 
					@@ -20,8 +20,30 @@
					 apter 1.
					+%0A%0AIt's pretty great.%0A%0A
				`),
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
		var outputPDFFile = "../../test/output/" + time.Now().Format(time.RFC3339) + "-out.pdf"
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
			fmt.Println(text)
			Expect(text).To(Equal("NEW TITLE PAGE1Dedicated to my fellows.2This is chapter 1.It's pretty great.3"))
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
