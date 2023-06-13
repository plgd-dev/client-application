import React, { FC, useEffect, useMemo, useState } from 'react'
import { Props, defaultProps, onboardingDataDefault } from './IncompleteOnboardingDataModal.types'
import { useIntl } from 'react-intl'
import isFunction from 'lodash/isFunction'
import classNames from 'classnames'
import { validate as isValidUUID } from 'uuid'

import Modal from '@shared-ui/components/Atomic/Modal'
import Button from '@shared-ui/components/Atomic/Button'
import CopyBox from '@shared-ui/components/Atomic/CopyBox'
import Label from '@shared-ui/components/Atomic/Label'
import { WellKnownConfigType } from '@shared-ui/common/hooks'
import TextField from '@shared-ui/components/Atomic/TextField'

import { messages as t } from '@/containers/Devices/Devices.i18n'
import './IncompleteOnboardingDataModal.scss'

export const getOnboardingDataFromConfig = (wellKnowConfig: WellKnownConfigType) => ({
    coapGatewayAddress: wellKnowConfig?.remoteProvisioning?.coapGateway || '',
    authorizationProviderName: wellKnowConfig?.remoteProvisioning?.deviceOauthClient.providerName || '',
    hubId: wellKnowConfig?.remoteProvisioning?.id || '',
    certificateAuthorities: wellKnowConfig?.remoteProvisioning?.certificateAuthorities || '',
    authorizationCode: '',
})

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

    const renderBody = () => {
        return (
            <div>
                <Label inline title={_(t.onboardingFieldDeviceId)}>
                    <div className='auth-code-box'>
                        <TextField onChange={() => {}} readOnly={true} value={deviceId} />
                        <CopyBox textToCopy={deviceId} />
                    </div>
                </Label>
                <Label inline title={_(t.onboardingFieldHubId)}>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: !isValidUUID(onboardingData.hubId || '') })}
                            onChange={(e) => handleInputChange(e.target.value, 'hubId')}
                            value={onboardingData.hubId || ''}
                        />
                        <CopyBox textToCopy={onboardingData.hubId || ''} />
                    </div>
                </Label>
                <Label inline title={_(t.onboardingFieldDeviceEndpoint)}>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.coapGatewayAddress === '' })}
                            onChange={(e) => handleInputChange(e.target.value, 'coapGatewayAddress')}
                            value={onboardingData.coapGatewayAddress || ''}
                        />
                        <CopyBox textToCopy={onboardingData.coapGatewayAddress || ''} />
                    </div>
                </Label>
                <Label inline title={_(t.onboardingFieldAuthorizationCode)}>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.authorizationCode === '' })}
                            onChange={(e) => handleInputChange(e.target.value, 'authorizationCode')}
                            value={onboardingData.authorizationCode || ''}
                        />
                        <CopyBox textToCopy={onboardingData.authorizationCode || ''} />
                    </div>
                </Label>
                <Label inline title={_(t.onboardingFieldAuthorizationProvider)}>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.authorizationProviderName === '' })}
                            onChange={(e) => handleInputChange(e.target.value, 'authorizationProviderName')}
                            value={onboardingData.authorizationProviderName || ''}
                        />
                        <CopyBox textToCopy={onboardingData.authorizationProviderName || ''} />
                    </div>
                </Label>
                <Label inline title={_(t.onboardingFieldCertificateAuthority)}>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.certificateAuthorities === '' })}
                            onChange={(e) => handleInputChange(e.target.value, 'certificateAuthorities')}
                            value={onboardingData.certificateAuthorities || ''}
                        />
                        <CopyBox textToCopy={onboardingData.certificateAuthorities || ''} />
                    </div>
                </Label>
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
        const { coapGatewayAddress, authorizationProviderName, hubId, authorizationCode, certificateAuthorities } =
            onboardingData

        return (
            !coapGatewayAddress ||
            !authorizationProviderName ||
            !hubId ||
            !isValidUUID(hubId) ||
            !authorizationCode ||
            !certificateAuthorities
        )
    }, [onboardingData])

    const renderFooter = () => (
        <div className='w-100 d-flex justify-content-end'>
            <Button onClick={handleClose} variant='secondary'>
                {_(t.cancel)}
            </Button>

            <Button disabled={hasError} onClick={handleSubmit} variant='primary'>
                {_(t.onboardDevice)}
            </Button>
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
