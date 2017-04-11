package general_api

import (
    "log"
    "time"
    "net/url"
    "strings"
    
    "github.com/tmaiaroto/aegis/lambda"
)



func logger(inner lambda.RouteHandler, name string) lambda.RouteHandler {
    return lambda.RouteHandler(func(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
        start := time.Now()

        inner(ctx, evt, res, params)
        
        
        index       := strings.IndexRune(evt.Path, '/')
        requestURI  := evt.Path[index:]
        
        
        log.Printf(
            "%s\t%s\t%s\t%s",
            evt.HTTPMethod,
            requestURI,
            name,
            time.Since(start),
        )
    })
}