import { WellKnownConfigType } from '@shared-ui/common/hooks'

export type Props = {
    initializedByAnother: boolean
    suspectedUnauthorized: boolean
    mockApp: boolean
    setInitialize: (isInitialize?: boolean) => void
    wellKnownConfig?: WellKnownConfigType
}
