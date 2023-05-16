import { FC, memo, useState } from 'react'
import { useIntl } from 'react-intl'
import isFunction from 'lodash/isFunction'

import TableActionButton from '@shared-ui/components/organisms/TableActionButton'
import { ItemType } from '@shared-ui/components/organisms/TableActionButton/TableActionButton.types'

import { messages as t } from '../../Devices.i18n'
import { getDevicesResourcesAllApi } from '@/containers/Devices/rest'
import { canSetDPSEndpoint } from '@/containers/Devices/utils'
import { Props } from './DevicesListActionButton.types'

const DevicesListActionButton: FC<Props> = memo(
    ({ deviceId, onView, isOwned, onOwnChange, showDpsModal, resourcesLoadedCallback }) => {
        const getDefaultItems = () => {
            const defaultItems: ItemType[] = [
                {
                    id: 'detail',
                    onClick: () => onView(deviceId),
                    label: _(t.details),
                    icon: 'icon-show-password',
                },
                {
                    id: 'delete',
                    onClick: () => onView(deviceId),
                    label: _(t.delete),
                    icon: 'trash',
                },
                {
                    id: 'own',
                    onClick: () => onOwnChange(),
                    label: isOwned ? _(t.disOwnDevice) : _(t.ownDevice),
                    icon: isOwned ? 'close' : 'plus',
                },
            ]

            if (isOwned) {
                defaultItems.push({
                    id: 'dps',
                    onClick: () => showDpsModal(deviceId),
                    label: _(t.setDpsEndpoint),
                    icon: 'network',
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
                isFunction(resourcesLoadedCallback) && resourcesLoadedCallback(data.resources)
                const hasDPS = canSetDPSEndpoint(data.resources)

                setItems(() => {
                    const newItems: ItemType[] = []
                    items.forEach((item) => {
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

        return <TableActionButton items={items} onToggle={handleToggle} />
    }
)

DevicesListActionButton.displayName = 'DevicesListActionButton'

export default DevicesListActionButton
