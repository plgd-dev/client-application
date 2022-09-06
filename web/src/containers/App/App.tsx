import { useContext, useEffect, useState } from 'react'
import classNames from 'classnames'
import { Router } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import { Helmet } from 'react-helmet'
import appConfig from '@/config'
import {
  ToastContainer,
  BrowserNotificationsContainer,
} from '@shared-ui/components/new/Toast/Toast'
import LeftPanel from '@shared-ui/components/new/LeftPanel'
import Menu from '@shared-ui/components/new/Menu'
import StatusBar from '@shared-ui/components/new/StatusBar'
import Footer from '@shared-ui/components/new/Footer'
import { useLocalStorage } from '@shared-ui/common/hooks'
import { Routes } from '@/routes'
import { history } from '@/store/history'
import AppContext from './AppContext'
import { security } from '@shared-ui/common/services/security'
import './App.scss'
import { fetchApi } from '@shared-ui/common/services'
import { Props, BuildInformationType } from '@/containers/App/App.types'

const App = ({ config }: Props) => {
  const [collapsed, setCollapsed] = useLocalStorage('leftPanelCollapsed', true)
  const [buildInformation, setBuildInformation] =
    useState<BuildInformationType | null>(null)
  security.setGeneralConfig(config)

  useEffect(() => {
    try {
      fetchApi(`${config.httpGatewayAddress}/api/v1/information`).then(
        result => {
          setBuildInformation(result.data)
        }
      )
    } catch (e) {
      console.error(e)
    }
  }, []) // eslint-disable-line

  return (
    <AppContext.Provider
      value={{
        collapsed,
        ...config,
        buildInformation: buildInformation || null,
      }}
    >
      <Router history={history}>
        {/*<InitServices />*/}
        <Helmet
          defaultTitle={appConfig.appName}
          titleTemplate={`%s | ${appConfig.appName}`}
        />
        <Container fluid id="app" className={classNames({ collapsed })}>
          <StatusBar />
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
        <ToastContainer />
        <BrowserNotificationsContainer />
      </Router>
    </AppContext.Provider>
  )
}

export const useAppConfig = () => useContext(AppContext)

export default App
