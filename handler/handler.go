package handler

import (
  "bufio"
  "io"
  "fmt"
  "encoding/json"
  "github.com/dpb587/logsearch-checks/check"
)

type Emitter func(check.Status)
type HandlerCallback func(check.Status, Emitter)

type Handler struct {
    Reader io.Reader
    Processor HandlerCallback
}

func emit(cs check.Status) {
    csjs, err := json.Marshal(cs)

    if nil != err {
        panic(err)
    }

    fmt.Println(string(csjs))
}

func (sh *Handler) Run() {
    stdin := bufio.NewScanner(sh.Reader)

    var i int

    for stdin.Scan() {
        i = i + 1

        var mc check.Status

        if err := json.Unmarshal([]byte(stdin.Text()), &mc); err != nil {
            panic(err)
        }

        sh.Processor(mc, emit)
    }

    if nil != stdin.Err() {
        panic(stdin.Err())
    }
}
