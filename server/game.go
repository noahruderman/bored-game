package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/coder/websocket"
)

type Game struct {
	board [][]Tile
	seed, chunkSize float64

	players []*Player
	playerTurn int
}

func (g *Game) resizeBoard(length int) {
	g.board = make([][]Tile, length)
	for i := range length {
		g.board[i] = make([]Tile, length)
	}
}

func (g *Game) generateMap(length int) {
	g.resizeBoard(length)

	for i := range length {
		for j := range length {
			noise := perlinNoise(float64(j) / g.chunkSize, float64(i) / g.chunkSize, g.seed)
			tile := &g.board[i][j].tiletype

			if noise < .25 {
				*tile = OCEAN
			} else if noise < .4 {
				*tile = WATER
			} else if noise < .45 {
				*tile = COAST
			} else if noise < .52 {
				*tile = PLAINS
			} else if noise < .62 {
				*tile = GRASS
			} else if noise < .72 {
				*tile = FOREST
			} else if noise < .8 {
				*tile = MOUNTAIN
			} else {
				*tile = VOLCANO
			}
		}
	}
}

func (g *Game) broadcastMessage(skipPlayer *Player, ctx context.Context, msg string) {
	toSend := fmt.Append(nil, "broadcast\n", msg)
	for _, p := range g.players {
		if p != skipPlayer && p.connected {
			p.c.Write(ctx, websocket.MessageText, toSend)
		}
	}
}

func (g *Game) updateWithGameState(player *Player, ctx context.Context) error {
	// player index
	response_player_index := fmt.Sprintf("player-index %d\n", player.player_index)

	// map
	response_map := fmt.Sprintf("map %d %f %f\n", len(g.board), g.chunkSize, g.seed)

	// troops
	response_troops := strings.Builder{}
	fmt.Fprintf(&response_troops, "troops %d\n", len(g.players))
	for _, p := range g.players {
		fmt.Fprintf(&response_troops, "%d ", len(p.troops))
	}
	response_troops.WriteByte('\n')

	for _, p := range g.players {
		for _, t := range p.troops {
			fmt.Fprintf(&response_troops, "%d %d,", t.x, t.y)
		}
		response_troops.WriteByte('\n')
	}

	// modified tiles
	response_tiles := strings.Builder{}
	response_tiles.WriteString("modified-tiles\n")
	for i := range g.board {
		for j := range g.board[i] {
			if g.board[i][j].modified {
				response_tiles.WriteByte('m')
			} else {
				response_tiles.WriteByte('.')
			}
		}
		response_tiles.WriteByte('\n')
	}

	response := response_player_index + response_map + response_troops.String() + response_tiles.String()
	return player.c.Write(ctx, websocket.MessageText, []byte(response))
}

func (g *Game) updateWithMap(player *Player, ctx context.Context) error {
	response := fmt.Sprintf("broadcast\nupdate-map %d %f %f\n", len(g.board), g.chunkSize, g.seed)
	return player.c.Write(ctx, websocket.MessageText, []byte(response))
}

type Player struct {
	c *websocket.Conn
	connected bool

	player_index int
	troops []Troop
	wood, stone int
}
