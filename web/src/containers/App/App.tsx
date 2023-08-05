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
import AppContext from '@shared-ui/app/clientApp/App/AppContext'
import { DEVICE_AUTH_MODE } from '@shared-ui/app/clientApp/constants'
import { getParentAppWellKnownConfiguration } from '@shared-ui/app/clientApp/App/AppRest'
import { Props } from '@shared-ui/app/clientApp/App/App.types'

import AppLoader from '@/containers/App/AppLoader/AppLoader'
import AppInner from '@/containers/App/AppInner/AppInner'
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
        if (isIframe) {
            const queryParams = new URLSearchParams(window.location.search)
            const parentWellKnownConfigUrl = queryParams.get('wellKnownConfigUrl')

            if (wellKnownConfig && parentWellKnownConfigUrl) {
                if (wellKnownConfigState.current === WellKnownConfigurationState.UNUSED) {
                    wellKnownConfigState.current = WellKnownConfigurationState.REQUESTED

                    getParentAppWellKnownConfiguration(parentWellKnownConfigUrl).then((response) => {
                        if (response) {
                            wellKnownConfigState.current = WellKnownConfigurationState.RECEIVED
                            const data = response.data

                            setWellKnownConfig({
                                ...wellKnownConfig,
                                remoteProvisioning: mergeConfig(wellKnownConfig.remoteProvisioning!, data),
                            })
                        }
                    })
                }
            }
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
