package venti

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterCharFunc(keys.Venti, NewChar)
}

type char struct {
	*tmpl.Character
	qInfuse             attributes.Element
	infuseCheckLocation combat.AttackPattern
	aiAbsorb            combat.AttackInfo
	snapAbsorb          combat.Snapshot
	c4bonus             []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.infuseCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableObject)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 4 {
		c.c4bonus = make([]float64, attributes.EndStatType)
		c.c4bonus[attributes.AnemoP] = 0.25
	}
	return nil
}

func (c *char) ReceiveParticle(p character.Particle, isActive bool, partyCount int) {
	c.Character.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 4 {
		//only pop this if active
		if !isActive {
			return
		}
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("venti-c4", 600),
			AffectedStat: attributes.AnemoP,
			Amount: func() ([]float64, bool) {
				return c.c4bonus, true
			},
		})
	}
}
