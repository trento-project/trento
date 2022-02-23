context('Homepage', () => {
  before(() => {
    cy.visit('/');
  });

  describe('The homepage has the global health component visible', () => {
    it('should display the global health chart', () => {
      cy.get('#homepage-component').contains('At a glance');
      cy.get('.health-summary-row').should('have.length', 3);
      cy.get('.health-summary-id').should('have.length', 4); // The first is the header
      cy.get('.health-summary-icon.text-success').should('have.length', 9);
      cy.get('.health-summary-icon.text-danger').should('have.length', 3);
    });
  });
});
