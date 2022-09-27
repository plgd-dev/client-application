import { useEffect, useMemo, useState } from 'react'
import { useIntl } from 'react-intl'
import { useParams } from 'react-router-dom'
import classNames from 'classnames'
import { history } from '@/store'
import ConfirmModal from '@shared-ui/components/new/ConfirmModal'
import Layout from '@shared-ui/components/new/Layout'
import NotFoundPage from '@/containers/NotFoundPage'
import { useIsMounted } from '@shared-ui/common/hooks'
import { messages as menuT } from '@shared-ui/components/new/Menu/Menu.i18n'
import { showSuccessToast } from '@shared-ui/components/new/Toast/Toast'
import DevicesDetails from '../DevicesDetails'
import DevicesResources from '../../Resources/DevicesResources'
import DevicesDetailsHeader from '../DevicesDetailsHeader'
import DevicesDetailsTitle from '../DevicesDetailsTitle'
import DevicesResourcesModal from '../../Resources/DevicesResourcesModal'
import DevicesDPSModal from '../../DevicesDPSModal'
import {
    devicesStatuses,
    defaultNewResource,
    resourceModalTypes,
    NO_DEVICE_NAME,
    devicesOwnerships,
} from '../../constants'
import {
    handleCreateResourceErrors,
    handleFetchResourceErrors,
    handleDeleteResourceErrors,
    handleDeleteDevicesErrors,
    updateResourceMethod,
    handleUpdateResourceErrors,
} from '../../utils'
import {
    getDevicesResourcesApi,
    createDevicesResourceApi,
    deleteDevicesResourceApi,
    ownDeviceApi,
    disownDeviceApi,
} from '../../rest'
import { useDeviceDetails, useDevicesResources } from '../../hooks'
import { messages as t } from '../../Devices.i18n'
import './DevicesDetailsPage.scss'
import { disOwnDevice, ownDevice } from '@/containers/Devices/slice'
import { useDispatch } from 'react-redux'
import { BreadcrumbItem } from '@shared-ui/components/new/Breadcrumbs/Breadcrumbs.types'
import omit from 'lodash/omit'
import { DevicesDetailsResourceModalData } from '@/containers/Devices/Detail/DevicesDetailsPage/DevicesDetailsPage.types'
import { DevicesResourcesModalParamsType } from '@/containers/Devices/Resources/DevicesResourcesModal/DevicesResourcesModal.types'

