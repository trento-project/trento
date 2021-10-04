import React from 'react';

const PASSING = 'passing';
const WARNING = 'warning';
const CRITICAL = 'critical';
const SKIPPED = 'skipped';

const CheckResultIcon = ({ result }) => {
  switch (result) {
    case PASSING:
      return <i className="eos-icons eos-18 text-success">check_circle</i>;
    case WARNING:
      return <i className="eos-icons eos-18 text-warning">warning</i>;
    case CRITICAL:
      return <i className="eos-icons eos-18 text-danger">error</i>;
    case SKIPPED:
      return <i className="eos-icons eos-18 text-muted">fiber_manual_record</i>;
    default:
      return <i className="eos-icons eos-18 text-muted">fiber_manual_record</i>;
  }
};

export default CheckResultIcon;
