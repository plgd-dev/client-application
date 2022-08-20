import React, { useEffect, useRef, useState } from 'react'
import { useIntl } from 'react-intl'
import Button from '@shared-ui/components/new/Button'
import { Modal } from '@shared-ui/components/old/modal'
import { TextField } from '@shared-ui/components/old/text-field'
import Label from '@shared-ui/components/new/Label'
import {
  showErrorToast,
  showSuccessToast,
} from '@shared-ui/components/old/toast'
import { addDeviceByIp } from './rest'
import { messages as t } from './devices-i18n'
import { useIsMounted } from '@shared-ui/common/hooks'
import { addDevice } from '@/containers/devices/slice'
import { useDispatch } from 'react-redux'

const FindNewDevice = ({ disabled }) => {
  const [fetching, setFetching] = useState(false)
  const [show, setShow] = useState(false)
  const [error, setError] = useState(false)
  const [deviceIp, setDeviceIp] = useState('')
  const { formatMessage: _ } = useIntl()
  const baseInputRef = useRef(undefined)
  const isMounted = useIsMounted()
  const dispatch = useDispatch()

  useEffect(() => {
    if (deviceIp !== '') {
      // TODO: validation
      // !isIP(deviceIp) && !error && setError(true)
      // isIP(deviceIp) && error && setError(false)
    } else {
      error && setError(false)
    }
  }, [deviceIp, error])

  useEffect(() => {
    show && baseInputRef?.current?.focus()
  }, [show])

  const onClose = () => {
    !fetching && setShow(false) && setDeviceIp('')
  }

  const renderBody = () => (
    <Label
      title={_(t.deviceIp)}
      required={true}
      errorMessage={error ? _(t.invalidIp) : undefined}
    >
      <TextField
        value={deviceIp}
        onChange={e => setDeviceIp(e.target.value.trim())}
        placeholder={_(t.enterDeviceIp)}
        disabled={fetching}
        inputRef={baseInputRef}
        onKeyPress={e => (e.charCode === 13 ? handleFetch() : undefined)}
      />
    </Label>
  )

  const handleFetch = async () => {
    setFetching(true)

    try {
      const promise = addDeviceByIp(deviceIp)
      promise.then(response => {
        if (isMounted) {
          setFetching(false)
          const deviceData = response.data.result

          dispatch(addDevice(deviceData))

          showSuccessToast({
            title: _(t.deviceAddByIpSuccess),
            message: deviceData.data.content.n,
          })

          setDeviceIp('')
          setShow(false)
        }
      })
    } catch (e) {
      showErrorToast({
        title: _(t.deviceAddByIpError),
        message: e.message,
      })

      isMounted && setFetching(false)
    }
  }

  const renderFooter = () => {
    return (
      <div className="w-100 d-flex justify-content-end align-items-center">
        <Button variant="secondary" onClick={onClose} disabled={fetching}>
          {_(t.cancel)}
        </Button>

        <Button
          variant="primary"
          onClick={handleFetch}
          loading={fetching}
          disabled={fetching || error || deviceIp === ''}
        >
          {_(t.addDevice)}
        </Button>
      </div>
    )
  }

  return (
    <>
      <Button
        onClick={() => setShow(true)}
        className="m-r-10"
        icon="fa-plus"
        disabled={disabled}
      >
        {_(t.deviceByIp)}
      </Button>

      <Modal
        show={show}
        onClose={onClose}
        title={_(t.findDeviceByIp)}
        renderBody={renderBody}
        renderFooter={renderFooter}
        closeButton={!fetching}
      />
    </>
  )
}

export default FindNewDevice
