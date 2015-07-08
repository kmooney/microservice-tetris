package restlite

import (
    "net/http"
    "net/url"
    "fmt"
    "encoding/json"
    "io"

)

type Resource interface { 
    Get(values url.Values, body io.Reader) (int, interface{})
    Post(values url.Values, body io.Reader) (int, interface{})
    Put(values url.Values, body io.Reader) (int, interface{})
    Delete(values url.Values, body io.Reader) (int, interface{})
}

type GetNotSupported struct {}
func (GetNotSupported) Get(values url.Values, body io.Reader) (int, interface{}) { 
    return 405, ""
}

type PostNotSupported struct {}
func (PostNotSupported) Post(values url.Values, body io.Reader) (int, interface{}) {
    return 405, ""
}

type PutNotSupported struct {}
func (PutNotSupported) Put(values url.Values, body io.Reader) (int, interface{}) { 
    return 405, ""
}

type DeleteNotSupported struct {}
func (DeleteNotSupported) Delete(values url.Values, body io.Reader) (int, interface{}) { 
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
        var err error

        request.ParseForm()
        fmt.Println(request)
        method := request.Method
        values := request.Form


        // TODO if values are empty, and method is post, 
        // use the body of the post.
        switch method { 
            case "GET":
                code, data = resource.Get(values, request.Body)
            case "POST":
                code, data = resource.Post(values, request.Body)
            case "PUT":
                code, data = resource.Put(values, request.Body)
            case "DELETE":
                code, data = resource.Delete(values, request.Body)
            default:
                api.Abort(rw, 405)
        } 
        content, err := json.Marshal(data)
        if err != nil { 
            fmt.Println("Error: ", err.Error())
            api.Abort(rw, 500)
        }
        rw.WriteHeader(code)
        rw.Write(content)
    }
}
