package db

import (
	"github.com/timshannon/bolthold"
)

type Resource struct {
	Name string `boltholdKey:"Name"`
}

type ResourceRepository interface {
	Add(resource *Resource) error
	Remove(resource *Resource) error
	Get(name string) (*Resource, error)
	GetAll() ([]*Resource, error)
	Count() (int, error)
}

type resourceRepository struct {
	store *bolthold.Store
}

func newResourceRepository(store *bolthold.Store) ResourceRepository {
	return &resourceRepository{store: store}
}

func (r *resourceRepository) Add(resource *Resource) error {
	return r.store.Insert(resource.Name, resource)
}

func (r *resourceRepository) Remove(resource *Resource) error {
	return r.store.Delete(resource.Name, resource)
}

func (r *resourceRepository) Get(name string) (*Resource, error) {
	resource := &Resource{}
	err := r.store.FindOne(resource, bolthold.Where("Name").Eq(name))
	return resource, err
}

func (r *resourceRepository) GetAll() ([]*Resource, error) {
	resources := make([]*Resource, 0)
	err := r.store.Find(&resources, nil)
	return resources, err
}

func (r *resourceRepository) Count() (int, error) {
	return r.store.Count(&Resource{}, nil)
}
