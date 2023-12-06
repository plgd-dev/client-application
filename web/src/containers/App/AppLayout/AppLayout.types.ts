import { WellKnownConfigType } from '@shared-ui/common/hooks'

export type AppLayoutRefType = {
    getAuthProviderRef(): any
}

export type Props = {
    initializedByAnother: boolean
    mockApp: boolean
    reFetchConfig: () => Promise<any>
    suspectedUnauthorized: boolean
    wellKnownConfig: WellKnownConfigType
}
