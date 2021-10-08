import React, { useEffect, useState } from 'react';
import ReactDOM from 'react-dom';
import { get } from 'axios';

import Accordion from '@components/Accordion';
import ChecksTable from '@components/ChecksTable';

const clusterId = window.location.pathname.split('/').pop();

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

const ClustersChecks = ({ clusterId }) => {
  const [results, setResults] = useState([]);
  const [hosts, setHosts] = useState({});

  useEffect(() => {
    get(`/api/clusters/${clusterId}/results`).then(({ data }) => {
      const groupedChecks = groupChecks(data.checks);
      setResults(groupedChecks);
      setHosts(data.hosts);
    });
  }, []);

  return (
    <div>
      {results.map((section) => {
        return (
          <Accordion
            className="checks-results-accordion"
            key={section.name}
            title={section.name}
          >
            <ChecksTable checks={section.checks} clusterHosts={hosts} />
          </Accordion>
        );
      })}
    </div>
  );
};

ReactDOM.render(
  <ClustersChecks clusterId={clusterId} />,
  document.getElementById('cluster-checks-results')
);
