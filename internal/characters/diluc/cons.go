package diluc

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("diluc-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if t.HP()/t.MaxHP() > 0.5 {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	stack := 0
	last := 0
	c.Core.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
		if last != 0 && c.Core.F-last < 90 {
			return false
		}
		//last time is more than 10 seconds ago, reset stacks back to 0
		if c.Core.F-last > 600 {
			stack = 0
		}
		stack++
		if stack > 3 {
			stack = 3
		}
		m[attributes.ATKP] = 0.1 * float64(stack)
		m[attributes.AtkSpd] = 0.05 * float64(stack)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("diluc-c2", 600),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, "diluc-c2")
}

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.4
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("diluc-c4", -1),
		AffectedStat: attributes.DmgP,
		Amount: func() ([]float64, bool) {
			if c.Core.Status.Duration("dilucc4") > 0 {
				return m, true
			}
			return nil, false
		},
	})
}
