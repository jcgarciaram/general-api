package general_api

import (
    "net/url"
    
    "github.com/tmaiaroto/aegis/lambda"
)



func cors(inner lambda.RouteHandler) lambda.RouteHandler {
    return lambda.RouteHandler(func(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {

        res.Headers = map[string]string{ 
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE, HEAD",
            "Access-Control-Allow-Credentials": "true",
            "Access-Control-Max-Age": "86400",
            "Access-Control-Allow-Headers":
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-XSRF-Token, X-HTTP-Method-Override, X-Requested-With",
        }
        
        inner(ctx, evt, res, params)

    })
}