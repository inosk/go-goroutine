package goroutinepractice

import (
  "log"
)

func worker ( queue chan string, ctrl chan bool ) {
  for {
    select {
    case message := <-queue:
      log.Println(message)
    case _ctrl := <-ctrl:
      if _ctrl {
        log.Println("worker killed")
        return
      }
    }
  }
}

func main() {
  queue := make(chan string)
  ctrl := make(chan bool)

  go worker(queue, ctrl)

  queue <- "hoge"
  queue <- "fuga"
  queue <- "piyo"
  ctrl  <- false
  queue <- "moge"
  ctrl  <- true
  // NOTE:
  // channelを読み出すgoroutingがいないのでdeadlockエラー
  //queue <- "hoge"
}
