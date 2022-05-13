package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status bool   `json:"status"`
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world")
	})

	router.GET("/pokemon", getPokemons)
	router.GET("/pokemon/:id", getPokemonById)

	router.Run("localhost:8000")
}

func getPokemons(context *gin.Context) {
	pokemons, err := getCsvData()
	if err != nil {
		fmt.Printf("Error: %T %v \n", err, err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, pokemons)
}

func getPokemonById(context *gin.Context) {
	id, ok := context.Params.Get("id")
	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter."})
		return
	}

	int_id, error := strconv.Atoi(id)

	if error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Provided id is not valid."})
		return
	}

	pokemons, err := getCsvData()
	if err != nil {
		fmt.Printf("Error: %T %v \n", err, err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for _, pokemon := range pokemons {
		if pokemon.ID == int_id {
			context.IndentedJSON(http.StatusOK, pokemon)
			return
		}
	}

	context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Pokemon not found."})

}

func getCsvData() ([]Pokemon, error) {
	var pokemons []Pokemon
	f, err := os.Open("pokemon-data.csv")
	if err != nil {
		return nil, errors.New("CSV file not found")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.New("unabled to read csv file")
	}

	for i, line := range data {
		if i > 0 {
			poke_id, err := strconv.Atoi(line[0])
			if err != nil {
				return nil, errors.New("all IDs on pokemon database must be Integers")
			}
			pokemons = append(pokemons, Pokemon{
				ID:     poke_id,
				Name:   line[1],
				Type:   line[2],
				Status: line[3] != "0",
			})
		}
	}

	return pokemons, nil
}
