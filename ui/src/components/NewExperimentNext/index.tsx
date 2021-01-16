import { Box, Grid } from '@material-ui/core'

import AdjustIcon from '@material-ui/icons/Adjust'
import Alert from '@material-ui/lab/Alert'
import LoadFrom from './LoadFrom'
import React from 'react'
import Step1 from './Step1'
import Step2 from './Step2'
import Step3 from './Step3'
import { useStoreSelector } from 'store'

const NewExperiment = () => {
  const { node } = useStoreSelector((state) => state.nodes)

  return (
    <Grid container spacing={6}>
      <Grid item xs={12} md={8}>
        {node && (
          <Alert severity="info" icon={<AdjustIcon />}>
            当前处于 {node}
          </Alert>
        )}
        <Box mt={6}>
          <Step1 />
        </Box>
        <Box mt={6}>
          <Step2 />
        </Box>
        <Box mt={6}>
          <Step3 />
        </Box>
      </Grid>
      <Grid item xs={12} md={4}>
        <LoadFrom />
      </Grid>
    </Grid>
  )
}

export default NewExperiment
