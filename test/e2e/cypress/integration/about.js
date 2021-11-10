context('Trento About page', () => {
  beforeEach(() => {
    cy.visit('http://demo.trento-project.io');
  });

  it('should contain all relevant information', () => {
    cy.get('.js-sidebar-toggle').click();
    cy.get('.menu-title').contains('About').click();
    cy.url().should('include', '/about');


    cy.get('dl').should('contain', 'Web version');
    cy.get('dl').should('contain', 'Github repository');
    cy.get('dl').should('contain', 'https://github.com/trento-project/trento');
    cy.get('dl').should('contain', 'Web version');
    cy.get('dl').should('contain', 'SLES_SAP machines');
    cy.get('dl').should('contain', 'Subscription');
    cy.get('dl').should('contain', 'Premium');
  });
});
