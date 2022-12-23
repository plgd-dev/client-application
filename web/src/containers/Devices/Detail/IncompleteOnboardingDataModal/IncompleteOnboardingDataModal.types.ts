export type OnboardingDataType = {
    authorizationCode?: string
    certificateAuthority?: string
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
    certificateAuthority: '',
    coapGateway: '',
    id: '',
    providerName: '',
}

export const defaultProps = {
    show: false,
    onboardingData: onboardingDataDefault,
}
