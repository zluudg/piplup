package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.zluudg.se/piplup/internal/api"
	"git.zluudg.se/piplup/internal/app"
	"git.zluudg.se/piplup/internal/common"
	"git.zluudg.se/piplup/internal/logger"
)

/* Rewritten if building with make */
var version = "BAD-BUILD"
var commit = "BAD-BUILD"

type conf struct {
	Address         string   `json:"address"`
	UdpPort         string   `json:"udp_port"`
	TlsPort         string   `json:"tls_port"`
	UpstreamAddress string   `json:"upstream_address"`
	UpstreamPort    string   `json:"upstream_port"`
	Inject          string   `json:"inject"`
	MatchSuffix     string   `json:"match_suffix"`
	CertDir         string   `json:"cert_dir"`
	Api             api.Conf `json:"api"`
}

func main() {
	var configFile string
	var runVersionCmd bool
	var debugFlag bool
	var mainConf conf

	flag.BoolVar(&runVersionCmd,
		"version",
		false,
		"Print version then exit",
	)
	flag.StringVar(&configFile,
		"config",
		"",
		"Configuration file to use",
	)
	flag.BoolVar(&debugFlag,
		"debug",
		false,
		"Enable DEBUG logs",
	)
	flag.Parse()

	log, err := logger.Create(
		logger.Conf{
			Debug: debugFlag,
		})
	if err != nil {
		panic(fmt.Sprintf("Could not create logger, err: '%s'", err))
	}

	log.Info("piplup version: '%s', commit: '%s'", version, commit)
	if runVersionCmd {
		os.Exit(0)
	}

	log.Debug("Debug logging enabled")

	if configFile == "" {
		log.Error("No config file specified, exiting...")
		os.Exit(-1)
	}

	file, err := os.Open(configFile)
	if err != nil {
		log.Error("Couldn't open config file '%s', exiting...", configFile)
		os.Exit(-1)
	}

	confDecoder := json.NewDecoder(file)
	if confDecoder == nil {
		log.Error("Problem decoding config file '%s', exiting...", configFile)
		os.Exit(-1)
	}

	confDecoder.DisallowUnknownFields()
	confDecoder.Decode(&mainConf)

	application, err := app.NewBuilder().
		Logger(log).
		Address(mainConf.Address).
		UdpPort(mainConf.UdpPort).
		TlsPort(mainConf.TlsPort).
        Upstream(mainConf.UpstreamAddress, mainConf.UpstreamPort).
		CertDir(mainConf.CertDir).
		MatchSuffix(mainConf.MatchSuffix).
		Inject(mainConf.Inject).
		Build()
	if err != nil {
		log.Error("Error creating application: '%s'", err)
		os.Exit(-1)
	}

	mainConf.Api.Log = log
	mainConf.Api.Application = application
	appApi, err := api.Create(mainConf.Api)
	if err != nil {
		log.Error("Error creating API: '%s'", err)
		os.Exit(-1)
	}

	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	exitCh := make(chan common.Exit)

	go application.Run(ctx, exitCh)

	go appApi.Run(ctx, exitCh)

	exitLoop := false
	for {
		select {
		case s := <-sigChan:
			log.Info("Got signal '%s'", s)
			exitLoop = true
		case exit := <-exitCh:
			if exit.Err != nil {
				log.Error("%s exited with error: '%s'", exit.ID, exit.Err)
				if exit.Err == common.ErrFatal {
					exitLoop = true
				}
			} else {
				log.Info("%s done!", exit.ID)
			}
		}
		if exitLoop {
			break
		}
	}

	cancel()
	log.Info("Cancelling, giving threads some time to finish...")
	time.Sleep(2 * time.Second)
	close(exitCh)
	log.Info("Exiting...")
	os.Exit(0)
}
