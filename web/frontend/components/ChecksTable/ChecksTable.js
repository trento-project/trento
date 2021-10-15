import React from 'react';

import RowGroup from './RowGroup';
import CheckResultIcon from './CheckResultIcon';

const groupChecks = (checks) => {
  const groups = checks.reduce((accumulator, current) => {
    const { group } = current;
    return accumulator[group]
      ? { ...accumulator, [group]: [...accumulator[group], current] }
      : { ...accumulator, [group]: [current] };
  }, {});

  return Object.keys(groups).map((key) => {
    return { name: key, checks: groups[key] };
  });
};

const ChecksTable = ({ checks, clusterHosts }) => {
  const checkGroups = groupChecks(checks);

  return (
    <div className="table-responsive">
      <table className="table eos-table">
        <thead>
          <tr>
            <th>Description</th>
            <th>Test ID</th>
            {Object.keys(clusterHosts).map((label) => (
              <th key={label} scope="col" style={{ textAlign: 'center' }}>
                {!(clusterHosts[label].reachable) &&
                  <CheckResultIcon result="warning" tooltip={clusterHosts[label].msg}/>}{label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {checkGroups.map(({ name, checks }) => (
            <RowGroup
              key={name}
              id={name}
              name={name}
              checks={checks}
              clusterHosts={clusterHosts}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ChecksTable;
