import { FC, useContext } from 'react'
import { UserManagerSettings } from 'oidc-client-ts'
import { AuthProvider, UserManager } from 'oidc-react'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import { useWellKnownConfiguration, WellKnownConfigType } from '@shared-ui/common/hooks/useWellKnownConfiguration'
import { security } from '@shared-ui/common/services'

import AppContext from './AppContext'
import { DEVICE_AUTH_MODE } from '@/constants'
import AppLoader from '@/containers/App/AppLoader/AppLoader'
import AppInner from '@/containers/App/AppInner/AppInner'
import { Props } from './App.types'
import './App.scss'

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
        scope: wellKnownConfig?.remoteProvisioning?.webOauthClient?.scopes?.join?.(' ') || '',
    })

    const userManager = new UserManager({
        ...getOidcCommonSettings(),
        automaticSilentRenew: true,
        client_id: wellKnownConfig?.remoteProvisioning?.webOauthClient?.clientId || '',
        redirect_uri: window.location.origin,
        extraQueryParams: {
            audience: wellKnownConfig?.remoteProvisioning?.webOauthClient?.audience || false,
        },
    } as UserManagerSettings)

    security.setUserManager(userManager)

    const Wrapper = (child: any) => (
        <AuthProvider
            {...getOidcCommonSettings()}
            automaticSilentRenew={true}
            clientId={wellKnownConfig?.remoteProvisioning?.webOauthClient?.clientId || ''}
            onSignIn={async (userData) => {
                // remove auth params
                window.location.hash = ''
                window.location.href = window.location.origin
            }}
            redirectUri={window.location.origin}
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
                configError={wellKnownConfigError}
                mockApp={props.mockApp}
                reFetchConfig={reFetchConfig}
                setInitialize={setInitialize}
                wellKnownConfig={wellKnownConfig}
            />
        </ConditionalWrapper>
    )
}

export const useAppConfig = () => useContext(AppContext)

export default App
