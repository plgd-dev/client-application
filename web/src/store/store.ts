import { configureStore } from '@reduxjs/toolkit'
import { setupListeners } from '@reduxjs/toolkit/query'
import { persistReducer } from 'redux-persist'
import storage from 'redux-persist/lib/storage'

import { StoreType as DeviceStoreType } from '@shared-ui/app/clientApp/Devices/slice'

import { createRootReducer } from './reducers'
import { StoreType as AppStoreType } from '../containers/App/slice'

export type CombinedStoreType = {
    activeNotifications: DeviceStoreType
    app: AppStoreType
}

const persistConfig = {
    key: 'root',
    storage: storage,
    blacklist: [],
}

const rootReducer = createRootReducer()

const persistedReducer = persistReducer(persistConfig, rootReducer)

const store = configureStore({
    reducer: persistedReducer,
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
            serializableCheck: false,
            immutableCheck: false,
        }),
})

setupListeners(store.dispatch)

export default store
