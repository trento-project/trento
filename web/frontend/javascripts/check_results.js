import React, { useEffect, useState } from 'react';
import ReactDOM from 'react-dom';
import { get } from 'axios';
import classNames from 'classnames';

import ChecksTable, { CheckResultIcon } from '@components/ChecksTable';
import Dropdown, { DropdownItem } from '@components/Dropdown';

const clusterId = window.location.pathname.split('/').pop();

const toggleFilter = (filter, selectedFilters) =>
  selectedFilters.includes(filter)
    ? selectedFilters.filter((string) => string !== filter)
    : [...selectedFilters, filter];

const ClustersChecks = ({ clusterId }) => {
  const [results, setResults] = useState([]);
  const [hosts, setHosts] = useState({});
  const [filters, setFilters] = useState([]);

  const displayedResults = results.filter(({ hosts }) => {
    if (filters.length === 0) {
      return true;
    }
    return Object.keys(hosts).reduce((accumulator, hostname) => {
      return accumulator ? true : filters.includes(hosts[hostname].result);
    }, false);
  });

  useEffect(() => {
    get(`/api/checks/${clusterId}/results`).then(({ data }) => {
      setResults(data.checks);
      setHosts(data.hosts);
    });
  }, []);

  return (
    <div>
      <div className="checks-filters-row">
        <Dropdown multi label="Filter the checks list">
          <DropdownItem
            className={classNames({ selected: filters.includes('warning') })}
            onClick={() => setFilters(toggleFilter('warning', filters))}
          >
            <CheckResultIcon result="warning" />
            Warning
          </DropdownItem>
          <DropdownItem
            className={classNames({ selected: filters.includes('critical') })}
            onClick={() => setFilters(toggleFilter('critical', filters))}
          >
            <CheckResultIcon result="critical" />
            Critical
          </DropdownItem>
          <DropdownItem onClick={() => setFilters([])}>See all</DropdownItem>
        </Dropdown>
      </div>
      <ChecksTable checks={displayedResults} clusterHosts={hosts} />
    </div>
  );
};

ReactDOM.render(
  <ClustersChecks clusterId={clusterId} />,
  document.getElementById('cluster-checks-results')
);
