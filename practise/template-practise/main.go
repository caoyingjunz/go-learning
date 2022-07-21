package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	pt "go-learning/practise/template-practise/template"
)

// 内置一些用于比较的函数：
// eq arg1 arg2：
// arg1 == arg2时为true
// ne arg1 arg2：
// arg1 != arg2时为true
// lt arg1 arg2：
// arg1 < arg2时为true
// le arg1 arg2：
// arg1 <= arg2时为true
// gt arg1 arg2：
// arg1 > arg2时为true
// ge arg1 arg2：
// arg1 >= arg2时为true
// 更多可以 参考 https://www.cnblogs.com/f-ck-need-u/p/10053124.html
// 渲染 + 通用模板合并

const (
	// YAMLDocumentSeparator is the separator for YAML documents
	YAMLDocumentSeparator = "---\n"
)

type Friend struct {
	Name string `json:"name"`
}

type Person struct {
	UserName string            `json:"user_name"`
	Emails   []string          `json:"emails"`
	Friends  []Friend          `json:"friends"`
	Mods     map[string]string `json:"mods"`
}

func MergeBytes(mBytes ...[]byte) []byte {
	return bytes.Join(mBytes, []byte(YAMLDocumentSeparator))
}

func main() {
	fri := Friend{
		Name: "name",
	}

	emails := make([]string, 0)
	emails = append(emails, "test1@gmail.com")
	emails = append(emails, "test2@gmail.com")

	tpl := template.New("test")
	tpl = template.Must(tpl.Parse(pt.ServiceTemplate))

	m := make(map[string]string)
	m["m1"] = "v1"
	m["m2"] = "v2"

	p := Person{
		UserName: "caoyingjun",
		Emails:   emails,
		Friends:  []Friend{fri},
		Mods:     m,
	}

	f, err := os.Create("./service.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 写入指定文件
	tpl.Execute(f, p)

	// 写入标准输出
	//tpl.Execute(os.Stdout, p)

	// 写入变量1
	var buf bytes.Buffer
	tpl.Execute(&buf, p)
	fmt.Println(string(buf.Bytes()))

	// 写入变量2
	//bs := strings.Builder{}
	//tpl.Execute(&bs, p)
	//fmt.Println(bs.String())

	b, err := ioutil.ReadFile("./service.yaml")
	if err != nil {
		panic(err)
	}

	nb := MergeBytes([]byte(pt.CommonService), b)
	if err = ioutil.WriteFile("./merge-service.yaml", nb, 0640); err != nil {
		panic(err)
	}
}
