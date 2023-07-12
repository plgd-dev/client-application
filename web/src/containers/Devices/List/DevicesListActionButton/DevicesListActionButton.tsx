import { FC, memo, useState } from 'react'
import { useIntl } from 'react-intl'
import isFunction from 'lodash/isFunction'

import TableActionButton from '@shared-ui/components/Organisms/TableActionButton'
import { ItemType } from '@shared-ui/components/Organisms/TableActionButton/TableActionButton.types'
import { IconClose, IconNetwork, IconPlus, IconShowPassword } from '@shared-ui/components/Atomic/Icon'

import { messages as t } from '../../Devices.i18n'
import { getDevicesResourcesAllApi } from '@/containers/Devices/rest'
import { canSetDPSEndpoint } from '@/containers/Devices/utils'
import { Props } from './DevicesListActionButton.types'

const DevicesListActionButton: FC<Props> = memo(
    ({ deviceId, onView, onDelete, isOwned, onOwnChange, showDpsModal, resourcesLoadedCallback }) => {
        const getDefaultItems = () => {
            const defaultItems: ItemType[] = [
                {
                    id: 'detail',
                    onClick: () => onView(deviceId),
                    label: _(t.details),
                    icon: <IconShowPassword />,
                },
                {
                    id: 'own',
                    onClick: () => onOwnChange(),
                    label: isOwned ? _(t.disOwnDevice) : _(t.ownDevice),
                    icon: isOwned ? <IconClose /> : <IconPlus />,
                },
            ]

            if (isOwned) {
                defaultItems.push({
                    id: 'dps',
                    onClick: () => showDpsModal(deviceId),
                    label: _(t.setDpsEndpoint),
                    icon: <IconNetwork />,
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
