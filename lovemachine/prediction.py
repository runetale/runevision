import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.feature_extraction import DictVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import classification_report

# データの読み込み
df = pd.read_csv('processed_data.csv')

# ラベルとクラスの設定
label_column = 'label'
class_column = 'class'

# 特徴量とターゲットの分離
X = df.drop(columns=[label_column, class_column])
y_label = df[label_column]
y_class = df[class_column]

# データの分割
X_train, X_test, y_label_train, y_label_test, y_class_train, y_class_test = train_test_split(
    X, y_label, y_class, test_size=0.3, random_state=42
)

# 特徴量の変換
vectorizer = DictVectorizer(sparse=False)
X_train_transformed = vectorizer.fit_transform(X_train.fillna(0).to_dict(orient='records'))
X_test_transformed = vectorizer.transform(X_test.fillna(0).to_dict(orient='records'))

# ラベルのモデルのトレーニング
label_model = LogisticRegression(max_iter=1000)
label_model.fit(X_train_transformed, y_label_train)

# クラスのモデルのトレーニング
class_model = LogisticRegression(max_iter=1000)
class_model.fit(X_train_transformed, y_class_train)

# モデルの予測と評価
y_label_pred = label_model.predict(X_test_transformed)
y_class_pred = class_model.predict(X_test_transformed)

print("ラベルの分類レポート:")
print(classification_report(y_label_test, y_label_pred))

print("クラスの分類レポート:")
print(classification_report(y_class_test, y_class_pred))
