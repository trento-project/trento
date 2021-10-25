package datapipeline

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

// TestProjectorWorkersPool tests that the worker pool correctly spawns workers
// when new events are added to the channel.
func TestProjectorWorkersPool(t *testing.T) {
	workersNumber = 2

	var wg sync.WaitGroup
	wg.Add(2)

	projector := new(MockProjector)
	projector.On("Project", mock.Anything).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(nil)

	projectorRegistry := []Projector{
		projector,
	}

	projectorsWorkersPool := NewProjectorsWorkerPool(projectorRegistry)
	ctx, cancel := context.WithCancel(context.Background())
	go projectorsWorkersPool.Run(ctx)

	ch := projectorsWorkersPool.GetChannel()
	ch <- &DataCollectedEvent{}
	ch <- &DataCollectedEvent{}

	wg.Wait()

	projector.AssertNumberOfCalls(t, "Project", 2)
	cancel()
}

// TestProjectorWorkersPool_BoundedParallelism tests that no more than the workersNumber limit
// of workers are spawned.
func TestProjectorWorkersPool_BoundedParallelism(t *testing.T) {
	workersNumber = 2
	quit := make(chan struct{})

	projector := new(MockProjector)
	projector.On("Project", mock.Anything).Run(func(args mock.Arguments) {
		<-quit
	}).Return(nil)

	projectorRegistry := []Projector{
		projector,
	}

	projectorsWorkersPool := NewProjectorsWorkerPool(projectorRegistry)
	ctx, cancel := context.WithCancel(context.Background())
	go projectorsWorkersPool.Run(ctx)

	go func() {
		ch := projectorsWorkersPool.GetChannel()
		ch <- &DataCollectedEvent{}
		ch <- &DataCollectedEvent{}
		ch <- &DataCollectedEvent{}
	}()

	time.Sleep(100 * time.Millisecond)
	projector.AssertNumberOfCalls(t, "Project", 2)

	quit <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	projector.AssertNumberOfCalls(t, "Project", 3)

	cancel()
}

// TestProjectorWorkersPool_Drain tests that the workers are drained when the context is canceled
// and that the worker pool shuts down gracefully.
func TestProjectorWorkersPool_Drain(t *testing.T) {
	workersNumber = 2
	drainTimeout = 200 * time.Millisecond
	done1 := false
	done2 := false

	startProcessing := make(chan struct{})

	projector := new(MockProjector)
	projector.On("Project", mock.Anything).Run(func(args mock.Arguments) {
		<-startProcessing
		time.Sleep(drainTimeout)
		if args.Get(0).(*DataCollectedEvent).ID == 1 {
			done1 = true
		} else {
			done2 = true
		}
	}).Return(nil)

	projectorRegistry := []Projector{
		projector,
	}

	projectorsWorkersPool := NewProjectorsWorkerPool(projectorRegistry)

	ctx, cancel := context.WithCancel(context.Background())
	go projectorsWorkersPool.Run(ctx)

	ch := projectorsWorkersPool.GetChannel()
	ch <- &DataCollectedEvent{ID: 1}
	ch <- &DataCollectedEvent{ID: 2}

	startProcessing <- struct{}{}
	startProcessing <- struct{}{}

	time.Sleep(100 * time.Millisecond)
	cancel()

	time.Sleep(drainTimeout)
	projector.AssertNumberOfCalls(t, "Project", 2)
	assert.True(t, done1)
	assert.True(t, done2)
}
