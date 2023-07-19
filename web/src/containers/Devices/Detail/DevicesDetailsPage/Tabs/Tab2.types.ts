export type Props = {
    closeDpsModal: () => void
    deviceName: string
    deviceStatus: string
    isActiveTab: boolean
    isOnline: boolean
    isOwned: boolean
    isUnregistered: boolean
    loadingResources: boolean
    resourcesData?: any
    refreshResources: () => void
    showDpsModal: boolean
}
