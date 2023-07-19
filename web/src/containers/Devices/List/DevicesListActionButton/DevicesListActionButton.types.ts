import { devicesOwnerships } from '@/containers/Devices/constants'

export type Props = {
    deviceId: string
    ownershipStatus: keyof typeof devicesOwnerships
    onOwnChange: () => void
    onView: (deviceId: string) => void
    resourcesLoadedCallback: (data: any) => void
    showDpsModal: (deviceId: string) => void
    onDelete: () => void
}
