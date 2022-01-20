context('Trento About page', () => {
  before(() => {
    cy.visit('/');
    cy.navigateToItem(['Settings', 'About'])
    cy.url().should('include', '/about');
  });

  it('should contain all relevant information', () => {
    cy.get('dl').should('contain', 'Trento flavor');
    cy.get('dl').should('contain', 'Server version');
    cy.get('dl').should('contain', 'Github repository');
    cy.get('dl').should('contain', 'https://github.com/trento-project/trento');
    cy.get('dl').should('contain', 'SLES for SAP subscriptions');
  });
});
