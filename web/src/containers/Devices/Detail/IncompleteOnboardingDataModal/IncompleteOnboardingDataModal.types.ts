export type OnboardingDataType = {
    authorizationCode?: string
    certificateAuthorities?: string
    coapGatewayAddress?: string
    hubId?: string
    authorizationProviderName?: string
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
    coapGatewayAddress: '',
    hubId: '',
    authorizationProviderName: '',
}

export const defaultProps = {
    show: false,
    onboardingData: onboardingDataDefault,
}
