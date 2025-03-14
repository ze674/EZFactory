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
	{
		Name: "Товар 1",
		GTIN: "1234567890123",
	},
	{
		Name: "Товар 2",
		GTIN: "1234567890124",
	},
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		page(home()).Render(r.Context(), w)
	})

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		page(productList(products)).Render(r.Context(), w)
	})
	r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		newProduct := Product{
			Name: r.FormValue("name"),
			GTIN: r.FormValue("gtin"),
		}

		products = append(products, newProduct)

		page(productList(products)).Render(r.Context(), w)
	})
	http.ListenAndServe(":8080", r)
}
