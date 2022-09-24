import { WellKnownConfigType } from '@shared-ui/common/hooks'

export type Props = {
    configError: Error | undefined
    reFetchConfig: any
    setInitialize: (isInitialize?: boolean) => void
    wellKnownConfig: WellKnownConfigType
}
