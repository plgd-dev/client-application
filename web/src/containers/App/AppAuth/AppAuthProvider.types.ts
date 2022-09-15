import { WellKnownConfigType } from '@/containers/App/App.types'
import { ReactNode } from 'react'

export type Props = {
  wellKnownConfig?: WellKnownConfigType
  children: ReactNode
  setAuthError: (error: string) => void
}
