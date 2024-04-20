package consts

var (
	GOOD_ROLE = []string{
		"Merlin",
		"Percival",
		"Loyal Servants of Arthur",
		"Loyal Servants of Arthur",
		"Loyal Servants of Arthur",
		"Loyal Servants of Arthur",
		"Loyal Servants of Arthur",
	}
	EVIL_ROLE = []string{
		"Assassin",
		"Mordred",
		"Morgana",
		"Oberon",
		"Minions of Mordred",
		"Minions of Mordred",
		"Minions of Mordred",
	}
	ROLE_AVAILABLE   = append(GOOD_ROLE, EVIL_ROLE...)
	ROLE_SHUFFLE_MAP = map[int]map[string]int{
		5: {
			"GOOD": 3,
			"EVIL": 2,
		},
		6: {
			"GOOD": 4,
			"EVIL": 2,
		},
		7: {
			"GOOD": 4,
			"EVIL": 3,
		},
		8: {
			"GOOD": 5,
			"EVIL": 3,
		},
		9: {
			"GOOD": 6,
			"EVIL": 3,
		},
		10: {
			"GOOD": 6,
			"EVIL": 4,
		},
	}
	ROUND_ASSIGNMENT_MAP = map[int][]int{
		5:  {2, 3, 2, 3, 3},
		6:  {2, 3, 4, 3, 4},
		7:  {2, 3, 3, 4, 4},
		8:  {3, 4, 4, 5, 5},
		9:  {3, 4, 4, 5, 5},
		10: {3, 4, 4, 5, 5},
	}
)
