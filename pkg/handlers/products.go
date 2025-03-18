package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"Factory/pkg/db"
	"Factory/pkg/models"
	"Factory/templates"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := db.GetProducts()
	if err != nil {
		http.Error(w, "Ошибка получения продуктов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Header.Get("HX-Request") == "true" {
		templates.ProductList(products).Render(r.Context(), w)
	} else {
		templates.Page(templates.ProductList(products)).Render(r.Context(), w)
	}
}

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p := models.Product{Name: r.FormValue("name"), GTIN: r.FormValue("gtin")}
	err := db.AddProduct(p)
	if err != nil {
		http.Error(w, "Ошибка добавления продукта: "+err.Error(), http.StatusInternalServerError)
		return
	}
	products, err := db.GetProducts()
	if err != nil {
		http.Error(w, "Ошибка получения продуктов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ProductItems(products).Render(r.Context(), w)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	deleted, err := db.DeleteProduct(id)
	if err != nil {
		http.Error(w, "Ошибка удаления продукта: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Продукт не найден", http.StatusNotFound)
		return
	}
	products, err := db.GetProducts()
	if err != nil {
		http.Error(w, "Ошибка получения продуктов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ProductItems(products).Render(r.Context(), w)
}

func AddFormHandler(w http.ResponseWriter, r *http.Request) {
	templates.AddProductForm().Render(r.Context(), w)
}

func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}
