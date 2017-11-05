/*
*
* Reply Resources File
* Copyright Azareal 2016 - 2018
*
 */
package main

import (
	"errors"
	"time"
)

// ? - Should we add a reply store to centralise all the reply logic? Would this cover profile replies too or would that be separate?
var rstore ReplyStore
var prstore ProfileReplyStore

type ReplyUser struct {
	ID                int
	ParentID          int
	Content           string
	ContentHtml       string
	CreatedBy         int
	UserLink          string
	CreatedByName     string
	Group             int
	CreatedAt         time.Time
	RelativeCreatedAt string
	LastEdit          int
	LastEditBy        int
	Avatar            string
	ClassName         string
	ContentLines      int
	Tag               string
	URL               string
	URLPrefix         string
	URLName           string
	Level             int
	IPAddress         string
	Liked             bool
	LikeCount         int
	ActionType        string
	ActionIcon        string
}

type Reply struct {
	ID                int
	ParentID          int
	Content           string
	CreatedBy         int
	Group             int
	CreatedAt         time.Time
	RelativeCreatedAt string
	LastEdit          int
	LastEditBy        int
	ContentLines      int
	IPAddress         string
	Liked             bool
	LikeCount         int
}

var ErrAlreadyLiked = errors.New("You already liked this!")

// TODO: Write tests for this
// TODO: Wrap these queries in a transaction to make sure the state is consistent
func (reply *Reply) Like(uid int) (err error) {
	var rid int // unused, just here to avoid mutating reply.ID
	err = stmts.hasLikedReply.QueryRow(uid, reply.ID).Scan(&rid)
	if err != nil && err != ErrNoRows {
		return err
	} else if err != ErrNoRows {
		return ErrAlreadyLiked
	}

	score := 1
	_, err = stmts.createLike.Exec(score, reply.ID, "replies", uid)
	if err != nil {
		return err
	}
	_, err = stmts.addLikesToReply.Exec(1, reply.ID)
	return err
}

// TODO: Write tests for this
func (reply *Reply) Delete() error {
	_, err := stmts.deleteReply.Exec(reply.ID)
	if err != nil {
		return err
	}
	_, err = stmts.removeRepliesFromTopic.Exec(1, reply.ParentID)
	tcache, ok := topics.(TopicCache)
	if ok {
		tcache.CacheRemove(reply.ParentID)
	}
	return err
}

// Copy gives you a non-pointer concurrency safe copy of the reply
func (reply *Reply) Copy() Reply {
	return *reply
}

// TODO: Refactor this to stop hitting the global stmt store
type ReplyStore interface {
	Get(id int) (*Reply, error)
	Create(tid int, content string, ipaddress string, fid int, uid int) (id int, err error)
}

type SQLReplyStore struct {
}

func NewSQLReplyStore() *SQLReplyStore {
	return &SQLReplyStore{}
}

func (store *SQLReplyStore) Get(id int) (*Reply, error) {
	reply := Reply{ID: id}
	err := stmts.getReply.QueryRow(id).Scan(&reply.ParentID, &reply.Content, &reply.CreatedBy, &reply.CreatedAt, &reply.LastEdit, &reply.LastEditBy, &reply.IPAddress, &reply.LikeCount)
	return &reply, err
}

// TODO: Write a test for this
func (store *SQLReplyStore) Create(tid int, content string, ipaddress string, fid int, uid int) (id int, err error) {
	wcount := wordCount(content)
	res, err := stmts.createReply.Exec(tid, content, parseMessage(content, fid, "forums"), ipaddress, wcount, uid)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = stmts.addRepliesToTopic.Exec(1, uid, tid)
	if err != nil {
		return int(lastID), err
	}
	tcache, ok := topics.(TopicCache)
	if ok {
		tcache.CacheRemove(tid)
	}
	return int(lastID), err
}

type ProfileReplyStore interface {
	Get(id int) (*Reply, error)
}

// TODO: Refactor this to stop using the global stmt store
type SQLProfileReplyStore struct {
}

func NewSQLProfileReplyStore() *SQLProfileReplyStore {
	return &SQLProfileReplyStore{}
}

func (store *SQLProfileReplyStore) Get(id int) (*Reply, error) {
	reply := Reply{ID: id}
	err := stmts.getUserReply.QueryRow(id).Scan(&reply.ParentID, &reply.Content, &reply.CreatedBy, &reply.CreatedAt, &reply.LastEdit, &reply.LastEditBy, &reply.IPAddress)
	return &reply, err
}
