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
	playerName = formatStringWithNumbers(playerName)
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
var normalCharsWithNumbers = regexp.MustCompile(`[^\p{L}0-9_\s]`)
var normalChars = regexp.MustCompile(`[^\p{L}\s]`)

func formatString(text string) string {
	text = normalChars.ReplaceAllString(text,"")
	return strings.Join(strings.Fields(text)," ")
}

func formatStringWithNumbers(text string) string {
	text = normalCharsWithNumbers.ReplaceAllString(text,"")
	return strings.Join(strings.Fields(text)," ")
}

func normalizeString(text string) string {
	return formatString(strings.ToLower(text))
}

func normalizeStringWithNumbers(text string) string {
	return formatStringWithNumbers(strings.ToLower(text))
}

func aiFallbackRecognizeBarrel(message string, aiKey string, aiUrl string, aiModel string, prompt string,oldErr error) (*RecognizedBarrelItem, error) {
    if aiKey == "" || aiUrl == "" || aiModel == "" {
        return nil, oldErr
    }

    return aiBarrelRecognition(aiKey, aiUrl, aiModel, message, prompt)
}