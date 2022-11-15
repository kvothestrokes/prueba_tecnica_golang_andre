package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	// Obtener 25 objetos
	// obtenerRegistros()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jokes := obtenerRegistros()
		json.NewEncoder(w).Encode(jokes)
	})
	http.ListenAndServe(":8080", nil)
}

// Se debe desarrollar un API RESTful que contenga un único endpoint el cual retorne un array con un rango definido de 25 objetos,
//  diferentes  unos de otros, es decir no debe existir ningún objeto que contenga el mismo id
// Los objetos serán obtenidos de un API de terceros, teniendo como restricción consumir estrictamente el siguiente
// endpoint: https://api.chucknorris.io/jokes/random
// https://api.chucknorris.io/

type jokeResponse struct {
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	IconUrl    string   `json:"icon_url"`
	Id         string   `json:"id"`
	UpdatedAt  string   `json:"updated_at"`
	Url        string   `json:"url"`
	Value      string   `json:"value"`
}

func obtenerRegistros() []jokeResponse {
	i := 1
	contador := 25
	var jokes []jokeResponse
	for i <= contador {

		currentJoke, err := getJoke()
		if err != nil {
			fmt.Printf("Errror al obtener objeto %v", err.Error())
		}

		repetido, _ := verificarRepetido(currentJoke, jokes)

		if !repetido {
			jokes = append(jokes, currentJoke)
			i++
		}
	}

	// Imprimir objetos
	for i, joke := range jokes {
		fmt.Printf("\n Joke n:%v  ID: %v", i+1, joke.Id)
	}
	return jokes
}

func verificarRepetido(appendJoke jokeResponse, jokeArray []jokeResponse) (bool, jokeResponse) {

	repetido := false
	var returnJoke jokeResponse
	for _, joke := range jokeArray {
		if joke.Id == appendJoke.Id {
			repetido = true
			returnJoke = appendJoke
			break
		}
	}

	return repetido, returnJoke
}

func getJoke() (jokeResponse, error) {
	endpoint := "https://api.chucknorris.io/jokes/random"
	var joke jokeResponse

	resp, errGet := http.Get(endpoint)
	if errGet != nil {
		fmt.Printf("Error al obtener broma %v", errGet.Error())
		return joke, errGet
	}
	defer resp.Body.Close()

	response, errParse := ioutil.ReadAll(resp.Body)
	if errParse != nil {
		fmt.Printf("Error al parsear broma %v", errParse.Error())
		return joke, errParse
	}

	json.Unmarshal(response, &joke)
	return joke, nil
}
