package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"example/pokemon/structs"
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

	for i := range pokemon.Abilities {
		abilityName := pokemon.Abilities[i].Ability.Name
		effect, err := getAbilityEffect(abilityName)
		if err != nil {
			log.Println(err.Error())
		} else {
			pokemon.Abilities[i].Ability.Effect = effect
		}
	}

	pokemonData := structs.FlattenedPokemon{
		ID:       pokemon.ID,
		Name:     pokemon.Name,
		Abilities: flattenAbilities(pokemon.Abilities),
		Stats:     flattenStats(pokemon.Stats),
		Sprites:   pokemon.Sprites,
	}

	c.IndentedJSON(http.StatusOK, pokemonData)
}

func getAbilityEffect(abilityName string) (string, error) {
    url := fmt.Sprintf("https://pokeapi.co/api/v2/ability/%s", abilityName)

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error fetching ability from PokeAPI: %v", err)
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error fetching ability for %v, status code: %d", abilityName, resp.StatusCode)
        return "", fmt.Errorf("failed to fetch ability: %s", abilityName)
    }

    var effects structs.EffectEntry

    if err := json.NewDecoder(resp.Body).Decode(&effects); err != nil {
        log.Printf("Error decoding ability response: %v", err)
        return "", err
    }

    for _, entry := range effects.EffectEntries {
        if entry.Language.Name == "en" {
            return entry.Effect, nil
        }
    }

    return "", fmt.Errorf("no english effect found for ability: %s", abilityName)
}

func flattenAbilities(abilities []structs.Ability) []structs.FlattenedAbility {
	var flattenedAbilities []structs.FlattenedAbility
	for _, ability := range abilities {
		flattenedAbilities = append(flattenedAbilities, structs.FlattenedAbility{
			Name:   ability.Ability.Name,
			Effect: ability.Ability.Effect,
		})
	}
	return flattenedAbilities
}

func flattenStats(stats []structs.Stat) []structs.FlattenedStat {
	var flattenedStats []structs.FlattenedStat
	for _, stat := range stats {
		flattenedStats = append(flattenedStats, structs.FlattenedStat{
			Name:  stat.Stat.Name,
			Value: stat.Value,
		})
	}
	return flattenedStats
}
