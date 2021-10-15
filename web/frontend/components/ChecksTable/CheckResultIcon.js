import React from 'react';

const PASSING = 'passing';
const WARNING = 'warning';
const CRITICAL = 'critical';
const SKIPPED = 'skipped';

const CheckResultIcon = ({ result, tooltip = null }) => {
  const tooltipData = tooltip ? { datatoggle: 'tooltip', title: tooltip } : {};

  switch (result) {
    case PASSING:
      return (
        <i className="eos-icons eos-18 text-success" {...tooltipData}>
          check_circle
        </i>
      );
    case WARNING:
      return (
        <i className="eos-icons eos-18 text-warning" {...tooltipData}>
          warning
        </i>
      );
    case CRITICAL:
      return (
        <i className="eos-icons eos-18 text-danger" {...tooltipData}>
          error
        </i>
      );
    case SKIPPED:
      return (
        <i className="eos-icons eos-18 text-muted" {...tooltipData}>
          fiber_manual_record
        </i>
      );
    default:
      return (
        <i className="eos-icons eos-18 text-muted" {...tooltipData}>
          fiber_manual_record
        </i>
      );
  }
};

export default CheckResultIcon;
