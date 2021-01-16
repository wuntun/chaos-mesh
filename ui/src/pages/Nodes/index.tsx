import { Box, Button, Grid, Typography } from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useStoreDispatch, useStoreSelector } from 'store'

import AddIcon from '@material-ui/icons/Add'
import AddNode from './AddNode'
import AutoSizer from 'react-virtualized-auto-sizer'
import BubbleView from './BubbleView'
import ConfirmDialog from 'components-mui/ConfirmDialog'
import Node from './Node'
import Paper from 'components-mui/Paper'
import PaperTop from 'components-mui/PaperTop'
import T from 'components/T'
import { getNodes } from 'slices/nodes'

const Nodes = () => {
  const { nodes } = useStoreSelector((state) => state.nodes)
  const k8sNodes = nodes.filter((node) => node.kind === 'k8s')
  const physicNodes = nodes.filter((node) => node.kind === 'physic')
  const dispatch = useStoreDispatch()

  const [openAddNode, setOpenAddNode] = useState(false)

  useEffect(() => {
    dispatch(getNodes())
  }, [dispatch])

  const handleAddNodeSubmitCallback = () => {
    dispatch(getNodes())

    setOpenAddNode(false)
  }

  return (
    <Grid container spacing={6} style={{ height: '100%' }}>
      <Grid item sm={12} md={4}>
        <Paper>
          <PaperTop title={T('common.status')} subtitle={T('nodes.status.subtitle')}>
            <Button
              variant="outlined"
              size="small"
              color="primary"
              startIcon={<AddIcon />}
              onClick={() => setOpenAddNode(true)}
            >
              {T('nodes.add.title')}
            </Button>
          </PaperTop>
          <Box height={450}>
            <AutoSizer>{({ width, height }) => <BubbleView width={width} height={height} />}</AutoSizer>
          </Box>
        </Paper>
        <ConfirmDialog
          open={openAddNode}
          setOpen={setOpenAddNode}
          title={T('nodes.add.title')}
          dialogProps={{
            PaperProps: {
              variant: 'outlined',
              style: { width: 500, minWidth: 300 },
            },
          }}
        >
          <AddNode onSubmitCallback={handleAddNodeSubmitCallback} />
        </ConfirmDialog>
      </Grid>
      <Grid item sm={12} md={8}>
        {k8sNodes.length > 0 && (
          <Box mb={6}>
            <Box mb={6}>
              <Typography variant="button">{T('nodes.list.k8s')}</Typography>
            </Box>
            <Grid container spacing={3}>
              {k8sNodes.map((node) => (
                <Grid key={node.name} item xs={12}>
                  <Node data={node} />
                </Grid>
              ))}
            </Grid>
          </Box>
        )}
        {physicNodes.length > 0 && (
          <Box mb={6}>
            <Box mb={6}>
              <Typography variant="button">{T('nodes.list.physic')}</Typography>
            </Box>
            <Grid container spacing={3}>
              {physicNodes.map((node) => (
                <Grid key={node.name} item xs={12}>
                  <Node data={node} />
                </Grid>
              ))}
            </Grid>
          </Box>
        )}
      </Grid>
    </Grid>
  )
}

export default Nodes
