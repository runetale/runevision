import json
import pandas as pd

def preprocess_data(json_file_path):
    """
    JSON ファイルを読み込み、特徴量とラベルに変換する関数

    Parameters:
    json_file_path (str): 入力 JSON ファイルのパス

    Returns:
    pd.DataFrame: 特徴量とラベルを含む DataFrame
    """
    with open(json_file_path, 'r', encoding='utf-8') as file:
        data = json.load(file)

    rows = []
    for entry in data:
        row = {
            'url': entry.get('url'),
            'domain': entry.get('domain'),
            'status_code': entry.get('status_code'),
            'server': entry.get('server'),
            'content_type': entry.get('content_type'),
            'port_http': entry.get('port_http'),
            'port_https': entry.get('port_https'),
            'label': entry.get('label'),
            'class': entry.get('class')
        }

        # tech_stack を特徴量に変換
        tech_stack = entry.get('tech_stack', {})
        for key, values in tech_stack.items():
            for value in values:
                row[f'tech_stack_{key}_{value}'] = 1

        # headers を特徴量に変換
        headers = entry.get('headers', {})
        for key, value in headers.items():
            row[f'header_{key}'] = value

        rows.append(row)

    df = pd.DataFrame(rows)
    return df

# 使用例
if __name__ == "__main__":
    json_file_path = 'data2.json'  # JSON ファイルのパス
    df = preprocess_data(json_file_path)
    df.to_csv('processed_data.csv', index=False)  # 結果を CSV ファイルとして保存
