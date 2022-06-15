import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'
import { ActionButton } from '@/components/action-button'
import { messages as t } from './devices-i18n'

export const DevicesListActionButton = ({
  deviceId,
  onView,
  isOwned,
  onOwnChange,
}) => {
  const { formatMessage: _ } = useIntl()

  return (
    <ActionButton
      menuProps={{
        align: 'right',
      }}
      items={[
        {
          onClick: () => onView(deviceId),
          label: _(t.details),
          icon: 'fa-eye',
        },
        {
          onClick: () => onOwnChange(),
          label: isOwned ? _(t.disOwnDevice) : _(t.ownDevice),
          icon: isOwned ? 'fa-cloud-download-alt' : 'fa-cloud-upload-alt',
        },
      ]}
    >
      <i className="fas fa-ellipsis-h" />
    </ActionButton>
  )
}

DevicesListActionButton.propTypes = {
  deviceId: PropTypes.string.isRequired,
  onView: PropTypes.func.isRequired,
}
