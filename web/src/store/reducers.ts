import { combineReducers } from 'redux'
import { enableBatching } from 'redux-batched-actions'

import devicesReducer from '@/containers/Devices/slice'
import appReducer from '@/containers/App/slice'

export const createRootReducer = () =>
    enableBatching(
        combineReducers({
            app: appReducer,
            devices: devicesReducer,
        })
    )
