package pdfpatch_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPdfpatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pdfpatch Suite")
}
