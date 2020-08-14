import React, { useState, ReactEventHandler, ChangeEvent, useCallback, useEffect } from 'react';
import FileDropper from './FileDropper';
import axios from 'axios'
import fileDownload from 'js-file-download'
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import Button from '@material-ui/core/Button';
import { LinearProgress, Box, TextField, Container } from '@material-ui/core';
import { Alert, AlertTitle } from '@material-ui/lab'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      width: '100%',
    },
    button: {
      marginRight: theme.spacing(1),
    },
    instructions: {
      marginTop: theme.spacing(1),
      marginBottom: theme.spacing(1),
    },
  }),
)

const SELECT_PATCH_BUNDLE = 0
const SELECT_PDFS = 1
const SELECT_STYLE = 2
const REMIX_PDF = 3

const stepLabels = [
  'Upload patch bundle',
  'Upload source PDFs',
  'Select style',
  'Remix PDF'
]

function isLastStep(step: number) {
  return step === stepLabels.length - 1
}

function PatchStepper() {
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);
  const [outputStyle, setOutputStyle] = useState("")
  const [remixDownloaded, setRemixDownloaded] = useState(false)
  const [bundleFile, setBundleFile] = useState<File | undefined>()
  const [pdfsFile, setPdfsFile] = useState<File[]>([])
  const [downloadProgress, setDownloadProgress] = useState(0)
  const [patchFailure, setPatchFailure] = useState<String | null>(null)

  const onBundleDrop = useCallback(acceptedFiles => {
    setBundleFile(acceptedFiles[0])
  }, [])

  const onPdfsDrop = useCallback(acceptedFiles => {
    setPdfsFile(acceptedFiles)
  }, [])

  const handleNext = () => {
    if (!isStepComplete()) { //sanity check
      console.warn("PatchStepper#handleNext was called when step wasn't complete.")
      return
    }

    setActiveStep(activeStep + 1)
  };

  const isStepComplete = () => {
    return (
      (activeStep === 0 && bundleFile !== undefined) ||
      (activeStep === 1 && pdfsFile.length > 0) ||
      (activeStep === 2 && outputStyle !== "") ||
      (activeStep === 3 && remixDownloaded)
    )
  }

  const handleBack = () => {
    setActiveStep(activeStep - 1)
  }

  const handleReset = () => {
    setBundleFile(undefined)
    setPdfsFile([])
    setOutputStyle("")
    setRemixDownloaded(false)

    setActiveStep(0);
  };

  useEffect(() => {
    if (activeStep === 3 && downloadProgress === 0 && !patchFailure) {
      setDownloadProgress(1)
      submitPdfPatch()
    }
  }, [activeStep])

  function submitPdfPatch() {
    if (bundleFile === undefined) {
      console.warn("tried to submitPdfPatch without bundleFile")
      return
    }
    setPatchFailure(null)
    var formData = new FormData()
    formData.set('cssName', outputStyle)
    formData.append('bundle', bundleFile)
    pdfsFile.forEach(pdfFile => formData.append('pdfs', pdfFile))
    return axios.post('https://pdfpatch-gyeisy4svq-nn.a.run.app/api/v0/patch', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      responseType: 'blob'
    }).then(response => {
      fileDownload(response.data, 'patched.pdf')
      setDownloadProgress(100)
    }).catch(err => {
      err.response.data.text().then((error: String) => {
        console.log(error)
        setPatchFailure(error)
      })
    })
  }

  function getStepContent(step: number) {
    switch (step) {
      case 0:
        return (
          <FileDropper key="patch-bundle-dropper" multiple={false} accept={['application/zip']} onChange={onBundleDrop} files={bundleFile ? [bundleFile] : []} />
        )
      case 1:
        return (
          <FileDropper key="pdf-dropper" multiple={true} accept={['application/pdf']} onChange={onPdfsDrop} files={pdfsFile} />
        )
      case 2:
        return (
          <TextField
            required
            id="output-style"
            label="Style"
            variant="outlined"
            onChange={(event) => setOutputStyle(event.target.value)}
          />
        )
      case 3:
        if (patchFailure) {
          return (
            <>
              <LinearProgress variant="determinate" value={100} />
              <Alert severity="error">
                <AlertTitle>Error</AlertTitle>
                {patchFailure}<br />
                  You can go back, fix this error and then retry.
                </Alert>
            </>
          )
        } else if (downloadProgress < 100) {
          return (
            <>
              <LinearProgress />
              <Alert severity="info">
                <AlertTitle>Patching in progress</AlertTitle>
                This may take up to a minute.
              </Alert>
            </>
          )
        } else {
          return (
            <>
              <LinearProgress variant="determinate" value={100} />
              <Alert severity="success">
                <AlertTitle>Success</AlertTitle>
                  Your patched PDF is downloading now.
              </Alert>
            </>
          )
        }
      default:
        return 'Unknown step';
    }
  }

  function getActionButton() {
    switch (activeStep) {
      case SELECT_PATCH_BUNDLE:
      case SELECT_PDFS:
      case SELECT_STYLE:
        return (
          <Button
            variant="contained"
            color="primary"
            onClick={isLastStep(activeStep) ? handleReset : handleNext}
            className={classes.button}
            disabled={!isStepComplete()}
          >
            {isLastStep(activeStep) ? 'Remix' : 'Next'}
          </Button>
        )
      case REMIX_PDF:
        return (
          <Button
            variant="contained"
            color="primary"
            onClick={submitPdfPatch}
            className={classes.button}
            disabled={!patchFailure}
          >
            Retry
          </Button>
        )
    }
  }

  return (
    <div className={classes.root}>
      <Stepper activeStep={activeStep}>
        {stepLabels.map((label, index) => (
          <Step key={label}>
            <StepLabel>{label}</StepLabel>
          </Step>
        ))}
      </Stepper>
      <Container maxWidth="md">
        <Box my={4} >
          <Button disabled={activeStep === 0} onClick={handleBack} className={classes.button}>
            Back
          </Button>
          {getActionButton()}
        </Box>
        <Box m={1}>
          {getStepContent(activeStep)}
        </Box>
      </Container>
    </div>
  );
}

export default PatchStepper