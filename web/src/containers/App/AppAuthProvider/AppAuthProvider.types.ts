import { WellKnownConfigType } from '@/containers/App/App.types'
import { ReactElement } from 'react'

export type AppAuthProviderRefType = {
    getSignOutMethod(): () => Promise<void>
}

export type Props = {
    children: ReactElement
    setAuthError: (error: string) => void
    setInitialize: (isInitialize?: boolean) => void
    wellKnownConfig?: WellKnownConfigType
}
