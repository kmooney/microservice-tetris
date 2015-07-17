package main

import (
    "net/url"
    "net/http"
    "restlite"
    "fmt"
    "strconv"
    "errors"
    "encoding/json"
    "bytes"
    "io"
)

const WIDTH uint = 10
const HEIGHT uint = 20

/** Shape **/
type Shape struct { 
    Width int
    Data []bool
    Position [2]int // row, column
}

type ShapeResource struct { 
    restlite.DeleteNotSupported
    restlite.GetNotSupported
    Gameboards Gameboards
    Subscriptions Subscriptions
}

func GetParam(values url.Values, p string) ([]byte, error) { 
    param := values[p]
    if len(param) == 0 {
        return []byte{0}, errors.New("Parameter must exist")
    }

    return []byte(param[0]), nil
}

func (sr ShapeResource) Post(values url.Values, body io.Reader) (int, interface{}) {
    var game_id int
    var err error
    var param []byte
    var shape Shape 
    var valid bool
    
    param, err = GetParam(values, "game_id")

    if err != nil { 
        return 500, err.Error()
    }

    game_id, err = strconv.Atoi(string(param))
    
    if err != nil {
        return 500, fmt.Sprintf("Param is not int %v", err.Error())
    }
    
    decoder := json.NewDecoder(body)
    err = decoder.Decode(&shape)
    if err != nil {
        return 500, fmt.Sprintf("Invalid shapedata: %s", err.Error())
    }
    valid, err = sr.Gameboards[game_id].Valid(shape)
    if err != nil { 
        return 500, err.Error()
    }
    if valid {
        return 200, "OK"
    } else {
        return 412, "Precondition Failed - Precondition was shape is valid for board."
    }
}

// put the shape in the position
func (sr ShapeResource) Put(values url.Values, body io.Reader) (int, interface{}) { 
    var game_id int
    var err error
    var param []byte
    var shape Shape 
    var valid bool

    param, err = GetParam(values, "game_id")
    if err != nil { 
        return 500, err.Error()
    }

    game_id, err = strconv.Atoi(string(param))
    if err != nil { 
        return 500, err.Error()
    }

    decoder := json.NewDecoder(body)
    err = decoder.Decode(&shape)
    if err != nil { 
        return 500, fmt.Sprintf("Invalid shapedata: %s", err.Error())
    }
    valid, err = sr.Gameboards[game_id].Valid(shape)
    if err != nil { 
        return 500, err.Error()
    }
    if valid {
        // TODO check validity, as in GET, then place shape and trigger callback
        sr.Gameboards[game_id].Place(shape)
        return 200, fmt.Sprintf("OK %d, %v", game_id, shape)
    } else {
        return 412, "PUT Failed - Precondition was shape is valid for board."
    }
}

/** Ticks **/
type TickResource struct { 
    restlite.PutNotSupported
    restlite.DeleteNotSupported
    restlite.PostNotSupported
    Gameboards Gameboards
    Subscriptions Subscriptions
}


func (tr TickResource) Get(values url.Values, body io.Reader) (int, interface{}) {
    var game_id int
    var err error
    var changed bool
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
    changed, err = tr.Gameboards[game_id].Tick()
    if err != nil { 
        return 500, err.Error()
    }
    if changed { 
        fmt.Println("Ok, it changed - I'm calling notify on the subscriptions")
        go tr.Subscriptions.Notify(game_id, tr.Gameboards[game_id])
    }
    return 200, fmt.Sprintf("OK - Game # %d ticked", game_id)
}

type Shapedata_t [][]bool

/** Gameboard **/
type Gameboard struct {
    Level int
    Lines int
    Score int
    Shapedata Shapedata_t
    Gameover bool
}

func (gb *Gameboard) Tick () (bool, error) { 
    /**
       1.  Check to see if any rows are completed.
           a)  If yes, line count should be increased
           b)  If yes, score should be increased
           c)  If yes, completed lines should be cleared

       2. Has anything changed (gameover, lines, level, score, shapes)?
           a)  Notify callback listeners.
    **/
    var changed bool
    var complete_row bool
    var row_count int
    var row []bool
    var row_idx int

    for row_idx = 0; row_idx < 20; row_idx ++ {
        row = gb.Shapedata[row_idx]
        complete_row = true

        for _, cell := range row {
            if cell == false {
                complete_row = false
            }
        }

        if complete_row {
            gb.Shapedata = append(gb.Shapedata[:row_idx], gb.Shapedata[(row_idx+1):]...)
            gb.Shapedata = append([][]bool{{false,false,false,false,false,false,false,false,false,false}},gb.Shapedata...)
            row_count ++
        }
    }

    if row_count > 0 { 
        if row_count == 4 {
            gb.Score += 2500
        }
        if row_count > 1 {
            gb.Score += row_count * 500
        } else { 
            gb.Score += 100
        }
        changed = true
        gb.Lines += row_count
    }

    return changed, nil
}

