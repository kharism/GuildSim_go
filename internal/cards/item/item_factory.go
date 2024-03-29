package item

import (
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
	case name == "Greed Potion":
		h := NewGreedPotion(state)
		return &h
	case name == "RefreshPotion":
		h := NewRefreshPotion(state)
		return &h
	case name == "Recursion Potion":
		h := NewRecursionPotion(state)
		return &h
	case name == "Banish Potion":
		h := NewBanishPotion(state)
		return &h
	case name == "PotionofGreed":
		h := NewGreedPotion(state)
		return &h
	default:
		h := NewHealingPotion(state)
		return &h
	}
	return nil
}

func CreateRelicRandom(state cards.AbstractGamestate, rarity int) cards.Card {
	commonList := []string{"CombatGauntlet", "ExplorerBoots", "Regen Amulet"}
	rareList := []string{"CompanionBuckler", "BloodyCompass"}
	allList := []string{}
	if (rarity & cards.RARITY_COMMON) != 0 {
		allList = append(allList, commonList...)
	}
	if (rarity & cards.RARITY_RARE) != 0 {
		allList = append(allList, rareList...)
	}
	picked := rand.Int() % len(allList)
	newPotion := CreateRelic(allList[picked], state)
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
	case "CompanionBuckler":
		h := NewCompanionBuckler(state)
		return &h
	case "BloodyCompass":
		h := NewBloodyCompass(state)
		return &h
	}
	return nil
}
func CreatePotionRandom(state cards.AbstractGamestate, rarity int) cards.Card {
	commonList := []string{"Explore Potion", "Combat Potion", "Healing Potion", "Banish Potion", "PotionofGreed"}
	rareList := []string{"RefreshPotion", "Recursion Potion", "Uncurse potion"}
	allList := []string{}
	if (rarity & cards.RARITY_COMMON) != 0 {
		allList = append(allList, commonList...)
	}
	if (rarity & cards.RARITY_RARE) != 0 {
		allList = append(allList, rareList...)
	}
	picked := rand.Int() % len(allList)
	// fmt.Println("ALLLIST", allList, allList[picked])
	newPotion := CreatePotion(allList[picked], state)
	// fmt.Println(newPotion)
	return newPotion
}
