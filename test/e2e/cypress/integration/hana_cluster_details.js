import { availableHanaCluster } from '../fixtures/hana-cluster-details/available_hana_cluster';

context('HANA database details', () => {
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');
    cy.loadChecksCatalog('checks-catalog/catalog.json');
    cy.loadChecksResults(
      'hana-cluster-details/checks_results.json',
      '9c832998801e28cd70ad77380e82a5c0'
    );

    cy.visit(`/clusters/${availableHanaCluster.id}`);
    cy.url().should('include', `/clusters/${availableHanaCluster.id}`);
  });

  describe('HANA cluster details should be consistent with the state of the cluster', () => {
    it(`should have name ${availableHanaCluster.name}`, () => {
      cy.get('span').contains(availableHanaCluster.name);
    });

    it(`should have sid ${availableHanaCluster.sid}`, () => {
      cy.get('span').contains(availableHanaCluster.sid);
    });

    it(`should have cluster type ${availableHanaCluster.clusterType}`, () => {
      cy.get('span').contains(availableHanaCluster.clusterType);
    });

    it(`should have system replication mode ${availableHanaCluster.hanaSystemReplicationMode}`, () => {
      cy.get('span').contains(availableHanaCluster.hanaSystemReplicationMode);
    });

    it(`should have fencing type ${availableHanaCluster.fencingType}`, () => {
      cy.get('span').contains(availableHanaCluster.fencingType);
    });

    it(`should have sapHanaSRHealthState ${availableHanaCluster.hanaSecondarySyncState}`, () => {
      cy.get('span').contains(availableHanaCluster.hanaSecondarySyncState);
    });

    it(`should have sap hana sr health state ${availableHanaCluster.sapHanaSRHealthState}`, () => {
      cy.get('span').contains(availableHanaCluster.sapHanaSRHealthState);
    });

    it(`should have CIB last written ${availableHanaCluster.cibLastWritten}`, () => {
      cy.get('span').contains(availableHanaCluster.cibLastWritten);
    });

    it(`should have hana system replication operation mode ${availableHanaCluster.hanaSystemReplicationOperationMode}`, () => {
      cy.get('span').contains(
        availableHanaCluster.hanaSystemReplicationOperationMode
      );
    });
  });

  describe('Cluster health count should show the right checks status', () => {
    it('should show the expected passing checks count', () => {
      cy.get('.health-passing').contains(2);
    });

    it('should show the expected warning checks count', () => {
      cy.get('.health-warning').contains(2);
    });

    it('should show the expected critical checks count', () => {
      cy.get('.health-critical').contains(2);
    });
  });

  describe('Check results modal should show the expected checks', () => {
    before(() => {
      cy.get('button').contains('Show check results').click();
    });

    const checkResults = [
      { id: 'C620DC', icon: 'check_circle' },
      { id: '00081D', icon: 'warning', filter: 'Warning' },
      { id: '0B6DB2', icon: 'error', filter: 'Critical' },
    ];

    it('should show the expected checks', () => {
      cy.get('.modal-title').contains('Health details');
      checkResults.forEach(({ id, icon }) => {
        cy.get('.checks-table-row-group')
          .parent()
          .within(() => {
            cy.contains(id);
            cy.get(`i:contains("${icon}")`).should('have.length', 2);
          });
      });
    });

    it('should filter by checks status', () => {
      cy.get('.dropdown-toggle').click();
      checkResults.forEach(({ icon, filter }) => {
        if (filter) {
          cy.get('.dropdown-item').contains(filter).click();
          cy.get('.checks-table-row-group')
            .parent()
            .within(() => {
              cy.get(`i:contains("${icon}")`).should('have.length', 2);
              // Check that the other checks status are not present
              checkResults.forEach(({ icon: otherIcon }) => {
                if (icon !== otherIcon) {
                  cy.get(`i:contains("${otherIcon}")`).should('have.length', 0);
                }
              });
            });
          cy.get('.dropdown-item').contains(filter).click();
        }
      });
    });

    after(() => {
      cy.get('.modal-title')
        .contains('Health details')
        .parent()
        .within(() => {
          cy.get('.close').click();
        });
    });
  });

  describe('Cluster sites should have the expected hosts', () => {
    availableHanaCluster.sites.forEach((site) => {
      it(`should have ${site.name}`, () => {
        cy.get('td').contains(site.name);
      });

      site.hosts.forEach((host) => {
        it(`site ${site.name} should have host ${host.hostname}`, () => {
          cy.get('td').contains(host.hostname);
        });

        it(`${host.hostname} should have the expected IP addresses`, () => {
          host.ips.forEach((ip) => {
            cy.get('td').contains(ip);
          });
        });

        it(`${host.hostname} should have the expected virtual IP addresses`, () => {
          host.virtualIps.forEach((ip) => {
            cy.get('td').contains(ip);
          });
        });

        it(`${host.hostname} should have the expected role`, () => {
          cy.get('td').contains(host.role);
        });

        it(`${host.hostname} should have the expected attributes and resources`, () => {
          cy.get('tbody').within(() => {
            host.attributes.forEach(({ attribute, value }) => {
              cy.get('td').contains(attribute);
              cy.get('td').contains(value);
            });

            host.resources.forEach(({ id, type, role, status, failCount }) => {
              cy.get('td').contains(id);
              cy.get('td').contains(type);
              cy.get('td').contains(role);
              cy.get('td').contains(status);
              cy.get('td').contains(failCount);
            });
          });
        });
      });
    });
  });

  describe('Cluster SBD should have the expected devices with the correct status', () => {
    availableHanaCluster.sbd.forEach((item, index) => {
      it(`should have SBD device name "${item.deviceName}" and status "${item.status}"`, () => {
        cy.get('.eos-table')
          .eq(2)
          .find('tr')
          .eq(index + 1)
          .find('td')
          .as('tableCell');
        cy.get('@tableCell').eq(0).should('contain', item.status);
        cy.get('@tableCell').eq(1).should('contain', item.deviceName);
      });
    });
  });
});
