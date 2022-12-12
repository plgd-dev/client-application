import { useContext } from 'react'
import AppContext from './AppContext'
import './App.scss'
import { AuthProvider, UserManager } from 'oidc-react'
import { DEVICE_AUTH_MODE } from '@/constants'
import ConditionalWrapper from '@shared-ui/components/new/ConditionalWrapper'
import AppLoader from '@/containers/App/AppLoader/AppLoader'
import AppInner from '@/containers/App/AppInner/AppInner'
import { security } from '@shared-ui/common/services'
import { useWellKnownConfiguration, WellKnownConfigType } from '@shared-ui/common/hooks/useWellKnownConfiguration'
import { UserManagerSettings } from 'oidc-client-ts'

const App = () => {
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin
    const [wellKnownConfig, setWellKnownConfig, reFetchConfig, wellKnownConfigError] =
        useWellKnownConfiguration(httpGatewayAddress)

    security.setGeneralConfig({
        httpGatewayAddress,
    })

    const setInitialize = (value = true) => {
        setWellKnownConfig({
            ...wellKnownConfig,
            isInitialized: value,
        } as WellKnownConfigType)
    }

    if (!wellKnownConfig) {
        return <AppLoader />
    }

    const getOidcCommonSettings = () => ({
        authority: wellKnownConfig.remoteProvisioning?.authority || '',
        scope: wellKnownConfig.remoteProvisioning?.webOauthClient.scopes.join?.(' '),
    })

    return (
        <ConditionalWrapper
            condition={wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
            wrapper={(child: any) => (
                <AuthProvider
                    {...getOidcCommonSettings()}
                    clientId={wellKnownConfig?.remoteProvisioning?.webOauthClient.clientId || ''}
                    redirectUri={window.location.origin}
                    onSignIn={async () => {
                        // remove auth params
                        window.location.hash = ''
                        window.location.href = window.location.origin
                    }}
                    automaticSilentRenew={true}
                    userManager={
                        new UserManager({
                            ...getOidcCommonSettings(),
                            client_id: wellKnownConfig?.remoteProvisioning?.webOauthClient.clientId,
                            redirect_uri: window.location.origin,
                            extraQueryParams: {
                                audience: wellKnownConfig?.remoteProvisioning?.webOauthClient.audience || false,
                            },
                        } as UserManagerSettings)
                    }
                >
                    {child}
                </AuthProvider>
            )}
        >
            <AppInner
                wellKnownConfig={wellKnownConfig}
                configError={wellKnownConfigError}
                setInitialize={setInitialize}
                reFetchConfig={reFetchConfig}
            />
        </ConditionalWrapper>
    )
}

export const useAppConfig = () => useContext(AppContext)

export default App
