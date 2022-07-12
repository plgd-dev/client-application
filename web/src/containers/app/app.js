import { useContext, useEffect, useState } from 'react'
import classNames from 'classnames'
import { Router } from 'react-router-dom'
import Container from 'react-bootstrap/Container'
import { Helmet } from 'react-helmet'
import appConfig from '@/config'
import {
  ToastContainer,
  BrowserNotificationsContainer,
} from '@/components/toast'
import { LeftPanel } from '@/components/left-panel'
import { Menu } from '@/components/menu'
import { StatusBar } from '@/components/status-bar'
import { Footer } from '@/components/footer'
import { useLocalStorage } from '@/common/hooks'
import { Routes } from '@/routes'
import { history } from '@/store/history'
import { AppContext } from './app-context'
import { security } from '@/common/services/security'
import './app.scss'
import { fetchApi } from '@/common/services'

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
