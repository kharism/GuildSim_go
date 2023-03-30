package item

import (
	"fmt"
	"github/kharism/GuildSim_go/internal/cards"
	"math/rand"
)

func CreatePotion(name string, state cards.AbstractGamestate) cards.Card {
	switch {
	case name == "Combat Potion":
		h := NewCombatPotion(state)
		return &h
	case name == "Explore Potion":
		h := NewExplorePotion(state)
		return &h
	case name == "Healing Potion":
		h := NewHealingPotion(state)
		return &h
	default:
		h := NewHealingPotion(state)
		return &h
	}
	return nil
}
func CreateRelicRandom(state cards.AbstractGamestate, rarity int) cards.Card {
	commonList := []string{}
	rareList := []string{"CombatGauntlet", "ExplorerBoots", "Regen Amulet"}
	allList := []string{}
	if (rarity & cards.RARITY_COMMON) != 0 {
		allList = append(allList, commonList...)
	}
	if (rarity & cards.RARITY_RARE) != 0 {
		allList = append(allList, rareList...)
	}
	picked := rand.Int() % len(allList)
	newPotion := CreatePotion(allList[picked], state)
	return newPotion
}

func CreateRelic(name string, state cards.AbstractGamestate) cards.Card {
	switch name {
	case "CombatGauntlet":
		h := NewCombatGauntlet(state)
		return &h
	case "ExplorerBoots":
		h := NewExplorerBoots(state)
		return &h
	case "Regen Amulet":
		h := NewRegenAmulet(state)
		return &h
	}
	return nil
}
func CreatePotionRandom(state cards.AbstractGamestate, rarity int) cards.Card {
	commonList := []string{"Explore Potion", "Combat Potion", "Healing Potion"}
	rareList := []string{"RefreshPotion", "Recursion Potion"}
	allList := []string{}
	if (rarity & cards.RARITY_COMMON) != 0 {
		allList = append(allList, commonList...)
	}
	if (rarity & cards.RARITY_RARE) != 0 {
		allList = append(allList, rareList...)
	}
	picked := rand.Int() % len(allList)
	fmt.Println("ALLLIST", allList, allList[picked])
	newPotion := CreatePotion(allList[picked], state)
	fmt.Println(newPotion)
	return newPotion
}
