# CouchDB Hooks for [Logrus](https://github.com/sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Install
With `go get`:

```shell
$ go get github.com/codehakase/couchdbrus
```

## Usage

```go
package main

import (
	"github.com/sirupsen/logrus"
  "github.com/codehakase/couchdbrus"
	"github.com/ryanjyoder/couchdb"  
)

// see test for sample implementation of client
var client *couchdb.Client

func main() {
  log := logrus.New()
  couchdbHook, err := couchdbrus.NewHook(client, "mylogdatabaseName")
  if err != nil {
    // do proper error handling here...
    log.Fatalf("could not create couchdb hook: %v", 
    err)
  }
  log.Hooks.Add(couchdbHook)

  log.WithFields(logrus.Fields{
    "hostname": "hakaselabs",
    "source":   "spacex",
    "tag":      "test-tag",
   }).Info("Hello Captain we can see the moon!")
}
```


## Author
Francis Sunday - [@codehakase](https://twitter.com/codehakase)
