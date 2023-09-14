import React, { FC, useEffect, useMemo, useState } from 'react'
import {
    Props,
    defaultProps,
    onboardingDataDefault,
    CopyDataType,
    OnboardingDataType,
} from './IncompleteOnboardingDataModal.types'
import { useIntl } from 'react-intl'
import isFunction from 'lodash/isFunction'
import { validate as isValidUUID } from 'uuid'
import isEmpty from 'lodash/isEmpty'
import { motion, AnimatePresence } from 'framer-motion'

import Modal from '@shared-ui/components/Atomic/Modal'
import Button from '@shared-ui/components/Atomic/Button'
import { WellKnownConfigType } from '@shared-ui/common/hooks'
import FormLabel from '@shared-ui/components/Atomic/FormLabel'
import FormInput from '@shared-ui/components/Atomic/FormInput'
import FormGroup from '@shared-ui/components/Atomic/FormGroup'
import FormTextarea from '@shared-ui/components/Atomic/FormTextarea'

import * as copyStyles from '@shared-ui/components/Atomic/CopyElement/CopyElement.styles'

import { messages as t } from '@/containers/Devices/Devices.i18n'
import * as styles from './IncompleteOnboardingDataModal.styles'

export const getOnboardingDataFromConfig = (wellKnowConfig: WellKnownConfigType) => ({
    hubId: wellKnowConfig?.remoteProvisioning?.id || '',
    deviceEndpoint: wellKnowConfig?.remoteProvisioning?.coapGateway || '',
    authorizationProvider: wellKnowConfig?.remoteProvisioning?.deviceOauthClient.providerName || '',
    certificateAuthorities: wellKnowConfig?.remoteProvisioning?.certificateAuthorities || '',
    authorizationCode: '',
})

const ALLOWED_ATTRIBUTES = [
    'hubId',
    'deviceEndpoint',
    'authorizationCode',
    'authorizationProvider',
    'certificateAuthorities',
]

