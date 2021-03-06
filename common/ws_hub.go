package common

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// TODO: Rename this to WebSockets?
var WsHub WsHubImpl

// TODO: Make this an interface?
// TODO: Write tests for this
type WsHubImpl struct {
	// TODO: Implement some form of generics so we don't write as much odd-even sharding code
	evenOnlineUsers map[int]*WSUser
	oddOnlineUsers  map[int]*WSUser
	evenUserLock    sync.RWMutex
	oddUserLock     sync.RWMutex

	// TODO: Add sharding for this too?
	OnlineGuests map[*WSUser]bool
	GuestLock    sync.RWMutex

	lastTick      time.Time
	lastTopicList []*TopicsRow
}

func init() {
	// TODO: Do we really want to initialise this here instead of in main.go / general_test.go like the other things?
	WsHub = WsHubImpl{
		evenOnlineUsers: make(map[int]*WSUser),
		oddOnlineUsers:  make(map[int]*WSUser),
		OnlineGuests:    make(map[*WSUser]bool),
	}
}

func (hub *WsHubImpl) Start() {
	ticker := time.NewTicker(time.Minute * 5)
	defer func() {
		ticker.Stop()
	}()

	go func() {
		for {
			var item = func(lock *sync.RWMutex, userMap map[int]*WSUser) {
				lock.RLock()
				defer lock.RUnlock()
				// TODO: Copy to temporary slice for less contention?
				for _, user := range userMap {
					user.Ping()
				}
			}
			select {
			case <-ticker.C:
				item(&hub.evenUserLock, hub.evenOnlineUsers)
				item(&hub.oddUserLock, hub.oddOnlineUsers)
			}
		}
	}()
	if Config.DisableLiveTopicList {
		return
	}
	hub.lastTick = time.Now()
	AddScheduledSecondTask(hub.Tick)
}

// This Tick is separate from the admin one, as we want to process that in parallel with this due to the blocking calls to gopsutil
func (hub *WsHubImpl) Tick() error {
	return wsTopicListTick(hub)
}

func wsTopicListTick(hub *WsHubImpl) error {
	// Avoid hitting GetList when the topic list hasn't changed
	if !TopicListThaw.Thawed() && hub.lastTopicList != nil {
		return nil
	}
	
	// Don't waste CPU time if nothing has happened
	// TODO: Get a topic list method which strips stickies?
	tList, _, _, err := TopicList.GetList(1, "", nil)
	if err != nil {
		hub.lastTick = time.Now()
		return err // TODO: Do we get ErrNoRows here?
	}
	defer func() {
		hub.lastTick = time.Now()
		hub.lastTopicList = tList
	}()
	if len(tList) == 0 {
		return nil
	}

	// TODO: Optimise this by only sniffing the top non-sticky
	// TODO: Optimise this by getting back an unsorted list so we don't have to hop around the stickies
	// TODO: Add support for new stickies / replies to them
	if len(tList) == len(hub.lastTopicList) {
		var hasItem = false
		for j, tItem := range tList {
			if !tItem.Sticky {
				if tItem.ID != hub.lastTopicList[j].ID || !tItem.LastReplyAt.Equal(hub.lastTopicList[j].LastReplyAt) {
					hasItem = true
				}
			}
		}
		if !hasItem {
			return nil
		}
	}

	// TODO: Implement this for guests too? Should be able to optimise it far better there due to them sharing the same permission set
	// TODO: Be less aggressive with the locking, maybe use an array of sorts instead of hitting the main map every-time
	topicListMutex.RLock()
	if len(topicListWatchers) == 0 {
		topicListMutex.RUnlock()
		return nil
	}

	// Copy these over so we close this loop as fast as possible so we can release the read lock, especially if the group gets are backed by calls to the database
	var groupIDs = make(map[int]bool)
	var currentWatchers = make([]*WSUser, len(topicListWatchers))
	var i = 0
	for wsUser, _ := range topicListWatchers {
		currentWatchers[i] = wsUser
		groupIDs[wsUser.User.Group] = true
		i++
	}
	topicListMutex.RUnlock()

	var groups = make(map[int]*Group)
	var canSeeMap = make(map[string][]int)
	for groupID, _ := range groupIDs {
		group, err := Groups.Get(groupID)
		if err != nil {
			// TODO: Do we really want to halt all pushes for what is possibly just one user?
			return err
		}
		groups[group.ID] = group

		var canSee = make([]byte, len(group.CanSee))
		for i, item := range group.CanSee {
			canSee[i] = byte(item)
		}
		canSeeMap[string(canSee)] = group.CanSee
	}

	var canSeeRenders = make(map[string][]byte)
	for name, canSee := range canSeeMap {
		topicList, forumList, _, err := TopicList.GetListByCanSee(canSee, 1, "", nil)
		if err != nil {
			return err // TODO: Do we get ErrNoRows here?
		}
		if len(topicList) == 0 {
			continue
		}
		_ = forumList // Might use this later after we get the base feature working

		if topicList[0].Sticky {
			var lastSticky = 0
			for i, row := range topicList {
				if !row.Sticky {
					lastSticky = i
					break
				}
			}
			if lastSticky == 0 {
				continue
			}
			topicList = topicList[lastSticky:]
		}

		// TODO: Compare to previous tick to eliminate unnecessary work and data
		var wsTopicList = make([]*WsTopicsRow, len(topicList))
		for i, topicRow := range topicList {
			wsTopicList[i] = topicRow.WebSockets()
		}

		outBytes, err := json.Marshal(&WsTopicList{wsTopicList, 0})
		if err != nil {
			return err
		}
		canSeeRenders[name] = outBytes
	}

	// TODO: Use MessagePack for additional speed?
	//fmt.Println("writing to the clients")
	for _, wsUser := range currentWatchers {
		group := groups[wsUser.User.Group]
		var canSee = make([]byte, len(group.CanSee))
		for i, item := range group.CanSee {
			canSee[i] = byte(item)
		}

		//fmt.Println("writing to user #", wsUser.User.ID)
		outBytes := canSeeRenders[string(canSee)]
		//fmt.Println("outBytes: ", string(outBytes))
		err := wsUser.WriteToPageBytes(outBytes, "/topics/")
		if err == ErrNoneOnPage {
			//fmt.Printf("werr for #%d: %s\n", wsUser.User.ID, err)
			wsUser.FinalizePage("/topics/", func() {
				topicListMutex.Lock()
				delete(topicListWatchers, wsUser)
				topicListMutex.Unlock()
			})
			continue
		}
	}
	return nil
}

