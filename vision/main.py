from gensim.models import word2vec
from nltk import word_tokenize
import os
import nltk
nltk.download('punkt')

def get_all_log_files(directory):
    log_files = []
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".log"):
                log_files.append(os.path.join(root, file))
    return log_files

directory_path = "./loghub"
log_files = get_all_log_files(directory_path)

logs = []

for log_file in log_files:
    with open(log_file, 'r', encoding='utf-8', errors='ignore') as f:
        print(log_file)
        logs.extend(f.readlines())

tokenized_logs = [word_tokenize(log.lower()) for log in logs]

model = word2vec.Word2Vec(tokenized_logs, vector_size=100, window=5, min_count=1, workers=4)

model.save("vision.model")


load_model = word2vec.Word2Vec.load("vision.model")

print("model details:", load_model)