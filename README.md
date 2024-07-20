# runevision
automated red team tools

# v0.0.1
## application
- [ ] targetのDNSスキャン sub finderをベースに
    - crtshのpostgres, hacker targer, dnsdumpstarが精度高いのでそれで実装
- [ ] targetのPortスキャン naabuをベースに
    ```go
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

- [ ] targetのtempateの攻撃はnucleiを使用, vpn攻撃をメインにする
- [ ] nucleiやスキャニングで収集したlogをdbに保存、postgresql

- [x] ターゲットのサービスのコンテキストをGPTやcurlなどの文字列から取得
- [x] targetのサービスの特性を知る LOVE MACHINE
- [x] サービスの特性と文字コンテキストから可能性のあるCVEを提示

## lovemachine-context
サービスのコンテキストから学習
- [x] lovemachine, アプリケーションの特性を理解
- [x] アプリケーションの特性と近いexploit方法やCVEを提示
.pyで実装している

## lovemachine-log
サービスのログから学習
- [ ] データベースに保存しているhttpリクエストのログを使用して学習
love_machine.pyで実装している


## api
https://runetale.postman.co/workspace/runetale~0bf06704-a345-4663-8e4b-a807be69477e/collection/35986956-704630ab-04ee-4bb2-86e9-65e12fc697fe?action=share&creator=35986956

## Next TODO
- [ ] cloudlistを使った、内部アセットはスキャン
- [ ] 独自のyamlをparseしてhttpを実行 scenerio goを参考に
- [ ] parseしたyamlの値を使用して、httpリクエストを送る
- [ ] gptのfunction callを使って、サービスの脆弱性検知をより細かく行う