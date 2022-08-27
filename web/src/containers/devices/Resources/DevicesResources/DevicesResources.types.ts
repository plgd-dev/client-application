import { devicesStatuses } from '@/containers/devices/constants'
import {DeviceResourcesCrudType} from "@/containers/devices/Devices.types";

export type DevicesResourcesDeviceStatusType =
  typeof devicesStatuses[keyof typeof devicesStatuses]

export type Props = {
  data: {
    deviceId?: string
    href?: string
    interfaces: string[]
    resourceTypes: string[]
  }
  deviceId: string
  deviceStatus: DevicesResourcesDeviceStatusType
  isOwned: boolean
  loading: boolean
} & DeviceResourcesCrudType
