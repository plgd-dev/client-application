import { ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
    deviceId: string
    deviceName: string
    isOwned: boolean
    isUnregistered: boolean
    onOwnChange: () => void
    openDpsModal: () => void
    resources: ResourcesType[]
    onboardResourceLoading: boolean
    onboardButtonCallback?: () => void
    deviceOnboardingResourceData: any
}
