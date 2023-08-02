import { useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { useIntl } from 'react-intl'
import { toast } from 'react-toastify'
import { useDispatch, useSelector } from 'react-redux'
import { Link, useNavigate } from 'react-router-dom'

import Notification from '@shared-ui/components/Atomic/Notification/Toast'
import ConfirmModal from '@shared-ui/components/Atomic/ConfirmModal'
import PageLayout from '@shared-ui/components/Atomic/PageLayout'
import { useIsMounted } from '@shared-ui/common/hooks'
import { messages as menuT } from '@shared-ui/components/Atomic/Menu/Menu.i18n'
import { getApiErrorMessage } from '@shared-ui/common/utils'
import Footer from '@shared-ui/components/Layout/Footer'
import DevicesList from '@shared-ui/components/Organisms/DevicesList/DevicesList'
import { DevicesResourcesModalParamsType } from '@shared-ui/components/Organisms/DevicesResourcesModal/DevicesResourcesModal.types'
import Badge from '@shared-ui/components/Atomic/Badge'
import Tag from '@shared-ui/components/Atomic/Tag'

import { useDevicesList } from '../../hooks'
import DevicesListHeader from '../DevicesListHeader'
import { deleteDevicesApi, disownDeviceApi, ownDeviceApi } from '../../rest'
import {
    handleDeleteDevicesErrors,
    handleOwnDevicesErrors,
    handleUpdateResourceErrors,
    sleep,
    updateResourceMethod,
} from '../../utils'
import { messages as t } from '../../Devices.i18n'
import { getDevices, updateDevices, flushDevices, ownDevice, disOwnDevice } from '@/containers/Devices/slice'
import DevicesTimeoutModal from '../DevicesTimeoutModal'
import DevicesDPSModal from '../../DevicesDPSModal'
import { DeviceDataType, ResourcesType } from '@/containers/Devices/Devices.types'
import { DpsDataType } from '@/containers/Devices/List/DevicesListPage/DevicesListPage.types'
import { DEVICE_TYPE_OIC_WK_D, devicesOwnerships, NO_DEVICE_NAME } from '@/containers/Devices/constants'
import DevicesListActionButton from '@/containers/Devices/List/DevicesListActionButton'
import AppContext from '@/containers/App/AppContext'

const { OWNED, UNSUPPORTED } = devicesOwnerships

const DevicesListPage = () => {
    const { formatMessage: _ } = useIntl()
    const { data, loading, error: deviceError, refresh } = useDevicesList()
    const [deleteModalOpen, setDeleteModalOpen] = useState(false)
    const [timeoutModalOpen, setTimeoutModalOpen] = useState(false)
    const [selectedDevices, setSelectedDevices] = useState([])
    const [isAllSelected, setIsAllSelected] = useState(false)
    const [deleting, setDeleting] = useState(false)
    const [owning, setOwning] = useState(false)
    const [showDpsModal, setShowDpsModal] = useState(false)
    const [singleDevice, setSingleDevice] = useState<null | string>(null)
    const [dpsData, setDpsData] = useState<DpsDataType>({
        deviceId: '',
        resources: undefined,
    })
    const isMounted = useIsMounted()
    const dispatch = useDispatch()
    const navigate = useNavigate()
    const dataToDisplay: DeviceDataType = useSelector(getDevices)
    const { collapsed, iframeMode } = useContext(AppContext)

    useEffect(() => {
        deviceError && toast.error(getApiErrorMessage(deviceError))
    }, [deviceError])

    useEffect(() => {
        // @ts-ignore
        data && dispatch(updateDevices(data))
    }, [data, dispatch])

    const handleOpenDeleteModal = useCallback(
        (deviceId?: string) => {
            if (typeof deviceId === 'string') {
                setSingleDevice(deviceId)
            } else if (singleDevice && !deviceId) {
                setSingleDevice(null)
            }

            setDeleteModalOpen(true)
        },
        [singleDevice]
    )

    const handleCloseDeleteModal = useCallback(() => {
        setDeleteModalOpen(false)
    }, [])

    const handleRefresh = useCallback(() => {
        refresh()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const deleteDevices = async () => {
        try {
            setDeleting(true)
            await deleteDevicesApi()
            await sleep(200)

            if (isMounted.current) {
                Notification.success({
                    title: _(t.devicesDeleted),
                    message: _(t.devicesDeletedMessage),
                })

                // @ts-ignore
                dispatch(flushDevices(data))

                setDeleting(false)
                setDeleteModalOpen(false)
                handleCloseDeleteModal()
            }
        } catch (error) {
            setDeleting(false)
            handleDeleteDevicesErrors(error, _)
        }
    }

    const handleOwnDevice = async (isOwned: boolean, deviceId: string, deviceName: string) => {
        try {
            setOwning(true)
            isOwned ? await disownDeviceApi(deviceId) : await ownDeviceApi(deviceId)

            if (isMounted.current) {
                Notification.success({
                    title: isOwned ? _(t.deviceDisOwned) : _(t.deviceOwned),
                    message: isOwned
                        ? _(t.deviceWasDisOwned, { name: deviceName })
                        : _(t.deviceWasOwned, { name: deviceName }),
                })

                if (!isOwned) {
                    // @ts-ignore
                    dispatch(ownDevice(deviceId))
                } else {
                    // @ts-ignore
                    dispatch(disOwnDevice(deviceId))
                    refresh()
                }

                setOwning(false)
            }
        } catch (error) {
            handleOwnDevicesErrors(error, _)
            refresh()
            setOwning(false)
        }
    }

    // Updates the resource through rest API
    const updateResource = async (
        { href, currentInterface = '' }: DevicesResourcesModalParamsType,
        resourceDataUpdate: any
    ) => {
        await updateResourceMethod(
            { deviceId: dpsData.deviceId, href, currentInterface },
            resourceDataUpdate,
            () => {
                Notification.success({
                    title: _(t.resourceUpdateSuccess),
                    message: _(t.resourceWasUpdated),
                })
                setShowDpsModal(false)
                setDpsData({ deviceId: '', resources: undefined })
            },
            (error: any) => {
                handleUpdateResourceErrors(error, { id: dpsData.deviceId, href }, _)
            }
        )
    }

    const handleFlashDevices = () => {
        // @ts-ignore
        dispatch(flushDevices())
    }

    const loadingOrOwning = useMemo(() => loading || owning, [loading, owning])
    const loadingOrDeletingOrOwning = useMemo(() => loadingOrOwning || deleting, [loadingOrOwning, deleting])

    const handleOpenTimeoutModal = useCallback(() => {
        setTimeoutModalOpen(true)
    }, [])

    const columns = useMemo(
        () => [
            {
                Header: _(t.name),
                accessor: 'data.content.n',
                Cell: ({ value, row }: { value: any; row: any }) => {
                    const deviceName = value || NO_DEVICE_NAME
                    return (
                        <Link data-test-id={deviceName} to={`/devices/${row.original?.id}`}>
                            <span className='no-wrap-text'>{deviceName}</span>
                        </Link>
                    )
                },
                style: { width: '100%' },
            },
            {
                Header: 'Types',
                accessor: 'types',
                style: { maxWidth: '350px', width: '100%' },
                Cell: ({ value }: { value: any }) => {
                    if (!value) {
                        return null
                    }
                    return value
                        .filter((i: string) => i !== DEVICE_TYPE_OIC_WK_D)
                        .map((i: string) => <Tag key={i}>{i}</Tag>)
                },
            },
            {
                Header: 'ID',
                accessor: 'id',
                style: { maxWidth: '350px', width: '100%' },
                Cell: ({ value }: { value: any }) => {
                    return <span className='no-wrap-text'>{value}</span>
                },
            },
            {
                Header: _(t.ownershipStatus),
                accessor: 'ownershipStatus',
                style: { width: '250px' },
                Cell: ({ value }: { value: any }) => {
                    const isOwned = OWNED === value

                    if (UNSUPPORTED === value) {
                        return <Badge className='grey'>{_(t.unsupported)}</Badge>
                    }

                    return <Badge className={isOwned ? 'green' : 'red'}>{isOwned ? _(t.owned) : _(t.unowned)}</Badge>
                },
            },
            {
                Header: _(t.actions),
                accessor: 'actions',
                disableSortBy: true,
                Cell: ({ row }: { row: any }) => {
                    const {
                        original: { id, ownershipStatus, data },
                    } = row
                    const isOwned = ownershipStatus === OWNED

                    return (
                        <DevicesListActionButton
                            deviceId={id}
                            onDelete={deleteDevices}
                            onOwnChange={() => handleOwnDevice(isOwned, id, data.content.name)}
                            onView={(deviceId) => navigate(`/devices/${deviceId}`)}
                            ownershipStatus={ownershipStatus}
                            resourcesLoadedCallback={(resources) => {
                                setDpsData((prevData: DpsDataType) => ({
                                    ...prevData,
                                    resources,
                                }))
                            }}
                            showDpsModal={(deviceId: string) => {
                                setDpsData((prevData: DpsDataType) => ({ ...prevData, deviceId }))
                                setShowDpsModal(true)
                            }}
                        />
                    )
                },
                className: 'actions',
            },
        ],
        [loading] // eslint-disable-line
    )

    return (
        <PageLayout
            breadcrumbs={[
                {
                    label: _(menuT.devices),
                },
            ]}
            footer={<Footer footerExpanded={false} paginationComponent={<div id='paginationPortalTarget'></div>} />}
            header={
                <DevicesListHeader
                    handleFlashDevices={handleFlashDevices}
                    i18n={{
                        flushCache: _(t.flushCache),
                    }}
                    loading={loadingOrOwning}
                    openTimeoutModal={handleOpenTimeoutModal}
                    refresh={handleRefresh}
                />
            }
            loading={loading || owning}
            title={_(menuT.devices)}
        >
            <DevicesList
                collapsed={collapsed ?? false}
                columns={columns}
                data={dataToDisplay}
                i18n={{
                    delete: _(t.delete),
                    search: _(t.search),
                    select: _(t.select),
                }}
                iframeMode={iframeMode}
                isAllSelected={isAllSelected}
                loading={loadingOrDeletingOrOwning}
                onDeleteClick={handleOpenDeleteModal}
                selectedDevices={selectedDevices}
                setIsAllSelected={setIsAllSelected}
                setSelectedDevices={setSelectedDevices}
            />

            <ConfirmModal
                body={_(t.flushDevicesMessage)}
                confirmButtonText={_(t.flushCache)}
                loading={loadingOrDeletingOrOwning}
                onClose={handleCloseDeleteModal}
                onConfirm={deleteDevices}
                show={deleteModalOpen}
                title={<>{_(t.flushDevices)}</>}
            >
                {_(t.flushCache)}
            </ConfirmModal>

            <DevicesTimeoutModal onClose={() => setTimeoutModalOpen(false)} show={timeoutModalOpen} />

            <DevicesDPSModal
                onClose={() => setShowDpsModal(false)}
                resources={dpsData.resources as ResourcesType[]}
                show={showDpsModal}
                updateResource={updateResource}
            />
        </PageLayout>
    )
}

DevicesListPage.displayName = 'DevicesListPage'

export default DevicesListPage
