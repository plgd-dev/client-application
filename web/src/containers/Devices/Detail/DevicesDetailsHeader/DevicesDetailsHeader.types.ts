import { ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
    buttonsLoading?: boolean
    deviceId: string
    deviceName: string
    deviceOnboardingResourceData: any
    handleOpenEditDeviceNameModal: () => void
    incompleteOnboardingData: boolean
    isOwned: boolean
    isUnregistered: boolean
    onOwnChange: () => void
    onboardButtonCallback?: () => void
    onboardResourceLoading: boolean
    onboarding: boolean
    openDpsModal: () => void
    openOnboardingModal: () => void
    resources: ResourcesType[]
}
