import { useContext, useEffect, useState } from 'react'
import AppContext from './AppContext'
import './App.scss'
import { WellKnownConfigType } from '@/containers/App/App.types'
import { getAppWellKnownConfiguration } from '@/containers/App/AppRest'
import { AuthProvider } from 'oidc-react'
import { DEVICE_AUTH_MODE } from '@/constants'
import ConditionalWrapper from '@shared-ui/components/new/ConditionalWrapper'
import AppLoader from '@/containers/App/AppLoader/AppLoader'
import AppInner from '@/containers/App/AppInner/AppInner'
import { User } from 'oidc-client-ts'
import { security } from '@shared-ui/common/services'

const App = () => {
  const [wellKnownConfig, setWellKnownConfig] = useState<
    WellKnownConfigType | undefined
  >(undefined)

  const [configError, setConfigError] = useState<any>(null)

  security.setGeneralConfig({
    httpGatewayAddress:
      process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin,
  })

  useEffect(() => {
    try {
      getAppWellKnownConfiguration().then(result => {
        setWellKnownConfig({
          ...result.data,
          // deviceAuthenticationMode: DEVICE_AUTH_MODE.PRE_SHARED_KEY,
        })
      })
    } catch (e) {
      setConfigError(
        new Error('Could not retrieve the well-known configuration.')
      )
    }
  }, []) // eslint-disable-line

  const setInitialize = (value = true) => {
    // @ts-ignore
    setWellKnownConfig({
      ...wellKnownConfig,
      isInitialized: value,
    })
  }

  if (!wellKnownConfig) {
    return <AppLoader />
  }

  return (
    <ConditionalWrapper
      condition={
        wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
      }
      wrapper={child => (
        <AuthProvider
          authority={
            `${wellKnownConfig?.remoteProvisioning.authorization.authority}` ||
            ''
          }
          clientId={
            wellKnownConfig?.remoteProvisioning.authorization.clientId || ''
          }
          redirectUri={window.location.origin}
          onSignIn={async (user: User | null) => {
            // remove auth params
            window.location.hash = ''
            window.location.href = window.location.origin
          }}
          automaticSilentRenew={true}
        >
          {child}
        </AuthProvider>
      )}
    >
      <AppInner
        wellKnownConfig={wellKnownConfig}
        configError={configError}
        setInitialize={setInitialize}
      />
    </ConditionalWrapper>
  )
}

export const useAppConfig = () => useContext(AppContext)

export default App
