package serve

import (
	"fmt"
	"io"
	"log"
	"net/http"

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

//RegisterTo is
func (u UserResource) RegisterTo(container *restful.Container) {
	loginless := new(restful.WebService)
	loginless.
		Path("").
		Consumes("*/*").
		Produces("*/*")
	loginless.Route(loginless.GET("/callback").To(u.callback))

	ws := new(restful.WebService)
	ws.
		Path("/user").
		Consumes("*/*").
		Produces("*/*")

	ws.Filter(checkCookie)

	ws.Route(ws.GET("/{user-id}").To(u.nop))
	ws.Route(ws.POST("").To(u.nop))
	ws.Route(ws.PUT("/{user-id}").To(u.nop))
	ws.Route(ws.DELETE("/{user-id}").To(u.nop))

	container.Add(ws)
	container.Add(loginless)
}

// if check cookie failed, redirect to login page
func checkCookie(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	cookie, err := req.Request.Cookie("user")
	if err != nil || cookie == nil {
		fmt.Println("login please : ", err, req.Request.URL.String())
		http.Redirect(resp, req.Request, GetLoginURL(req.Request.URL.String()), http.StatusMovedPermanently)
		return
	}
	chain.ProcessFilter(req, resp)
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

	// Set cookie
	cookie := http.Cookie{Name: "user", Value: user.Login, Path: "/", MaxAge: 86400}
	http.SetCookie(response, &cookie)

	request.Request.AddCookie(&cookie)

	state := request.QueryParameter("state")
	fmt.Println("redirect url is : ", state)
	//redirect back to user request
	http.Redirect(response, request.Request, "http://localhost:8001"+state, http.StatusMovedPermanently)

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
