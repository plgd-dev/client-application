import React, { FC, useEffect, useRef, useState } from 'react'
import { useIntl } from 'react-intl'
import Button from '@shared-ui/components/Atomic/Button'
import Modal from '@shared-ui/components/Atomic/Modal'
import Notification from '@shared-ui/components/Atomic/Notification/Toast'
import { convertSize, IconPlus } from '@shared-ui/components/Atomic/Icon'
import FormGroup from '@shared-ui/components/Atomic/FormGroup'
import FormLabel from '@shared-ui/components/Atomic/FormLabel'
import FormInput from '@shared-ui/components/Atomic/FormInput'

import { addDeviceByIp } from '../../rest'
import { messages as t } from '../../Devices.i18n'
import { useIsMounted } from '@shared-ui/common/hooks'
import { addDevice } from '@/containers/Devices/slice'
import { useDispatch } from 'react-redux'
import { Props } from './FindNewDeviceByIp.types'
import * as styles from '@/containers/Devices/Detail/IncompleteOnboardingDataModal/IncompleteOnboardingDataModal.styles'

const FindNewDeviceByIp: FC<Props> = ({ disabled }) => {
    const [fetching, setFetching] = useState<boolean>(false)
    const [show, setShow] = useState<boolean>(false)
    const [error, setError] = useState<boolean>(false)
    const [deviceIp, setDeviceIp] = useState<string>('')
    const { formatMessage: _ } = useIntl()
    const baseInputRef = useRef<HTMLInputElement | undefined>(undefined)
    const isMounted = useIsMounted()
    const dispatch = useDispatch()

    useEffect(() => {
        if (deviceIp !== '') {
            // validation ?
        } else {
            error && setError(false)
        }
    }, [deviceIp, error])

    useEffect(() => {
        show && baseInputRef?.current?.focus()
    }, [show])

    const onClose = () => {
        if (!fetching) {
            setShow(false)
            setDeviceIp('')
        }
    }

    const renderBody = () => (
        <FormGroup error={error ? _(t.invalidIp) : undefined} id='form-group-device-ip'>
            <FormLabel required={true} text={_(t.deviceIp)} />
            <FormInput
                disabled={fetching}
                inputRef={baseInputRef}
                inputWrapperStyle={styles.inputWrapper}
                onChange={(e) => setDeviceIp(e.target.value.trim())}
                onKeyPress={(e) => (e.charCode === 13 ? handleFetch() : undefined)}
                placeholder={_(t.enterDeviceIp) as string}
                value={deviceIp}
            />
        </FormGroup>
    )

    const handleFetch = async () => {
        setFetching(true)

        try {
            const promise = addDeviceByIp(deviceIp)
            promise.then((response) => {
                if (isMounted) {
                    setFetching(false)
                    const deviceData = response.data.result

                    dispatch(addDevice(deviceData))

                    Notification.success({
                        title: _(t.deviceAddByIpSuccess),
                        message: deviceData.data.content.n,
                    })

                    setDeviceIp('')
                    setShow(false)
                }
            })
        } catch (e: any) {
            Notification.error({
                title: _(t.deviceAddByIpError),
                message: e.message,
            })

            isMounted && setFetching(false)
        }
    }

    const renderFooter = () => {
        return (
            <div className='w-100 d-flex justify-content-between align-items-center'>
                <div />
                <div className='modal-buttons'>
                    <Button className='modal-button' disabled={fetching} onClick={onClose} variant='secondary'>
                        {_(t.cancel)}
                    </Button>

                    <Button
                        className='modal-button'
                        disabled={fetching || error || deviceIp === ''}
                        loading={fetching}
                        onClick={handleFetch}
                        variant='primary'
                    >
                        {_(t.addDevice)}
                    </Button>
                </div>
            </div>
        )
    }

    return (
        <>
            <Button
                className='m-r-10'
                disabled={disabled}
                icon={<IconPlus {...convertSize(20)} />}
                onClick={() => setShow(true)}
            >
                {_(t.deviceByIp)}
            </Button>

            <Modal
                closeButton={!fetching}
                onClose={onClose}
                renderBody={renderBody}
                renderFooter={renderFooter}
                show={show}
                title={_(t.findDeviceByIp)}
            />
        </>
    )
}

export default FindNewDeviceByIp
