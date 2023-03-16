# GuildSim

## Intro

Play as guild manager, dispatch hero to explore area, defeat monster, gains reputation and be the best guild in town

This is just my interpertation of how ascension should be. Their gameplay loop is fine but sometimes I feel the pacing can be chore because no player can acquire cards in center row.

## Objective of the game

Defeat 3 bosses in the game to clear the game

## Gameplay loop

1. explore area to gain reputation, gold, recruit better heroes and unlocks harder monster
2. defeat monsters to gain more reputation. If monsters is not immediately taken care of there will be some side effect each turn
3. once certain conditions are met (reputation level, gold level, deck size etc), boss or mini boss is can appear into the center row to be defeated
4. repeat

## Code
I'll try to make the code as robust as I can. I only have minor experience making game and I don't even know what the game will looks like at the end. Probably over-engineer something for no good reason.

For now I only support textUI only while I'm also planning to make some simple GUI. But I haven't decided what UI framework to choose at the moment

## installation

### TextUI
pull or clone the repo and just run `go build main.go` and run the executable

### GUI
still pending

## How to play

### TextUI
Run the program and you should see something like this
```
Resource
HP: 60
Cards In center Row:
[] RookieNurse (Hero) [Exploration:1]:Heal 2 Hp
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
Cards in hand:
[0] RookieCombatant []:Add 1 Combat point
[1] RookieAdventurer []:Add 1 exlporation point
[2] RookieCombatant []:Add 1 Combat point
[3] RookieAdventurer []:Add 1 exlporation point
[4] RookieAdventurer []:Add 1 exlporation point

<prompt>
```
- HP is your hit point, when it goes to 0 you loose.

- cards in center row is available hero to recruit, monster to defeat, and area to explore. use 'qwerty' to pick which card to recruit/defeat/explore. 

The format of the cards on the Center row is [\<key to interact>] \<card name> (\<card type>) [\<card cost>]: <card eff/desc>

for example from the view above:
```
[] RookieNurse (Hero) [Exploration:1]:Heal 2 Hp
```
means the card Hero RookieNurse has cost exploration point of 1, and can heal 2Hp when played from hand. We don't have enough resource to recruit it since it has empty square bracket. Once we have enough resource the square bracket will have something inside to indicate which key you need to input to recruit the hero, explore the area or defeat the monster.

Explored area and defeated monsters will be banished, but recruited hero will be sent to cooldown/discard pile. Cards will be moved from discard pile to deck when the deck run out of cards.

cards in hand is card in your hand, input number 0-4 to play the card. or input 'A' to play all cards from hand, from the largest index to index 0.

after you play some cards the screen will turn into something like 
```
Resource
HP: 60 Combat:2 Exploration:3
Cards In center Row:
[q] RookieNurse (Hero) [Exploration:1]:Heal 2 Hp
[w] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[e] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[r] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[t] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
Cards in hand:
<prompt>
```
as you can see the square bracket has letter in them. indicating you can interact with some of the cards. This is due to you have generated resource which you can see on the right side of HP.
```
Resource
HP: 60 Combat:2 Exploration:3
```
in this situation we have 3 exploration and 2 combat point. We then distribute the point to recruit hero or defeat monster accordingly. If we input q, then we can see RookieNurse is replaced by something else, probably something like
```
Resource
HP: 60 Combat:2 Exploration:2
Cards In center Row:
[q] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[w] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[e] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[r] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[t] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
Cards in hand:
```
Now we spent 1 exploration to recruit rookienurse and it now sit in our discard/cooldown pile, waiting for our player deck to empty.

since we have a lot of GoblinMonster, and combat point we should defeat some of them. A monster can do something nasty on our endphase if not dealt with immediately. As you can see on goblin monster, each of them can inflict 1 damage to us. We spent 2 of our combat point to deals with them and reduce the possible damage to ourself into only taking 3 instead of 5.

After defeated 2 of them we should have something like

```
Resource
HP: 60 Combat:0 Exploration:2 Reputation:2
Cards In center Row:
[] GoblinSmallLairArea (Area) [Exploration:3]:Reward: 100Money and 2 Reputation also shuffle goblinwolfraider into center deck
[w] RookieNurse (Hero) [Exploration:1]:Heal 2 Hp
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
[] GoblinMonster (Mons) [Combat:1]:1 dmg per turn, Reward: 1 Reputation
Cards in hand:
```

if you want to end your turn input empty string. On end of turn you will take damage accordingly, and whatever resource you have except reputation and money will be reset to 0. Then on the start of new turn you draw 5 cards. It is good idea to play as many card as you can.
