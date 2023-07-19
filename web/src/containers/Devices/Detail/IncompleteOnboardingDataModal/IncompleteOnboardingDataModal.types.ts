export type CopyDataType = {
    attribute: string
    attributeKey: keyof typeof onboardingDataDefault
    value: string
}

export type OnboardingDataType = typeof onboardingDataDefault

export type Props = {
    deviceId: string
    onClose?: () => void
    onSubmit?: (onboardingData: OnboardingDataType) => void
    onboardingData: OnboardingDataType
    show: boolean
}

export const onboardingDataDefault = {
    hubId: '',
    deviceEndpoint: '',
    authorizationCode: '',
    certificateAuthorities: '',
    authorizationProvider: '',
}

export const defaultProps = {
    show: false,
    onboardingData: onboardingDataDefault,
}
