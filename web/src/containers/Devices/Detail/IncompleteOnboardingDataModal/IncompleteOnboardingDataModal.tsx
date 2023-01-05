import React, { FC, useEffect, useMemo, useState } from 'react'
import { Props, defaultProps, onboardingDataDefault } from './IncompleteOnboardingDataModal.types'
import Modal from '@shared-ui/components/new/Modal'
import { messages as t } from '@/containers/Devices/Devices.i18n'
import { useIntl } from 'react-intl'
import Button from '@shared-ui/components/new/Button'
import CopyBox from '@shared-ui/components/new/CopyBox'
import isFunction from 'lodash/isFunction'
import TextField from '../../../../../shared-ui/src/components/new/TextField'
import classNames from 'classnames'
import Label from '../../../../../shared-ui/src/components/new/Label'
import { WellKnownConfigType } from '@shared-ui/common/hooks'
import './IncompleteOnboardingDataModal.scss'
import { validate as isValidUUID } from 'uuid'

export const getOnboardingDataFromConfig = (wellKnowConfig: WellKnownConfigType) => ({
    coapGatewayAddress: wellKnowConfig?.remoteProvisioning?.coapGateway || '',
    authorizationProviderName: wellKnowConfig?.remoteProvisioning?.deviceOauthClient.providerName || '',
    hubId: wellKnowConfig?.remoteProvisioning?.id || '',
    certificateAuthorities: wellKnowConfig?.remoteProvisioning?.certificateAuthorities || '',
    authorizationCode: '',
})

const IncompleteOnboardingDataModal: FC<Props> = (props) => {
    const {
        show,
        onClose,
        onSubmit,
        onboardingData: onboardingDataProps,
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
                <Label title={_(t.onboardingFieldHubId)} inline>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: !isValidUUID(onboardingData.hubId || '') })}
                            value={onboardingData.hubId || ''}
                            onChange={(e) => handleInputChange(e.target.value, 'hubId')}
                        />
                        <CopyBox textToCopy={onboardingData.hubId || ''} />
                    </div>
                </Label>
                <Label title={_(t.onboardingFieldDeviceEndpoint)} inline>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.coapGatewayAddress === '' })}
                            value={onboardingData.coapGatewayAddress || ''}
                            onChange={(e) => handleInputChange(e.target.value, 'coapGatewayAddress')}
                        />
                        <CopyBox textToCopy={onboardingData.coapGatewayAddress || ''} />
                    </div>
                </Label>
                <Label title={_(t.onboardingFieldAuthorizationCode)} inline>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.authorizationCode === '' })}
                            value={onboardingData.authorizationCode || ''}
                            onChange={(e) => handleInputChange(e.target.value, 'authorizationCode')}
                        />
                        <CopyBox textToCopy={onboardingData.authorizationCode || ''} />
                    </div>
                </Label>
                <Label title={_(t.onboardingFieldAuthorizationProvider)} inline>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.authorizationProviderName === '' })}
                            value={onboardingData.authorizationProviderName || ''}
                            onChange={(e) => handleInputChange(e.target.value, 'authorizationProviderName')}
                        />
                        <CopyBox textToCopy={onboardingData.authorizationProviderName || ''} />
                    </div>
                </Label>
                <Label title={_(t.onboardingFieldCertificateAuthority)} inline>
                    <div className='auth-code-box'>
                        <TextField
                            className={classNames({ error: onboardingData.certificateAuthorities === '' })}
                            value={onboardingData.certificateAuthorities || ''}
                            onChange={(e) => handleInputChange(e.target.value, 'certificateAuthorities')}
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
            <Button variant='secondary' onClick={handleClose}>
                {_(t.cancel)}
            </Button>

            <Button variant='primary' onClick={handleSubmit} disabled={hasError}>
                {_(t.onboardDevice)}
            </Button>
        </div>
    )

    return (
        <Modal
            show={show}
            onClose={onClose}
            title={_(t.onboardIncompleteModalTitle)}
            renderBody={renderBody}
            renderFooter={renderFooter}
        />
    )
}

IncompleteOnboardingDataModal.displayName = 'IncompleteOnboardingDataModal'
IncompleteOnboardingDataModal.defaultProps = defaultProps

export default IncompleteOnboardingDataModal
