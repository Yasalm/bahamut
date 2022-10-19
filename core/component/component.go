package component

/*
this package is responsible for defining what a Component should be
its purpose and how its configured and monitored
purpuse: Database, BI tool, integration tool, and other specific purpueses in defining and building data products
should a run time engine be included within this category of Component

what defines an Component? for example if you are tasked to manually create install a metabase tool what your intented to do is not simply to copy docker image and run an voila it done? no you need to encoberate it with other tools, it needs certian aspects of its own to run. besides enviroment variables
*/

import (
	"bahamut/core/docker"
	"bahamut/core/types"

	"github.com/google/uuid"
)

type Api struct {
	Name string
}

/*
Need to address dependices for example Metabase to be production ready it needs to have a database. Wheather to use its own or the main USER database
*/
// because it is of type database it loads implementation of the interface of given type.
type Component struct {
	Name   string // this name will be used as docker name added UUID.
	ID     uuid.UUID
	State  State
	Docker *docker.Docker
	OfType types.MakerType // TODO: add logic to handle specific Component behavior or features such as Component:Integration should be differnent than Component:Database in method they perform.
	// Api Api // TODO: add Api to expose both,  1) Our own internal features to handle this Component, 2)  Api of that given service for instance metabase to query or post questions
}

func New(Name string, OfType types.MakerType, options types.Options) *Component {
	config := docker.NewConig(Name, options.Image, options.Env, options.RestartPolicy)
	docker := docker.NewDocker(config)

	comp := &Component{
		Name:   Name,
		ID:     uuid.New(),
		State:  Pending,
		Docker: docker,
		OfType: OfType,
	}
	return comp
}

func (comp *Component) ConfigureDocker() {
	// TODO: Configuration to the Component.
}

// Append the Client instance to the Component
func (comp *Component) Init(config docker.Config) {
	docker := docker.NewDocker(&config)
	comp.Docker = docker
}
