package dori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillRelease = 16
)

var skillSalesHitmarks = []int{46, 59, 59} // counted starting from skill hitmark

func init() {
	skillFrames = frames.InitAbilSlice(44) // E -> Q
	skillFrames[action.ActionDash] = 43    // E -> D
	skillFrames[action.ActionSwap] = 43    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Troubleshooter Shot",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	afterSalesCB := func(_ combat.AttackCB) { // executes after the troublshooter shot hits
		c.afterSales(travel)
	}

	// When the Troubleshooter Shots or After-Sales Service Rounds from
	// Spirit-Warding Lamp: Troubleshooter Cannon hit opponents, Dori will
	// restore 5 Elemental Energy for every 100% Energy Recharge possessed. Per
	// Spirit-Warding Lamp: Troubleshooter Cannon, only one instance of Energy
	// restoration can be triggered and a maximum of 15 Energy can be restored
	// this way.
	done := false
	a4CB := func(a combat.AttackCB) {
		if done {
			return
		}
		a4energy := a.AttackEvent.Snapshot.Stats[attributes.ER] * 5
		if a4energy > 15 {
			a4energy = 15
		}
		c.AddEnergy("dori-a4", a4energy)
		done = true
	}

	c.Core.Tasks.Add(func() {
		// C6
		if c.Base.Cons >= 6 {
			c.Core.Player.AddWeaponInfuse(
				c.Index,
				c6key,
				attributes.Electro,
				180,
				true,
				combat.AttackTagNormal,
				combat.AttackTagExtra,
				combat.AttackTagPlunge,
			)
		}
	}, skillFrames[action.ActionAttack]) // TODO:It activates on the attack cancel frames?

	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
		0,
		skillRelease+travel,
		afterSalesCB,
		a4CB,
	)

	c.SetCDWithDelay(action.ActionSkill, 9*60, 16)
	c.Core.QueueParticle("dori", 2, attributes.Electro, skillRelease+travel+c.ParticleDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) afterSales(travel int) func() {
	return func() {
		ae := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "After-Sales Service Round",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillAfter[c.TalentLvlSkill()],
		}
		for i := 0; i < c.afterCount; i++ {
			c.Core.QueueAttack(
				ae,
				combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
				0,
				skillSalesHitmarks[i],
			)
		}
	}
}
