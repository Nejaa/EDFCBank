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
	//r.RegisterPath("bank remove :resourceName :count", onBankRemove)
	//r.RegisterPath("bank list", onBankList)
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

//
//func onBankRemove(ec *router.EventContext) {
//	resourceName, found := ec.StringParam("resourceName")
//	if !found {
//		ec.Answer("No resource name given")
//	}
//
//	err := db.Resources.Remove(&db.Resource{Name: resourceName})
//	if err != bolthold.ErrKeyExists {
//		utils.PanicOnError(err)
//	}
//
//	ec.Answer("Ressource " + resourceName + " not managed anymore")
//}
//
//func onBankList(ec *router.EventContext) {
//	if count, err := db.Resources.Count(); err != nil {
//		ec.Answer("No resource created")
//	} else {
//		message := fmt.Sprintf("%d known resources:\n", count)
//		elems, err := db.Resources.GetAll()
//		utils.PanicOnError(err)
//		users, err := db.Users.GetAll()
//		utils.PanicOnError(err)
//		resourceMap := map[db.Resource]int{}
//		for _, elem := range elems {
//			resourceMap[*elem] = 0
//			for _, user := range users {
//				resourceMap[*elem] += user.GetResourceCount(*elem)
//			}
//		}
//
//		for resource, count := range resourceMap {
//			message += fmt.Sprintf("- %s : %d \n", resource.Name, count)
//		}
//
//		ec.Answer(message)
//	}
//}
