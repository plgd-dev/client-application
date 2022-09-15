import { fetchApi, security } from '@shared-ui/common/services'
import { devicesApiEndpoints } from './constants'
import { interfaceGetParam } from './utils'
import { signIdentityCsr } from '@/containers/App/AppRest'

type SecurityConfig = {
  httpGatewayAddress: string
  authority: string
}

const getConfig = () => security.getGeneralConfig() as SecurityConfig

/**
 * Get a single thing by its ID Rest Api endpoint
 */
export const getDeviceApi = (deviceId: string) =>
  fetchApi(
    `${getConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}`
  )

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
    `${getConfig().httpGatewayAddress}${
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
}: {
  deviceId: string
  href: string
  currentInterface?: string
}) =>
  fetchApi(
    `${getConfig().httpGatewayAddress}${
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
    }/${deviceId}/resource-links${href}?${interfaceGetParam(currentInterface)}`,
    { method: 'POST', body: data }
  )

/**
 * Delete devices RESOURCE Rest Api endpoint
 */
export const deleteDevicesResourceApi = ({
  deviceId,
  href,
}: {
  deviceId: string
  href: string
}) =>
  fetchApi(
    `${getConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/resource-links${href}`,
    { method: 'DELETE' }
  )

/**
 * Add device by IP
 */
export const addDeviceByIp = (deviceIp: string) =>
  fetchApi(
    `${getConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }?useEndpoints=${deviceIp}`
  )

/**
 * Own device by deviceId
 */
export const ownDeviceApi = (deviceId: string) => {
  fetchApi(
    `${getConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/own`,
    { method: 'POST' }
  ).then(result => {
    if (result?.data?.identityCertificateChallenge) {
      const state = result.data.identityCertificateChallenge.state
      //owning with csr
      // @ts-ignore
      const { authority } = security.getWebOAuthConfig()
      signIdentityCsr(
        authority,
        result.data.identityCertificateChallenge.certificateSigningRequest
      ).then(result => {
        fetchApi(
          `${getConfig().httpGatewayAddress}${
            devicesApiEndpoints.DEVICES
          }/${deviceId}/own/${state}`,
          {
            method: 'POST',
            body: {
              certificate: result.data.certificate,
            },
          }
        ).then(r => r)
      })
    } else {
      return result
    }
  })
}

/**
 * DisOwn device by deviceId
 */
export const disownDeviceApi = (deviceId: string) =>
  fetchApi(
    `${getConfig().httpGatewayAddress}${
      devicesApiEndpoints.DEVICES
    }/${deviceId}/disown`,
    { method: 'POST' }
  )
