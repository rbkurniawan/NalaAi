package prompts

import (
    "os"
    "strings"
)

type PromptManager struct {
    prompts map[string]string
}

func NewPromptManager() *PromptManager {
    pm := &PromptManager{
        prompts: make(map[string]string),
    }
    pm.loadPrompts()
    return pm
}

func (pm *PromptManager) loadPrompts() {
    for _, env := range os.Environ() {
        pair := strings.SplitN(env, "=", 2)
        if len(pair) != 2 {
            continue
        }
        
        key, value := pair[0], pair[1]
        if strings.Contains(key, "_PROMPT") || key == "INSTRUCTION_PROMPT" {
            pm.parsePrompt(value)
        }
    }
}

func (pm *PromptManager) parsePrompt(content string) {
    // Format: [nama prompt]content[nama prompt-end]
    startIdx := strings.Index(content, "[")
    if startIdx == -1 {
        return
    }
    
    endIdx := strings.Index(content, "]")
    if endIdx == -1 {
        return
    }
    
    promptName := content[startIdx+1 : endIdx]
    contentStart := endIdx + 1
    
    endMarker := "[" + promptName + "-end]"
    contentEnd := strings.Index(content, endMarker)
    if contentEnd == -1 {
        return
    }
    
    promptContent := strings.TrimSpace(content[contentStart:contentEnd])
    pm.prompts[promptName] = promptContent
}

func (pm *PromptManager) GetPrompt(name string) string {
    if prompt, exists := pm.prompts[name]; exists {
        return prompt
    }
    // Default ke Instruction jika tidak ditemukan
    if instruction, exists := pm.prompts["Instruction"]; exists {
        return instruction
    }
    return ""
}

func (pm *PromptManager) GetInstruction() string {
    return pm.GetPrompt("Instruction")
}