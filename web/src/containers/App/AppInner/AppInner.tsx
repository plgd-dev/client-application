import { useIntl } from 'react-intl'
import { DEVICE_AUTH_MODE } from '@/constants'
import PreSharedKeySetup from '@/containers/App/PreSharedKeySetup/PreSharedKeySetup'
import Container from 'react-bootstrap/Container'
import classNames from 'classnames'
import StatusBar from '@shared-ui/components/new/StatusBar'
import LeftPanel from '@shared-ui/components/new/LeftPanel'
import Menu from '@shared-ui/components/new/Menu'
import { Routes } from '@/routes'
import Footer from '@shared-ui/components/new/Footer'
import AppContext from '@/containers/App/AppContext'
import { Router } from 'react-router-dom'
import { history } from '@/store'
import { Helmet } from 'react-helmet'
import appConfig from '@/config'
import { BrowserNotificationsContainer, Button, ToastContainer } from '@shared-ui/components/new'
import { useLocalStorage, WellKnownConfigType } from '@shared-ui/common/hooks'
import AppAuthProvider from '@/containers/App/AppAuthProvider/AppAuthProvider'
import { ReactElement, useRef, useState } from 'react'
import ConditionalWrapper from '@shared-ui/components/new/ConditionalWrapper'
import { messages as t } from '../App.i18n'
import { security } from '@shared-ui/common/services'
import { Props } from './AppInner.types'
import { reset } from '@/containers/App/AppRest'
import UserWidget from '@shared-ui/components/new/UserWidget'
import { AppAuthProviderRefType } from '@/containers/App/AppAuthProvider/AppAuthProvider.types'
import InitializedByAnother from '@/containers/App/AppInner/InitializedByAnother/InitializedByAnother'
import { User } from 'oidc-react'
import jwtDecode from 'jwt-decode'
import get from 'lodash/get'

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
    const [authError, setAuthError] = useState<string | undefined>(undefined)
    const [collapsed, setCollapsed] = useLocalStorage('leftPanelCollapsed', true)
    const { formatMessage: _ } = useIntl()
    const authProviderRef = useRef<AppAuthProviderRefType | null>(null)

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

    const unauthorizedCallback = () => {
        setSuspectedUnauthorized(true)

        reFetchConfig().then((newWellKnownConfig: WellKnownConfigType) => {
            if (authProviderRef) {
                const userData: User = authProviderRef?.current?.getUserData()
                const parsedData = jwtDecode(userData.access_token)
                const ownerId = get(parsedData, newWellKnownConfig.remoteProvisioning?.jwtOwnerClaim as string, '')

                if (ownerId !== newWellKnownConfig?.owner) {
                    setInitializedByAnother(true)
                }
            }

            setSuspectedUnauthorized(false)
        })
    }

    const AppLayout = () => {
        const handleLogout = () => {
            if (authProviderRef) {
                const signOut = authProviderRef?.current?.getSignOutMethod

                if (signOut) {
                    if (!initializedByAnother) {
                        reset().then((_r) => {
                            signOut().then((_r: void) => {
                                setInitialize(false)
                            })
                        })
                    } else {
                        // s remoteProvisioning vsetko nad
                        // bez remoteProvisioning
                        signOut().then()
                    }
                } else {
                    // preshared mode
                    reset().then(() => {
                        setInitialize(false)
                    })
                }
            }
        }

        if (
            !wellKnownConfig?.isInitialized &&
            wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.PRE_SHARED_KEY
        ) {
            return <PreSharedKeySetup setInitialize={setInitialize} />
        }

        return (
            <ConditionalWrapper
                condition={!props.mockApp && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
                wrapper={(child: ReactElement) => (
                    <AppAuthProvider
                        wellKnownConfig={wellKnownConfig}
                        setAuthError={setAuthError}
                        setInitialize={setInitialize}
                        ref={authProviderRef}
                    >
                        {child}
                    </AppAuthProvider>
                )}
            >
                <Container fluid id='app' className={classNames({ collapsed })}>
                    <StatusBar>
                        {!props.mockApp &&
                            wellKnownConfig &&
                            wellKnownConfig.remoteProvisioning &&
                            wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509 && (
                                <UserWidget logout={handleLogout} />
                            )}
                        {wellKnownConfig &&
                            wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.PRE_SHARED_KEY && (
                                <Button className='m-l-15' onClick={handleLogout}>
                                    Logout
                                </Button>
                            )}
                    </StatusBar>
                    <LeftPanel>
                        <Menu
                            menuItems={[
                                {
                                    to: '/',
                                    icon: 'fa-list',
                                    nameKey: 'devices',
                                    className: 'devices',
                                },
                            ]}
                            collapsed={!!collapsed}
                            toggleCollapsed={() => setCollapsed(!collapsed)}
                            initializedByAnother={initializedByAnother}
                        />
                    </LeftPanel>
                    <div id='content'>
                        <InitializedByAnother show={initializedByAnother} logout={handleLogout} />
                        {!initializedByAnother && !suspectedUnauthorized && <Routes />}
                        <Footer
                            links={[
                                {
                                    to: 'https://github.com/plgd-dev/client-application/blob/main/pb/service.swagger.json',
                                    i18key: 'API',
                                },
                                {
                                    to: 'https://docs.plgd.dev/',
                                    i18key: 'docs',
                                },
                                {
                                    to: 'https://discord.gg/Pcusx938kg',
                                    i18key: 'contribute',
                                },
                            ]}
                        />
                    </div>
                </Container>
            </ConditionalWrapper>
        )
    }

    // Render an error box with a config error
    if (configError || authError) {
        return <div className='client-error-message'>{`${_(t.authError)}: ${configError?.message || authError}`}</div>
    }

    return (
        <AppContext.Provider
            value={{
                collapsed,
                unauthorizedCallback,
                buildInformation: buildInformation || undefined,
            }}
        >
            <Router history={history}>
                <Helmet defaultTitle={appConfig.appName} titleTemplate={`%s | ${appConfig.appName}`} />
                <AppLayout />
                <ToastContainer />
                <BrowserNotificationsContainer />
            </Router>
        </AppContext.Provider>
    )
}

AppInner.displayName = 'AppInner'

export default AppInner
