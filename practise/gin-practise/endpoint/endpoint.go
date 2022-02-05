package endpoint

import (
	"context"

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

func GetPractise(c *gin.Context) {
	r := GinResp{}

	log.Glog.Info("get GetPractise log")

	r.SetMessage("Get Practise Message")
	r.SetCode(200)

	c.JSON(200, r)
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
	WorkerSet  = worker.NewWorker()
	KubeEngine = k8s.NewKubeEngine()
)

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
