import {fetchApi, security} from '@shared-ui/common/services'
import {DEVICE_AUTH_CODE_SESSION_KEY} from '@/constants'
import {devicesApiEndpoints} from './constants'
import {interfaceGetParam} from './utils'

type SecurityConfig = {
    httpGatewayAddress: string
    authority: string
}

const config: SecurityConfig = security.getGeneralConfig() as SecurityConfig

/**
 * Get a single thing by its ID Rest Api endpoint
 */
export const getDeviceApi = (deviceId: string) =>
    fetchApi(
        `${config.httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
        }/${deviceId}`
    )

/**
 * Delete a set of devices by their IDs Rest Api endpoint
 */
export const deleteDevicesApi = () => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }`,
    {
        method: 'DELETE',
    }
)


/**
 * Get devices RESOURCES Rest Api endpoint
 * @param deviceId
 */
export const getDevicesResourcesAllApi = (deviceId: string) =>
    fetchApi(
        `${config.httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
        }/${deviceId}/${devicesApiEndpoints.DEVICES_RESOURCES_SUFFIX}`
    )

/**
 * Get devices RESOURCES Rest Api endpoint
 */
export const getDevicesResourcesApi = ({
                                           deviceId,
                                           href,
                                           currentInterface = '',
                                       }: { deviceId: string; href: string; currentInterface?: string }) =>
    fetchApi(
        `${config.httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
        }/${deviceId}/resources${href}?${interfaceGetParam(currentInterface)}`
    )

/**
 * Update devices RESOURCE Rest Api endpoint
 */
export const updateDevicesResourceApi = (
    {
        deviceId,
        href,
        currentInterface = '',
        ttl
    }: { deviceId: string, href: string, currentInterface?: string; ttl: string | number },
    data: any
) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }/${deviceId}/resources${href}?timeToLive=${ttl}&${interfaceGetParam(
        currentInterface
    )}`,
    {method: 'PUT', body: data, timeToLive: ttl}
)


/**
 * Create devices RESOURCE Rest Api endpoint
 */
export const createDevicesResourceApi = (
    {
        deviceId,
        href,
        currentInterface = '',
        ttl
    }: { deviceId: string, href: string, currentInterface?: string; ttl: string | number },
    data: any
) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }/${deviceId}/resource-links${href}?timeToLive=${ttl}&${interfaceGetParam(
        currentInterface
    )}`,
    {method: 'POST', body: data, timeToLive: ttl}
)


/**
 * Delete devices RESOURCE Rest Api endpoint
 * @param {*} params { deviceId, href - resource href}
 */
export const deleteDevicesResourceApi = ({deviceId, href}: { deviceId: string; href: string }) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }/${deviceId}/resource-links${href}`,
    {method: 'DELETE'}
)


/**
 * Update the shadowSynchronization of one Thing Rest Api endpoint
 */
export const updateDeviceShadowSynchronizationApi = (
    deviceId: string,
    shadowSynchronization: any
) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }/${deviceId}/metadata`,
    {method: 'PUT', body: {shadowSynchronization}}
)


/**
 * Returns an async function which resolves with a authorization code gathered from a rendered iframe, used for onboarding of a device.
 * @param {*} deviceId
 */
export const getDeviceAuthCode = (deviceId: string) => new Promise((resolve, reject) => {
    const {authority} = config
    const {clientId, audience, scopes = []} = security.getDeviceOAuthConfig() as any

    if (!clientId) {
        return reject(
            new Error(
                'clientId is missing from the deviceOauthClient configuration'
            )
        )
    }

    let timeout: NodeJS.Timeout | string | number | undefined = undefined
    const iframe = document.createElement('iframe')
    iframe.src = `${authority}/authorize?response_type=code&client_id=${clientId}&scope=${scopes}&audience=${
        audience || ''
    }&redirect_uri=${window.location.origin}/devices&device_id=${deviceId}`

    const destroyIframe = () => {
        sessionStorage.removeItem(DEVICE_AUTH_CODE_SESSION_KEY)
        iframe && iframe.parentNode && iframe.parentNode.removeChild(iframe)
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

/**
 * Add device by IP
 * @param {*} deviceIp
 */
export const addDeviceByIp = (deviceIp: string) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }?useEndpoints=${deviceIp}`
)


/**
 * Own device by deviceId
 */
export const ownDeviceApi = (deviceId: string) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }/${deviceId}/own`,
    {method: 'POST'}
)


/**
 * DisOwn device by deviceId
 */
export const disownDeviceApi = (deviceId: string) => fetchApi(
    `${config.httpGatewayAddress}${
        devicesApiEndpoints.DEVICES
    }/${deviceId}/disown`,
    {method: 'POST'}
)

