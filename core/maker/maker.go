package maker

import (
	"bahamut/core/component"
	"bahamut/core/docker"
	"bahamut/core/types"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type (
	Maker struct {
		ID             uuid.UUID
		Name           string
		Queue          *queue.Queue
		Db             map[uuid.UUID]*component.Component // TODO: Replace it with persistent datastore.
		OfType         types.MakerType
		ComponentCount uint64
	}
)

// Return [] of Maker
func New(Name string, OfType types.MakerType, NoOfMakers int) []*Maker {
	// TODO: Load from previous persistent queue stack, and Db.
	var Makers []*Maker

	for i := 0; i < NoOfMakers; i++ {
		Maker := &Maker{
			ID:             uuid.New(),
			Name:           Name,
			Queue:          queue.New(),
			Db:             make(map[uuid.UUID]*component.Component),
			OfType:         OfType,
			ComponentCount: 0,
		}
		Makers = append(Makers, Maker)
	}
	return Makers
}

func (r *Maker) Start() {
	for {
		if r.Queue.Len() != 0 {
			result := r.run()
			if result.Error != nil {
				log.Print("Encountered error running component: %v", result.Error)
			}
		} else {
			log.Printf("No component scheduled to run at Maker: %v", r.Name)
		}
		// log.Println("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}

func (r *Maker) run() docker.DockerResult {
	comp := r.Queue.Dequeue()

	// interface to component type.
	compQueued := comp.(component.Component)
	log.Println("Processing a component in Queued, with state: %v", compQueued.State)

	compPersisted := r.Db[compQueued.ID] // if task exists in DB should be retirved then checked to see if Maker should transist it to recived state.

	if compPersisted == nil {
		compPersisted = &compQueued
	}

	r.sync(&compQueued)

	var result docker.DockerResult
	if component.ShouldTransition(compPersisted.State, compQueued.State) {
		switch compQueued.State {
		case component.Scheduled:
			result = r.startComponent(&compQueued)
		case component.Completed:
			result = r.stopComponent(&compQueued)
		default:
			result.Error = errors.New("Encountered unexpected component state transistion error.")
		}
	} else {
		err := fmt.Errorf("Invalid state transition from %v to %v", compPersisted.State, compQueued.State)
		result.Error = err
	}
	return result
}

// Start a component by calling its docker.Run method. and updating the state to state.Running.
func (r *Maker) startComponent(c *component.Component) docker.DockerResult {
	result := c.Docker.Run()
	if result.Error != nil {
		log.Print("Encountered error running component: %v", result.Error)
		c.State = component.Failed
		r.sync(c)
		return result
	}
	c.State = component.Running
	log.Print("Running compenet with state %v", c.State)
	r.sync(c)
	return result
}

// Stop a component by running its docker.Stop method. and updating the state to state.Completed.
func (r *Maker) stopComponent(c *component.Component) docker.DockerResult {
	result := c.Docker.Stop()
	if result.Error != nil {
		log.Printf("Error stopping component %s: %v", c.Docker.ContainerId, result.Error)
		c.State = component.Completed
		r.sync(c)
		return result
	}
	c.State = component.Completed
	r.sync(c)
	log.Printf("Stopped and removed container %v for component %v", c.Docker.ContainerId, c.ID)
	return result
}

// Ensure that components states are up to date.
func (r *Maker) sync(c *component.Component) {
	// TODO: Should be rewritten to sync to persisitent datastore.
	r.Db[c.ID] = c
}

// Schedule a component to run, State.Pending -> State.Scheduled
func (r *Maker) Schedule(c component.Component) component.Component {
	if c.State == component.Pending {
		c.State = component.Scheduled
		r.ComponentCount += 1
	}
	r.Queue.Enqueue(c)
	return c
}
