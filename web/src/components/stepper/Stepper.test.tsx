import React from 'react';
import { render, RenderResult, fireEvent, act } from '@testing-library/react';
import Stepper from './Stepper';
import mockAxios from 'jest-mock-axios';

describe('on initial render', () => {
  let stepper: RenderResult
  let bundleFileUpload: HTMLElement
  let pdfFilesUpload: HTMLElement
  let cssStyleSelect: HTMLElement
  let submitButton
  let bundleFile: File
  let pdfFiles: File[]
  let cssSelection: string
  beforeEach(async () => {
    await new Promise(resolve => setImmediate(resolve))
    bundleFile = new File(['...'], 'patch-bundle.zip')
    pdfFiles = [
      new File(['...'], '1.pdf'),
      new File(['...'], '2.pdf')
    ]
    cssSelection = 'foo.css'
    stepper = render(<Stepper />)
  });

  afterEach(() => {
    mockAxios.reset();
  })

  it('asks for the bundle file', () => {
    bundleFileUpload = stepper.getByText(/upload patch-bundle/i)
  });

  describe('after adding a patch bundle', () => {
    beforeEach(async () => {
      bundleFileUpload = stepper.getByLabelText(/upload patch-bundle/i)
      fireEvent.change(bundleFileUpload, { target: { files: [bundleFile] } })
      await stepper.findByText(bundleFile.name, { exact: false })
    })

    describe('after adding all of the required PDFs', () => {
      beforeEach(async () => {
        pdfFilesUpload = stepper.getByLabelText(/upload original PDFs/i)
        fireEvent.change(pdfFilesUpload, { target: { files: pdfFiles } })
        await stepper.findByText(bundleFile.name, { exact: false })
      })

      describe('after specifying the style', () => {
        beforeEach(() => {
          cssStyleSelect = stepper.getByLabelText(/Style/)
          fireEvent.change(cssStyleSelect, { target: { value: 'large-print' } })
        })

        describe('after submitting', () => {
          beforeEach(() => {
            submitButton = stepper.getByText(/Patch/i, {selector: "input[type=submit]"})
            fireEvent.click(submitButton)
          })

          it('POST to the pdfpatch API', () => {
            let formData = new FormData()
            formData.set('cssName', 'large-print.css')
            formData.append('bundle', bundleFile)
            formData.append('pdfs', pdfFiles[0])
            formData.append('pdfs', pdfFiles[1])
            expect(mockAxios.post).toHaveBeenCalledWith('https://pdfpatch-gyeisy4svq-nn.a.run.app/api/v0/patch', formData, {
              // data: bodyFormData,
              headers: {'Content-Type': 'multipart/form-data' }
            })
          })
        })
      })
    })
  })
})