const IncompleteOnboardingDataModal: FC<Props> = (props) => {
    const {
        deviceId,
        onClose,
        onSubmit,
        onboardingData: onboardingDataProps,
        show,
    } = {
        ...defaultProps,
        ...props,
    }

    const [onboardingData, setOnboardingData] = useState(onboardingDataProps || onboardingDataDefault)
    const [showImport, setShowImport] = useState(false)
    const [textareaValue, setTextareaValue] = useState('')

    useEffect(() => {
        setOnboardingData(onboardingDataProps)
    }, [onboardingDataProps])

    const { formatMessage: _ } = useIntl()

    const handleInputChange = (value: string, key: string) => {
        let dataForSave = value
        if (dataForSave.at(0) === '"' && dataForSave.at(-1) === '"') {
            dataForSave = dataForSave.substring(1)
            dataForSave = dataForSave.substring(0, dataForSave.length - 1)
        }
        setOnboardingData({ ...onboardingData, [key]: dataForSave })
    }

    const parsePasteData = (value: string) => {
        const dataOnSave: OnboardingDataType = {} as OnboardingDataType

        JSON.parse(value).forEach((item: CopyDataType) => {
            if (ALLOWED_ATTRIBUTES.includes(item.attributeKey)) {
                dataOnSave[item.attributeKey] = item.value
            }
        })

        if (!isEmpty(dataOnSave)) {
            setOnboardingData({ ...onboardingData, ...dataOnSave })
            setTextareaValue('')
            setShowImport(false)
        }
    }

    const renderBody = () => {
        return (
            <div>
                <div css={styles.topLine}>
                    <a
                        css={copyStyles.copy}
                        href='#'
                        onClick={(e) => {
                            e.preventDefault()
                            setShowImport(!showImport)
                        }}
                    >
                        <span css={copyStyles.text}>{_(t.pasteAll)}</span>
                    </a>
                </div>
                <AnimatePresence>
                    {showImport && (
                        <motion.div
                            animate={{ height: 120 }}
                            exit={{
                                height: 0,
                            }}
                            initial={{ height: 0 }}
                            transition={{
                                duration: 0.3,
                            }}
                        >
                            <FormTextarea
                                autoFocus={true}
                                onChange={(e: any) => setTextareaValue(e.target.value)}
                                onPaste={(e: any) => parsePasteData(e.clipboardData.getData('text'))}
                                value={textareaValue}
                            />
                        </motion.div>
                    )}
                </AnimatePresence>
                <div css={styles.row}>
                    <FormGroup inline id='form-group-device-id' inlineJustifyContent='space-between'>
                        <FormLabel text={_(t.onboardingFieldDeviceId)} />
                        <FormInput
                            copy={true}
                            disabled={true}
                            inputWrapperStyle={styles.inputWrapper}
                            value={deviceId}
                        />
                    </FormGroup>
                </div>
                <div css={styles.row}>
                    <FormGroup
                        inline
                        error={!isValidUUID(onboardingData.hubId || '')}
                        id='form-group-hubid'
                        inlineJustifyContent='space-between'
                    >
                        <FormLabel text={_(t.onboardingFieldHubId)} />
                        <FormInput
                            copy={true}
                            inputWrapperStyle={styles.inputWrapper}
                            onChange={(e) => handleInputChange(e.target.value, 'hubId')}
                            value={onboardingData.hubId || ''}
                        />
                    </FormGroup>
                </div>
                <div css={styles.row}>
                    <FormGroup
                        inline
                        error={onboardingData.deviceEndpoint === ''}
                        id='form-group-coapGatewayAddress'
                        inlineJustifyContent='space-between'
                    >
                        <FormLabel text={_(t.onboardingFieldDeviceEndpoint)} />
                        <FormInput
                            copy={true}
                            inputWrapperStyle={styles.inputWrapper}
                            onChange={(e) => handleInputChange(e.target.value, 'deviceEndpoint')}
                            value={onboardingData.deviceEndpoint || ''}
                        />
                    </FormGroup>
                </div>
                <div css={styles.row}>
                    <FormGroup
                        inline
                        error={onboardingData.authorizationCode === ''}
                        id='form-group-authorizationCode'
                        inlineJustifyContent='space-between'
                    >
                        <FormLabel text={_(t.onboardingFieldAuthorizationCode)} />
                        <FormInput
                            copy={true}
                            inputWrapperStyle={styles.inputWrapper}
                            onChange={(e) => handleInputChange(e.target.value, 'authorizationCode')}
                            value={onboardingData.authorizationCode || ''}
                        />
                    </FormGroup>
                </div>
                <div css={styles.row}>
                    <FormGroup
                        inline
                        error={onboardingData.authorizationProvider === ''}
                        id='form-group-authorizationProviderName'
                        inlineJustifyContent='space-between'
                    >
                        <FormLabel text={_(t.onboardingFieldAuthorizationProvider)} />
                        <FormInput
                            copy={true}
                            inputWrapperStyle={styles.inputWrapper}
                            onChange={(e) => handleInputChange(e.target.value, 'authorizationProvider')}
                            value={onboardingData.authorizationProvider || ''}
                        />
                    </FormGroup>
                </div>
                <div css={styles.row}>
                    <FormGroup
                        inline
                        error={onboardingData.certificateAuthorities === ''}
                        id='form-group-certificateAuthorities'
                        inlineJustifyContent='space-between'
                    >
                        <FormLabel text={_(t.onboardingFieldCertificateAuthority)} />
                        <FormInput
                            copy={true}
                            // inputWrapperStyle={styles.inputWrapper}
                            onChange={(e) => handleInputChange(e.target.value, 'certificateAuthorities')}
                            value={onboardingData.certificateAuthorities || ''}
                        />
                    </FormGroup>
                </div>
            </div>
        )
    }

    const handleClose = () => {
        isFunction(onClose) && onClose && onClose()
    }

    const handleSubmit = () => {
        isFunction(onClose) && onClose && onClose()
        isFunction(onSubmit) && onSubmit && onSubmit(onboardingData)
    }

    const hasError = useMemo(() => {
        const { deviceEndpoint, authorizationProvider, hubId, authorizationCode, certificateAuthorities } =
            onboardingData

        return (
            !deviceEndpoint ||
            !authorizationProvider ||
            !hubId ||
            !isValidUUID(hubId) ||
            !authorizationCode ||
            !certificateAuthorities
        )
    }, [onboardingData])

    const renderFooter = () => (
        <div className='w-100 d-flex justify-content-between align-items-center'>
            <div />
            <div className='modal-buttons'>
                <Button className='modal-button' onClick={handleClose} variant='secondary'>
                    {_(t.cancel)}
                </Button>

                <Button className='modal-button' disabled={hasError} onClick={handleSubmit} variant='primary'>
                    {_(t.onboardDevice)}
                </Button>
            </div>
        </div>
    )

    return (
        <Modal
            onClose={onClose}
            renderBody={renderBody}
            renderFooter={renderFooter}
            show={show}
            title={_(t.onboardIncompleteModalTitle)}
        />
    )
}

IncompleteOnboardingDataModal.displayName = 'IncompleteOnboardingDataModal'
IncompleteOnboardingDataModal.defaultProps = defaultProps

export default IncompleteOnboardingDataModal
