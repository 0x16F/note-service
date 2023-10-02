package fsscanner

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindDirectory(startDir, name string) (string, error) {
	currentDir := startDir
	for {
		// Проверяем, есть ли директория с именем name в текущей директории
		path := filepath.Join(currentDir, name)
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			// Директория найдена
			return path, nil
		}

		// Если достигли корневой директории, останавливаем поиск
		if currentDir == "/" || len(currentDir) == 0 {
			break
		}

		// Переходим на уровень вверх
		currentDir = filepath.Dir(currentDir)
	}

	return "", fmt.Errorf("directory %s not found", name)
}
