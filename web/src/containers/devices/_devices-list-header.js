import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'
import FindNewDeviceByIp from './LIst/FindNewDeviceByIp/FindNewDeviceByIp'
import { messages as t } from './Devices.i18n'
import SplitButton from '@shared-ui/components/new/SplitButton'

export const DevicesListHeader = ({ loading, refresh, openTimeoutModal }) => {
  const { formatMessage: _ } = useIntl()

  return (
    <div className="d-flex align-items-center">
      <FindNewDeviceByIp disabled={loading} />
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
