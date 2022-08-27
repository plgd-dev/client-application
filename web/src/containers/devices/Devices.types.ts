export type ResourcesType = {
  deviceId: string
  href: string
  interfaces: string[]
  resourceTypes: string[]
}

export type StreamApiPropsType = {
  data: any
  updateData: any
  loading: boolean
  error: string | null
  refresh: () => void
}

export type DeviceResourcesCrudType = {
  onCreate: (href: string) => Promise<void>
  onDelete: (href: string) => void
  onUpdate: ({
               deviceId,
               href,
             }: {
    deviceId?: string
    currentInterface?: string
    href: string
  }) => void | Promise<void>
}
