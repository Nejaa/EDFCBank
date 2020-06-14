package db

import (
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

type User struct {
	ID            snowflake.ID `boltholdKey:"name"`
	Username      string
	discriminator disgord.Discriminator
	Bank          map[Resource]int64
}
