import { useIntl } from 'react-intl'
import FindNewDeviceByIp from '../FindNewDeviceByIp'
import { messages as t } from '../../Devices.i18n'
import SplitButton from '@shared-ui/components/new/SplitButton'
import { FC } from 'react'
import { Props } from './DevicesListHeader.types'

const DevicesListHeader: FC<Props> = ({
  loading,
  refresh,
  openTimeoutModal,
}) => {
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

DevicesListHeader.displayName = 'DevicesListHeader'

export default DevicesListHeader
