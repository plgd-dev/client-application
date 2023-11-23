import { useRef, useState, useMemo, useCallback } from 'react'
import { Helmet } from 'react-helmet'
import { BrowserRouter } from 'react-router-dom'
import { User } from 'oidc-react'
import jwtDecode from 'jwt-decode'
import get from 'lodash/get'

import { ToastContainer } from '@shared-ui/components/Atomic/Notification'
import { BrowserNotificationsContainer } from '@shared-ui/components/Atomic/Toast'
import { useLocalStorage, WellKnownConfigType } from '@shared-ui/common/hooks'
import { security } from '@shared-ui/common/services'
import AppContext from '@shared-ui/app/share/AppContext'
import { DEVICE_AUTH_MODE } from '@shared-ui/app/clientApp/constants'

import appConfig from '@/config'
import { Props } from './AppInner.types'
import AppLayout from '@/containers/App/AppLayout/AppLayout'
import { AppLayoutRefType } from '@/containers/App/AppLayout/AppLayout.types'
import { storeUserWellKnownConfig } from '@/containers/App/slice'

const getBuildInformation = (wellKnownConfig: WellKnownConfigType) => ({
    buildDate: wellKnownConfig?.buildDate || '',
    commitHash: wellKnownConfig?.commitHash || '',
    commitDate: wellKnownConfig?.commitDate || '',
    releaseUrl: wellKnownConfig?.releaseUrl || '',
    version: wellKnownConfig?.version || '',
})

const AppInner = (props: Props) => {
    const {
        wellKnownConfig,
        configError,
        reFetchConfig,
        setInitialize,
        initializedByAnother: initializedByAnotherProp,
        updateWellKnownConfig,
    } = props
    const buildInformation = getBuildInformation(wellKnownConfig)

    const appLayoutRef = useRef<AppLayoutRefType | null>(null)

    if (wellKnownConfig && wellKnownConfig.remoteProvisioning) {
        security.setWebOAuthConfig({
            authority: wellKnownConfig.remoteProvisioning.authority,
            certificateAuthority: wellKnownConfig.remoteProvisioning.certificateAuthority,
            clientId: wellKnownConfig.remoteProvisioning.webOauthClient?.clientId,
            redirect_uri: window.location.origin,
        })
    }

    const isInitializedByAnother = useMemo(
        () =>
            initializedByAnotherProp &&
            wellKnownConfig.isInitialized &&
            wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509,
        [wellKnownConfig, initializedByAnotherProp]
    )

    const [initializedByAnother, setInitializedByAnother] = useState(
        initializedByAnotherProp === undefined ? false : isInitializedByAnother
    )
    const [suspectedUnauthorized, setSuspectedUnauthorized] = useState(false)
    const [collapsed, setCollapsed] = useLocalStorage('leftPanelCollapsed', true)
    const unauthorizedCallback = useCallback(() => {
        setSuspectedUnauthorized(true)

        reFetchConfig().then((newWellKnownConfig: WellKnownConfigType) => {
            if (appLayoutRef.current) {
                const userData: User = appLayoutRef.current?.getAuthProviderRef().getUserData()
                const parsedData = jwtDecode(userData.access_token)
                const ownerId = get(parsedData, newWellKnownConfig.remoteProvisioning?.jwtOwnerClaim as string, '')

                if (ownerId !== newWellKnownConfig?.owner) {
                    setInitializedByAnother(true)
                }
            }

            setSuspectedUnauthorized(false)
        })
    }, [reFetchConfig])

    const contextValue = useMemo(
        () => ({
            unauthorizedCallback,
            collapsed,
            setCollapsed,
            buildInformation: buildInformation || undefined,
            isHub: false,
            updateAppWellKnownConfig: storeUserWellKnownConfig,
        }),
        [buildInformation, collapsed, setCollapsed, unauthorizedCallback]
    )

    // Render an error box with a config error
    if (configError) {
        return <div className='client-error-message'>{configError?.message}</div>
    }

    return (
        <AppContext.Provider value={contextValue}>
            <BrowserRouter>
                <Helmet defaultTitle={appConfig.appName} titleTemplate={`%s | ${appConfig.appName}`} />
                <AppLayout
                    initializedByAnother={!!initializedByAnother}
                    mockApp={props.mockApp}
                    ref={appLayoutRef}
                    setInitialize={setInitialize}
                    suspectedUnauthorized={suspectedUnauthorized}
                    updateWellKnownConfig={updateWellKnownConfig}
                    wellKnownConfig={wellKnownConfig}
                />
                <ToastContainer />
                <BrowserNotificationsContainer />
            </BrowserRouter>
        </AppContext.Provider>
    )
}

AppInner.displayName = 'AppInner'

export default AppInner
