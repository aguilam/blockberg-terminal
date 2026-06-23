package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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