import { selectedHost } from '../fixtures/host-details/selected_host';

context('Host Details', () => {
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');

    cy.task('startAgentHeartbeat', [selectedHost.agentId]);

    cy.visit('/');
    cy.navigateToItem('Hosts');
    cy.intercept('GET', `/hosts/${selectedHost.agentId}`).as('getHostDetails');
    cy.get(`#host-${selectedHost.hostName} > .tn-hostname > a`).click();
    cy.wait('@getHostDetails');
    cy.url().should('include', `/hosts/${selectedHost.agentId}`);
  });

  describe('Detailed view for a specific host should be available', () => {
    it('should show the host I clicked on in the overview', () => {
      cy.get('.tn-host-details-container .tn-hostname').should(
        'contain',
        selectedHost.hostName
      );
    });
  });

  describe('SAP instances for this host should be displayed', () => {
    it(`should show a link to the SAP System details view for ${selectedHost.sapSystem}`, () => {
      cy.get(':nth-child(2) > .text-muted > a')
        .should('contain', selectedHost.sapSystem)
        .invoke('attr', 'href')
        .should('include', `/sapsystems/${selectedHost.sapInstanceId}`);
    });
    it(`should show SAP instance with ID ${selectedHost.sapInstanceId}`, () => {
      cy.get(':nth-child(12) > .table > tbody > tr > :nth-child(1)').should(
        'contain',
        selectedHost.sapInstanceId
      );
    });
  });

  describe('Cluster details for this host should be displayed', () => {
    it(`should show a link to the cluster details view for ${selectedHost.clusterName}`, () => {
      cy.get(':nth-child(3) > .text-muted > a')
        .should('contain', selectedHost.clusterName)
        .invoke('attr', 'href')
        .should('include', `/clusters/${selectedHost.clusterId}`);
    });
  });

  describe('Cloud details for this host should be displayed', () => {
    it(`should show ${selectedHost.hostName} under the VM Name`, () => {
      cy.get(
        ':nth-child(6) > :nth-child(1) > .col-sm-12 > :nth-child(1) > :nth-child(2) > .text-muted'
      ).should('contain', selectedHost.hostName);
    });
    it(`should show ${selectedHost.resourceGroup} under the Resource group label`, () => {
      cy.get(
        ':nth-child(6) > :nth-child(1) > .col-sm-12 > :nth-child(1) > :nth-child(3) > .text-muted'
      ).should('contain', selectedHost.resourceGroup);
    });
  });

  describe("Trento agent status should be 'running'", () => {
    it("should show the status as 'running'", () => {
      cy.get('.badge').should('contain', 'running');
    });
  });
});
