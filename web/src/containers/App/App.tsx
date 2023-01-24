import { FC, useContext } from 'react'
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
import { Props } from './App.types'

const App: FC<Props> = (props) => {
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin
    const [wellKnownConfig, setWellKnownConfig, reFetchConfig, wellKnownConfigError] =
        useWellKnownConfiguration(httpGatewayAddress)

    security.setGeneralConfig({
        httpGatewayAddress,
    })

    security.setWellKnowConfig(wellKnownConfig)

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
        authority: wellKnownConfig?.remoteProvisioning?.authority || '',
        scope: wellKnownConfig?.remoteProvisioning?.webOauthClient.scopes.join?.(' '),
    })

    const userManager = new UserManager({
        ...getOidcCommonSettings(),
        automaticSilentRenew: true,
        client_id: wellKnownConfig?.remoteProvisioning?.webOauthClient.clientId,
        redirect_uri: window.location.origin,
        extraQueryParams: {
            audience: wellKnownConfig?.remoteProvisioning?.webOauthClient.audience || false,
        },
    } as UserManagerSettings)

    security.setUserManager(userManager)

    const Wrapper = (child: any) => (
        <AuthProvider
            {...getOidcCommonSettings()}
            automaticSilentRenew={true}
            clientId={wellKnownConfig?.remoteProvisioning?.webOauthClient.clientId || ''}
            redirectUri={window.location.origin}
            onSignIn={async (userData) => {
                // remove auth params
                window.location.hash = ''
                window.location.href = window.location.origin
            }}
            userManager={userManager}
        >
            {child}
        </AuthProvider>
    )

    return (
        <ConditionalWrapper
            condition={!props.mockApp && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
            wrapper={Wrapper}
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
