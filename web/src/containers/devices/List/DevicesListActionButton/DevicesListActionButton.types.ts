export type Props = {
  deviceId: string
  isOwned: boolean
  onOwnChange: () => void
  onView: (deviceId: string) => void
  resourcesLoadedCallback: () => void
  showDpsModal: (deviceId: string) => void
}
