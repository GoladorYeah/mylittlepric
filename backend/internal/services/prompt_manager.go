package services

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type PromptManager struct {
	prompts map[string]string
	mu      sync.RWMutex
}

func NewPromptManager() *PromptManager {
	pm := &PromptManager{
		prompts: make(map[string]string),
	}
	pm.loadPrompts()
	return pm
}

func (pm *PromptManager) loadPrompts() {
	promptFiles := map[string]string{
		"master":                  "internal/services/prompts/master_prompt.txt",
		"specialized_electronics": "internal/services/prompts/specialized_electronics.txt",
		"specialized_parametric":  "internal/services/prompts/specialized_parametric.txt",
		"specialized_generic":     "internal/services/prompts/specialized_generic_model.txt",
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	for key, path := range promptFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("⚠️  Failed to load prompt %s: %v\n", key, err)
			continue
		}
		pm.prompts[key] = string(content)
	}

	fmt.Printf("✅ Loaded %d prompts\n", len(pm.prompts))
}

func (pm *PromptManager) GetPrompt(key, country, language, currency, category string) string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	prompt, exists := pm.prompts[key]
	if !exists {
		return ""
	}

	languageName := getLanguageName(language)

	prompt = strings.ReplaceAll(prompt, "{country}", country)
	prompt = strings.ReplaceAll(prompt, "{language}", languageName)
	prompt = strings.ReplaceAll(prompt, "{currency}", currency)
	prompt = strings.ReplaceAll(prompt, "{category}", category)

	return prompt
}
