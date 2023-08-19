package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/config"
	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/STARRY-S/bilibili-danmaku-client/pkg/voice"
	"github.com/STARRY-S/bilibili-danmaku-client/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

const (
	bufferLength = 1024
)

func (c *Client) Connect() error {
	if c.roomID <= 0 {
		return fmt.Errorf("invalid room ID [%d]", c.roomID)
	}

	err := c.buildWsConnection()
	if err != nil {
		return fmt.Errorf("Connect: %w", err)
	}

	// prepare go routines
	c.prepareRoutines()

	// waiting all routine stop
	c.wg.Wait()

	logrus.Infof("Client stopped gracefully")

	return nil
}

func (c *Client) buildWsConnection() error {
	var err error
	c.ws, err = websocket.Dial(utils.DanmakuURL, "", utils.DanmakuOrigin)
	if err != nil {
		return err
	}
	// send first data
	d := data.GetFirstData(c.roomID)
	if d == nil {
		return fmt.Errorf("failed to get first data")
	}
	pkg := data.NewPackage(d, data.PV_1, data.O_7)
	_, err = c.ws.Write(pkg.Encode())
	if err != nil {
		return err
	}

	// server connected
	if c.ws.IsClientConn() {
		logrus.Infof("Server connected")
	}

	return nil
}

func (c *Client) prepareRoutines() {
	ctx, stop := context.WithCancel(context.Background())
	go c.sendPackageRoutine(ctx)
	go c.sendHeartBeatRoutine(ctx)
	go c.readPackageRoutine(ctx)
	go c.handleAudioRoutine(ctx)
	go c.handleErrorRoutine(ctx)

	c.wg.Add(5)
	// do not wait readWsRoutine since it may in blocked status
	go c.readWsRoutine(ctx)

	// handle SIGINT gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			s := <-sig
			logrus.Debugf("signal received: %v", s)
			if s != os.Interrupt {
				continue
			}
			stop()
			// force exit if not stopped gracefully
			<-sig
			os.Exit(1)
		}
	}()
}

func (c *Client) sendPackageRoutine(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case pkg := <-c.sendPackageCh:
			d := pkg.Encode()
			n, err := c.ws.Write(d)
			if err != nil {
				c.errorCh <- fmt.Errorf("sendPackageRoutine: %v", err)
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
		// send heartbeat package every 30 seconds
		pkg := data.NewPackage(data.HeartBeatData, data.PV_1, data.O_2)
		select {
		case c.sendPackageCh <- pkg:
		case <-time.After(time.Millisecond * 100):
			logrus.Warnf("send package failed")
		}

		select {
		case <-time.After(time.Second * 30):
			continue
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
				c.errorCh <- fmt.Errorf("readPackageRoutine: %w", err)
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
		if err != nil && errors.Is(err, io.EOF) {
			logrus.Warnf("Disconnected, reconnecting...")
			if err := c.buildWsConnection(); err != nil {
				c.errorCh <- err
				ctx.Done()
				return
			}
		}
		if pkg == nil {
			continue
		}

		select {
		case c.readPackageCh <- pkg:
		case <-ctx.Done():
			return
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
		pkg.Data = append(pkg.Data, buff[:n]...)
		readLen += n
	}

	return pkg, nil
}

func (c *Client) handleAudioRoutine(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case str := <-c.voiceStringCh:
			go func() {
				if !config.GetBool("voice") {
					return
				}
				if err := voice.NewVoice(str).Say(); err != nil {
					c.errorCh <- err
				}
			}()
		case <-ctx.Done():
			logrus.Debugf("handleAudioRoutine exited gracefully")
			return
		}
	}
}

func (c *Client) handleErrorRoutine(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case e := <-c.errorCh:
			logrus.Error(e)
		case <-ctx.Done():
			logrus.Debugf("handleErrorRoutine exited gracefully")
			return
		}
	}
}
