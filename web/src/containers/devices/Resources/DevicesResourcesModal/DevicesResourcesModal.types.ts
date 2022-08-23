import { ReactNode } from 'react'
import { resourceModalTypes } from '../../constants'

export type DevicesResourcesModalType =
  typeof resourceModalTypes[keyof typeof resourceModalTypes]

type DevicesResourcesModalParamsType = {
  href?: string
  currentInterface: string
}

export type Props = {
  confirmDisabled: boolean
  createResource: (
    { href, currentInterface }: DevicesResourcesModalParamsType,
    jsonData?: object
  ) => void
  data?: {
    href: string
    types: string[]
    interfaces: string[]
  }
  deviceId?: string
  fetchResource: ({
    href,
    currentInterface,
  }: DevicesResourcesModalParamsType) => void
  isDeviceOnline: boolean
  isUnregistered: boolean
  loading: boolean
  onClose?: () => void
  resourceData?: object
  retrieving: boolean
  ttlControl: ReactNode
  type?: DevicesResourcesModalType
  updateResource: (
    { href, currentInterface }: DevicesResourcesModalParamsType,
    jsonData?: object
  ) => void
}

export const defaultProps = {
  type: resourceModalTypes.UPDATE_RESOURCE,
}
