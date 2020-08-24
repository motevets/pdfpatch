import React from 'react';
import PatchStepper from './components/stepper/PatchStepper';

function App() {
  const remixApiUrl = process.env.REACT_APP_REMIX_API_HOST
  if(remixApiUrl === undefined) {
    throw new Error("missing REACT_APP_REMIX_API")
  }
  return (
    <PatchStepper remixApiHost={remixApiUrl}/>
  );
}

export default App;
