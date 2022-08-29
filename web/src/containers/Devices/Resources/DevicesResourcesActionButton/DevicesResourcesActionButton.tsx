import { useIntl } from 'react-intl'

import ActionButton from '@shared-ui/components/new/ActionButton'
import { canCreateResource, canBeResourceEdited } from '../../utils'
import { messages as t } from '../../Devices.i18n'
import { FC } from 'react'
import { defaultProps, Props } from './DevicesResourcesActionButton.types'

export const DevicesResourcesActionButton: FC<Props> = ({
  href,
  deviceId,
  interfaces,
  onCreate,
  onUpdate,
  onDelete,
  isOwned,
  endpointInformations,
}) => {
  const { formatMessage: _ } = useIntl()

  const create = interfaces && canCreateResource(interfaces) && isOwned
  const edit =
    (endpointInformations && canBeResourceEdited(endpointInformations)) ||
    isOwned

  if (!create && !edit) {
    return null
  }

  return (
    <ActionButton
      menuProps={{
        align: 'start',
      }}
      items={[
        {
          onClick: () => onCreate(href),
          label: _(t.create),
          icon: 'fa-plus',
          hidden: !create,
        },
        {
          onClick: () => onUpdate({ deviceId, href }),
          label: _(t.update),
          icon: 'fa-pen',
          hidden: !edit,
        },
        {
          onClick: () => onDelete(href),
          label: _(t.delete),
          icon: 'fa-trash-alt',
          hidden: !edit,
        },
      ]}
    >
      <i className="fas fa-ellipsis-h" />
    </ActionButton>
  )
}

DevicesResourcesActionButton.displayName = 'DevicesResourcesActionButton'
DevicesResourcesActionButton.defaultProps = defaultProps

export default DevicesResourcesActionButton
