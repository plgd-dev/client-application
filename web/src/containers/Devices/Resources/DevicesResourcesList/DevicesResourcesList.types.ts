import { DeviceResourcesCrudType } from '@/containers/Devices/Devices.types'

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
