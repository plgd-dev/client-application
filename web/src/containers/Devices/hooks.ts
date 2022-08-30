import debounce from 'lodash/debounce'
import { useStreamApi, useEmitter } from '@shared-ui/common/hooks'
import { useAppConfig } from '@/containers/App'

import {
  devicesApiEndpoints,
  DEVICES_STATUS_WS_KEY,
  resourceEventTypes,
  TIMEOUT_UNIT_PRECISION,
} from './constants'
import { getResourceRegistrationNotificationKey } from './utils'
import { useSelector } from 'react-redux'
import { getDevicesDiscoveryTimeout } from '@/containers/Devices/slice'
import { StreamApiPropsType } from '@/containers/Devices/Devices.types'

export const useDevicesList = () => {
  const { httpGatewayAddress } = useAppConfig()
  const discoveryTimeout = useSelector(getDevicesDiscoveryTimeout)

  // Fetch the data
  const { data, updateData, ...rest } = useStreamApi(
    `${httpGatewayAddress}${devicesApiEndpoints.DEVICES}`,
    {
      shadowQueryParameter: `?timeout=${
        discoveryTimeout / TIMEOUT_UNIT_PRECISION
      }`,
    }
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

export const useDeviceDetails = (deviceId: string) => {
  const { httpGatewayAddress } = useAppConfig()

  // Fetch the data
  const { data, updateData, ...rest }: StreamApiPropsType = useStreamApi(
    `${httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}`,
    { streamApi: false }
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

export const useDevicesResources = (deviceId: string) => {
  const { httpGatewayAddress } = useAppConfig()

  // Fetch the data
  const { data, updateData, ...rest }: StreamApiPropsType = useStreamApi(
    `${httpGatewayAddress}${devicesApiEndpoints.DEVICES}/${deviceId}/${devicesApiEndpoints.DEVICES_RESOURCES_SUFFIX}`,
    { parseResult: 'json' }
  )

  useEmitter(
    getResourceRegistrationNotificationKey(deviceId),
    ({
      event,
      resources: updatedResources,
    }: {
      event: any
      resources: any
    }) => {
      if (data?.resources) {
        const resources = data.resources // get the first set of resources from an array, since it came from a stream of data
        let updatedLinks: any = []

        updatedResources.forEach((resource: any) => {
          if (event === resourceEventTypes.ADDED) {
            const linkExists =
              resources.findIndex(
                (link: any) => link.href === resource.href
              ) !== -1
            if (linkExists) {
              // Already exists, update
              updatedLinks = resources.map((link: any) =>
                link.href === resource.href ? resource : link
              )
            } else {
              updatedLinks = resources.concat(resource)
            }
          } else {
            updatedLinks = resources.filter(
              (link: any) => link.href !== resource.href
            )
          }
        })

        updateData([{ ...data, resources: updatedLinks }])
      }
    }
  )

  return { data, updateData, ...rest }
}
