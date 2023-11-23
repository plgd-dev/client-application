import { FC, useCallback, useContext, useMemo } from 'react'
import { UserManagerSettings } from 'oidc-client-ts'
import { AuthProvider, UserManager } from 'oidc-react'
import { useIntl } from 'react-intl'
import { useDispatch, useSelector } from 'react-redux'
import { ThemeProvider } from '@emotion/react'
import get from 'lodash/get'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import {
    useWellKnownConfiguration,
    WellKnownConfigType,
    mergeConfig,
} from '@shared-ui/common/hooks/useWellKnownConfiguration'
import { security } from '@shared-ui/common/services'
import AppContext from '@shared-ui/app/share/AppContext'
import { DEVICE_AUTH_MODE } from '@shared-ui/app/clientApp/constants'
import { Props } from '@shared-ui/app/clientApp/App/App.types'
import AppLoader from '@shared-ui/app/clientApp/App/AppLoader'
import { useAppTheme } from '@shared-ui/common/hooks/use-app-theme'
import { getTheme } from '@shared-ui/app/clientApp/App/AppRest'

import AppInner from '@/containers/App/AppInner/AppInner'
import { messages as t } from './App.i18n'
import './App.scss'
import { setTheme, setThemes, storeWellKnownConfig } from '@/containers/App/slice'
import { CombinedStoreType } from '@/store/store'

const App: FC<Props> = (props) => {
    const { formatMessage: _ } = useIntl()
    const httpGatewayAddress = process.env.REACT_APP_HTTP_GATEWAY_ADDRESS || window.location.origin
    const dispatch = useDispatch()

    const appStore = useSelector((state: CombinedStoreType) => state.app)
    const [wellKnownConfig, setWellKnownConfig, reFetchConfig, wellKnownConfigError] = useWellKnownConfiguration(
        httpGatewayAddress,
        {
            onSuccess: (wellKnownCfg) => {
                dispatch(storeWellKnownConfig(wellKnownCfg))
            },
        }
    )

    const [theme, themeError] = useAppTheme({
        getTheme,
        setTheme,
        setThemes,
    })

    const currentTheme = useMemo(() => appStore.configuration?.theme ?? 'plgd', [appStore.configuration?.theme])

    const getThemeData = useCallback(() => {
        if (theme) {
            const index = theme.findIndex((i: any) => Object.keys(i)[0] === currentTheme)
            if (index >= 0) {
                return get(theme[index], `${currentTheme}`, {})
            }
        }

        return {}
    }, [theme, currentTheme])

    const mergedWellKnownConfig = useMemo(() => {
        if (wellKnownConfig) {
            return mergeConfig(wellKnownConfig, appStore.userWellKnownConfig)
        }
        return undefined
    }, [appStore.userWellKnownConfig, wellKnownConfig])

    security.setGeneralConfig({
        httpGatewayAddress,
    })

    const authority = useMemo(
        () => mergedWellKnownConfig?.remoteProvisioning?.authority,
        [mergedWellKnownConfig?.remoteProvisioning?.authority]
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

    const useAuthLib = useMemo(
        () =>
            !props.mockApp && !!authority && mergedWellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509,
        [authority, props.mockApp, mergedWellKnownConfig?.deviceAuthenticationMode]
    )

    const getScope = useCallback((scope: string | []) => {
        if (typeof scope === 'string') {
            return scope
        }

        return scope?.join?.(' ')
    }, [])

    if (!wellKnownConfig || !theme) {
        return <AppLoader i18n={{ loading: _(t.loading) }} />
    } else {
        security.setWellKnowConfig(wellKnownConfig)
    }

    // Render an error box
    if (themeError) {
        return <div className='client-error-message'>{`${_(t.authError)}: ${themeError?.message}`}</div>
    }

    const getOidcCommonSettings = () => ({
        authority: authority || '',
        scope: getScope(mergedWellKnownConfig?.remoteProvisioning?.webOauthClient?.scopes) || '',
    })

    const userManager = new UserManager({
        ...getOidcCommonSettings(),
        automaticSilentRenew: true,
        client_id: mergedWellKnownConfig?.remoteProvisioning?.webOauthClient?.clientId || '',
        redirect_uri: window.location.href,
        extraQueryParams: {
            audience: mergedWellKnownConfig?.remoteProvisioning?.webOauthClient?.audience || false,
        },
    } as UserManagerSettings)

    if (useAuthLib) {
        security.setUserManager(userManager)
    }

    const onSignIn = async () => {
        window.location.href = window.location.href.split('?')[0]
    }

    const Wrapper = (child: any) => (
        <AuthProvider
            {...getOidcCommonSettings()}
            automaticSilentRenew={true}
            clientId={mergedWellKnownConfig?.remoteProvisioning?.webOauthClient?.clientId || ''}
            onSignIn={onSignIn}
            redirectUri={window.location.href}
            userManager={userManager}
        >
            {child}
        </AuthProvider>
    )

    return (
        <ThemeProvider theme={getThemeData()}>
            <ConditionalWrapper condition={useAuthLib} wrapper={Wrapper}>
                <AppInner
                    configError={wellKnownConfigError}
                    initializedByAnother={!authority}
                    mockApp={props.mockApp}
                    reFetchConfig={reFetchConfig}
                    setInitialize={setInitialize}
                    updateWellKnownConfig={updateWellKnownConfig}
                    wellKnownConfig={mergedWellKnownConfig as WellKnownConfigType}
                />
            </ConditionalWrapper>
        </ThemeProvider>
    )
}

export const useAppConfig = () => useContext(AppContext)

export default App
