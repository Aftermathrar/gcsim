package shenhe

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Shenhe, NewChar)
}

type char struct {
	*tmpl.Character
	skillBuff []float64
	c4bonus   []float64
	c4count   int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	c.c4count = 0

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillBuff = make([]float64, attributes.EndStatType)
	c.skillBuff[attributes.DmgP] = 0.15
	c.quillDamageMod()
	if c.Base.Cons >= 4 {
		c.c4bonus = make([]float64, attributes.EndStatType)
		c.c4()
	}
	return nil
}
