import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { store } from '@/store'
import { App } from '@/containers/app'
import { IntlProvider } from '@shared-ui/components/old/intl-provider'
import reportWebVitals from './reportWebVitals'

fetch('/web_configuration.json')
  .then(response => response.json())
  .then(config => {
    const { httpGatewayAddress } = config
    if (!httpGatewayAddress) {
      throw new Error(
        'httpGatewayAddress must be set in web_configuration.json'
      )
    }

    // If you want to start measuring performance in your app, pass a function
    // to log results (for example: reportWebVitals(console.log))
    // or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
    reportWebVitals()

    ReactDOM.render(
      <Provider store={store}>
        <IntlProvider>
          <App config={config} />
        </IntlProvider>
      </Provider>,
      document.getElementById('root')
    )
  })
  .catch(error => {
    const rootDiv = document.getElementById('root')
    rootDiv.innerHTML = `<div class="client-error-message">${error.message}</div>`
  })
