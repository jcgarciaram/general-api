package routes

import (
    "net/url"
    
    "github.com/tmaiaroto/aegis/lambda"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc lambda.RouteHandler
}

type Routes []Route


func (routes *Routes) AppendRoutes(newRoutes Routes) {
    for _, r := range newRoutes {
        *routes = append(*routes, r)
    }
}


func OptionsHandler() lambda.RouteHandler {
    return lambda.RouteHandler(func(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

        
        res.Headers = map[string]string{
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, HEAD",
            "Access-Control-Allow-Credentials": "true",
            "Access-Control-Max-Age": "86400",
            "Access-Control-Allow-Headers":
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-XSRF-Token, X-HTTP-Method-Override, X-Requested-With",
        }

  
    })
}
