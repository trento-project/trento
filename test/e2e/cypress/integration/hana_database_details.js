import { agents } from '../fixtures/hosts-overview/available_hosts';
import {
  selectedDatabase,
  attachedHosts,
} from '../fixtures/hana-database-details/selected_database';

context('HANA database details', () => {
  before(() => {
    cy.resetDatabase();
    cy.loadScenario('healthy-27-node-SAP-cluster');

    cy.task('startAgentHeartbeat', agents());

    cy.visit(`/databases/${selectedDatabase.Id}`);
    cy.url().should('include', `/databases/${selectedDatabase.Id}`);
  });

  describe('HANA database details page is available', () => {
    it(`should display the "${selectedDatabase.Sid}" database details page`, () => {
      cy.visit(`/databases/${selectedDatabase.Id}`);
      cy.url().should('include', `/databases/${selectedDatabase.Id}`);
      cy.get('h1').should('contain', 'HANA Database details');
      cy.get('dd').eq(0).should('contain', selectedDatabase.Sid);
      cy.get('dd').eq(1).should('contain', selectedDatabase.Type);
    });

    it(`should display "Not found" page when HANA database doesn't exist`, () => {
      cy.visit(`/databases/other`, { failOnStatusCode: false });
      cy.url().should('include', `/databases/other`);
      cy.get('h1').should('contain', 'Not Found');
      cy.get('p').should('contain', "The requested URL doesn't exist");
    });
  });

  describe('The database layout shows all the running instances', () => {
    before(() => {
      cy.visit(`/databases/${selectedDatabase.Id}`);
      cy.url().should('include', `/databases/${selectedDatabase.Id}`);
    });

    selectedDatabase.Hosts.forEach((instance, index) => {
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
      it(`should show ${state} badge in instace when SAPControl-GRAY state is received`, () => {
        cy.loadScenario(`hana-database-detail-${state}`);
        cy.visit(`/databases/${selectedDatabase.Id}`);
        // using row 2 as the changed instance is the 2nd in order based on instance_number
        cy.get('.eos-table').eq(0).find('tr').eq(2).find('td').as('tableCell');
        cy.get('@tableCell').eq(6).should('contain', `SAPControl-${state}`);
        cy.get('@tableCell')
          .eq(6)
          .find('span')
          .should('have.class', `badge-${badge}`);
      });
    });

    it(`should show a new instance when an event with a new SAP instance is received`, () => {
      cy.loadScenario(`hana-database-detail-NEW`);
      cy.visit(`/databases/${selectedDatabase.Id}`);
      cy.get('.eos-table').eq(0).find('tr').should('have.length', 4);
      cy.get('.eos-table').eq(0).find('tr').eq(-1).find('td').as('tableCell');
      cy.get('@tableCell').eq(0).should('contain', 'newinstance');
      cy.get('@tableCell').eq(1).should('contain', '99');
    });
  });

  describe('The hosts table shows the attached hosts to this HANA database', () => {
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
        cy.get('@tableCell')
          .eq(4)
          .find('a')
          .should('have.attr', 'href', `/clusters/${host.ClusterId}`);
        cy.get('@tableCell').eq(5).should('contain', host.Version);
      });
    });
  });
});
