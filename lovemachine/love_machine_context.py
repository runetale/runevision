import requests
from bs4 import BeautifulSoup
import re
import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report
import joblib
import json
import random
import time
from datetime import datetime, timedelta

# todo postgresに格納して、更新があれば追加するような仕組みにする。毎回取得はアホ
# CVEデータベースを取得
print("Fetching CVE data...")
def fetch_cve_data():
    end_date = datetime.now()
    start_date = end_date - timedelta(days=365)  # 過去1年間の開始日
    descriptions = []
    ids = []

    while start_date < end_date:
        week_end_date = start_date + timedelta(days=7)
        if week_end_date > end_date:
            week_end_date = end_date
            print(f"coming to end_date, => {week_end_date}")

        api = (
            f"https://services.nvd.nist.gov/rest/json/cves/2.0?resultsPerPage=100&pubStartDate={start_date.isoformat()}Z&pubEndDate={week_end_date.isoformat()}Z"
        )
        print(f"Fetching data from {start_date} to {week_end_date}...")

        success = False
        for _ in range(3):  # 最大3回のリトライ
            try:
                response = requests.get(api, timeout=10)
                response.raise_for_status()
                json_data = response.json()
                success = True
                break
            except requests.RequestException as e:
                print(f"Request error: {e}")
            except json.JSONDecodeError as e:
                print(f"JSON decode error: {e}")

            wait_time = random.randint(5, 10)
            print(f"Retrying in {wait_time} seconds...")
            time.sleep(wait_time)

        if not success:
            print(f"Failed to fetch data for the period {start_date} to {week_end_date}. Skipping this period.")
            start_date = week_end_date
            continue

        if 'vulnerabilities' in json_data:
            descriptions.extend([item['cve']['descriptions'][0]['value'] for item in json_data['vulnerabilities']])
            ids.extend([item['cve']['id'] for item in json_data['vulnerabilities']])

        time.sleep(random.randint(3, 8))

        start_date = week_end_date

    return descriptions, ids

# データを取得
descriptions, ids = fetch_cve_data()

# 結果の一部を表示
print(f"Total CVEs fetched: {len(ids)}")


# URLからデータを取得
print("Fetching URL content...")
def fetch_url_content(url):
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        return response.text
    except requests.RequestException as e:
        print(f"Error fetching {url}: {e}")
        return None

# HTMLタグを除去しテキストをクリーニング
def clean_html(html):
    soup = BeautifulSoup(html, 'html.parser')
    text = soup.get_text()
    text = re.sub(r'\s+', ' ', text).strip()
    return text

print("Defining topics...")

# todo もっと正確なトピックデータを使用
# トピックのデモデータ（実際には適切なデータを使用）
topics = [["security", "vulnerability", "attack"], ["encryption", "data", "privacy"], ["network", "protocol", "internet"]]

print("Fetching content from URLs...")
# URLリストの取得（例としてハードコードされたリスト）
urls = [
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
    "https://www.runetale.com/",
]

# URLコンテンツの取得とクリーニング
url_contents = [clean_html(fetch_url_content(url)) for url in urls if fetch_url_content(url)]

# CVEデータの取得とクリーニング
print("Cleaning and processing CVE data...")

# テキストデータのベクトル化
def vectorize_texts(texts):
    vectorizer = TfidfVectorizer(stop_words='english')
    vectors = vectorizer.fit_transform(texts)
    return vectorizer, vectors

# トピック、URLコンテンツ、CVE説明文を一度にベクトル化
print("Vectorizing texts...")
all_texts = [' '.join(topic) for topic in topics] + url_contents + descriptions
vectorizer, all_vectors = vectorize_texts(all_texts)

# ベクトルの分割
num_topics = len(topics)
num_urls = len(url_contents)
topic_vectors = all_vectors[:num_topics]
url_vectors = all_vectors[num_topics:num_topics + num_urls]
cve_vectors = all_vectors[num_topics + num_urls:]

# 特徴量としてコサイン類似度を計算
def calculate_cosine_similarity(topic_vectors, url_vectors):
    similarities = []
    for topic_vector in topic_vectors:
        similarity = cosine_similarity(url_vectors, topic_vector)
        similarities.append(similarity.flatten())
    return np.array(similarities).T

print("Calculating cosine similarities...")

# 特徴量としてのコサイン類似度を計算
X = calculate_cosine_similarity(topic_vectors, url_vectors)
y = np.random.randint(0, 2, size=len(url_contents))  # ここではランダムにラベルを付与

if len(url_contents) < 2:
    raise ValueError("Not enough URL contents to split into training and test sets. Please provide more URL contents.")

# データをトレーニングセットとテストセットに分割
print("Splitting data into training and test sets...")
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.5, random_state=42)

print("Training the model...")

# モデルの学習
model = RandomForestClassifier(n_estimators=100, random_state=42)
model.fit(X_train, y_train)

# モデルの評価
print("Evaluating the model...")
y_pred = model.predict(X_test)
print(classification_report(y_test, y_pred))

print("Saving the model...")

# モデルの保存
joblib.dump(model, 'cve_predictor_model.pkl')
joblib.dump(vectorizer, 'tfidf_vectorizer.pkl')

print("Preparing to predict CVEs for a specific URL...")
# 特定のURLに対する予測
def predict_cve_for_url(url, model, vectorizer, topics, cve_vectors, cve_ids):
    content = fetch_url_content(url)
    if not content:
        return []
    
    cleaned_content = clean_html(content)
    content_vector = vectorizer.transform([cleaned_content])
    topic_vectors = vectorizer.transform([' '.join(topic) for topic in topics])
    content_similarities = calculate_cosine_similarity(topic_vectors, content_vector)
    
    predictions = model.predict(content_similarities)
    related_cves = [cve_ids[i] for i, pred in enumerate(predictions) if pred == 1]
    return related_cves

# 予測モデルのロード
model = joblib.load('cve_predictor_model.pkl')
vectorizer = joblib.load('tfidf_vectorizer.pkl')

# 特定のURLに対するCVEの予測
test_url = "https://www.runetale.com/"
print(f"Predicting CVEs for {test_url}...")
related_cves = predict_cve_for_url(test_url, model, vectorizer, topics, cve_vectors, ids)
print(f"related cves for {test_url}: {related_cves}")
