import {
    forwardRef,
    memo,
    ReactElement,
    SyntheticEvent,
    useCallback,
    useContext,
    useEffect,
    useImperativeHandle,
    useMemo,
    useRef,
    useState,
} from 'react'
import { useIntl } from 'react-intl'
import { useNavigate, useLocation } from 'react-router-dom'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import Layout from '@shared-ui/components/Layout'
import Header from '@shared-ui/components/Layout/Header'
import LeftPanel, { parseActiveItem } from '@shared-ui/components/Layout/LeftPanel'
import VersionMark, { getVersionMarkData } from '@shared-ui/components/Atomic/VersionMark'
import UserWidget from '@/containers/App/UserWidget/UserWidget'
import Button from '@shared-ui/components/Atomic/Button'
import { MenuItem } from '@shared-ui/components/Layout/LeftPanel/LeftPanel.types'
import { getMinutesBetweenDates } from '@shared-ui/common/utils'
import { severities } from '@shared-ui/components/Atomic/VersionMark/constants'

import { DEVICE_AUTH_MODE, GITHUB_VERSION_REQUEST_INTERVAL } from '@/constants'
import AppAuthProvider from '@/containers/App/AppAuthProvider/AppAuthProvider'
import InitializedByAnother from '@/containers/App/AppInner/InitializedByAnother/InitializedByAnother'
import { mather, menu, Routes } from '@/routes'
import { getVersionNumberFromGithub, reset } from '@/containers/App/AppRest'
import { AppAuthProviderRefType } from '@/containers/App/AppAuthProvider/AppAuthProvider.types'
import { messages as t } from '../App.i18n'
import { AppLayoutRefType, Props } from './AppLayout.types'
import AppContext from '@/containers/App/AppContext'
import { useDispatch, useSelector } from 'react-redux'
import { CombinedStoreType } from '@/store/store'
import { setVersion } from '@/containers/App/slice'

const AppLayout = forwardRef<AppLayoutRefType, Props>((props, ref) => {
    const { mockApp, wellKnownConfig, setInitialize, initializedByAnother, suspectedUnauthorized } = props
    const { formatMessage: _ } = useIntl()
    const location = useLocation()
    const dispatch = useDispatch()
    const navigate = useNavigate()

    const [authError, setAuthError] = useState<string | undefined>(undefined)
    const [activeItem, setActiveItem] = useState(parseActiveItem(location.pathname, menu, mather))

    const authProviderRef = useRef<AppAuthProviderRefType | null>(null)

    const { collapsed, setCollapsed, iframeMode } = useContext(AppContext)

    const appStore = useSelector((state: CombinedStoreType) => state.app)

    useImperativeHandle(ref, () => ({
        getAuthProviderRef: () => authProviderRef.current,
    }))

    const requestVersion = useCallback((now: Date) => {
        getVersionNumberFromGithub().then((ret) => {
            dispatch(
                setVersion({
                    requestedDatetime: now,
                    latest: ret.data.tag_name.replace('v', ''),
                    latest_url: ret.data.html_url,
                })
            )
        })

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        const now: Date = new Date()

        if (
            !appStore.version.requestedDatetime ||
            getMinutesBetweenDates(new Date(appStore.version.requestedDatetime), now) > GITHUB_VERSION_REQUEST_INTERVAL
        ) {
            requestVersion(now)
        }

        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const versionMarkData = useMemo(
        () =>
            getVersionMarkData({
                buildVersion: wellKnownConfig?.version || '',
                githubVersion: appStore.version.latest || '',
                i18n: {
                    version: _(t.version),
                    newUpdateIsAvailable: _(t.newUpdateIsAvailable),
                },
            }),
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [appStore.version.latest, wellKnownConfig]
    )

    const handleItemClick = (item: MenuItem, e: SyntheticEvent) => {
        e.preventDefault()

        setActiveItem(item.id)
        item.link && navigate(item.link)
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
                isIframeMode={iframeMode}
                leftPanel={
                    <LeftPanel
                        activeId={activeItem}
                        collapsed={collapsed}
                        menu={menu}
                        onItemClick={handleItemClick}
                        setCollapsed={setCollapsed}
                        versionMark={
                            wellKnownConfig && (
                                <VersionMark
                                    severity={versionMarkData.severity}
                                    update={
                                        wellKnownConfig &&
                                        versionMarkData.severity !== severities.SUCCESS &&
                                        appStore.version.latest_url
                                            ? {
                                                  text: _(t.clickHere),
                                                  onClick: (e) => {
                                                      e.preventDefault()
                                                      window.open(appStore.version.latest_url, '_blank')
                                                  },
                                              }
                                            : undefined
                                    }
                                    versionText={versionMarkData.text}
                                />
                            )
                        }
                    />
                }
            />
        </ConditionalWrapper>
    )
})

AppLayout.displayName = 'AppLayout'

export default AppLayout
