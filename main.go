package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/martin-bernal/go-bootcamp-project/internal/controller"
	"github.com/martin-bernal/go-bootcamp-project/internal/interfaces"
	"github.com/martin-bernal/go-bootcamp-project/internal/repository"
	"github.com/martin-bernal/go-bootcamp-project/internal/service"
)

func main() {
	r := gin.Default()
	pr := repository.NewPokemonRepo("pokemon-data.csv")

	pokeClient, _ := interfaces.NewPokemonClient("https://pokeapi.co/api/v2/")
	ps := service.NewPokemonService(pr, pokeClient)
	pc := controller.NewPokemonController(ps)

	r.GET("/pokemon", pc.GetPokemons)
	r.GET("/pokemon/:id", pc.FindPokemon)
	r.GET("/get_pokemon/:pokemon_name", pc.CreatePokemon)
	r.GET("/pokemon/async", pc.GetPokemonsAsync)
	err := r.Run("localhost:8000")
	if err != nil {
		fmt.Println("Error running HTTP server")
	}
}
