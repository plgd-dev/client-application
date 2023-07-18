import { useRef, useState, useMemo, useCallback } from 'react'
import { Helmet } from 'react-helmet'
import { BrowserRouter } from 'react-router-dom'
import { User } from 'oidc-react'
import jwtDecode from 'jwt-decode'
import get from 'lodash/get'
import { ThemeProvider } from '@emotion/react'

import { ToastContainer } from '@shared-ui/components/Atomic/Notification'
import { BrowserNotificationsContainer } from '@shared-ui/components/Atomic/Toast'
import { useLocalStorage, WellKnownConfigType } from '@shared-ui/common/hooks'
import { security } from '@shared-ui/common/services'
import light from '@shared-ui/components/Atomic/_theme/light'

import AppContext from '@/containers/App/AppContext'
import appConfig from '@/config'
import { Props } from './AppInner.types'
import AppLayout from '@/containers/App/AppLayout/AppLayout'
import { AppLayoutRefType } from '@/containers/App/AppLayout/AppLayout.types'

const getBuildInformation = (wellKnownConfig: WellKnownConfigType) => ({
    buildDate: wellKnownConfig?.buildDate || '',
    commitHash: wellKnownConfig?.commitHash || '',
    commitDate: wellKnownConfig?.commitDate || '',
    releaseUrl: wellKnownConfig?.releaseUrl || '',
    version: wellKnownConfig?.version || '',
})

const AppInner = (props: Props) => {
    const { wellKnownConfig, configError, reFetchConfig, setInitialize } = props
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

    const [initializedByAnother, setInitializedByAnother] = useState(false)
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
        }),
        [buildInformation, collapsed, setCollapsed, unauthorizedCallback]
    )

    // Render an error box with a config error
    if (configError) {
        return <div className='client-error-message'>{configError?.message}</div>
    }

    console.log({ initializedByAnother })

    return (
        <AppContext.Provider value={contextValue}>
            <ThemeProvider theme={light}>
                <BrowserRouter>
                    <Helmet defaultTitle={appConfig.appName} titleTemplate={`%s | ${appConfig.appName}`} />
                    <AppLayout
                        initializedByAnother={initializedByAnother}
                        mockApp={props.mockApp}
                        ref={appLayoutRef}
                        setInitialize={setInitialize}
                        suspectedUnauthorized={suspectedUnauthorized}
                        wellKnownConfig={wellKnownConfig}
                    />
                    <ToastContainer />
                    <BrowserNotificationsContainer />
                </BrowserRouter>
            </ThemeProvider>
        </AppContext.Provider>
    )
}

AppInner.displayName = 'AppInner'

export default AppInner
