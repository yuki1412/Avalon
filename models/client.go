package models

import (
	"sort"
)

type Clients struct {
	clients []uint32
}

func ClientsInit(clients []uint32) Clients {
	return Clients{
		clients: clients,
	}
}

func (c *Clients) GetClientList() []uint32 {
	keys := make([]uint32, 0, len(c.clients))
	for k := range c.clients {
		keys = append(keys, uint32(k))
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return c.clients
}

func (c *Clients) Len() int {
	return len(c.clients)
}
