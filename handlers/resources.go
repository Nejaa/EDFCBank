package handlers

import (
	"EDFCBank/db"
	"EDFCBank/router"
	"EDFCBank/utils"
	"fmt"
	"github.com/timshannon/bolthold"
)

func RegisterResourceRoutes(r router.Router) {
	r.RegisterPath("resources add :resourceName", onResourceAdd)
	r.RegisterPath("resources remove :resourceName", onResourceRemove)
	r.RegisterPath("resources list ", onResourceList)
}

func onResourceList(ec *router.EventContext) {

	if count, err := db.Resources.Count(); err != nil {
		ec.Answer("No resource created")
	} else {
		message := fmt.Sprintf("%d known resources:\n", count)
		elems, err := db.Resources.GetAll()
		utils.PanicOnError(err)
		for _, elem := range elems {
			message = message + "- " + elem.Name + "\n"
		}

		ec.Answer(message)

	}
}

func onResourceAdd(ec *router.EventContext) {
	resourceName, found := ec.StringParam("resourceName")
	if !found {
		ec.Answer("No resource name given")
	}

	err := db.Resources.Add(&db.Resource{Name: resourceName})
	if err == bolthold.ErrKeyExists {
		ec.Answer("Resource " + resourceName + " already exists")
	} else if err != nil {
		utils.PanicOnError(err)
	} else {
		ec.Answer("Resource " + resourceName + " now managed")
	}
}

func onResourceRemove(ec *router.EventContext) {
	resourceName, found := ec.StringParam("resourceName")
	if !found {
		ec.Answer("No resource name given")
	}

	err := db.Resources.Remove(&db.Resource{Name: resourceName})
	if err != bolthold.ErrKeyExists {
		utils.PanicOnError(err)
	}

	ec.Answer("Ressource " + resourceName + " not managed anymore")
}
