import { useAuth } from 'oidc-react'
import { security } from '@shared-ui/common/services'
import { forwardRef, useEffect, useImperativeHandle } from 'react'
import { REMOTE_PROVISIONING_MODE } from '@/constants'
import {
  getJwksData,
  getOpenIdConfiguration,
  initializeFinal,
  initializeJwksData,
  signIdentityCsr,
} from '@/containers/App/AppRest'
import AppLoader from '@/containers/App/AppLoader/AppLoader'

const AppAuthProvider = forwardRef((props: any, ref) => {
  const { wellKnownConfig, children, setAuthError, setInitialize } = props
  const {
    isLoading,
    userData,
    signOutRedirect: signOutMethod,
    userManager,
  } = useAuth()

  if (userData) {
    security.setAccessToken(userData.access_token)

    if (userManager) {
      security.setUserManager(userManager)
    }
  }

  useImperativeHandle(ref, () => ({
    getSignOutMethod() {
      return signOutMethod
    },
    getLoading() {
      return isLoading
    },
    getUserData() {
      return userData
    },
  }))

  useEffect(() => {
    if (
      !isLoading &&
      wellKnownConfig &&
      !wellKnownConfig.isInitialized &&
      wellKnownConfig.remoteProvisioning.mode ===
        REMOTE_PROVISIONING_MODE.USER_AGENT
    ) {
      try {
        getOpenIdConfiguration(
          wellKnownConfig.remoteProvisioning.authorization.authority
        ).then(result => {
          getJwksData(result.data.jwks_uri).then(result => {
            initializeJwksData(result.data).then(result => {
              const state = result.data.identityCertificateChallenge.state

              signIdentityCsr(
                wellKnownConfig?.remoteProvisioning.userAgent
                  .certificateAuthorityAddress,
                result.data.identityCertificateChallenge
                  .certificateSigningRequest
              ).then(result => {
                initializeFinal(state, result.data.certificate).then(() => {
                  setInitialize(true)
                })
              })
            })
          })
        })
      } catch (e) {
        console.error(e)
        setAuthError(e)
      }
    }
  }, [wellKnownConfig, isLoading, setAuthError, setInitialize])

  if (isLoading || !wellKnownConfig || !wellKnownConfig?.isInitialized) {
    return <AppLoader />
  }

  return children
})

export default AppAuthProvider