export type OnboardingDataType = {
    authorizationCode?: string
    certificateAuthorities?: string
    coapGateway?: string
    id?: string
    providerName?: string
}

export type Props = {
    onClose?: () => void
    onSubmit?: (onboardingData: OnboardingDataType) => void
    show: boolean
    onboardingData: OnboardingDataType
}

export const onboardingDataDefault = {
    authorizationCode: '',
    certificateAuthorities: '',
    coapGateway: '',
    id: '',
    providerName: '',
}

export const defaultProps = {
    show: false,
    onboardingData: onboardingDataDefault,
}
