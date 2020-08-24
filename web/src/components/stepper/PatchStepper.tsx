import React, { useState, useEffect, useCallback } from 'react'
import FileDropper from './FileDropper'
import axios from 'axios'
import fileDownload from 'js-file-download'
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles'
import Stepper from '@material-ui/core/Stepper'
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel'
import Button from '@material-ui/core/Button'
import { LinearProgress, Box, Container, ListItem, List, ListItemText, ListItemIcon, FormControl, FormLabel, FormControlLabel, RadioGroup, Radio, Snackbar } from '@material-ui/core'
import { Alert, AlertTitle } from '@material-ui/lab'
import JSZip from 'jszip'
import yaml from 'js-yaml'
import CheckBoxOutlineBlankIcon from '@material-ui/icons/CheckBoxOutlineBlank';
import CheckBoxIcon from '@material-ui/icons/CheckBox';

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
    fileCheckbox: {
      minWidth: 'auto',
      paddingRight: theme.spacing(1)
    },
    radioArea: {
      backgroundColor: '#fff',
      width: '100%',
      padding: theme.spacing(2)
    },
    step: {
      "&$completed": {
        color: "#09af00"
      },
    },
    completed: {}
  }),
)

const SELECT_PATCH_BUNDLE = 0
const SELECT_PDFS = 1
const SELECT_STYLE = 2
const REMIX_PDF = 3

const stepLabels = [
  'Upload remix bundle',
  'Upload source PDFs',
  'Select style',
  'Remix PDF'
]

function isLastStep(step: number) {
  return step === stepLabels.length - 1
}

type Style = {
  name: string
  description: string
  styleSheet: string
}

type Source = {
  url: string
  fileName: string
  md5sum: string
}

type Manifest = {
  sources: Source[]
  styles: Style[]
}

type RawStyle = {
  name?: string
  description?: string
  style_sheet?: string
}
type RawSource = {
  url?: string
  file_name?: string
  md5sum?: string
}

type RawManifest = {
  sources?: Source[]
  styles?: Style[]
}

type Snack = {
  severity: "error" | "warning" | "info" | "success"
  message: string
}

class UploadedFilesList {
  _sourcesFilesMap: {
    [filename: string]: {
      file?: File
      source: Source
    }
  } = {}

  constructor(sources?: Source[]) {
    if(sources === undefined) {
      sources = []
    }
    sources.forEach(source => {
      this._sourcesFilesMap[source.fileName] = { source }
    })
  }

  clone(): UploadedFilesList {
    const clonedList = new UploadedFilesList([])
    for(const fileName in this._sourcesFilesMap) {
      clonedList._sourcesFilesMap[fileName] = Object.assign({}, this._sourcesFilesMap[fileName])
    }
    return clonedList
  }

  get fileNames(): string[] {
    return Object.keys(this._sourcesFilesMap)
  }

  get files(): File[] {
    const fileList: File[] = []
    for(const fileName in this._sourcesFilesMap) {
      const file = this._sourcesFilesMap[fileName].file
      if(file !== undefined) {
        fileList.push(file)
      } else {
        throw new Error(`required file ${fileName} is missing`)
      }
    }
    return fileList
  }

  addFile(file: File) {
    if(this._sourcesFilesMap[file.name] === undefined) {
      throw new Error(`${file.name} is not one of the required source PDF files`)
    } else {
      this._sourcesFilesMap[file.name].file = file
    }
  }

  isFilePresent(fileName: string): boolean {
    return this._sourcesFilesMap[fileName].file !== undefined
  }

  areAllFilesPresent(): boolean {
    for(const fileName in this._sourcesFilesMap) {
      if(!this.isFilePresent(fileName)) {
        return false
      }
    }
    return true
  }
}

function parseManifest(maybeManifest: unknown): Manifest {
  const manifest: Manifest = {
    sources: [],
    styles: []
  }
  if (!(maybeManifest instanceof Object)) {
    throw new Error('manifest is not an object')
  }
  const rawManifest = maybeManifest as RawManifest
  if (rawManifest.sources instanceof Array) {
    rawManifest.sources.forEach(rawSource => manifest.sources.push(parseSource(rawSource)))
  } else {
    throw new Error('manifest is missing sources key')
  }
  if (rawManifest.styles instanceof Array) {
    rawManifest.styles.forEach(rawStyle => manifest.styles.push(parseStyle(rawStyle)))
  } else {
    throw new Error('manifest is missing style key')
  }
  return manifest
}

