import PropTypes from 'prop-types'

export const Display = ({ when, children }) => (when ? children : null)

Display.propTypes = {
  children: PropTypes.node.isRequired,
  when: PropTypes.bool.isRequired,
}
