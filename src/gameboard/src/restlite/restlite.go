package restlite

import (
    "net/http"
    "net/url"
    "fmt"
    "encoding/json"
)

type Resource interface { 
    Get(values url.Values) (int, interface{})
    Post(values url.Values) (int, interface{})
    Put(values url.Values) (int, interface{})
    Delete(values url.Values) (int, interface{})
}

type GetNotSupported struct {}
func (GetNotSupported) Get(values url.Values) (int, interface{}) { 
    return 405, ""
}

type PostNotSupported struct {}
func (PostNotSupported) Post(values url.Values) (int, interface{}) {
    return 405, ""
}

type PutNotSupported struct {}
func (PutNotSupported) Put(values url.Values) (int, interface{}) { 
    return 405, ""
}

type DeleteNotSupported struct {}
func (DeleteNotSupported) Delete(values url.Values) (int, interface{}) { 
    return 405, ""
}

type API struct {}

func (api *API) AddResource(resource Resource, path string) {
        http.HandleFunc(path, api.requestHandler(resource))
}

func (api *API) Start(port int) { 
    portString := fmt.Sprintf(":%d", port)
    http.ListenAndServe(portString, nil)
}

func (api *API) Abort(rw http.ResponseWriter, statusCode int) { 
    rw.WriteHeader(statusCode)
}

func (api *API) requestHandler(resource Resource) http.HandlerFunc { 
    return func(rw http.ResponseWriter, request *http.Request) { 
        var data interface{}
        var code int

        request.ParseForm()
        method := request.Method
        values := request.Form
        switch method { 
            case "GET":
                code, data = resource.Get(values)
            case "POST":
                code, data = resource.Post(values)
            case "PUT":
                code, data = resource.Post(values)
            case "DELETE":
                code, data = resource.Delete(values)
            default:
                api.Abort(rw, 405)
        } 
        content, err := json.Marshal(data)
        if err != nil { 
            api.Abort(rw, 500)
        }
        rw.WriteHeader(code)
        rw.Write(content)
    }
}