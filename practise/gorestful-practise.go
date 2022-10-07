package main

// repo: https://github.com/emicklei/go-restful
// Refer to https://www.kubernetes.org.cn/1788.html
// https://github.com/emicklei/go-restful/tree/master/examples
import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

type Userg struct {
	Id, Name, Age string
}

type UserResource struct {
	// normally one would use DAO (data access object)
	users map[string]Userg
}

type RequestBody struct {
	Name string
	Age  string
}

func (u UserResource) Register(container *restful.Container) {
	// 创建新的WebService
	ws := new(restful.WebService)

	// 设定WebService对应的路径("/users")和支持的MIME类型(restful.MIME_XML/ restful.MIME_JSON)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	// 添加路由： GET /{user-id} --> u.findUser
	// 带括号的表示变量，不带的是固定路径，如果不写 则404
	//ws.Route(ws.GET("/{user-id}/test/{user-name}").To(u.findUser))
	ws.Route(ws.GET("/{user-id}").To(u.findUser))

	// 添加路由： POST / --> u.updateUser
	ws.Route(ws.POST("").To(u.createUser))

	// 添加路由： PUT /{user-id} --> u.createUser
	ws.Route(ws.PUT("/{user-id}").To(u.updateUser))

	// 添加路由： DELETE /{user-id} --> u.removeUser
	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser))

	// 将初始化好的WebService添加到Container中
	container.Add(ws)
}

// GET http://127.0.0.1:8080/users/test?name=caoyingjun
// GET http://127.0.0.1:8080/users/id2/test/user-name?k1=111&k2=222
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	// 获取user-id 和 url 中的变量方法
	// 获取指定 path 的路径变量
	userId := request.PathParameter("user-id")

	// 获取 k v 变量
	query := request.Request.URL.Query()
	name := query.Get("name")

	usr := u.users[userId]
	if len(usr.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(name)
	}
}

//POST http://127.0.0.1:8080/users
//Content-Type: application/json
//
//{
//"name": "caoyingjun",
//"age": "18"
//}
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := Userg{
		Id: request.PathParameter("user-id"),
	}

	err := request.ReadEntity(&usr)
	if err != nil {
		response.WriteEntity(usr)
		return
	}

	u.users[usr.Id] = usr
	// superfluous  多余的设置
	//response.WriteHeader(200)
	response.WriteEntity(usr)
}

//PUT http://127.0.0.1:8080/users/test
//Content-Type: application/json
//
//{
//"name": "caoyingjun",
//"age": "18"
//}
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(Userg)
	userId := request.PathParameter("user-id")
	// 将请求中的body获取
	err := request.ReadEntity(&usr)
	if err != nil {
		response.WriteEntity(usr)
		return
	}

	id := usr.Id
	name := usr.Name
	age := usr.Age
	fmt.Println(userId, id, name, age)
	u.users[usr.Id] = *usr
	response.WriteEntity(usr)
}

//DELETE  http://127.0.0.1:8080/users/test
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
	// 设置返回的格式
	response.WriteEntity(map[string]string{"delete": "true"})
}

func main() {
	// 创建一个空的Container
	wsContainer := restful.NewContainer()

	// 设定路由为CurlyRouter(快速路由)
	wsContainer.Router(restful.CurlyRouter{})

	// 创建自定义的Resource Handle(此处为UserResource)
	u := UserResource{map[string]Userg{}}

	// 创建WebService，并将WebService加入到Container中
	u.Register(wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}

	// 启动服务
	log.Fatal(server.ListenAndServe())
}
