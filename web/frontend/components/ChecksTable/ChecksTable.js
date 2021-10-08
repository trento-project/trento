import React from 'react';
import CheckResultIcon from './CheckResultIcon';

const ChecksTable = ({ checks, clusterHosts }) => {
  return (
    <div className="table-responsive">
      <table className="table eos-table">
        <thead>
          <tr>
            <th>Test ID</th>
            <th>Description</th>
            {Object.keys(clusterHosts).map((label) => (
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
                {Object.keys(clusterHosts).map((hostname) => (
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
