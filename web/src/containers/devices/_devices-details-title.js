import { useState, useMemo } from 'react'
import PropTypes from 'prop-types'
import classNames from 'classnames'
import Form from 'react-bootstrap/Form'
import { useIntl } from 'react-intl'

import Button from '@shared-ui/components/new/Button'
import { showErrorToast } from '@shared-ui/components/new/Toast/Toast'
import { useIsMounted } from '@shared-ui/common/hooks'
import { getApiErrorMessage } from '@shared-ui/common/utils'
import { updateDevicesResourceApi } from './rest'
import { canChangeDeviceName, getDeviceChangeResourceHref } from './utils'
import { deviceResourceShape } from './shapes'
import { messages as t } from './devices-i18n'

export const DevicesDetailsTitle = ({
  className,
  deviceName,
  deviceId,
  updateDeviceName,
  isOwned,
  resources,
  ttl,
  loading,
  ...rest
}) => {
  const { formatMessage: _ } = useIntl()
  const [inputTitle, setInputTitle] = useState('')
  const [edit, setEdit] = useState(false)
  const [saving, setSaving] = useState(false)
  const isMounted = useIsMounted()
  const canUpdate = useMemo(
    () => canChangeDeviceName(resources) && isOwned,
    [resources, isOwned]
  )

  const onEditClick = () => {
    setInputTitle(deviceName || '')
    setEdit(true)
  }

  const onCloseClick = () => {
    setEdit(false)
  }

  const cancelSave = () => {
    setSaving(false)
    setEdit(false)
  }

  const onSave = async () => {
    if (inputTitle.trim() !== '' && inputTitle !== deviceName && canUpdate) {
      const href = getDeviceChangeResourceHref(resources)

      setSaving(true)

      try {
        const { data } = await updateDevicesResourceApi(
          { deviceId, href, ttl },
          {
            n: inputTitle,
          }
        )

        if (isMounted.current) {
          cancelSave()
          updateDeviceName(data?.n || inputTitle)
        }
      } catch (error) {
        if (error && isMounted.current) {
          showErrorToast({
            title: _(t.deviceNameChangeFailed),
            message: getApiErrorMessage(error),
          })
          cancelSave()
        }
      }
    } else {
      cancelSave()
    }
  }

  const handleKeyDown = e => {
    if (e.keyCode === 13) {
      // Enter
      onSave()
    } else if (e.keyCode === 27) {
      // Esc
      cancelSave()
    }
  }

  if (edit) {
    return (
      <div className="form-control-with-button h2-input">
        <Form.Control
          type="text"
          placeholder={`${_(t.enterDeviceName)}...`}
          value={inputTitle}
          onChange={e => setInputTitle(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={saving}
          autoFocus
        />
        <Button
          variant="primary"
          onClick={onSave}
          disabled={saving}
          loading={saving}
        >
          {saving ? _(t.saving) : _(t.save)}
        </Button>
        <Button
          className="close-button"
          variant="secondary"
          onClick={onCloseClick}
          disabled={saving}
        >
          <i className="fas fa-times" />
        </Button>
      </div>
    )
  }

  return (
    <h2
      {...rest}
      className={classNames(className, 'd-inline-flex align-items-center', {
        'title-with-icon': canUpdate,
      })}
      onClick={canUpdate ? onEditClick : null}
    >
      <span
        className={canUpdate ? 'link reveal-icon-on-hover icon-visible' : null}
      >
        {deviceName}
      </span>
      {canUpdate && <i className="fas fa-pen" />}
    </h2>
  )
}

DevicesDetailsTitle.propTypes = {
  className: PropTypes.string,
  deviceName: PropTypes.string,
  deviceId: PropTypes.string,
  loading: PropTypes.bool.isRequired,
  updateDeviceName: PropTypes.func.isRequired,
  isOwned: PropTypes.bool.isRequired,
  resources: PropTypes.arrayOf(deviceResourceShape),
  ttl: PropTypes.number,
}

DevicesDetailsTitle.defaultProps = {
  className: null,
  deviceName: null,
  deviceId: null,
  resources: [],
  ttl: 0,
}
