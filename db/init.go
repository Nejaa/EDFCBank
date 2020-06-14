package db

import (
	"EDFCBank/utils"
	"github.com/timshannon/bolthold"
)

var st *bolthold.Store

var Resources ResourceRepository

func init() {
	var err error
	st, err = bolthold.Open("test.db", 0666, nil)
	utils.FatalOnError(err)

	Resources = newResourceRepository(st)

}