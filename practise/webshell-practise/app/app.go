package app

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/igm/sockjs-go/v3/sockjs"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

type WebShell struct {
	Conn     sockjs.Session
	SizeChan chan *remotecommand.TerminalSize

	Namespace string
	Pod       string
	Container string
}

func (w *WebShell) Write(p []byte) (int, error) {
	err := w.Conn.Send(string(p))
	return len(p), err
}

func (w *WebShell) Read(p []byte) (int, error) {
	var msg map[string]uint16
	reply, err := w.Conn.Recv()
	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal([]byte(reply), &msg); err != nil {
		return copy(p, reply), nil
	} else {
		w.SizeChan <- &remotecommand.TerminalSize{
			Height: msg["rows"],
			Width:  msg["cols"],
		}
		return 0, nil
	}
}

func (w *WebShell) Next() *remotecommand.TerminalSize {
	size := <-w.SizeChan
	return size
}

func WebShellHandler(w *WebShell, cmd string) error {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return err
	}

	gv := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config.GroupVersion = &gv
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = scheme.Codecs
	clientSet, err := rest.RESTClientFor(config)
	if err != nil {
		return err
	}
	req := clientSet.Post().
		Resource("pods").
		Name(w.Pod).
		Namespace(w.Namespace).
		SubResource("exec").
		Param("container", w.Container).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("command", cmd).Param("tty", "true")
	req.VersionedParams(
		&v1.PodExecOptions{
			Container: w.Container,
			Command:   []string{},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		scheme.ParameterCodec,
	)
	executor, err := remotecommand.NewSPDYExecutor(
		config, http.MethodPost, req.URL(),
	)
	if err != nil {
		return err
	}
	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             w,
		Stdout:            w,
		Stderr:            w,
		Tty:               true,
		TerminalSizeQueue: w,
	})
}
