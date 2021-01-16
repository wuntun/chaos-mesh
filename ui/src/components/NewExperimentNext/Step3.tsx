import { Box, Button, Typography } from '@material-ui/core'
import { setAlert, setAlertOpen } from 'slices/globalStatus'
import { useStoreDispatch, useStoreSelector } from 'store'

import DoneAllIcon from '@material-ui/icons/DoneAll'
import Paper from 'components-mui/Paper'
import PaperTop from 'components-mui/PaperTop'
import PublishIcon from '@material-ui/icons/Publish'
import React from 'react'
import T from 'components/T'
import api from 'api'
import { parseSubmit } from 'lib/formikhelpers'
import { resetNewExperiment } from 'slices/experiments'
import { useHistory } from 'react-router-dom'
import { useIntl } from 'react-intl'

const Step3 = () => {
  const { step1, step2, basic, target } = useStoreSelector((state) => state.experiments)
  const { kind } = useStoreSelector((state) => state.nodes)
  const dispatch = useStoreDispatch()

  const history = useHistory()
  const intl = useIntl()

  const submitExperiment = () => {
    const parsedValues = parseSubmit({
      ...basic,
      target,
    })

    if (process.env.NODE_ENV === 'development') {
      console.debug('Current node kind', kind)
      console.debug('Debug parsedValues:', parsedValues)
    }

    const newFunc = kind === 'physic' ? api.experiments.physicNew : api.experiments.newExperiment

    newFunc(parsedValues)
      .then(() => {
        dispatch(
          setAlert({
            type: 'success',
            message: intl.formatMessage({ id: 'common.createSuccessfully' }),
          })
        )
        dispatch(setAlertOpen(true))
        dispatch(resetNewExperiment())

        history.push('/experiments')
      })
      .catch(console.error)
  }

  return (
    <>
      {step1 && step2 && (
        <Paper>
          <PaperTop title={T('common.submit')} />
          <Box p={6} textAlign="center">
            <DoneAllIcon fontSize="large" />
            <Typography>{T('newE.complete')}</Typography>
            <Box mt={6} textAlign="right">
              <Button variant="contained" color="primary" startIcon={<PublishIcon />} onClick={submitExperiment}>
                {T('common.submit')}
              </Button>
            </Box>
          </Box>
        </Paper>
      )}
    </>
  )
}

export default Step3
