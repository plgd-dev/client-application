import { devicesStatuses, shadowSynchronizationStates } from '../../constants'
import {DeviceDataType, ResourcesType} from '@/containers/devices/Devices.types'

export type DevicesDetailMetaDataStatusValueType =
  typeof devicesStatuses[keyof typeof devicesStatuses]

export type DevicesDetailMetaDataStatusShadowSynchronizationType =
  typeof shadowSynchronizationStates[keyof typeof shadowSynchronizationStates]


export type Props = {
  data: DeviceDataType
  loading: boolean
  isOwned: boolean
  resources: ResourcesType[]
  deviceId?: string
}
