context('SAP Systems Overview', () => {
  const availableSAPSystems = [
    {
      sid: 'NWD',
      id: 'a1e80e3e152a903662f7882fb3f8a851',
      attachedDatabase: {
        sid: 'HDD',
        id: 'fd44c254ccb14331e54015c720c7a1f2',
        tenant: 'HDD',
        dbAddress: '10.100.1.13',
      },
      instances: [
        {
          sid: 'NWD',
          features: 'MESSAGESERVER|ENQUE',
          instanceNumber: '00',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: 'netweaver_cluster',
          clusterID: '',
          hostname: 'vmnwdev01',
          hostID: '7269ee51-5007-5849-aaa7-7c4a98b0c9ce',
        },
        {
          sid: 'NWD',
          features: 'ABAP|GATEWAY|ICMAN|IGS',
          instanceNumber: '01',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: '',
          clusterID: '',
          hostname: 'vmnwdev03',
          hostID: '9a3ec76a-dd4f-5013-9cf0-5eb4cf89898f',
        },
        {
          sid: 'NWD',
          features: 'ABAP|GATEWAY|ICMAN|IGS',
          instanceNumber: '02',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: '',
          clusterID: '',
          hostname: 'vmnwdev04',
          hostID: '1b0e9297-97dd-55d6-9874-8efde4d84c90',
        },
        {
          sid: 'NWD',
          features: 'ENQREP',
          instanceNumber: '10',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: 'netweaver_cluster',
          clusterID: '',
          hostname: 'vmnwdev02',
          hostID: 'fb2c6b8a-9915-5969-a6b7-8b5a42de1971',
        },
        {
          sid: 'HDD',
          features: 'HDB|HDB_WORKER',
          instanceNumber: '10',
          systemReplication: 'HANA Secondary',
          systemReplicationStatus: '',
          clusterName: 'hana_cluster',
          clusterID: '04b8f8c21f9fd8991224478e8c4362f8',
          hostname: 'vmhdbdev02',
          hostID: '0a055c90-4cb6-54ce-ac9c-ae3fedaf40d4',
        },
        {
          sid: 'HDD',
          features: 'HDB|HDB_WORKER',
          instanceNumber: '10',
          systemReplication: 'HANA Primary',
          systemReplicationStatus: 'SOK',
          clusterName: 'hana_cluster',
          clusterID: '04b8f8c21f9fd8991224478e8c4362f8',
          hostname: 'vmhdbdev01',
          hostID: '13e8c25c-3180-5a9a-95c8-51ec38e50cfc',
        },
      ],
      tag: 'env3',
    },
    {
      id: '97a1e70aeff3c0685d65c4c3d32d533b',
      sid: 'NWP',
      attachedDatabase: {
        sid: 'HDP',
        id: '1154f7678ac587e5f0f242830a5201f1',
        tenant: 'HDP',
        dbAddress: '10.80.1.13',
      },
      instances: [
        {
          sid: 'NWP',
          features: 'MESSAGESERVER|ENQUE',
          instanceNumber: '00',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: 'netweaver_cluster',
          clusterID: '',
          hostname: 'vmnwprd01',
          hostID: '116d49bd-85e1-5e59-b820-83f66db8800c',
        },
        {
          sid: 'NWP',
          features: 'ABAP|GATEWAY|ICMAN|IGS',
          instanceNumber: '01',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: '',
          clusterID: '',
          hostname: 'vmnwprd03',
          hostID: 'a3297d85-5e8b-5ac5-b8a3-55eebc2b8d12',
        },
        {
          sid: 'NWP',
          features: 'ABAP|GATEWAY|ICMAN|IGS',
          instanceNumber: '02',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: '',
          clusterID: '',
          hostname: 'vmnwprd04',
          hostID: '0fc07435-7ee2-54ca-b0de-fb27ffdc5deb',
        },
        {
          sid: 'NWP',
          features: 'ENQREP',
          instanceNumber: '10',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: 'netweaver_cluster',
          clusterID: '',
          hostname: 'vmnwprd02',
          hostID: '4b30a6af-4b52-5bda-bccb-f2248a12c992',
        },
        {
          sid: 'HDP',
          features: 'HDB|HDB_WORKER',
          instanceNumber: '10',
          systemReplication: 'HANA Primary',
          systemReplicationStatus: 'SOK',
          clusterName: 'hana_cluster',
          clusterID: '04b8f8c21f9fd8991224478e8c4362f8',
          hostname: 'vmhdbprd01',
          hostID: '9cd46919-5f19-59aa-993e-cf3736c71053',
        },
        {
          sid: 'HDP',
          features: 'HDB|HDB_WORKER',
          instanceNumber: '10',
          systemReplication: 'HANA ',
          systemReplicationStatus: '',
          clusterName: 'hana_cluster',
          clusterID: '04b8f8c21f9fd8991224478e8c4362f8',
          hostname: 'vmhdbprd02',
          hostID: 'b767b3e9-e802-587e-a442-541d093b86b9',
        },
      ],
      tag: 'env1',
    },
    {
      id: 'd01fdc69aeba7bd5133b210eb2884853',
      sid: 'NWQ',
      attachedDatabase: {
        sid: 'HDQ',
        id: '9953878f07bb54cac20d5d5d7ff08af2',
        tenant: 'HDQ',
        dbAddress: '10.90.1.13',
      },
      instances: [
        {
          sid: 'NWQ',
          features: 'MESSAGESERVER|ENQUE',
          instanceNumber: '00',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: 'netweaver_cluster',
          clusterID: '',
          hostname: 'vmnwqas01',
          hostID: '25677e37-fd33-5005-896c-9275b1284534',
        },
        {
          sid: 'NWQ',
          features: 'ABAP|GATEWAY|ICMAN|IGS',
          instanceNumber: '01',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: '',
          clusterID: '',
          hostname: 'vmnwqas03',
          hostID: '098fc159-3ed6-58e7-91be-38fda8a833ea',
        },
        {
          sid: 'NWQ',
          features: 'ABAP|GATEWAY|ICMAN|IGS',
          instanceNumber: '02',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: '',
          clusterID: '',
          hostname: 'vmnwqas04',
          hostID: '81e9b629-c1e7-538f-bff1-47d3a6580522',
        },
        {
          sid: 'NWQ',
          features: 'ENQREP',
          instanceNumber: '10',
          systemReplication: '',
          systemReplicationStatus: '',
          clusterName: 'netweaver_cluster',
          clusterID: '',
          hostname: 'vmnwqas02',
          hostID: '3711ea88-9ccc-5b07-8f9d-042be449d72b',
        },
        {
          sid: 'HDQ',
          features: 'HDB|HDB_WORKER',
          instanceNumber: '10',
          systemReplication: 'HANA Primary',
          systemReplicationStatus: 'SOK',
          clusterName: 'hana_cluster',
          clusterID: '04b8f8c21f9fd8991224478e8c4362f8',
          hostname: 'vmhdbqas01',
          hostID: '99cf8a3a-48d6-57a4-b302-6e4482227ab6',
        },
        {
          sid: 'HDQ',
          features: 'HDB|HDB_WORKER',
          instanceNumber: '10',
          systemReplication: 'HANA Secondary',
          systemReplicationStatus: '',
          clusterName: 'hana_cluster',
          clusterID: '04b8f8c21f9fd8991224478e8c4362f8',
          hostname: 'vmhdbqas02',
          hostID: 'e0c182db-32ff-55c6-a9eb-2b82dd21bc8b',
        },
      ],
      tag: 'env2',
    },
  ];

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
          cy.get('.collapse-toggle').each(($el, index) => {
            cy.wrap($el).click({ force: true });
          });

          cy.get(`#instances-${id}`)
            .find('tr')
            .each((row, index) => {
              cy.wrap(row).within(() => {
                cy.get('td').eq(1).should('contain', instances[index].sid);
                cy.get('td').eq(2).should('contain', instances[index].features);
                cy.get('td')
                  .eq(3)
                  .should('contain', instances[index].instanceNumber);
                cy.get('td')
                  .eq(4)
                  .should('contain', instances[index].systemReplication);
                cy.get('td')
                  .eq(4)
                  .should('contain', instances[index].systemReplicationStatus);
                cy.get('td')
                  .eq(5)
                  .should('contain', instances[index].clusterName);
                cy.get('td').eq(6).should('contain', instances[index].hostname);
              });
            });
        });

        it(`should have a link to known type clusters`, () => {
          cy.on('uncaught:exception', (err, runnable) => {
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
  });
});
