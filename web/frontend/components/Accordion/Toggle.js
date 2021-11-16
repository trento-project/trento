import React from 'react';
import { useAccordionToggle } from 'react-bootstrap/AccordionToggle';

const AccordionToggle = ({ className = '', eventKey, callback = () => {} }) => {
  const decoratedOnClick = useAccordionToggle(
    eventKey,
    () => callback && callback(eventKey)
  );

  return (
    <i
      className={`eos-icons eos-18 trento-pointer ${className}`}
      onClick={decoratedOnClick}
    >
      keyboard_arrow_down
    </i>
  );
};

export default AccordionToggle;
