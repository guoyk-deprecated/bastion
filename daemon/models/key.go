package models

import (
	"github.com/jinzhu/copier"
	"github.com/guoyk93/bastion/types"
)

type Key struct {
	Fingerprint string `storm:"id"`
	Account     string `storm:"index"`
	Source      string `storm:"index"`
	Name        string
	CreatedAt   int64
	ViewedAt    int64
}

func (k Key) ToGRPCKey() *types.Key {
	n := types.Key{}
	copier.Copy(&n, &k)
	return &n
}
