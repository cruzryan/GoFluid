package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	N     = 128
	iter  = 16
	SCALE = 5
)

func constrain(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

func IX(x, y int) int {
	x = constrain(x, 0, N-1)
	y = constrain(y, 0, N-1)
	k := x + (y * N)
	if k >= 16384 {
		fmt.Println("FLUID: OUT OF BOUNDS BRO, REQUESTED: ", k)
		return 1
	} else {
		return k
	}
}

type Fluid struct {
	size int
	dt   float32
	diff float32
	visc float32

	s       []float32
	density []float32

	Vx []float32
	Vy []float32

	Vx0 []float32
	Vy0 []float32
}

func (f *Fluid) setup(dt, diffusion, viscosity float32) {

	f.size = N
	f.dt = dt
	f.diff = diffusion
	f.visc = viscosity

	f.s = make([]float32, N*N)
	f.density = make([]float32, N*N)

	f.Vx = make([]float32, N*N)
	f.Vy = make([]float32, N*N)

	f.Vx0 = make([]float32, N*N)
	f.Vy0 = make([]float32, N*N)

	fmt.Println("FLUID: Setup done!")
}

func (f *Fluid) step() {

	//El algoritmo
	diffuse(1, f.Vx0, f.Vx, f.visc, f.dt)
	diffuse(2, f.Vy0, f.Vy, f.visc, f.dt)

	project(f.Vx0, f.Vy0, f.Vx, f.Vy)

	advect(1, f.Vx, f.Vx0, f.Vx0, f.Vy0, f.dt)
	advect(2, f.Vy, f.Vy0, f.Vx0, f.Vy0, f.dt)

	project(f.Vx, f.Vy, f.Vx0, f.Vy0)

	diffuse(0, f.s, f.density, f.diff, f.dt)
	advect(0, f.density, f.s, f.Vx, f.Vy, f.dt)
}

func (f *Fluid) addDensity(x, y int, amount float32) {
	index := IX(x, y)
	f.density[index] += amount
}

func (f *Fluid) addVelocity(x, y int, amountX, amountY float32) {
	index := IX(x, y)
	// fmt.Println("\nFLUID: VELOCITY INDEX RETURNED: ", index)
	f.Vx[index] += amountX
	f.Vy[index] += amountY
}

func (f *Fluid) renderD() {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			x := i * SCALE
			y := j * SCALE

			d := uint8(f.density[IX(i, j)])
			// fmt.Println("DENSITY:", d)

			// alpha := uint8(int((d + 50)) % 25)
			// fmt.Println("ALPHA:", alpha)
			// alpha := uint8(255)
			// fmt.Println("X: ", x, "Y: ", y)
			if d == 1 || d == 0 {
				rl.DrawRectangle(int32(x), int32(y), int32(SCALE), int32(SCALE), rl.NewColor(0, 0, 0, 255))
			} else {

				slope := (255 - 0) / (180 - 0)
				output := uint8(0 + uint8(slope)*(d-0))
				rl.DrawRectangle(int32(x), int32(y), int32(SCALE), int32(SCALE), rl.NewColor(output, output, output-140, output))
			}

		}
	}
}
