import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'
import { ActionButton } from '@/components/action-button'
import { messages as t } from './devices-i18n'
import { useState } from 'react'
import { getDevicesResourcesAllApi } from '@/containers/devices/rest'
import { canSetDPSEndpoint } from '@/containers/devices/utils'
import isFunction from 'lodash/isFunction'

export const DevicesListActionButton = ({
  deviceId,
  onView,
  isOwned,
  onOwnChange,
  showDpsModal,
  resourcesLoadedCallback,
}) => {
  const { formatMessage: _ } = useIntl()
  const [resources, setResources] = useState(undefined)
  const [items, setItems] = useState([
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
    {
      id: 'dps',
      onClick: () => showDpsModal(deviceId),
      label: _(t.setDpsEndpoint),
      icon: 'fa-bacon',
      loading: true,
    },
  ])

  const handleToggle = async isOpen => {
    if (isOpen && !resources) {
      const { data } = await getDevicesResourcesAllApi(deviceId)

      setResources(data.resources)
      isFunction(resourcesLoadedCallback) &&
        resourcesLoadedCallback(data.resources)
      const hasDPS = canSetDPSEndpoint(data.resources)

      setItems(() => {
        const newItems = []
        items.map(item => {
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
      onToggle={handleToggle}
      menuProps={{
        align: 'right',
      }}
      items={items}
    >
      <i className="fas fa-ellipsis-h" />
    </ActionButton>
  )
}

DevicesListActionButton.propTypes = {
  deviceId: PropTypes.string.isRequired,
  onView: PropTypes.func.isRequired,
}
