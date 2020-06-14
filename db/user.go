package db

import (
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"github.com/timshannon/bolthold"
)

type User struct {
	ID            snowflake.ID `boltholdKey:"id"`
	Username      string
	Discriminator disgord.Discriminator
	bank          map[Resource]int
}

func (u *User) GetResourceCount(resource Resource) int {
	if u.bank == nil {
		return 0
	}

	return u.bank[resource]
}

func (u *User) AddResource(resource Resource, count int) {
	if u.bank == nil {
		u.bank = map[Resource]int{}
	}

	u.bank[resource] = u.bank[resource] + count
}

func (u *User) RemoveResource(resource Resource, count int) {
	if u.bank == nil {
		u.bank = map[Resource]int{}
		return
	}

	banked := u.bank[resource]
	banked -= count
	if banked < 0 {
		banked = 0
	}

	u.bank[resource] = banked
}

type UserRepository interface {
	Add(user *User) error
	Remove(userId snowflake.ID) error
	Get(userId snowflake.ID) (*User, error)
	Update(user *User) error
	GetAll() ([]*User, error)
}

type userRepository struct {
	store *bolthold.Store
}

func newUserRepository(store *bolthold.Store) UserRepository {
	return &userRepository{store: store}
}

func (r *userRepository) Add(user *User) error {
	return r.store.Insert(user.ID, user)
}

func (r *userRepository) Remove(userId snowflake.ID) error {
	return r.store.Delete(userId, User{})
}

func (r *userRepository) Get(userId snowflake.ID) (*User, error) {
	var resource *User
	err := r.store.Find(resource, bolthold.Where("id").Eq(userId))
	return resource, err
}
func (r *userRepository) Update(user *User) error {
	return r.store.Update(user.ID, user)
}

func (r *userRepository) GetAll() ([]*User, error) {
	users := make([]*User, 0)
	err := r.store.Find(&users, nil)
	return users, err
}
