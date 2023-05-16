import React, { FC } from 'react'
import { Props, defaultProps } from './FirstTimeOnboardingModal.types'
import Modal from '@shared-ui/components/new/Modal'
import { messages as t } from '@/containers/Devices/Devices.i18n'
import { useIntl } from 'react-intl'
import isFunction from 'lodash/isFunction'
import Button from '@shared-ui/components/new/Button'
import testId from '@/testId'

const FirstTimeOnboardingModal: FC<Props> = (props) => {
    const { show, onClose, onSubmit } = { ...defaultProps, ...props }
    const { formatMessage: _ } = useIntl()
    const { firstTimeModalButton } = testId.devices.detail

    const handleSubmit = () => {
        isFunction(onClose) && onClose && onClose()
        isFunction(onSubmit) && onSubmit && onSubmit()
    }

    const renderFooter = () => (
        <div className='w-100 d-flex justify-content-end'>
            <Button dataTestId={firstTimeModalButton} onClick={handleSubmit} variant='primary'>
                {_(t.ok)}
            </Button>
        </div>
    )

    return (
        <Modal
            onClose={onClose}
            renderBody={() => <div>{_(t.firstTimeDescription)}</div>}
            renderFooter={renderFooter}
            show={show}
            title={_(t.firstTimeTitle)}
        />
    )
}

FirstTimeOnboardingModal.displayName = 'FirstTimeOnboardingModal'
FirstTimeOnboardingModal.defaultProps = defaultProps

export default FirstTimeOnboardingModal
