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
from datetime import datetime, timedelta

print("Fetching CVE data...")
# CVEデータベースを取得
def fetch_cve_data():
    end_date = datetime.now()
    # todo もっと長い期間のCVEデータを取得
    start_date = end_date - timedelta(days=7)  # 過去7日間のCVEを取得
    api = (
        f"https://services.nvd.nist.gov/rest/json/cves/2.0?resultsPerPage=100&pubStartDate={start_date.isoformat()}Z&pubEndDate={end_date.isoformat()}Z"
    )
    print(api)
    response = requests.get(api)
    json_data = json.loads(response.text)

    descriptions = [item['cve']['descriptions'][0]['value'] for item in json_data['vulnerabilities']]
    ids = [item['cve']['id'] for item in json_data['vulnerabilities']]
    return descriptions, ids

print("Fetching URL content...")
# URLからデータを取得
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

print("Cleaning and processing CVE data...")
# CVEデータの取得とクリーニング
cve_descriptions, cve_ids = fetch_cve_data()

# テキストデータのベクトル化
def vectorize_texts(texts):
    vectorizer = TfidfVectorizer(stop_words='english')
    vectors = vectorizer.fit_transform(texts)
    return vectorizer, vectors

print("Vectorizing texts...")
# トピック、URLコンテンツ、CVE説明文を一度にベクトル化
all_texts = [' '.join(topic) for topic in topics] + url_contents + cve_descriptions
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

print("Splitting data into training and test sets...")
# データをトレーニングセットとテストセットに分割
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.5, random_state=42)

print("Training the model...")

# モデルの学習
model = RandomForestClassifier(n_estimators=100, random_state=42)
model.fit(X_train, y_train)

print("Evaluating the model...")
# モデルの評価
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
related_cves = predict_cve_for_url(test_url, model, vectorizer, topics, cve_vectors, cve_ids)
print(f"related cves for {test_url}: {related_cves}")
