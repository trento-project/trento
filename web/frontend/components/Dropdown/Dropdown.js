import React, { cloneElement } from 'react';
import classNames from 'classnames';

const Dropdown = ({
  children,
  label = 'Click here',
  className,
  multi = false,
}) => {
  const dropdownClasses = classNames('dropdown', className);
  const menuClasses = classNames('dropdown-menu');

  return (
    <div className={dropdownClasses}>
      <button
        className="btn btn-secondary dropdown-toggle"
        data-toggle="dropdown"
      >
        <i className="eos-icons eos-18 icon-reset">keyboard_arrow_down</i>
        {label}
      </button>
      <div className={menuClasses} aria-labelledby="dropdownMenuButton">
        {React.Children.map(children, (childNode) =>
          React.isValidElement(childNode)
            ? cloneElement(childNode, { multi })
            : childNode
        )}
      </div>
    </div>
  );
};

export default Dropdown;
