package main

import (
    "net/url"
    "restlite"
    "fmt"
    "strconv"
    "errors"
    "encoding/json"
)

const WIDTH uint = 10
const HEIGHT uint = 20

// TODO: Create validation and shape placement resources.


/** Shape **/
type Shape struct { 
    Width int
    Data []bool
    Position [2]int
}

type ShapeResource struct { 
    restlite.DeleteNotSupported
    restlite.PostNotSupported
    Gameboards Gameboards
}

func GetParam(values url.Values, p string) ([]byte, error) { 
    param := values[p]
    if len(param) == 0 {
        return []byte{0}, errors.New("Parameter must exist")
    }

    return []byte(param[0]), nil
}

func (vr ShapeResource) Get(values url.Values) (int, interface{}) {
    var game_id int
    var err error
    var param []byte
    var shape Shape 

    param, err = GetParam(values, "game_id")
    if err != nil { 
        return 500, err.Error()
    }
    game_id, err = strconv.Atoi(string(param))
    if err != nil { 
        return 500, err.Error()
    }

    param, err = GetParam(values, "shapedata")
    if err != nil { 
        return 500, err.Error()
    }

    err = json.Unmarshal(param, shape)
    if err != nil { 
        return 500, fmt.Sprintf("Invalid shapedata: %s", err.Error())
    }
    // check validity, else return error
    return 200, fmt.Sprintf("OK %i, %v", game_id, shape)
}

// put the shape in the position 
func (sr ShapeResource) Put(values url.Values) (int, interface{}) { 
    var game_id int
    var err error
    var param []byte
    var shape Shape 

    param, err = GetParam(values, "game_id")
    if err != nil { 
        return 500, err.Error()
    }
    game_id, err = strconv.Atoi(string(param))
    if err != nil { 
        return 500, err.Error()
    }

    param, err = GetParam(values, "shapedata")
    if err != nil { 
        return 500, err.Error()
    }

    err = json.Unmarshal(param, shape)
    if err != nil { 
        return 500, fmt.Sprintf("Invalid shapedata: %s", err.Error())
    }
    // check validity, as in GET, then place shape and trigger callback
    return 200, fmt.Sprintf("OK %i, %v", game_id, shape)
}


/** Ticks **/

type TickResource struct { 
    restlite.PutNotSupported
    restlite.DeleteNotSupported
    restlite.PostNotSupported
    Gameboards Gameboards
}

/** THINK: What's the semantically correct HTTP verb to perform a tick - probably
    not a GET, why would you "GET" a Tick resource? **/
func (tr TickResource) Get(values url.Values) (int, interface{}) {
    var game_id int
    var err error
    game_id_param := values["game_id"]
    if len(game_id_param) != 1 { 
        return 500, "Bad game_id parameter"
    } else { 
        game_id, err = strconv.Atoi(game_id_param[0])
        if err != nil { 
            return 500, "Cannot convert game_id to int"
        }
    }
    // maybe this is a goroutine?  we can return quick from this
    // and rely on the callback to notify clients....
    tr.Gameboards[game_id].Tick()
    return 200, fmt.Sprintf("OK - Game # %i ticked", game_id)
}


/** Gameboard **/
type Gameboard struct {
    Level int
    Lines int
    Score int
    Shapedata map[int]map[int]bool
    CurrentShape Shape
    Gameover bool
}

func (gb Gameboard) Tick () (error) { 
    /**
       1.  Check to see if any rows are completed.
           a)  If yes, line count should be increased
           b)  If yes, score should be increased
           c)  If yes, completed lines should be cleared

       2.  Should the game be over?

       3.  Should the game be deleted?  (how to tell?)
           a)  If yes, delete the game from the map
           b)  Notify callback listeners.

       4. Has anything changed (gameover, lines, level, score, shapes)?
           a)  Notify callback listeners.
    **/
    return nil
}

func (gb Gameboard) ValidShapePosition(s Shape) (bool) {
    /**
        Make sure the shape and its position fall within
        the bounds of the gameboard and do not collide with any
        of the placed-shape data.
    **/
    return false
}

type Gameboards map[int]Gameboard

type GameboardResource struct {
    restlite.PutNotSupported
    restlite.DeleteNotSupported
    Gameboards Gameboards
}

func (GameboardResource) Get(values url.Values) (int, interface{}) { 
    data := map[string]string{"game": "board"}
    return 200, data
}

