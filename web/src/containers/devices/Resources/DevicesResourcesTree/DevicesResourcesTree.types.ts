import {DeviceResourcesCrudType} from "@/containers/devices/Devices.types";

export type Props = {
  data: {
    deviceId?: string
    href?: string
    interfaces: string[]
    resourceTypes: string[]
  }
  deviceId: string
  isOwned: boolean
  loading: boolean
} & DeviceResourcesCrudType
