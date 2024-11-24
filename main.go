package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"pokedex/pokemon/structs"
	"pokedex/pokemon/manipulations"
)

func main() {
    router := gin.Default()

	router.Use(cors.Default())

    router.GET("/pokemon/:name", getPokemon)

	// router.GET("/stats/:name", getStats)

    router.Run("localhost:8080")
}

func getPokemon(c *gin.Context) {
	name := c.Param("name")

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching data from PokeAPI: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from PokeAPI"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received non-OK status code %d", resp.StatusCode)
		c.JSON(resp.StatusCode, gin.H{"error": "Pokemon not found"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var pokemon structs.Pokemon

	if err := json.Unmarshal(body, &pokemon); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Pokémon data"})
		return
	}

	manipulations.AddAbilityEffect(pokemon.Abilities)

	pokemonData := manipulations.FlattenPokemonData(pokemon)

	c.IndentedJSON(http.StatusOK, pokemonData)
}