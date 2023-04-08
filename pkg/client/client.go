package client

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/STARRY-S/bilibili-danmaku-client/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Client struct {
	roomID        int
	ws            data.WsConnection
	sendPackageCh chan *data.Package

	OnMessage func(*data.Package) error

	exitCh chan struct{}
}

func NewClient(rid int) *Client {
	return &Client{
		roomID: rid,

		OnMessage: defaultOnMessage,
		exitCh:    make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	if c.roomID <= 0 {
		return fmt.Errorf("invalid room ID [%d]", c.roomID)
	}

	var err error
	c.ws, err = websocket.Dial(utils.DanmakuURL, "", utils.DanmakuOrigin)
	if err != nil {
		return fmt.Errorf("Connect: %w", err)
	}

	// ==========================
	// send first data
	fd := data.GetFirstData(c.roomID)
	if fd == nil {
		return fmt.Errorf("Connect: failed to get first data")
	}
	logrus.Infof("fd: %v", string(fd))
	pkg := data.NewPackage(fd, data.PV_1, data.O_7)
	d := pkg.Encode()
	n, err := c.ws.Write(d)
	if err != nil {
		return err
	}
	logrus.Infof("Send response: %v", n)

	// server connected
	if c.ws.IsClientConn() {
		logrus.Infof("Server connected")
	}

	// prepare go routines
	c.onConnect()

	// main routine blocks here
	<-time.After(time.Minute * 30)
	close(c.exitCh)

	logrus.Warnf("!!!EXIT SIGNAL TRIGGERED!!!")
	<-time.After(time.Second * 30)

	return nil
}

func (c *Client) onConnect() {
	go c.sendPackageRoutine()
	go c.sendHeartBeatRoutine()
	go c.readWsRoutine()
}

func (c *Client) sendData(d []byte, p data.Protocol, o data.Operation) error {
	pkg := data.NewPackage(d, p, o)
	c.sendPackageCh <- pkg

	return nil
}

func (c *Client) handleMessage(msg string) {

}

func defaultOnMessage(pkg *data.Package) error {
	msg := fmt.Sprintf("{PackageLength: %v, HeaderLength: %v, PV: %v, OP: %v}\n",
		pkg.PackageLength, pkg.HeaderLength, pkg.ProtocolVersion,
		pkg.Operation)
	logrus.Debugf("defaultOnMessage %v", msg)

	switch pkg.ProtocolVersion {
	case data.PV_0: // JSON plantext
		// broadcast junk message, discard it
		// logrus.Debugf("OnMessage: %v", string(pkg.Data))
	case data.PV_1: // uncompressed data and popularity value
		switch pkg.Operation {
		case data.O_3: // 气人值
			var u uint32
			r := bytes.NewReader(pkg.Data)
			binary.Read(r, binary.BigEndian, &u)
			logrus.Infof("气人值: %v", u)
		case data.O_8: // 进房间成功
			logrus.Infof("Response: %v", string(pkg.Data))
		}
	case data.PV_2: // gzip compressed JSON
		if pkg.Operation != 5 {
			return fmt.Errorf("unknow data: protocol 2, operation %v",
				pkg.Operation)
		}
		zr, err := zlib.NewReader(bytes.NewReader(pkg.Data))
		if err != nil {
			return fmt.Errorf("defaultOnMessage: gzip.NewReader: %w", err)
		}
		data, err := io.ReadAll(zr)
		if err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("defaultOnMessage: io.ReadAll: %w", err)
		}
		// danmaku data
		// TODO: EXPERIMENTAL
		str := string(data)
		start := strings.Index(str, "{")
		end := strings.LastIndex(str, "}")
		str = str[start : end+1]
		str = strings.Replace(str, ">", "", -1)
		fmt.Printf("XXXX %v\n", str)

		var obj interface{}
		d := json.NewDecoder(bytes.NewBufferString(str))
		err = d.Decode(&obj)
		if err != nil {
			logrus.Error(err)
		}
		// logrus.Infof("XXXXXX: %++v", obj)

		out, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("XXXX %v", string(out))

	case data.PV_3: // brotli compressed data
		logrus.Infof("unsupported bortli compressed data")
	}

	return nil
}
