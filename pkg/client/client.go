package client

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/data"
	"github.com/andybalholm/brotli"
	"github.com/sirupsen/logrus"
)

type Client struct {
	roomID        int
	ws            data.WsConnection
	sendPackageCh chan *data.Package
	readPackageCh chan *data.Package
	voiceStringCh chan string
	errorCh       chan error

	wg *sync.WaitGroup

	// Callback function when receive one websocket package
	OnPackage func(*data.Package) error

	popularity int      // 气人值
	rankList   []string // 高能榜 (正在观看的人)
	watched    int      // 看过的人数

	mutex *sync.RWMutex
}

func NewClient(rid int) *Client {
	c := &Client{
		roomID:        rid,
		sendPackageCh: make(chan *data.Package),
		readPackageCh: make(chan *data.Package),
		voiceStringCh: make(chan string),
		errorCh:       make(chan error),
		mutex:         new(sync.RWMutex),
		wg:            new(sync.WaitGroup),
	}
	c.OnPackage = c.defaultOnPackage
	return c
}

func (c *Client) SetPopolarity(i int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.popularity = i
}

func (c *Client) Popularity() int {
	c.mutex.RLock()
	defer c.mutex.Unlock()
	return c.popularity
}

func (c *Client) SetRankList() {
	return
}

func (c *Client) GetRankList() {
	return
}

func (c *Client) defaultOnPackage(pkg *data.Package) error {
	dataLength := int(pkg.PackageLength - uint32(pkg.HeaderLength))
	if len(pkg.Data) != dataLength && len(pkg.Data) != 19 {
		// heart beat respond package length is 19 (actual 4), ignore it
		logrus.Debugf("Warning: data length: %d, should be %d",
			len(pkg.Data), dataLength)
	}
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
			c.SetPopolarity(int(u))
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
		if pkg.Operation != 5 {
			return fmt.Errorf("unknow data: protocol 3, operation %v",
				pkg.Operation)
		}
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
