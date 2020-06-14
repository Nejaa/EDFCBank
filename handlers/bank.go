package handlers

import (
	"EDFCBank/db"
	"EDFCBank/router"
	"EDFCBank/utils"
	"fmt"
	"github.com/timshannon/bolthold"
)

func RegisterBankRoutes(r router.Router) {
	r.RegisterPath("bank add :resourceName :count", onBankAdd)
	r.RegisterPath("bank remove :resourceName :count", onBankRemove)
	r.RegisterPath("bank balance", onBankBalance)
}

func onBankAdd(ec *router.EventContext) {
	user, err := db.Users.Get(ec.Event.Message.Author)
	utils.PanicOnError(err)

	resourceName, found := ec.StringParam("resourceName")
	if !found {
		ec.Answer("No resource name given")
	}

	count, found := ec.IntParam("count")
	if !found {
		ec.Answer("No resource count given")
	}

	resource, err := db.Resources.Get(resourceName)
	if err == bolthold.ErrNotFound {
		ec.Answer("Resource " + resourceName + " is not managed")
	} else if err != nil {
		utils.PanicOnError(err)
	} else {
		user.AddResource(*resource, count)
		err := db.Users.Update(user)
		utils.PanicOnError(err)
		ec.Answer(fmt.Sprintf("<@%d> you know own %d %s", user.Id, user.GetResourceCount(*resource), resource.Name))
	}
}

func onBankRemove(ec *router.EventContext) {
	user, err := db.Users.Get(ec.Event.Message.Author)
	utils.PanicOnError(err)

	resourceName, found := ec.StringParam("resourceName")
	if !found {
		ec.Answer("No resource name given")
	}

	count, found := ec.IntParam("count")
	if !found {
		ec.Answer("No resource count given")
	}

	resource, err := db.Resources.Get(resourceName)
	if err == bolthold.ErrNotFound {
		resource = &db.Resource{Name: resourceName} // allow to remove elements from an un-managed resource for cleanup
	} else if err != nil {
		utils.PanicOnError(err)
	}

	user.RemoveResource(*resource, count)
	err = db.Users.Update(user)
	utils.PanicOnError(err)
	ec.Answer(fmt.Sprintf("<@%d> you now own %d %s", user.Id, user.GetResourceCount(*resource), resource.Name))
}

func onBankBalance(ec *router.EventContext) {

	user, err := db.Users.Get(ec.Event.Message.Author)
	utils.PanicOnError(err)

	resourceSet := map[db.Resource]struct{}{}
	resources, err := db.Resources.GetAll()
	for _, resource := range resources {
		resourceSet[*resource] = struct{}{}
	}

	message := fmt.Sprintf("<@%d> your balance is:\n", user.Id)
	if len(user.Bank) == 0 {
		message += "Empty"
	} else {
		for resource, count := range user.Bank {
			managedStatus := "not managed"
			if _, found := resourceSet[resource]; found {
				managedStatus = "managed"
			}
			message += fmt.Sprintf("\t- %s : %d (%s)\n", resource.Name, count, managedStatus)
		}
	}
	ec.Answer(message)
}
