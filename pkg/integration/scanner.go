package integration

import (
	"fmt"
	"time"
)

var (
	// DirectoryPath путь к директории с файлами интеграции
	DirectoryPath = "./integration_files"

	// ScanInterval интервал сканирования директории (в секундах)
	ScanInterval = 60

	// ScannerRunning флаг работы сканера
	ScannerRunning = false
)

// StartScanner запускает периодическое сканирование директории
func StartScanner() {
	if ScannerRunning {
		return
	}

	ScannerRunning = true

	go func() {
		for ScannerRunning {
			// Сканируем директорию
			processedFiles, err := ScanDirectory(DirectoryPath)
			if err != nil {
				fmt.Printf("Ошибка сканирования директории: %v\n", err)
			} else if len(processedFiles) > 0 {
				fmt.Printf("Обработано файлов: %d\n", len(processedFiles))
			}

			// Ждем указанный интервал
			time.Sleep(time.Duration(ScanInterval) * time.Second)
		}
	}()

	fmt.Printf("Сканер запущен. Путь: %s, интервал: %d сек.\n", DirectoryPath, ScanInterval)
}

// StopScanner останавливает сканер
func StopScanner() {
	ScannerRunning = false
	fmt.Println("Сканер остановлен.")
}

// SetScannerConfig устанавливает конфигурацию сканера
func SetScannerConfig(directory string, interval int) {
	DirectoryPath = directory
	ScanInterval = interval

	fmt.Printf("Конфигурация сканера обновлена. Путь: %s, интервал: %d сек.\n", DirectoryPath, ScanInterval)
}
