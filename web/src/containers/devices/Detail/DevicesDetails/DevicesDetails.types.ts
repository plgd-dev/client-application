import { devicesStatuses, shadowSynchronizationStates } from '../../constants'
import { ResourcesType } from '@/containers/devices/Devices.types'

export type DevicesDetailMetaDataStatusValueType =
  typeof devicesStatuses[keyof typeof devicesStatuses]

export type DevicesDetailMetaDataStatusShadowSynchronizationType =
  typeof shadowSynchronizationStates[keyof typeof shadowSynchronizationStates]

export type Props = {
  data: {
    id: string
    types: string[]
    endpoints: string[]
    name: string
    metadata: {
      status: {
        value: DevicesDetailMetaDataStatusValueType
      }
      shadowSynchronization: DevicesDetailMetaDataStatusShadowSynchronizationType
    }
  }
  loading: boolean
  isOwned: boolean
  resources: ResourcesType[]
  deviceId?: string
}
