import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'
import { Button } from '@/components/button'
import FindNewDeviceByIp from './find-new-device-by-ip'
import { messages as t } from './devices-i18n'

export const DevicesListHeader = ({ loading, refresh }) => {
  const { formatMessage: _ } = useIntl()

  return (
    <div className="d-flex align-items-center">
      <FindNewDeviceByIp />
      <Button disabled={loading} onClick={refresh} icon="fa-sync">
        {`${_(t.discovery)}`}
      </Button>
    </div>
  )
}

DevicesListHeader.propTypes = {
  loading: PropTypes.bool.isRequired,
  refresh: PropTypes.func.isRequired,
}
