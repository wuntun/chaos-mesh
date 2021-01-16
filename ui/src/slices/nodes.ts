import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'

import { Node } from 'api/nodes'
import api from 'api'

export const getNodes = createAsyncThunk('node/list', async () => (await api.nodes.nodes()).data)

const initialState: {
  nodes: Node[]
} = {
  nodes: [],
}

const nodesSlice = createSlice({
  name: 'nodes',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(getNodes.fulfilled, (state, action) => {
      state.nodes = Object.values(action.payload)
    })
  },
})

export default nodesSlice.reducer
