import React, { FC, useMemo, useState } from 'react'
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
import Tooltip from 'react-bootstrap/Tooltip'
import { v4 as uuidv4, validate as isValidUUID } from 'uuid'
import OverlayTrigger from 'react-bootstrap/OverlayTrigger'

type CopyEditableBlockType = {
    title: string
    data?: string
    onChange: (value: string) => void
    validator?: (data: string) => boolean
}

type EditBoxType = Omit<CopyEditableBlockType, 'title'>

const CopyEditableBlock = (props: CopyEditableBlockType) => {
    const { title, data, onChange, validator } = props
    const { formatMessage: _ } = useIntl()
    const validate = validator ? validator : (d: string) => d === ''
    const [editMode, setEditMode] = useState(validate(data || ''))

    const EditBox = (props: EditBoxType) => {
        const { data: defaultData, validator } = props
        const [data, setData] = useState(defaultData)
        const validate = validator ? validator : (d: string) => d === ''

        if (editMode) {
            return (
                <>
                    <TextField
                        className={classNames({ error: validate(data || '') })}
                        value={data || ''}
                        onChange={(e) => setData(e.target.value)}
                    />
                    <div className='copy-box'>
                        <OverlayTrigger
                            placement='right'
                            overlay={
                                <Tooltip id={`menu-item-tooltip-${uuidv4()}`} className='plgd-tooltip'>
                                    {_(t.save)}
                                </Tooltip>
                            }
                        >
                            <div
                                className='box m-l-10'
                                onClick={() => {
                                    onChange(data || '')
                                    setEditMode(false)
                                }}
                            >
                                <i className='fa fa-check' />
                            </div>
                        </OverlayTrigger>
                    </div>
                    <div className='copy-box'>
                        <OverlayTrigger
                            placement='right'
                            overlay={
                                <Tooltip id={`menu-item-tooltip-${uuidv4()}`} className='plgd-tooltip'>
                                    {_(t.cancel)}
                                </Tooltip>
                            }
                        >
                            <div
                                className='box m-l-10'
                                onClick={() => {
                                    setData(defaultData)
                                    setEditMode(false)
                                }}
                            >
                                <i className='fa fa-times' />
                            </div>
                        </OverlayTrigger>
                    </div>
                </>
            )
        }

        return (
            <>
                <span>{data || '-'}</span>
                <div className='copy-box'>
                    <OverlayTrigger
                        placement='right'
                        overlay={
                            <Tooltip id={`menu-item-tooltip-${uuidv4()}`} className='plgd-tooltip'>
                                {_(t.edit)}
                            </Tooltip>
                        }
                    >
                        <div className='box m-l-10' onClick={() => setEditMode(true)}>
                            <i className='fa fa-pen' />
                        </div>
                    </OverlayTrigger>
                </div>
                <CopyBox textToCopy={data} />
            </>
        )
    }

    return (
        <Label title={title} inline>
            <div className='auth-code-box'>
                <EditBox data={data} onChange={onChange} validator={validator} />
            </div>
        </Label>
    )
}

export const getOnboardingDataFromConfig = (wellKnowConfig: WellKnownConfigType) => ({
    coapGateway: wellKnowConfig?.remoteProvisioning?.coapGateway || '',
    providerName: wellKnowConfig?.remoteProvisioning?.deviceOauthClient.providerName || '',
    id: wellKnowConfig?.remoteProvisioning?.id || '',
    certificateAuthorities: wellKnowConfig?.remoteProvisioning?.certificateAuthorities || '',
    authorizationCode: '',
})

const IncompleteOnboardingDataModal: FC<Props> = (props) => {
    const { show, onClose, onSubmit } = {
        ...defaultProps,
        ...props,
    }

    const [onboardingData, setOnboardingData] = useState(onboardingDataDefault)

    const { formatMessage: _ } = useIntl()

    const handleInputChange = (value: string, key: string) => {
        setOnboardingData({ ...onboardingData, [key]: value })
    }

    const renderBody = () => {
        return (
            <div>
                <CopyEditableBlock
                    title={_(t.onboardingFieldHubId)}
                    data={onboardingData.id}
                    onChange={(value: string) => handleInputChange(value, 'id')}
                    validator={(d) => !isValidUUID(d)}
                />
                <CopyEditableBlock
                    title={_(t.onboardingFieldDeviceEndpoint)}
                    data={onboardingData.coapGateway}
                    onChange={(value: string) => handleInputChange(value, 'coapGateway')}
                />
                <CopyEditableBlock
                    title={_(t.onboardingFieldAuthorizationCode)}
                    data={onboardingData.authorizationCode}
                    onChange={(value: string) => handleInputChange(value, 'authorizationCode')}
                />
                <CopyEditableBlock
                    title={_(t.onboardingFieldAuthorizationProvider)}
                    data={onboardingData.providerName}
                    onChange={(value: string) => handleInputChange(value, 'providerName')}
                />
                <CopyEditableBlock
                    title={_(t.onboardingFieldCertificateAuthority)}
                    data={onboardingData.certificateAuthorities}
                    onChange={(value: string) => handleInputChange(value, 'certificateAuthorities')}
                />
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
        const { coapGateway, providerName, id, authorizationCode, certificateAuthorities } = onboardingData

        return !coapGateway || !providerName || !id || !isValidUUID(id) || !authorizationCode || !certificateAuthorities
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
