import React from 'react'
import { createRoot } from 'react-dom/client'
import { Provider } from 'react-redux'
import { persistStore } from 'redux-persist'
import { PersistGate } from 'redux-persist/integration/react'

import IntlProvider from '@shared-ui/components/Atomic/IntlProvider'

import { store } from '@/store'
import { App } from '@/containers/App'
import reportWebVitals from './reportWebVitals'
import { DEVICE_AUTH_CODE_SESSION_KEY } from '@/constants'

// @ts-ignore
import languages from './languages/languages.json'
import appConfig from '@/config'

let persistor = persistStore(store)

reportWebVitals()

const BaseComponent = () => {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const isMockApp = window.location.pathname === '/devices-code-redirect' && !!code

    if (window.location.pathname === '/devices' && !!code) {
        window.location.hash = ''
        window.location.href = `${window.location.origin}/devices-code-redirect?code=${code}`
        return null
    }

    if (isMockApp) {
        localStorage.setItem(DEVICE_AUTH_CODE_SESSION_KEY, code)

        window.addEventListener('load', function () {
            setInterval(() => {
                if (localStorage.getItem(DEVICE_AUTH_CODE_SESSION_KEY)) {
                    window.close()
                }
            }, 200)
        })
    }

    return (
        <Provider store={store}>
            <PersistGate persistor={persistor}>
                <IntlProvider defaultLanguage={appConfig.defaultLanguage} languages={languages}>
                    <App mockApp={isMockApp} />
                </IntlProvider>
            </PersistGate>
        </Provider>
    )
}

const root = createRoot(document.getElementById('root') as Element)
root.render(<BaseComponent />)
