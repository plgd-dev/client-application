import { fetchApi, security } from '@shared-ui/common/services'
import { devicesApiEndpoints } from './constants'
import { interfaceGetParam } from './utils'
import { signIdentityCsr } from '@/containers/App/AppRest'
import { WellKnownConfigType } from '@shared-ui/common/hooks'
import { DEVICE_AUTH_MODE, DEVICE_AUTH_CODE_SESSION_KEY } from '@/constants'

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

export const PLGD_BROWSER_USED = 'plgdBrowserUsed'

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
            const audienceParam = audience ? `&audience=${audience}` : ''

            const win = window.open(
                `${authorizationEndpoint}?response_type=code&client_id=${clientId}&scope=${scopes}${audienceParam}&redirect_uri=${window.location.origin}/devices&device_id=${deviceId}`,
                'thePopUp',
                ''
            )
            const pollTimer = window.setInterval(function () {
                if (win && win.closed) {
                    window.clearInterval(pollTimer)
                    // find code after close
                    const code = localStorage.getItem(DEVICE_AUTH_CODE_SESSION_KEY)

                    if (code) {
                        localStorage.removeItem(DEVICE_AUTH_CODE_SESSION_KEY)
                        return doResolve(code)
                    } else {
                        return reject('user-cancel')
                    }
                }
            }, 200)

            const doResolve = (value: any) => {
                clearTimeout(timeout)
                resolve(value)
            }

            const doReject = () => {
                clearTimeout(timeout)

                reject(new Error('Failed to get the device auth code.'))
            }

            let attempts = 0
            const maxAttempts = 60

            const getCode = () => {
                attempts += 1
                const code = localStorage.getItem(DEVICE_AUTH_CODE_SESSION_KEY)

                if (code) {
                    localStorage.removeItem(DEVICE_AUTH_CODE_SESSION_KEY)
                    return doResolve(code)
                }

                if (attempts > maxAttempts && !IS_PRE_SHARED_KEY_MOD) {
                    return doReject()
                }

                timeout = setTimeout(getCode, 500)
            }

            // scan for the code
            getCode()
        })
    })
}
