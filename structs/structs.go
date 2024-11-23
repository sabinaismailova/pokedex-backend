package structs

type Sprites struct {
	BackDefault        *string `json:"back_default"`         
	BackFemale         *string `json:"back_female"`          
	BackShiny          *string `json:"back_shiny"`           
	BackShinyFemale    *string `json:"back_shiny_female"`    
	FrontDefault       *string `json:"front_default"`        
	FrontFemale        *string `json:"front_female"`         
	FrontShiny         *string `json:"front_shiny"`          
	FrontShinyFemale   *string `json:"front_shiny_female"`   
}

type Ability struct {
	Ability  struct {
        Name string `json:"name"`
		Effect string `json:"effect"`
    } `json:"ability"`
}

type EffectEntry struct {
	EffectEntries []struct {
        Effect  string `json:"short_effect"`
        Language struct {
            Name string `json:"name"`
        } `json:"language"`
    } `json:"effect_entries"`
}

type Stat struct {
	Value int `json:"base_stat"`
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Pokemon struct {
    ID     int  `json:"id"`
    Name  string  `json:"name"`
    Sprites Sprites  `json:"sprites"`
	Abilities []Ability `json:"abilities"`
	Stats []Stat `json:"stats"`
}

type FlattenedStat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type FlattenedAbility struct {
	Name   string `json:"name"`
	Effect string `json:"effect"`
}

type FlattenedPokemon struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Abilities  []FlattenedAbility `json:"abilities"`
	Stats      []FlattenedStat   `json:"stats"`
	Sprites    Sprites  `json:"sprites"`
}
