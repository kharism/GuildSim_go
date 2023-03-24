package cards_test

import (
	"github/kharism/GuildSim_go/internal/cards"
	"testing"
)

func TestResource(t *testing.T) {
	res := cards.NewResource()
	res.AddResource(cards.RESOURCE_NAME_MONEY, 100)
	if res.Detail[cards.RESOURCE_NAME_MONEY] != 100 {
		t.Error("Not same")
	}
	res.AddResource(cards.RESOURCE_NAME_MONEY, 10)
	if res.Detail[cards.RESOURCE_NAME_MONEY] != 110 {
		t.Error("Not same")
	}
	res.RemoveResource(cards.RESOURCE_NAME_MONEY, 90)
	if res.Detail[cards.RESOURCE_NAME_MONEY] != 110-90 {
		t.Error("Not same")
	}
}

func TestCompareCostResource(t *testing.T) {
	cost := cards.NewCost()
	cost.AddResource(cards.RESOURCE_NAME_COMBAT, 10)
	cost.AddResource(cards.RESOURCE_NAME_MONEY, 10)

	res := cards.NewResource()
	res.AddResource(cards.RESOURCE_NAME_COMBAT, 1)
	res.AddResource(cards.RESOURCE_NAME_MONEY, 1)

	if cost.IsEnough(res) {
		t.Error("it should be not enough")
	}
	res.AddResource(cards.RESOURCE_NAME_COMBAT, 10)
	res.AddResource(cards.RESOURCE_NAME_MONEY, 10)
	if !cost.IsEnough(res) {
		t.Error("it should be enough")
	}

	emptyCost := cards.NewCost()
	if !emptyCost.IsEnough(res) {
		t.Error("it should be enough")
	}
}
