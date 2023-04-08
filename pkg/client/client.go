package client

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/STARRY-S/bilibili-danmaku-client/utils"
	"github.com/andybalholm/brotli"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Client struct {
	roomID        int
	ws            data.WsConnection
	sendPackageCh chan *data.Package
	readPackageCh chan *data.Package

	wg *sync.WaitGroup

	// Callback function when receive one websocket package
	OnPackage func(*data.Package) error
}

func NewClient(rid int) *Client {
	c := &Client{
		roomID:        rid,
		sendPackageCh: make(chan *data.Package),
		readPackageCh: make(chan *data.Package),
		wg:            &sync.WaitGroup{},
	}
	c.OnPackage = c.defaultOnPackage
	return c
}

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
	c.wg.Add(3)
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

func (c *Client) sendData(d []byte, p data.Protocol, o data.Operation) error {
	c.sendPackageCh <- data.NewPackage(d, p, o)
	return nil
}

func (c *Client) defaultOnPackage(pkg *data.Package) error {
	switch pkg.ProtocolVersion {
	case data.PV_0: // JSON plantext
		// Can be broadcast junk message or decompressed danmaku message
		c.handleMessage(pkg.Data)
	case data.PV_1: // uncompressed data and popularity value
		switch pkg.Operation {
		case data.O_3: // 气人值
			var u uint32
			r := bytes.NewReader(pkg.Data)
			binary.Read(r, binary.BigEndian, &u)
			logrus.Infof("气人值: %v", u)
		case data.O_8: // 进房间成功
			logrus.Debugf("Response: %v", string(pkg.Data))
		}
	case data.PV_2: // gzip compressed JSON
		if pkg.Operation != 5 {
			return fmt.Errorf("unknow data: protocol 2, operation %v",
				pkg.Operation)
		}
		zr, err := zlib.NewReader(bytes.NewReader(pkg.Data))
		if err != nil {
			return fmt.Errorf("defaultOnPackage: %w", err)
		}
		d, err := io.ReadAll(zr)
		if err != nil {
			return fmt.Errorf("defaultOnPackage: %w", err)
		}
		c.handleDecompressedPackage(d)
	case data.PV_3: // brotli compressed data
		logrus.Debugf("unsupported bortli compressed data")
		br := brotli.NewReader(bytes.NewReader(pkg.Data))
		d, err := io.ReadAll(br)
		if err != nil {
			return fmt.Errorf("defaultOnPackage: %w", err)
		}
		c.handleDecompressedPackage(d)
	}

	return nil
}

func (c *Client) handleDecompressedPackage(d []byte) error {
	pkg := &data.Package{}
	err := pkg.DecodeHead(d)
	if err != nil {
		return fmt.Errorf("handleDecompressedMessage: %w", err)
	}

	pkg.Data = make([]byte, pkg.PackageLength-data.PkgHeaderLen)
	copy(pkg.Data, d[data.PkgHeaderLen:pkg.PackageLength+1])

	c.OnPackage(pkg)

	// use recursion to handle multiple decompressed package
	if len(d) > int(pkg.PackageLength) {
		return c.handleDecompressedPackage(d[pkg.PackageLength:])
	}

	return nil
}
