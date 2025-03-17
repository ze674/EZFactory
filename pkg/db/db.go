package db

import (
	"Factory/pkg/models"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite3", "./factory.db")
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            gtin TEXT NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	log.Println("Таблица products успешно создана или уже существует!")
	return nil
}

func GetProducts() ([]models.Product, error) {
	query := "SELECT id,name,gtin FROM products"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.GTIN)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func AddProduct(p models.Product) error {
	query := "INSERT INTO products (name,gtin) VALUES (?,?)"

	result, err := DB.Exec(query, p.Name, p.GTIN)
	if err != nil {
		return err
	}

	//Проверка успешности добавления
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil

}

func DeleteProduct(id int) (bool, error) {
	query := "DELETE FROM products WHERE id = ?"

	result, err := DB.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func SearchProduct(entered string) ([]models.Product, error) {
	query := "SELECT id,name,gtin FROM products WHERE name LIKE ? OR gtin LIKE ?"

	rows, err := DB.Query(query, "%"+entered+"%", "%"+entered+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err = rows.Scan(&p.ID, &p.Name, &p.GTIN)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
