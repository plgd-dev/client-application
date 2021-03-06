import { useMemo } from 'react'
import PropTypes from 'prop-types'
import { useIntl } from 'react-intl'

import { TreeExpander } from '@/components/tree-expander'
import { TreeTable } from '@/components/table'
import { Badge } from '@/components/badge'
import { DevicesResourcesActionButton } from './_devices-resources-action-button'
import { RESOURCE_TREE_DEPTH_SIZE } from './constants'
import {
  canBeResourceEdited,
  createNestedResourceData,
  getLastPartOfAResourceHref,
} from './utils'
import { deviceResourceShape } from './shapes'
import { messages as t } from './devices-i18n'

export const DevicesResourcesTree = ({
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
        Cell: ({ value, row }) => {
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
                  expanded={!!row.isExpanded}
                  style={{
                    marginLeft: `${row.depth * RESOURCE_TREE_DEPTH_SIZE}px`,
                  }}
                />
                <span
                  className={!row.canExpand ? 'link reveal-icon-on-hover' : ''}
                  onClick={onLinkClick}
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
              <span className="link reveal-icon-on-hover" onClick={onLinkClick}>
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
        Cell: ({ value }) => {
          if (!deviceId) {
            return null
          }

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

DevicesResourcesTree.propTypes = {
  data: PropTypes.arrayOf(deviceResourceShape),
  onCreate: PropTypes.func.isRequired,
  onUpdate: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  loading: PropTypes.bool.isRequired,
}

DevicesResourcesTree.defaultProps = {
  data: null,
}
