import React from 'react';
import { DropzoneArea } from 'material-ui-dropzone';
import styled from 'styled-components';
import { makeStyles, createStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(theme => createStyles({
  previewChip: {
    minWidth: 160,
    maxWidth: 210
  },
  dropzoneRoot: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center'
  }
}))

type FileDropperProps = {
  multiple: boolean
  accept: string[]
  onChange: (acceptedFiles: File[]) => void
  files: File[]
}

export default function FileDropper(props: FileDropperProps) {
  const classes = useStyles()
  const { multiple, accept, onChange, files } = props

  return(
    <DropzoneArea
      initialFiles={files}
      filesLimit={multiple ? 100 : 1}
      acceptedFiles={accept}
      onChange={onChange}
      showPreviews={true}
      showPreviewsInDropzone={false}
      useChipsForPreview
      dropzoneClass={classes.dropzoneRoot}
      previewGridProps={{container: { spacing: 1, direction: 'row' }}}
      previewChipProps={{classes: { root: classes.previewChip } }}
      previewText="Selected files"
      showAlerts={["error"]}
    />
  )
}