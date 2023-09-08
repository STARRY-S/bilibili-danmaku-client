# bilibili-danmaku-client

A simple [Bilibili](https://bilibili.com) danmaku CLI (for [Bilibili Live](https://live.bilibili.com)) written in Go.

## Usage

Ensure Go 1.20+ installed.

```console
$ git clone https://github.com/STARRY-S/bilibili-danmaku-client.git && cd bilibili-danmaku-client
$ go build -o bilibili-danmaku-client .

$ ./bilibili-danmaku-client -h
Bilibili Danmaku Client Go
......

$ ./bilibili-danmaku-client -r [ROOM ID] -u [USER ID] --voice
INFO[0000] Voice output enabled
INFO[0000] Server connected
INFO[0000] 气人值: 1
......
```

## License

```
Copyright (c) 2023 STARRY-S

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
