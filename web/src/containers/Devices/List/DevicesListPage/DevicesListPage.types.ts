import { ResourcesType } from '@/containers/Devices/Devices.types'

export type DpsDataType = {
    deviceId: string
    resources: ResourcesType[] | undefined
}
