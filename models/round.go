package models

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	missionTrue  = true
	missionFalse = false
)

type Round struct {
	ID                 int             `json:"id"`
	Clients            Clients         `json:"clients"`
	RoundAssigner      uint32          `json:"round_assigner"`
	RoundAssignerIndex int             `json:"round_assigner_index"`
	MissionTracker     int             `json:"mission_tracker"`
	MissionBoard       []*bool         `json:"mission_board"`
	VoteList           map[string]bool `json:"vote_list"`
	Status             string          `json:"status"`
	SelectedClient     map[string]bool `json:"selected_client"`
}

func RoundInit(clients Clients) *Round {
	return &Round{
		Clients:        clients,
		MissionTracker: 0,
		MissionBoard:   make([]*bool, 5),
		Status:         "lobby",
		VoteList:       make(map[string]bool),
		SelectedClient: make(map[string]bool, clients.Len()),
	}
}

func (r *Round) InitRoundAssigner() {
	if r.RoundAssigner == 0 {
		r.RoundAssigner = r.Clients.clients[0]
	}
}

func (r *Round) UpdateStatus(status string) {
	r.Status = status
}

func (r *Round) ChangeStatus(msg string) {
	r.Status = msg
}

func (r *Round) ClearRecord(action string, gStatus string) {
	r.MissionBoard = make([]*bool, 5)
	r.MissionTracker = 0
	r.VoteList = make(map[string]bool)
	r.ChangeStatus(action)
	r.UpdateStatus(gStatus)
}

func (r *Round) StopGame() {
	r.ClearRecord("stop", "Game Stopped")
}

func (r *Round) nextRound() {
	if r.isMissionOver() {
		r.Status = "lobby"
		r.ChangeStatus(r.Status)
		return
	}

	r.MissionTracker++
	r.VoteList = make(map[string]bool)
	r.Status = "vote"
}

func (r *Round) isMissionOver() bool {
	return r.MissionTracker >= len(r.MissionBoard)-1
}

func (r *Round) handlerVote() string {
	// if all clients have voted, return the number of votes
	if r.IsVotingOver() {
		r.UpdateStatus("Vote Result")
	}

	return fmt.Sprintf("Votes: [%d/%d]\n", len(r.VoteList), r.Clients.Len())
}

func (r *Round) getVoteResult() string {
	var result string = "Success"
	r.MissionBoard[r.MissionTracker] = &missionTrue
	for _, v := range r.VoteList {
		if !v {
			result = "Failed"
			r.MissionBoard[r.MissionTracker] = &missionFalse
		}
	}
	return fmt.Sprintf("Vote Result: %s\n", result)
}

func (r *Round) AnyMajority() bool {
	trueCount := 0
	falseCount := 0
	for _, v := range r.MissionBoard {
		if v == nil {
			continue
		}
		if *v {
			trueCount++
		} else {
			falseCount++
		}
		if trueCount > len(r.MissionBoard)/2 || falseCount > len(r.MissionBoard)/2 {
			return true
		}
	}
	return false
}

func (r *Round) nextRoundAssigner() {
	// Round Robin Assign
	r.RoundAssignerIndex = (r.RoundAssignerIndex + 1) % r.Clients.Len()
	r.RoundAssigner = r.Clients.clients[r.RoundAssignerIndex]
}

func (r *Round) IsVotingOver() bool {
	return len(r.VoteList) >= r.Clients.Len()
}

func (r *Round) HandlerSelectedClient(client string) string {
	clientID, err := strconv.ParseUint(client, 10, 32)
	if err != nil {
		return ""
	}
	_, exist := r.SelectedClient[string(clientID)]
	if !exist {
		r.SelectedClient[string(clientID)] = true
	} else {
		delete(r.SelectedClient, string(clientID))
	}
	return fmt.Sprintf("clientList;%s", r.GetClientList())
}

func (r *Round) GetClientList() string {
	stringList := ""
	for _, value := range r.Clients.clients {
		stringList += fmt.Sprintf("Client ID: %d\n", value)
	}
	return fmt.Sprintf("Number of Players: %d\n%s", r.Clients.Len(), stringList)
}

func (r *Round) GetSystemInfo() string {
	stringList := fmt.Sprintf("Status:%s\n", strings.ToTitle(r.Status))
	stringList += fmt.Sprintf("GStatus:%s\n", strings.ToTitle(r.Status))
	stringList += fmt.Sprintf("Round:%d\n", r.MissionTracker+1)
	stringList += fmt.Sprintf("Round Assigner:%d\n", r.RoundAssigner)

	for _, id := range strings.Split(r.GetClientList(), ";") {
		if id == "true" || id == "false" {
			continue
		}
		stringList += fmt.Sprintf("- %s\n", id)
	}
	if r.Status != "start" && r.Status != "stop" {
		missionStatus := "Mission ["
		for _, v := range r.MissionBoard {
			if v == &missionTrue || v == &missionFalse {
				missionStatus += fmt.Sprintf("%t ", *v)
			} else {
				missionStatus += fmt.Sprintf("%t ", false)
			}
		}
		missionStatus += "]\n"
		stringList += missionStatus
	}
	return fmt.Sprintf("Number of Players: %d\n%s", r.Clients.Len(), stringList)
}
