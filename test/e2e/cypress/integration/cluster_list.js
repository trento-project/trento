context('Cluster list', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('should contain both hana and netweaver cluster', () => {
    cy.get('.js-sidebar-toggle').click();
    cy.get('.menu-title').contains('Pacemaker Clusters').click();
    cy.url().should('include', '/clusters');


    cy.get('td').should('contain', 'hana_cluster');
    cy.get('td').should('contain', 'netweaver_cluster');

    cy.get('table >tbody >tr td:nth-child(2)').each((el, index, _) => {
      const text = el.text();

      if (text.includes('hana_cluster')) {
        cy.get('table >tbody >tr td:nth-child(4)').eq(index).then(function (clusterType) {
          expect(clusterType.text().trim()).to.equal('HANA scale-up');
        })
        cy.get('table >tbody >tr td:nth-child(5)').eq(index).then(function (sid) {
          expect(sid.text().trim()).to.equal('PRD');
        })
      }

      if (text.includes('netweaver')) {
        cy.get('table >tbody >tr td:nth-child(4)').eq(index).then(function (clusterType) {
          expect(clusterType.text().trim()).to.equal('Unknown');
        })

      }
    })
  })
});
