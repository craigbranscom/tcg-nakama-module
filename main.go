package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/heroiclabs/nakama-common/runtime"
)

type CardResource struct {
	CardType string `json:"card_type"`
	Quantity int    `json:"quantity"`
}

// TODO: make inventory a map?
type PlayerInventory struct {
	PlayerId  string         `json:"player_id"`
	Inventory []CardResource `json:"inventory"`
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	err := initializer.RegisterRpc("AddCardToInventory", AddCardToInventory)
	if err != nil {
		log.Printf("Error registering rpc: %v", err)
		return err
	}
	return nil
}

func AddCardToInventory(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	//get player id from payload
	// playerId := payload["player_id"].(string)
	playerId := "test_player"

	//query db for player
	query := fmt.Sprintf("SELECT inventory FROM players WHERE player_id = '%s'", playerId)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error querying database for player inventory:", err)
		return "", err
	}
	defer rows.Close()

	//scan rows into json string
	var inventoryJSON string
	for rows.Next() {
		err := rows.Scan(&inventoryJSON)
		if err != nil {
			log.Println("Error scanning db rows:", err)
			return "", err
		}
	}

	//parse player inventory from json
	var playerInventory PlayerInventory
	err = json.Unmarshal([]byte(inventoryJSON), &playerInventory)
	if err != nil {
		log.Println("Error parsing inventory json:", err)
		return "", err
	}

	//TODO: add card to inventory
	var updatedInventoryJSON = inventoryJSON
	// updatedInventoryJSON, err = json.Marshal(playerInventory)
	// if err != nil {
	// 	log.Println("Error marshalling updated player inventory json:", err)
	// 	return "", err
	// }

	//update db for player
	updateQuery := fmt.Sprintf("UPDATE players SET inventory = '%s' WHERE player_id = '%s'", updatedInventoryJSON, playerId)
	_, err = db.ExecContext(ctx, updateQuery)
	if err != nil {
		log.Println("Error updating database for player inventory:", err)
		return "", err
	}

	return "Player Inventory Updated", nil
}
