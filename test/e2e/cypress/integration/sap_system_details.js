import { agents } from '../fixtures/hosts-overview/available_hosts';
import {
  selectedSystem,
  attachedHosts,
} from '../fixtures/sap-system-details/selected_system';

context('SAP system details', () => {
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');

    cy.task('startAgentHeartbeat', agents());

    cy.visit(`/sapsystems/${selectedSystem.Id}`);
    cy.url().should('include', `/sapsystems/${selectedSystem.Id}`);
  });

  describe('SAP system details page is available', () => {
    it(`should display the "${selectedSystem.Sid}" system details page`, () => {
      cy.visit(`/sapsystems/${selectedSystem.Id}`);
      cy.url().should('include', `/sapsystems/${selectedSystem.Id}`);
      cy.get('h1').should('contain', 'SAP System details');
      cy.get('dd').eq(0).should('contain', selectedSystem.Sid);
      cy.get('dd').eq(1).should('contain', selectedSystem.Type);
    });

    it(`should display "Not found" page when SAP system doesn't exist`, () => {
      cy.visit(`/sapsystems/other`, { failOnStatusCode: false });
      cy.url().should('include', `/sapsystems/other`);
      cy.get('h1').should('contain', 'Not Found');
      cy.get('p').should('contain', "The requested URL doesn't exist");
    });
  });

  describe('The system layout shows all the running instances', () => {
    before(() => {
      cy.visit(`/sapsystems/${selectedSystem.Id}`);
      cy.url().should('include', `/sapsystems/${selectedSystem.Id}`);
    });

    selectedSystem.Hosts.forEach((instance, index) => {
      it(`should show hostname "${instance.Hostname}" with the correct values`, () => {
        cy.get('.eos-table')
          .eq(0)
          .find('tr')
          .eq(index + 1)
          .find('td')
          .as('tableCell');
        cy.get('@tableCell').eq(0).should('contain', instance.Hostname);
        cy.get('@tableCell').eq(1).should('contain', instance.Instance);
        cy.get('@tableCell').eq(2).should('contain', instance.Features);
        cy.get('@tableCell').eq(3).should('contain', instance.HttpPort);
        cy.get('@tableCell').eq(4).should('contain', instance.HttpsPort);
        cy.get('@tableCell').eq(5).should('contain', instance.StartPriority);
        cy.get('@tableCell').eq(6).should('contain', instance.Status);
        cy.get('@tableCell')
          .eq(6)
          .find('span')
          .should('have.class', instance.StatusBadge);
      });
    });

    const states = [
      ['GRAY', 'secondary'],
      ['GREEN', 'primary'],
      ['YELLOW', 'warning'],
      ['RED', 'danger'],
    ];

    states.forEach(([state, badge]) => {
      it(`should show ${state} badge in instace when SAPControl-${state} state is received`, () => {
        cy.loadScenario(`sap-system-detail-${state}`);
        cy.visit(`/sapsystems/${selectedSystem.Id}`);
        // using row 3 as the changed instance is the 3rd in order based on instance_number
        cy.get('.eos-table').eq(0).find('tr').eq(3).find('td').as('tableCell');
        cy.get('@tableCell').eq(6).should('contain', `SAPControl-${state}`);
        cy.get('@tableCell')
          .eq(6)
          .find('span')
          .should('have.class', `badge-${badge}`);
      });
    });

    it(`should show a new instance when an event with a new SAP instance is received`, () => {
      cy.loadScenario(`sap-system-detail-NEW`);
      cy.visit(`/sapsystems/${selectedSystem.Id}`);
      cy.get('.eos-table').eq(0).find('tr').should('have.length', 6);
      cy.get('.eos-table').eq(0).find('tr').eq(-1).find('td').as('tableCell');
      cy.get('@tableCell').eq(0).should('contain', 'newinstance');
      cy.get('@tableCell').eq(1).should('contain', '99');
    });
  });

  describe('The hosts table shows the attached hosts to this SAP system', () => {
    attachedHosts.forEach((host, index) => {
      it(`should show ${host.Name} with the correct link and data`, () => {
        cy.get('.eos-table')
          .eq(1)
          .find('tr')
          .eq(index + 1)
          .find('td')
          .as('tableCell');
        cy.get('@tableCell').eq(1).should('contain', host.Name);
        cy.get('@tableCell')
          .eq(1)
          .find('a')
          .should('have.attr', 'href', `/hosts/${host.AgentId}`);
        cy.get('@tableCell').eq(2).contains(host.Address);
        cy.get('@tableCell').eq(3).should('contain', host.Provider);
        cy.get('@tableCell').eq(4).should('contain', host.Cluster);
        cy.get('@tableCell').eq(5).should('contain', host.Version);
      });
    });
  });
});
