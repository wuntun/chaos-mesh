import { Box, Button, IconButton, Typography } from '@material-ui/core'
import { getNodes, setKind, setNode } from 'slices/nodes'
import { useStoreDispatch, useStoreSelector } from 'store'

import { Node as APINode } from 'api/nodes'
import DeleteOutlinedIcon from '@material-ui/icons/DeleteOutlined'
import Paper from 'components-mui/Paper'
import React from 'react'
import Space from 'components-mui/Space'
import T from 'components/T'
import api from 'api'

interface Props {
  data: APINode
}

const Node: React.FC<Props> = ({ data }) => {
  const { name, kind, config } = data

  const { node } = useStoreSelector((state) => state.nodes)
  const dispatch = useStoreDispatch()

  const handleDeleteNode = (name: string) => () => {
    api.nodes.del(name).then(() => dispatch(getNodes()))
  }

  const handleUseNode = (d: APINode) => () => {
    dispatch(setNode(d.name))
    dispatch(setKind(d.kind))

    api.auth.node(d.name)
  }

  return (
    <Paper>
      <Box display="flex" justifyContent="space-between" alignItems="center" p={3}>
        <Space display="flex" alignItems="center">
          <Typography component="div">{name}</Typography>
          {kind === 'physic' && (
            <Typography variant="body2" color="textSecondary">
              {window.atob(config)}
            </Typography>
          )}
        </Space>
        <Space>
          <IconButton
            color="primary"
            component="span"
            size="small"
            disabled={name === node}
            onClick={handleDeleteNode(name)}
          >
            <DeleteOutlinedIcon />
          </IconButton>
          <Button
            variant="outlined"
            color="primary"
            size="small"
            disabled={name === node}
            onClick={handleUseNode(data)}
          >
            {name === node ? T('common.using') : T('common.use')}
          </Button>
        </Space>
      </Box>
    </Paper>
  )
}

export default Node
