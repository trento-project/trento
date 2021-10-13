import React, { useState, useCallback, cloneElement } from 'react';
import classNames from 'classnames';

const Dropdown = ({
  children,
  label = 'Click here',
  className,
  multi = false,
}) => {
  const [isOpen, setOpen] = useState(false);
  const dropdownClasses = classNames('dropdown', className, { show: isOpen });
  const menuClasses = classNames('dropdown-menu', { show: isOpen });

  const closeDropdown = useCallback(() => setOpen(false), [setOpen]);

  return (
    <div className={dropdownClasses}>
      <button
        className="btn btn-secondary dropdown-toggle"
        data-toggle="dropdown"
      >
        <i
          className="eos-icons eos-18 icon-reset"
          onClick={() => setOpen(!isOpen)}
        >
          keyboard_arrow_down
        </i>
        {label}
      </button>
      <div className={menuClasses} aria-labelledby="dropdownMenuButton">
        {React.Children.map(children, (childNode) =>
          React.isValidElement(childNode)
            ? cloneElement(childNode, { closeDropdown, multi })
            : childNode
        )}
      </div>
    </div>
  );
};

export default Dropdown;
