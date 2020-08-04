package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"runtime/debug"

	"github.com/motevets/pdfpatch/pkg/pdfpatch"
)

func patch(w http.ResponseWriter, r *http.Request) {
	var (
		err                error
		pdfFilesHeaders    []*multipart.FileHeader
		bundleFileHeaders  []*multipart.FileHeader
		uploadedBundleFile multipart.File
		assetsDir          string
		pdfsDir            string
		bundleFilePath     string
		bundleFile         *os.File
		cssName            string
		outputPDFPath      string
		outputPDFFile      *os.File
	)

	err = r.ParseMultipartForm(10 << 20) // use max 10mb of memory for upload
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	cssName = r.FormValue("cssName")
	if cssName == "" {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("Missing \"cssName\" field"))
		return
	}

	pdfFilesHeaders = r.MultipartForm.File["pdfs"]
	if pdfFilesHeaders == nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("Missing \"pdfs\" files field"))
		return
	}

	bundleFileHeaders = r.MultipartForm.File["bundle"]
	if bundleFileHeaders == nil || len(bundleFileHeaders) == 0 {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("Missing \"bundle\" file field"))
		return
	}

	assetsDir, err = ioutil.TempDir("", "bundle-assets-")
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	pdfsDir = path.Join(assetsDir, "pdfs")
	err = os.Mkdir(pdfsDir, 0755)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	for _, pdfFileHeader := range pdfFilesHeaders {
		var (
			uploadedPdfFile multipart.File
			savedPdfFile    *os.File
		)
		uploadedPdfFile, err = pdfFileHeader.Open()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err)
			return
		}
		defer uploadedPdfFile.Close()

		savedPdfFile, err = os.OpenFile(path.Join(pdfsDir, pdfFileHeader.Filename), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err)
			return
		}
		defer savedPdfFile.Close()
		_, err = io.Copy(savedPdfFile, uploadedPdfFile)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err)
			return
		}
	}
	log.Printf("pdfs written to %s", pdfsDir)

	uploadedBundleFile, err = bundleFileHeaders[0].Open()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	defer uploadedBundleFile.Close()

	bundleFilePath = path.Join(assetsDir, bundleFileHeaders[0].Filename)
	bundleFile, err = os.OpenFile(bundleFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	defer bundleFile.Close()

	_, err = io.Copy(bundleFile, uploadedBundleFile)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	log.Printf("bundle written to %s", bundleFile.Name())

	outputPDFPath = path.Join(assetsDir, "output.pdf")

	err = pdfpatch.PatchBundle(bundleFilePath, pdfsDir, cssName, outputPDFPath)
	if err != nil {
		writeErr(w, http.StatusUnprocessableEntity, err)
		return
	}
	log.Println("output PDF written to " + outputPDFPath)

	outputPDFFile, err = os.Open(outputPDFPath)
	if err != nil {
		writeErr(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, outputPDFFile)
}

func ServeApi(port string) (err error) {
	serverAddress := fmt.Sprintf(":%s", port)
	http.HandleFunc("/api/v0/patch", patch)
	log.Printf("pdfpatch server running and listening on %s", port)
	return http.ListenAndServe(serverAddress, nil)
}

type Handler interface {
	RequestParmeters() interface{}
	Exec(interface{}) (interface{}, error)
}

func writeErr(w http.ResponseWriter, statusCode int, err error) {
	var msg string
	if statusCode == http.StatusInternalServerError {
		msg = http.StatusText(statusCode)
		log.Println(err)
		debug.PrintStack()
	} else {
		msg = err.Error()
	}
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, msg)
}
