import { combineReducers } from 'redux'
import { enableBatching } from 'redux-batched-actions'

import devicesReducer from '@shared-ui/app/clientApp/Devices/slice'

import appReducer from '@/containers/App/slice'

export const createRootReducer = () =>
    enableBatching(
        combineReducers({
            app: appReducer,
            devices: devicesReducer,
        })
    )
