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
import {
  BrowserNotificationsContainer,
  ToastContainer,
} from '@shared-ui/components/new'
import { useLocalStorage } from '@shared-ui/common/hooks'
import AppAuthProvider from '@/containers/App/AppAuth/AppAuthProvider'
import { ReactNode, useRef, useState } from 'react'
import ConditionalWrapper from '@shared-ui/components/new/ConditionalWrapper'
import { messages as t } from '../App.i18n'
import { security } from '@shared-ui/common/services'
import { Props } from './AppInner.types'
import { reset } from '@/containers/App/AppRest'

const AppInner = (props: Props) => {
  const { wellKnownConfig, configError, setInitialize } = props
  const buildInformation = {
    version: wellKnownConfig?.version,
    buildDate: wellKnownConfig?.buildDate,
    commitHash: wellKnownConfig?.commitHash,
    commitDate: wellKnownConfig?.commitDate,
    releaseUrl: wellKnownConfig?.releaseUrl,
  }
  const [authError, setAuthError] = useState<string | undefined>(undefined)
  const [collapsed, setCollapsed] = useLocalStorage('leftPanelCollapsed', true)
  const { formatMessage: _ } = useIntl()
  const authProviderRef = useRef<any>(undefined)

  if (wellKnownConfig) {
    security.setWebOAuthConfig({
      authority: wellKnownConfig.remoteProvisioning.authorization.authority,
      clientId: wellKnownConfig.remoteProvisioning.authorization.clientId,
      redirect_uri: window.location.origin,
    })
  }

  const AppLayout = () => {
    if (
      wellKnownConfig?.deviceAuthenticationMode ===
      DEVICE_AUTH_MODE.PRE_SHARED_KEY
    ) {
      return <PreSharedKeySetup />
    }

    return (
      <ConditionalWrapper
        condition={
          wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
        }
        wrapper={(child: ReactNode) => (
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
        <Container fluid id="app" className={classNames({ collapsed })}>
          <StatusBar>
            <button
              className="m-l-10"
              onClick={() => {
                if (authProviderRef) {
                  const signOutMethod =
                    authProviderRef?.current?.getSignOutMethod()

                  reset().then(_r => {
                    signOutMethod().then((r: void) => {
                      setInitialize(false)
                    })
                  })
                }
              }}
            >
              signOut
            </button>
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
              collapsed={collapsed}
              toggleCollapsed={() => setCollapsed(!collapsed)}
            />
          </LeftPanel>
          <div id="content">
            <Routes />
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
    return (
      <div className="client-error-message">
        {`${_(t.authError)}: ${configError?.message || authError}`}
      </div>
    )
  }

  return (
    <AppContext.Provider
      value={{
        collapsed,
        buildInformation: buildInformation || undefined,
      }}
    >
      <Router history={history}>
        <Helmet
          defaultTitle={appConfig.appName}
          titleTemplate={`%s | ${appConfig.appName}`}
        />
        <AppLayout />
        <ToastContainer />
        <BrowserNotificationsContainer />
      </Router>
    </AppContext.Provider>
  )
}

AppInner.displayName = 'AppInner'

export default AppInner
