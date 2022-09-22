import { BuildInformationType } from '@/containers/App/App.types'

export type AppContextType = {
    buildInformation?: BuildInformationType | null
    collapsed: boolean
    setInitializedByAnother?: () => void
}
