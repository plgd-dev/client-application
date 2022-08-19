import { useMemo, useRef } from 'react'
import { useIntl } from 'react-intl'
import { useSelector } from 'react-redux'
import classNames from 'classnames'
import PropTypes from 'prop-types'
import Button from '@shared-ui/components/new/Button'
import { canSetDPSEndpoint, getDeviceNotificationKey } from './utils'
import { isNotificationActive } from './slice'
import { messages as t } from './devices-i18n'

export const DevicesDetailsHeader = ({
  deviceId,
  isUnregistered,
  onOwnChange,
  isOwned,
  resources,
  openDpsModal,
}) => {
  const { formatMessage: _ } = useIntl()
  const deviceNotificationKey = getDeviceNotificationKey(deviceId)
  const notificationsEnabled = useRef(false)
  notificationsEnabled.current = useSelector(
    isNotificationActive(deviceNotificationKey)
  )

  const greyedOutClassName = classNames({
    'grayed-out': isUnregistered,
  })

  const hasDPS = useMemo(() => canSetDPSEndpoint(resources), [resources])

  return (
    <div
      className={classNames('d-flex align-items-center', greyedOutClassName)}
    >
      <Button
        variant="secondary"
        icon={isOwned ? 'fa-cloud-download-alt' : 'fa-cloud-upload-alt'}
        onClick={onOwnChange}
        disabled={isUnregistered}
      >
        {isOwned ? _(t.disOwnDevice) : _(t.ownDevice)}
      </Button>
      {hasDPS && (
        <Button
          icon="fa-bacon"
          variant="secondary"
          disabled={!isOwned}
          className="m-l-10"
          onClick={openDpsModal}
        >
          {_(t.setDpsEndpoint)}
        </Button>
      )}
    </div>
  )
}

DevicesDetailsHeader.propTypes = {
  deviceId: PropTypes.string,
  deviceName: PropTypes.string,
  isUnregistered: PropTypes.bool.isRequired,
}

DevicesDetailsHeader.defaultProps = {
  deviceId: null,
  deviceName: null,
}
