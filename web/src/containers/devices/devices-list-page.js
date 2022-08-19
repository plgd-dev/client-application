import { useEffect, useState } from 'react'
import { useIntl } from 'react-intl'
import { showSuccessToast } from '@shared-ui/components/old/toast'
import ConfirmModal from '@shared-ui/components/new/ConfirmModal'
import { Layout } from '@shared-ui/components/old/layout'
import { useIsMounted } from '@shared-ui/common/hooks'
import { messages as menuT } from '@shared-ui/components/old/menu/menu-i18n'
import { useDevicesList } from './hooks'
import { DevicesList } from './_devices-list'
import { DevicesListHeader } from './_devices-list-header'
import { deleteDevicesApi, disownDeviceApi, ownDeviceApi } from './rest'
import {
  handleDeleteDevicesErrors,
  handleOwnDevicesErrors,
  handleUpdateResourceErrors,
  sleep,
  updateResourceMethod,
} from './utils'
import { messages as t } from './devices-i18n'
import { toast } from 'react-toastify'
import { getApiErrorMessage } from '@shared-ui/common/utils'
import {
  getDevices,
  updateDevices,
  flushDevices,
  ownDevice,
  disOwnDevice,
} from '@/containers/devices/slice'
import { useDispatch, useSelector } from 'react-redux'
import { DevicesTimeoutModal } from './_devices-timeout-modal'
import { DevicesDPSModal } from '@/containers/devices/_devices-dps-modal'

export const DevicesListPage = () => {
  const { formatMessage: _ } = useIntl()
  const { data, loading, error: deviceError, refresh } = useDevicesList()
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [timeoutModalOpen, setTimeoutModalOpen] = useState(false)
  const [selectedDevices, setSelectedDevices] = useState([])
  const [deleting, setDeleting] = useState(false)
  const [owning, setOwning] = useState(false)
  const [showDpsModal, setShowDpsModal] = useState(false)
  const [dpsData, setDpsData] = useState({
    deviceId: undefined,
    resources: undefined,
  })
  const isMounted = useIsMounted()
  const dispatch = useDispatch()
  const dataToDisplay = useSelector(getDevices)

  useEffect(() => {
    deviceError && toast.error(getApiErrorMessage(deviceError))
  }, [deviceError])

  useEffect(() => {
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

  const handleOwnDevice = async (isOwned, deviceId, deviceName) => {
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
          dispatch(ownDevice(deviceId))
        } else {
          dispatch(disOwnDevice(deviceId))
          refresh()
        }

        setOwning(false)
      }
    } catch (error) {
      handleOwnDevicesErrors(error, _, true)
      refresh()
      setOwning(false)
    }
  }

  // Updates the resource through rest API
  const updateResource = async (
    { href, currentInterface = '' },
    resourceDataUpdate
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
        setDpsData({ deviceId: undefined, resources: undefined })
      },
      error => {
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
        showDpsModal={deviceId => {
          setDpsData(prevData => ({ ...prevData, deviceId }))
          setShowDpsModal(true)
        }}
        resourcesLoadedCallback={resources => {
          setDpsData(prevData => ({
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
            <i className="fas fa-trash-alt" />
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

      <DevicesTimeoutModal
        show={timeoutModalOpen}
        onClose={() => setTimeoutModalOpen(false)}
      />

      <DevicesDPSModal
        show={showDpsModal}
        onClose={() => setShowDpsModal(false)}
        updateResource={updateResource}
        resources={dpsData.resources}
      />
    </Layout>
  )
}
