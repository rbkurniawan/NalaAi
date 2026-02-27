package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "github.com/rbkurniawan/NalaAi/config"
    "github.com/rbkurniawan/NalaAi/prompts"
    "github.com/rbkurniawan/NalaAi/utils"
)

type ChatHandler struct {
    config    *config.Config
    prompts   *prompts.PromptManager
    logger    *utils.Logger
}

type ChatRequest struct {
    Messages []Message `json:"messages"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatResponse struct {
    Choices []struct {
        Message Message `json:"message"`
    } `json:"choices"`
}

func NewChatHandler(cfg *config.Config, pm *prompts.PromptManager, logger *utils.Logger) *ChatHandler {
    return &ChatHandler{
        config:  cfg,
        prompts: pm,
        logger:  logger,
    }
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Read request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        h.logger.Log("ERROR", fmt.Sprintf("Failed to read request body: %v", err))
        http.Error(w, "Failed to read request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Parse request
    var chatReq ChatRequest
    if err := json.Unmarshal(body, &chatReq); err != nil {
        h.logger.Log("ERROR", fmt.Sprintf("Failed to parse request: %v", err))
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }

    // Add instruction prompt to messages
    instruction := h.prompts.GetInstruction()
    messages := []Message{
        {Role: "developer", Content: instruction},
    }
    messages = append(messages, chatReq.Messages...)

    // Prepare request for Azure OpenAI
    azureReq := map[string]interface{}{
        "messages": messages,
        "model":    h.config.AzureDeployment,
    }

    azureBody, err := json.Marshal(azureReq)
    if err != nil {
        h.logger.Log("ERROR", fmt.Sprintf("Failed to marshal Azure request: %v", err))
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Call Azure OpenAI
    client := &http.Client{}
    azureURL := fmt.Sprintf("%s/chat/completions", h.config.AzureEndpoint)
    
    req, err := http.NewRequest("POST", azureURL, bytes.NewBuffer(azureBody))
    if err != nil {
        h.logger.Log("ERROR", fmt.Sprintf("Failed to create Azure request: %v", err))
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("api-key", h.config.AzureAPIKey)

    resp, err := client.Do(req)
    if err != nil {
        h.logger.Log("ERROR", fmt.Sprintf("Failed to call Azure OpenAI: %v", err))
        http.Error(w, "Failed to call Azure service", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read Azure response
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        h.logger.Log("ERROR", fmt.Sprintf("Failed to read Azure response: %v", err))
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Log request and response
    h.logger.LogRequestResponse(string(body), string(respBody))

    // Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(respBody)
}