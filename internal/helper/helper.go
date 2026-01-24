package helper

import (
	"encoding/json"
	"log"
)

func TransformStruct(src, dest any) error {
	jsonData, err := json.Marshal(src)
	if err != nil {
		log.Default().Println("Failed to marshal struct:", err)
		return err
	}
	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		log.Default().Println("Failed to unmarshal struct:", err)
		return err
	}
	return nil
}
