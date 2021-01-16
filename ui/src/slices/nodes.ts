import { PayloadAction, createAsyncThunk, createSlice } from '@reduxjs/toolkit'

import LS from 'lib/localStorage'
import { Node } from 'api/nodes'
import api from 'api'

export const getNodes = createAsyncThunk('node/list', async () => (await api.nodes.nodes()).data)

const initialState: {
  nodes: Node[]
  node: string
  kind: Node['kind']
} = {
  nodes: [],
  node: '',
  kind: 'k8s',
}

const nodesSlice = createSlice({
  name: 'nodes',
  initialState,
  reducers: {
    setNode(state, action: PayloadAction<string>) {
      const name = action.payload

      state.node = name

      LS.set('node', name)
    },
    setKind(state, action: PayloadAction<Node['kind']>) {
      const kind = action.payload

      state.kind = kind

      LS.set('node-kind', kind)
    },
  },
  extraReducers: (builder) => {
    builder.addCase(getNodes.fulfilled, (state, action) => {
      state.nodes = Object.values(action.payload)
    })
  },
})

export const { setNode, setKind } = nodesSlice.actions

export default nodesSlice.reducer
