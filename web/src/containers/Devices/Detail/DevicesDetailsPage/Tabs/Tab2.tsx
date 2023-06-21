import { FC, useEffect, useMemo, useState } from 'react'
import { useResizeDetector } from 'react-resize-detector'
import { useParams } from 'react-router-dom'
import { useIntl } from 'react-intl'
import omit from 'lodash/omit'

import { useIsMounted } from '@shared-ui/common/hooks'
import DevicesResourcesModal from '@shared-ui/components/Organisms/DevicesResourcesModal'
import Notification from '@shared-ui/components/Atomic/Notification/Toast'
import { DevicesResourcesModalParamsType } from '@shared-ui/components/Organisms/DevicesResourcesModal/DevicesResourcesModal.types'
import DeleteModal from '@shared-ui/components/Atomic/Modal/components/DeleteModal'

import DevicesResources from '@/containers/Devices/Resources/DevicesResources'
import { Props } from './Tab2.types'
import { DevicesDetailsResourceModalData } from '../DevicesDetailsPage.types'
import { deviceResourceUpdateListener } from '@/containers/Devices/websockets'
import { messages as t } from '@/containers/Devices/Devices.i18n'
import { isNotificationActive } from '@/containers/Devices/slice'
import { createDevicesResourceApi, deleteDevicesResourceApi, getDevicesResourcesApi } from '@/containers/Devices/rest'
import {
    handleCreateResourceErrors,
    handleDeleteResourceErrors,
    handleFetchResourceErrors,
    handleUpdateResourceErrors,
    updateResourceMethod,
} from '@/containers/Devices/utils'
import DevicesDPSModal from '@/containers/Devices/DevicesDPSModal'
import { history } from '@/store'
import { defaultNewResource, resourceModalTypes } from '@/containers/Devices/constants'

