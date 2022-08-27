import {FC, useMemo, useState} from 'react'
import Modal from '@shared-ui/components/new/Modal'
import {messages as t} from '@/containers/devices/Devices.i18n'
import Button from '@shared-ui/components/new/Button'
import {useIntl} from 'react-intl'
import isFunction from 'lodash/isFunction'
import classNames from 'classnames'
import TextField from '@shared-ui/components/new/TextField'
import Label from '@shared-ui/components/new/Label'
import {isValidEndpoint} from '@/containers/devices/utils'
import {knownResourceTypes} from '@/containers/devices/constants'
import {Props, defaultProps} from './DevicesDPSModal.types'

const DevicesDPSModal: FC<Props> = (props) => {
    const {
        show,
        onClose,
        updateResource,
        resources,
    } = {...defaultProps, ...props}
    const {formatMessage: _} = useIntl()
    const [inputValue, setInputValue] = useState('')
    const [hasError, setHasError] = useState(false)
    const DpsResource = useMemo(
        () =>
            resources &&
            resources.find(resource =>
                resource.resourceTypes.includes(knownResourceTypes.X_PLGD_DPS_CONF)
            ),
        [resources]
    )

    const handleInputChange = (e: any) => {
        const value = e.target.value
        const isValid = isValidEndpoint(value)

        !hasError && !isValid && setHasError(true)
        hasError && isValid && setHasError(false)
        setInputValue(e.target.value)
    }

    const renderBody = () => {
        return (
            <Label
                title={_(t.deviceProvisioningServiceEndpoint)}
                onClick={e => e.preventDefault()}
            >
                <TextField
                    className={classNames({error: hasError})}
                    value={inputValue}
                    onChange={handleInputChange}
                />
            </Label>
        )
    }

    const handleSubmit = () => {
        isFunction(onClose) && onClose && onClose()
        isFunction(updateResource) &&
        DpsResource &&
        updateResource(
            {href: DpsResource.href},
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
        <div className="w-100 d-flex justify-content-end">
            <Button variant="secondary" onClick={handleClose}>
                {_(t.cancel)}
            </Button>

            <Button
                variant="primary"
                onClick={handleSubmit}
                disabled={hasError || inputValue === ''}
            >
                {_(t.save)}
            </Button>
        </div>
    )

    return (
        <Modal
            show={show}
            onClose={onClose}
            title={_(t.provisionNewDeviceTitle)}
            renderBody={renderBody}
            renderFooter={renderFooter}
        />
    )
}

DevicesDPSModal.displayName = 'DevicesDPSModal'
DevicesDPSModal.defaultProps = defaultProps

export default DevicesDPSModal