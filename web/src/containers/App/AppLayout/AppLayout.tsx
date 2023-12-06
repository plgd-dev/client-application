import React, { forwardRef, ReactElement, useCallback, useEffect, useImperativeHandle, useMemo, useRef } from 'react'
import { useIntl } from 'react-intl'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { useTheme } from '@emotion/react'
import isEmpty from 'lodash/isEmpty'
import isEqual from 'lodash/isEqual'

import ConditionalWrapper from '@shared-ui/components/Atomic/ConditionalWrapper'
import Layout from '@shared-ui/components/Layout'
import Header from '@shared-ui/components/Layout/Header'
import UserWidget from '@/containers/App/UserWidget/UserWidget'
import Button from '@shared-ui/components/Atomic/Button'
import { reset } from '@shared-ui/app/clientApp/App/AppRest'
import { DEVICE_AUTH_MODE } from '@shared-ui/app/clientApp/constants'
import AppAuthProvider from '@shared-ui/app/clientApp/App/AppAuthProvider'
import { AppAuthProviderRefType } from '@shared-ui/app/clientApp/App/AppAuthProvider/AppAuthProvider.types'
import { ThemeType } from '@shared-ui/components/Atomic/_theme'
import { hasDifferentOwner } from '@shared-ui/common/services/api-utils'
import { useAppInitialization } from '@shared-ui/app/clientApp/Devices/hooks'
import { useAppVersion } from '@shared-ui/common/hooks/use-app-version'
import Logo from '@shared-ui/components/Atomic/Logo'

import { Routes } from '@/routes'
import { messages as t } from '../App.i18n'
import { AppLayoutRefType, Props } from './AppLayout.types'
import { CombinedStoreType } from '@/store/store'
import { setVersion, storeUserWellKnownConfig } from '@/containers/App/slice'
import AppConfig from '../AppConfig/AppConfig'

const AppLayout = forwardRef<AppLayoutRefType, Props>((props, ref) => {
    const { mockApp, wellKnownConfig, initializedByAnother, suspectedUnauthorized, reFetchConfig } = props
    const { formatMessage: _ } = useIntl()
    const dispatch = useDispatch()
    const navigate = useNavigate()
    const theme: ThemeType = useTheme()

    const authProviderRef = useRef<AppAuthProviderRefType | null>(null)

    const appStore = useSelector((state: CombinedStoreType) => state.app)

    useImperativeHandle(ref, () => ({
        getAuthProviderRef: () => authProviderRef.current,
    }))

    const [version] = useAppVersion({ requestedDatetime: appStore.version.requestedDatetime })

    useEffect(() => {
        if (version && !isEqual(appStore.version, version)) {
            dispatch(setVersion(version))
        }
    }, [appStore.version, dispatch, version])

    const handleLogout = () => {
        if (authProviderRef?.current) {
            const signOut = authProviderRef?.current?.getSignOutMethod

            if (signOut) {
                if (!initializedByAnother) {
                    reset().then((_r) => {
                        reFetchConfig().then(() => {
                            dispatch(storeUserWellKnownConfig({}))
                            signOut().then((_r: void) => {})
                        })
                    })
                } else {
                    signOut().then()
                }
            }
        } else {
            // PSK mode
            reset().then(() => {
                reFetchConfig().then(() => {
                    dispatch(storeUserWellKnownConfig({}))
                })
            })
        }
    }

    const getUserWidgetComponent = useCallback(() => {
        if (isEmpty(appStore.userWellKnownConfig)) {
            return <div />
        }

        if (
            !mockApp &&
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
    }, [handleLogout, initializedByAnother, mockApp, wellKnownConfig, appStore.userWellKnownConfig])

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
        clientData: appStore.userWellKnownConfig,
        reInitialize: diffOwner && sameOwnerDiffAuth && !wellKnownConfig.isInitialized,
        reFetchConfig,
    })

    return (
        <ConditionalWrapper
            condition={
                !props.mockApp &&
                !initializedByAnother &&
                wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.X509
            }
            wrapper={(child: ReactElement) => <AppAuthProvider ref={authProviderRef}>{child}</AppAuthProvider>}
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
                        contentLeft={theme.logo && <Logo logo={theme.logo} onClick={() => navigate(`/`)} />}
                        userWidget={getUserWidgetComponent()}
                    />
                }
            />
        </ConditionalWrapper>
    )
})

AppLayout.displayName = 'AppLayout'

export default AppLayout