func (gb Gameboard) Valid(s Shape) (bool, error) {
    /**
        Make sure the shape and its position fall within
        the bounds of the gameboard and do not collide with any
        of the placed-shape data.
    **/
    var row int
    var col int
    var k int
    var l = len(s.Data)
    for k = 0; k < l; k++ {
        col = (k % s.Width) + s.Position[1]
        row = (k / s.Width) + s.Position[0]
        if row >= 20 || row < 0 || col >= 10 || col < 0 {
            return false, nil
        }
        if gb.Shapedata[row][col] == true && s.Data[k] == true {
            return false, nil 
        }
        
    }
    return true, nil
}

func (gb Gameboard) Place(s Shape) (error) {
    var valid bool 
    var err error
    var row int
    var col int
    var k int
    var l = len(s.Data)

    valid, err = gb.Valid(s)
    if err != nil { 
        return err
    }
    if !valid {
        return errors.New("Bad board data - shape doesn't fit!")
    }

    for k = 0; k < l; k++ {
        row = (k / s.Width) + s.Position[0]
        col = (k % s.Width) + s.Position[1] 
        if (s.Data[k]) {
            gb.Shapedata[row][col] = true
        }
    }
    return nil
}

type Gameboards map[int]*Gameboard

type GameboardResource struct {
    restlite.PutNotSupported
    restlite.DeleteNotSupported
    Gameboards Gameboards
}

func (gr GameboardResource) Get(values url.Values, body io.Reader) (int, interface{}) { 
    game_id_data, err := GetParam(values, "game_id")
    game_id_string := string(game_id_data)
    var game_id int
    game_id, err  = strconv.Atoi(game_id_string)

    if err != nil {
        return 500, "Can't convert Game ID"
    }
    gameboard, exists := gr.Gameboards[game_id]
    if exists == false {
        return 404, "Game ID not found"
    }
    return 200, gameboard
}

func (gr GameboardResource) Post(values url.Values, body io.Reader) (int, interface{}) { 
    game_id_data, err := GetParam(values, "game_id")
    game_id_string := string(game_id_data)
    var game_id int
    game_id, err  = strconv.Atoi(game_id_string)
    if err != nil { 
        return 500, "Game ID must be int value"
    }
    _, check := gr.Gameboards[game_id]
    if check != false { 
        return 500, "Game is already created, cannot replace"
    }
    shapeData := make([][]bool, 20)
    shapeValues := make([]bool, 10*20)
    for i := range shapeData {
        shapeData[i], shapeValues = shapeValues[:10], shapeValues[10:]
    }
    gr.Gameboards[game_id] = &Gameboard{1, 0, 0, shapeData, false}
    return 200, "OK"
}

/** Subscriptions **/

type Subscriptions map[int]map[string]bool

func (s Subscriptions) Notify(game_id int, g *Gameboard) {
    //var err error
    fmt.Printf("Called notify! %v", s[game_id])
    
    for url := range(s[game_id]) {
        var buf bytes.Buffer
        encoder := json.NewEncoder(&buf)
        encoder.Encode(g)
        go http.Post(url, "application/octet-stream", &buf)
    }
} 

type SubscriptionResource struct {
    Subscriptions Subscriptions
    restlite.PutNotSupported
}

func (sr SubscriptionResource) Get(values url.Values, body io.Reader) (int, interface{}) {
    var game_id int
    game_id_param := values["game_id"]
    if len(game_id_param) == 1 {
        game_id, _ = strconv.Atoi(game_id_param[0])
        return 200, sr.Subscriptions[game_id]
    }
    return 500, fmt.Sprintf("Need to define game_id parameter")
}

func (sr SubscriptionResource) Post(values url.Values, body io.Reader) (int, interface{}) {
    var game_id int
    var response_url string
    var err error

    game_id_param := values["game_id"]
    response_url_param := values["response_url"]

    if len(game_id_param) != 1 {
        return 500, "Need to define game_id"
    } else {
        game_id, err = strconv.Atoi(game_id_param[0])
        if (err != nil) {
            return 500, fmt.Sprintf("Game ID must be an integer: %v", err)
        }
    }
    if len(response_url_param) != 1 {
        return 500, "Need to define response_url"
    } else { 
        response_url = response_url_param[0]
    }
    _, exists := sr.Subscriptions[game_id]
    if !exists { 
        sr.Subscriptions[game_id] = make(map[string]bool)
    }

    sr.Subscriptions[game_id][response_url]=true
    return 200, "Posted"
}


func (sr SubscriptionResource) Delete(values url.Values, body io.Reader) (int, interface{}) { 
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
    subscriptions := make(Subscriptions)

    gameboardResource := new(GameboardResource)
    gameboardResource.Gameboards = gameboards

    tickResource := new(TickResource)
    tickResource.Gameboards = gameboards
    tickResource.Subscriptions = subscriptions

    subscriptionResource := new(SubscriptionResource)
    subscriptionResource.Subscriptions = subscriptions

    shapeResource := new(ShapeResource)
    shapeResource.Gameboards = gameboards
    shapeResource.Subscriptions = subscriptions

    var api = new (restlite.API)

    api.AddResource(tickResource, "/tick")
    api.AddResource(gameboardResource, "/game")
    api.AddResource(subscriptionResource, "/subscribe")
    api.AddResource(shapeResource, "/shape")

    api.Start(8000)
}
