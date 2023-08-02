import { FC, useCallback, useContext, useEffect, useMemo, useRef } from 'react'
import { UserManagerSettings } from 'oidc-client-ts'
import { AuthProvider, UserManager } from 'oidc-react'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import {
    mergeConfig,
    useWellKnownConfiguration,
    WellKnownConfigType,
    WellKnownConfigurationState,
} from '@shared-ui/common/hooks/useWellKnownConfiguration'
import { security } from '@shared-ui/common/services'

import AppContext from './AppContext'
import { DEVICE_AUTH_MODE } from '@/constants'
import AppLoader from '@/containers/App/AppLoader/AppLoader'
import AppInner from '@/containers/App/AppInner/AppInner'
import { Props } from './App.types'
import './App.scss'

const App: FC<Props> = (props) => {
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin
    const [wellKnownConfig, setWellKnownConfig, reFetchConfig, wellKnownConfigError] = useWellKnownConfiguration(
        httpGatewayAddress,
        () => (wellKnownConfigState.current = WellKnownConfigurationState.MERGED)
    )
    const wellKnownConfigState = useRef(WellKnownConfigurationState.UNUSED)

    const inIframe = useCallback(() => {
        try {
            return window.self !== window.top
        } catch (e) {
            return true
        }
    }, [])

    const isIframe = useMemo(() => inIframe(), [inIframe])

    security.setGeneralConfig({
        httpGatewayAddress,
    })

    security.setWellKnowConfig(wellKnownConfig)

    useEffect(() => {
        if (isIframe && wellKnownConfig) {
            // send message that client-app is ready

            if (wellKnownConfigState.current === WellKnownConfigurationState.UNUSED) {
                wellKnownConfigState.current = WellKnownConfigurationState.REQUESTED
                // @ts-ignore
                window.top.postMessage(
                    {
                        key: 'PLGD_EVENT_MESSAGE',
                        clientReady: true,
                    },
                    '*'
                )
            }

            // listen on message
            window.addEventListener('message', function (event) {
                if (
                    wellKnownConfigState.current === WellKnownConfigurationState.REQUESTED &&
                    event.data.hasOwnProperty('key') &&
                    event.data.key === 'PLGD_EVENT_MESSAGE' &&
                    event.data.hasOwnProperty('PLGD_HUB_REMOTE_PROVISIONING_DATA')
                ) {
                    if (wellKnownConfig?.remoteProvisioning) {
                        setWellKnownConfig({
                            ...wellKnownConfig,
                            remoteProvisioning: mergeConfig(
                                wellKnownConfig.remoteProvisioning,
                                event.data.PLGD_HUB_REMOTE_PROVISIONING_DATA
                            ),
                        })
                    }
                }
            })
        }
    }, [wellKnownConfig, isIframe, setWellKnownConfig])

    if (wellKnownConfig && process.env.REACT_APP_TEST_PROVIDER_NAME) {
        // @ts-ignore
        wellKnownConfig.remoteProvisioning.deviceOauthClient.providerName = process.env.REACT_APP_TEST_PROVIDER_NAME
    }

    const setInitialize = (value = true) => {
        setWellKnownConfig({
            ...wellKnownConfig,
            isInitialized: value,
        } as WellKnownConfigType)
    }

    if (
        (!wellKnownConfig && !isIframe) ||
        (isIframe && wellKnownConfigState.current !== WellKnownConfigurationState.MERGED)
    ) {
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
        redirect_uri: window.location.href,
        extraQueryParams: {
            audience: wellKnownConfig?.remoteProvisioning?.webOauthClient?.audience || false,
        },
    } as UserManagerSettings)

    security.setUserManager(userManager)

    const onSignIn = async () => {
        window.location.href = window.location.href.split('?')[0]
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
            condition={!props.mockApp && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
            wrapper={Wrapper}
        >
            <AppInner
                configError={wellKnownConfigError}
                isIframe={isIframe}
                mockApp={props.mockApp}
                reFetchConfig={reFetchConfig}
                setInitialize={setInitialize}
                wellKnownConfig={wellKnownConfig as WellKnownConfigType}
            />
        </ConditionalWrapper>
    )
}

export const useAppConfig = () => useContext(AppContext)

export default App
