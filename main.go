package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(800, 450, "Deck Crawler")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawText("Welcome to the Deck Crawler!", 190, 200, 20, rl.LightGray)
		rl.EndDrawing()
	}
}
