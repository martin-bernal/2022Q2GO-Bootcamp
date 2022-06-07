package repository

import (
	"encoding/csv"
	"errors"
	"github.com/martin-bernal/go-bootcamp-project/internal/entity"
	"log"
	"os"
	"strconv"
)

type PokemonRepo interface {
	//Reads pokemons from the csv file
	ReadPokemons() ([]entity.Pokemon, error)

	WritePokemon(poke *entity.Pokemon) error
}

type pokemonRepo struct {
	filepath string
}

func NewPokemonRepo(filepath string) PokemonRepo {
	return &pokemonRepo{filepath}
}

func (pr *pokemonRepo) ReadPokemons() ([]entity.Pokemon, error) {
	log.Println("Enters to Repo/ReadPokemons function")

	var pokemons []entity.Pokemon
	f, err := os.Open(pr.filepath)
	if err != nil {
		return nil, errors.New("CSV file not found")
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.New("Unable to read csv file")
	}

	for i, line := range data {
		if i > 0 {
			pokeId, err := strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				return nil, errors.New("All IDs on pokemon database must be Integers")
			}

			pokeHeight, err := strconv.ParseInt(line[3], 10, 64)
			if err != nil {
				return nil, errors.New("All heights on pokemon database must be Integers")
			}

			baseExp, err := strconv.ParseInt(line[4], 10, 64)
			if err != nil {
				return nil, errors.New("All experience data on pokemon database must be Integers")
			}

			pokemons = append(pokemons, entity.Pokemon{
				ID:             pokeId,
				Name:           line[1],
				Type:           line[2],
				Height:         pokeHeight,
				BaseExperience: baseExp,
				Status:         line[5] != "0",
			})
		}
	}
	return pokemons, nil
}

func (pr *pokemonRepo) WritePokemon(poke *entity.Pokemon) error {
	f, err := os.OpenFile(pr.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.New("CSV file not found")
	}

	pokeStatus := 0
	if poke.Status {
		pokeStatus = 1
	}

	row := []string{
		strconv.FormatInt(poke.ID, 10),
		poke.Name,
		poke.Type,
		strconv.FormatInt(poke.Height, 10),
		strconv.FormatInt(poke.BaseExperience, 10),
		strconv.Itoa(pokeStatus),
	}

	w := csv.NewWriter(f)
	err = w.Write(row)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
