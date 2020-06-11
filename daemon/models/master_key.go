package models

import (
	"github.com/jinzhu/copier"
	"github.com/guoyk93/bastion/types"
)

type MasterKey struct {
	Fingerprint string `storm:"id"`
	PublicKey   string
}

func (m MasterKey) ToGRPCModel() *types.MasterKey {
	o := types.MasterKey{}
	copier.Copy(&o, &m)
	return &o
}
