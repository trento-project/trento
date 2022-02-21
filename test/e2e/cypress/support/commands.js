// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })

const initializeOpenSidebar = () => cy.setCookie('collapsedSidebar', 'false');

const selectPagination = (itemsPerPage) => {
  const pagination = [10, 25, 50, 100];
  cy.get('.pagination-actions button.dropdown-toggle').click();
  cy.get(
    `.pagination-actions .dropdown-menu .dropdown-item:nth-child(${
      pagination.indexOf(itemsPerPage) + 1
    })`
  ).click();
  // eslint-disable-next-line cypress/no-unnecessary-waiting
  cy.wait(100);
};

Cypress.Commands.add('navigateToItem', (item) => {
  initializeOpenSidebar();
  const items = Array.isArray(item) ? item : [item];
  items.forEach((it) => cy.get('.menu-title').contains(it).click());
});

Cypress.Commands.add('reloadList', (listName, itemsPerPage) => {
  cy.intercept('GET', `/${listName}?per_page=${itemsPerPage}`).as('reloadList');
  selectPagination(itemsPerPage);
  cy.wait('@reloadList');
});

Cypress.Commands.add('resetDatabase', () => {
  cy.log('Resetting DB...');
  cy.exec(
    `${Cypress.env('trento_binary')} ctl db-reset --db-host=${Cypress.env(
      'db_host'
    )} --db-port=${Cypress.env('db_port')}`
  );
});

Cypress.Commands.add('pruneChecksResults', () => {
  cy.log('Resetting DB...');
  cy.exec(
    `${Cypress.env(
      'trento_binary'
    )} ctl prune-checks-results --db-host=${Cypress.env(
      'db_host'
    )} --db-port=${Cypress.env('db_port')}`
  );
});

Cypress.Commands.add('loadScenario', (scenario) => {
  const [fixturesPath, photofinishBinary, collectorHost, collectorPort] = [
    Cypress.env('fixtures_path'),
    Cypress.env('photofinish_binary'),
    Cypress.env('collector_host'),
    Cypress.env('collector_port'),
  ];
  cy.log(`Loading scenario "${scenario}"...`);
  cy.exec(
    `cd ${fixturesPath} && ${photofinishBinary} run --url "http://${collectorHost}:${collectorPort}/api/collect" ${scenario}`
  );
});

Cypress.Commands.add('loadChecksCatalog', (catalog) => {
  const [webApiHost, webApiPort] = [
    Cypress.env('web_api_host'),
    Cypress.env('web_api_port'),
  ];
  cy.log(`Loading checks catalog "${catalog}"...`);
  cy.fixture(catalog).then((file) => {
    cy.request({
      method: 'PUT',
      url: `http://${webApiHost}:${webApiPort}/api/checks/catalog`,
      body: file,
    });
  });
});

Cypress.Commands.add('loadChecksResults', (results, clusterId) => {
  const [webApiHost, webApiPort] = [
    Cypress.env('web_api_host'),
    Cypress.env('web_api_port'),
  ];
  cy.log(`Loading checks results "${results}"...`);
  cy.fixture(results).then((file) => {
    cy.request({
      method: 'POST',
      url: `http://${webApiHost}:${webApiPort}/api/checks/${clusterId}/results`,
      body: file,
    });
  });
});
