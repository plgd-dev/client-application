import { ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
  className?: string
  deviceId: string
  deviceName?: string
  isOwned: boolean
  loading: boolean
  resources: ResourcesType[]
  updateDeviceName: (title: string) => void
}

export const defaultProps = {
  resources: [],
}
