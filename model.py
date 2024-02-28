from sentence_transformers import SentenceTransformer
model = SentenceTransformer('BAAI/bge-base-en-v1.5')

def encodeToVector(doc):
    return model.encode(doc)