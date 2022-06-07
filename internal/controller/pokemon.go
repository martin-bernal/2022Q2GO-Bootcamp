package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/martin-bernal/go-bootcamp-project/internal/service"
)

type PokemonController interface {
	CreatePokemon(c *gin.Context)
	FindPokemon(c *gin.Context)
	GetPokemons(c *gin.Context)
	GetPokemonsAsync(c *gin.Context)
}

type pokemonController struct {
	pokeService service.PokemonService
}

func NewPokemonController(pokeService service.PokemonService) PokemonController {
	return &pokemonController{pokeService}
}

func (pc *pokemonController) CreatePokemon(c *gin.Context) {
	pokeName, ok := c.Params.Get("pokemon_name")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing pokemon name parameter."})
		return
	}

	err := pc.pokeService.CreatePokemon(pokeName)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pokemon added to database."})
}

func (pc *pokemonController) FindPokemon(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter."})
		return
	}
	poke, err := pc.pokeService.FindPokemon(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, poke)
}

func (pc *pokemonController) GetPokemons(c *gin.Context) {
	pokes, err := pc.pokeService.GetPokemons()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pokes)
}

func (pc *pokemonController) GetPokemonsAsync(c *gin.Context) {
	typeParam, ok := c.GetQuery("type")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing type query parameter."})
		return
	}

	if typeParam != "odd" && typeParam != "even" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Type parameter only accepts 'odd' or 'even' values."})
		return
	}

	items, ok := c.GetQuery("items")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing items query parameter."})
		return
	}

	intItems, err := strconv.Atoi(items)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "items param must be a valid integer."})
		return
	}

	itemsPerWorkers, ok := c.GetQuery("items_per_workers")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing items_per_workers query parameter."})
		return
	}

	intItemsPerWorkers, err := strconv.Atoi(itemsPerWorkers)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "items_per_workers param must be a valid integer."})
		return
	}
	pokes, err := pc.pokeService.GetPokemonsAsync(typeParam, intItems, intItemsPerWorkers)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get pokemons async.", "data": pokes})
}
