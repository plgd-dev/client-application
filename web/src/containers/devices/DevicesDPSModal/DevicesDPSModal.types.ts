import {
    DevicesResourcesModalParamsType
} from "@/containers/devices/Resources/DevicesResourcesModal/DevicesResourcesModal.types";
import {ResourcesType} from "@/containers/devices/Devices.types";

export type Props = {
    onClose?: () => void
    resources?: ResourcesType[]
    show: boolean
    updateResource?: (params: DevicesResourcesModalParamsType, resourceDataUpdate: any) => void | Promise<void>
}

export const defaultProps = {
    show: false
}