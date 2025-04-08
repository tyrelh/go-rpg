package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
}

type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

func NewTilemapJSON(path string) (*TilemapJSON, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tilemap json file: %w", err)
	}
	defer jsonFile.Close()

	var tilemap TilemapJSON
	err = json.NewDecoder(jsonFile).Decode(&tilemap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tilemap json file: %w", err)
	}
	return &tilemap, nil
}
