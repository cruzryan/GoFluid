package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawNozzle() {
	for i := 0; i < 200; i++ {
		if i < (N/2)*SCALE { //Remove Scale

			ln := int32(math.Log(float64(i)) * 14.2)
			half := (int32(N) / int32(2)) * int32(SCALE) //REMOVE SCALE
			rl.DrawCircle(int32(i), int32(ln+half), 2, rl.DarkGray)
			rl.DrawCircle(int32(i), int32(-ln+half), 2, rl.DarkGray)

		}
	}

}

var (
	fluid Fluid
)

func main() {
	fmt.Println("Simulation started!")

	fluid.setup(0.2, 0, 0.0000001)

	// for i := 0; i < N; i++ {
	// 	for j := 0; j < N; j++ {
	// 		fluid.addDensity(i, j, 1)
	// 	}
	// }

	//Window Set-up
	rl.InitWindow(int32(N*SCALE), int32(N*SCALE), "Rocket Water")
	rl.SetTargetFPS(120)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		for i := 0; i < 5; i++ {
			for j := (N / 2) - 6; j < (N/2)+6; j++ {
				fluid.addDensity(i, j, 30)
				// fluid.addDensity(i+2, j, 30)
			}
		}
		// fluid.addVelocity(int(rl.GetMouseX()/int32(SCALE)), int(rl.GetMouseY()/int32(SCALE)), 2, 0)

		fluid.step()
		fluid.renderD()

		// rl.DrawRectangle(50, int32((N)/2)*3, 5, 100, rl.Red)
		// drawNozzle()

		// fluid.addVelocity(2, (N/2)+2, 1, 0)
		fluid.addVelocity(2, (N/2)+1, 4, 0)
		fluid.addVelocity(2, N/2, 4, 0)
		fluid.addVelocity(2, (N/2)-1, 4, 0)
		// fluid.addVelocity(2, (N/2)-2, 1, 0)
		rl.EndDrawing()

	}

	rl.CloseWindow()

}
