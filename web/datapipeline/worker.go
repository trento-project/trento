package datapipeline

import (
	"context"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"gorm.io/gorm"
)

// TODO: tune bufferSize and workersNumber
var workersNumber int64 = 100
var bufferSize int64 = workersNumber * 100

func initProjectorsRegistry(db *gorm.DB) []*Projector {
	clusterListProjector := NewProjector("cluster_list", db)
	clusterListProjector.AddHandler(ClusterDiscovery, ClusterListHandler)

	return []*Projector{
		clusterListProjector,
	}
}

// StartProjectorsWorkerPool starts a pool of workers to process events
func StartProjectorsWorkerPool(db *gorm.DB) chan *DataCollectedEvent {
	ch := make(chan *DataCollectedEvent, bufferSize)
	projectorsRegistry := initProjectorsRegistry(db)

	log.Infof("Starting projector pool. Workers limit: %d", workersNumber)
	go workerPool(ch, projectorsRegistry)

	return ch
}

// workerPool starts a worker everytime a new event is received
// and limits concurrency to workersNumber by using a semaphore
func workerPool(ch chan *DataCollectedEvent, projectorsRegistry []*Projector) {
	ctx := context.Background()
	sem := semaphore.NewWeighted(workersNumber)

	for event := range ch {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
		}
		log.Debugf("Semaphore acquired, starting projector worker")

		go func(event *DataCollectedEvent) {
			defer sem.Release(1)
			for _, projector := range projectorsRegistry {
				projector.Project(event)
			}
		}(event)
	}
}
