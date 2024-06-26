# runevision
automated red team tools

# todo v0.0.1
## application
- [ ] web ui
    - [ ] input target url
    - [ ] test recommended
    - [ ] discovered threats
    - [ ] display recommended exploit or cve サービスの特性に合わせて
    - [ ] list targets logs
- [ ] targetのDNSスキャン
    - crtshのpostgres, hacker targer, dnsdumpstarが精度高いのでそれで実装
- [ ] targetのPortスキャン

    ``` go
    // port scanはnaabuのこの辺りがポートスキャンの結果を得ている
    go s.ICMPResultWorker(ctx)
	go s.TCPResultWorker(ctx)
	go s.UDPResultWorker(ctx)
    ```
    ``` go
    // こいつが
    go s.ICMPResultWorker(ctx)
	go s.TCPResultWorker(ctx)
	go s.UDPResultWorker(ctx)

    // こいつが👆の関数で待機しているchannelに対して,トランスポート層から得たデータをloopBackScanCaseCallbackかtransportReaderCallbackを使ってスキャンの結果送信している
    func TransportReadWorker() {
    // このコードがnet.Dialをつかって接続できるか確認、接続できたらそのポートをScanResultsに返す
    func (r *Runner) handleHostPort(ctx context.Context, host string, p *port.Port) {
    ```

- [ ] yamlをparseしてhttpを実行 scenerio goを参考に
- [ ] parseしたyamlの値を使用して、httpリクエストを送る
- [ ] 監視して収集したlogをdbに保存、postgresql
- [ ] ターゲットのサービスのコンテキストをGPTやcurlなどの文字列から取得
- [ ] targetのサービスの特性を知る by vision
- [ ] サービスの特性と文字コンテキストから可能性のあるCVEを提示

## automatd metasploit
- [ ] pentest gptを使った、対話式ハッキング
    - `pentestgpt/utils/API` 周りは参考になりそう


## vision
- [x] vision, learning big logs
指定したlogを使用して学習
- [ ] sense, 取得したログとvisionを使って、アプリケーションの特性を理解
- [ ] アプリケーションの特性と近いexploit方法やCVEを提示