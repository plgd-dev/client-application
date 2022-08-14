import { useContext, useEffect, useState } from 'react'
import classNames from 'classnames'
import { Router } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import { Helmet } from 'react-helmet'
import appConfig from '@/config'
import {
  ToastContainer,
  BrowserNotificationsContainer,
} from '@shared-ui/components/old/toast'
import { LeftPanel } from '@shared-ui/components/old/left-panel'
import { Menu } from '@shared-ui/components/old/menu'
import { StatusBar } from '@shared-ui/components/old/status-bar'
import { Footer } from '@shared-ui/components/old/footer'
import { useLocalStorage } from '@shared-ui/common/hooks'
import { Routes } from '@/routes'
import { history } from '@/store/history'
import { AppContext } from './app-context'
import { security } from '@shared-ui/common/services/security'
import './app.scss'
import { fetchApi } from '@shared-ui/common/services'

const menuItems = [
  {
    to: '/',
    icon: 'fa-list',
    nameKey: 'devices',
    className: 'devices',
  },
]

const App = ({ config }) => {
  const [collapsed, setCollapsed] = useLocalStorage('leftPanelCollapsed', true)
  const [buildInformation, setBuildInformation] = useState(undefined)
  security.setGeneralConfig(config)

  useEffect(() => {
    fetchApi(`${config.httpGatewayAddress}/api/v1/information`).then(result => {
      setBuildInformation(result.data)
    })
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
              menuItems={{ menuItems }}
              collapsed={collapsed}
              toggleCollapsed={() => setCollapsed(!collapsed)}
            />
          </LeftPanel>
          <div id="content">
            <Routes />
            <Footer />
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
