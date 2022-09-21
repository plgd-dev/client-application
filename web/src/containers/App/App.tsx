import { useContext, useEffect, useState } from 'react'
import AppContext from './AppContext'
import './App.scss'
import { WellKnownConfigType} from '@/containers/App/App.types'
import { getAppWellKnownConfiguration } from '@/containers/App/AppRest'
import {AuthProvider, UserManager} from 'oidc-react'
import { DEVICE_AUTH_MODE } from '@/constants'
import ConditionalWrapper from '@shared-ui/components/new/ConditionalWrapper'
import AppLoader from '@/containers/App/AppLoader/AppLoader'
import AppInner from '@/containers/App/AppInner/AppInner'
import { security } from '@shared-ui/common/services'

const App = () => {
    const [wellKnownConfig, setWellKnownConfig] = useState<WellKnownConfigType | undefined>(undefined)
    const [configError, setConfigError] = useState<any>(null)

    security.setGeneralConfig({
        httpGatewayAddress: process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin,
    })

    useEffect(() => {
        try {
            getAppWellKnownConfiguration().then((result) => {
                setWellKnownConfig({
                    ...result.data,
                })
            })
        } catch (e) {
            setConfigError(new Error('Could not retrieve the well-known configuration.'))
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

    const oidcCommonSettings = {
        authority: wellKnownConfig?.remoteProvisioning.authorization.authority || '',
        scope: wellKnownConfig?.remoteProvisioning.authorization.scopes.join?.(' ') || 'openid',
    }

    return (
        <ConditionalWrapper
            condition={wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
            wrapper={(child) => (
                <AuthProvider
                    {...oidcCommonSettings}
                    clientId={wellKnownConfig?.remoteProvisioning.authorization.clientId || ''}
                    redirectUri={window.location.origin}
                    onSignIn={async () => {
                        // remove auth params
                        window.location.hash = ''
                        window.location.href = window.location.origin
                    }}
                    automaticSilentRenew={true}
                    userManager={
                        new UserManager({
                            ...oidcCommonSettings,
                            client_id: wellKnownConfig?.remoteProvisioning.authorization.clientId,
                            redirect_uri: window.location.origin,
                            extraQueryParams: {
                                audience: wellKnownConfig?.remoteProvisioning.authorization.audience || false
                            }
                        })
                    }
                >
                    {child}
                </AuthProvider>
            )}
        >
            <AppInner wellKnownConfig={wellKnownConfig} configError={configError} setInitialize={setInitialize} />
        </ConditionalWrapper>
    )
}

export const useAppConfig = () => useContext(AppContext)

export default App
