package main

import (
  "os"
  "github.com/dpb587/logsearch-checks/check"
  "github.com/dpb587/logsearch-checks/handler"
)

func mhandle(cs check.Status, emit handler.Emitter) {
    if check.CHECK_OKAY != cs.Check.GetStatus() {
        emit(cs)
    }
}

func main() {
    v := handler.Handler{os.Stdin, mhandle}
    v.Run()
}
