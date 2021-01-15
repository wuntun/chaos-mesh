import { Box, Button, Typography } from '@material-ui/core'

import Paper from 'components-mui/Paper'
import React from 'react'
import Space from 'components-mui/Space'
import T from 'components/T'

const Node = () => {
  return (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={3}>
        <Space display="flex" alignItems="center">
          <Typography component="div">k8s-cluster-1</Typography>
        </Space>
        <Button variant="outlined" color="primary" size="small">
          {T('common.use')}
        </Button>
      </Box>
    </Paper>
  )
}

export default Node
