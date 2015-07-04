package main

import (
    "net/url"
    "restlite"
    "fmt"
    "strconv"
)

/** Ticks **/

type TickResource struct { 
    restlite.PutNotSupported
    restlite.DeleteNotSupported
    restlite.PostNotSupported
}

func (TickResource) Get(values url.Values) (int, interface{}) { 
    return 200, fmt.Sprintf("OK - Game # %v", values["game_id"])
}

/** Gameboard **/

type Gameboard struct { 
    Level int
    Lines int
    Score int
    Shapedata map[int]map[int]bool
}

type Gameboards []Gameboard

type GameboardResource struct { 
    restlite.PutNotSupported
    restlite.DeleteNotSupported
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

func (SubscriptionResource) Post(values url.Values) (int, interface{}) { 
    return 400, "Not Defined Yet"
}

func (SubscriptionResource) Put(values url.Values) (int, interface{}) { 
    return 400, "Not Defined Yet"
}

func (SubscriptionResource) Delete(values url.Values) (int, interface{}) { 
    return 400, "Not Defined Yet"
}

func main() {
    
    gameboardResource := new(GameboardResource)
    tickResource := new(TickResource)
    subscriptionResource := new(SubscriptionResource)
    subscriptionResource.Subscriptions = make(Subscriptions)

    var api = new (restlite.API)

    api.AddResource(tickResource, "/tick")
    api.AddResource(gameboardResource, "/game")
    api.AddResource(subscriptionResource, "/subscribe")

    api.Start(8000)
}
