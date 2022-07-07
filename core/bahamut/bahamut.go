package bahamut

import (
	"bahamut/core/component"
	"bahamut/core/runner"
	"bahamut/core/types"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type (
	Bahamut struct {
		ID               uuid.UUID
		Pending          *queue.Queue
		Runners          map[types.RunnerGroup][]*runner.Runner
		Db               map[uuid.UUID]*component.Component // TODO: replace with persistent datastore.
		RunnerComponents map[uuid.UUID][]uuid.UUID          // TODO
		CompoenentRunner map[uuid.UUID]uuid.UUID
		Assigned         map[uuid.UUID]uuid.UUID // component ID to Runner ID
	}
)

// Create an instance of bahamut
func New(settings types.BahamutSettings) *Bahamut {
	db := make(map[uuid.UUID]*component.Component)
	compoenentRunner := make(map[uuid.UUID]uuid.UUID)
	runnerComponents := make(map[uuid.UUID][]uuid.UUID)
	runners := make(map[types.RunnerGroup][]*runner.Runner)
	assigned := make(map[uuid.UUID]uuid.UUID)

	for typeOfRunner, nOfRunners := range settings {
		r := runner.New("", typeOfRunner, nOfRunners)
		// Check if runners are running, ideally it needs to be rewritten.
		for _, r_ := range r {
			go r_.Start() // We should also consider writing a method on runner that check whether it is running or not. This will be visited once we write stats package.
		}
		runners[typeOfRunner] = r
	}

	return &Bahamut{
		ID:               uuid.New(),
		Pending:          queue.New(),
		Runners:          runners,
		Db:               db,
		RunnerComponents: runnerComponents,
		CompoenentRunner: compoenentRunner,
		Assigned:         assigned,
	}
}

// Allow Bahamut to schedule & dispatch components.
func (b *Bahamut) Start() {
	for {
		if b.Pending.Len() != 0 {
			fmt.Print("[Bahamut] Processing Queued Components %$v\n", b.Pending.Len())
			b.Dispatch()
		}
		// log.Printf("[Bahamut] Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}

// Dispatch a component from pending queue to appropriate runner
func (b *Bahamut) Dispatch() {
	// Since we arrive after a component is queued to the Pending queue. we can assumed it has been assigned to a runner through our simple -yet- method.
	// Given these assumptions we can dequeue the component fron Pending list
	// and pass its ID to assigned map which would give us the runner that is assigned to it.
	enQueued := b.Pending.Dequeue()
	comp := enQueued.(component.Component) // type casting given Queue package return an interface.

	assignedRunnerID := b.Assigned[comp.ID]
	gRunners := b.Runners[comp.OfType]
	var runner *runner.Runner
	for _, r := range gRunners {
		if r.ID == assignedRunnerID {
			runner = r
		}
	}
	log.Printf("Dispatching component %v to runner Group %v:", comp.Name,
		runner.OfType)

	runner.Schedule(comp)

}

// Assign a component to appropriate runner by checking the OfType fields
func (b *Bahamut) Assign(c component.Component) {
	groupRunners := b.Runners[c.OfType]
	selectedRunner := b.SelectRunner(groupRunners)
	b.Assigned[c.ID] = selectedRunner
}

// medicore way of selecting a runner
// TODO: Implement a better algo to select runner ideally by checking its stats and load.
// But thats shall be done in optimization phase.
func (b *Bahamut) SelectRunner(toSelect []*runner.Runner) uuid.UUID {
	randRunnerIndex := rand.Intn(len(toSelect))
	selectedRunner := toSelect[randRunnerIndex]
	return selectedRunner.ID
}

// Schedule a component to be dispatched to the appropriate runner.
func (b *Bahamut) Schedule(c component.Component) {
	b.Pending.Enqueue(c)
	b.Assign(c)
}

func (b *Bahamut) Sync() {
	// TODO: Sync bahamut state and its runner to persistent DB.
}
