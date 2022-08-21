import { ResourcesType } from '@/containers/devices/Devices.types'

export type Props = {
  deviceId: string
  deviceName: string
  isOwned: boolean
  isUnregistered: boolean
  onOwnChange: () => void
  openDpsModal: () => void
  resources: ResourcesType[]
}
