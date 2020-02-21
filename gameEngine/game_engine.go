package gameEngine

import (
	"github.com/jsterling7/gcp-battle-peach/model"
)

//my url
const myURL string = "https://gcp-battle-peach-j6hslvdfdq-ue.a.run.app/battle"

var myScoreFiveRoundsAgo int
var myScoreFourRoundsAgo int
var myScoreThreeRoundsAgo int
var myScoreTwoRoundsAgo int
var myScoreOneRoundAgo int

var lastPostion model.Space
var lastMove string

var round int = 0

func Play(gameState model.GameState) string {

	round = round + 1

	myPlayerMap := getMyPlayerMap(gameState)

	myState := myPlayerMap[myURL]

	lastPostion = myState.Space

	updateSavedScores(myState.Score)

	otherPlayers := myPlayerMap
	delete(otherPlayers, myState.Name)

	//RULES

	// I can move forward one space
	// I can turn right
	// I can turn left
	// I can throw a peach
	// a thrown peach will travel up to three spaces in the direction the player is facing

	// LOGIC

	// if there is someone up to three spaces infront of me, I throw

	//if my score is not increasing after five rounds or has decreased move forward if possible and if not turn right

	if lastMove == "T" && myState.Score <= myScoreThreeRoundsAgo {
		// if forward is a safe and unoccupided space I move forward
		if isSafe := isForwardSafe(myState, otherPlayers, gameState.Arena.Dims); isSafe {
			lastMove = "F"
			return "F"
		}

		// else I turn to face the most optimal position to set up the next throw
		lastMove = "R"
		return "R"
	}

	// get hit spaces
	hitSpace1, hitSpace2, hitSpace3 := getHitSpaces(gameState.Arena.Dims, myState)

	// look for player in those spaces
	if _, inRange := playerInRange(otherPlayers, hitSpace1, hitSpace2, hitSpace3); inRange {
		lastMove = "T"
		return "T"
	}

	// if forward is a safe and unoccupided space I move forward
	if isSafe := isForwardSafe(myState, otherPlayers, gameState.Arena.Dims); isSafe {
		lastMove = "F"
		return "F"
	}

	// else I turn to face the most optimal position to set up the next throw
	lastMove = "R"
	return "R"
}

func updateSavedScores(myScore int) {
	myScoreFiveRoundsAgo = myScoreFourRoundsAgo
	myScoreFourRoundsAgo = myScoreThreeRoundsAgo
	myScoreTwoRoundsAgo = myScoreOneRoundAgo
	myScoreOneRoundAgo = myScore
}

func isForwardSafe(player model.Player, otherPlayers map[string]model.Player, areanaDims []int) bool {
	potentialSpace := getForwardSpace(player, areanaDims)
	if potentialSpace == nil {
		return false
	}

	for _, player := range otherPlayers {
		//anothe player is already in the space
		if player.Space.X == potentialSpace.X && player.Space.Y == potentialSpace.Y {
			return false
		}
	}

	// 12 potential postions around player where another player could be and be in range
	dangerNorthSpaces := []model.Space{}
	dangerSouthSpaces := []model.Space{}
	dangerEastSpaces := []model.Space{}
	dangerWestSpaces := []model.Space{}

	//north
	for i := 1; i < 3; i++ {
		currSpace := model.Space{
			X: potentialSpace.X,
			Y: potentialSpace.Y + i,
		}
		if inMap := isSpaceInMap(&currSpace, areanaDims); inMap {
			dangerNorthSpaces = append(dangerNorthSpaces, currSpace)
		}
	}

	//south
	for i := 1; i < 3; i++ {
		currSpace := model.Space{
			X: potentialSpace.X,
			Y: potentialSpace.Y - i,
		}
		if inMap := isSpaceInMap(&currSpace, areanaDims); inMap {
			dangerSouthSpaces = append(dangerSouthSpaces, currSpace)
		}
	}

	//east
	for i := 1; i < 3; i++ {
		currSpace := model.Space{
			X: potentialSpace.X + i,
			Y: potentialSpace.Y,
		}
		if inMap := isSpaceInMap(&currSpace, areanaDims); inMap {
			dangerEastSpaces = append(dangerEastSpaces, currSpace)
		}
	}
	//west
	for i := 1; i < 3; i++ {
		currSpace := model.Space{
			X: potentialSpace.X - i,
			Y: potentialSpace.Y,
		}
		if inMap := isSpaceInMap(&currSpace, areanaDims); inMap {
			dangerWestSpaces = append(dangerWestSpaces, currSpace)
		}
	}

	//if another play is in one of these positions and faceing the correct direction then forward is not safe

	for _, player := range otherPlayers {
		for _, space := range dangerNorthSpaces {
			if player.Space.X == space.X && player.Space.Y == space.Y && player.Direction == "S" {
				return false
			}
		}

		for _, space := range dangerSouthSpaces {
			if player.Space.X == space.X && player.Space.Y == space.Y && player.Direction == "N" {
				return false
			}
		}

		for _, space := range dangerEastSpaces {
			if player.Space.X == space.X && player.Space.Y == space.Y && player.Direction == "W" {
				return false
			}
		}

		for _, space := range dangerWestSpaces {
			if player.Space.X == space.X && player.Space.Y == space.Y && player.Direction == "E" {
				return false
			}
		}
	}

	return true
}

