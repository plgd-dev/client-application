import { WellKnownConfigType } from '@shared-ui/common/hooks'

export type Props = {
    configError: Error | undefined
    initializedByAnother?: boolean
    mockApp: boolean
    reFetchConfig: any
    setInitialize: (isInitialize?: boolean) => void
    updateWellKnownConfig: (data: WellKnownConfigType) => void
    wellKnownConfig: WellKnownConfigType
}
