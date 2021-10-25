package datapipeline

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// TODO: tune workersNumber
var workersNumber int64 = 100
var drainTimeout = time.Second * 5

type ProjectorsWorkerPool struct {
	ch                 chan *DataCollectedEvent
	projectorsRegistry ProjectorRegistry
}

func NewProjectorsWorkerPool(projectorsRegistry ProjectorRegistry) *ProjectorsWorkerPool {
	return &ProjectorsWorkerPool{
		projectorsRegistry: projectorsRegistry,
		ch:                 make(chan *DataCollectedEvent),
	}
}

// Run runs a pool of workers to process events
func (p *ProjectorsWorkerPool) Run(ctx context.Context) {
	log.Infof("Starting projector pool. Workers limit: %d", workersNumber)
	sem := semaphore.NewWeighted(workersNumber)

	for {
		select {
		case event := <-p.ch:
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Debugf("Discarding event: %d, shutting down already.", event.ID)
				break
			}
			log.Infof("Projecting event: %d", event.ID)

			go func() {
				defer sem.Release(1)
				for _, projector := range p.projectorsRegistry {
					projector.Project(event)
				}
			}()
		case <-ctx.Done():
			log.Infof("Projectors worker pool is shutting down... Waiting for active workers to drain.")

			ctx, cancel := context.WithTimeout(context.Background(), drainTimeout)
			defer cancel()

			if err := sem.Acquire(ctx, workersNumber); err != nil {
				log.Warnf("Timed out while draining workers: %v", err)
			}

			return
		}
	}
}

// GetChannel returns the channel used by the worker pool
func (p *ProjectorsWorkerPool) GetChannel() chan *DataCollectedEvent {
	return p.ch
}
