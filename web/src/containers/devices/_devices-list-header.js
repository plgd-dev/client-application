import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'
import FindNewDeviceByIp from './find-new-device-by-ip'
import { messages as t } from './devices-i18n'
import { SplitButton } from '@/components/split-button'

export const DevicesListHeader = ({ loading, refresh, openTimeoutModal }) => {
  const { formatMessage: _ } = useIntl()

  return (
    <div className="d-flex align-items-center">
      <FindNewDeviceByIp />
      <SplitButton
        disabled={loading}
        onClick={refresh}
        icon="fa-sync"
        items={[
          {
            onClick: openTimeoutModal,
            label: _(t.changeTimeout),
            icon: 'fa-pen',
          },
        ]}
      >
        {`${_(t.discovery)}`}
      </SplitButton>
    </div>
  )
}

DevicesListHeader.propTypes = {
  loading: PropTypes.bool.isRequired,
  refresh: PropTypes.func.isRequired,
  openTimeoutModal: PropTypes.func.isRequired,
}
