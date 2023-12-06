import { useRef, useState, useMemo, useCallback, useEffect } from 'react'
import { Helmet } from 'react-helmet'
import { BrowserRouter } from 'react-router-dom'
import { useSelector } from 'react-redux'
import isEmpty from 'lodash/isEmpty'

import { BrowserNotificationsContainer } from '@shared-ui/components/Atomic/Toast'
import { useLocalStorage, WellKnownConfigType } from '@shared-ui/common/hooks'
import { security } from '@shared-ui/common/services'
import AppContext from '@shared-ui/app/share/AppContext'
import { DEVICE_AUTH_MODE } from '@shared-ui/app/clientApp/constants'
import { hasDifferentOwner } from '@shared-ui/common/services/api-utils'
import App from '@shared-ui/components/Atomic/App/App'

import appConfig from '@/config'
import { Props } from './AppInner.types'
import AppLayout from '@/containers/App/AppLayout/AppLayout'
import { AppLayoutRefType } from '@/containers/App/AppLayout/AppLayout.types'
import { storeUserWellKnownConfig } from '@/containers/App/slice'
import { CombinedStoreType } from '@/store/store'

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

    const appStore = useSelector((state: CombinedStoreType) => state.app)

    const differentOwner = useCallback(
        (wellKnownConfig: WellKnownConfigType, userWellKnownConfig: any) =>
            hasDifferentOwner(wellKnownConfig, userWellKnownConfig, true),
        []
    )

    const unauthorizedCallback = useCallback(() => {
        setSuspectedUnauthorized(true)

        reFetchConfig()
            .then((newWellKnownConfig: WellKnownConfigType) => {
                if (differentOwner(newWellKnownConfig, appStore.userWellKnownConfig)) {
                    setInitializedByAnother(true)
                }
            })
            .then(() => {
                setSuspectedUnauthorized(false)
            })
    }, [differentOwner, reFetchConfig, appStore.userWellKnownConfig])

    // check on load
    useEffect(() => {
        if (!isEmpty(appStore.userWellKnownConfig)) {
            unauthorizedCallback()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [appStore.userWellKnownConfig])

    const contextValue = useMemo(
        () => ({
            unauthorizedCallback,
            collapsed,
            setCollapsed,
            buildInformation: buildInformation || undefined,
            isHub: false,
            updateAppWellKnownConfig: storeUserWellKnownConfig,
            reFetchConfig,
        }),
        [buildInformation, collapsed, setCollapsed, unauthorizedCallback, reFetchConfig]
    )

    // Render an error box with a config error
    if (configError) {
        return <div className='client-error-message'>{configError?.message}</div>
    }

    return (
        <AppContext.Provider value={contextValue}>
            <BrowserRouter>
                <Helmet defaultTitle={appConfig.appName} titleTemplate={`%s | ${appConfig.appName}`} />
                <App toastContainerPortalTarget={document.getElementById('toast-root')}>
                    <AppLayout
                        initializedByAnother={!!initializedByAnother}
                        mockApp={props.mockApp}
                        reFetchConfig={reFetchConfig}
                        ref={appLayoutRef}
                        suspectedUnauthorized={suspectedUnauthorized}
                        updateWellKnownConfig={updateWellKnownConfig}
                        wellKnownConfig={wellKnownConfig}
                    />
                </App>
                <BrowserNotificationsContainer />
            </BrowserRouter>
        </AppContext.Provider>
    )
}

AppInner.displayName = 'AppInner'

export default AppInner
