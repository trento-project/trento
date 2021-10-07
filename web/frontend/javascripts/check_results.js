import React, { useEffect, useState } from 'react';
import ReactDOM from 'react-dom';
import { get } from 'axios';

import ChecksTable from '@components/ChecksTable';

const clusterId = window.location.pathname.split('/').pop();

const ClustersChecks = ({ clusterId }) => {
  const [results, setResults] = useState([]);
  const [hosts, setHosts] = useState({});

  useEffect(() => {
    get(`/api/clusters/${clusterId}/results`).then(({ data }) => {
      setResults(data.checks);
      setHosts(data.hosts);
    });
  }, []);

  return (
    <div>
      <ChecksTable checks={results} clusterHosts={hosts} />
    </div>
  );
};

ReactDOM.render(
  <ClustersChecks clusterId={clusterId} />,
  document.getElementById('cluster-checks-results')
);
