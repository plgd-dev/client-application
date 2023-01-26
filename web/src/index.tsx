import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { store } from '@/store'
import { App } from '@/containers/App'
import IntlProvider from '@shared-ui/components/new/IntlProvider'
import reportWebVitals from './reportWebVitals'
import { DEVICE_AUTH_CODE_SESSION_KEY } from '@/constants'

reportWebVitals()

const BaseComponent = () => {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const isMockApp = window.location.pathname === '/devices' && !!code

    if (isMockApp) {
        const { detect } = require('detect-browser')
        const browser = detect()
        browser && console.log(browser.name)
        localStorage.setItem(DEVICE_AUTH_CODE_SESSION_KEY, code)

        window.addEventListener('load', function () {
            setTimeout(() => {
                if (localStorage.getItem(DEVICE_AUTH_CODE_SESSION_KEY)) {
                    // safari cant close window - show mocApp
                    if (browser && browser.name === 'safari') {
                        window.location.href = `${window.location.origin}/devices`
                    } else {
                        window.close()
                    }
                }
            }, 200)
        })
    }

    return (
        <Provider store={store}>
            <IntlProvider>
                <App mockApp={isMockApp} />
            </IntlProvider>
        </Provider>
    )
}

ReactDOM.render(<BaseComponent />, document.getElementById('root'))
