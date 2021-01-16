import { Box, Button, Typography } from '@material-ui/core'

import { Node as APINode } from 'api/nodes'
import Paper from 'components-mui/Paper'
import React from 'react'
import Space from 'components-mui/Space'
import T from 'components/T'
import api from 'api'

interface Props {
  data: APINode
}

const Node: React.FC<Props> = ({ data }) => {
  const { name } = data

  return (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={3}>
        <Space display="flex" alignItems="center">
          <Typography component="div">{name}</Typography>
        </Space>
        <Button variant="outlined" color="primary" size="small">
          {T('common.use')}
        </Button>
      </Box>
    </Paper>
  )
}

export default Node
