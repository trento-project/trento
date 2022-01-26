context('Trento About page', () => {
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');
    cy.visit('/');
    cy.navigateToItem(['Settings', 'About']);
    cy.url().should('include', '/about');
  });

  it('should contain all relevant information', () => {
    cy.get('dl').should('contain', 'Trento flavor');
    cy.get('dl').should('contain', 'Community');

    cy.get('dl').should('contain', 'Server version');
    cy.exec(`${Cypress.env('trento_binary')} version`).then(({ stdout }) => {
      const version = stdout.split('version ').pop().split('\n')[0];

      cy.get('dl').should('contain', version);
    });

    cy.get('dl').should('contain', 'Github repository');
    cy.get('dl').should('contain', 'https://github.com/trento-project/trento');
    cy.get('dl').should('contain', 'SLES for SAP subscriptions');
    cy.get('dl').should('contain', '27 Found');
  });
});
