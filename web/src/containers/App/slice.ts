import { createSlice } from '@reduxjs/toolkit'

export type StoreType = {
    version: {
        latest?: string
        latest_url?: string
        requestedDatetime?: string
    }
    wellKnownConfig: any
    userWellKnownConfig: any
    configuration: {
        theme: string
        themes: string[]
    }
}

const initialState: StoreType = {
    version: {},
    wellKnownConfig: {},
    userWellKnownConfig: {},
    configuration: {
        theme: '',
        themes: ['plgd'],
    },
}

const { reducer, actions } = createSlice({
    name: 'app',
    initialState,
    reducers: {
        setVersion(state, { payload }) {
            state.version = payload
        },
        storeWellKnownConfig(state, { payload }) {
            state.wellKnownConfig = payload
        },
        storeUserWellKnownConfig(state, { payload }) {
            state.userWellKnownConfig = payload
        },
        setThemes(state, { payload }) {
            state.configuration.themes = payload
        },
        setTheme(state, { payload }) {
            state.configuration.theme = payload
        },
    },
})

// Actions
export const { setVersion, storeWellKnownConfig, storeUserWellKnownConfig, setTheme, setThemes } = actions

// Reducer
export default reducer