function parseSource(sourceObj: RawSource): Source {
  const source: Source = {
    url: "",
    fileName: "",
    md5sum: ""
  }
  try {
    source.fileName = parseStringAttribute(sourceObj.file_name)
  } catch(e) {
    throw new Error(`${JSON.stringify(source)}\n${e}`)
  }
  try {
    source.url = parseStringAttribute(sourceObj.url)
  } catch(e) {
    throw new Error(`${source.fileName} url: ${e}`)
  }
  try {
    source.md5sum = parseStringAttribute(sourceObj.md5sum)
  } catch(e) {
    throw new Error(`${source.fileName} md5sum: ${e}`)
  }
  return source
}

function parseStyle(styleObj: RawStyle): Style {
  const style: Style = {
    name: "",
    description: "",
    styleSheet: ""
  }
  try {
    style.name = parseStringAttribute(styleObj.name)
  } catch(e) {
    throw new Error(`${JSON.stringify(styleObj)}\n${e}`)
  }
  try {
    style.description = parseStringAttribute(styleObj.description)
  } catch (e) {
    throw new Error(`${style.name} description: ${e}`)
  }
  try{
    style.styleSheet = parseStringAttribute(styleObj.style_sheet)
  } catch(e) {
    throw new Error(`${style.name} stylesheet: ${e}`)
  }
  return style
}

function parseStringAttribute(thing: any): string {
  if (typeof (thing) === 'string') {
    return thing
  } else {
    throw new Error(thing + ' is not a string')
  }
}

type PatchStepperProps = {
  remixApiHost: string
}

