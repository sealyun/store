package serve

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
)

// Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on a web page
// to make XMLHttpRequests to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
//
// GET http://localhost:8080/users
//
// GET http://localhost:8080/users/1
//
// PUT http://localhost:8080/users/1
//
// DELETE http://localhost:8080/users/1
//
// OPTIONS http://localhost:8080/users/1  with Header "Origin" set to some domain and

//UserResource  s
type UserResource struct{}

const (
	clientID     = "89c1b05d77fb1c92a1ef"
	clientSecret = "541ddd76e65abeabd12ad9f8b02f6601394d3ad0"
)

//RegisterTo is
func (u UserResource) RegisterTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("").
		Consumes("*/*").
		Produces("*/*")

	ws.Filter(webserviceLogging).Filter(checkCookie)

	ws.Route(ws.GET("/{user-id}").To(u.nop))
	ws.Route(ws.GET("/callback").To(u.callback))
	ws.Route(ws.POST("").To(u.nop))
	ws.Route(ws.PUT("/{user-id}").To(u.nop))
	ws.Route(ws.DELETE("/{user-id}").To(u.nop))

	container.Add(ws)
}

// if check cookie failed, redirect to login page
func checkCookie(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()
	chain.ProcessFilter(req, resp)
	log.Printf("[webservice-filter (timer)] %v\n", time.Now().Sub(now))
}

func (u UserResource) nop(request *restful.Request, response *restful.Response) {
	io.WriteString(response.ResponseWriter, "this would be a normal response")
}

func (u UserResource) callback(request *restful.Request, response *restful.Response) {
	code := request.QueryParameter("code")
	accessToken, err := GetGithubAccessToken(clientID, clientSecret, code)
	if err != nil {
		io.WriteString(response.ResponseWriter, "fetch token failed"+accessToken)
	}
	user, err := GetUserInfo(accessToken)
	if err != nil {
		fmt.Println(err)
	}
	_ = user
	io.WriteString(response.ResponseWriter, "code is : "+code)
}

//Run is
func Run() {
	wsContainer := restful.NewContainer()
	u := UserResource{}
	u.RegisterTo(wsContainer)

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	log.Print("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
