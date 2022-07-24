import PropTypes from 'prop-types'
import BDropdown from 'react-bootstrap/Dropdown'
import omit from 'lodash/omit'

import { dropdownTypes } from './constants'

const { PRIMARY, SECONDARY, EMPTY } = dropdownTypes

export const ActionButton = ({ type, menuProps, items, onToggle, ...rest }) => {
  const getIcon = item => {
    if (item.loading) {
      return <i className={`fas fa-spinner m-r-10`} />
    } else if (item.icon) {
      return <i className={`fas ${item.icon} m-r-10`} />
    }

    return null
  }
  return (
    <BDropdown className="action-button" onToggle={onToggle}>
      <BDropdown.Toggle variant={type} {...omit(rest, 'children')}>
        <span />
        <span />
        <span />
      </BDropdown.Toggle>

      <BDropdown.Menu
        renderOnMount={true}
        {...menuProps}
        popperConfig={{
          strategy: 'fixed',
          modifiers: [
            {
              name: 'offset',
              options: {
                offset: [-9, -15],
              },
            },
          ],
        }}
      >
        {items
          .filter(item => !item.hidden)
          .map(item => {
            return (
              item.component || (
                <BDropdown.Item
                  className="btn btn-secondary"
                  key={item.id || item.label}
                  onClick={item.onClick}
                  disabled={item.loading}
                >
                  {getIcon(item)}
                  {!item.loading && item.label}
                </BDropdown.Item>
              )
            )
          })}
      </BDropdown.Menu>
    </BDropdown>
  )
}

ActionButton.propTypes = {
  children: PropTypes.node.isRequired,
  type: PropTypes.oneOf([PRIMARY, SECONDARY, EMPTY]),
  items: PropTypes.arrayOf(
    PropTypes.shape({
      onClick: PropTypes.func,
      label: PropTypes.string,
      id: PropTypes.string,
      hidden: PropTypes.bool,
      component: PropTypes.node,
      loading: PropTypes.bool,
    })
  ).isRequired,
  menuProps: PropTypes.shape({
    align: PropTypes.string,
    flip: PropTypes.bool,
  }),
}

ActionButton.defaultProps = {
  type: EMPTY,
  menuProps: {},
}
