package component

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Failed
	Completed
)

//TODO: State handling helper functions
var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

// validate passed state should transisted to.
func Validate(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

/// validate if it should transition
func ShouldTransition(from State, to State) bool {
	return Validate(stateTransitionMap[from], to)
}
