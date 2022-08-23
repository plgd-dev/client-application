import { FC, useMemo } from 'react'
import { useIntl } from 'react-intl'
import Badge from '@shared-ui/components/new/Badge'
import Table from '@shared-ui/components/new/Table'
import DevicesResourcesActionButton from '../DevicesResourcesActionButton'
import { RESOURCES_DEFAULT_PAGE_SIZE } from '../../constants'
import { messages as t } from '../../Devices.i18n'
import { canBeResourceEdited } from '@/containers/devices/utils'
import { Props } from './DevicesResourcesList.types'

const DevicesResourcesList: FC<Props> = ({
  data,
  onUpdate,
  onCreate,
  onDelete,
  deviceId,
  loading,
  isOwned,
}) => {
  const { formatMessage: _ } = useIntl()

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
        style: { width: '70%' },
      },
      {
        Header: _(t.types),
        accessor: 'resourceTypes',
        style: { width: '20%' },
        Cell: ({ value }: { value: any }) => {
          return (
            <div className="badges-box-horizontal">
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
        style: { textAlign: 'center' },
        Cell: ({ row }: { row: any }) => {
          const {
            original: { href, interfaces, endpointInformations },
          } = row
          return (
            <DevicesResourcesActionButton
              disabled={loading}
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
    [onUpdate, onCreate, onDelete, loading] //eslint-disable-line
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
    />
  )
}

DevicesResourcesList.displayName = 'DevicesResourcesList'

export default DevicesResourcesList
