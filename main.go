package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Product struct {
	Name string
	GTIN string
}

var products = []Product{
	{Name: "Шоколад", GTIN: "1234567890123"},
	{Name: "Молоко", GTIN: "9876543210987"},
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") == "true" {
			home().Render(r.Context(), w)
		} else {
			page(home()).Render(r.Context(), w)
		}
	})

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") == "true" {
			productList(products).Render(r.Context(), w)
		} else {
			page(productList(products)).Render(r.Context(), w)
		}
	})

	r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		newProduct := Product{
			Name: r.FormValue("name"),
			GTIN: r.FormValue("gtin"),
		}
		products = append(products, newProduct)
		w.Header().Set("Content-Type", "text/html")
		productItems(products).Render(r.Context(), w) // Возвращаем только <ul> для Htmx
	})

	http.ListenAndServe(":8080", r)
}
