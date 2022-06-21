import { useMemo } from 'react'
import PropTypes from 'prop-types'
import { useIntl } from 'react-intl'
import classNames from 'classnames'
import { Badge } from '@/components/badge'
import { Table } from '@/components/table'
import { DevicesResourcesActionButton } from './_devices-resources-action-button'
import { RESOURCES_DEFAULT_PAGE_SIZE, devicesStatuses } from './constants'
import { deviceResourceShape } from './shapes'
import { messages as t } from './devices-i18n'
import { canBeResourceEdited } from '@/containers/devices/utils'

export const DevicesResourcesList = ({
  data,
  onUpdate,
  onCreate,
  onDelete,
  deviceStatus,
  deviceId,
  loading,
  isOwned,
}) => {
  const { formatMessage: _ } = useIntl()
  const isUnregistered = deviceStatus === devicesStatuses.UNREGISTERED
  const greyedOutClassName = classNames({ 'grayed-out': isUnregistered })

  const columns = useMemo(
    () => [
      {
        Header: _(t.href),
        accessor: 'href',
        Cell: ({ value, row }) => {
          const {
            original: { deviceId: deviceIdOrigin, href, endpointInformations },
          } = row

          const edit = canBeResourceEdited(endpointInformations)

          if (!edit) {
            return <span>{value}</span>
          }
          return (
            <div className="tree-expander-container">
              <span
                className="link reveal-icon-on-hover"
                onClick={() => onUpdate({ deviceId: deviceIdOrigin, href })}
              >
                {value}
              </span>
              <i className="fas fa-pen" />
            </div>
          )
        },
        style: { width: '100%' },
      },
      {
        Header: _(t.types),
        accessor: 'resourceTypes',
        Cell: ({ value }) => {
          return (
            <div className="badges-box-horizontal">
              {value?.map?.(type => (
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
        Cell: ({ row }) => {
          const {
            original: { href, interfaces, endpointInformations },
          } = row
          return (
            <DevicesResourcesActionButton
              disabled={isUnregistered || loading}
              href={href}
              deviceId={deviceId}
              interfaces={interfaces}
              onCreate={onCreate}
              onUpdate={onUpdate}
              onDelete={onDelete}
              isOwned={isOwned}
              endpointInformations={endpointInformations}
            />
          )
        },
        className: 'actions',
      },
    ],
    [onUpdate, onCreate, onDelete, isUnregistered, loading] //eslint-disable-line
  )

  return (
    <Table
      columns={columns}
      data={data || []}
      defaultSortBy={[
        {
          id: 'href',
          desc: false,
        },
      ]}
      defaultPageSize={RESOURCES_DEFAULT_PAGE_SIZE}
      autoFillEmptyRows
      className={greyedOutClassName}
      paginationProps={{
        className: greyedOutClassName,
        disabled: isUnregistered,
      }}
    />
  )
}

DevicesResourcesList.propTypes = {
  data: PropTypes.arrayOf(deviceResourceShape),
  onCreate: PropTypes.func.isRequired,
  onUpdate: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  loading: PropTypes.bool.isRequired,
  deviceStatus: PropTypes.oneOf(Object.values(devicesStatuses)),
}

DevicesResourcesList.defaultProps = {
  data: null,
  deviceStatus: null,
}
