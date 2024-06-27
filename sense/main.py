from gensim.models import word2vec
from nltk import word_tokenize

model = word2vec.Word2Vec.load("../lovemachine/lovemachine.model")

dictionary = ['auth', 'authenticated', 'call', 'client', 'consoles', 'core', 'db', 'jobs', 'login', 'logout', 'modules', 'plugins',
'port', 'server', 'token', 'sessions', 'ssl', 'uri']

def check_related_with_word2vec(text, dictionary, model, threshold=0.7):
    tokens = word_tokenize(text.lower())
    for token in tokens:
        for word in dictionary:
            if token in model.wv and word in model.wv:
                similarity = model.wv.similarity(token, word)
                if similarity >= threshold:
                    return True
    return False

# 実際にtargetから取得したログになる
new_logs = [
    "GET /new_page.html HTTP/1.1",
    "DELETE /api/remove HTTP/1.1",
]

related_results = [check_related_with_word2vec(log, dictionary, model) for log in new_logs]
print(related_results)
