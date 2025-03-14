package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Product struct {
	Name string
	GTIN string
}

var products = []Product{}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Главная страница Factory"))
	})

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Виды продуктов"))
	})

	http.ListenAndServe(":8080", r)
}
