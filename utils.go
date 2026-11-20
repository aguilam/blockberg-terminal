package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func getUserMinecraftUUID(playerName string) (string,error) {
	var target struct {
		ID string `json:"id"`
	}
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.minecraftservices.com/minecraft/profile/lookup/name/%s", playerName)
	resp, err := client.Get(url)
	if err != nil{

		return "",err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		return "", fmt.Errorf("api error: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return "", err
	}
	return target.ID, nil
}

func nilCordFlag(intFlag *int) *int {
	if intFlag == nil || *intFlag == unsetCoord{
		return nil
	}
	return intFlag
}

func normalizeString(text string) string {
	newString := strings.ToLower(strings.TrimSpace(text))
	re := regexp.MustCompile(`[a-zA-Z0-9_\s]`)
	newString = re.ReplaceAllString(newString,"")
	re = regexp.MustCompile(`\s{2,}`)
	return re.ReplaceAllString(newString," ")
}

func aiFallbackRecognizeBarrel(message string, aiKey string, aiUrl string, aiModel string, prompt string,oldErr error) (*RecognizedBarrelItem, error) {
    if aiKey == "" || aiUrl == "" || aiModel == "" {
        return nil, oldErr
    }

    return aiBarrelRecognition(aiKey, aiUrl, aiModel, message, prompt)
}