package integration

import (
	"Factory/pkg/db"
	"Factory/pkg/models"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// OutMarkFile представляет структуру XML-файла OUT_MARK
type OutMarkFile struct {
	XMLName        xml.Name     `xml:"root"`
	DocumentID     string       `xml:"document_id"`
	Source         Source       `xml:"Source"`
	Destination    Destination  `xml:"destination"`
	GTIN           string       `xml:"gtin"`
	Date           string       `xml:"data"`
	CodeDivision   CodeDivision `xml:"code_division"`
	Batch          string       `xml:"batch"`
	ExpirationDate string       `xml:"expirationdate"`
	InUseDate      string       `xml:"ineusedate"`
	Labels         Labels       `xml:"labels"`
}

type Source struct {
	OrgID  string `xml:"org_id"`
	NodeID string `xml:"node_id"`
}

type Destination struct {
	OrgID  string `xml:"org_id"`
	NodeID string `xml:"node_id"`
}

type CodeDivision struct {
	L00All  string `xml:"l_00_all"`
	L00Task string `xml:"l_00_task"`
	L01All  string `xml:"l_01_all"`
	L01Task string `xml:"l_01_task"`
	L02All  string `xml:"l_02_all"`
	L02Task string `xml:"l_02_task"`
}

type Labels struct {
	Labels []string `xml:"label"`
}

// ParseOutMarkFile разбирает XML-файл из 1С
func ParseOutMarkFile(filePath string) (*OutMarkFile, error) {
	// Открываем файл
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer xmlFile.Close()

	// Читаем содержимое файла
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	// Парсим XML
	var outMarkFile OutMarkFile
	err = xml.Unmarshal(xmlData, &outMarkFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга XML: %v", err)
	}

	return &outMarkFile, nil
}

// ProcessOutMarkFile обрабатывает файл OUT_MARK и сохраняет его в базу данных
func ProcessOutMarkFile(filePath string) (int, error) {
	// Парсим файл
	outMarkFile, err := ParseOutMarkFile(filePath)
	if err != nil {
		return 0, err
	}

	// Проверяем, не обрабатывали ли мы уже этот файл
	exists, err := db.FileExists(outMarkFile.DocumentID)
	if err != nil {
		return 0, fmt.Errorf("ошибка проверки существования файла: %v", err)
	}
	if exists {
		return 0, fmt.Errorf("файл с UUID %s уже обработан", outMarkFile.DocumentID)
	}

	// Получаем количество кодов
	codesCount, err := strconv.Atoi(outMarkFile.CodeDivision.L00Task)
	if err != nil {
		// Если не удалось преобразовать, пробуем просто подсчитать количество кодов
		codesCount = len(outMarkFile.Labels.Labels)
	}

	// Поиск продукта по GTIN
	product, err := GetProductByGTIN(outMarkFile.GTIN)
	var productID int
	if err != nil {
		// Если продукт не найден, устанавливаем ID в 0
		productID = 0
	} else {
		productID = product.ID
	}

	// Формируем путь для архивного файла
	_, fileName := filepath.Split(filePath)
	archivePath := filepath.Join("archive", fileName)

	// Создаем запись о файле интеграции
	file := models.IntegrationFile{
		UUID:         outMarkFile.DocumentID,
		Filename:     fileName,
		FilePath:     archivePath,
		GTIN:         outMarkFile.GTIN,
		ProductID:    productID,
		BatchNumber:  outMarkFile.Batch,
		Date:         outMarkFile.Date,
		CodesCount:   codesCount,
		Status:       models.FileStatusNew,
		ErrorMessage: "",
	}

	// Сохраняем файл в базу данных
	fileID, err := db.AddIntegrationFile(file)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления файла в базу данных: %v", err)
	}

	// Сохраняем коды в базу данных
	_, err = db.AddIntegrationCodes(int(fileID), outMarkFile.Labels.Labels)
	if err != nil {
		// Если не удалось сохранить коды, удаляем запись о файле и возвращаем ошибку
		// [Здесь должна быть функция удаления файла из базы данных]
		return 0, fmt.Errorf("ошибка добавления кодов в базу данных: %v", err)
	}

	// Перемещаем файл в архивную директорию
	err = moveFileToArchive(filePath, archivePath)
	if err != nil {
		// Ошибка перемещения файла не критична, просто логируем её
		fmt.Printf("Ошибка перемещения файла в архив: %v\n", err)
	}

	return int(fileID), nil
}

// ScanDirectory сканирует указанную директорию на наличие новых файлов OUT_MARK
func ScanDirectory(directory string) ([]int, error) {
	// Проверяем существование директории
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("директория не существует: %s", directory)
	}

	// Проверяем существование архивной директории, создаем если не существует
	archiveDir := filepath.Join(directory, "archive")
	_, err = os.Stat(archiveDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(archiveDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать архивную директорию: %v", err)
		}
	}

	// Получаем список файлов в директории
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения директории: %v", err)
	}

	// Массив ID обработанных файлов
	var processedFileIDs []int

	// Обрабатываем каждый файл
	for _, file := range files {
		// Пропускаем директории и файлы, не соответствующие шаблону OUT_MARK_*.xml
		if file.IsDir() || !strings.HasPrefix(file.Name(), "OUT_MARK_") || !strings.HasSuffix(file.Name(), ".xml") {
			continue
		}

		// Формируем полный путь к файлу
		filePath := filepath.Join(directory, file.Name())

		// Обрабатываем файл
		fileID, err := ProcessOutMarkFile(filePath)
		if err != nil {
			fmt.Printf("Ошибка обработки файла %s: %v\n", file.Name(), err)
			continue
		}

		// Добавляем ID файла в список обработанных
		processedFileIDs = append(processedFileIDs, fileID)
	}

	return processedFileIDs, nil
}

// GetProductByGTIN - функция для поиска продукта по GTIN
func GetProductByGTIN(gtin string) (models.Product, error) {
	return db.GetProductByGTIN(gtin)
}

// moveFileToArchive перемещает файл в архивную директорию
func moveFileToArchive(sourcePath, destPath string) error {
	// Создаем архивную директорию, если она не существует
	archiveDir := filepath.Dir(destPath)
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать архивную директорию: %v", err)
	}

	// Копируем файл
	sourceData, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать исходный файл: %v", err)
	}

	err = ioutil.WriteFile(destPath, sourceData, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать в архивный файл: %v", err)
	}

	// Удаляем исходный файл
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("не удалось удалить исходный файл: %v", err)
	}

	return nil
}
