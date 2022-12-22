import { FC, useEffect, useMemo, useState } from 'react'
import { Props, defaultProps } from './IncompleteOnboardingDataModal.types'
import Modal from '@shared-ui/components/new/Modal'
import { messages as t } from '@/containers/Devices/Devices.i18n'
import { useIntl } from 'react-intl'
import Button from '@shared-ui/components/new/Button'
import isFunction from 'lodash/isFunction'
import TextField from '../../../../../shared-ui/src/components/new/TextField'
import classNames from 'classnames'
import Label from '../../../../../shared-ui/src/components/new/Label'

const IncompleteOnboardingDataModal: FC<Props> = (props) => {
    const {
        show,
        onClose,
        onboardingData: onboardingDataProps,
    } = {
        ...defaultProps,
        ...props,
    }
    const [onboardingData, setOnboardingData] = useState(onboardingDataProps)
    const { formatMessage: _ } = useIntl()

    useEffect(() => {
        setOnboardingData(onboardingDataProps)
    }, [onboardingDataProps])

    const handleInputChange = (e: any, key: string) => {
        const value = e.target.value

        setOnboardingData({ ...onboardingData, [key]: value })
    }

    const renderBody = () => {
        return (
            <div>
                <Label title={_(t.onboardingFieldAuthority)} onClick={(e) => e.preventDefault()}>
                    <TextField
                        className={classNames({ error: onboardingData.authority === '' })}
                        value={onboardingData.authority || ''}
                        onChange={(e) => handleInputChange(e, 'authority')}
                    />
                </Label>
                <Label title={_(t.onboardingFieldCoapGateway)} onClick={(e) => e.preventDefault()}>
                    <TextField
                        className={classNames({ error: onboardingData.coapGateway === '' })}
                        value={onboardingData.coapGateway || ''}
                        onChange={(e) => handleInputChange(e, 'coapGateway')}
                    />
                </Label>
                <Label title={_(t.onboardingFieldClientId)} onClick={(e) => e.preventDefault()}>
                    <TextField
                        className={classNames({ error: onboardingData.clientId === '' })}
                        value={onboardingData.clientId || ''}
                        onChange={(e) => handleInputChange(e, 'clientId')}
                    />
                </Label>
                <Label title={_(t.onboardingFieldProviderName)} onClick={(e) => e.preventDefault()}>
                    <TextField
                        className={classNames({ error: onboardingData.providerName === '' })}
                        value={onboardingData.providerName || ''}
                        onChange={(e) => handleInputChange(e, 'providerName')}
                    />
                </Label>
                <Label title={_(t.onboardingFieldScopes)} onClick={(e) => e.preventDefault()}>
                    <TextField
                        className={classNames({ error: onboardingData.scopes === '' })}
                        value={onboardingData.scopes || ''}
                        onChange={(e) => handleInputChange(e, 'scopes')}
                    />
                </Label>
                <Label title={_(t.onboardingFieldId)} onClick={(e) => e.preventDefault()}>
                    <TextField
                        className={classNames({ error: onboardingData.id === '' })}
                        value={onboardingData.id || ''}
                        onChange={(e) => handleInputChange(e, 'id')}
                    />
                </Label>
            </div>
        )
    }

    const handleClose = () => {
        isFunction(onClose) && onClose && onClose()
    }

    const handleSubmit = () => {
        isFunction(onClose) && onClose && onClose()
    }

    const hasError = useMemo(() => {
        const { authority, coapGateway, clientId, providerName, scopes, id } = onboardingData

        return !authority || !coapGateway || !clientId || !providerName || !scopes || !id
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
