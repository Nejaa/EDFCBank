package db

import (
	"EDFCBank/params"
	"EDFCBank/utils"
	"github.com/timshannon/bolthold"
)

var st *bolthold.Store

var Resources ResourceRepository
var Users UserRepository

func init() {
	var err error
	st, err = bolthold.Open(*params.DbFile, 0666, nil)
	utils.FatalOnError(err)

	Resources = newResourceRepository(st)
	Users = newUserRepository(st)

}
