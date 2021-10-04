import React from 'react';
import CheckResultIcon from './CheckResultIcon';

const ChecksTable = ({ checks }) => {
  const hostnames = Object.keys(checks[0].hosts);

  return (
    <div className="table-responsive">
      <table className="table eos-table">
        <thead>
          <tr>
            <th>Test ID</th>
            <th>Description</th>
            {hostnames.map((label) => (
              <th key={label} scope="col" style={{ textAlign: 'center' }}>
                {label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {checks.map(({ id, description, hosts }) => {
            return (
              <tr key={id}>
                <td>{id}</td>
                <td>{description}</td>
                {Object.keys(hosts).map((hostname) => (
                  <td key={hostname} className="align-center">
                    <CheckResultIcon result={hosts[hostname].result} />
                  </td>
                ))}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

export default ChecksTable;
