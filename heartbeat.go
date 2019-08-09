package juice

import "time"

type heartbeat struct {
	HeartbeatCheckInterval time.Duration
	HeartbeatIdleTime      time.Duration
}

type source interface {
	Close(client *Client)
	GetClients() map[uint32]*Client
}

func (h *heartbeat) run(s source) error {
	go func() {
		for range time.Tick(h.HeartbeatCheckInterval) {
			if err := h.check(s); err != nil {
				panic(err)
			}
		}
	}()
	return nil
}

//记录检查一次的消耗时间
func (h *heartbeat) check(s source) error {
	for _, cli := range s.GetClients() {
		if time.Now().Sub(cli.LastTime) > h.HeartbeatIdleTime {
			s.Close(cli)
		}
	}
	return nil
}
