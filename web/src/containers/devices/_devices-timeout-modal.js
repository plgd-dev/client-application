import { useState } from 'react'
import { Modal } from '@shared-ui/components/old/modal'
import { messages as t } from '@/containers/devices/devices-i18n'
import { Button } from '@shared-ui/components/old/button'
import PropTypes from 'prop-types'
import { useIntl } from 'react-intl'
import { CommanTimeoutControl } from './_command-timeout-control'
import { DISCOVERY_DEFAULT_TIMEOUT } from '@/containers/devices/constants'
import { useDispatch, useSelector } from 'react-redux'
import {
  getDevicesDiscoveryTimeout,
  setDiscoveryTimeout,
} from '@/containers/devices/slice'
import isFunction from 'lodash/isFunction'

export const DevicesTimeoutModal = ({ show, onClose }) => {
  const { formatMessage: _ } = useIntl()
  const dispatch = useDispatch()
  const discoveryTimeout = useSelector(getDevicesDiscoveryTimeout)

  const [userValue, setUserValue] = useState(discoveryTimeout)
  const [ttlHasError, setTtlHasError] = useState(false)

  const renderBody = () => {
    return (
      <CommanTimeoutControl
        title={_(t.discoveryTimeout)}
        defaultValue={discoveryTimeout}
        defaultTtlValue={DISCOVERY_DEFAULT_TIMEOUT}
        onChange={val => setUserValue(val)}
        ttlHasError={ttlHasError}
        onTtlHasError={setTtlHasError}
        disabled={false}
      />
    )
  }

  const handleSubmit = () => {
    if (userValue !== discoveryTimeout) {
      dispatch(setDiscoveryTimeout(userValue))
    }

    isFunction(onClose) && onClose()
  }

  const renderFooter = () => (
    <div className="w-100 d-flex justify-content-end">
      <Button variant="secondary" onClick={() => onClose()}>
        {_(t.cancel)}
      </Button>

      <Button variant="primary" onClick={handleSubmit} disabled={ttlHasError}>
        {_(t.save)}
      </Button>
    </div>
  )

  return (
    <Modal
      show={show}
      onClose={onClose}
      title={_(t.changeDiscoveryTimeout)}
      renderBody={renderBody}
      renderFooter={renderFooter}
      onExited={() => setUserValue(discoveryTimeout)}
    />
  )
}

DevicesTimeoutModal.propTypes = {
  onClose: PropTypes.func,
  show: PropTypes.bool.isRequired,
}

DevicesTimeoutModal.defaultProps = {
  show: false,
}
