package extractor_test

import (
	"github.com/motevets/pdfpatch/pkg/extractor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextFromPDF", func() {

	When("the folder and PDF files exist", func() {
		It("returns the text extracted from all of the PDFs in order of the filenames in the array", func() {
			text, err := extractor.TextFromPdf("../../test/fixtures/", []string{"hello_from_page_1.pdf", "hallo_von_seite_2.pdf"})
			Expect(err).NotTo(HaveOccurred())
			Expect(text).To(Equal("Hello from page 1.\nHallo von Seite 2.\n"))
		})
	})
})
