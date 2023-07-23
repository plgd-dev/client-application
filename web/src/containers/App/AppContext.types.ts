import { BuildInformationType } from '@shared-ui/common/hooks'

export type AppContextType = {
    buildInformation?: BuildInformationType | null
    collapsed?: boolean
    iframeMode?: boolean
    setCollapsed?: (collapsed: boolean) => void
    unauthorizedCallback?: () => void
}
