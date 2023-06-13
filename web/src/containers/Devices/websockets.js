import { store, history } from '@/store'
import { Emitter } from '@shared-ui/common/services/emitter'
import { showInfoToast } from '@shared-ui/components/Atomic/Toast/Toast'
import {
    devicesStatuses,
    resourceEventTypes,
    DEVICES_STATUS_WS_KEY,
    DEVICES_REGISTERED_UNREGISTERED_COUNT_EVENT_KEY,
} from './constants'
import {
    getDeviceNotificationKey,
    getResourceRegistrationNotificationKey,
    getResourceUpdateNotificationKey,
} from './utils'
import { isNotificationActive } from './slice'
import { getDeviceApi } from './rest'
import { messages as t } from './Devices.i18n'

const { ONLINE, REGISTERED, UNREGISTERED } = devicesStatuses
const DEFAULT_NOTIFICATION_DELAY = 500

const getDeviceIds = (deviceId, deviceRegistered, deviceUnregistered) => {
    const _deviceRegistered = deviceRegistered ? deviceRegistered.deviceIds : deviceUnregistered.deviceIds

    return deviceId ? [deviceId] : _deviceRegistered
}

const getEventType = (deviceUnregistered) => {
    const _deviceUnregistered = deviceUnregistered ? UNREGISTERED : null

    return _deviceUnregistered ? REGISTERED : _deviceUnregistered
}

// WebSocket listener for device status change.
export const deviceStatusListener = async ({ deviceMetadataUpdated, deviceRegistered, deviceUnregistered }) => {
    if (deviceMetadataUpdated || deviceRegistered || deviceUnregistered) {
        const notificationsEnabled = isNotificationActive(DEVICES_STATUS_WS_KEY)(store.getState())

        setTimeout(
            async () => {
                const {
                    deviceId,
                    status: { value: deviceStatus } = {},
                    shadowSynchronization,
                } = deviceMetadataUpdated || {}

                const eventType = getEventType(deviceUnregistered)
                const deviceIds = getDeviceIds(deviceId, deviceRegistered, deviceUnregistered)

                const status = deviceStatus || eventType

                try {
                    deviceIds.forEach(async (deviceId) => {
                        // Emit an event: things.status.{deviceId}
                        Emitter.emit(`${DEVICES_STATUS_WS_KEY}.${deviceId}`, {
                            deviceId,
                            status,
                            shadowSynchronization,
                        })

                        // Get the notification state of a single device from redux store
                        const currentDeviceNotificationsEnabled = isNotificationActive(
                            getDeviceNotificationKey(deviceId)
                        )(store.getState())

                        // Show toast
                        if ((notificationsEnabled || currentDeviceNotificationsEnabled) && status !== UNREGISTERED) {
                            const { data: { name } = {} } = await getDeviceApi(deviceId)
                            const toastMessage = status === ONLINE ? t.deviceWentOnline : t.deviceWentOffline
                            showInfoToast(
                                {
                                    title: t.devicestatusChange,
                                    message: { message: toastMessage, params: { name } },
                                },
                                {
                                    onClick: () => {
                                        history.push(`/devices/${deviceId}`)
                                    },
                                    isNotification: true,
                                }
                            )
                        }
                    })
                } catch (error) {} // ignore error

                // If the event was registered or unregistered, emit an event with the number to increment by
                if ([REGISTERED, UNREGISTERED].includes(status)) {
                    // Emit an event: things-registered-unregistered-count
                    Emitter.emit(DEVICES_REGISTERED_UNREGISTERED_COUNT_EVENT_KEY, deviceIds.length)
                }
            },
            notificationsEnabled ? DEFAULT_NOTIFICATION_DELAY : 0
        )
    }
}

const getResources = (resourcePublished, resourceUnpublished) =>
    resourcePublished
        ? resourcePublished.resources // if resource was published, use the resources list from the event
        : resourceUnpublished.hrefs.map((href) => ({ href })) // if the resource was unpublished, create an array of ojects contaning hrefs, so that it matches the resources object

const getEvent = (resourcePublished) => (resourcePublished ? resourceEventTypes.ADDED : resourceEventTypes.REMOVED)

const getToastText = (isNew, s1, s2) => (isNew ? s1 : s2)

export const deviceResourceRegistrationListener =
    ({ deviceId, deviceName }) =>
    ({ resourcePublished, resourceUnpublished }) => {
        if (resourcePublished || resourceUnpublished) {
            // Device notifications must be enabled to see a toast message
            const notificationsEnabled = isNotificationActive(getDeviceNotificationKey(deviceId))(store.getState())

            const resources = getResources(resourcePublished, resourceUnpublished)
            const resourceRegistrationObservationWSKey = getResourceRegistrationNotificationKey(deviceId)
            const event = getEvent(resourcePublished)

            // Emit an event: things.resource.registration.{deviceId}
            Emitter.emit(`${resourceRegistrationObservationWSKey}.${event}`, {
                event,
                resources,
            })

            if (notificationsEnabled) {
                const isNew = event === resourceEventTypes.ADDED

                // If 5 or more resources came in the WS, show only one notification message
                if (resources.length >= 5) {
                    const toastTitle = getToastText(isNew, t.newResources, t.resourcesDeleted)
                    const toastMessage = getToastText(isNew, t.resourcesAdded, t.resourcesWereDeleted)
                    const onClickAction = () => {
                        history.push(`/devices/${deviceId}`)
                    }
                    // Show toast
                    showInfoToast(
                        {
                            title: toastTitle,
                            message: {
                                message: toastMessage,
                                params: { deviceName, deviceId, count: resources.length },
                            },
                        },
                        {
                            onClick: onClickAction,
                            isNotification: true,
                        }
                    )
                } else {
                    resources.forEach(({ href }) => {
                        const toastTitle = getToastText(isNew, t.newResource, t.resourceDeleted)
                        const toastMessage = getToastText(isNew, t.resourceAdded, t.resourceWithHrefWasDeleted)
                        const onClickAction = () => {
                            if (isNew) {
                                // redirect to resource and open resource modal
                                history.push(`/devices/${deviceId}${href}`)
                            } else {
                                // redirect to device
                                history.push(`/devices/${deviceId}`)
                            }
                        }
                        // Show toast
                        showInfoToast(
                            {
                                title: toastTitle,
                                message: {
                                    message: toastMessage,
                                    params: { href, deviceName, deviceId },
                                },
                            },
                            {
                                onClick: onClickAction,
                                isNotification: true,
                            }
                        )
                    })
                }
            }
        }
    }

export const deviceResourceUpdateListener =
    ({ deviceId, href, deviceName }) =>
    ({ resourceChanged }) => {
        if (resourceChanged) {
            const eventKey = getResourceUpdateNotificationKey(deviceId, href)
            const notificationsEnabled = isNotificationActive(eventKey)(store.getState())

            // Emit an event: things.resource.update.{deviceId}.{href}
            Emitter.emit(`${eventKey}`, resourceChanged.content)

            if (notificationsEnabled) {
                // Show toast
                showInfoToast(
                    {
                        title: t.resourceUpdated,
                        message: {
                            message: t.resourceUpdatedDesc,
                            params: { href, deviceName },
                        },
                    },
                    {
                        onClick: () => {
                            history.push(`/devices/${deviceId}${href}`)
                        },
                        isNotification: true,
                    }
                )
            }
        }
    }
