import { ResourcesType } from '@/containers/Devices/Devices.types'

export type Props = {
  className?: string
  deviceId: string
  deviceName?: string
  isOwned: boolean
  loading: boolean
  resources: ResourcesType[]
  ttl: number
  updateDeviceName: (title: string) => void
}

export const defaultProps = {
  resources: [],
  ttl: 0,
}
