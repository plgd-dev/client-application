import React, { FC, useMemo, useState } from 'react'
import { useIntl } from 'react-intl'
import isFunction from 'lodash/isFunction'

import Modal from '@shared-ui/components/new/Modal'
import Button from '@shared-ui/components/new/Button'
import FormGroup from '@shared-ui/components/new/FormGroup'
import FormLabel from '@shared-ui/components/new/FormLabel'
import FormInput from '@shared-ui/components/new/FormInput'

import { messages as t } from '@/containers/Devices/Devices.i18n'
import { isValidEndpoint } from '@/containers/Devices/utils'
import { knownResourceTypes } from '@/containers/Devices/constants'
import { Props, defaultProps } from './DevicesDPSModal.types'

const DevicesDPSModal: FC<Props> = (props) => {
    const { show, onClose, updateResource, resources } = {
        ...defaultProps,
        ...props,
    }
    const { formatMessage: _ } = useIntl()
    const [inputValue, setInputValue] = useState('')
    const [hasError, setHasError] = useState(false)
    const DpsResource = useMemo(
        () =>
            resources &&
            resources.find((resource) => resource.resourceTypes.includes(knownResourceTypes.X_PLGD_DPS_CONF)),
        [resources]
    )

    const handleInputChange = (e: any) => {
        const value = e.target.value
        const isValid = isValidEndpoint(value)

        !hasError && !isValid && setHasError(true)
        hasError && isValid && setHasError(false)
        setInputValue(e.target.value)
    }

    const renderBody = () => (
        <div>
            <FormGroup id='device-name'>
                <FormLabel text={_(t.deviceProvisioningServiceEndpoint)} />
                <FormInput onChange={handleInputChange} value={inputValue} />
            </FormGroup>
        </div>
    )

    const handleSubmit = () => {
        isFunction(onClose) && onClose && onClose()
        isFunction(updateResource) &&
            DpsResource &&
            updateResource(
                { href: DpsResource.href, currentInterface: '' },
                {
                    endpoint: inputValue,
                }
            )
    }

    const handleClose = () => {
        setInputValue('')
        setHasError(false)
        isFunction(onClose) && onClose && onClose()
    }

    const renderFooter = () => (
        <div className='w-100 d-flex justify-content-end'>
            <div />
            <div className='modal-buttons'>
                <Button className='modal-button' onClick={handleClose} variant='secondary'>
                    {_(t.cancel)}
                </Button>

                <Button
                    className='modal-button'
                    disabled={hasError || inputValue === ''}
                    onClick={handleSubmit}
                    variant='primary'
                >
                    {_(t.save)}
                </Button>
            </div>
        </div>
    )

    return (
        <Modal
            appRoot={document.getElementById('root')}
            onClose={onClose}
            portalTarget={document.getElementById('modal-root')}
            renderBody={renderBody}
            renderFooter={renderFooter}
            show={show}
            title={_(t.provisionNewDeviceTitle)}
        />
    )
}

DevicesDPSModal.displayName = 'DevicesDPSModal'
DevicesDPSModal.defaultProps = defaultProps

export default DevicesDPSModal
