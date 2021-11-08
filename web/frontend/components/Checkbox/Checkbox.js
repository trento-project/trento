import React from 'react';
import Form from 'react-bootstrap/Form';

const noop = () => {};

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
