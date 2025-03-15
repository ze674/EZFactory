package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	ID   int
	Name string
	GTIN string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./factory.db")
	if err != nil {
		log.Fatal("Ошибка подключения к БД ", err)
	}
	defer db.Close()

	//Создание таблицы если её нет
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT NOT NULL,
	    gtin TEXT NOT NULL 
	)
	`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы ", err)
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") == "true" {
			home().Render(r.Context(), w)
		} else {
			page(home()).Render(r.Context(), w)
		}
	})

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		products, err := getProducts()
		if err != nil {
			http.Error(w, "Ошибка получения списка продуктов "+err.Error(), http.StatusInternalServerError)
		}
		if r.Header.Get("HX-Request") == "true" {
			productList(products).Render(r.Context(), w)
		} else {
			page(productList(products)).Render(r.Context(), w)
		}
	})

	r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		product := Product{
			Name: r.FormValue("name"),
			GTIN: r.FormValue("gtin"),
		}

		err = addProduct(product)
		if err != nil {
			http.Error(w, "Ошибка добавления продукта "+err.Error(), http.StatusInternalServerError)
			return
		}

		products, err := getProducts()
		if err != nil {
			http.Error(w, "Ошибка получения списка продуктов "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		productItems(products).Render(r.Context(), w) // Возвращаем только <ul> для Htmx
	})

	r.Delete("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		if id >= 0 && id < len(products) {
			products = append(products[:id], products[id+1:]...)
		}
		productItems(products).Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", r)
}

func getProducts() ([]Product, error) {
	query := "SELECT id,name,gtin FROM products"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.GTIN)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func addProduct(p Product) error {
	query := "INSERT INTO products (name,gtin) VALUES (?,?)"

	result, err := db.Exec(query, p.Name, p.GTIN)
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
