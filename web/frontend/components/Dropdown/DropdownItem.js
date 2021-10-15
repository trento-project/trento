import React from 'react';
import classNames from 'classnames';

const DropdownItem = ({ children, className, onClick, multi = false }) => (
  <a
    className={classNames('dropdown-item', className)}
    href="#"
    onClick={(e) => {
      if (multi) {
        e.preventDefault();
        e.stopPropagation();
      }
      onClick(e);
    }}
  >
    {children}
  </a>
);

export default DropdownItem;
