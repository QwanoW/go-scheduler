package internal

import (
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/replit/database-go"
)

// SaveSchedule saves the schedule to the Replit database
func SaveSchedule(schedule map[string]string) error {
	// Convert the schedule map to a JSON string
	jsonSchedule, err := json.Marshal(schedule)
	if err != nil {
		return fmt.Errorf("failed to marshal schedule: %w", err)
	}

	// Save the JSON string to the Replit database
	err = db.Set("schedule", string(jsonSchedule))
	if err != nil {
		return fmt.Errorf("failed to set schedule in db: %w", err)
	}

	return nil
}

// GetSchedule retrieves the schedule from the Replit database
func GetSchedule() (map[string]string, error) {
	// Retrieve the schedule as a JSON string from the Replit database
	jsonSchedule, err := db.Get("schedule")
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get schedule from db: %w", err)
	}

	// Convert the JSON string back to a map
	var schedule map[string]string
	err = json.Unmarshal([]byte(jsonSchedule), &schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule: %w", err)
	}

	return schedule, nil
}

// SaveMsgIds saves the schedule to the Replit database
func SaveMsgIds(msgIds []int) error {
	// Convert the schedule's msgIds slice to a JSON string
	jsonSchedule, err := json.Marshal(msgIds)
	if err != nil {
		return fmt.Errorf("failed to marshal schedule's msgIds: %w", err)
	}

	// Save the JSON string to the Replit database
	err = db.Set("messages", string(jsonSchedule))
	if err != nil {
		return fmt.Errorf("failed to set schedule's msgIds in db: %w", err)
	}

	return nil
}

// GetMsgIds retrieves the schedule's msgIds from the Replit database
func GetMsgIds() ([]int, error) {
	// Retrieve the schedule as a JSON string from the Replit database
	jsonMessages, err := db.Get("messages")
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get schedule's msgIds from db: %w", err)
	}

	// Convert the JSON string back to a slice
	var msgIds []int
	err = json.Unmarshal([]byte(jsonMessages), &msgIds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule's msgIds: %w", err)
	}

	return msgIds, nil
}

//// ClearDB clears all records in the Replit database
//func ClearDB() error {
//	// Delete all keys in the database
//	keys, err := db.ListKeys("")
//	if err != nil {
//		return fmt.Errorf("failed to list keys from db: %w", err)
//	}
//
//	for _, key := range keys {
//		err = db.Delete(key)
//		if err != nil {
//			return fmt.Errorf("failed to delete key '%s' from db: %w", key, err)
//		}
//	}
//
//	return nil
//}