func getForwardSpace(player model.Player, areanaDims []int) *model.Space {
	forwardSpace := model.Space{}

	switch player.Direction {
	case "N":
		forwardSpace.X = player.Space.X
		forwardSpace.Y = player.Space.Y + 1
	case "S":
		forwardSpace.X = player.Space.X
		forwardSpace.Y = player.Space.Y - 1
	case "E":
		forwardSpace.X = player.Space.X + 1
		forwardSpace.Y = player.Space.Y
	case "W":
		forwardSpace.X = player.Space.X - 1
		forwardSpace.Y = player.Space.Y
	}

	if inMap := isSpaceInMap(&forwardSpace, areanaDims); !inMap {
		return nil
	}

	return &forwardSpace
}

func playerInRange(playerMap map[string]model.Player, hitSpace1, hitSpace2, hitSpace3 *model.Space) (*model.Player, bool) {
	for _, player := range playerMap {
		playerInHitSpaceOne := hitSpace1 != nil && player.Space.X == hitSpace1.X && player.Space.Y == hitSpace1.Y
		playerInHitSpaceTwo := hitSpace2 != nil && player.Space.X == hitSpace2.X && player.Space.Y == hitSpace2.Y
		playerInHitSpaceThree := hitSpace3 != nil && player.Space.X == hitSpace3.X && player.Space.Y == hitSpace3.Y

		if playerInHitSpaceOne || playerInHitSpaceTwo || playerInHitSpaceThree {
			return &player, true
		}
	}

	return nil, false
}

func getHitSpaces(areanaDims []int, player model.Player) (*model.Space, *model.Space, *model.Space) {
	// if direction is north
	// add 1 to y
	// south
	// subtract 1 from y
	// e
	// add 1 to x
	// w
	// subtract 1 from x

	hitSpace1 := &model.Space{}
	hitSpace2 := &model.Space{}
	hitSpace3 := &model.Space{}

	switch player.Direction {
	case "N":
		hitSpace1.X = player.Space.X
		hitSpace1.Y = player.Space.Y + 1
		hitSpace2.X = player.Space.X
		hitSpace2.Y = player.Space.Y + 2
		hitSpace3.X = player.Space.X
		hitSpace3.Y = player.Space.Y + 3
	case "S":
		hitSpace1.X = player.Space.X
		hitSpace1.Y = player.Space.Y - 1
		hitSpace2.X = player.Space.X
		hitSpace2.Y = player.Space.Y - 2
		hitSpace3.X = player.Space.X
		hitSpace3.Y = player.Space.Y - 3
	case "E":
		hitSpace1.X = player.Space.X + 1
		hitSpace1.Y = player.Space.Y
		hitSpace2.X = player.Space.X + 2
		hitSpace2.Y = player.Space.Y
		hitSpace3.X = player.Space.X + 3
		hitSpace3.Y = player.Space.Y
	case "W":
		hitSpace1.X = player.Space.X - 1
		hitSpace1.Y = player.Space.Y
		hitSpace2.X = player.Space.X - 2
		hitSpace2.Y = player.Space.Y
		hitSpace3.X = player.Space.X - 3
		hitSpace3.Y = player.Space.Y
	}

	if !isSpaceInMap(hitSpace1, areanaDims) {
		hitSpace1 = nil
	}
	if !isSpaceInMap(hitSpace2, areanaDims) {
		hitSpace2 = nil
	}
	if !isSpaceInMap(hitSpace3, areanaDims) {
		hitSpace3 = nil
	}

	return hitSpace1, hitSpace2, hitSpace3
}

func isSpaceInMap(space *model.Space, areanaDims []int) bool {
	return space.X >= 0 && space.X <= areanaDims[0] && space.Y >= 0 && space.Y <= areanaDims[1]
}

func getMyPlayerMap(gameState model.GameState) map[string]model.Player {
	playerMap := gameState.Arena.State.(map[string]interface{})

	myPlayerMap := make(map[string]model.Player)

	for playerKey, player := range playerMap {

		currPlayer := player.(map[string]interface{})

		newPlayer := model.Player{
			Name: playerKey,
		}

		for key, attribute := range currPlayer {
			switch attributeTypeValue := attribute.(type) {
			case string:
				if key == "direction" {
					newPlayer.Direction = attributeTypeValue
				}
			case float64:
				if key == "x" {
					newPlayer.Space.X = int(attributeTypeValue)
				}
				if key == "y" {
					newPlayer.Space.Y = int(attributeTypeValue)
				}
				if key == "score" {
					newPlayer.Score = int(attributeTypeValue)
				}
			case bool:
				if key == "wasHit" {
					newPlayer.WasHit = attributeTypeValue
				}
			}
		}
		myPlayerMap[newPlayer.Name] = newPlayer
	}

	return myPlayerMap
}
