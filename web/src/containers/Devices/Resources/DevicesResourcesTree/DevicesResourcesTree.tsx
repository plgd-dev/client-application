import { FC, useMemo } from 'react'
import { useIntl } from 'react-intl'

import TreeExpander from '@shared-ui/components/new/TreeExpander'
import TreeTable from '@shared-ui/components/new/Table'
import Badge from '@shared-ui/components/new/Badge'
import DevicesResourcesActionButton from '../DevicesResourcesActionButton'
import { RESOURCE_TREE_DEPTH_SIZE } from '../../constants'
import {
  canBeResourceEdited,
  createNestedResourceData,
  getLastPartOfAResourceHref,
} from '../../utils'
import { messages as t } from '../../Devices.i18n'
import { Props } from './DevicesResourcesTree.types'

const DevicesResourcesTree: FC<Props> = ({
  data: rawData,
  onUpdate,
  onCreate,
  onDelete,
  loading,
  isOwned,
  deviceId,
}) => {
  const { formatMessage: _ } = useIntl()
  const data = useMemo(() => createNestedResourceData(rawData), [rawData])

  const columns = useMemo(
    () => [
      {
        Header: _(t.href),
        accessor: 'href',
        Cell: ({ value, row }: { value: any; row: any }) => {
          const {
            original: { href, endpointInformations },
          } = row

          const lastValue = getLastPartOfAResourceHref(value)
          const onLinkClick = deviceId
            ? () => onUpdate({ deviceId, href: href.replace(/\/$/, '') })
            : null

          const edit = canBeResourceEdited(endpointInformations) || isOwned

          if (!edit) {
            return <span>{lastValue}</span>
          }

          if (row.canExpand) {
            return (
              <div className="tree-expander-container">
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
                {!row.canExpand && <i className="fas fa-pen" />}
              </div>
            )
          }

          return (
            <div
              className="tree-expander-container"
              style={{
                marginLeft: `${
                  row.depth === 0
                    ? 0
                    : (row.depth + 1) * RESOURCE_TREE_DEPTH_SIZE
                }px`,
              }}
            >
              <span
                className="link reveal-icon-on-hover"
                onClick={() => onLinkClick}
              >
                {`/${lastValue}`}
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
        Cell: ({ value }: { value: any }) => {
          if (!deviceId) {
            return null
          }

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
        Cell: ({ row }: { row: any }) => {
          if (row.canExpand) {
            return null
          }

          const {
            original: { href, interfaces, endpointInformations },
          } = row
          const cleanHref = href.replace(/\/$/, '') // href without a trailing slash
          return (
            <DevicesResourcesActionButton
              disabled={loading}
              href={cleanHref}
              deviceId={deviceId}
              interfaces={interfaces}
              onCreate={onCreate}
              onUpdate={onUpdate}
              onDelete={onDelete}
              isOwned={isOwned}
              endpointInformations={endpointInformations || []}
            />
          )
        },
      },
    ],
    [onUpdate, onCreate, onDelete, loading] //eslint-disable-line
  )

  return <TreeTable columns={columns} data={data || []} />
}

DevicesResourcesTree.displayName = 'DevicesResourcesTree'

export default DevicesResourcesTree
