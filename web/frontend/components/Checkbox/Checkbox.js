import React from 'react';
import Form from 'react-bootstrap/Form';

const noop = () => {};

// Custom Checkbox component. Why do we need this? Because our design system doesn't play
// super-well with React. That said, our CSS styles the :before of the checkbox' label,
// and places it _over_ the actual checkbox input. The result? The input doesn't trigger any
// onChange event. 12 hours of my life wasted debugging this.
const Checkbox = ({
  className = '',
  checked,
  onChange,
  label = '',
  inline = false,
  ...props
}) => {
  return (
    <Form.Check inline={inline}>
      <Form.Check.Input
        className={`eos-checkbox ${className}`}
        checked={checked}
        onChange={noop}
        {...props}
      />
      <Form.Check.Label onClick={onChange}>{label}</Form.Check.Label>
    </Form.Check>
  );
};

export default Checkbox;
