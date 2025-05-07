package db

import "Factory/pkg/models"

func GetProducts() ([]models.Product, error) {
	query := "SELECT id,name,gtin,label_data FROM products"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.GTIN, &p.LabelData)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func GetProductByID(id int) (models.Product, error) {
	query := "SELECT id,name,gtin,label_data FROM products WHERE id = ?"

	row := DB.QueryRow(query, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.GTIN, &p.LabelData)
	if err != nil {
		return p, err
	}

	return p, nil
}

func AddProduct(p models.Product) error {
	query := "INSERT INTO products (name,gtin,label_data) VALUES (?,?,?)"

	result, err := DB.Exec(query, p.Name, p.GTIN, p.LabelData)
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

	products := []models.Product{}
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

func UpdateLabelData(id int, labelData string) error {
	query := "UPDATE products SET label_data = ? WHERE id = ?"
	_, err := DB.Exec(query, labelData, id)
	if err != nil {
		return err
	}

	return nil
}

// GetProductByGTIN находит продукт по GTIN
func GetProductByGTIN(gtin string) (models.Product, error) {
	query := "SELECT id, name, gtin, label_data FROM products WHERE gtin = ?"

	row := DB.QueryRow(query, gtin)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.GTIN, &p.LabelData)
	if err != nil {
		return p, err
	}

	return p, nil
}
