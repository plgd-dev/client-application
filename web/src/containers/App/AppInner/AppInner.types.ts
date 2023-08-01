import { WellKnownConfigType } from '@shared-ui/common/hooks'

export type Props = {
    configError: Error | undefined
    isIframe: boolean
    mockApp: boolean
    reFetchConfig: any
    setInitialize: (isInitialize?: boolean) => void
    wellKnownConfig: WellKnownConfigType
}
