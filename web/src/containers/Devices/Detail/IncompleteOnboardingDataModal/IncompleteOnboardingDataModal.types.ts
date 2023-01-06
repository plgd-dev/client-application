export type OnboardingDataType = {
    authorizationCode?: string
    certificateAuthorities?: string
    coapGatewayAddress?: string
    hubId?: string
    authorizationProviderName?: string
}

export type Props = {
    deviceId: string
    onClose?: () => void
    onSubmit?: (onboardingData: OnboardingDataType) => void
    onboardingData: OnboardingDataType
    show: boolean
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
