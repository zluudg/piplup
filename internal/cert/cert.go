package cert

import (
	"context"
	"crypto/tls"
	"math"
	"path/filepath"
	"sync"
	"time"

	"git.zluudg.se/piplup/internal/common"
	"git.zluudg.se/piplup/internal/logger"
)

type Conf struct {
	Active   bool   `json:"active"`
	Debug    bool   `json:"debug"`
	Interval int    `json:"interval"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
	Log      common.Logger
}

type certHandle struct {
	active    bool
	id        string
	log       common.Logger
	ticker    *time.Ticker
	certPath  string
	keyPath   string
	certStore certs
}

type certs struct {
	sync.RWMutex
	certificate *tls.Certificate
}

func Create(conf Conf) (*certHandle, error) {
	c := new(certHandle)
	c.id = "cert manager"
	c.active = conf.Active

	if conf.Log == nil {
		log := logger.New(
			logger.Conf{
				Debug: conf.Debug,
			})
		c.log = log
	} else {
		c.log = conf.Log
	}
	c.log.Debug("Debug logging enabled for %s", c.id)

	if conf.Cert == "" || conf.Key == "" {
		c.log.Error("Missing path to cert or key in config")
		return nil, common.ErrBadParam
	}
	c.certPath = filepath.Clean(conf.Cert)
	c.keyPath = filepath.Clean(conf.Key)

	if conf.Interval > 0 {
		c.ticker = time.NewTicker(time.Duration(conf.Interval) * time.Second)
	} else {
		c.ticker = time.NewTicker(time.Duration(math.MaxInt32) * time.Second)
		c.ticker.Stop()
		c.log.Warning("No interval set for scanning cert directory. Won't be refreshing.")
	}

	err := c.scanCert(context.Background())
	if err != nil {
		c.log.Error("First time scanning cert failed: %s", err)
		return nil, err
	}

	return c, nil
}

func (c *certHandle) Run(ctx context.Context, exitCh chan<- common.Exit) {
	defer c.ticker.Stop()
	if !c.active {
		exitCh <- common.Exit{ID: c.id, Err: nil}
		return
	}

CERT_LOOP:
	for {
		select {
		case <-ctx.Done():
			c.log.Info("Shutting down %s", c.id)
			break CERT_LOOP
		case t, ok := <-c.ticker.C:
			if ok {
				c.log.Debug("Cert scan tick %s", t)
				err := c.scanCert(ctx)
				if err != nil {
					if err == common.ErrFatal {
						exitCh <- common.Exit{ID: c.id, Err: err}
						return
					} else {
						c.log.Error("Failed scanning cert directory: %s", err)
					}
				} else {
					c.log.Debug("Re-scan of cert dir done!")
				}
			} else {
				c.log.Error("Ticker channel closed unexpectedly, exiting")
				exitCh <- common.Exit{ID: c.id, Err: common.ErrFatal}
				return
			}
		}
	}

	exitCh <- common.Exit{ID: c.id, Err: nil}
	c.log.Info("Shutdown done for %s", c.id)
	return
}

func (c *certHandle) GetCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	c.certStore.RLock()
	defer c.certStore.RUnlock()
	cert := c.certStore.certificate
	return cert, nil
}

func (c *certHandle) scanCert(ctx context.Context) error {
	certificate, err := tls.LoadX509KeyPair(c.certPath, c.keyPath)
	if err != nil {
		c.log.Error("Could not read certificate: %s", err)
		return common.ErrFatal
	}

	c.certStore.Lock()
	defer c.certStore.Unlock()
	c.certStore.certificate = &certificate

	return nil
}
