import { DeviceDataType, ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
    deviceId: string
    data: DeviceDataType
    isActiveTab: boolean
    isOwned: boolean
    onboardResourceLoading: boolean
    deviceOnboardingResourceData: any
    resources: ResourcesType[]
}
