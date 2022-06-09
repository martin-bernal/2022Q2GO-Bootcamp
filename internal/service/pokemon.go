package service

import (
	"errors"
	"github.com/martin-bernal/go-bootcamp-project/internal/interfaces"
	"strconv"
	"sync"

	"github.com/martin-bernal/go-bootcamp-project/internal/entity"
	"github.com/martin-bernal/go-bootcamp-project/internal/repository"
)

//LAYER DEFINITION
type PokemonService interface {
	// CreatePokemon Function to create new pokemon
	CreatePokemon(pokeName string) error
	// FindPokemon Function to find pokemon
	FindPokemon(pokeId string) (entity.Pokemon, error)
	// GetPokemons Function to get all pokemons
	GetPokemons() ([]entity.Pokemon, error)
	//	GetPokemonsAsync Function to get all pokemons using worker pool
	GetPokemonsAsync(typeParam string, items int, itemsPerWorkers int) ([]*entity.Pokemon, error)
}

// pokemonService add the necesary fields for this file,
// repo reffers to the Pokemonrepository the file uses
// pokemonClient reffers to the Pokemon Client the file uses
type pokemonService struct {
	repo          repository.PokemonRepo
	pokemonClient interfaces.Client
}

// pokeTask struct to save relevant info of jobs sended to the worker pool
type pokeTask struct {
	Pokemon *entity.Pokemon
	TypeID  string
	TaskRun func(pokemon *entity.Pokemon, typeId string)
}

func (t pokeTask) Run() {
	t.TaskRun(t.Pokemon, t.TypeID)
}

// NewPokemonService returns the service interface
func NewPokemonService(repo repository.PokemonRepo, pokeClient interfaces.Client) PokemonService {
	return &pokemonService{repo: repo, pokemonClient: pokeClient}
}

//IMPLEMENTATIONS

func (ps *pokemonService) CreatePokemon(pokeName string) error {

	pokemon, _ := ps.pokemonClient.GetPokemonFromApi(pokeName)

	err := ps.repo.WritePokemon(pokemon)

	if err != nil {
		return err
	}
	return nil
}

func (ps *pokemonService) FindPokemon(pokeId string) (entity.Pokemon, error) {
	pokemons, err := ps.repo.ReadPokemons()
	var pokemon entity.Pokemon

	pokeIntId, err2 := strconv.ParseInt(pokeId, 10, 64)
	pokeFound := false

	if err2 != nil {
		return pokemon, errors.New("No valid int")
	}

	if err != nil {
		return entity.Pokemon{}, err
	}

	for _, poke := range pokemons {
		if poke.ID == pokeIntId {
			pokemon = poke
			pokeFound = true
		}
	}

	if !pokeFound {
		return pokemon, errors.New("Pokemon not found")
	}

	return pokemon, nil
}

func (ps *pokemonService) GetPokemons() ([]entity.Pokemon, error) {
	pokemons, err := ps.repo.ReadPokemons()
	if err != nil {
		return nil, err
	}
	return pokemons, nil
}

func (ps *pokemonService) GetPokemonsAsync(typeParam string, items int, itemsPerWorkers int) ([]*entity.Pokemon, error) {
	pokemons, err := ps.repo.ReadPokemons()
	if err != nil {
		return nil, err
	}

	pool := NewGoroutinePool(3, itemsPerWorkers)

	var res []*entity.Pokemon

	wg := &sync.WaitGroup{}

	taskFn := func(pokemon *entity.Pokemon, typeId string) {
		defer wg.Done()

		switch {
		case typeId == "odd":
			if pokemon.ID%2 != 0 {
				res = append(res, pokemon)
			}
		case typeId == "even":
			if pokemon.ID%2 == 0 {
				res = append(res, pokemon)
			}
		}
	}

	var tasks []pokeTask
	for index := range pokemons {
		if index < items {
			tasks = append(tasks, pokeTask{
				Pokemon: &pokemons[index],
				TypeID:  typeParam,
				TaskRun: taskFn,
			})
		}
	}

	for _, task := range tasks {
		wg.Add(1)
		pool.ScheduleWork(task)
	}
	pool.Close()

	wg.Wait()

	return res, nil
}
