package api

import (
	"context"
	//"errors"
	//"fmt"
	//"net"
	"net/http"
	//"sync"

	"git.zluudg.se/piplup/internal/common"
)

type Conf struct {
	Active      bool   `json:"active"`
	Address     string `json:"address"`
	Port        string `json:"port"`
	Log         common.Logger
	Application application
}

type myApi struct {
	id              string
	listenInterface string
	app             application
	log             common.Logger
	srv             http.Server
}

type application interface {
}

func Create(conf Conf) (*myApi, error) {
	a := new(myApi)
	a.id = "api"

	//if !conf.Active {
	//	return nil, errors.New("attempting to create inactive api")
	//}

	//if conf.Log == nil {
	//	return nil, errors.New("attempting to create api without logger")
	//}

	//if conf.Address == "" {
	//	return nil, errors.New("attempting to create api without address")
	//}

	//if conf.Port == "" {
	//	return nil, errors.New("attempting to create api without port")
	//}

	//a.listenInterface = net.JoinHostPort(conf.Address, conf.Port)
	//a.log = conf.Log
	return a, nil
}

func (ma *myApi) Run(ctx context.Context, exitCh chan<- common.Exit) {
	//if parent == nil {
	//	ma.log.Error("Bad parent handle for api")
	//	return
	//}
	//ma.parent = parent

	//http.HandleFunc("GET /api/ping", ma.handlerPing)
	//http.HandleFunc("GET /api/ok", ma.handlerOk)
	//http.HandleFunc("GET /api/err", ma.handlerErr)

	//server := http.Server{
	//	Addr: ma.listenInterface,
	//}

	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()

	//	err := server.ListenAndServe()
	//	if errors.Is(err, http.ErrServerClosed) {
	//		ma.log.Info("API server closing")
	//	} else {
	//		ma.log.Error("Unexpected API server shutdown: '%s'", err)
	//	}

	//}()

	//<-ctx.Done()
	//ma.log.Info("Shutting down API")
	//server.Shutdown(ctx)
	//wg.Wait()

	exitCh <- common.Exit{ID: ma.id, Err: nil}
	return
}

//func (ma *myApi) handlerPing(rw http.ResponseWriter, r *http.Request) {
//	n := ma.parent.GetPmPing()
//	msg := fmt.Sprintf("{\"n_ping\": %d}", n)
//	rw.Write([]byte(msg))
//}
//
//func (ma *myApi) handlerOk(rw http.ResponseWriter, r *http.Request) {
//	n := ma.parent.GetPmOk()
//	msg := fmt.Sprintf("{\"n_ok\": %d}", n)
//	rw.Write([]byte(msg))
//}
//
//func (ma *myApi) handlerErr(rw http.ResponseWriter, r *http.Request) {
//	n := ma.parent.GetPmErr()
//	msg := fmt.Sprintf("{\"n_err\": %d}", n)
//	rw.Write([]byte(msg))
//}
