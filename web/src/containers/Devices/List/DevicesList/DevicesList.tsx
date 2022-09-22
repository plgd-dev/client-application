import { FC, useMemo } from 'react'
import { useIntl } from 'react-intl'
import { Link, useHistory } from 'react-router-dom'
import classNames from 'classnames'
import Button from '@shared-ui/components/new/Button'
import Badge from '@shared-ui/components/new/Badge'
import Table from '@shared-ui/components/new/Table'
import DevicesListActionButton from '../DevicesListActionButton'
import { DEVICES_DEFAULT_PAGE_SIZE, NO_DEVICE_NAME, devicesOwnerships, DEVICE_TYPE_OIC_WK_D } from '../../constants'
import { messages as t } from '../../Devices.i18n'
import { Props, defaultProps } from './DevicesList.types'

const { OWNED, UNOWNED } = devicesOwnerships

const DevicesList: FC<Props> = (props) => {
    const {
        data,
        loading,
        setSelectedDevices,
        onDeleteClick,
        unselectRowsToken,
        ownDevice,
        showDpsModal,
        resourcesLoadedCallback,
    } = { ...defaultProps, ...props }
    const { formatMessage: _ } = useIntl()
    const history = useHistory()

    const columns = useMemo(
        () => [
            {
                Header: _(t.name),
                accessor: 'data.content.n',
                Cell: ({ value, row }: { value: any; row: any }) => {
                    const deviceName = value || NO_DEVICE_NAME
                    return (
                        <Link to={`/devices/${row.original?.id}`}>
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
                        .map((i: string) => <Badge key={i}>{i}</Badge>)
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
                            onView={(deviceId) => history.push(`/devices/${deviceId}`)}
                            onDelete={onDeleteClick}
                            isOwned={isOwned}
                            onOwnChange={() => ownDevice(isOwned, id, data.content.name)}
                            showDpsModal={showDpsModal}
                            resourcesLoadedCallback={resourcesLoadedCallback}
                        />
                    )
                },
                className: 'actions',
            },
        ],
        [] // eslint-disable-line
    )

    const validData = (data: any) => (!data || data[0] === undefined ? [] : data)

    return (
        <Table
            className='with-selectable-rows'
            columns={columns}
            data={validData(data)}
            defaultSortBy={[
                {
                    id: 'name',
                    desc: false,
                },
            ]}
            autoFillEmptyRows
            defaultPageSize={DEVICES_DEFAULT_PAGE_SIZE}
            getRowProps={(row) => ({
                className: classNames({
                    'grayed-out': row.original?.status === UNOWNED,
                }),
            })}
            getColumnProps={(column) => {
                if (column.id === 'actions') {
                    return { style: { textAlign: 'center' } }
                }

                return {}
            }}
            primaryAttribute='id'
            onRowsSelect={setSelectedDevices}
            bottomControls={
                <Button onClick={onDeleteClick} variant='secondary' icon='fa-trash-alt' disabled={loading}>
                    {_(t.flushCache)}
                </Button>
            }
            unselectRowsToken={unselectRowsToken}
        />
    )
}

DevicesList.displayName = 'DevicesList'
DevicesList.defaultProps = defaultProps

export default DevicesList
