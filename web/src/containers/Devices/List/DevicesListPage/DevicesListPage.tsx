import { useEffect, useState } from 'react'
import { useIntl } from 'react-intl'
import { showSuccessToast } from '@shared-ui/components/new/Toast/Toast'
import ConfirmModal from '@shared-ui/components/new/ConfirmModal'
import Layout from '@shared-ui/components/new/Layout'
import { useIsMounted } from '@shared-ui/common/hooks'
import { messages as menuT } from '@shared-ui/components/new/Menu/Menu.i18n'
import { useDevicesList } from '../../hooks'
import DevicesList from '../DevicesList'
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
import { toast } from 'react-toastify'
import { getApiErrorMessage } from '@shared-ui/common/utils'
import { getDevices, updateDevices, flushDevices, ownDevice, disOwnDevice } from '@/containers/Devices/slice'
import { useDispatch, useSelector } from 'react-redux'
import DevicesTimeoutModal from '../DevicesTimeoutModal'
import DevicesDPSModal from '../../DevicesDPSModal'
import { DeviceDataType, ResourcesType } from '@/containers/Devices/Devices.types'
import { DpsDataType } from '@/containers/Devices/List/DevicesListPage/DevicesListPage.types'
import { DevicesResourcesModalParamsType } from '@/containers/Devices/Resources/DevicesResourcesModal/DevicesResourcesModal.types'

const DevicesListPage = () => {
    const { formatMessage: _ } = useIntl()
    const { data, loading, error: deviceError, refresh } = useDevicesList()
    const [deleteModalOpen, setDeleteModalOpen] = useState(false)
    const [timeoutModalOpen, setTimeoutModalOpen] = useState(false)
    const [selectedDevices, setSelectedDevices] = useState([])
    const [deleting, setDeleting] = useState(false)
    const [owning, setOwning] = useState(false)
    const [showDpsModal, setShowDpsModal] = useState(false)
    const [dpsData, setDpsData] = useState<DpsDataType>({
        deviceId: '',
        resources: undefined,
    })
    const isMounted = useIsMounted()
    const dispatch = useDispatch()
    const dataToDisplay: DeviceDataType = useSelector(getDevices)

    useEffect(() => {
        deviceError && toast.error(getApiErrorMessage(deviceError))
    }, [deviceError])

    useEffect(() => {
        // @ts-ignore
        data && dispatch(updateDevices(data))
    }, [data, dispatch])

    const handleOpenDeleteModal = () => {
        setDeleteModalOpen(true)
    }

    const handleCloseDeleteModal = () => {
        setDeleteModalOpen(false)
    }

    const handleRefresh = () => {
        refresh()
    }

    const deleteDevices = async () => {
        try {
            setDeleting(true)
            await deleteDevicesApi()
            await sleep(200)

            if (isMounted.current) {
                showSuccessToast({
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
                showSuccessToast({
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
                showSuccessToast({
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

    const loadingOrDeleting = loading || deleting || owning

    return (
        <Layout
            title={_(menuT.devices)}
            breadcrumbs={[
                {
                    label: _(menuT.devices),
                },
            ]}
            loading={loading || owning}
            header={
                <DevicesListHeader
                    loading={loading || owning}
                    refresh={handleRefresh}
                    openTimeoutModal={() => setTimeoutModalOpen(true)}
                />
            }
        >
            <DevicesList
                data={dataToDisplay}
                selectedDevices={selectedDevices}
                setSelectedDevices={setSelectedDevices}
                loading={loadingOrDeleting}
                onDeleteClick={handleOpenDeleteModal}
                ownDevice={handleOwnDevice}
                showDpsModal={(deviceId: string) => {
                    setDpsData((prevData: DpsDataType) => ({ ...prevData, deviceId }))
                    setShowDpsModal(true)
                }}
                resourcesLoadedCallback={(resources) => {
                    setDpsData((prevData: DpsDataType) => ({
                        ...prevData,
                        resources,
                    }))
                }}
            />

            <ConfirmModal
                onConfirm={deleteDevices}
                show={deleteModalOpen}
                title={
                    <>
                        <i className='fas fa-trash-alt' />
                        {_(t.flushDevices)}
                    </>
                }
                body={_(t.flushDevicesMessage)}
                confirmButtonText={_(t.flushCache)}
                loading={loadingOrDeleting}
                onClose={handleCloseDeleteModal}
            >
                {_(t.flushCache)}
            </ConfirmModal>

            <DevicesTimeoutModal show={timeoutModalOpen} onClose={() => setTimeoutModalOpen(false)} />

            <DevicesDPSModal
                show={showDpsModal}
                onClose={() => setShowDpsModal(false)}
                updateResource={updateResource}
                resources={dpsData.resources as ResourcesType[]}
            />
        </Layout>
    )
}

DevicesListPage.displayName = 'DevicesListPage'

export default DevicesListPage
