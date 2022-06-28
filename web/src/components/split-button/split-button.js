import PropTypes from 'prop-types'
import RBDropdown from 'react-bootstrap/Dropdown'
import classNames from 'classnames'

import { buttonVariants, iconPositions } from '../button/constants'
import { Button } from '@/components/button'

const { PRIMARY, SECONDARY } = buttonVariants
const { ICON_LEFT, ICON_RIGHT } = iconPositions

export const SplitButton = ({
  children,
  variant,
  className,
  menuProps,
  items,
  disabled,
  ...rest
}) => (
  <RBDropdown className="split-button">
    <Button
      {...rest}
      variant={variant}
      disabled={disabled}
      className={classNames('split-button-left', className)}
    >
      {children}
    </Button>
    <RBDropdown.Toggle
      variant={variant}
      disabled={disabled}
      className="split-button-right"
    />

    <RBDropdown.Menu
      {...menuProps}
      popperConfig={{
        strategy: 'fixed',
        modifiers: [
          {
            name: 'offset',
            options: {
              offset: [0, 5],
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
              <RBDropdown.Item
                className="btn btn-secondary"
                key={item.id || item.label}
                onClick={item.onClick}
              >
                {item.icon && <i className={`fas ${item.icon} m-r-10`} />}
                {item.label}
              </RBDropdown.Item>
            )
          )
        })}
    </RBDropdown.Menu>
  </RBDropdown>
)

SplitButton.propTypes = {
  variant: PropTypes.oneOf([PRIMARY, SECONDARY]),
  icon: PropTypes.string,
  iconPosition: PropTypes.oneOf([ICON_LEFT, ICON_RIGHT]),
  onClick: PropTypes.func,
  loading: PropTypes.bool,
  className: PropTypes.string,
  items: PropTypes.arrayOf(
    PropTypes.shape({
      onClick: PropTypes.func,
      label: PropTypes.string,
      id: PropTypes.string,
      hidden: PropTypes.bool,
      component: PropTypes.node,
    })
  ).isRequired,
  menuProps: PropTypes.shape({
    align: PropTypes.string,
    flip: PropTypes.bool,
  }),
}

SplitButton.defaultProps = {
  variant: SECONDARY,
  icon: null,
  iconPosition: ICON_LEFT,
  onClick: null,
  loading: false,
  className: null,
  menuProps: {},
}