func (hub *WsHubImpl) GuestCount() int {
	defer hub.GuestLock.RUnlock()
	hub.GuestLock.RLock()
	return len(hub.OnlineGuests)
}

func (hub *WsHubImpl) UserCount() (count int) {
	hub.evenUserLock.RLock()
	count += len(hub.evenOnlineUsers)
	hub.evenUserLock.RUnlock()
	hub.oddUserLock.RLock()
	count += len(hub.oddOnlineUsers)
	hub.oddUserLock.RUnlock()
	return count
}

func (hub *WsHubImpl) HasUser(uid int) (exists bool) {
	hub.evenUserLock.RLock()
	_, exists = hub.evenOnlineUsers[uid]
	hub.evenUserLock.RUnlock()
	if exists {
		return exists
	}
	hub.oddUserLock.RLock()
	_, exists = hub.oddOnlineUsers[uid]
	hub.oddUserLock.RUnlock()
	return exists
}

func (hub *WsHubImpl) broadcastMessage(msg string) error {
	var userLoop = func(users map[int]*WSUser, mutex *sync.RWMutex) error {
		defer mutex.RUnlock()
		for _, wsUser := range users {
			err := wsUser.WriteAll(msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// TODO: Can we move this RLock inside the closure safely?
	hub.evenUserLock.RLock()
	err := userLoop(hub.evenOnlineUsers, &hub.evenUserLock)
	if err != nil {
		return err
	}
	hub.oddUserLock.RLock()
	return userLoop(hub.oddOnlineUsers, &hub.oddUserLock)
}

func (hub *WsHubImpl) getUser(uid int) (wsUser *WSUser, err error) {
	var ok bool
	if uid%2 == 0 {
		hub.evenUserLock.RLock()
		wsUser, ok = hub.evenOnlineUsers[uid]
		hub.evenUserLock.RUnlock()
	} else {
		hub.oddUserLock.RLock()
		wsUser, ok = hub.oddOnlineUsers[uid]
		hub.oddUserLock.RUnlock()
	}
	if !ok {
		return nil, errWsNouser
	}
	return wsUser, nil
}

// Warning: For efficiency, some of the *WSUsers may be nil pointers, DO NOT EXPORT
func (hub *WsHubImpl) getUsers(uids []int) (wsUsers []*WSUser, err error) {
	if len(uids) == 0 {
		return nil, errWsNouser
	}
	var appender = func(lock *sync.RWMutex, users map[int]*WSUser) {
		lock.RLock()
		defer lock.RUnlock()
		// We don't want to keep a lock on this for too long, so we'll accept some nil pointers
		for _, uid := range uids {
			wsUsers = append(wsUsers, users[uid])
		}
	}
	appender(&hub.evenUserLock, hub.evenOnlineUsers)
	appender(&hub.oddUserLock, hub.oddOnlineUsers)
	if len(wsUsers) == 0 {
		return nil, errWsNouser
	}
	return wsUsers, nil
}

// For Widget WOL, please avoid using this as it might wind up being really long and slow without the right safeguards
func (hub *WsHubImpl) AllUsers() (users []*User) {
	var appender = func(lock *sync.RWMutex, userMap map[int]*WSUser) {
		lock.RLock()
		defer lock.RUnlock()
		for _, user := range userMap {
			users = append(users, user.User)
		}
	}
	appender(&hub.evenUserLock, hub.evenOnlineUsers)
	appender(&hub.oddUserLock, hub.oddOnlineUsers)
	return users
}

func (hub *WsHubImpl) removeUser(uid int) {
	if uid%2 == 0 {
		hub.evenUserLock.Lock()
		delete(hub.evenOnlineUsers, uid)
		hub.evenUserLock.Unlock()
	} else {
		hub.oddUserLock.Lock()
		delete(hub.oddOnlineUsers, uid)
		hub.oddUserLock.Unlock()
	}
}

func (hub *WsHubImpl) AddConn(user User, conn *websocket.Conn) (*WSUser, error) {
	if user.ID == 0 {
		wsUser := new(WSUser)
		wsUser.User = new(User)
		*wsUser.User = user
		wsUser.AddSocket(conn, "")
		WsHub.GuestLock.Lock()
		WsHub.OnlineGuests[wsUser] = true
		WsHub.GuestLock.Unlock()
		return wsUser, nil
	}

	// TODO: How should we handle user state changes if we're holding a pointer which never changes?
	userptr, err := Users.Get(user.ID)
	if err != nil && err != ErrStoreCapacityOverflow {
		return nil, err
	}

	var mutex *sync.RWMutex
	var theMap map[int]*WSUser
	if user.ID%2 == 0 {
		mutex = &hub.evenUserLock
		theMap = hub.evenOnlineUsers
	} else {
		mutex = &hub.oddUserLock
		theMap = hub.oddOnlineUsers
	}

	mutex.Lock()
	wsUser, ok := theMap[user.ID]
	if !ok {
		wsUser = new(WSUser)
		wsUser.User = userptr
		wsUser.Sockets = []*WSUserSocket{&WSUserSocket{conn, ""}}
		theMap[user.ID] = wsUser
		mutex.Unlock()
		return wsUser, nil
	}
	mutex.Unlock()
	wsUser.AddSocket(conn, "")
	return wsUser, nil
}

func (hub *WsHubImpl) RemoveConn(wsUser *WSUser, conn *websocket.Conn) {
	wsUser.RemoveSocket(conn)
	wsUser.Lock()
	if len(wsUser.Sockets) == 0 {
		hub.removeUser(wsUser.User.ID)
	}
	wsUser.Unlock()
}

func (hub *WsHubImpl) PushMessage(targetUser int, msg string) error {
	wsUser, err := hub.getUser(targetUser)
	if err != nil {
		return err
	}
	return wsUser.WriteAll(msg)
}

func (hub *WsHubImpl) pushAlert(targetUser int, alert Alert) error {
	wsUser, err := hub.getUser(targetUser)
	if err != nil {
		return err
	}
	astr, err := BuildAlert(alert, *wsUser.User)
	if err != nil {
		return err
	}
	return wsUser.WriteAll(astr)
}

func (hub *WsHubImpl) pushAlerts(users []int, alert Alert) error {
	wsUsers, err := hub.getUsers(users)
	if err != nil {
		return err
	}

	var errs []error
	for _, wsUser := range wsUsers {
		if wsUser == nil {
			continue
		}
		alert, err := BuildAlert(alert, *wsUser.User)
		if err != nil {
			errs = append(errs, err)
		}
		err = wsUser.WriteAll(alert)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Return the first error
	if len(errs) != 0 {
		for _, err := range errs {
			return err
		}
	}
	return nil
}
