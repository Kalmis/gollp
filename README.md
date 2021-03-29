# gollp

Converts MLLP to HTTP. Inspired by the [mllp-http](https://github.com/rivethealth/mllp-http) package.

# Overview

Listens a TCP socket for MLLP messages and sends them via HTTP post to a given URL. Message in POST response is sent back as an MLLP message. Multiple messages can be sent in one TCP connection, each message is handled separately one-by-one.

# Usage

From Github releases you can download 64bit binaries for Linux, OS X and Windows, or simply clone and build.

CLI Commands
```
  -help
        Show help
  -ip string
        Address to be listened, e.g. localhost or 0.0.0.0. (default "localhost")
  -port int
        Port to be listened.  (default 2575)
  -url string
        Target URL where data is sent as HTTP POST
```

Listen to `localhost` on port `2575` and route messages to `https://ptsv2.com/t/tz0jp-1616963104/post` (you can view results in [here](https://ptsv2.com/t/tz0jp-1616963104))

```
$ ./gollp --ip localhost --port 2575 --url https://ptsv2.com/t/tz0jp-1616963104/post

Starting to listen localhost:2575 and routing messages to https://ptsv2.com/t/tz0jp-1616963104/post
INFO: Processing message
MSH||||
INFO: Processing message
MSH||||
INFO: EOF reached, closing connection
```

## Docker

Gollp can also be run in Docker container

```
docker build -t gollp .
docker run -it -p 3010:3010 gollp --ip 0.0.0.0 --port 3010 --url=https://ptsv2.com/t/bq56z-1617043531/post
```

# HTTP Data format

Gollp sends messages forward in JSON format and assumes to receive response in same format. Currently the JSON is very simple. Message does not include the MLLP starting or ending block (including the last CR)
```
{
    "Message": "MSH|||||"
}
```