package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/martin-bernal/go-bootcamp-project/internal/entity"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// DefaultTimeout is http client timeout and 60 seconds
	DefaultTimeout = 60 * time.Second
)

type client struct {
	endpoint   string
	restClient *http.Client
}

// Client Pokemon client contract
type Client interface {
	// GetPokemonFromApi Creates the request to be executed
	GetPokemonFromApi(pokemonName string) (*entity.Pokemon, error)
	// getPokemonRequest executed the create request
	getPokemonRequest(endpoint string, resource string) (response *Response, err error)
}

type Response struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Height         int64  `json:"height"`
	BaseExperience int64  `json:"base_experience"`
}

// NewPokemonClient initialize the pokemon client struct
func NewPokemonClient(endpoint string) (Client, error) {
	c := &client{
		endpoint: endpoint,
		restClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	return c, nil
}

func (c *client) GetPokemonFromApi(pokemonName string) (*entity.Pokemon, error) {
	response, err := c.getPokemonRequest(c.endpoint+"pokemon/", pokemonName)
	if err != nil {
		fmt.Println("error: ", err)
		return nil, errors.New("error getting pokemon")
	}

	pokemon := entity.Pokemon{
		ID:             response.ID,
		Name:           response.Name,
		Type:           "Normal",
		Height:         response.Height,
		BaseExperience: response.BaseExperience,
		Status:         true,
	}

	fmt.Print("response", response)
	return &pokemon, nil
}

func (c *client) getPokemonRequest(endpoint string, resource string) (response *Response, err error) {
	fmt.Println("pokemonApiClient/Get ", resource)
	defer fmt.Println("DONE pokemonApiClient/Get ", resource)
	req, err := http.NewRequest("GET", endpoint+resource, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return response, nil
}
