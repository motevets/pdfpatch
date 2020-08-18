import React from 'react';
import { DropzoneAreaBase } from 'material-ui-dropzone';
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
  onChange: (acceptedFiles: File[]) => void
}

export default function FileDropper(props: FileDropperProps) {
  const classes = useStyles()
  const { multiple, onChange } = props

  return(
    <DropzoneAreaBase
      fileObjects={[]}
      filesLimit={multiple ? 100 : 1}
      onDrop={onChange}
      showPreviews={false}
      showPreviewsInDropzone={false}
      dropzoneClass={classes.dropzoneRoot}
      showAlerts={["error"]}
    />
  )
}