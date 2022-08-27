import {DeviceDataType} from "@/containers/devices/Devices.types";

export type Props = {
    data: DeviceDataType
    selectedDevices: string[]
    setSelectedDevices: (data?: any) => void
    loading: boolean
    onDeleteClick: () => void
    unselectRowsToken: string | number
    ownDevice: (isOwned: boolean, id: string, name: string) => void
    showDpsModal: () => void
    resourcesLoadedCallback: () => void
}

export const defaultProps = {
    data: undefined
}