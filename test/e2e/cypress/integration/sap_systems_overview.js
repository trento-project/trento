import {
  availableSAPSystems,
  isHanaPrimary,
  isHanaSecondary,
} from '../fixtures/sap-systems-overview/available_sap_systems';

context('SAP Systems Overview', () => {
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');

    cy.visit('/');
    cy.navigateToItem('SAP Systems');
    cy.url().should('include', '/sapsystems');
  });

  describe('Registered SAP Systems should be available in the overview', () => {
    describe('Discovered SID are the expected ones', () => {
      availableSAPSystems.forEach(({ sid: sid }) => {
        it(`should have a sid named ${sid}`, () => {
          cy.get('td').contains(sid);
        });
      });
    });

    describe('System healths are the expected ones', () => {
      availableSAPSystems.forEach(({ sid: sid, health: health }, index) => {
        it(`should have a health ${health} for sid ${sid}`, () => {
          cy.get('.eos-table')
            .eq(0)
            .find('tr')
            .filter(':visible')
            .eq(index + 1)
            .find('td')
            .as('tableCell');
          cy.get('@tableCell').eq(0).should('contain', health);
        });
      });
    });

    describe('Links to the details page are the expected ones', () => {
      availableSAPSystems.forEach(({ sid: sid, id: id }) => {
        it(`should have a link to the SAP System with id: ${id}`, () => {
          cy.get('td').contains(sid).click();
          cy.location('pathname').should('eq', `/sapsystems/${id}`);
          cy.go('back');
        });
      });
    });

    describe('Attached databases are the expected ones', () => {
      availableSAPSystems.forEach(
        ({ sid: sid, attachedDatabase: attachedDatabase }) => {
          it(`should show the expected attached database details`, () => {
            cy.get('td')
              .contains(sid)
              .parent('td')
              .parent('tr')
              .within(() => {
                cy.get('td').eq(4).contains(attachedDatabase.sid);
                cy.get('td').eq(5).contains(attachedDatabase.tenant);
                cy.get('td').eq(6).contains(attachedDatabase.dbAddress);
              });
          });
          it(`should have a link to the attached HANA database with id: ${attachedDatabase.id}`, () => {
            cy.contains(attachedDatabase.sid).click();
            cy.location('pathname').should(
              'eq',
              `/databases/${attachedDatabase.id}`
            );
            cy.go('back');
          });
        }
      );
    });

    describe('Instances are the expected ones', () => {
      availableSAPSystems.forEach(({ id: id, instances: instances }) => {
        it(`should show the expected instances details`, () => {
          cy.get('.collapse-toggle').each(($el) => {
            cy.wrap($el).click({ force: true });
          });

          cy.get(`#instances-${id}`)
            .find('tr')
            .each((row, index) => {
              cy.wrap(row).within(() => {
                cy.get('td').eq(0).should('contain', instances[index].health);
                cy.get('td').eq(1).should('contain', instances[index].sid);
                cy.get('td').eq(2).should('contain', instances[index].features);
                cy.get('td')
                  .eq(3)
                  .should('contain', instances[index].instanceNumber);
                cy.get('td')
                  .eq(4)
                  .should('contain', instances[index].systemReplication);
                if (isHanaPrimary(instances[index])) {
                  cy.get('td')
                    .eq(4)
                    .should(
                      'not.contain',
                      instances[index].systemReplicationStatus
                    );
                }
                if (isHanaSecondary(instances[index])) {
                  cy.get('td')
                    .eq(4)
                    .should(
                      'contain',
                      instances[index].systemReplicationStatus
                    );
                }
                cy.get('td')
                  .eq(5)
                  .should('contain', instances[index].clusterName);
                cy.get('td').eq(6).should('contain', instances[index].hostname);
              });
            });
        });

        it(`should have a link to known type clusters`, () => {
          cy.on('uncaught:exception', () => {
            // do not fail on XHR requests in the cluster page
            return false;
          });

          cy.get(`#instances-${id}`)
            .find('tr')
            .each((row, index) => {
              if (
                instances[index].clusterName !== '' &&
                instances[index].clusterID !== ''
              ) {
                cy.wrap(row)
                  .get('td')
                  .contains(instances[index].clusterName)
                  .click({ force: true });
                cy.location('pathname').should(
                  'eq',
                  `/clusters/${instances[index].clusterID}`
                );
                cy.go('back');
              }
            });
        });

        it(`should have a link to the hosts`, () => {
          cy.get(`#instances-${id}`)
            .find('tr')
            .each((row, index) => {
              cy.wrap(row)
                .get('td')
                .contains(instances[index].hostname)
                .click({ force: true });
              cy.location('pathname').should(
                'eq',
                `/hosts/${instances[index].hostID}`
              );
              cy.go('back');
            });
        });
      });
    });
  });

  describe('SAP Systems Tagging', () => {
    before(() => {
      cy.get('body').then(($body) => {
        const deleteTag = '.tn-sap-systems-tags x';
        if ($body.find(deleteTag).length > 0) {
          cy.get(deleteTag).then(($deleteTag) =>
            cy.wrap($deleteTag).click({ multiple: true })
          );
        }
      });
    });

    availableSAPSystems.forEach(({ sid, tag }) => {
      describe(`Add tag '${tag}' to SAP System with sid: '${sid}'`, () => {
        it(`should tag SAP System '${sid}'`, () => {
          cy.get('td')
            .contains(sid)
            .parent('td')
            .parent('tr')
            .within(() => {
              cy.get('.tagify').type(`${tag}{enter}`);
            });
        });
      });
    });
  });

  describe('Filtering the SAP Systems overview', () => {
    describe('Filtering by SIDs', () => {
      before(() => {
        cy.get('.filter-option').eq(0).click();
      });
      availableSAPSystems.forEach(({ sid }) => {
        it(`should have SAP Systems ${sid}'`, () => {
          cy.intercept('GET', `/sapsystems?sids=${sid}`).as('filterBySIDs');
          cy.get('.dropdown-item').contains(sid).click();

          cy.wait('@filterBySIDs').then(() => {
            cy.get('td').should('contain', sid);
            availableSAPSystems
              .filter(({ sid: s }) => s !== sid)
              .forEach(({ sid: s }) => {
                cy.get('td').should('not.contain', s);
              });
            cy.intercept('GET', `/sapsystems`).as('resetFilter');
            cy.get('.dropdown-item').contains(sid).click();
            cy.wait('@resetFilter');
          });
        });
      });
    });

    describe('Filtering by tags', () => {
      before(() => {
        cy.get('.filter-option').eq(1).click();
      });
      availableSAPSystems.forEach(({ sid, tag }) => {
        it(`should have SAP Systems ${sid} tagged with tag '${tag}'`, () => {
          cy.intercept('GET', `/sapsystems?tags=${tag}`).as('filterByTags');
          cy.get('.dropdown-item').contains(tag).click();

          cy.wait('@filterByTags').then(() => {
            cy.get('td').should('contain', tag);
            availableSAPSystems
              .filter(({ tag: t }) => t !== tag)
              .forEach(({ tag: t }) => {
                cy.get('td').should('not.contain', t);
              });
            cy.intercept('GET', `/sapsystems`).as('resetFilter');
            cy.get('.dropdown-item').contains(tag).click();
            cy.wait('@resetFilter');
          });
        });
      });
    });

    describe('Health states are updated', () => {
      const states = [
        ['GRAY', 'fiber_manual_record'],
        ['YELLOW', 'warning'],
        ['RED', 'error'],
      ];

      states.forEach(([state, health], index) => {
        it(`should have ${state} health in SAP system and instance ${
          index + 1
        } when SAPControl-${state} state is received`, () => {
          cy.loadScenario(`sap-systems-overview-${state}`);
          cy.visit(`/sapsystems`);

          cy.get('.eos-table')
            .eq(0)
            .find('tr')
            .filter(':visible')
            .eq(1)
            .find('td')
            .as('tableCell');
          cy.get('@tableCell').eq(0).should('contain', health);

          cy.get('.eos-table')
            .eq(0)
            .find('tr')
            .eq(index + 4) // + 4 moves selects the row within the collpased table
            .find('td')
            .as('instanceTableCell');

          cy.get('@instanceTableCell').eq(0).should('contain', health);
        });
      });

      it(`should have RED health in SAP system and HANA instance 1 when SAPControl-RED state is received`, () => {
        cy.loadScenario('healthy-27-node-SAP-cluster');
        cy.loadScenario(`sap-systems-overview-hana-RED`);
        cy.visit(`/sapsystems`);

        cy.get('.eos-table')
          .eq(0)
          .find('tr')
          .filter(':visible')
          .eq(1)
          .find('td')
          .as('tableCell');
        cy.get('@tableCell').eq(0).should('contain', 'error');

        cy.get('.eos-table')
          .eq(0)
          .find('tr')
          .eq(8) // + 4 moves selects the row within the collpased table
          .find('td')
          .as('instanceTableCell');

        cy.get('@instanceTableCell').eq(0).should('contain', 'error');
      });
    });
    describe('SAP diagnostics agent', () => {
      it(`should skip SAP diagnostics agent discovery visualization`, () => {
        cy.loadScenario('sap-systems-overview-DAA');
        cy.visit(`/sapsystems`);
        cy.get('.eos-table').should('not.contain', 'DAA');
      });
    });
  });
});
