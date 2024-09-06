package ai

import (
    "fmt"
    "math"
    "sort"
    "strings"
)

// GOAPState represents the world state as a map of string keys to interface{} values
type GOAPState map[string]interface{}

// GOAPAction represents an action that can be taken in the GOAP system
type GOAPAction struct {
    Name          string
    CostFunc      func(state GOAPState) float64
    Preconditions GOAPState
    Effects       GOAPState
}

// GOAPPlanner is the main struct for the GOAP system
type GOAPPlanner struct {
    Actions []GOAPAction
}

// NewGOAPPlanner creates a new GOAPPlanner
func NewGOAPPlanner() *GOAPPlanner {
    return &GOAPPlanner{
        Actions: make([]GOAPAction, 0),
    }
}

// AddAction adds a new action to the planner
func (p *GOAPPlanner) AddAction(action GOAPAction) {
    p.Actions = append(p.Actions, action)
}

// Plan finds the optimal sequence of actions to reach the goal state from the start state
func (p *GOAPPlanner) Plan(start, goal GOAPState) []GOAPAction {
    openList := []GOAPState{start}
    cameFrom := make(map[string]GOAPState)
    gScore := make(map[string]float64)
    fScore := make(map[string]float64)
    actionTaken := make(map[string]GOAPAction)

    startKey := stateToString(start)
    gScore[startKey] = 0
    fScore[startKey] = p.heuristic(start, goal)

    for len(openList) > 0 {
        current := p.lowestFScore(openList, fScore)
        currentKey := stateToString(current)
        if p.stateEquals(current, goal) {
            return p.reconstructPath(cameFrom, actionTaken, current)
        }

        openList = p.removeState(openList, current)

        for _, action := range p.Actions {
            if p.canExecuteAction(current, action) {
                neighbor := p.applyEffects(current, action.Effects)
                neighborKey := stateToString(neighbor)
                tentativeGScore := gScore[currentKey] + action.CostFunc(current)

                if _, exists := gScore[neighborKey]; !exists || tentativeGScore < gScore[neighborKey] {
                    cameFrom[neighborKey] = current
                    actionTaken[neighborKey] = action
                    gScore[neighborKey] = tentativeGScore
                    fScore[neighborKey] = gScore[neighborKey] + p.heuristic(neighbor, goal)
                    if !p.containsState(openList, neighbor) {
                        openList = append(openList, neighbor)
                    }
                }
            }
        }
    }

    return nil // No plan found
}

// PlanWithStateSelection plans actions using state-dependent action selection
func (p *GOAPPlanner) PlanWithStateSelection(start, goal GOAPState) []GOAPAction {
    relevantActions := p.SelectActions(start)
    oldActions := p.Actions
    p.Actions = relevantActions
    plan := p.Plan(start, goal)
    p.Actions = oldActions
    return plan
}

// SelectActions filters actions based on the current state
func (p *GOAPPlanner) SelectActions(state GOAPState) []GOAPAction {
    selectedActions := make([]GOAPAction, 0)
    for _, action := range p.Actions {
        if p.isActionRelevant(state, action) {
            selectedActions = append(selectedActions, action)
        }
    }
    return selectedActions
}

// isActionRelevant determines if an action is relevant based on the current state
func (p *GOAPPlanner) isActionRelevant(state GOAPState, action GOAPAction) bool {
    // This is a placeholder implementation. You should customize this based on your game's logic.
    return true
}

// heuristic estimates the cost from a state to the goal state
func (p *GOAPPlanner) heuristic(state, goal GOAPState) float64 {
    // Simple heuristic: count the number of mismatched goals
    count := 0
    for k, v := range goal {
        if state[k] != v {
            count++
        }
    }
    return float64(count)
}

// lowestFScore finds the state with the lowest fScore in the open list
func (p *GOAPPlanner) lowestFScore(list []GOAPState, fScore map[string]float64) GOAPState {
    lowestF := math.Inf(1)
    var lowestState GOAPState
    for _, state := range list {
        key := stateToString(state)
        if f, exists := fScore[key]; exists && f < lowestF {
            lowestF = f
            lowestState = state
        }
    }
    return lowestState
}

// stateEquals checks if two states are equal
func (p *GOAPPlanner) stateEquals(current, goal GOAPState) bool {
    for k, v := range goal {
        if bv := current[k]; bv != v {
            return false
        }
    }
    return true
}

// reconstructPath builds the sequence of actions from start to goal
func (p *GOAPPlanner) reconstructPath(cameFrom map[string]GOAPState, actionTaken map[string]GOAPAction, current GOAPState) []GOAPAction {
    var path []GOAPAction
    for {
        currentKey := stateToString(current)
        action, exists := actionTaken[currentKey]
        if !exists {
            break
        }
        path = append([]GOAPAction{action}, path...)
        current = cameFrom[currentKey]
    }
    return path
}

// removeState removes a state from the list of states
func (p *GOAPPlanner) removeState(list []GOAPState, state GOAPState) []GOAPState {
    for i, s := range list {
        if p.stateEquals(s, state) {
            return append(list[:i], list[i+1:]...)
        }
    }
    return list
}

// canExecuteAction checks if an action's preconditions are met in the given state
func (p *GOAPPlanner) canExecuteAction(state GOAPState, action GOAPAction) bool {
    for k, v := range action.Preconditions {
        if state[k] != v {
            return false
        }
    }
    return true
}

// applyEffects applies an action's effects to a state
func (p *GOAPPlanner) applyEffects(state GOAPState, effects GOAPState) GOAPState {
    newState := make(GOAPState)
    for k, v := range state {
        newState[k] = v
    }
    for k, v := range effects {
        newState[k] = v
    }
    return newState
}

// containsState checks if a state is in the list of states
func (p *GOAPPlanner) containsState(list []GOAPState, state GOAPState) bool {
    for _, s := range list {
        if p.stateEquals(s, state) {
            return true
        }
    }
    return false
}

// stateToString converts a GOAPState to a string representation
func stateToString(state GOAPState) string {
    var pairs []string
    for k, v := range state {
        pairs = append(pairs, fmt.Sprintf("%s:%v", k, v))
    }
    sort.Strings(pairs) // Sort to ensure consistent string representation
    return strings.Join(pairs, "|")
}
