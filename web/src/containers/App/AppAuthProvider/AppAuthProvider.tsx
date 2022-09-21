import {useAuth} from 'oidc-react'
import {security} from '@shared-ui/common/services'
import {forwardRef, useEffect, useImperativeHandle} from 'react'
import {REMOTE_PROVISIONING_MODE} from '@/constants'
import {
    getJwksData,
    getOpenIdConfiguration,
    initializeFinal,
    initializeJwksData,
    signIdentityCsr,
} from '@/containers/App/AppRest'
import AppLoader from '@/containers/App/AppLoader/AppLoader'
import {AppAuthProviderRefType, Props} from './AppAuthProvider.types'

const AppAuthProvider = forwardRef<AppAuthProviderRefType, Props>((props, ref) => {
    const {wellKnownConfig, children, setAuthError, setInitialize} = props
    const {isLoading, userData, signOutRedirect, userManager} = useAuth()

    if (userData) {
        security.setAccessToken(userData.access_token)

        if (userManager) {
            security.setUserManager(userManager)
        }
    }

    useImperativeHandle(ref, () => ({
        getSignOutMethod: () => signOutRedirect({
            post_logout_redirect_uri: window.location.origin,
        }),
    }))

    useEffect(() => {
        if (
            !isLoading &&
            wellKnownConfig &&
            !wellKnownConfig.isInitialized &&
            wellKnownConfig.remoteProvisioning.mode === REMOTE_PROVISIONING_MODE.USER_AGENT
        ) {
            try {
                getOpenIdConfiguration(wellKnownConfig.remoteProvisioning.authorization.authority).then((result) => {
                    getJwksData(result.data.jwks_uri).then((result) => {
                        initializeJwksData(result.data).then((result) => {
                            const state = result.data.identityCertificateChallenge.state

                            signIdentityCsr(
                                wellKnownConfig?.remoteProvisioning.userAgent.certificateAuthorityAddress,
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
        return <AppLoader/>
    }

    return children
})

AppAuthProvider.displayName = 'AppAuthProvider'

export default AppAuthProvider
