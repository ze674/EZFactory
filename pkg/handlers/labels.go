package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Factory/pkg/db"
	"Factory/templates"
	"github.com/go-chi/chi/v5"
)

func LabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	p, err := db.GetProductByID(id)
	if err != nil {
		http.Error(w, "Ошибка получения продукта: "+err.Error(), http.StatusNotFound)
		return
	}

	values := make(map[string]string)
	if p.LabelData != "" {
		var data map[string]interface{}
		err = json.Unmarshal([]byte(p.LabelData), &data)
		if err == nil {
			for _, key := range []string{"header", "label_name", "standard", "article", "unit_weight", "box_quantity", "box_weight"} {
				if val, ok := data[key]; ok {
					values[key] = fmt.Sprintf("%v", val)
				} else {
					values[key] = ""
				}
			}
		}
	}

	templates.LabelForm(p, values).Render(r.Context(), w)
}

func UpdateLabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	r.ParseForm()
	labelData := map[string]string{
		"header":       r.FormValue("header"),
		"label_name":   r.FormValue("label_name"),
		"standard":     r.FormValue("standard"),
		"article":      r.FormValue("article"),
		"unit_weight":  r.FormValue("unit_weight"),
		"box_quantity": r.FormValue("box_quantity"),
		"box_weight":   r.FormValue("box_weight"),
	}

	jsonData, err := json.Marshal(labelData)
	if err != nil {
		http.Error(w, "Ошибка сериализации данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.UpdateLabelData(id, string(jsonData))
	if err != nil {
		http.Error(w, "Ошибка обновления данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}
