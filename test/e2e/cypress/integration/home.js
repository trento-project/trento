context('Homepage', () => {
  before(() => {
    cy.visit('/');
  });

  describe('The main links in the homepage should work (means the links are the expected ones)', () => {
    it('should provide correct link to Blue Horizon for SAP', () => {
      cy.get('#tn-blue-horizon-link')
        .should('have.attr', 'target', 'blank')
        .and(
          'have.attr',
          'href',
          'https://github.com/SUSE/blue-horizon-for-sap'
        );
    });
    it('should provide correct link to the Hosts page', () => {
      cy.get('#tn-hosts-link').should('have.attr', 'href', '/hosts');
    });
    it('should provide correct link to the Pacemaker Clusters', () => {
      cy.get('#tn-clusters-link').should('have.attr', 'href', '/clusters');
    });
    it('should provide correct link to the SAP Systems', () => {
      cy.get('#tn-sapsystems-link').should('have.attr', 'href', '/sapsystems');
    });
  });

  describe('The documentation links available in the homepage should work (means the links are the expected ones)', () => {
    it('should provide correct link to the Scope document', () => {
      cy.get('#tn-scope-link').should(
        'have.attr',
        'href',
        'https://github.com/trento-project/trento/blob/main/docs/scope.md'
      );
    });
    it('should provide correct link to the Readme', () => {
      cy.get('#tn-readme-link').should(
        'have.attr',
        'href',
        'https://github.com/trento-project/trento/blob/main/README.md'
      );
    });
    it('should provide correct link to the Architecture document', () => {
      cy.get('#tn-architecture-link').should(
        'have.attr',
        'href',
        'https://github.com/trento-project/trento/blob/main/docs/trento-architecture.md'
      );
    });
  });
});
