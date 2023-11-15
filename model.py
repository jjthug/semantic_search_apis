from sentence_transformers import SentenceTransformer
model = SentenceTransformer('multi-qa-MiniLM-L6-cos-v1')

def encodeToVector(doc):
    return model.encode(doc)