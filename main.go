package main

import (
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/Vpinezi/go_crud/docs" // Importa o arquivo docs gerado pelo swag
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("secret")

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var items []Item

func main() {
	router := mux.NewRouter()

	// Rota principal para a documentação Swagger
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	router.HandleFunc("/items", GetAllItems).Methods("GET")
	router.HandleFunc("/items/{id}", GetItem).Methods("GET")
	router.HandleFunc("/items", CreateItem).Methods("POST")
	router.HandleFunc("/items/{id}", UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", DeleteItem).Methods("DELETE")
	router.HandleFunc("/login", Login).Methods("POST")

	http.ListenAndServe(":8000", router)
}

// GetAllItems retorna uma lista de todos os itens.
// @Summary Lista todos os itens
// @Description Retorna uma lista de todos os itens
// @Produce json
// @Success 200 {array} Item
// @Router /items [get]
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

// GetItem retorna um item com o ID especificado.
// @Summary Retorna um item
// @Description Retorna um item com o ID especificado
// @Produce json
// @Param id path string true "ID do item"
// @Success 200 {object} Item
// @Router /items/{id} [get]
func GetItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range items {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Item{})
}

// CreateItem cria um novo item.
// @Summary Cria um novo item
// @Description Cria um novo item
// @Accept json
// @Produce json
// @Param item body Item true "Novo item a ser criado"
// @Success 200 {object} Item
// @Router /items [post]
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	items = append(items, item)
	json.NewEncoder(w).Encode(item)
}

// UpdateItem atualiza um item existente com o ID especificado.
// @Summary Atualiza um item existente
// @Description Atualiza um item existente com o ID especificado
// @Accept json
// @Produce json
// @Param id path string true "ID do item"
// @Param item body Item true "Dados do item a serem atualizados"
// @Success 200 {object} Item
// @Router /items/{id} [put]
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range items {
		if item.ID == params["id"] {
			items = append(items[:index], items[index+1:]...)
			var newItem Item
			_ = json.NewDecoder(r.Body).Decode(&newItem)
			newItem.ID = params["id"]
			items = append(items, newItem)
			json.NewEncoder(w).Encode(newItem)
			return
		}
	}
	json.NewEncoder(w).Encode(items)
}

// DeleteItem exclui um item com o ID especificado.
// @Summary Exclui um item
// @Description Exclui um item com o ID especificado
// @Param id path string true "ID do item"
// @Success 200 {string} string "Item excluído com sucesso"
// @Router /items/{id} [delete]
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range items {
		if item.ID == params["id"] {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(items)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Mock user. In a real-world scenario, you would authenticate against a database.
	username := "user"
	password := "password"

	// Retrieve the username and password from the request body
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	_ = json.NewDecoder(r.Body).Decode(&creds)

	// Check if the provided credentials are valid
	if creds.Username == username && creds.Password == password {
		// If valid, create a JWT token
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(mySigningKey)

		// Send the token back as a response
		w.Write([]byte(tokenString))
	} else {
		// If credentials are invalid, return unauthorized status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
