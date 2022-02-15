import React, { Fragment, useState, useEffect, useCallback } from 'react';
import ReactDOM from 'react-dom';
import { get, post } from 'axios';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Accordion from 'react-bootstrap/Accordion';
import Card from 'react-bootstrap/Card';
import Table from 'react-bootstrap/Table';
import Form from 'react-bootstrap/Form';
import Spinner from 'react-bootstrap/Spinner';

import { logError } from '@lib/log';
import { toggle, hasOne, remove } from '@lib/lists';
import Checkbox from '@components/Checkbox';
import { AccordionToggle } from '@components/Accordion';
import { showSuccessToast, showErrorToast } from '@components/Toast';

const clusterId = window.location.pathname.split('/').pop();

const getChecksIds = (checks) => checks.map(({ id }) => id);

const mergeConnectionSettings = (hostnames, connectionSettings) =>
  hostnames.reduce(
    (accumulator, current) =>
      connectionSettings[current]
        ? { ...accumulator, [current]: connectionSettings[current] }
        : { ...accumulator, [current]: '' },
    {}
  );

const SettingsButton = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [checksCatalog, setChecksCatalog] = useState([]);
  const [selectedChecks, setSelectedChecks] = useState([]);
  const [settings, setSettings] = useState({});
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    get('/api/checks/catalog')
      .then(({ data }) => {
        data ? setChecksCatalog(data) : setChecksCatalog([]);
      })
      .catch((error) => {
        logError(error);
        setChecksCatalog([]);
        showErrorToast({
          content: 'Error fetching checks catalog data.',
        });
      });

    get(`/api/checks/${clusterId}/settings`)
      .then(({ data }) => {
        const {
          hostnames,
          connection_settings: connectionSettings,
          selected_checks: selectedChecks,
        } = data;
        const newSettings = mergeConnectionSettings(
          hostnames,
          connectionSettings
        );
        setSettings(newSettings);
        setSelectedChecks(selectedChecks);
        setLoading(false);
      })
      .catch((error) => {
        logError(error);
        setSelectedChecks([]);
        setSettings({});
        setLoading(false);
        showErrorToast({
          content: 'Error fetching the checks data, please refresh.',
        });
      });
  }, [modalOpen]);

  const submit = useCallback(() => {
    const payload = {
      selected_checks: selectedChecks,
      connection_settings: settings,
    };
    setLoading(true);
    post(`/api/checks/${clusterId}/settings`, payload)
      .then(() => {
        setLoading(false);
        setModalOpen(false);
        showSuccessToast({
          content: 'Cluster settings successfully saved.',
        });
      })
      .catch((err) => {
        logError(err);
        setLoading(false);
        showErrorToast({
          content: 'Error saving the checks settings, please retry',
        });
      });
  }, [selectedChecks, settings]);

  return (
    <Fragment>
      <Button variant="secondary" size="sm" onClick={() => setModalOpen(true)}>
        <i className="eos-icons eos-18">settings</i>Settings
      </Button>
      <Modal show={modalOpen} onHide={() => setModalOpen(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Cluster settings</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <h6>Connection settings</h6>
          <Accordion>
            <Card>
              <Card.Header>
                Host connection settings
                <AccordionToggle
                  className="float-right"
                  eventKey="connection-settings"
                />
              </Card.Header>
              <Accordion.Collapse eventKey="connection-settings">
                <Card.Body className="card-check-selection">
                  <Table>
                    <thead>
                      <tr>
                        <th>Host</th>
                        <th>Connection user</th>
                        <th>Default user</th>
                      </tr>
                    </thead>
                    <tbody>
                      {Object.keys(settings).map((host) => (
                        <tr key={host}>
                          <td>{host}</td>
                          <td>
                            <Form.Control
                              size="sm"
                              value={settings[host]}
                              onChange={({ target: { value } }) =>
                                setSettings({ ...settings, [host]: value })
                              }
                            />
                          </td>
                          <td>root</td>
                        </tr>
                      ))}
                    </tbody>
                  </Table>
                </Card.Body>
              </Accordion.Collapse>
            </Card>
          </Accordion>
          <h6>Checks selection</h6>
          <Accordion>
            {checksCatalog.map(({ group, checks }) => (
              <Card key={group}>
                <Card.Header>
                  <Checkbox
                    label={group}
                    inline
                    checked={hasOne(getChecksIds(checks), selectedChecks)}
                    onChange={() => {
                      const checksIds = getChecksIds(checks);
                      const newSelectionSet = hasOne(checksIds, selectedChecks)
                        ? remove(checksIds, selectedChecks)
                        : [...selectedChecks, ...checksIds];
                      setSelectedChecks(newSelectionSet);
                    }}
                  />
                  <AccordionToggle className="float-right" eventKey={group} />
                </Card.Header>

                <Accordion.Collapse eventKey={group}>
                  <Card.Body className="card-check-selection">
                    <Table>
                      <thead>
                        <tr>
                          <th className="header-check-selection"></th>
                          <th className="header-check-selection">Test ID</th>
                          <th className="header-check-selection">
                            Description
                          </th>
                        </tr>
                      </thead>
                      <tbody>
                        {checks.map(({ id, description }) => (
                          <tr key={id}>
                            <td className="row-status">
                              <Checkbox
                                label={id}
                                checked={selectedChecks.includes(id)}
                                onChange={() => {
                                  setSelectedChecks(toggle(id, selectedChecks));
                                }}
                              />
                            </td>
                            <td>{id}</td>
                            <td>{description}</td>
                          </tr>
                        ))}
                      </tbody>
                    </Table>
                  </Card.Body>
                </Accordion.Collapse>
              </Card>
            ))}
          </Accordion>
        </Modal.Body>
        <Modal.Footer>
          <Button
            variant="secondary"
            disabled={loading}
            onClick={() => setModalOpen(false)}
          >
            Close
          </Button>
          <Button variant="primary" disabled={loading} onClick={submit}>
            {loading && (
              <Spinner animation="border" role="status" as="span" size="sm" />
            )}{' '}
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>
    </Fragment>
  );
};

ReactDOM.render(
  <SettingsButton clusterId={clusterId} />,
  document.getElementById('cluster-settings-button')
);
