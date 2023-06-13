import { FC, memo } from 'react'
import { useIntl } from 'react-intl'

import { IconRefresh, IconTrash } from '@shared-ui/components/Atomic/Icon'
import SplitButton from '@shared-ui/components/Atomic/SplitButton'
import Button from '@shared-ui/components/Atomic/Button'

import FindNewDeviceByIp from '../FindNewDeviceByIp'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesListHeader.types'

const DevicesListHeader: FC<Props> = memo((props) => {
    const { loading, refresh, openTimeoutModal, handleFlashDevices, i18n } = props
    const { formatMessage: _ } = useIntl()

    return (
        <div className='d-flex align-items-center'>
            <FindNewDeviceByIp disabled={loading} />
            <SplitButton
                disabled={loading}
                icon={<IconRefresh />}
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
            <Button className='m-l-10' disabled={loading} icon={<IconTrash />} onClick={handleFlashDevices}>
                {i18n.flushCache}
            </Button>
        </div>
    )
})

DevicesListHeader.displayName = 'DevicesListHeader'

export default DevicesListHeader
