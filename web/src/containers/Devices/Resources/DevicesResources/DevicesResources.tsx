import { FC, useMemo } from 'react'
import { useIntl } from 'react-intl'
import classNames from 'classnames'

import Switch from '@shared-ui/components/Atomic/Switch'
import { useLocalStorage } from '@shared-ui/common/hooks'
import DevicesResourcesList from '@shared-ui/components/Organisms/DevicesResourcesList'
import TableActionButton from '@shared-ui/components/Organisms/TableActionButton'
import DevicesResourcesTree from '@shared-ui/components/Organisms/DevicesResourcesTree'
import TreeExpander from '@shared-ui/components/Atomic/TreeExpander'
import Badge from '@shared-ui/components/Atomic/Badge'
import { IconEdit, IconPlus, IconTrash } from '@shared-ui/components/Atomic/Icon'

import { devicesStatuses, RESOURCE_TREE_DEPTH_SIZE } from '../../constants'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesResources.types'
import { canBeResourceEdited, canCreateResource, getLastPartOfAResourceHref } from '@/containers/Devices/utils'

export const DevicesResources: FC<Props> = ({
    data,
    deviceId,
    deviceStatus,
    isActiveTab,
    isOwned,
    loading,
    onCreate,
    onDelete,
    onUpdate,
    pageSize,
}) => {
    const { formatMessage: _ } = useIntl()
    const [treeViewActive, setTreeViewActive] = useLocalStorage('treeViewActive', false)
    const isUnregistered = devicesStatuses.UNREGISTERED === deviceStatus
    const greyedOutClassName = classNames({
        'grayed-out': isUnregistered,
    })

    const columns = useMemo(
        () => [
            {
                Header: _(t.href),
                accessor: 'href',
                Cell: ({ value, row }: { value: any; row: any }) => {
                    const {
                        original: { deviceId: deviceIdOrigin, href, endpointInformations },
                    } = row

                    const edit = canBeResourceEdited(endpointInformations) || isOwned

                    if (!edit) {
                        return <span>{value}</span>
                    }
                    return (
                        <div className='tree-expander-container'>
                            <span
                                className='link reveal-icon-on-hover'
                                onClick={() => onUpdate({ deviceId: deviceIdOrigin, href })}
                            >
                                {value}
                            </span>
                        </div>
                    )
                },
                style: { width: '300px' },
            },
            {
                Header: _(t.types),
                accessor: 'resourceTypes',
                style: { width: '100%' },
                Cell: ({ value }: { value: any }) => value.join(', '),
            },
            {
                Header: _(t.actions),
                accessor: 'actions',
                disableSortBy: true,
                style: { textAlign: 'center' },
                Cell: ({ row }: { row: any }) => {
                    const {
                        original: { href, interfaces, endpointInformations },
                    } = row

                    const create = interfaces && canCreateResource(interfaces) && isOwned
                    const edit = (endpointInformations && canBeResourceEdited(endpointInformations)) || isOwned

                    return (
                        <TableActionButton
                            disabled={loading}
                            items={[
                                {
                                    onClick: () => onCreate(href),
                                    label: _(t.create),
                                    icon: <IconPlus />,
                                    hidden: !create,
                                },
                                {
                                    onClick: () => onUpdate({ deviceId, href }),
                                    label: _(t.update),
                                    icon: <IconEdit />,
                                    hidden: !edit,
                                },
                                {
                                    onClick: () => onDelete(href),
                                    label: _(t.delete),
                                    icon: <IconTrash />,
                                    hidden: !edit,
                                },
                            ]}
                        />
                    )
                },
                className: 'actions',
            },
        ],
        [onUpdate, onCreate, onDelete, loading] //eslint-disable-line
    )

    const treeColumns = useMemo(
        () => [
            {
                Header: _(t.href),
                accessor: 'href',
                Cell: ({ value, row }: { value: any; row: any }) => {
                    const {
                        original: { href, endpointInformations },
                    } = row

                    const lastValue = getLastPartOfAResourceHref(value)
                    const onLinkClick = deviceId ? () => onUpdate({ deviceId, href: href.replace(/\/$/, '') }) : null

                    const edit = canBeResourceEdited(endpointInformations) || isOwned

                    if (!edit) {
                        return <span>{lastValue}</span>
                    }

                    if (row.canExpand) {
                        return (
                            <div className='tree-expander-container'>
                                <TreeExpander
                                    {...row.getToggleRowExpandedProps({ title: null })}
                                    expanded={row.isExpanded}
                                    style={{
                                        marginLeft: `${row.depth * RESOURCE_TREE_DEPTH_SIZE}px`,
                                    }}
                                />
                                <span
                                    className={!row.canExpand ? 'link reveal-icon-on-hover' : ''}
                                    onClick={() => onLinkClick}
                                >
                                    {`/${lastValue}/`}
                                </span>
                                {!row.canExpand && <i className='fas fa-pen' />}
                            </div>
                        )
                    }

                    return (
                        <div
                            className='tree-expander-container'
                            style={{
                                marginLeft: `${row.depth === 0 ? 0 : (row.depth + 1) * RESOURCE_TREE_DEPTH_SIZE}px`,
                            }}
                        >
                            <span className='link reveal-icon-on-hover' onClick={() => onLinkClick}>
                                {`/${lastValue}`}
                            </span>
                            <i className='fas fa-pen' />
                        </div>
                    )
                },
                style: { width: '100%' },
            },
            {
                Header: _(t.types),
                accessor: 'resourceTypes',
                Cell: ({ value }: { value: any }) => {
                    if (!deviceId) {
                        return null
                    }

                    return (
                        <div className='badges-box-horizontal'>
                            {value?.map?.((type: string) => (
                                <Badge key={type}>{type}</Badge>
                            ))}
                        </div>
                    )
                },
            },
            {
                Header: _(t.actions),
                accessor: 'actions',
                disableSortBy: true,
                Cell: ({ row }: { row: any }) => {
                    if (row.canExpand) {
                        return null
                    }

                    const {
                        original: { href, interfaces, endpointInformations },
                    } = row
                    const cleanHref = href.replace(/\/$/, '') // href without a trailing slash
                    const create = interfaces && canCreateResource(interfaces) && isOwned
                    const edit = (endpointInformations && canBeResourceEdited(endpointInformations)) || isOwned

                    return (
                        <TableActionButton
                            disabled={loading}
                            items={[
                                {
                                    onClick: () => onCreate(cleanHref),
                                    label: _(t.create),
                                    icon: <IconPlus />,
                                    hidden: !create,
                                },
                                {
                                    onClick: () => onUpdate({ deviceId, href: cleanHref }),
                                    label: _(t.update),
                                    icon: <IconEdit />,
                                    hidden: !edit,
                                },
                                {
                                    onClick: () => onDelete(cleanHref),
                                    label: _(t.delete),
                                    icon: <IconTrash />,
                                    hidden: !edit,
                                },
                            ]}
                        />
                    )
                },
            },
        ],
        [onUpdate, onCreate, onDelete, loading] //eslint-disable-line
    )

    return (
        <>
            <div className={classNames('d-flex justify-content-between align-items-center', greyedOutClassName)}>
                <div />
                <div className='d-flex justify-content-end align-items-center' style={{ paddingBottom: 24 }}>
                    <Switch
                        checked={treeViewActive}
                        disabled={isUnregistered}
                        id='toggle-tree-view'
                        label={_(t.treeView)}
                        onChange={() => setTreeViewActive(!treeViewActive)}
                    />
                </div>
            </div>

            {treeViewActive ? (
                <DevicesResourcesTree columns={treeColumns} data={data} />
            ) : (
                <DevicesResourcesList
                    columns={columns}
                    data={data}
                    i18n={{
                        search: _(t.search),
                    }}
                    isActiveTab={isActiveTab}
                    pageSize={pageSize}
                />
            )}
        </>
    )
}

DevicesResources.displayName = 'DevicesResources'

export default DevicesResources
