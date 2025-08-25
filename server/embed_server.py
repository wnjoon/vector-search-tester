from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer
import numpy as np

# 1. Load model: model will be loaded only once and will be stored in memory
print("Loading Sentence-BERT model...")

# 2. Choice model
# - 'jhgan/ko-sroberta-multask': Korean model
# - 'all-MiniLM-L6-v2': English model

# model = SentenceTransformer('jhgan/ko-sroberta-multask') 
en_model = SentenceTransformer('all-MiniLM-L6-v2')
ko_model = SentenceTransformer('jhgan/ko-sroberta-multitask')
print("Model loaded successfully.")

# 3. Create Flask app
app = Flask(__name__)

# 4. Create POST endpoint to handle '/embed' path
@app.route('/embed/en', methods=['POST'])
def embed_en():
    # 4.1. Extract JSON data from request body
    data = request.get_json()
    if not data or 'text' not in data:
        return jsonify({"error": "text field is required"}), 400

    text_to_embed = data['text']
    
    # 4.2. Convert text to embedding vector
    embedding = en_model.encode(text_to_embed)
    
    # 4.3. Convert numpy array to list for JSON response
    response_data = {
        "text": text_to_embed,
        "embedding": embedding.tolist()
    }
    
    return jsonify(response_data)

@app.route('/embed/ko', methods=['POST'])
def embed_ko():
    # 4.1. Extract JSON data from request body
    data = request.get_json()
    if not data or 'text' not in data:
        return jsonify({"error": "text field is required"}), 400

    text_to_embed = data['text']
    
    # 4.2. Convert text to embedding vector
    embedding = ko_model.encode(text_to_embed)
    
    # 4.3. Convert numpy array to list for JSON response
    response_data = {
        "text": text_to_embed,
        "embedding": embedding.tolist()
    }
    
    return jsonify(response_data)

# 5. Run server
if __name__ == '__main__':
    # 5.1. Open port 6600 (modifiable)
    app.run(host='0.0.0.0', port=6600)