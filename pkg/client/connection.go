package client

import (
	"errors"
	"io"
	"time"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/sirupsen/logrus"
)

func (c *Client) sendPackageRoutine() {
	c.sendPackageCh = make(chan *data.Package)
	for {
		select {
		case pkg := <-c.sendPackageCh:
			d := pkg.Encode()
			n, err := c.ws.Write(d)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("Sent package length: %v", n)
		case <-c.exitCh:
			logrus.Debug("Send package routine exited gracefully.")
			return
		}
	}
}

func (c *Client) sendHeartBeatRoutine() {
	// send heartbeat package first
	err := c.sendData(data.HeartBeatData, data.PV_1, data.O_2)
	if err != nil {
		logrus.Error(err)
	}
	for {
		select {
		case <-time.After(time.Second * 30):
			// send heartbeat package every 30 seconds
			err = c.sendData(data.HeartBeatData, data.PV_1, data.O_2)
			if err != nil {
				logrus.Error(err)
			}
		case <-c.exitCh:
			logrus.Debugf("heartbeat routine exited gracefully")
			return
		}
	}
}

func (c *Client) readWsRoutine() {
	// read message loop
	for {
		pkg, err := c.readWsPackage()
		if err != nil {
			logrus.Error(err)
		}
		if pkg == nil {
			continue
		}
		// OnMessage callback to handle package
		if err := c.OnMessage(pkg); err != nil {
			logrus.Error(err)
		}

		select {
		case <-c.exitCh:
			logrus.Debugf("readWsRoutine exited gracefully")
			return
		default:
		}
	}
}

func (c *Client) readWsPackage() (*data.Package, error) {
	buffLen := 64
	buff := make([]byte, buffLen)
	n, err := c.ws.Read(buff)
	if err != nil && errors.Is(err, io.EOF) {
		// websocket connection closed
		return nil, err
	}

	if n < int(data.PkgHeaderLen) {
		// unknow data, discard
		return nil, nil
	}

	pkg := &data.Package{}
	err = pkg.DecodeHead(buff)
	if err != nil {
		return nil, err
	}
	pkg.Data = make([]byte, n-int(data.PkgHeaderLen))
	copy(pkg.Data, buff[data.PkgHeaderLen:])

	readLen := n
	for readLen < int(pkg.PackageLength) {
		n, err = c.ws.Read(buff)
		if err != nil && errors.Is(err, io.EOF) {
			// websocket connection closed
			return nil, err
		}
		pkg.Data = append(pkg.Data, buff...)
		readLen += n
	}

	return pkg, nil
}
