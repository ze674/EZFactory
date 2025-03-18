package handlers

import (
	"net/http"

	"Factory/pkg/db"
	"Factory/templates"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	if query == "" {
		products, err := db.GetProducts()
		if err != nil {
			http.Error(w, "Ошибка получения продуктов: "+err.Error(), http.StatusInternalServerError)
			return
		}
		templates.ProductItems(products).Render(r.Context(), w)
		return
	}

	products, err := db.SearchProduct(query)
	if err != nil {
		http.Error(w, "Ошибка поиска продуктов: "+err.Error(), http.StatusInternalServerError)
		return
	}
	templates.ProductItems(products).Render(r.Context(), w)
}
