import { fetchApi, security } from '@/common/services'
import { DEVICE_AUTH_CODE_SESSION_KEY } from '@/constants'
import { devicesApiEndpoints, knownResourceTypes } from './constants'
import { interfaceGetParam, resourceToUrl } from './utils'

/**
 * Get a single thing by its ID Rest Api endpoint
 * @param {*} params { deviceId }
 * @param {*} data
 */
export const getDeviceApi = deviceId =>
  fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}`
  )

/**
 * Delete a set of devices by their IDs Rest Api endpoint
 */
export const deleteDevicesApi = () => {
  // We split the fetch into multiple chunks due to the URL being too long for the browser to handle

  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }`,
    {
      method: 'DELETE',
    }
  )
}

/**
 * Get devices RESOURCES Rest Api endpoint
 * @param deviceId
 */
export const getDevicesResourcesAllApi = deviceId =>
  fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/${devicesApiEndpoints.DEVICES_RESOURCES_SUFFIX}`
  )

/**
 * Get devices RESOURCES Rest Api endpoint
 * @param {*} params { deviceId, href - resource href, currentInterface - interface }
 * @param {*} data
 */
export const getDevicesResourcesApi = ({
  deviceId,
  href,
  currentInterface = null,
}) =>
  fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/resources${href}?${interfaceGetParam(currentInterface)}`
  )

/**
 * Update devices RESOURCE Rest Api endpoint
 * @param {*} params { deviceId, href - resource href, currentInterface - interface, ttl - timeToLive }
 * @param {*} data
 */
export const updateDevicesResourceApi = (
  { deviceId, href, currentInterface = null, ttl },
  data
) => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/resources${href}?timeToLive=${ttl}&${interfaceGetParam(
      currentInterface
    )}`,
    { method: 'PUT', body: data, timeToLive: ttl }
  )
}

/**
 * Create devices RESOURCE Rest Api endpoint
 * @param {*} params { deviceId, href - resource href, currentInterface - interface, ttl - timeToLive }
 * @param {*} data
 */
export const createDevicesResourceApi = (
  { deviceId, href, currentInterface = null, ttl },
  data
) => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/resource-links${href}?timeToLive=${ttl}&${interfaceGetParam(
      currentInterface
    )}`,
    { method: 'POST', body: data, timeToLive: ttl }
  )
}

/**
 * Delete devices RESOURCE Rest Api endpoint
 * @param {*} params { deviceId, href - resource href}
 */
export const deleteDevicesResourceApi = ({ deviceId, href }) => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/resource-links${href}`,
    { method: 'DELETE' }
  )
}

/**
 * Update the shadowSynchronization of one Thing Rest Api endpoint
 * @param {*} deviceId
 * @param {*} shadowSynchronization
 */
export const updateDeviceShadowSynchronizationApi = (
  deviceId,
  shadowSynchronization
) => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/metadata`,
    { method: 'PUT', body: { shadowSynchronization } }
  )
}

/**
 * Returns an async function which resolves with a authorization code gathered from a rendered iframe, used for onboarding of a device.
 * @param {*} deviceId
 */
export const getDeviceAuthCode = deviceId => {
  return new Promise((resolve, reject) => {
    const { authority } = security.getGeneralConfig()
    const { clientId, audience, scopes = [] } = security.getDeviceOAuthConfig()

    if (!clientId) {
      return reject(
        new Error(
          'clientId is missing from the deviceOauthClient configuration'
        )
      )
    }

    let timeout = null
    const iframe = document.createElement('iframe')
    iframe.src = `${authority}/authorize?response_type=code&client_id=${clientId}&scope=${scopes}&audience=${
      audience || ''
    }&redirect_uri=${window.location.origin}/devices&device_id=${deviceId}`

    const destroyIframe = () => {
      sessionStorage.removeItem(DEVICE_AUTH_CODE_SESSION_KEY)
      iframe.parentNode.removeChild(iframe)
    }

    const doResolve = value => {
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
}

/**
 * Add device by IP
 * @param {*} deviceId
 */
export const addDeviceByIp = deviceId => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }?useEndpoints=${deviceId}`
  )
}

/**
 * Own device by deviceId
 * @param {*} params deviceId
 */
export const ownDeviceApi = deviceId => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/own`,
    { method: 'POST' }
  )
}

/**
 * DisOwn device by deviceId
 * @param {*} params deviceId
 */
export const disownDeviceApi = deviceId => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/disown`,
    { method: 'POST' }
  )
}

/**
 * Check if device is owned by user
 * @param {*} deviceId
 */
export const checkDeviceOwnerApi = deviceId => {
  return fetchApi(
    `${security.getGeneralConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/resources/${resourceToUrl(knownResourceTypes.OIC_SEC_DOXM)}`
  )
}
