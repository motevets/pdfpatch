package pdfpatch_test

import (
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
})
