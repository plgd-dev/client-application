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

    const oidcCommonSettings = {
        authority: wellKnownConfig.remoteProvisioning.authorization.authority || '',
        scope: wellKnownConfig.remoteProvisioning.authorization.scopes.join?.(' '),
    }

    return (
        <ConditionalWrapper
            condition={wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
            wrapper={(child: any) => (
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
                                audience: wellKnownConfig?.remoteProvisioning.authorization.audience || false,
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
