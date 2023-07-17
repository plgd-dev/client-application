import { DeviceDataType, ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
    deviceId: string
    data: DeviceDataType
    isActiveTab: boolean
    isOwned: boolean
    isUnsupported: boolean
    onboardResourceLoading: boolean
    deviceOnboardingResourceData: any
    resources: ResourcesType[]
}
