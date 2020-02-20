package model



//{
//	"_links": {
//		"self": {
//			"href": "https://YOUR_SERVICE_URL"
//		}
//	},
//	"arena": {
//		"dims": [4,3], // width, height
//		"state": {
//			"https://A_PLAYERS_URL": {
//				"x": 0, // zero-based x position, where 0 = left
//				"y": 0, // zero-based y position, where 0 = top
//				"direction": "N", // N = North, W = West, S = South, E = East
//				"wasHit": false,
//				"score": 0
//			}
//	... // also you and the other players
//		}
//	}
//}


type GameState struct {
	Arena Arena `json:"arena"`
}


type Arena struct {
	Dims []int `json:"dims"`
	State interface{} `json:"state"`
}