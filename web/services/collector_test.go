package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/datapipeline"
)

func TestCollectorService_StoreEvent(t *testing.T) {
	db := helpers.SetupTestDatabase()
	tx := db.Begin()
	tx.AutoMigrate(&datapipeline.DataCollectedEvent{})
	defer tx.Rollback()

	ch := make(chan *datapipeline.DataCollectedEvent, 1)
	collectorService := NewCollectorService(tx, ch)

	collectorService.StoreEvent(&datapipeline.DataCollectedEvent{
		AgentID:       "agent_id",
		DiscoveryType: "test_discovery_type",
		Payload:       []byte("{}"),
	})

	eventFromChannel := <-ch
	var eventFromDB datapipeline.DataCollectedEvent
	tx.First(&eventFromDB)

	assert.EqualValues(t, eventFromChannel.ID, eventFromDB.ID)
	assert.EqualValues(t, eventFromChannel.AgentID, eventFromDB.AgentID)
	assert.EqualValues(t, eventFromChannel.DiscoveryType, eventFromDB.DiscoveryType)
	assert.EqualValues(t, eventFromChannel.Payload, eventFromDB.Payload)
}
