import React, { useState, ReactEventHandler, ChangeEvent, useCallback } from 'react';
import FileDropper from './FileDropper';
import axios from 'axios'
import fileDownload from 'js-file-download'

function Stepper() {
  const [cssStyle, setCssStyle] = useState("")
  const [bundleFile, setBundleFile] = useState<File | undefined>()
  const [pdfsFile, setPdfsFile] = useState<File[]>([])

  const onBundleDrop = useCallback(acceptedFiles => {
    setBundleFile(acceptedFiles[0])
  }, [])

  const onPdfsDrop = useCallback(acceptedFiles => {
    setPdfsFile(acceptedFiles)
  }, [])

  return (
    <React.Fragment>
      <FileDropper label='Step 1: Upload patch-bundle (.zip file)' multiple={false} accept='application/zip' onDrop={onBundleDrop}/>
      <FileDropper label='Step 2: Upload original PDFs (.pdf files)' multiple={true} accept='application/pdf' onDrop={onPdfsDrop}/>
      <label>
        Style: 
        <input type="text" onChange={(event) => setCssStyle(event.target.value)}/>
      </label>
      <input type="submit" value="Patch" onClick={submitFn(bundleFile, pdfsFile, cssStyle)} />
    </React.Fragment>
  );
}

function submitFn(bundleFile : File | undefined, pdfFiles : File[], cssStyle : string) {
  return function() {
    submitPdfPatch(bundleFile as File, pdfFiles, `${cssStyle}.css`)
  }
}

function submitPdfPatch(bundleFile : File, pdfFiles : File[], cssFile : string) {
  var formData = new FormData()
  formData.set('cssName', cssFile)
  formData.append('bundle', bundleFile)
  pdfFiles.forEach(pdfFile => formData.append('pdfs', pdfFile))
  return axios.post('https://pdfpatch-gyeisy4svq-nn.a.run.app/api/v0/patch', formData, {
    headers: {'Content-Type': 'multipart/form-data' },
    responseType: 'blob'
  }).then(response => {
    fileDownload(response.data, 'patched.pdf')
  })
}

export default Stepper;