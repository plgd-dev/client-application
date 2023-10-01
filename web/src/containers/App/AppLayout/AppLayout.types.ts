import { WellKnownConfigType } from '@shared-ui/common/hooks'

export type AppLayoutRefType = {
    getAuthProviderRef(): any
}

export type Props = {
    initializedByAnother: boolean
    suspectedUnauthorized: boolean
    mockApp: boolean
    setInitialize: (isInitialize?: boolean) => void
    updateWellKnownConfig: (data: WellKnownConfigType) => void
    wellKnownConfig?: WellKnownConfigType
}
