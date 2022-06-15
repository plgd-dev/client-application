import Form from 'react-bootstrap/Form'

export const TextField = ({ value, inputRef, ...rest }) => (
  <Form.Control {...rest} ref={inputRef} type="text" value={value} />
)
