import {
  allClusterNames,
  allClusterIds,
  clusterIdByName,
} from '../fixtures/clusters-overview/available_clusters';

context('Clusters Overview', () => {
  const availableClusters = allClusterNames();
  const availableClustersId = allClusterIds();
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');
    cy.loadChecksCatalog('checks-catalog/catalog.json');
    cy.visit('/');
    cy.navigateToItem('Clusters');
    cy.url().should('include', '/clusters');
  });

  describe('Registered Clusters should be available in the overview', () => {
    it('should show all of the registered clusters with default pagination settings', () => {
      cy.get('.tn-clustername')
        .its('length')
        .should('eq', availableClusters.length);
    });
    it('should show 9 as total items in the pagination controls', () => {
      cy.get('.pagination-count').should('contain', '9 items');
    });
    it('should have 1 pages', () => {
      cy.get('.page-item').its('length').should('eq', 3); // We add +2 to the page count because of the first and last page
    });
    describe('Discovered clusternames are the expected ones', () => {
      availableClusters.forEach((clusterName) => {
        it(`should have a cluster named ${clusterName}`, () => {
          cy.get('.tn-clustername').each(($link) => {
            const displayedClusterName = $link.text().trim();
            expect(availableClusters).to.include(displayedClusterName);
          });
        });
      });
    });
  });

  describe('Health State', () => {
    before(() => {
      cy.pruneChecksResults();
    });

    const healthStates = [
      ['check_circle', '4', 'SOK'],
      ['error', '1', 'SOK'],
      ['error', '4', 'SFAIL'],
      ['error', '1', 'SFAIL'],
    ];
    healthStates.forEach(([health, srState, syncState]) => {
      it(`should show a ${health} state when SR health state is ${srState} and sync state is ${syncState}`, () => {
        cy.loadScenario(`cluster-${srState}-${syncState}`);
        cy.visit(`/clusters`);

        cy.get(`#cluster-${clusterIdByName('hana_cluster_1')}`)
          .find('td')
          .as('tableCell');
        cy.get('@tableCell').eq(0).should('contain', health);
      });
    });

    const healthStatesWithChecks = [
      ['check_circle', '4', 'SOK', 'passing'],
      ['error', '1', 'SFAIL', 'passing'],
      ['warning', '4', 'SOK', 'warning'],
      ['error', '4', 'SOK', 'critical'],
    ];

    healthStatesWithChecks.forEach(
      ([health, srState, syncState, checksResult]) => {
        it(`should show a ${health} state when SR health state is ${srState} and sync state is ${syncState} and checks results are ${checksResult}`, () => {
          cy.loadScenario(`cluster-${srState}-${syncState}`);
          cy.loadChecksResults(
            `clusters-overview/checks_results_${checksResult}.json`,
            '04b8f8c21f9fd8991224478e8c4362f8'
          );

          cy.visit(`/clusters`);

          cy.get(`#cluster-${clusterIdByName('hana_cluster_1')}`)
            .find('td')
            .as('tableCell');
          cy.get('@tableCell').eq(0).should('contain', health);
        });
      }
    );

    it('should show an unknown state when the cluster is not of HANA type', () => {
      cy.visit(`/clusters`);

      cy.get(`#cluster-${clusterIdByName('drbd_cluster')}`)
        .find('td')
        .as('tableCell');
      cy.get('@tableCell').eq(0).should('contain', 'fiber_manual_record');
    });
  });

  describe('Health Container', () => {
    before(() => {
      cy.pruneChecksResults();
      cy.loadChecksResults(
        'clusters-overview/checks_results_critical.json',
        '04b8f8c21f9fd8991224478e8c4362f8'
      );
      cy.loadChecksResults(
        'clusters-overview/checks_results_warning.json',
        '4e905d706da85f5be14f85fa947c1e39'
      );
      cy.loadChecksResults(
        'clusters-overview/checks_results_passing.json',
        '9c832998801e28cd70ad77380e82a5c0'
      );
    });

    describe('Health Container shows the health overview of all Clusters', () => {
      it('should show health status of the entire cluster when set to paginate by 10', () => {
        cy.reloadList('clusters', 10);
        cy.get('.health-container .health-passing').should('contain', 1);
        cy.get('.health-container .health-warning').should('contain', 1);
        cy.get('.health-container .health-critical').should('contain', 1);
      });
      it('should show health status of the entire cluster when set to paginate by 100', () => {
        cy.reloadList('clusters', 100);
        cy.get('.health-container .health-passing').should('contain', 1);
        cy.get('.health-container .health-warning').should('contain', 1);
        cy.get('.health-container .health-critical').should('contain', 1);
      });
    });
  });

  describe('Clusters Tagging', () => {
    before(() => {
      cy.get('body').then(($body) => {
        const deleteTag = '.tn-cluster-tags x';
        if ($body.find(deleteTag).length > 0) {
          cy.get(deleteTag).then(($deleteTag) =>
            cy.wrap($deleteTag).click({ multiple: true })
          );
        }
      });
    });
    const clustersByMatchingPattern = (pattern) => (clusterName) =>
      clusterName.includes(pattern);
    const taggingRules = [
      ['hana_cluster_1', 'env1'],
      ['hana_cluster_2', 'env2'],
      ['hana_cluster_3', 'env3'],
    ];

    taggingRules.forEach(([pattern, tag]) => {
      describe(`Add tag '${tag}' to all clusters with '${pattern}' in the cluster name`, () => {
        availableClusters
          .filter(clustersByMatchingPattern(pattern))
          .forEach((clusterName) => {
            it(`should tag cluster '${clusterName}'`, () => {
              cy.get(
                `#cluster-${clusterIdByName(
                  clusterName
                )} > .tn-cluster-tags > .tagify`
              )
                .type(tag)
                .trigger('change');
            });
          });
      });
    });
  });

  describe('Filtering the Clusters overview', () => {
    before(() => {
      cy.reloadList('clusters', 100);
    });

    const resetFilter = (option) => {
      cy.intercept('GET', `/clusters?per_page=100`).as('resetFilter');
      cy.get(option).click();
      cy.wait('@resetFilter');
    };

    describe('Filtering by health', () => {
      before(() => {
        cy.get('.tn-filters > :nth-child(2) > .btn').click();
      });
      const healthScenarios = [
        ['passing', 1],
        ['warning', 1],
        ['critical', 1],
      ];
      healthScenarios.forEach(
        ([health, expectedClustersWithThisHealth], index) => {
          it(`should show ${expectedClustersWithThisHealth} clusters when filtering by health '${health}'`, () => {
            cy.intercept('GET', `/clusters?per_page=100&health=${health}`).as(
              'filterByHealthStatus'
            );
            const selectedOption = `#bs-select-1-${index}`;
            cy.get(selectedOption).click();
            cy.wait('@filterByHealthStatus').then(() => {
              cy.get('.tn-clustername')
                .its('length')
                .should('eq', expectedClustersWithThisHealth);
              cy.get('.pagination-count').should(
                'contain',
                `${expectedClustersWithThisHealth} items`
              );
              cy.get('.page-item')
                .its('length')
                .should(
                  'eq',
                  Math.ceil(expectedClustersWithThisHealth / 100) + 2
                );
              resetFilter(selectedOption);
            });
          });
        }
      );
    });

    describe('Filtering by SAP system', () => {
      before(() => {
        cy.get('.tn-filters > :nth-child(4) > .btn').click();
      });
      const SAPSystemsScenarios = [
        ['HDD', 1],
        ['HDP', 1],
        ['HDQ', 1],
      ];
      SAPSystemsScenarios.forEach(
        ([sapsystem, expectedRelatedClusters], index) => {
          it(`should have ${expectedRelatedClusters} clusters related to SAP system '${sapsystem}'`, () => {
            cy.intercept('GET', `/clusters?per_page=100&sids=${sapsystem}`).as(
              'filterBySAPSystem'
            );
            const selectedOption = `#bs-select-3-${index}`;
            cy.get(selectedOption).click();
            cy.wait('@filterBySAPSystem').then(() => {
              cy.get('.tn-clustername')
                .its('length')
                .should('eq', expectedRelatedClusters);
              cy.get('.pagination-count').should(
                'contain',
                `${expectedRelatedClusters} items`
              );
              cy.get('.page-item')
                .its('length')
                .should('eq', Math.ceil(expectedRelatedClusters / 100) + 2);
            });
            resetFilter(selectedOption);
          });
        }
      );
    });

    describe('Filtering by Cluster name', () => {
      before(() => {
        cy.get('.tn-filters > :nth-child(3) > .btn').click();
      });
      const clusterNameScenarios = [
        ['drbd_cluster', 3],
        ['hana_cluster_1', 1],
        ['hana_cluster_2', 1],
        ['hana_cluster_3', 1],
        ['netweaver_cluster', 3],
      ];
      clusterNameScenarios.forEach(
        ([clusterName, expectedRelatedClusters], index) => {
          it(`should have ${expectedRelatedClusters} clusters related to name '${clusterName}'`, () => {
            cy.intercept(
              'GET',
              `/clusters?per_page=100&name=${clusterName}`
            ).as('filterByClusterName');
            const selectedOption = `#bs-select-2-${index}`;
            cy.get(selectedOption).click();
            cy.wait('@filterByClusterName').then(() => {
              cy.get('.tn-clustername')
                .its('length')
                .should('eq', expectedRelatedClusters);
              cy.get('.pagination-count').should(
                'contain',
                `${expectedRelatedClusters} items`
              );
              cy.get('.page-item')
                .its('length')
                .should('eq', Math.ceil(expectedRelatedClusters / 100) + 2);
            });
            resetFilter(selectedOption);
          });
        }
      );
    });

    describe('Filtering by tags', () => {
      before(() => {
        cy.get('.tn-filters > :nth-child(6) > .btn').click();
      });
      const tagsScenarios = [
        ['env1', 1],
        ['env2', 1],
        ['env3', 1],
      ];
      tagsScenarios.forEach(([tag, expectedTaggedClusters], index) => {
        it(`should have ${expectedTaggedClusters} clusters tagged with tag '${tag}'`, () => {
          cy.intercept('GET', `/clusters?per_page=100&tags=${tag}`).as(
            'filterByTags'
          );
          const selectedOption = `#bs-select-5-${index}`;
          cy.get(selectedOption).click();
          cy.wait('@filterByTags').then(() => {
            cy.get('.tn-clustername')
              .its('length')
              .should('eq', expectedTaggedClusters);
            cy.get('.pagination-count').should(
              'contain',
              `${expectedTaggedClusters} items`
            );
            cy.get('.page-item')
              .its('length')
              .should('eq', Math.ceil(expectedTaggedClusters / 100) + 2);
            resetFilter(selectedOption);
          });
        });
      });
    });

    describe('Removing filtered tags', () => {
      const tag = 'tag1';

      before(() => {
        // Wait for the POST that gets triggered on tag submission
        cy.intercept('/api/clusters/**').as('tagPosted');
        // Wait for the filter reset after the first POST finishes
        cy.intercept('GET', '/api/tags?resource_type=clusters').as(
          'filterRefreshed'
        );

        // Clear the filter dropdown by clicking outside
        cy.get('h1').click();

        for (let i = 0; i < 2; i++) {
          cy.get(
            `#cluster-${availableClustersId[i]} > .tn-cluster-tags > .tagify`
          )
            .type(tag)
            .click();

          // We wait for the POST to finish and the filter to be refreshed
          cy.wait('@tagPosted');
          cy.wait('@filterRefreshed');
        }

        cy.get('.dropdown-item').contains(tag);
        cy.get('.tn-filters > :nth-child(6) > .btn').click();
      });

      it(`should reload the clusters table when filtered tags are removed`, () => {
        cy.intercept('GET', `/clusters?per_page=100&tags=${tag}`).as(
          'filterByTags'
        );
        cy.get('.dropdown-item').contains(tag).click();
        cy.wait('@filterByTags').then(() => {
          cy.get('.tn-clustername').should('have.length', 2);
        });

        cy.get('.tn-filters > :nth-child(6) > .btn').click();

        cy.intercept('DELETE', `${tag}`).as('firstTagRemoved');
        cy.get(
          `#cluster-${availableClustersId[0]} > .tn-cluster-tags > .tagify > .tagify__tag > .tagify__tag__removeBtn`
        ).click();
        cy.wait('@firstTagRemoved').then(() => {
          cy.get('.tn-clustername').should('have.length', 1);
        });

        cy.intercept('DELETE', `${tag}`).as('secondTagRemoved');

        cy.get(
          `#cluster-${availableClustersId[1]} > .tn-cluster-tags > .tagify > .tagify__tag > .tagify__tag__removeBtn`
        ).click();
        cy.wait('@secondTagRemoved').then(() => {
          cy.get('.dropdown-item').contains(tag).should('not.exist');
          cy.get('.tn-clustername').should(
            'have.length',
            availableClustersId.length
          );
        });
      });
    });
  });
});
