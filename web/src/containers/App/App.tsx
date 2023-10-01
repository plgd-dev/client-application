import { FC, useContext, useMemo, useState } from 'react'
import { UserManagerSettings } from 'oidc-client-ts'
import { AuthProvider, UserManager } from 'oidc-react'
import { useIntl } from 'react-intl'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import { useWellKnownConfiguration, WellKnownConfigType } from '@shared-ui/common/hooks/useWellKnownConfiguration'
import { security } from '@shared-ui/common/services'
import AppContext from '@shared-ui/app/clientApp/App/AppContext'
import { DEVICE_AUTH_MODE } from '@shared-ui/app/clientApp/constants'
import { Props } from '@shared-ui/app/clientApp/App/App.types'
import AppLoader from '@shared-ui/app/clientApp/App/AppLoader'

import AppInner from '@/containers/App/AppInner/AppInner'
import { messages as t } from './App.i18n'
import './App.scss'
import UnInitializedScreen from '@shared-ui/app/clientApp/App/UnInitializedScreen'

const App: FC<Props> = (props) => {
    const { formatMessage: _ } = useIntl()
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [test, setTest] = useState(1)
    const [wellKnownConfig, setWellKnownConfig, reFetchConfig, wellKnownConfigError] = useWellKnownConfiguration(
        httpGatewayAddress,
        undefined,
        (w) => {
            setTest(2)
        }
    )

    security.setGeneralConfig({
        httpGatewayAddress,
    })

    security.setWellKnowConfig(wellKnownConfig)

    const authority = useMemo(
        () => wellKnownConfig?.remoteProvisioning?.authority,
        [wellKnownConfig?.remoteProvisioning?.authority]
    )

    const setInitialize = (value = true) => {
        setWellKnownConfig({
            ...wellKnownConfig,
            isInitialized: value,
        } as WellKnownConfigType)
    }

    const updateWellKnownConfig = (data: WellKnownConfigType) => {
        setWellKnownConfig(data)
    }

    if (!wellKnownConfig) {
        return <AppLoader i18n={{ loading: _(t.loading) }} />
    }

    const getOidcCommonSettings = () => ({
        authority: authority || '',
        scope: wellKnownConfig?.remoteProvisioning?.webOauthClient?.scopes?.join?.(' ') || '',
    })

    const userManager = new UserManager({
        ...getOidcCommonSettings(),
        automaticSilentRenew: true,
        client_id: wellKnownConfig?.remoteProvisioning?.webOauthClient?.clientId || '',
        redirect_uri: window.location.href,
        extraQueryParams: {
            audience: wellKnownConfig?.remoteProvisioning?.webOauthClient?.audience || false,
        },
    } as UserManagerSettings)

    security.setUserManager(userManager)

    const onSignIn = async () => {
        window.location.href = window.location.href.split('?')[0]
    }

    if (!wellKnownConfig.isInitialized) {
        return <UnInitializedScreen updateWellKnownConfig={setWellKnownConfig} wellKnownConfig={wellKnownConfig} />
    }

    const Wrapper = (child: any) => (
        <AuthProvider
            {...getOidcCommonSettings()}
            automaticSilentRenew={true}
            clientId={wellKnownConfig?.remoteProvisioning?.webOauthClient?.clientId || ''}
            onSignIn={onSignIn}
            redirectUri={window.location.href}
            userManager={userManager}
        >
            {child}
        </AuthProvider>
    )

    return (
        <ConditionalWrapper
            condition={
                !props.mockApp && !!authority && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
            }
            wrapper={Wrapper}
        >
            <AppInner
                configError={wellKnownConfigError}
                initializedByAnother={!authority}
                mockApp={props.mockApp}
                reFetchConfig={reFetchConfig}
                setInitialize={setInitialize}
                updateWellKnownConfig={updateWellKnownConfig}
                wellKnownConfig={wellKnownConfig as WellKnownConfigType}
            />
        </ConditionalWrapper>
    )
}

export const useAppConfig = () => useContext(AppContext)

export default App
