package general_api

import (
    // "fmt"
    "net/url"
    
    "github.com/tmaiaroto/aegis/lambda"
    r "github.com/jcgarciaram/general-api/routes"
    
)

func NewRouter(routes r.Routes) *lambda.Router {
    

    // Handle with a URL reqeust path Router
	router := lambda.NewRouter(fallThrough)

    for _, route := range routes {
        
        handler := route.HandlerFunc
        handler = cors(handler)
        handler = logger(handler, route.Name)


        router.Handle(route.Method, route.Pattern, handler)

        
    }
    return router
}




func fallThrough(ctx *lambda.Context, evt *lambda.Event, res *lambda.ProxyResponse, params url.Values) {
	res.SetStatus(404)
}
