package endpoint

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"go-learning/practise/gin-practise/hander"
	"go-learning/practise/gin-practise/k8s"
	"go-learning/practise/gin-practise/log"
	"go-learning/practise/gin-practise/worker"
)

type GinResp struct {
	Code    int         `json:"code"`
	Resp    interface{} `json:"resp,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (g *GinResp) SetCode(c int) {
	g.Code = c
}
func (g *GinResp) SetMessage(msg string) {
	g.Message = msg
}

type practiseDB struct {
	db interface{}
	s  string
}

func (c *practiseDB) getPractise(db interface{}) *practiseDB {
	c.db = db
	return c
}

func (c *practiseDB) get() string {
	return c.s
}

var pdb = practiseDB{s: "ceshi db orm"}

func GetPractise(c *gin.Context) {
	r := GinResp{}

	desc := pdb.getPractise("db-driver").get()
	log.Glog.Info(desc)

	log.Glog.Info("get GetPractise log")

	r.SetMessage("Get Practise Message")
	r.SetCode(200)

	c.JSON(200, r)
}

func PostMid(c *gin.Context) {
	fmt.Println("start mid test")
	n := struct {
		Name string `json:"name"`
		Age  string `json:"age"`
	}{}

	if err := c.ShouldBindJSON(&n); err != nil {
		c.AbortWithStatusJSON(400, "bad request")
		return
	}

	c.JSON(200, "ok")
}

func PostPractise(c *gin.Context) {
	r := GinResp{}

	var gr hander.GinRequest
	var err error
	err = c.ShouldBindJSON(&gr)
	if err != nil {
		c.AbortWithStatusJSON(500, "bad request")
		return
	}

	if err = hander.Dohandler(context.TODO(), gr); err != nil {
		r.SetMessage(err.Error())
		c.AbortWithStatusJSON(500, "bad request")
		return
	}

	r.SetMessage("Get Practise Message")
	r.SetCode(200)

	log.Glog.Info("Post Practise log")

	r.Resp = gr
	c.JSON(200, r)
}

var (
	WorkerSet  worker.WorkerInterface
	KubeEngine k8s.EngineInterface
)

func RegisterWorker(w worker.WorkerInterface) {
	WorkerSet = w
}

func RegisterEngine(e k8s.EngineInterface) {
	KubeEngine = e
}

func TestQueue(c *gin.Context) {
	r := GinResp{}

	q := c.Query("queue")
	if err := WorkerSet.DoTest(context.TODO(), q); err != nil {
		r.SetMessage("test error queue")
		c.JSON(400, r)
		return
	}

	r.Resp = "test ok queue"
	c.JSON(200, r)
}

func TestAfterQueue(c *gin.Context) {
	r := GinResp{}

	q := c.Query("after")
	if err := WorkerSet.DoAfterTest(context.TODO(), q); err != nil {
		r.SetMessage("test error after queue")
		c.JSON(400, r)
		return
	}

	r.Resp = "test ok after queue"
	c.JSON(200, r)
}

func TestPod(c *gin.Context) {
	r := GinResp{}

	name := c.Query("name")
	namespace := c.Query("namespace")
	key := c.Query("key")

	pod, err := KubeEngine.GetPod(context.TODO(), key, name, namespace)
	if err != nil {
		r.SetMessage(err.Error())
		c.JSON(400, r)
		return
	}

	r.Resp = pod
	c.JSON(200, r)
}

func TestOptimise(c *gin.Context) {
	r := GinResp{}
	// 在 Create 里面进行参数的解析和异常处理
	if err := Practise.Optimise(c).Create("dd"); err != nil {
		return
	}

	c.JSON(200, r)
}

func GetOptimise(c *gin.Context) {
	r := GinResp{}
	var err error
	// 在 get 里设置参数
	if r.Resp, err = Practise.Optimise(c).Get(); err != nil {
		return
	}
}

func Download(c *gin.Context) {
	filename := c.Query("filename")

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename+".tar")
	c.File("/Users/xxx/yyy.tar")
}
