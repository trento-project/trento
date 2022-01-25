import { availableChecks } from '../fixtures/checks-catalog/available_checks'

context('Checks catalog', () => {
    before(() => {
        cy.resetDatabase()
        cy.loadChecksCatalog('checks-catalog/catalog.json')

        cy.visit('/catalog');
        cy.url().should('include', '/catalog');
    })

    describe('Checks catalog should be available', () => {
        it('should show 5 check groups in the catalog', () => {
            cy.get('div.check-group').should('have.length', 5)
        })
        it('should show 35 checks in the catalog', () => {
            cy.get('tr.check-row').should('have.length', 35)
        })
    })

    describe('Checks grouping and identification is correct', () => {
        availableChecks.forEach((checks, group) => {
            it(`should include group '${group}'`, () => {
                cy.get('.check-group > h4').should('contain', group)
            })
            checks.forEach((check_id) => {
                it(`should include check '${check_id}'`, () => {
                    cy.get('.check-row').should('contain', check_id)
                })
            })
        })
    })

    describe('Individual checks data is expanded', () => {
        it('should expand check data when clicked', () => {
          const firstCheck = availableChecks.get('Corosync')[0]
          cy.get('.check-row').filter(
            `[id="${firstCheck}"]`).find('td > a.link-dark').click()
          cy.get('.check-row').find(
            `#collapse-${firstCheck}`).should('be.visible')

          cy.get('.check-row').filter(
            `[id="${firstCheck}"]`).find('td > a.link-dark').click()
          cy.get('.check-row').find(
            `#collapse-${firstCheck}`).should('not.be.visible')
        })
        it('should expand check data when id is added in the url', () => {
            const firstCheck = availableChecks.get('Corosync')[0]
            cy.visit(`/catalog#${firstCheck}`)
            cy.get('.check-row').find(
              `#collapse-${firstCheck}`).should('be.visible')
        })
    })
})