const Tab2: FC<Props> = (props) => {
    const {
        closeDpsModal,
        deviceName,
        deviceStatus,
        isActiveTab,
        isOnline,
        isOwned,
        isUnregistered,
        loadingResources,
        resourcesData,
        refreshResources,
        showDpsModal,
    } = props
    const {
        id,
        href: hrefParam,
    }: {
        id: string
        href: string
    } = useParams()

    const { formatMessage: _ } = useIntl()
    const { ref, width, height } = useResizeDetector()
    const isMounted = useIsMounted()

    const [resourceModalData, setResourceModalData] = useState<DevicesDetailsResourceModalData | undefined>(undefined)
    const [loadingResource, setLoadingResource] = useState(false)
    const [savingResource, setSavingResource] = useState(false)
    const [deleteResourceHref, setDeleteResourceHref] = useState<string>('')
    const [resourceModal, setResourceModal] = useState(false)

    const resources = useMemo(() => resourcesData?.resources || [], [resourcesData])

    // Open the resource modal when href is present
    useEffect(
        () => {
            if (hrefParam && !loadingResources) {
                openUpdateModal({ href: `/${hrefParam}` })
            }
        },
        [hrefParam, loadingResources] // eslint-disable-line
    )

    // Fetches the resources supported types and sets its values to the modal data, which opens the modal.
    const openCreateModal = async (href: string) => {
        // If there is already a fetch for a resource, disable the next attempt for a fetch until the previous fetch finishes
        if (loadingResource) {
            return
        }

        setLoadingResource(true)

        try {
            const { data: deviceData } = await getDevicesResourcesApi({
                deviceId: id,
                href,
            })
            const supportedTypes = deviceData?.data?.content?.rts || []

            if (isMounted.current) {
                setLoadingResource(false)

                // Setting the data and opening the modal
                setResourceModalData({
                    data: {
                        href,
                        types: supportedTypes,
                    },
                    resourceData: {
                        ...defaultNewResource,
                        rt: supportedTypes,
                    },
                    type: resourceModalTypes.CREATE_RESOURCE,
                })
                setResourceModal(true)
            }
        } catch (error) {
            if (error && isMounted.current) {
                setLoadingResource(false)
                handleFetchResourceErrors(error, _)
            }
        }
    }

    // Created a new resource through rest API
    const createResource = async (
        { href, currentInterface = '' }: DevicesResourcesModalParamsType,
        resourceDataCreate: object
    ) => {
        setSavingResource(true)

        try {
            await createDevicesResourceApi({ deviceId: id, href, currentInterface }, resourceDataCreate)

            if (isMounted.current) {
                Notification.success({
                    title: _(t.resourceCreateSuccess),
                    message: _(t.resourceWasCreated),
                })

                refreshResources()
                setResourceModalData(undefined) // close modal
                setSavingResource(false)
            }
        } catch (error) {
            if (error && isMounted.current) {
                handleCreateResourceErrors(error, { id, href }, _)
                setSavingResource(false)
            }
        }
    }

    // Updates the resource through rest API
    const updateResource = async (
        { href, currentInterface = '' }: DevicesResourcesModalParamsType,
        resourceDataUpdate: any
    ) => {
        setSavingResource(true)

        await updateResourceMethod(
            { deviceId: id, href, currentInterface },
            resourceDataUpdate,
            () => {
                Notification.success({
                    title: _(t.resourceUpdateSuccess),
                    message: _(t.resourceWasUpdated),
                })
                handleCloseUpdateModal()
                setSavingResource(false)
            },
            (error: any) => {
                setSavingResource(false)
                handleUpdateResourceErrors(error, { id, href }, _)
            }
        )
    }

    // Fetches the resource and sets its values to the modal data, which opens the modal.
    const openUpdateModal = async ({ href, currentInterface = '' }: { href: string; currentInterface?: string }) => {
        // If there is already a fetch for a resource, disable the next attempt for a fetch until the previous fetch finishes
        if (loadingResource) {
            return
        }

        setLoadingResource(true)

        try {
            const { data: resourceData } = await getDevicesResourcesApi({
                deviceId: id,
                href,
                currentInterface,
            })

            omit(resourceData, ['data.content.if', 'data.content.rt'])

            if (isMounted.current) {
                setLoadingResource(false)

                // Retrieve the types and interfaces of this resource
                const { resourceTypes: types = [], interfaces = [] } =
                    resources?.find?.((link: { href: string }) => link.href === href) || {}

                // Setting the data and opening the modal
                setResourceModalData({
                    data: {
                        href,
                        types,
                        interfaces,
                    },
                    resourceData,
                })
                setResourceModal(true)
            }
        } catch (error) {
            if (error && isMounted.current) {
                setLoadingResource(false)
                handleFetchResourceErrors(error, _)
            }
        }
    }

    const handleCloseUpdateModal = () => {
        setResourceModalData(undefined)

        if (hrefParam) {
            // Remove the href from the URL when the update modal is closed
            history.replace(window.location.pathname.replace(`/${hrefParam}`, ''))
        }
    }

    const deleteResource = async () => {
        setLoadingResource(true)

        try {
            await deleteDevicesResourceApi({
                deviceId: id,
                href: deleteResourceHref || '',
            })

            if (isMounted.current) {
                Notification.success({
                    title: _(t.resourceDeleteSuccess),
                    message: _(t.resourceWasDeleted),
                })

                refreshResources()
                setLoadingResource(false)
                closeDeleteModal()
            }
        } catch (error) {
            if (error && isMounted.current) {
                handleDeleteResourceErrors(error, { id, href: deleteResourceHref }, _)
                setLoadingResource(false)
                closeDeleteModal()
            }
        }
    }

    const openDeleteModal = (href: string) => {
        setDeleteResourceHref(href)
    }

    const closeDeleteModal = () => {
        setDeleteResourceHref('')
    }

    console.log(resources)

    return (
        <div
            ref={ref}
            style={{
                height: '100%',
            }}
        >
            <DevicesResources
                data={resources}
                deviceId={id}
                deviceStatus={deviceStatus}
                isActiveTab={isActiveTab}
                isOwned={isOwned}
                loading={loadingResource}
                onCreate={openCreateModal}
                onDelete={openDeleteModal}
                onUpdate={openUpdateModal}
                pageSize={{ width, height: height ? height - 32 : 0 }} // tree switch
            />

            <DevicesResourcesModal
                {...resourceModalData}
                createResource={createResource}
                deviceId={id}
                deviceName={deviceName}
                deviceResourceUpdateListener={deviceResourceUpdateListener}
                fetchResource={openUpdateModal}
                i18n={{
                    close: _(t.close),
                    commandTimeout: _(t.commandTimeout),
                    create: _(t.create),
                    creating: _(t.creating),
                    deviceId: _(t.deviceId),
                    interfaces: _(t.interfaces),
                    notifications: _(t.notifications),
                    off: _(t.off),
                    on: _(t.on),
                    resourceInterfaces: _(t.resourceInterfaces),
                    retrieve: _(t.retrieve),
                    retrieving: _(t.retrieving),
                    supportedTypes: _(t.supportedTypes),
                    types: _(t.types),
                    update: _(t.update),
                    updating: _(t.updating),
                }}
                isDeviceOnline={isOnline}
                isNotificationActive={isNotificationActive}
                isUnregistered={isUnregistered}
                loading={savingResource}
                onClose={() => setResourceModal(false)}
                retrieving={loadingResource}
                show={resourceModal}
                updateResource={updateResource}
            />

            <DeleteModal
                footerActions={[
                    {
                        label: _(t.cancel),
                        onClick: closeDeleteModal,
                        variant: 'tertiary',
                    },
                    {
                        label: _(t.delete),
                        onClick: deleteResource,
                        variant: 'primary',
                    },
                ]}
                onClose={closeDeleteModal}
                show={!!deleteResourceHref}
                subTitle={_(t.deleteResourceMessageSubtitle)}
                title={_(t.deleteResourceMessage)}
            />

            <DevicesDPSModal
                onClose={closeDpsModal}
                resources={resources}
                show={showDpsModal}
                updateResource={updateResource}
            />
        </div>
    )
}

Tab2.displayName = 'Tab2'

export default Tab2
