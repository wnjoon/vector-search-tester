from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer
import numpy as np

# 1. 모델 로드 (서버 시작 시 한 번만 실행되어 메모리에 상주합니다)
print("Loading Sentence-BERT model...")
# 한국어 모델 예시 (원하는 모델로 변경 가능)
# model = SentenceTransformer('jhgan/ko-sroberta-multask') 
model = SentenceTransformer('all-MiniLM-L6-v2')
print("Model loaded successfully.")

# 2. Flask 앱 생성
app = Flask(__name__)

# 3. '/embed' 경로로 POST 요청을 처리할 API 엔드포인트 생성
@app.route('/embed', methods=['POST'])
def embed():
    # 요청 본문에서 JSON 데이터 추출
    data = request.get_json()
    if not data or 'text' not in data:
        return jsonify({"error": "text field is required"}), 400

    text_to_embed = data['text']
    
    # 텍스트를 임베딩 벡터로 변환
    embedding = model.encode(text_to_embed)
    
    # JSON으로 반환하기 위해 numpy array를 list로 변환
    response_data = {
        "text": text_to_embed,
        "embedding": embedding.tolist()
    }
    
    return jsonify(response_data)

# 4. 서버 실행
if __name__ == '__main__':
    # 외부에서 접근 가능하도록 0.0.0.0으로 설정
    app.run(host='0.0.0.0', port=6600)