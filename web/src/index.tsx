import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { store } from '@/store'
import { App } from '@/containers/App'
import IntlProvider from '@shared-ui/components/new/IntlProvider'
import reportWebVitals from './reportWebVitals'

reportWebVitals()

ReactDOM.render(
    <Provider store={store}>
        <IntlProvider>
            <App />
        </IntlProvider>
    </Provider>,
    document.getElementById('root')
)
