package main

// global connection counter
var connIDCounter uint32
var status string = "lobby"
var gameStatus string = "lobby"
var voteList map[uint32]bool = make(map[uint32]bool)
var missionBoard []*bool = make([]*bool, 5)
var missionRound int = 0
var missionTrue bool = true
var missionFalse bool = false
var selectedClient map[uint32]bool = make(map[uint32]bool)
var roundAssigner uint32 = 0