function PatchStepper(props: PatchStepperProps) {
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);
  const [outputStyle, setOutputStyle] = useState("")
  const [bundleFile, setBundleFile] = useState<File | undefined>()
  const [downloadProgress, setDownloadProgress] = useState(0)
  const [patchFailure, setPatchFailure] = useState<String | null>(null)
  const [availableStyles, setAvailableStyles] = useState<Style[]>([])
  const [snack, setSnack] = useState<Snack>()
  const [sourcesFilesMap, setSourcesFilesMap] = useState(new UploadedFilesList())

  const handleReset = () => {
    setActiveStep(0);
    setOutputStyle("")
    setBundleFile(undefined)
    setDownloadProgress(0)
    setPatchFailure(null)
    setAvailableStyles([])
    setSnack(undefined)
    setSourcesFilesMap(new UploadedFilesList())
  };

  const onBundleDrop = async (droppedFiles: File[]) => {
    const bundleFile = droppedFiles[0]
    let zip: JSZip
    let manifestText: string
    let manifest: Manifest
    let rawManifest: string | object | undefined
    try {
      zip = await JSZip.loadAsync(bundleFile)
    } catch (e) {
      setSnack({severity: 'error', message: bundleFile.name + ' is not a valid zip file'})
      return
    }
    const manifestFile = zip.file('manifest.yml')
    if (manifestFile === null) {
      setSnack({severity: 'error', message: "Invalid patch bundle, does not contain manifest.yml"})
      return
    }
    try {
      manifestText = await manifestFile.async("text")
      rawManifest = yaml.safeLoad(manifestText)
    } catch(e) {
      console.error(e)
      setSnack({severity: 'error', message: "Could not read manifest.yml in zip file"})
      return
    }
    try {
      manifest = parseManifest(rawManifest)
    } catch (e) {
      setSnack({severity: 'error', message: "error parsing manifest: " + e.toString()})
      return
    }
    setAvailableStyles(manifest.styles)
    setSourcesFilesMap(new UploadedFilesList(manifest.sources))
    setBundleFile(bundleFile)
    setActiveStep(activeStep + 1)
  }

  const onPdfsDrop = async (acceptedFiles: File[]) => {
    const nextSourcesFilesMap = sourcesFilesMap.clone()
    const errors: string[] = []
    acceptedFiles.forEach(file => {
      try {
        nextSourcesFilesMap.addFile(file)
      } catch(e) {
        errors.push(e.toString())
      }
    })
    if(errors.length > 0) {
      setSnack({severity: "error", message: errors.join(", ")})
      return
    }
    setSourcesFilesMap(nextSourcesFilesMap)
    await new Promise(res => setTimeout(res, 50)) //hack: dropzone will try to update state after after this component is otherwise unmounted due to race condition
    if(nextSourcesFilesMap.areAllFilesPresent()) {
      setActiveStep(activeStep + 1)
    }
  }

  const handleNext = () => {
    if (!isStepComplete()) { //sanity check
      console.warn("PatchStepper#handleNext was called when step wasn't complete.")
      return
    }

    setActiveStep(activeStep + 1)
  };

  const isStepComplete = (): boolean => {
    return !!(
      (activeStep === 0 && bundleFile !== undefined) ||
      (activeStep === 1 && sourcesFilesMap.areAllFilesPresent()) ||
      (activeStep === 2 && outputStyle !== "") ||
      (activeStep === 3 && (downloadProgress === 100 || patchFailure))
    )
  }

  const handleBack = () => {
    setActiveStep(activeStep - 1)
  }

  const submitPdfPatch = useCallback(() => {
    if (bundleFile === undefined) {
      console.warn("tried to submitPdfPatch without bundleFile")
      return
    }
    setPatchFailure(null)
    var formData = new FormData()
    formData.set('cssName', outputStyle)
    formData.append('bundle', bundleFile)
    sourcesFilesMap.files.forEach(pdfFile => formData.append('pdfs', pdfFile))
    return axios.post(`${props.remixApiHost}/api/v0/patch`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      responseType: 'blob'
    }).then(response => {
      fileDownload(response.data, 'patched.pdf')
      setDownloadProgress(100)
    }).catch(err => {
      err.response.data.text().then((error: String) => {
        setPatchFailure(error)
      })
    })
  }, [bundleFile, outputStyle, sourcesFilesMap, props.remixApiHost])

  useEffect(() => {
    if (activeStep === 3 && downloadProgress === 0 && !patchFailure) {
      setDownloadProgress(1)
      submitPdfPatch()
    }
  }, [activeStep, downloadProgress, patchFailure, submitPdfPatch])

  function getStepContent(step: number) {
    switch (step) {
      case SELECT_PATCH_BUNDLE:
        return (
          <>
            <FileDropper key="patch-bundle-dropper" multiple={false} onChange={onBundleDrop} />
            <List dense>
                <ListItem>
                  <ListItemIcon className={classes.fileCheckbox}>
                    {bundleFile ? <CheckBoxIcon /> : <CheckBoxOutlineBlankIcon /> }
                  </ListItemIcon>
                  <ListItemText>REMIX BUNDLE ({bundleFile ? bundleFile.name : ".zip"})</ListItemText>
                </ListItem>
            </List>
          </>
        )
      case SELECT_PDFS:
        return (
          <>
            <FileDropper key="pdf-dropper" multiple={true} onChange={onPdfsDrop} />
            <List dense>
              {sourcesFilesMap.fileNames.map(fileName =>
                <ListItem key={fileName}>
                  <ListItemIcon className={classes.fileCheckbox}>
                    {sourcesFilesMap.isFilePresent(fileName) ? <CheckBoxIcon /> : <CheckBoxOutlineBlankIcon /> }
                  </ListItemIcon>
                  <ListItemText>{fileName}</ListItemText>
                </ListItem>
              )}
            </List>
          </>
        )
      case SELECT_STYLE:
        return (
          <Box className={classes.radioArea}>
            <FormControl>
              <FormLabel>Style</FormLabel>
              <RadioGroup aria-label="style" name="style" className={classes.radioArea} value={outputStyle} onChange={(event) => setOutputStyle(event.target.value)}>
                {availableStyles.map(style => (
                  <FormControlLabel key={style.styleSheet} value={style.styleSheet} control={<Radio />} label={style.name} />
                ))}
              </RadioGroup>
            </FormControl>
          </Box>
        )
      case REMIX_PDF:
        if (patchFailure) {
          return (
            <>
              <LinearProgress variant="determinate" value={100} />
              <Alert severity="error">
                <AlertTitle>Error</AlertTitle>
                {patchFailure}<br />
                  You can start over, fix this error, and then retry. If the error persists, contact the author of this remix.
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
        console.error("Unexpected step while rendering step")
        return null
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
            {activeStep === REMIX_PDF - 1 ? 'Remix' : 'Next'}
          </Button>
        )
      case REMIX_PDF:
        return (
          <Button
            variant="contained"
            color="secondary"
            onClick={handleReset}
            className={classes.button}
            disabled={!isStepComplete()}
          >
            Start Over
          </Button>
        )
      default:
        console.error("Unexpected step while getting action button")
        return null
    }
  }

  return (
    <div className={classes.root}>
      <Stepper activeStep={activeStep}>
        {stepLabels.map((label, index) => (
          <Step completed={activeStep > index || (index === REMIX_PDF && downloadProgress === 100)} key={label}>
            <StepLabel error={!!(index === REMIX_PDF && patchFailure)} StepIconProps={{classes: {root: classes.step, completed: classes.completed}}}>{label}</StepLabel>
          </Step>
        ))}
      </Stepper>
      <Container maxWidth="md">
        <Box my={4} >
          <Button disabled={activeStep === 0 || activeStep === REMIX_PDF} onClick={handleBack} className={classes.button}>
            Back
          </Button>
          {getActionButton()}
        </Box>
        <Box m={1}>
          {getStepContent(activeStep)}
        </Box>
      </Container>
      <Snackbar open={!!snack} autoHideDuration={6000} onClose={() => setSnack(undefined)} >
        {snack && <Alert severity={snack.severity} onClose={() => setSnack(undefined)}>{snack.message}</Alert>}
      </Snackbar>
    </div>
  );
}

export default PatchStepper