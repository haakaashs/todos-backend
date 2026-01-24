package helper

import (
	"encoding/json"

	// log "github.com/haakaashs/todos-backend/pkg"
	"log"
)

func TransformStruct(src, dest any) error {
	jsonData, err := json.Marshal(src)
	if err != nil {
		// log..ERROR("Failed to marshal struct:", err)
		log.Default().Println("Failed to marshal struct:", err)
		return err
	}
	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		// log.Generic.ERROR("Failed to unmarshal struct:", err)
		log.Default().Println("Failed to unmarshal struct:", err)
		return err
	}
	return nil
}
