import React, { Fragment, useState } from 'react';
import CheckResultIcon from './CheckResultIcon';

const RowGroup = ({ name, checks }) => {
  const [open, setOpen] = useState(true);
  const emptyCells = Object.keys(checks[0].hosts)
    .map((key) => <td key={key} />)
    .concat(<td key="emptycell" />);

  return (
    <Fragment>
      <tr className="checks-table-row-group" onClick={() => setOpen(!open)}>
        <td className="checks-table-row-group-label">{name}</td>
        {emptyCells}
      </tr>
      {open &&
        checks.map(({ id, description, hosts }) => {
          return (
            <tr key={id}>
              <td>{description}</td>
              <td>{id}</td>
              {Object.keys(hosts).map((hostname) => (
                <td key={hostname} className="align-center">
                  <CheckResultIcon result={hosts[hostname].result} />
                </td>
              ))}
            </tr>
          );
        })}
    </Fragment>
  );
};

export default RowGroup;
