import { FC, memo } from 'react'
import { useIntl } from 'react-intl'

import { Icon } from '@shared-ui/components/new/Icon'
import SplitButton from '@shared-ui/components/new/SplitButton'

import FindNewDeviceByIp from '../FindNewDeviceByIp'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesListHeader.types'

const DevicesListHeader: FC<Props> = memo(({ loading, refresh, openTimeoutModal }) => {
    const { formatMessage: _ } = useIntl()

    return (
        <div className='d-flex align-items-center'>
            <FindNewDeviceByIp disabled={loading} />
            <SplitButton
                disabled={loading}
                icon={<Icon icon='refresh' />}
                items={[
                    {
                        onClick: openTimeoutModal,
                        label: _(t.changeTimeout),
                        icon: 'edit',
                    },
                ]}
                onClick={refresh}
            >
                {_(t.discovery)}
            </SplitButton>
        </div>
    )
})

DevicesListHeader.displayName = 'DevicesListHeader'

export default DevicesListHeader
