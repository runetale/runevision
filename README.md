# runevision
automated red team tools

# todo
- [] what logs to use? like apache? http request? heade? request parameter? 
- [] launch meteasploit db
- [] pymetasploit3 connect to metasploit db
- [] go application database
- [] skipfish or some analyzing web tool
- [x] vision, learning big logs
- [] nmap golang
- [x] sense, 取得したログとvisionを使って、pymetasploit3で使用するエクスプロイトコマンドを発行
- [] automate pymetasploit3 by analyzing result(nltk token & nmap) 
- [] 実行したエクスプロイトをdbに保存、今後の学習につながる
----- in here automated red team -----
- [] interactive pentest gpt # interactive hacking tool by web, is the current situation more accurate?
- [] application api server
- [] web ui
- [] slack integration

# a.i
after todo
- skipfish or some analyzing web tool
- got nltk token by analyzing log

今後重要なのはこの2点、スキャンした対象がどのようなサービスの性質でどのようなハッキング被害が多いかを適切なフォーマットのlogファイルで取得し、vision modelを強化学習する。