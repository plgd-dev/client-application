import React, {
    forwardRef,
    ReactElement,
    useCallback,
    useEffect,
    useImperativeHandle,
    useMemo,
    useRef,
    useState,
} from 'react'
import { useIntl } from 'react-intl'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { useTheme } from '@emotion/react'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import Layout from '@shared-ui/components/Layout'
import Header from '@shared-ui/components/Layout/Header'
import UserWidget from '@/containers/App/UserWidget/UserWidget'
import Button from '@shared-ui/components/Atomic/Button'
import { getMinutesBetweenDates } from '@shared-ui/common/utils'
import { getVersionNumberFromGithub, reset } from '@shared-ui/app/clientApp/App/AppRest'
import { DEVICE_AUTH_MODE, GITHUB_VERSION_REQUEST_INTERVAL } from '@shared-ui/app/clientApp/constants'
import AppAuthProvider from '@shared-ui/app/clientApp/App/AppAuthProvider'
import { AppAuthProviderRefType } from '@shared-ui/app/clientApp/App/AppAuthProvider/AppAuthProvider.types'
import { ThemeType } from '@shared-ui/components/Atomic/_theme'
import { hasDifferentOwner } from '@shared-ui/common/services/api-utils'
import { useAppInitialization } from '@shared-ui/app/clientApp/Devices/hooks'

import { Routes } from '@/routes'
import { messages as t } from '../App.i18n'
import { AppLayoutRefType, Props } from './AppLayout.types'
import { CombinedStoreType } from '@/store/store'
import { setVersion } from '@/containers/App/slice'
import AppConfig from '../AppConfig/AppConfig'

const LogoElement = (props: any) => {
    const { css, logo, className, onClick } = props
    return (
        <img
            alt=''
            className={className}
            css={css}
            height={logo.height}
            onClick={onClick}
            src={logo.source}
            width={logo.width}
        />
    )
}

const AppLayout = forwardRef<AppLayoutRefType, Props>((props, ref) => {
    const { mockApp, wellKnownConfig, setInitialize, initializedByAnother, suspectedUnauthorized } = props
    const { formatMessage: _ } = useIntl()
    const dispatch = useDispatch()
    const navigate = useNavigate()
    const theme: ThemeType = useTheme()

    const [authError, setAuthError] = useState<string | undefined>(undefined)

    const authProviderRef = useRef<AppAuthProviderRefType | null>(null)

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

    const handleLogout = () => {
        if (authProviderRef?.current) {
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
            }
        } else {
            // PSK mode
            reset().then(() => {
                setInitialize(false)
            })
        }
    }

    const getUserWidgetComponent = useCallback(() => {
        if (
            !mockApp &&
            !initializedByAnother &&
            wellKnownConfig &&
            wellKnownConfig.remoteProvisioning &&
            wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
        ) {
            return <UserWidget logoutTitle={_(t.logOut)} onLogout={handleLogout} />
        }
        if (wellKnownConfig && wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.PRE_SHARED_KEY) {
            return (
                <Button className='m-l-15' onClick={handleLogout}>
                    {_(t.logOut)}
                </Button>
            )
        }

        return <div />
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [handleLogout, initializedByAnother, mockApp, wellKnownConfig])

    const diffOwner = useMemo(
        () => hasDifferentOwner(wellKnownConfig, appStore.userWellKnownConfig, true),
        [wellKnownConfig, appStore.userWellKnownConfig]
    )

    // previous same owner, can be called reset
    const sameOwnerDiffAuth = useMemo(
        () => !hasDifferentOwner(wellKnownConfig, appStore.userWellKnownConfig),
        [wellKnownConfig, appStore.userWellKnownConfig]
    )

    const [initializationLoading, reInitializeLoading] = useAppInitialization({
        wellKnownConfig,
        loading: diffOwner && !sameOwnerDiffAuth,
        clientData: appStore.userWellKnownConfig,
        reInitialize: diffOwner && sameOwnerDiffAuth && !wellKnownConfig.isInitialized,
        changeInitialize: setInitialize,
    })

    if (authError) {
        return <div className='client-error-message'>{`${_(t.authError)}: ${authError}`}</div>
    }

    return (
        <ConditionalWrapper
            condition={
                !props.mockApp &&
                !initializedByAnother &&
                wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
            }
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
                    <Routes
                        initializedByAnother={
                            (initializedByAnother && !suspectedUnauthorized) || !wellKnownConfig.isInitialized
                        }
                        loading={initializationLoading || reInitializeLoading}
                    />
                }
                header={
                    <Header
                        breadcrumbs={<div id='breadcrumbsPortalTarget'></div>}
                        configButton={<AppConfig />}
                        contentLeft={theme.logo && <LogoElement logo={theme.logo} onClick={() => navigate(`/`)} />}
                        userWidget={getUserWidgetComponent()}
                    />
                }
            />
        </ConditionalWrapper>
    )
})

AppLayout.displayName = 'AppLayout'

export default AppLayout
