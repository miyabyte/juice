package juice

import (
	"sync"
	"time"
)

var ucm *userClientManager
var OnlineUsers map[int]time.Time

var changeUserStatusChan chan struct {
	uid    int
	status bool
}

var onceUCM sync.Once

type userClientManager struct {
	ucs map[int]*UserClient
	// future : create a user lock-pool don't depend this
	sync.RWMutex
}

type UserClient struct {
	sync.Mutex

	Uid     int
	Clients map[uint32]*Client
}

func GetUserCliManager() *userClientManager {
	onceUCM.Do(func() {

		ucm = &userClientManager{
			ucs: make(map[int]*UserClient),
		}

		OnlineUsers = make(map[int]time.Time)
		changeUserStatusChan = make(chan struct {
			uid    int
			status bool
		}, 100)

		//状态更改器
		go ucm.userStatusCustomer()
	})
	return ucm
}

func (ucm *userClientManager) GetOnlineUsers() map[int]time.Time {
	ucm.RLock()
	defer ucm.RUnlock()

	return OnlineUsers
}

func (ucm *userClientManager) GetUserClients() map[int]*UserClient {
	ucm.RLock()
	defer ucm.RUnlock()

	return ucm.ucs
}

func (ucm *userClientManager) GetUserClient(uid int) (uc *UserClient, ok bool) {
	ucm.RLock()
	defer ucm.RUnlock()

	uc, ok = ucm.ucs[uid]
	return
}

func (ucm *userClientManager) getOnlineUsers() map[int]time.Time {
	return OnlineUsers
}

func (ucm *userClientManager) getUserClient(uid int) (ucs *UserClient, ok bool) {
	ucs, ok = ucm.ucs[uid]
	return
}

// up  observer mode	|	locked from cliManager ->func addClient
func (ucm *userClientManager) AddClient(cli *Client) {
	if uc, ok := ucm.getUserClient(cli.Uid); ok {
		uc.Clients[cli.UUID] = cli
	} else {
		ucm.addUserClient(cli)
	}
}
func (ucm *userClientManager) addUserClient(cli *Client) {
	ucm.ucs[cli.Uid] = &UserClient{
		Uid: cli.Uid,
		Clients: map[uint32]*Client{
			cli.UUID: cli,
		},
	}

	//user online
	ucm.getOnlineUsers()[cli.Uid] = time.Now()

	changeUserStatusChan <- struct {
		uid    int
		status bool
	}{uid: cli.Uid, status: true}
}

// down observer mode	|	locked from cliManager ->func addClient
func (ucm *userClientManager) RemoveClient(cli *Client) {
	delete(ucm.ucs[cli.Uid].Clients, cli.UUID)

	if len(ucm.ucs[cli.Uid].Clients) == 0 {
		//user offline
		delete(ucm.ucs, cli.Uid)
		delete(ucm.getOnlineUsers(), cli.Uid)
		changeUserStatusChan <- struct {
			uid    int
			status bool
		}{uid: cli.Uid, status: false}
	}
}

// future : buffer
func (ucm *userClientManager) userStatusCustomer() {
	for c := range changeUserStatusChan {
		if c.status {
			// user up
		} else {
			// user down
		}
	}
}
