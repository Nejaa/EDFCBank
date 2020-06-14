package db

import (
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"github.com/timshannon/bolthold"
)

type User struct {
	Id            snowflake.ID `boltholdKey:"Id"`
	Username      string
	Discriminator disgord.Discriminator
	Bank          map[Resource]int
}

func newUserFromAuthor(author *disgord.User) *User {
	return &User{
		Id:            author.ID,
		Username:      author.Username,
		Discriminator: author.Discriminator,
		Bank:          map[Resource]int{},
	}
}

func (u *User) GetResourceCount(resource Resource) int {
	if u.Bank == nil {
		return 0
	}

	return u.Bank[resource]
}

func (u *User) AddResource(resource Resource, count int) {
	if count < 0 {
		panic("trying to add a negative amount")
	}
	if u.Bank == nil {
		u.Bank = map[Resource]int{}
	}

	u.Bank[resource] += count
}

func (u *User) RemoveResource(resource Resource, count int) {
	if count < 0 {
		panic("trying to remove a negative amount")
	}

	if u.Bank == nil {
		u.Bank = map[Resource]int{}
		return
	}

	banked := u.Bank[resource]
	banked -= count
	if banked < 0 {
		banked = 0
	}

	if banked == 0 {
		banked = 0
		delete(u.Bank, resource)
	} else {
		u.Bank[resource] = banked
	}
}

type UserRepository interface {
	Add(user *User) error
	Remove(userId snowflake.ID) error
	Get(author *disgord.User) (*User, error)
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
	return r.store.Insert(user.Id, user)
}

func (r *userRepository) Remove(userId snowflake.ID) error {
	return r.store.Delete(userId, User{})
}

func (r *userRepository) Get(author *disgord.User) (*User, error) {
	user := &User{}
	err := r.store.FindOne(user, bolthold.Where("Id").Eq(author.ID))
	if err != nil && err == bolthold.ErrNotFound {
		user = newUserFromAuthor(author)
		if err = r.Add(user); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return user, err
}
func (r *userRepository) Update(user *User) error {
	return r.store.Update(user.Id, user)
}

func (r *userRepository) GetAll() ([]*User, error) {
	users := make([]*User, 0)
	err := r.store.Find(&users, nil)
	return users, err
}
