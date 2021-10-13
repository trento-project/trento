import React from 'react';
import classNames from 'classnames';

const DropdownItem = ({
  children,
  className,
  onClick,
  closeDropdown,
  multi = false,
}) => (
  <a
    className={classNames('dropdown-item', className)}
    href="#"
    onClick={(e) => {
      e.preventDefault();
      e.stopPropagation();
      if (closeDropdown && !multi) {
        closeDropdown();
      }
      onClick(e);
    }}
  >
    {children}
  </a>
);

export default DropdownItem;
