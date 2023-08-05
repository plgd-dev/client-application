import { useAuth, User } from 'oidc-react'
import { forwardRef, useEffect, useImperativeHandle, useRef } from 'react'

import { security } from '@shared-ui/common/services'
import {
    getJwksData,
    getOpenIdConfiguration,
    initializeFinal,
    initializeJwksData,
    signIdentityCsr,
} from '@shared-ui/app/clientApp/App/AppRest'
import { REMOTE_PROVISIONING_MODE } from '@shared-ui/app/clientApp/constants'

import AppLoader from '@/containers/App/AppLoader/AppLoader'
import { AppAuthProviderRefType, Props } from './AppAuthProvider.types'

const AppAuthProvider = forwardRef<AppAuthProviderRefType, Props>((props, ref) => {
    const { wellKnownConfig, children, setAuthError, setInitialize } = props
    const { isLoading, userData, signOutRedirect } = useAuth()
    const userDataRef = useRef<User | null>(null)

    if (userData) {
        security.setAccessToken(userData.access_token)
        userDataRef.current = userData
    }

    useImperativeHandle(ref, () => ({
        getSignOutMethod: () =>
            signOutRedirect({
                post_logout_redirect_uri: window.location.origin,
            }),
        getUserData: () => userDataRef.current,
    }))

    useEffect(() => {
        if (
            !isLoading &&
            wellKnownConfig &&
            !wellKnownConfig.isInitialized &&
            wellKnownConfig.remoteProvisioning?.mode === REMOTE_PROVISIONING_MODE.USER_AGENT
        ) {
            try {
                getOpenIdConfiguration(wellKnownConfig.remoteProvisioning?.authority).then((result) => {
                    getJwksData(result.data.jwks_uri).then((result) => {
                        initializeJwksData(result.data).then((result) => {
                            const state = result.data.identityCertificateChallenge.state

                            signIdentityCsr(
                                wellKnownConfig.remoteProvisioning?.certificateAuthority as string,
                                result.data.identityCertificateChallenge.certificateSigningRequest
                            ).then((result) => {
                                initializeFinal(state, result.data.certificate).then(() => {
                                    setInitialize(true)
                                })
                            })
                        })
                    })
                })
            } catch (e) {
                console.error(e)
                setAuthError(e as string)
            }
        }
    }, [wellKnownConfig, isLoading, setAuthError, setInitialize])

    if (isLoading || !wellKnownConfig || !wellKnownConfig?.isInitialized) {
        return <AppLoader />
    }

    return children
})

AppAuthProvider.displayName = 'AppAuthProvider'

export default AppAuthProvider
