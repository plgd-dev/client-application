import { useIntl } from 'react-intl'
import ActionButton from '@shared-ui/components/new/ActionButton'
import { messages as t } from '../../Devices.i18n'
import { FC, useState } from 'react'
import { getDevicesResourcesAllApi } from '@/containers/devices/rest'
import { canSetDPSEndpoint } from '@/containers/devices/utils'
import isFunction from 'lodash/isFunction'
import { Props } from './DevicesListActionButton.types'
import { ActionButtonItemType } from '@shared-ui/components/new/ActionButton/ActionButton.types'

const DevicesListActionButton: FC<Props> = ({
  deviceId,
  onView,
  isOwned,
  onOwnChange,
  showDpsModal,
  resourcesLoadedCallback,
}) => {
  const getDefaultItems = () => {
    const defaultItems: ActionButtonItemType[] = [
      {
        id: 'detail',
        onClick: () => onView(deviceId),
        label: _(t.details),
        icon: 'fa-eye',
      },
      {
        id: 'own',
        onClick: () => onOwnChange(),
        label: isOwned ? _(t.disOwnDevice) : _(t.ownDevice),
        icon: isOwned ? 'fa-cloud-download-alt' : 'fa-cloud-upload-alt',
      },
    ]

    if (isOwned) {
      defaultItems.push({
        id: 'dps',
        onClick: () => showDpsModal(deviceId),
        label: _(t.setDpsEndpoint),
        icon: 'fa-bacon',
        loading: true,
      })
    }

    return defaultItems
  }
  const { formatMessage: _ } = useIntl()
  const [resources, setResources] = useState(undefined)
  const [items, setItems] = useState(getDefaultItems())

  const handleToggle = async (isOpen: boolean) => {
    if (isOpen && isOwned && !resources) {
      const { data } = await getDevicesResourcesAllApi(deviceId)

      setResources(data.resources)
      isFunction(resourcesLoadedCallback) &&
        resourcesLoadedCallback(data.resources)
      const hasDPS = canSetDPSEndpoint(data.resources)

      setItems(() => {
        const newItems: any = []
        items.forEach(item => {
          if (item.id === 'dps') {
            if (hasDPS) {
              newItems.push({
                ...item,
                loading: false,
              })
            }
          } else {
            newItems.push(item)
          }
        })

        return newItems
      })
    }
  }

  return (
    <ActionButton
      onToggle={handleToggle as any}
      menuProps={{
        align: 'end',
      }}
      items={items}
    >
      <i className="fas fa-ellipsis-h" />
    </ActionButton>
  )
}

DevicesListActionButton.displayName = 'DevicesListActionButton'

export default DevicesListActionButton
