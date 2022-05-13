package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Pokemon struct {
	ID     string `json:"id"`
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
			pokemons = append(pokemons, Pokemon{
				ID:     line[0],
				Name:   line[1],
				Type:   line[2],
				Status: line[3] != "0",
			})
		}
	}

	return pokemons, nil
}
