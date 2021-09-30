import React, { useState } from 'react';
import classNames from 'classnames';

const Accordion = ({ className, title, children }) => {
  const [isOpen, setOpen] = useState(true);
  const collapseClassNames = classNames('collapse', { show: isOpen });
  const iconClassNames = classNames(
    'eos-icons',
    'eos-18',
    'collapse-toggle',
    'clickable',
    { collapsed: !isOpen }
  );

  return (
    <div className={`${className} card`}>
      <div className="card-header">
        <h4 className="float-left">{title}</h4>
        <i className={iconClassNames} onClick={() => setOpen(!isOpen)}></i>
      </div>
      <div className={collapseClassNames}>
        <div className="card-body">{children}</div>
      </div>
    </div>
  );
};

export default Accordion;
