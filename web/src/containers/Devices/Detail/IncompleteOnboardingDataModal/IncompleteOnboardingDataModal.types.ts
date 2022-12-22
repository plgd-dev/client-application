export type OnboardingDataType = {
    authority?: string
    coapGateway?: string
    clientId?: string
    providerName?: string
    scopes?: string
    id?: string
}

export type Props = {
    onClose?: () => void
    show: boolean
    onboardingData: OnboardingDataType
}

export const onboardingDataDefault = {
    authority: '',
    coapGateway: '',
    clientId: '',
    providerName: '',
    scopes: '',
    id: '',
}

export const defaultProps = {
    show: false,
    onboardingData: onboardingDataDefault,
}
