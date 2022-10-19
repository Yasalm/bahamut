package bahamut

import (
	"bahamut/core/component"
	"bahamut/core/maker"
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
		Makers           map[types.MakerType][]*maker.Maker
		Db               map[uuid.UUID]*component.Component // TODO: replace with persistent datastore.
		MakersComponents map[uuid.UUID][]uuid.UUID
		CompoenentMakers map[uuid.UUID]uuid.UUID
		Assigned         map[uuid.UUID]uuid.UUID // component ID to Makers ID
	}
)

// Create an instance of bahamut
func New(settings types.BahamutSettings) *Bahamut {
	db := make(map[uuid.UUID]*component.Component)
	compoenentMakers := make(map[uuid.UUID]uuid.UUID)
	MakersComponents := make(map[uuid.UUID][]uuid.UUID)
	Makers := make(map[types.MakerType][]*maker.Maker)
	assigned := make(map[uuid.UUID]uuid.UUID)

	for typeOfMakers, nOfMakers := range settings {
		m := maker.New("", typeOfMakers, nOfMakers)
		// Check if Makers are running, ideally it needs to be rewritten.
		for _, m_ := range m {
			go m_.Start() // We should also consider writing a method on Makers that check whether it is running or not. This will be visited once we write stats package.§
		}
		Makers[typeOfMakers] = m
	}

	return &Bahamut{
		ID:               uuid.New(),
		Pending:          queue.New(),
		Makers:           Makers,
		Db:               db,
		MakersComponents: MakersComponents,
		CompoenentMakers: compoenentMakers,
		Assigned:         assigned,
	}
}

// Allow Bahamut to schedule & dispatch components.
func (b *Bahamut) Start() {
	for {
		if b.Pending.Len() != 0 {
			fmt.Print("[Bahamut] Processing Queued Components %$v\n", b.Pending.Len())
			b.dispatch()
		}
		// log.Printf("[Bahamut] Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}

// Dispatch a component from pending queue to appropriate Makers
func (b *Bahamut) dispatch() {
	// Since we arrive after a component is queued to the Pending queue. we can assumed it has been assigned to a Makers through our simple -yet- method.
	// Given these assumptions we can dequeue the component fron Pending list
	// and pass its ID to assigned map which would give us the Makers that is assigned to it.
	enQueued := b.Pending.Dequeue()
	comp := enQueued.(component.Component) // type casting given Queue package return an interface.

	assignedMakersID := b.Assigned[comp.ID]
	gMakers := b.Makers[comp.OfType]
	var Makers *maker.Maker
	for _, r := range gMakers {
		if r.ID == assignedMakersID {
			Makers = r
		}
	}
	log.Printf("Dispatching component %v to Makers Group %v:", comp.Name,
		Makers.OfType)

	Makers.Schedule(comp)

}

// Assign a component to appropriate Makers by checking the OfType fields
func (b *Bahamut) assign(c component.Component) {
	groupMakers := b.Makers[c.OfType]
	selectedMakers := b.selectMakers(groupMakers)
	b.Assigned[c.ID] = selectedMakers
}

// medicore way of selecting a Makers
// TODO: Implement a better algo to select Makers ideally by checking its stats and load.
// But thats shall be done in optimization phase.
func (b *Bahamut) selectMakers(toSelect []*maker.Maker) uuid.UUID {
	randMakersIndex := rand.Intn(len(toSelect))
	selectedMakers := toSelect[randMakersIndex]
	return selectedMakers.ID
}

// Schedule a component to be dispatched to the appropriate Makers.
func (b *Bahamut) Schedule(c component.Component) {
	b.Pending.Enqueue(c)
	b.assign(c)
}

func (b *Bahamut) sync() {
	// TODO: Sync bahamut state and its Makers to persistent DB.
}
