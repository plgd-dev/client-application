import { FC, useMemo, useRef } from 'react'
import { useIntl } from 'react-intl'
import { useSelector } from 'react-redux'
import classNames from 'classnames'
import Button from '@shared-ui/components/new/Button'
import { canSetDPSEndpoint, getDeviceNotificationKey } from '../../utils'
import { isNotificationActive } from '../../slice'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesDetailsHeader.types'

export const DevicesDetailsHeader: FC<Props> = ({
    deviceId,
    isUnregistered,
    onOwnChange,
    isOwned,
    resources,
    openDpsModal,
    onboardResourceLoading,
    onboardButtonCallback,
    deviceOnboardingResourceData,
}) => {
    const { formatMessage: _ } = useIntl()
    const deviceNotificationKey = getDeviceNotificationKey(deviceId)
    const notificationsEnabled = useRef(false)
    notificationsEnabled.current = useSelector(isNotificationActive(deviceNotificationKey))

    const greyedOutClassName = classNames({
        'grayed-out': isUnregistered,
    })

    const hasDPS = useMemo(() => canSetDPSEndpoint(resources), [resources])
    const onboardButton = deviceOnboardingResourceData?.content?.cps

    return (
        <div className={classNames('d-flex align-items-center', greyedOutClassName)}>
            {onboardButton && (
                <Button
                    icon='fa-plus'
                    variant='secondary'
                    disabled={!isOwned || onboardResourceLoading}
                    className='m-r-10'
                    loading={onboardResourceLoading}
                    onClick={onboardButtonCallback}
                >
                    {onboardButton === 'uninitialized' ? _(t.onboardDevice) : _(t.offboardDevice)}
                </Button>
            )}
            <Button
                variant='secondary'
                icon={isOwned ? 'fa-cloud-download-alt' : 'fa-cloud-upload-alt'}
                onClick={onOwnChange}
                disabled={isUnregistered}
            >
                {isOwned ? _(t.disOwnDevice) : _(t.ownDevice)}
            </Button>
            {hasDPS && (
                <Button
                    icon='fa-bacon'
                    variant='secondary'
                    disabled={!isOwned}
                    className='m-l-10'
                    onClick={openDpsModal}
                >
                    {_(t.setDpsEndpoint)}
                </Button>
            )}
        </div>
    )
}

DevicesDetailsHeader.displayName = 'DevicesDetailsHeader'

export default DevicesDetailsHeader
