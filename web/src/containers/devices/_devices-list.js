import PropTypes from 'prop-types'
import { useMemo } from 'react'
import { useIntl } from 'react-intl'
import { Link, useHistory } from 'react-router-dom'
import classNames from 'classnames'
import { Button } from '@/components/button'
import { Badge } from '@/components/badge'
import { Table } from '@/components/table'
import { DevicesListActionButton } from './_devices-list-action-button'
import {
  DEVICES_DEFAULT_PAGE_SIZE,
  NO_DEVICE_NAME,
  devicesOwnerships,
  DEVICE_TYPE_OIC_WK_D,
} from './constants'
import { deviceShape } from './shapes'
import { messages as t } from './devices-i18n'

const { OWNED, UNOWNED } = devicesOwnerships

export const DevicesList = ({
  data,
  loading,
  setSelectedDevices,
  onDeleteClick,
  unselectRowsToken,
  ownDevice,
  showDpsModal,
  resourcesLoadedCallback,
}) => {
  const { formatMessage: _ } = useIntl()
  const history = useHistory()

  const columns = useMemo(
    () => [
      {
        Header: _(t.name),
        accessor: 'data.content.n',
        Cell: ({ value, row }) => {
          const deviceName = value || NO_DEVICE_NAME
          return (
            <Link to={`/devices/${row.original?.id}`}>
              <span className="no-wrap-text">{deviceName}</span>
            </Link>
          )
        },
        style: { width: '100%' },
      },
      {
        Header: 'Types',
        accessor: 'types',
        style: { maxWidth: '350px', width: '100%' },
        Cell: ({ value }) => {
          if (!value) {
            return null
          }
          return value
            .filter(i => i !== DEVICE_TYPE_OIC_WK_D)
            .map(i => <Badge key={i}>{i}</Badge>)
        },
      },
      {
        Header: 'ID',
        accessor: 'id',
        style: { maxWidth: '350px', width: '100%' },
        Cell: ({ value }) => {
          return <span className="no-wrap-text">{value}</span>
        },
      },
      {
        Header: _(t.ownershipStatus),
        accessor: 'ownershipStatus',
        style: { width: '250px' },
        Cell: ({ value }) => {
          const isOwned = OWNED === value
          return (
            <Badge className={isOwned ? 'green' : 'red'}>
              {isOwned ? _(t.owned) : _(t.unowned)}
            </Badge>
          )
        },
      },
      {
        Header: _(t.actions),
        accessor: 'actions',
        disableSortBy: true,
        Cell: ({ row }) => {
          const {
            original: { id, ownershipStatus, data },
          } = row
          const isOwned = ownershipStatus === OWNED
          return (
            <DevicesListActionButton
              deviceId={id}
              onView={deviceId => history.push(`/devices/${deviceId}`)}
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

  return (
    <Table
      className="with-selectable-rows"
      columns={columns}
      data={data || []}
      defaultSortBy={[
        {
          id: 'name',
          desc: false,
        },
      ]}
      autoFillEmptyRows
      defaultPageSize={DEVICES_DEFAULT_PAGE_SIZE}
      getRowProps={row => ({
        className: classNames({
          'grayed-out': row.original?.status === UNOWNED,
        }),
      })}
      getColumnProps={column => {
        if (column.id === 'actions') {
          return { style: { textAlign: 'center' } }
        }

        return {}
      }}
      primaryAttribute="id"
      onRowsSelect={setSelectedDevices}
      bottomControls={
        <Button
          onClick={onDeleteClick}
          variant="secondary"
          icon="fa-trash-alt"
          disabled={loading}
        >
          {_(t.flushCache)}
        </Button>
      }
      unselectRowsToken={unselectRowsToken}
    />
  )
}

DevicesList.propTypes = {
  data: PropTypes.arrayOf(deviceShape),
  selectedDevices: PropTypes.arrayOf(PropTypes.string).isRequired,
  setSelectedDevices: PropTypes.func.isRequired,
  loading: PropTypes.bool.isRequired,
  onDeleteClick: PropTypes.func.isRequired,
  unselectRowsToken: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
}

DevicesList.defaultProps = {
  data: [],
  unselectRowsToken: null,
}
