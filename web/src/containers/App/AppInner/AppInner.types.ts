import { WellKnownConfigType } from '@/containers/App/App.types'

export type Props = {
  wellKnownConfig: WellKnownConfigType
  configError?: {
    message: string
  }
  setInitialize: (value: boolean) => void
}
