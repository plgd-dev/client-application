import { DevicesResourcesModalParamsType } from '@shared-ui/components/Organisms/DevicesResourcesModal/DevicesResourcesModal.types'
import { ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
    onClose?: () => void
    resources?: ResourcesType[]
    show: boolean
    updateResource?: (params: DevicesResourcesModalParamsType, resourceDataUpdate: any) => void | Promise<void>
}

export const defaultProps = {
    show: false,
}
