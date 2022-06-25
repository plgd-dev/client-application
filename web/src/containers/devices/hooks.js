import debounce from 'lodash/debounce'
import { useApi, useStreamApi, useEmitter } from '@/common/hooks'
import { useAppConfig } from '@/containers/app'

import {
  devicesApiEndpoints,
  DEVICES_STATUS_WS_KEY,
  resourceEventTypes,
} from './constants'
import { getResourceRegistrationNotificationKey } from './utils'

export const useDevicesList = () => {
  const { httpGatewayAddress } = useAppConfig()

  // Fetch the data
  const { data, updateData, ...rest } = useStreamApi(
    `${httpGatewayAddress}${devicesApiEndpoints.DEVICES}`,
    { telemetrySpan: 'get-devices' }
  )

  // Update the metadata when a WS event is emitted
  // useEmitter(DEVICES_STATUS_WS_KEY, newDeviceStatus => {
  //   if (data) {
  //     // Update the data with the current device status and shadowSynchronization
  //     updateData(updateDevicesDataStatus(data, newDeviceStatus))
  //   }
  // })

  return { data, updateData, ...rest }
}

export const useDeviceDetails = deviceId => {
  const { httpGatewayAddress } = useAppConfig()

  // Fetch the data
  const { data, updateData, ...rest } = useApi(
    `${httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}`,
    {
      telemetrySpan: 'get-device-detail',
    }
  )

  // Update the metadata when a WS event is emitted
  useEmitter(
    `${DEVICES_STATUS_WS_KEY}.${deviceId}`,
    debounce(({ status, shadowSynchronization }) => {
      if (data) {
        updateData({
          ...data,
          metadata: {
            ...data.metadata,
            shadowSynchronization,
            status: {
              ...data.metadata.status,
              value: status,
            },
          },
        })
      }
    }, 300)
  )

  return { data, updateData, ...rest }
}

export const useDeviceResources = deviceId => {
  const { httpGatewayAddress } = useAppConfig()

  // Fetch the data
  const { data, updateData, ...rest } = useStreamApi(
    `${httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/${devicesApiEndpoints.DEVICES_RESOURCES_SUFFIX}`,
    { parseResult: 'json', telemetrySpan: 'get-device-resources' }
  )

  useEmitter(
    getResourceRegistrationNotificationKey(deviceId),
    ({ event, resources: updatedResources }) => {
      if (data?.resources) {
        const resources = data.resources // get the first set of resources from an array, since it came from a stream of data
        let updatedLinks = []

        updatedResources.forEach(resource => {
          if (event === resourceEventTypes.ADDED) {
            const linkExists =
              resources.findIndex(link => link.href === resource.href) !== -1
            if (linkExists) {
              // Already exists, update
              updatedLinks = resources.map(link => {
                if (link.href === resource.href) {
                  return resource
                }

                return link
              })
            } else {
              updatedLinks = resources.concat(resource)
            }
          } else {
            updatedLinks = resources.filter(link => link.href !== resource.href)
          }
        })

        updateData([{ ...data, resources: updatedLinks }])
      }
    }
  )

  return { data, updateData, ...rest }
}
