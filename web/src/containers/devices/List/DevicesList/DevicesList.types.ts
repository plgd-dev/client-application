import {DeviceDataType, ResourcesType} from "@/containers/devices/Devices.types";

export type Props = {
    data: DeviceDataType
    loading: boolean
    onDeleteClick: () => void
    ownDevice: (isOwned: boolean, id: string, name: string) => void
    resourcesLoadedCallback: (resources: ResourcesType[]) => void
    selectedDevices: string[]
    setSelectedDevices: (data?: any) => void
    showDpsModal: (deviceId: string) => void
    unselectRowsToken?: string | number
}

export const defaultProps = {
    data: undefined
}