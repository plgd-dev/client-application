import { FC, useEffect } from 'react'
import { useIntl } from 'react-intl'
import { useSelector } from 'react-redux'

import { WebSocketEventClient, eventFilters } from '@shared-ui/common/services'
import Switch from '@shared-ui/components/Atomic/Switch'

import { getResourceUpdateNotificationKey } from '../../utils'
import { isNotificationActive } from '../../slice'
import { deviceResourceUpdateListener } from '../../websockets'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesResourcesModalNotifications.types'

const DevicesResourcesModalNotifications: FC<Props> = ({ deviceId, deviceName, href, isUnregistered }) => {
    const { formatMessage: _ } = useIntl()
    const resourceUpdateObservationWSKey = getResourceUpdateNotificationKey(deviceId, href)
    const notificationsEnabled = useSelector(isNotificationActive(resourceUpdateObservationWSKey))

    useEffect(() => {
        if (isUnregistered) {
            // Unregister the WS when the device is unregistered
            WebSocketEventClient.unsubscribe(resourceUpdateObservationWSKey)
        }
    }, [isUnregistered, resourceUpdateObservationWSKey])

    const toggleNotifications = (e: any) => {
        if (e.target.checked) {
            // Request browser notifications
            // (browsers will explicitly disallow notification permission requests not triggered in response to a user gesture,
            // so we must call it to make sure the user has received a notification request)
            Notification?.requestPermission?.()

            // Register the WS
            WebSocketEventClient.subscribe(
                {
                    eventFilter: [eventFilters.RESOURCE_CHANGED],
                    resourceIdFilter: [`${deviceId}${href}`],
                },
                resourceUpdateObservationWSKey,
                deviceResourceUpdateListener({ deviceId, href, deviceName })
            )
        } else {
            WebSocketEventClient.unsubscribe(resourceUpdateObservationWSKey)
        }

        // dispatch(toggleActiveNotification(resourceUpdateObservationWSKey))
    }

    return (
        <Switch
            checked={notificationsEnabled}
            disabled={isUnregistered}
            id='resource-update-notifications'
            label={_(t.notifications)}
            onChange={toggleNotifications}
        />
    )
}

DevicesResourcesModalNotifications.displayName = 'DevicesResourcesModalNotifications'

export default DevicesResourcesModalNotifications
