package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Response  string `json:"response"`
	Timestamp string `json:"timestamp"`
	Error     string `json:"error,omitempty"`
}

type BedrockRequest struct {
	Messages []BedrockMessage `json:"messages"`
}

type BedrockMessage struct {
	Role    string           `json:"role"`
	Content []ContentMessage `json:"content"`
}

type ContentMessage struct {
	Text string `json:"text"`
}

type BedrockResponse struct {
	Output struct {
		Message struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"message"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"inputTokens"`
		OutputTokens int `json:"outputTokens"`
	} `json:"usage"`
}

var startTime = time.Now()

func main() {
	// 로거 설정
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Bearer Token 확인
	bearerToken := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
	if bearerToken == "" {
		logrus.Warn("AWS_BEARER_TOKEN_BEDROCK not found")
	} else {
		logrus.Info("Bedrock Bearer Token loaded successfully")
	}

	// 포트 설정
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 라우터 설정
	r := mux.NewRouter()
	r.HandleFunc("/", chatPageHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/chat", chatHandler).Methods("POST")
	r.Use(loggingMiddleware)

	logrus.WithFields(logrus.Fields{
		"port":    port,
		"version": "1.0.0",
	}).Info("🚀 Server starting with Bedrock Claude 4.5")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start),
			"ip":       r.RemoteAddr,
		}).Info("Request processed")
	})
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Claude 4.5 AI Chat</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { 
            max-width: 900px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 20px; 
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
            height: 90vh;
            display: flex;
            flex-direction: column;
        }
        .header {
            background: linear-gradient(135deg, #ff6b6b, #ee5a24);
            color: white;
            padding: 25px;
            text-align: center;
        }
        .header h1 { font-size: 2em; margin-bottom: 10px; }
        .header p { opacity: 0.9; font-size: 1.1em; }
        
        .chat-container {
            flex: 1;
            overflow-y: auto;
            padding: 20px;
            background: #f8f9fa;
            display: flex;
            flex-direction: column;
            gap: 15px;
        }
        .message {
            padding: 15px 20px;
            border-radius: 20px;
            max-width: 75%;
            word-wrap: break-word;
            line-height: 1.5;
        }
        .user-message {
            background: linear-gradient(135deg, #007bff, #0056b3);
            color: white;
            margin-left: auto;
            border-bottom-right-radius: 5px;
        }
        .ai-message {
            background: white;
            color: #333;
            margin-right: auto;
            border: 1px solid #e9ecef;
            border-bottom-left-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        .input-container {
            padding: 25px;
            background: white;
            border-top: 1px solid #dee2e6;
        }
        .input-group {
            display: flex;
            gap: 15px;
            align-items: center;
        }
        #messageInput {
            flex: 1;
            padding: 15px 20px;
            border: 2px solid #dee2e6;
            border-radius: 25px;
            font-size: 16px;
            outline: none;
            transition: border-color 0.3s;
        }
        #messageInput:focus {
            border-color: #007bff;
        }
        #sendButton {
            padding: 15px 30px;
            background: linear-gradient(135deg, #28a745, #20c997);
            color: white;
            border: none;
            border-radius: 25px;
            cursor: pointer;
            font-size: 16px;
            font-weight: bold;
            transition: all 0.3s;
            min-width: 100px;
        }
        #sendButton:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(40, 167, 69, 0.4);
        }
        #sendButton:disabled {
            background: #6c757d;
            cursor: not-allowed;
            transform: none;
            box-shadow: none;
        }
        .loading {
            display: none;
            text-align: center;
            padding: 15px;
            color: #6c757d;
            font-style: italic;
        }
        .error {
            background: #f8d7da;
            color: #721c24;
            padding: 15px;
            border-radius: 10px;
            margin: 10px 0;
            border-left: 4px solid #dc3545;
        }
        .typing {
            display: none;
            padding: 15px 20px;
            background: white;
            border-radius: 20px;
            margin-right: auto;
            max-width: 75%;
            border: 1px solid #e9ecef;
            color: #6c757d;
        }
        .typing::after {
            content: '...';
            animation: typing 1.5s infinite;
        }
        @keyframes typing {
            0%, 60% { content: '...'; }
            30% { content: '..'; }
            90% { content: '.'; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🤖 Claude 4.5 AI Chat</h1>
            <p>Amazon Bedrock으로 구동되는 AI 어시스턴트</p>
        </div>
        
        <div class="chat-container" id="chatContainer">
            <div class="message ai-message">
                👋 안녕하세요! Claude 4.5입니다. 궁금한 것이 있으시면 언제든 물어보세요!
            </div>
        </div>
        
        <div class="typing" id="typing">
            Claude가 답변을 생각하고 있습니다
        </div>
        
        <div class="input-container">
            <div class="input-group">
                <input type="text" id="messageInput" placeholder="Claude 4.5에게 질문해보세요..." maxlength="2000">
                <button id="sendButton">전송</button>
            </div>
        </div>
    </div>

    <script>
        const chatContainer = document.getElementById('chatContainer');
        const messageInput = document.getElementById('messageInput');
        const sendButton = document.getElementById('sendButton');
        const typing = document.getElementById('typing');

        function addMessage(message, isUser = false) {
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message ' + (isUser ? 'user-message' : 'ai-message');
            messageDiv.innerHTML = message.replace(/\n/g, '<br>');
            chatContainer.appendChild(messageDiv);
            chatContainer.scrollTop = chatContainer.scrollHeight;
        }

        function showError(message) {
            const errorDiv = document.createElement('div');
            errorDiv.className = 'error';
            errorDiv.innerHTML = '❌ ' + message;
            chatContainer.appendChild(errorDiv);
            chatContainer.scrollTop = chatContainer.scrollHeight;
        }

        function showTyping() {
            typing.style.display = 'block';
            chatContainer.scrollTop = chatContainer.scrollHeight;
        }

        function hideTyping() {
            typing.style.display = 'none';
        }

        async function sendMessage() {
            const message = messageInput.value.trim();
            if (!message) return;

            // 사용자 메시지 추가
            addMessage(message, true);
            messageInput.value = '';
            
            // UI 상태 변경
            sendButton.disabled = true;
            sendButton.textContent = '전송 중...';
            showTyping();

            try {
                const response = await fetch('/api/chat', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ message: message })
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error('서버 오류 (' + response.status + '): ' + errorText);
                }

                const data = await response.json();
                
                if (data.error) {
                    showError(data.error);
                } else {
                    addMessage(data.response);
                }
            } catch (error) {
                showError('채팅 중 오류가 발생했습니다: ' + error.message);
            } finally {
                sendButton.disabled = false;
                sendButton.textContent = '전송';
                hideTyping();
            }
        }

        sendButton.addEventListener('click', sendMessage);
        messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                sendMessage();
            }
        });

        // 초기 포커스
        messageInput.focus();

        // 환영 메시지 애니메이션
        setTimeout(() => {
            addMessage('💡 팁: Shift+Enter로 줄바꿈, Enter로 전송할 수 있습니다.');
        }, 1000);
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.WithError(err).Error("Failed to decode request")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Bearer Token 확인
	bearerToken := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
	if bearerToken == "" {
		logrus.Error("AWS_BEARER_TOKEN_BEDROCK not found")
		response := ChatResponse{
			Error:     "Bedrock 인증 토큰이 설정되지 않았습니다",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Bedrock API 요청 준비
	bedrockReq := BedrockRequest{
		Messages: []BedrockMessage{
			{
				Role: "user",
				Content: []ContentMessage{
					{
						Text: req.Message,
					},
				},
			},
		},
	}

	reqBody, err := json.Marshal(bedrockReq)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal request")
		response := ChatResponse{
			Error:     "요청 처리 중 오류가 발생했습니다",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Bedrock API 호출 (Claude 3 Haiku 사용)
	bedrockURL := "https://bedrock-runtime.us-east-1.amazonaws.com/model/anthropic.claude-3-haiku-20240307-v1:0/converse"
	
	httpReq, err := http.NewRequest("POST", bedrockURL, bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.WithError(err).Error("Failed to create HTTP request")
		response := ChatResponse{
			Error:     "API 요청 생성 실패",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 헤더 설정
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+bearerToken)

	// HTTP 클라이언트로 요청
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		logrus.WithError(err).Error("Failed to call Bedrock API")
		response := ChatResponse{
			Error:     "AI 서비스에 연결할 수 없습니다",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"status": resp.StatusCode,
			"body":   string(bodyBytes),
		}).Error("Bedrock API error")
		
		response := ChatResponse{
			Error:     fmt.Sprintf("AI 서비스 오류 (상태: %d)", resp.StatusCode),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 응답 파싱
	var bedrockResp BedrockResponse
	if err := json.NewDecoder(resp.Body).Decode(&bedrockResp); err != nil {
		logrus.WithError(err).Error("Failed to decode Bedrock response")
		response := ChatResponse{
			Error:     "AI 응답 처리 중 오류가 발생했습니다",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 응답 텍스트 추출
	var responseText string
	if len(bedrockResp.Output.Message.Content) > 0 {
		responseText = bedrockResp.Output.Message.Content[0].Text
	} else {
		responseText = "죄송합니다. 응답을 생성할 수 없습니다."
	}

	logrus.WithFields(logrus.Fields{
		"input_tokens":  bedrockResp.Usage.InputTokens,
		"output_tokens": bedrockResp.Usage.OutputTokens,
	}).Info("Bedrock API call successful")

	response := ChatResponse{
		Response:  strings.TrimSpace(responseText),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	bearerToken := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
	
	health := map[string]interface{}{
		"status":       "healthy",
		"message":      "Claude 4.5 AI Chat Server",
		"timestamp":    time.Now(),
		"version":      "1.0.0",
		"uptime":       time.Since(startTime).Round(time.Second),
		"bedrock_auth": bearerToken != "",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}