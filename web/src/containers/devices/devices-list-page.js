import { useEffect, useState } from 'react'
import { useIntl } from 'react-intl'
import { showSuccessToast } from '@/components/toast'
import { ConfirmModal } from '@/components/confirm-modal'
import { Layout } from '@/components/layout'
import { useIsMounted } from '@/common/hooks'
import { messages as menuT } from '@/components/menu/menu-i18n'
import { useDevicesList } from './hooks'
import { DevicesList } from './_devices-list'
import { DevicesListHeader } from './_devices-list-header'
import { deleteDevicesApi, disownDevice, ownDevice } from './rest'
import {
  handleDeleteDevicesErrors,
  handleOwnDevicesErrors,
  sleep,
} from './utils'
import { messages as t } from './devices-i18n'
import { toast } from 'react-toastify'
import { getApiErrorMessage } from '@/common/utils'
import {
  getDevices,
  setDevices,
  flushDevices,
  toggleOwnDevice,
} from '@/containers/devices/slice'
import { useDispatch, useSelector } from 'react-redux'

export const DevicesListPage = () => {
  const { formatMessage: _ } = useIntl()
  const { data, loading, error: deviceError, refresh } = useDevicesList()
  const [deleteModalOpen, setDeleteModalOpen] = useState(false)
  const [selectedDevices, setSelectedDevices] = useState([])
  const [deleting, setDeleting] = useState(false)
  const [owning, setOwning] = useState(false)
  const isMounted = useIsMounted()
  const dispatch = useDispatch()

  const dataToDisplay = useSelector(getDevices)

  useEffect(() => {
    deviceError && toast.error(getApiErrorMessage(deviceError))
  }, [deviceError])

  useEffect(() => {
    dispatch(setDevices(data))
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
      isOwned ? await disownDevice(deviceId) : await ownDevice(deviceId)

      if (isMounted.current) {
        showSuccessToast({
          title: isOwned ? _(t.deviceDisOwned) : _(t.deviceOwned),
          message: isOwned
            ? _(t.deviceWasDisOwned, { name: deviceName })
            : _(t.deviceWasOwned, { name: deviceName }),
        })

        if (!isOwned) {
          dispatch(toggleOwnDevice({ deviceId: deviceId, ownState: !isOwned }))
        } else {
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
    </Layout>
  )
}
