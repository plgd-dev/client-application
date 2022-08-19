import { useIntl } from 'react-intl'
import PropTypes from 'prop-types'

import ActionButton from '@shared-ui/components/new/ActionButton'
import { canCreateResource, canBeResourceEdited } from './utils'
import { messages as t } from './devices-i18n'

export const DevicesResourcesActionButton = ({
  disabled,
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

  const create = canCreateResource(interfaces) && isOwned
  const edit = canBeResourceEdited(endpointInformations) || isOwned

  if (!create && !edit) {
    return null
  }

  return (
    <ActionButton
      disabled={disabled}
      menuProps={{
        align: 'right',
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

DevicesResourcesActionButton.propTypes = {
  disabled: PropTypes.bool.isRequired,
  href: PropTypes.string.isRequired,
  deviceId: PropTypes.string.isRequired,
  onCreate: PropTypes.func.isRequired,
  onUpdate: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  interfaces: PropTypes.arrayOf(PropTypes.string),
}

DevicesResourcesActionButton.defaultProps = {
  interfaces: [],
}
