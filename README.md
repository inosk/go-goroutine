# これはなに

go ルーチンの仕組みを理解するためのコード

# コードの説明

- goroutingepractice.go
  - 超単純にchannelを使ってgoルーチン間でデータをやり取りする
- dispatch_worker.go
  - worker poolの仕組みをgoで実装してみるテスト

# 幾つかのtips

- channel
  - go ルーチン間でデータをやり取りするための仕組み
  - x := <-channel でchannelからデータをよみだし
  - channel <-xでchannelにデータを渡す
  - channelにはmake した時に第二引数としてchannelの容量を指定できる
    - ex.) make(chan string, 2) 容量２でstring型が通れるchannel
    - 無指定の場合はbufferなし
  - 容量がいっぱいになるまでではchannelへの送信はブロックされない
  - 容量がいっぱいになると容量ができるばでブロック
  - channelのbufferが空の場合、読み出しは容量にかかわらずブロックされる
  - channelにはcapとlenという組み込み関数からbufferサイズと現在のbufferが埋まっている数をとれる
    - ex.) cap(ch1) -> makeの第二引数の数字が取れる。len(ch1) -> 埋まっているbufferの数が取れる

- interface
  - interfaceはstructと似ているがメソッド型しか定義できないという違いがある
  - interfaceはある構造体が同じメソッドを持っていれば代入可能
  - 空のinterface(interface{})はメソッドの定義がないので何でも突っ込める変数になる
  - 空のinterfaceでchannelを定義すればなんでも突っ込めるchannelができる(多分あまり推奨されない)