func (GameboardResource) Post(values url.Values) (int, interface{}) { 
    data := map[string]string{"you_posted": "a new gameboard"}
    return 200, data
}

/** Subscriptions **/

type Subscription struct { 
    GameID int
    Property string
    ResponseMethod string
    ResponseUrl string
}

// TODO Subscriptions should be indexed by GameID and ResponseUrl, not just GameID
// this way, we can have multiple observers on the same GameID
type Subscriptions map[int]Subscription

type SubscriptionResource struct {
    Subscriptions Subscriptions
}

func (sr SubscriptionResource) Get(values url.Values) (int, interface{}) { 
    var game_id int
    game_id_param := values["game_id"]
    if len(game_id_param) == 1 { 
        game_id, _ = strconv.Atoi(game_id_param[0])
        return 400, sr.Subscriptions[game_id]
    }
    return 500, fmt.Sprintf("Need to define game_id parameter")
}

func (sr SubscriptionResource) Post(values url.Values) (int, interface{}) { 
    var game_id int
    var response_method string
    var property string
    var response_url string
    var err error

    game_id_param := values["game_id"]
    property_param := values["property"]
    response_method_param := values["response_method"]
    response_url_param := values["response_url"]

    if len(game_id_param) != 1 {
        return 500, "Need to define game_id"
    } else { 
        game_id, err = strconv.Atoi(game_id_param[0])
        if (err != nil) { 
            return 500, fmt.Sprintf("Game ID must be an integer: %v", err)
        }
    }
    if len(response_method_param) != 1 { 
        response_method = "POST"
    } else { 
        response_method = response_method_param[0]
    }
    if len(response_url_param) != 1 { 
        return 500, "Need to define response_url"
    } else { 
        response_url = response_url_param[0]
    }
    if len(property_param) == 1 {
        property = property_param[0]
    }
    _, exists := sr.Subscriptions[game_id]
    if exists { 
        return 500, "Subscription already set - use PUT if you want to change it"
    }
    sr.Subscriptions[game_id] = Subscription{game_id, property, response_method, response_url} 
    return 200, "Posted"
}

func (sr SubscriptionResource) Put(values url.Values) (int, interface{}) { 
    var game_id int
    var response_method string
    var property string
    var response_url string
    var err error

    game_id_param := values["game_id"]
    property_param := values["property"]
    response_method_param := values["response_method"]
    response_url_param := values["response_url"]

    if len(game_id_param) != 1 {
        return 500, "Need to define game_id"
    } else { 
        game_id, err = strconv.Atoi(game_id_param[0])
        if (err != nil) { 
            return 500, fmt.Sprintf("Game ID must be an integer: %v", err)
        }
    }
    if len(response_method_param) != 1 { 
        response_method = "POST"
    } else { 
        response_method = response_method_param[0]
    }
    if len(response_url_param) != 1 { 
        return 500, "Need to define response_url"
    } else { 
        response_url = response_url_param[0]
    }
    if len(property_param) == 1 {
        property = property_param[0]
    }
    _, exists := sr.Subscriptions[game_id]
    if !exists { 
        return 404, "Subscription does not exist, POST to create one."
    }
    sr.Subscriptions[game_id] = Subscription{game_id, property, response_method, response_url} 
    return 200, "Put"
}

func (sr SubscriptionResource) Delete(values url.Values) (int, interface{}) { 
    var game_id int
    var err error
    game_id_param := values["game_id"]
    if len(game_id_param) == 1 { 
        game_id, err = strconv.Atoi(game_id_param[0])
        if err != nil { 
            return 500, fmt.Sprintf("%v", err)
        }
    } else { 
        return 500, "Game ID must be set exactly once"
    }

    _, exists := sr.Subscriptions[game_id]
    if (exists) { 
        delete(sr.Subscriptions, game_id)
        return 200, "Deleted"
    }
    return 404, "No such subscription"
}

func main() {
    gameboards := make(Gameboards)    

    gameboardResource := new(GameboardResource)
    gameboardResource.Gameboards = gameboards

    tickResource := new(TickResource)
    tickResource.Gameboards = gameboards

    subscriptionResource := new(SubscriptionResource)
    subscriptionResource.Subscriptions = make(Subscriptions)

    var api = new (restlite.API)

    api.AddResource(tickResource, "/tick")
    api.AddResource(gameboardResource, "/game")
    api.AddResource(subscriptionResource, "/subscribe")

    api.Start(8000)
}
