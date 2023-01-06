import { fetchApi, security } from '@shared-ui/common/services'
import { devicesApiEndpoints } from './constants'
import { interfaceGetParam } from './utils'
import { signIdentityCsr } from '@/containers/App/AppRest'
import { WellKnownConfigType } from '@shared-ui/common/hooks'
import { DEVICE_AUTH_MODE } from '@/constants'

type SecurityConfig = {
    httpGatewayAddress: string
    authority: string
}

const getConfig = () => security.getGeneralConfig() as SecurityConfig

/**
 * Get a single thing by its ID Rest Api endpoint
 */
export const getDeviceApi = (deviceId: string) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}`)

/**
 * Delete a set of devices by their IDs Rest Api endpoint
 */
export const deleteDevicesApi = () =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}`, {
        method: 'DELETE',
    })

/**
 * Get devices RESOURCES Rest Api endpoint
 * @param deviceId
 */
export const getDevicesResourcesAllApi = (deviceId: string) =>
    fetchApi(
        `${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/${
            devicesApiEndpoints.DEVICES_RESOURCES_SUFFIX
        }`
    )

/**
 * Get devices RESOURCES Rest Api endpoint
 */
export const getDevicesResourcesApi = ({
    deviceId,
    href,
    currentInterface = '',
}: {
    deviceId: string
    href: string
    currentInterface?: string
}) =>
    fetchApi(
        `${getConfig().httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
        }/${deviceId}/resources${href}${interfaceGetParam(currentInterface)}`
    )

/**
 * Update devices RESOURCE Rest Api endpoint
 */
export const updateDevicesResourceApi = (
    {
        deviceId,
        href,
        currentInterface = '',
    }: {
        deviceId: string
        href: string
        currentInterface?: string
    },
    data: any
) =>
    fetchApi(
        `${getConfig().httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
        }/${deviceId}/resources${href}${interfaceGetParam(currentInterface)}`,
        { method: 'PUT', body: data }
    )

/**
 * Create devices RESOURCE Rest Api endpoint
 */
export const createDevicesResourceApi = (
    {
        deviceId,
        href,
        currentInterface = '',
    }: {
        deviceId: string
        href: string
        currentInterface?: string
    },
    data: any
) =>
    fetchApi(
        `${getConfig().httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
        }/${deviceId}/resource-links${href}${interfaceGetParam(currentInterface)}`,
        { method: 'POST', body: data }
    )

/**
 * Delete devices RESOURCE Rest Api endpoint
 */
export const deleteDevicesResourceApi = ({ deviceId, href }: { deviceId: string; href: string }) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/resource-links${href}`, {
        method: 'DELETE',
    })

/**
 * Add device by IP
 */
export const addDeviceByIp = (deviceIp: string) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}?useEndpoints=${deviceIp}`)

/**
 * Own device by deviceId
 */
export const ownDeviceApi = (deviceId: string) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/own`, {
        method: 'POST',
    }).then((result) => {
        if (result?.data?.identityCertificateChallenge) {
            const state = result.data.identityCertificateChallenge.state
            //owning with csr
            // @ts-ignore
            const { certificateAuthority } = security.getWebOAuthConfig()
            signIdentityCsr(
                certificateAuthority,
                result.data.identityCertificateChallenge.certificateSigningRequest
            ).then((result) => {
                fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/own/${state}`, {
                    method: 'POST',
                    body: {
                        certificate: result.data.certificate,
                    },
                }).then((r) => r)
            })
        } else {
            return result
        }
    })

/**
 * DisOwn device by deviceId
 */
export const disownDeviceApi = (deviceId: string) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/disown`, { method: 'POST' })

export type OnboardDataType = {
    coapGatewayAddress: string
    authorizationCode: string
    authorizationProviderName: string
    hubId: string
    certificateAuthorities: string
}

export const onboardDeviceApi = (deviceId: string, data: OnboardDataType) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/onboard`, {
        method: 'POST',
        body: data,
    })

