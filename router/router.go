package router

import (
	"EDFCBank/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"strconv"
	"strings"
)

// EventContext provides the handler with the details of the received message as well as the occasional path variables
type EventContext struct {
	Session        disgord.Session
	Event          *disgord.MessageCreate
	pathVariables  map[string]interface{}
	deleteOnAnswer bool
}

// StringParam returns a path variable formatted as a string
func (ec *EventContext) StringParam(paramName string) (string, bool) {
	if val, found := ec.pathVariables[paramName]; found {
		return val.(string), true
	}
	return "", false
}

// IntParam returns a path variable formatted as an int or -1 if failed
func (ec *EventContext) IntParam(paramName string) (int, bool) {
	if val, found := ec.pathVariables[paramName]; found {
		if n, err := strconv.Atoi(val.(string)); err == nil {
			return n, true
		}
	}
	return -1, false
}

// Answer provides responds to the incoming command in the same channel.
// In case the router is defined to delete messages on answer the command message is removed from the channel
func (ec *EventContext) Answer(message string) {
	if ec.deleteOnAnswer {
		err := ec.Session.DeleteMessage(context.Background(),
			ec.Event.Message.ChannelID,
			ec.Event.Message.ID)
		utils.LogOnError(err)
	}
	_, err := ec.Session.CreateMessage(
		context.Background(),
		ec.Event.Message.ChannelID,
		&disgord.CreateMessageParams{Content: message})
	utils.LogOnError(err)
}

// RouterHandler is the signature of the router's handlers
type RouterHandler = func(ec *EventContext)

type routeNode struct {
	identifier string
	childs     map[string]*routeNode
	isParam    bool
	handler    RouterHandler
}

func newRouteNode(identifier string) *routeNode {
	isParam := strings.HasPrefix(identifier, ":")
	if isParam {
		identifier = strings.TrimPrefix(identifier, ":")
	}

	return &routeNode{
		identifier: identifier,
		childs:     map[string]*routeNode{},
		isParam:    isParam,
	}
}

func (rn *routeNode) getKey() string {
	if rn.isParam {
		return "*"
	}

	return rn.identifier
}

// RouterHandler receives MessageCreate events and dispatch them as required
type Router interface {
	RegisterPath(path string, handler RouterHandler)
	DeleteOnAnswer(deleteOnAnswer bool)
}

type router struct {
	RoutePrefix    string
	deleteOnAnswer bool
	routes         map[string]*routeNode
}

func NewRouter(client *disgord.Client, prefix string) Router {
	r := &router{
		RoutePrefix:    prefix,
		routes:         map[string]*routeNode{},
		deleteOnAnswer: true,
	}

	filter, err := std.NewMsgFilter(context.Background(), client)
	utils.FatalOnError(err)

	filter.SetPrefix(prefix)
	client.On(disgord.EvtMessageCreate, filter.HasPrefix, r.handleMessage)

	return r
}

func (r *router) DeleteOnAnswer(deleteOnAnswer bool) {
	r.deleteOnAnswer = deleteOnAnswer

}
func (r *router) RegisterPath(path string, handler RouterHandler) {
	path = strings.TrimSpace(path)
	pathElems := strings.Split(path, " ")
	pathLen := len(pathElems)

	curPaths := r.routes
	for idx, pathElem := range pathElems {
		lastElem := idx == pathLen-1

		node := newRouteNode(pathElem)

		// try to find if path elem already registered
		if n, found := curPaths[node.getKey()]; found {
			sameIdentifier := node.identifier == n.identifier
			if node.isParam && !sameIdentifier {
				panic("cannot register two path variables with different names at the same level " + node.identifier + "/" + n.identifier)
			} else if lastElem && sameIdentifier { // same route registered twice
				panic("route " + path + "already registered")
			}
			node = n
		} else { // else create new branch
			curPaths[node.getKey()] = node
		}

		if lastElem {
			node.handler = handler
		} else {
			curPaths = node.childs
		}
	}
}

func (r *router) handleMessage(session disgord.Session, evt *disgord.MessageCreate) {
	ec := &EventContext{
		Session:        session,
		Event:          evt,
		pathVariables:  map[string]interface{}{},
		deleteOnAnswer: r.deleteOnAnswer,
	}

	defer func() {
		if err := recover(); err != nil {
			ec.Answer(fmt.Sprintf("An error as occured : %v", err))
		}
	}()

	content := strings.TrimSpace(evt.Message.Content)
	pathElems := strings.Split(content, " ")
	pathLen := len(pathElems)
	if pathLen == 0 || pathLen == 1 && pathElems[0] == r.RoutePrefix {
		ec.Answer("Empty command")
	}

	if r.RoutePrefix != "" {
		pathElems = pathElems[1:]
		pathLen--
	}

	var handled bool
	curPaths := r.routes
	for idx, pathElem := range pathElems {
		lastElem := idx == pathLen-1

		var node *routeNode
		var found bool
		if node, found = curPaths[pathElem]; found { // try to find route based on the actual value of the path element
		} else if node, found = curPaths["*"]; found { // check if there is a pathvariable option registered for this level
			ec.pathVariables[node.identifier] = pathElem
		}

		if !found {
			break
		}

		if lastElem {
			handled = true
			node.handler(ec)
		} else {
			curPaths = node.childs
		}
	}

	if !handled {
		ec.Answer("Unrecognized command")
	}
}
