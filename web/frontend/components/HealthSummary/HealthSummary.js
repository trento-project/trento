import React from 'react';

const getHealthIcon = (health) => {
  switch (health) {
    case 'critical':
      return <i className="eos-icons health-summary-icon text-danger">error</i>;
    case 'warning':
      return (
        <i className="eos-icons health-summary-icon text-warning">warning</i>
      );
    case 'passing':
      return (
        <i className="eos-icons health-summary-icon text-success">
          check_circle
        </i>
      );
    default:
      return (
        <i className="eos-icons health-summary-icon text-muted">
          fiber_manual_record
        </i>
      );
  }
};

const any = (predicate, label) =>
  Object.keys(predicate).reduce((accumulator, key) => {
    if (accumulator) {
      return true;
    }
    return predicate[key] === label;
  }, false);

const getCounters = (data) =>
  data.reduce(
    (accumulator, element) => {
      if (any(element, 'critical')) {
        return { ...accumulator, critical: accumulator.critical + 1 };
      }

      if (any(element, 'warning')) {
        return { ...accumulator, warning: accumulator.warning + 1 };
      }

      if (any(element, 'unknown')) {
        return { ...accumulator, unknown: accumulator.unknown + 1 };
      }

      if (any(element, 'passing')) {
        return { ...accumulator, passing: accumulator.passing + 1 };
      }
      return accumulator;
    },
    { critical: 0, warning: 0, passing: 0, unknown: 0 }
  );

const HealthSummary = ({ data }) => {
  const counters = getCounters(data);

  return (
    <div>
      <div className="col">
        <div className="row">
          <div className="col">
            <h1>At a glance</h1>
          </div>
        </div>
        <hr className="margin-10px" />

        <h5>Global Health</h5>
        <div className="health-container horizontal-container">
          <div className="alert alert-inline alert-success health-passing">
            <i className="eos-icons-outlined eos-18 alert-icon">check_circle</i>
            <div className="alert-body">Passing</div>
            <span className="badge badge-secondary">{counters.passing}</span>
          </div>
          <div className="alert alert-inline alert-warning health-warning">
            <i className="eos-icons-outlined eos-18 alert-icon">warning</i>
            <div className="alert-body">Warning</div>
            <span className="badge badge-secondary">{counters.warning}</span>
          </div>
          <div className="alert alert-inline alert-danger health-critical">
            <i className="eos-icons-outlined eos-18 alert-icon">error</i>
            <div className="alert-body">Critical</div>
            <span className="badge badge-secondary">{counters.critical}</span>
          </div>
        </div>

        <div className="health-summary">
          <div className="health-summary-header">
            <span className="health-summary-id health-summary-cell">SID</span>
            <span className="health-summary-cell">SAP Instances</span>
            <span className="health-summary-cell">Database</span>
            <span className="health-summary-cell">Pacemaker Clusters</span>
            <span className="health-summary-cell">Hosts</span>
          </div>
          {data.map(
            ({
              id,
              sid,
              clusters_health,
              database_health,
              hosts_health,
              sapsystem_health,
            }) => {
              return (
                <div className="health-summary-row" key={id}>
                  <span className="health-summary-id health-summary-cell">
                    {sid}
                  </span>
                  <span className="health-summary-cell">
                    {getHealthIcon(sapsystem_health)}
                  </span>
                  <span className="health-summary-cell">
                    {getHealthIcon(database_health)}
                  </span>
                  <span className="health-summary-cell">
                    {getHealthIcon(clusters_health)}
                  </span>
                  <span className="health-summary-cell">
                    {getHealthIcon(hosts_health)}
                  </span>
                </div>
              );
            }
          )}
        </div>
      </div>
    </div>
  );
};

export default HealthSummary;
