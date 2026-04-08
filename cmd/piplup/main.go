package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.zluudg.se/piplup/internal/api"
	"git.zluudg.se/piplup/internal/app"
	"git.zluudg.se/piplup/internal/cert"
	"git.zluudg.se/piplup/internal/common"
	"git.zluudg.se/piplup/internal/logger"
)

const c_APP_IDENTIFIER = "piplup"

/* Rewritten if building with make */
var version = "BAD-BUILD"
var commit = "BAD-BUILD"

type conf struct {
	app.Conf
	ApiConf  api.Conf  `json:"api"`
	CertConf cert.Conf `json:"cert"`
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
		"config.json",
		"Configuration file to use",
	)
	flag.BoolVar(&debugFlag,
		"debug",
		false,
		"Enable DEBUG logs",
	)
	flag.Parse()

	log := logger.New(
		logger.Conf{
			Debug: debugFlag,
		})

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
	defer file.Close()

	confDecoder := json.NewDecoder(file)
	if confDecoder == nil {
		log.Error("Problem decoding config file '%s', exiting...", configFile)
		os.Exit(-1)
	}

	confDecoder.DisallowUnknownFields()
	err = confDecoder.Decode(&mainConf)
	if err != nil {
		log.Error("Problem decoding config file '%s': %s", configFile, err)
		os.Exit(-1)
	}

	mainConf.Debug = mainConf.Debug || debugFlag
	mainConf.CertConf.Debug = mainConf.CertConf.Debug || debugFlag

	/*
	 ******************************************************************
	 ********************** SET UP CERT *******************************
	 ******************************************************************
	 */
	certHandle, err := cert.Create(mainConf.CertConf)
	if err != nil {
		log.Error("Error creating cert handler: '%s'", err)
		os.Exit(-1)
	}

	/*
	 ******************************************************************
	 ********************** SET UP APP ********************************
	 ******************************************************************
	 */
	mainConf.Cert = certHandle
	appHandle, err := app.Create(mainConf.Conf)
	if err != nil {
		log.Error("Error creating application: '%s'", err)
		os.Exit(-1)
	}

	/*
	 ******************************************************************
	 ********************** SET UP METRICS ****************************
	 ******************************************************************
	 */
	//mainConf.Api.Application = application
	//appApi, err := api.Create(mainConf.Api)
	//if err != nil {
	//	log.Error("Error creating API: '%s'", err)
	//	os.Exit(-1)
	//}

	/*
	 ******************************************************************
	 ********************** START RUNNING STUFF ***********************
	 ******************************************************************
	 */
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigChan)
	defer signal.Stop(sigChan)

	ctx, cancel := context.WithCancel(context.Background())
	exitCh := make(chan common.Exit, 128)

	log.Info("Starting threads...")

	var wg sync.WaitGroup
	wg.Go(func() { certHandle.Run(ctx, exitCh) })
	wg.Go(func() { appHandle.Run(ctx, exitCh) })

MAIN_LOOP:
	for {
		select {
		case s, ok := <-sigChan:
			if ok {
				log.Info("Got signal '%s'", s)
			} else {
				log.Warning("Signal channel closed unexpectedly, exiting...")
			}
			break MAIN_LOOP
		case exit, ok := <-exitCh:
			if ok {
				if exit.Err != nil {
					log.Error("%s exited with error: '%s'", exit.ID, exit.Err)
					if exit.Err == common.ErrFatal {
						log.Error("%s encountered fatal error, exiting...", exit.ID)
						break MAIN_LOOP
					}
				} else {
					log.Info("%s done!", exit.ID)
				}
			} else {
				log.Warning("Exit channel closed unexpectedly, exiting...")
				break MAIN_LOOP
			}
		}
	}

	log.Info("Cancelling threads")
	cancel()

	log.Info("Waiting for threads to finish")
	wg.Wait()

	log.Info("Exiting...")
	os.Exit(0)
}
