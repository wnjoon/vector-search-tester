package main

func main() {
}

// 	// 임베딩을 요청할 텍스트
// 	text := "today's weather is wonderful"

// 	// 1. 요청 데이터 생성
// 	requestData := model.EmbeddingRequest{
// 		Text: text,
// 	}

// 	// 2. Go 구조체를 JSON으로 변환
// 	jsonData, err := json.Marshal(requestData)
// 	if err != nil {
// 		log.Fatalf("Error marshaling JSON: %v", err)
// 	}

// 	// 3. Python API 서버에 POST 요청 보내기
// 	resp, err := http.Post("http://localhost:6600/embed", "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		log.Fatalf("Error making HTTP request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		log.Fatalf("API server returned non-200 status: %d %s", resp.StatusCode, string(body))
// 	}

// 	// 4. 응답 본문을 읽어서 Go 구조체로 변환
// 	var embeddingResp EmbeddingResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
// 		log.Fatalf("Error decoding JSON response: %v", err)
// 	}

// 	// 5. 결과 확인
// 	fmt.Printf("Original Text: %s\n", embeddingResp.Text)
// 	// 벡터 전체를 출력하면 너무 기므로, 앞 5개 요소와 벡터 차원만 출력
// 	fmt.Printf("Embedding Vector (first 5 dims): %v\n", embeddingResp.Embedding[:5])
// 	fmt.Printf("Vector Dimension: %d\n", len(embeddingResp.Embedding))
// }
