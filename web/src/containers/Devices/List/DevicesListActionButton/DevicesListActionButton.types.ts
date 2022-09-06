export type Props = {
  deviceId: string
  isOwned: boolean
  onOwnChange: () => void
  onView: (deviceId: string) => void
  resourcesLoadedCallback: (data: any) => void
  showDpsModal: (deviceId: string) => void
  onDelete: () => void
}
