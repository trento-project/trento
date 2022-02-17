import React, { useState, useEffect } from 'react';
import ReactDOM from 'react-dom';
import { get } from 'axios';

import { logError } from '@lib/log';
import HealthSummary from '@components/HealthSummary';

const HomePage = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    get('/api/sapsystems/health')
      .then(({ data }) => {
        setLoading(false);
        data ? setData(data) : setData([]);
      })
      .catch((error) => {
        setLoading(false);
        logError(error);
        setData([]);
      });
  }, []);

  return loading ? (
    <div>Loading...</div>
  ) : (
    <div>
      <HealthSummary data={data} />
    </div>
  );
};

ReactDOM.render(<HomePage />, document.getElementById('homepage-component'));
