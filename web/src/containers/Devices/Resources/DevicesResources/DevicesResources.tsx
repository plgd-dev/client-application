import { useIntl } from 'react-intl'
import classNames from 'classnames'

import Switch from '@shared-ui/components/new/Switch'
import { useLocalStorage } from '@shared-ui/common/hooks'
import DevicesResourcesList from '../DevicesResourcesList'
import DevicesResourcesTree from '../DevicesResourcesTree'
import { devicesStatuses } from '../../constants'
import { messages as t } from '../../Devices.i18n'
import { FC } from 'react'
import { Props } from './DevicesResources.types'

export const DevicesResources: FC<Props> = ({
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
  const [treeViewActive, setTreeViewActive] = useLocalStorage(
    'treeViewActive',
    false
  )
  const isUnregistered = devicesStatuses.UNREGISTERED === deviceStatus
  const greyedOutClassName = classNames({
    'grayed-out': isUnregistered,
  })

  return (
    <>
      <div
        className={classNames(
          'd-flex justify-content-between align-items-center',
          greyedOutClassName
        )}
      >
        <h2>{_(t.resources)}</h2>
        <div className="d-flex justify-content-end align-items-center">
          <Switch
            id="toggle-tree-view"
            label={_(t.treeView)}
            checked={treeViewActive}
            onChange={() => setTreeViewActive(!treeViewActive)}
            disabled={isUnregistered}
          />
        </div>
      </div>

      {treeViewActive ? (
        <DevicesResourcesTree
          data={data}
          onUpdate={onUpdate}
          onCreate={onCreate}
          onDelete={onDelete}
          loading={loading}
          deviceId={deviceId}
          isOwned={isOwned}
        />
      ) : (
        <DevicesResourcesList
          data={data}
          onUpdate={onUpdate}
          onCreate={onCreate}
          onDelete={onDelete}
          loading={loading}
          deviceId={deviceId}
          isOwned={isOwned}
        />
      )}
    </>
  )
}

DevicesResources.displayName = 'DevicesResources'

export default DevicesResources
