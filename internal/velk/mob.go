package velk

type Mob struct {
	Name            string     `json:"name"`
	ShortDesc       string     `json:"shortdesc"`
	Room            *Room      `json:"-"`
	Resources       *Resources `json:"resources"`
	FightingTargets []*Player  `json:"-"`
}

func NewMob() *Mob {
	resources := &Resources{
		Health:    100,
		MaxHealth: 100,
		Mana:      100,
		MaxMana:   100,
		Move:      100,
		MaxMove:   100,
	}

	return &Mob{
		Name:            "Untitled Mob",
		ShortDesc:       "An untitled mob",
		Resources:       resources,
		FightingTargets: make([]*Player, 0),
	}
}