export const offboardDeviceApi = (deviceId: string) =>
    fetchApi(`${getConfig().httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/offboard`, {
        method: 'POST',
    })

export const DEVICE_AUTH_CODE_SESSION_KEY = 'tempDeviceAuthCode'

/**
 * Returns an async function which resolves with a authorization code gathered from a rendered iframe, used for onboarding of a device.
 * @param {*} deviceId
 */
export const getDeviceAuthCode = (deviceId: string) => {
    return new Promise((resolve, reject) => {
        const wellKnownConfig = security.getWellKnowConfig() as WellKnownConfigType

        if (!wellKnownConfig.remoteProvisioning) {
            return reject(new Error('remoteProvisioning is missing in wellKnowConfig'))
        }

        const { clientId, scopes = [], audience } = wellKnownConfig.remoteProvisioning.deviceOauthClient
        const IS_PRE_SHARED_KEY_MOD = wellKnownConfig?.deviceAuthenticationMode === DEVICE_AUTH_MODE.PRE_SHARED_KEY

        const AuthUserManager = security.getUserManager()

        AuthUserManager.metadataService.getAuthorizationEndpoint().then((authorizationEndpoint: string) => {
            let timeout: any = null

            const iframe = document.createElement('iframe')
            const audienceParam = audience ? `&audience=${audience}` : ''
            iframe.src = `${authorizationEndpoint}?response_type=code&client_id=${clientId}&scope=${scopes}${audienceParam}&redirect_uri=${window.location.origin}/devices&device_id=${deviceId}`
            iframe.className = IS_PRE_SHARED_KEY_MOD ? 'iframeAuthModalVisible' : 'iframeAuthModal'

            const closeWrapper = document.createElement('div')
            closeWrapper.className = 'iframeAuthModalClose'

            const closeElement = document.createElement('a')
            closeElement.innerHTML =
                '<svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" role="img"><path d="M16 29.333c7.333 0 13.333-6 13.333-13.334 0-7.333-6-13.333-13.333-13.333s-13.333 6-13.333 13.333c0 7.334 6 13.334 13.333 13.334ZM12.227 19.773l7.546-7.546M19.773 19.773l-7.546-7.546" stroke="currentcolor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>'
            closeElement.onclick = () => {
                destroyIframe()
                reject('user-cancel')
            }

            closeWrapper.appendChild(closeElement)

            const overlayElement = document.createElement('div')
            overlayElement.className = 'iframeAuthModalOverlay'

            const destroyIframe = () => {
                overlayElement.remove()
                closeWrapper.remove()
                sessionStorage.removeItem(DEVICE_AUTH_CODE_SESSION_KEY)
                iframe?.parentNode?.removeChild(iframe)
            }

            const doResolve = (value: any) => {
                destroyIframe()
                clearTimeout(timeout)
                resolve(value)
            }

            const doReject = () => {
                destroyIframe()
                clearTimeout(timeout)
                reject(new Error('Failed to get the device auth code.'))
            }

            iframe.onload = () => {
                let attempts = 0
                const maxAttempts = 40

                if (IS_PRE_SHARED_KEY_MOD) {
                    document.body.appendChild(overlayElement)
                    document.body.appendChild(closeWrapper)
                    iframe.className = 'iframeAuthModalVisible iframeAuthModalVisibleShadow'
                }

                const getCode = () => {
                    attempts += 1
                    const code = sessionStorage.getItem(DEVICE_AUTH_CODE_SESSION_KEY)

                    if (code) {
                        return doResolve(code)
                    }

                    if (attempts > maxAttempts && !IS_PRE_SHARED_KEY_MOD) {
                        return doReject()
                    }

                    timeout = setTimeout(getCode, 500)
                }

                getCode()
            }

            iframe.onerror = () => {
                doReject()
            }

            document.body.appendChild(iframe)
        })
    })
}
