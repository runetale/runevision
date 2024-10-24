import requests
from bs4 import BeautifulSoup
import socket
import builtwith
import json

# JSONデータをファイルから読み込み
with open('domains.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

def get_domain_info(url):
    try:
        # ドメイン名を抽出
        domain = url.split("//")[-1].split("/")[0]
        
        # ポート番号を確認 (HTTPポート80, HTTPSポート443)
        port_http = socket.getservbyname('http', 'tcp')
        port_https = socket.getservbyname('https', 'tcp')
        
        # ドメインで使用されている技術スタックを検出
        tech_stack = builtwith.builtwith(url)
        
        # ドメインにGETリクエストを送信
        response = requests.get(url)
        
        # HTTPステータスコードを取得
        status_code = response.status_code
        
        # HTTPレスポンスヘッダーを取得
        headers = response.headers
        
        # サーバーの種類を取得
        server = headers.get('Server', 'Unknown')
        
        # コンテンツタイプを取得
        content_type = headers.get('Content-Type', 'Unknown')
        
        # レスポンスボディの解析
        content = response.content
        if 'html' in content_type:
            # HTMLコンテンツの解析
            soup = BeautifulSoup(content, 'html.parser')
            title = soup.title.string if soup.title else 'No Title'
            content_summary = {
                'type': 'WebサイトまたはWebアプリケーション',
                'title': title
            }
        elif 'json' in content_type:
            # JSONコンテンツの解析
            json_content = response.json()
            content_summary = {
                'type': 'APIサーバー',
                'json_keys': list(json_content.keys())
            }
        else:
            content_summary = {
                'type': '不明なコンテンツタイプ'
            }
        
        # 結果をまとめる
        info = {
            'url': url,
            'domain': domain,
            'status_code': status_code,
            'server': server,
            'content_type': content_type,
            'headers': dict(headers),
            'port_http': port_http,
            'port_https': port_https,
            'tech_stack': tech_stack,
            'content_summary': content_summary
        }
        return info

    except requests.exceptions.RequestException as e:
        return {'error': str(e)}

# 使用例
output = []
for d in data:
    domain_info = get_domain_info(d['url'])
    print(domain_info)
    output.append(domain_info)

# JSONデータをファイルに保存
with open('data.json', 'w', encoding='utf-8') as f:
    json.dump(output, f, ensure_ascii=False, indent=4)
