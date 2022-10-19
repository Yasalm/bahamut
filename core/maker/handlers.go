package maker

import (
	"bahamut/core/component"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) scheduleHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	comp := component.Component{}

	err := d.Decode(&comp)

	if err != nil {
		msg := fmt.Sprintf("Error decoding component body, %v", err)
		log.Print(msg)

		w.WriteHeader(400)
		e := ErrResposne{
			HTTPStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	// TODO: fix pass by value in schedlue. this is solved by rewritten how to accpert components as events that desired state and id reference to the component id instead of passing the whole component back and forth.
	api.Maker.Schedule(comp)
	log.Printf("Added Component %v\n", comp.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(comp)
}

func (api *Api) stopHandler(w http.ResponseWriter, r *http.Request) {
	compID := chi.URLParam(r, "ID")
	if compID == "" {
		log.Print("No ID is provided")
		w.WriteHeader(400)
	}

	compId, _ := uuid.Parse(compID)
	comp := api.Maker.Db[compId]
	compCopy := *comp
	compCopy.State = component.Completed
	api.Maker.Schedule(compCopy)

	log.Printf("Scheduled a component:%v to stop at container %v", compId, comp.Docker.ContainerId)
	w.WriteHeader(204)
}

func (api *Api) getStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (api *Api) getComponents(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
