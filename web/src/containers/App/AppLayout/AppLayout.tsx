import { FC, memo, ReactElement, SyntheticEvent, useContext, useRef, useState } from 'react'
import { useIntl } from 'react-intl'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import Layout from '@shared-ui/components/Layout'
import Header from '@shared-ui/components/Layout/Header'
import LeftPanel, { parseActiveItem } from '@shared-ui/components/Layout/LeftPanel'
import VersionMark from '@shared-ui/components/Atomic/VersionMark'
import { severities } from '@shared-ui/components/Atomic/VersionMark/constants'
import UserWidget from '@/containers/App/UserWidget/UserWidget'
import Button from '@shared-ui/components/Atomic/Button'
import { WellKnownConfigType } from '@shared-ui/common/hooks'
import { MenuItem } from '@shared-ui/components/Layout/LeftPanel/LeftPanel.types'

import { DEVICE_AUTH_MODE } from '@/constants'
import AppAuthProvider from '@/containers/App/AppAuthProvider/AppAuthProvider'
import InitializedByAnother from '@/containers/App/AppInner/InitializedByAnother/InitializedByAnother'
import { mather, menu, Routes } from '@/routes'
import { reset } from '@/containers/App/AppRest'
import { AppAuthProviderRefType } from '@/containers/App/AppAuthProvider/AppAuthProvider.types'
import { history } from '@/store'
import { messages as t } from '../App.i18n'
import AppContext from '@/containers/App/AppContext'

type Props = {
    initializedByAnother: boolean
    suspectedUnauthorized: boolean
    mockApp: boolean
    setInitialize: (isInitialize?: boolean) => void
    wellKnownConfig?: WellKnownConfigType
}

const AppLayout: FC<Props> = (props) => {
    const { mockApp, wellKnownConfig, setInitialize, initializedByAnother, suspectedUnauthorized } = props
    const { formatMessage: _ } = useIntl()

    const [authError, setAuthError] = useState<string | undefined>(undefined)
    const [activeItem, setActiveItem] = useState(parseActiveItem(history.location.pathname, menu, mather))

    const authProviderRef = useRef<AppAuthProviderRefType | null>(null)

    const { collapsed, setCollapsed } = useContext(AppContext)

    const handleItemClick = (item: MenuItem, e: SyntheticEvent) => {
        e.preventDefault()

        setActiveItem(item.id)
        history.push(item.link)
    }

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

    const UserWidgetComponent = memo(() => {
        if (
            !mockApp &&
            wellKnownConfig &&
            wellKnownConfig.remoteProvisioning &&
            wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
        ) {
            return <UserWidget logout={handleLogout} />
        }
        if (wellKnownConfig && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.PRE_SHARED_KEY) {
            return (
                <Button className='m-l-15' onClick={handleLogout}>
                    Logout
                </Button>
            )
        }

        return <div />
    })

    if (authError) {
        return <div className='client-error-message'>{`${_(t.authError)}: ${authError}`}</div>
    }

    return (
        <ConditionalWrapper
            condition={!props.mockApp && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509}
            wrapper={(child: ReactElement) => (
                <AppAuthProvider
                    ref={authProviderRef}
                    setAuthError={setAuthError}
                    setInitialize={setInitialize}
                    wellKnownConfig={wellKnownConfig}
                >
                    {child}
                </AppAuthProvider>
            )}
        >
            <Layout
                content={
                    <>
                        <InitializedByAnother logout={handleLogout} show={initializedByAnother} />
                        {!initializedByAnother && !suspectedUnauthorized && <Routes />}
                    </>
                }
                header={
                    <Header
                        breadcrumbs={<div id='breadcrumbsPortalTarget'></div>}
                        userWidget={<UserWidgetComponent />}
                    />
                }
                leftPanel={
                    <LeftPanel
                        activeId={activeItem}
                        collapsed={collapsed}
                        menu={menu}
                        onItemClick={handleItemClick}
                        setCollapsed={setCollapsed}
                        versionMark={<VersionMark severity={severities.SUCCESS} versionText='Version 2.02' />}
                    />
                }
            />
        </ConditionalWrapper>
    )
}

AppLayout.displayName = 'AppLayout'

export default AppLayout
