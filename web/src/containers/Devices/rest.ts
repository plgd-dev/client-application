import { fetchApi, security } from '@shared-ui/common/services'
import { devicesApiEndpoints } from './constants'
import { interfaceGetParam } from './utils'
import { signIdentityCsr } from '@/containers/App/AppRest'
import { WellKnownConfigType } from '@shared-ui/common/hooks'

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
        const wellKnowConfig = security.getWellKnowConfig() as WellKnownConfigType

        console.log('getDeviceAuthCode!!')

        if (!wellKnowConfig.remoteProvisioning) {
            return reject(new Error('remoteProvisioning is missing in wellKnowConfig'))
        }

        console.log({ wellKnowConfig })

        const { clientId, scopes = [] } = wellKnowConfig.remoteProvisioning.deviceOauthClient
        const { audience } = wellKnowConfig.remoteProvisioning.webOauthClient

        const AuthUserManager = security.getUserManager()

        // AuthUserManager.signinPopup().then((u: User | null) => {
        //     console.log(u)
        //     console.log('signinRedirect done')
        //

        //
        // window.localStorage.setItem(
        //     `oidc.${state}`,
        //     JSON.stringify({
        //         authority: 'https://auth.plgd.cloud/realms/shared',
        //         client_id: 'LXZ9OhKWWRYqf12W0B5OXduqt02q0zjS',
        //         code_verifier:
        //             '5901ae2aa82942a888cfe35ddaf3923d350440855a584ba79e134a1407eb3d62b53d431e609f4782b00b1ec0050b5306',
        //         created: 1671655302,
        //         extraTokenParams: { plgd: 1 },
        //         id: state,
        //         redirect_uri: 'http://localhost:3000',
        //         request_type: 'si:s',
        //         response_mode: 'query',
        //         scope: 'openid',
        //     })
        // )

        AuthUserManager.metadataService.getAuthorizationEndpoint().then((authorizationEndpoint: string) => {
            let timeout: any = null

            const iframe = document.createElement('iframe')
            const audienceParam = audience ? `&audience=${audience}` : ''
            iframe.src = `${authorizationEndpoint}?response_type=code&client_id=${clientId}&scope=${scopes}${audienceParam}&redirect_uri=${window.location.origin}/devices&device_id=${deviceId}`
            iframe.className = 'iframeTestModal'

            const destroyIframe = () => {
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

                const getCode = () => {
                    attempts += 1
                    const code = sessionStorage.getItem(DEVICE_AUTH_CODE_SESSION_KEY)

                    if (code) {
                        return doResolve(code)
                    }

                    if (attempts > maxAttempts) {
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
