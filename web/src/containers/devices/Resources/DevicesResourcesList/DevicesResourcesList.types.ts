export type Props = {
  data: {
    deviceId?: string
    href?: string
    interfaces: string[]
    resourceTypes: string[]
  }
  deviceId: string
  isOwned: boolean
  loading: boolean
  onCreate: () => void
  onDelete: () => void
  onUpdate: ({ deviceId, href }: { deviceId: string; href: string }) => void
}
