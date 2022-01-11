context('Trento About page', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should contain all relevant information', () => {
    cy.get('.js-sidebar-toggle').click();
    cy.get('.menu-title').contains('Settings').click();
    cy.get('.menu-title').contains('About').click();
    cy.url().should('include', '/about');


    cy.get('dl').should('contain', 'Trento flavor');
    cy.get('dl').should('contain', 'Server version');
    cy.get('dl').should('contain', 'Github repository');
    cy.get('dl').should('contain', 'https://github.com/trento-project/trento');
    cy.get('dl').should('contain', 'SLES for SAP subscriptions');
  });
});
