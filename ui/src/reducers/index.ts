import { combineReducers } from 'redux'
import experiments from 'slices/experiments'
import globalStatus from 'slices/globalStatus'
import navigation from 'slices/navigation'
import nodes from 'slices/nodes'
import settings from 'slices/settings'

export default combineReducers({
  globalStatus,
  navigation,
  nodes,
  experiments,
  settings,
})
