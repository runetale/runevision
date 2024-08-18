import json

def update_labels_from_file(data_file_path, labels_file_path):
    """
    data.json を指定されたラベルとクラスの情報で更新する関数

    Parameters:
    data_file_path (str): 更新する data.json のファイルパス
    labels_file_path (str): ラベルとクラスの情報が含まれる JSON ファイルのパス
    """
    # ラベルとクラスの情報を読み込む
    with open(labels_file_path, 'r', encoding='utf-8') as file:
        labels_data = json.load(file)
    
    # data.json を読み込む
    with open(data_file_path, 'r', encoding='utf-8') as file:
        data = json.load(file)
    
    # ラベルとクラスの情報を URL に基づいてデータを更新
    url_label_class_map = {item['url']: {'label': item['label'], 'class': item['class']} for item in labels_data}

    for entry in data:
        url = entry.get('url')
        if url in url_label_class_map:
            entry['label'] = url_label_class_map[url]['label']
            entry['class'] = url_label_class_map[url]['class']
    
    # 更新されたデータを data.json に書き込む
    with open(data_file_path, 'w', encoding='utf-8') as file:
        json.dump(data, file, ensure_ascii=False, indent=4)

# 使用例
if __name__ == "__main__":
    # data.json のパス
    data_file_path = 'data2.json'
    
    # ラベルとクラスの情報が含まれる JSON ファイルのパス
    labels_file_path = 'domains.json'
    
    # 関数を呼び出して更新
    update_labels_from_file(data_file_path, labels_file_path)