const DevicesDetailsPage = () => {
    const { formatMessage: _ } = useIntl()
    const {
        id,
        href: hrefParam,
    }: {
        id: string
        href: string
    } = useParams()
    const [resourceModalData, setResourceModalData] = useState<DevicesDetailsResourceModalData | undefined>(undefined)
    const [loadingResource, setLoadingResource] = useState(false)
    const [savingResource, setSavingResource] = useState(false)
    const [showDpsModal, setShowDpsModal] = useState(false)
    const [deleteResourceHref, setDeleteResourceHref] = useState<string>('')
    const [ttlHasError] = useState(false)
    const isMounted = useIsMounted()
    const { data, updateData, loading, error: deviceError } = useDeviceDetails(id)
    const {
        data: resourcesData,
        loading: loadingResources,
        error: resourcesError,
        refresh: refreshResources,
    } = useDevicesResources(id)
    const dispatch = useDispatch()

    const isOwned = useMemo(() => data?.ownershipStatus === devicesOwnerships.OWNED, [data])

    const resources = resourcesData?.resources || []

    // Open the resource modal when href is present
    useEffect(
        () => {
            if (hrefParam && !loading && !loadingResources) {
                openUpdateModal({ href: `/${hrefParam}` })
            }
        },
        [hrefParam, loading, loadingResources] // eslint-disable-line
    )

    if (deviceError) {
        return <NotFoundPage title={_(t.deviceNotFound)} message={_(t.deviceNotFoundMessage, { id })} />
    }

    if (resourcesError) {
        return (
            <NotFoundPage title={_(t.deviceResourcesNotFound)} message={_(t.deviceResourcesNotFoundMessage, { id })} />
        )
    }

    const deviceStatus = data?.metadata?.status?.value
    const isOnline = true
    const isUnregistered = devicesStatuses.UNREGISTERED === deviceStatus
    const greyedOutClassName = classNames({
        'grayed-out': isUnregistered,
    })
    const deviceName = data?.data?.content?.n || NO_DEVICE_NAME
    const breadcrumbs: BreadcrumbItem[] = [
        {
            to: '/',
            label: _(menuT.devices),
        },
    ]

    if (deviceName) {
        breadcrumbs.push({ label: deviceName })
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
            }
        } catch (error) {
            if (error && isMounted.current) {
                setLoadingResource(false)
                handleFetchResourceErrors(error, _)
            }
        }
    }

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
            }
        } catch (error) {
            if (error && isMounted.current) {
                setLoadingResource(false)
                handleFetchResourceErrors(error, _)
            }
        }
    }

    const openDeleteModal = (href: string) => {
        setDeleteResourceHref(href)
    }

    const closeDeleteModal = () => {
        setDeleteResourceHref('')
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
                showSuccessToast({
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

    // Created a new resource through rest API
    const createResource = async (
        { href, currentInterface = '' }: DevicesResourcesModalParamsType,
        resourceDataCreate: object
    ) => {
        setSavingResource(true)

        try {
            await createDevicesResourceApi({ deviceId: id, href, currentInterface }, resourceDataCreate)

            if (isMounted.current) {
                showSuccessToast({
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

    const deleteResource = async () => {
        setLoadingResource(true)

        try {
            await deleteDevicesResourceApi({
                deviceId: id,
                href: deleteResourceHref || '',
            })

            if (isMounted.current) {
                showSuccessToast({
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

    // Handler which cleans up the resource modal data and updates the URL
    const handleCloseUpdateModal = () => {
        setResourceModalData(undefined)

        if (hrefParam) {
            // Remove the href from the URL when the update modal is closed
            history.replace(window.location.pathname.replace(`/${hrefParam}`, ''))
        }
    }

    // Update the device name in the data object
    const updateDeviceNameInData = (name: string) => {
        updateData({
            ...data,
            data: {
                ...data.data,
                content: {
                    ...data.data.content,
                    n: name,
                },
            },
        })
    }

    const handleOwnChange = async () => {
        try {
            isOwned ? await disownDeviceApi(id) : await ownDeviceApi(id)
            const newOwnState = !isOwned

            if (isMounted.current) {
                updateData({
                    ...data,
                    ownershipStatus: newOwnState ? devicesOwnerships.OWNED : devicesOwnerships.UNOWNED,
                })

                if (!newOwnState) {
                    // @ts-ignore
                    dispatch(disOwnDevice(id))
                    history.push('/')
                } else {
                    // @ts-ignore
                    dispatch(ownDevice(id))
                }

                showSuccessToast({
                    title: newOwnState ? _(t.deviceOwned) : _(t.deviceDisOwned),
                    message: newOwnState
                        ? _(t.deviceWasOwned, { name: deviceName })
                        : _(t.deviceWasDisOwned, { name: deviceName }),
                })
            }
        } catch (error) {
            handleDeleteDevicesErrors(error, _, true)
        }
    }

    return (
        <Layout
            title={`${deviceName ? deviceName + ' | ' : ''}${_(menuT.devices)}`}
            breadcrumbs={breadcrumbs}
            loading={loading || (!resourceModalData && loadingResource)}
            header={
                <DevicesDetailsHeader
                    deviceId={id}
                    deviceName={deviceName}
                    isOwned={isOwned}
                    onOwnChange={handleOwnChange}
                    isUnregistered={isUnregistered}
                    resources={resources}
                    openDpsModal={() => setShowDpsModal(true)}
                />
            }
        >
            <DevicesDetailsTitle
                className={classNames(
                    {
                        shimmering: loading,
                    },
                    greyedOutClassName
                )}
                updateDeviceName={updateDeviceNameInData}
                loading={loading}
                isOwned={isOwned}
                deviceName={deviceName}
                deviceId={id}
                resources={resources}
            />
            <DevicesDetails data={data} isOwned={isOwned} loading={loading} resources={resources} deviceId={id} />

            <DevicesResources
                data={resources}
                onUpdate={openUpdateModal}
                onCreate={openCreateModal}
                onDelete={openDeleteModal}
                deviceStatus={deviceStatus}
                loading={loadingResource}
                deviceId={id}
                isOwned={isOwned}
            />

            <DevicesResourcesModal
                {...resourceModalData}
                onClose={handleCloseUpdateModal}
                fetchResource={openUpdateModal}
                updateResource={updateResource}
                createResource={createResource}
                retrieving={loadingResource}
                loading={savingResource}
                isDeviceOnline={isOnline}
                isUnregistered={isUnregistered}
                deviceId={id}
                confirmDisabled={ttlHasError}
            />

            <ConfirmModal
                onConfirm={deleteResource}
                show={deleteResourceHref !== ''}
                title={
                    <>
                        <i className='fas fa-trash-alt' />
                        {`${_(t.delete)} ${deleteResourceHref}`}
                    </>
                }
                body={<>{_(t.deleteResourceMessage)}</>}
                confirmButtonText={_(t.delete)}
                loading={loadingResource}
                onClose={closeDeleteModal}
                confirmDisabled={ttlHasError}
            >
                {_(t.delete)}
            </ConfirmModal>

            <DevicesDPSModal
                show={showDpsModal}
                onClose={() => setShowDpsModal(false)}
                updateResource={updateResource}
                resources={resources}
            />
        </Layout>
    )
}

export default DevicesDetailsPage
