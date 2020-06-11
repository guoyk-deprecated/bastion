package models

import (
	"github.com/jinzhu/copier"
	"github.com/guoyk93/bastion/types"
)

type Session struct {
	Id         int64  `storm:"id,increment"`
	Account    string `storm:"index"`
	Command    string
	CreatedAt  int64
	FinishedAt int64
	IsRecorded bool
}

func (s Session) ToGRPCSession() *types.Session {
	n := types.Session{}
	copier.Copy(&n, &s)
	return &n
}
