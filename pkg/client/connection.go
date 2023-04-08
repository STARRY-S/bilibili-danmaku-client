package client

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/sirupsen/logrus"
)

const (
	bufferLength = 1024
)

func (c *Client) sendPackageRoutine(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case pkg := <-c.sendPackageCh:
			d := pkg.Encode()
			n, err := c.ws.Write(d)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("Sent package length: %v", n)
		case <-ctx.Done():
			close(c.sendPackageCh)
			logrus.Debug("sendPackageRoutine exited gracefully")
			return
		}
	}
}

func (c *Client) sendHeartBeatRoutine(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case <-time.After(time.Second * 30):
			// send heartbeat package every 30 seconds
			err := c.sendData(data.HeartBeatData, data.PV_1, data.O_2)
			if err != nil {
				logrus.Error(err)
			}
		case <-ctx.Done():
			logrus.Debugf("sendHeartBeatRoutine exited gracefully")
			return
		}
	}
}

func (c *Client) readPackageRoutine(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case pkg := <-c.readPackageCh:
			// OnPackage callback to handle package
			if err := c.OnPackage(pkg); err != nil {
				logrus.Error(err)
			}
		case <-ctx.Done():
			logrus.Debugf("readPackageRoutine exited gracefully")
			return
		}
	}
}

func (c *Client) readWsRoutine(ctx context.Context) {
	for {
		// This routine does not need to exit gracefully since it
		// will block here if no message read from ws connection
		pkg, err := c.readWsPackage()
		if err != nil {
			logrus.Error(err)
		}
		if pkg == nil {
			continue
		}

		c.readPackageCh <- pkg
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 100):
		}
	}
}

func (c *Client) readWsPackage() (*data.Package, error) {
	buff := make([]byte, bufferLength)
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
