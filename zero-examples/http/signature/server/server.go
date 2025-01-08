package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/zero-examples/http/signature/internal"
)

type Request struct {
	User string `form:"user,optional"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := httpx.Parse(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	io.Copy(w, r.Body)
}

func main() {
	flag.Parse()
	priFile, err := createTempFile(internal.PriKey)
	if err != nil {
		panic(err)
	}
	c := rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Verbose: true,
		Port:    3333,
		Signature: rest.SignatureConf{
			Strict: true,
			PrivateKeys: []rest.PrivateKeyConf{
				{
					Fingerprint: internal.Fingerprint,
					KeyFile:     priFile,
				},
			},
		},
	}
	engine := rest.MustNewServer(c)
	defer engine.Stop()

	engine.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/a/b",
		Handler: handle,
	}, rest.WithSignature(c.Signature))
	fmt.Println("start server...")
	engine.Start()
}

func createTempFile(body []byte) (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "go-unit-*.tmp")
	if err != nil {
		return "", err
	}

	tmpFile.Close()
	err = ioutil.WriteFile(tmpFile.Name(), body, os.ModePerm)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
