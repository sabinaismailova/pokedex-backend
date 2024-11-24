package manipulations 

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example/pokemon/structs"
)

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

func flattenTypes(types []structs.Type) []string {
	var flattenedTypes []string
	for _, t := range types {
		flattenedTypes = append(flattenedTypes, t.Type.Name)
	}
	return flattenedTypes
}

func AddAbilityEffect(abilities []structs.Ability) {
	for i := range abilities {
		abilityName := abilities[i].Ability.Name
		effect, err := getAbilityEffect(abilityName)
		if err != nil {
			log.Println(err.Error())
		} else {
			abilities[i].Ability.Effect = effect
		}
	}
}

func FlattenPokemonData(pokemon structs.Pokemon) structs.FlattenedPokemon {
	pokemonData := structs.FlattenedPokemon{
		ID:       pokemon.ID,
		Name:     pokemon.Name,
		Abilities: flattenAbilities(pokemon.Abilities),
		Stats:     flattenStats(pokemon.Stats),
		Cry:       pokemon.Cry.Audio,
		Types:     flattenTypes(pokemon.Types),
		Sprites:   pokemon.Sprites,
	}
	return pokemonData
}
