import React from 'react';
import { useDropzone } from 'react-dropzone';
import styled from 'styled-components';

type DropStyleProps = {
  isDragAccept?: boolean
  isDragReject?: boolean
  isDragActive?: boolean
}

const getColor = (props: DropStyleProps) => {
  if (props.isDragAccept) {
    return '#00e676';
  }
  if (props.isDragReject) {
    return '#ff1744';
  }
  if (props.isDragActive) {
    return '#2196f3';
  }
  return '#eeeeee';
}

const Container = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  border-width: 2px;
  border-radius: 2px;
  border-color: ${props => getColor(props as DropStyleProps)};
  border-style: dashed;
  background-color: #fafafa;
  color: #bdbdbd;
  outline: none;
  transition: border .24s ease-in-out;
`;

type FileDropperProps = {
  label: string
  multiple: boolean
  accept: string[] | string
  onDrop: (acceptedFiles: File[]) => void
}

export default function FileDropper(props: FileDropperProps) {
  const { label, multiple, accept, onDrop } = props

  const {
    getRootProps,
    getInputProps,
    acceptedFiles,
    isDragActive,
    isDragAccept,
    isDragReject
  } = useDropzone({ multiple, accept, onDrop });

  return (
    <div className="container">
      <Container {...getRootProps({ isDragActive, isDragAccept, isDragReject })}>
        <label>
          <input {...getInputProps()} />
          <p>{label}</p>
        </label>
        <p>{acceptedFiles.map(file => file.name).join(", ")}</p>
      </Container>
    </div>
  );
}