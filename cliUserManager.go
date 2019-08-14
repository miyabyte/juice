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

	uid     int
	clients map[uint32]*Client
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

func (ucm *userClientManager) GetUserClient(uid int) (ucs *UserClient, ok bool) {
	ucm.RLock()
	defer ucm.RUnlock()

	ucs, ok = ucm.ucs[uid]
	return
}

func getOnlineUsers() map[int]time.Time {
	return OnlineUsers
}

func getUserClient(uid int) (ucs *UserClient, ok bool) {
	ucs, ok = ucm.ucs[uid]
	return
}

func (ucm *userClientManager) AddUserClient(cli *Client) {
	ucm.ucs[cli.Uid] = &UserClient{
		uid: cli.Uid,
		clients: map[uint32]*Client{
			cli.UUID: cli,
		},
	}
	//user online
	changeUserStatusChan <- struct {
		uid    int
		status bool
	}{uid: cli.Uid, status: true}
}

// up  observer mode	|	locked from cliManager ->func addClient
func (ucm *userClientManager) AddClient(cli *Client) {
	if uc, ok := getUserClient(cli.Uid); ok {
		uc.clients[cli.UUID] = cli
	} else {
		ucm.AddUserClient(cli)
	}
}

// down observer mode	|	locked from cliManager ->func addClient
func (ucm *userClientManager) RemoveClient(cli *Client) {
	delete(ucm.ucs[cli.Uid].clients, cli.UUID)
	if len(ucm.ucs[cli.Uid].clients) == 0 {
		//user offline
		changeUserStatusChan <- struct {
			uid    int
			status bool
		}{uid: cli.Uid, status: false}
	}
}

func (ucm *userClientManager) userStatusCustomer() {
	for c := range changeUserStatusChan {
		if c.status {
			getOnlineUsers()[c.uid] = time.Now()
		} else {
			delete(getOnlineUsers(), c.uid)
		}
	}
}
