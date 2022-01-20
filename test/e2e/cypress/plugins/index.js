/// <reference types="cypress" />
// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)

/**
 * @type {Cypress.PluginConfig}
 */
// eslint-disable-next-line no-unused-vars

const http = require('http');

module.exports = (on, config) => {
  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config
  on('task', {
    startAgentHeartbeat(agents) {
      const {collector_host, collector_port, heartbeat_interval} = config.env
      agents.forEach((agentId) => {
        setInterval(() => {
          http.request({
            host: collector_host,
            path: `/api/hosts/${agentId}/heartbeat`,
            port: collector_port,
            method: 'POST'
          }).end()
        }, heartbeat_interval)
      })
      return null
    }
  })
}
