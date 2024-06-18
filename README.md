# runevision
automated red team tools

# todo v0.0.1
## application
- [ ] web ui
    - [ ] input target url
    - [ ] pentesting lists
    - [ ] discovered threats
    - [ ] executed exploits
    - [ ] list targets logs
- [ ] what logs to use? like apache? http request? heade? request parameter? 
read red teams book again
- [ ] skipfish or some analyzing web tool
- [ ] go application database
監視して収集したlogと実行したexploitコードをdbに保存、postgresql
- [ ] nmap golang
- [ ] 収集したログを吐いて、定期的にvisionに学習させる (一旦手動)
- [ ] application api server
- [ ] slack integration

## metasploit
- [ ] launch meteasploit db
docker?
- [ ] pymetasploit3 connect to metasploit db
pymetasploit3を使用してコンテナのmetasploit dbと繋げる
- [ ] interactive pentest gpt # interactive hacking tool by web, is the current situation more accurate?

## vision
- [x] vision, learning big logs
指定したlogを使用して学習
- [x] sense, 取得したログとvisionを使って、pymetasploit3で使用するエクスプロイトコマンドを発行
- [ ] pymetasploit3を使用して、発行されたコマンドを実行
- [ ] 実行したエクスプロイトをdbに保存 (今後の学習につながる)