package dispatch_worker

import(
  "log"
  "sync"
  "time"
  "math/rand"
  "runtime"
)

// worker管理のための構造体
// workerを予め生成しておいて、要求に応じて利用可能なworkerを渡す
type Dispatcher struct {
  // worker pool
  pool  chan *Worker
  // interface型はメソッドだけを定義するもの
  // 空のinterfaceはなんでも突っ込める
  // つまりchan interface{}はなんでも突っ込めるchannel
  queue chan interface{}
  // workers
  workers []*Worker
  // wait group
  wg sync.WaitGroup
  // 終了を受け付けるchannel
  quit chan struct{}
}

type Worker struct {
  // このworkerを管理しているdispatcher
  dispatcher *Dispatcher
  // workerに渡す処理
  data chan interface{}
  // 終了を受け付けるchannel
  quit chan struct{}
  // worker id
  id int
}

const maxWorkers = 5
const maxQueues = 10000

func NewDispatcher() *Dispatcher {
  // dispatcherの初期化
  d := &Dispatcher{
    // capacityはwokerの数
    pool:  make(chan *Worker, maxWorkers),
    // queueはメッセージのキューイングの数
    queue: make(chan interface{}, maxQueues),
    quit:  make(chan struct{}),
  }

  d.workers = make([]*Worker, cap(d.pool))
  for i := 0; i < cap(d.pool); i++ {
    w := Worker{
      dispatcher: d,
      data: make(chan interface{}),
      quit: make(chan struct{}),
      id: i,
    }
    d.workers[i] = &w
  }
  return d
}

// Dispaterのインスタンスメソッド的なやつ
func (d *Dispatcher) start() {
  // rangeはrubyで言うところのeachに近い
  // arrayにたいしてfor rangeした場合は、indexとvalueが渡ってくる
  // 以下のコードの場合は、_ にインデックス、wにvaleu(worker)がworkers分だけ渡る
  for _, w := range d.workers {
    w.start()
  }

  // goroutine
  go func() {
    // 無限ループ
    for {
      select {
      case v := <-d.queue:
        // poolからworkerをとりだして、dataにqueueから渡されたvを渡す
        (<-d.pool).data <- v
      case <-d.quit:
        log.Println("dispatcher killed.")
        return
      }
    }
  }()
}

// 処理を追加するメソッド
func (d *Dispatcher) add(v interface{}) {
  // キューイングされた場合に待機するためにWaitGroupをカウントアップ
  d.wg.Add(1)
  // キューイング
  d.queue <- v
}

// すべてのwgが終わるのを待つ
func (d *Dispatcher) wait() {
  d.wg.Wait()
}

func (w *Worker) start() {
  go func() {
    for {
      // worker poolに自分自身を追加する
      w.dispatcher.pool <- w

      select {
      case v:= <-w.data:
        time.Sleep(1 * time.Second)
        log.Printf("Worker ID:%d, Goroutine:%d, %s", w.id, runtime.NumGoroutine(), v)

        // WaitGroupのカウントダウン
        w.dispatcher.wg.Done()
      case <-w.quit:
        log.Println("worker killed.")
        return
      }
    }
  }()
}

func main() {
  rand.Seed(time.Now().UnixNano())
  d := NewDispatcher()
  d.start()
  for i := 0; i < 100; i++ {
    d.add("hogehoge")
  }

  d.wait()
}
