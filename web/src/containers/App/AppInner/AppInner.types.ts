import { WellKnownConfigType } from '@/containers/App/App.types'

export type Props = {
    configError?: {
        message: string
    }
    setInitialize: (isInitialize?: boolean) => void
    wellKnownConfig: WellKnownConfigType
}
