package files

import (
	"encoding/json"
	"fmt"
	"os"
)

func ExtractFilesBySuffix(path string, suffixes []string) []string {
	return []string{"./Dockerfile"}
}

func SaveObjectToFile(filePath string, obj interface{}) error {

	resultBytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error marshaling struct:", err)
		return err
	}

	err = os.WriteFile(filePath, resultBytes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}
