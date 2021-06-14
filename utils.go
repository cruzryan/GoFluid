package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func floor(val float32) float32 {
	//I love Go but this... I hate
	return float32(math.Floor(float64(val)))
}

func diffuse(b int, x []float32, x0 []float32, diff float32, dt float32) {
	a := dt * diff * (float32(N) - 2) * (float32(N) - 2)
	lin_solve(b, x, x0, a, 1+4*a)
}

func lin_solve(b int, x []float32, x0 []float32, a float32, c float32) {
	cRecip := 1.0 / c

	for k := 0; k < iter; k++ {
		for j := 1; j < N-1; j++ {
			for i := 1; i < N-1; i++ {
				x[IX(i, j)] = (x0[IX(i, j)] + a*(x[IX(i+1, j)]+x[IX(i-1, j)]+x[IX(i, j+1)]+x[IX(i, j-1)])) * cRecip
			}
		}

		//Sets border boundaries
		set_bnd(b, x)
	}
}

func project(velocX, velocY, p, div []float32) {
	for j := 1; j < N-1; j++ {
		for i := 1; i < N-1; i++ {
			div[IX(i, j)] = -0.5 * (velocX[IX(i+1, j)] - velocX[IX(i-1, j)] + velocY[IX(i, j+1)] - velocY[IX(i, j-1)]) / float32(N)
			p[IX(i, j)] = 0
		}
	}

	set_bnd(0, div)
	set_bnd(0, p)
	lin_solve(0, p, div, 1, 4)

	for j := 1; j < N-1; j++ {
		for i := 1; i < N-1; i++ {
			velocX[IX(i, j)] -= 0.5 * (p[IX(i+1, j)] - p[IX(i-1, j)]) * float32(N)
			velocY[IX(i, j)] -= 0.5 * (p[IX(i, j+1)] - p[IX(i, j-1)]) * float32(N)
		}
	}
	set_bnd(1, velocX)
	set_bnd(2, velocY)
}

func advect(b int, d []float32, d0 []float32, velocX []float32, velocY []float32, dt float32) {

	var i0, i1, j0, j1 float32

	dtx := dt * (float32(N) - 2)
	dty := dt * (float32(N) - 2)

	var s0, s1, t0, t1 float32
	var tmp1, tmp2, x, y float32

	Nfloat := float32(N)
	var ifloat, jfloat float32

	var i, j int

	for j, jfloat = 1, 1; j < N-1; j, jfloat = j+1, jfloat+1 {
		for i, ifloat = 1, 1; i < N-1; i, ifloat = i+1, ifloat+1 {
			tmp1 = dtx * velocX[IX(i, j)]
			tmp2 = dty * velocY[IX(i, j)]
			x = ifloat - tmp1
			y = jfloat - tmp2

			if x < 0.5 {
				x = 0.5
			}
			if x > Nfloat+0.5 {
				x = Nfloat + 0.5
			}

			i0 = floor(x)
			i1 = i0 + 1.0
			if y < 0.5 {
				y = 0.5
			}
			if y > Nfloat+0.5 {
				y = Nfloat + 0.5
			}
			j0 = floor(y)
			j1 = j0 + 1.0

			s1 = x - i0
			s0 = 1.0 - s1
			t1 = y - j0
			t0 = 1.0 - t1

			i0i := int(i0)
			i1i := int(i1)
			j0i := int(j0)
			j1i := int(j1)

			// bruh?
			d[IX(i, j)] = s0*(t0*d0[IX(i0i, j0i)]+t1*d0[IX(i0i, j1i)]) + s1*(t0*d0[IX(i1i, j0i)]+t1*d0[IX(i1i, j1i)])
		}
	}

	set_bnd(b, d)
}

//To stop it from coliding with itself
func set_bnd(b int, x []float32) {

	for i := 0; i < 100; i++ {
		if i < (N / 5) { //Remove Scale

			ln := int32(math.Log(float64(i)) * 3)
			half := (int32(N) / int32(2)) //REMOVE SCALE
			fluid.addDensity(i-30, int(ln+half), 0)
			fluid.addDensity(i-30, int(-ln+half), 0)
			rl.DrawCircle(int32(i)*int32(SCALE), int32(ln+half)*int32(SCALE), 5, rl.Red)
			rl.DrawCircle(int32(i)*int32(SCALE), int32(-ln+half)*int32(SCALE), 5, rl.Red)

			if b == 2 {
				x[IX(i, int(ln+half))] = -x[IX(i, int(ln+half))]
				x[IX(i, int(-ln+half))] = -x[IX(i, int(-ln+half))]
			}

		}
	}

	// //Nozzle Border
	// for l := 0; l < N; l++ {

	// 	if l < (N / 2) { //Remove Scale

	// 		ln := int32(math.Log(float64(l)) * 10)
	// 		half := (int32(N) / int32(2)) //REMOVE SCALE

	// 		y := int(ln + half)
	// 		y2 := int(-ln + half)

	// 		fluid.addDensity(l, y, 100)
	// 		fluid.addDensity(l, y2, 100)

	// 		x[IX(l, y)] = -4 * x[IX(l, y)]
	// 		x[IX(l, y2)] = -4 * x[IX(l, y2)]

	// 		// if b == 1 {
	// 		// 	x[IX(l, y)] = -4 * x[IX(l, y)]
	// 		// 	x[IX(l, y2)] = -4 * x[IX(l, y2)]
	// 		// } else {
	// 		// 	x[IX(l, y)] = -4 * x[IX(l, y)]
	// 		// 	x[IX(l, y2)] = -4 * x[IX(l, y2)]
	// 		// }

	// 		return
	// 	}
	// }

	for i := 1; i < N-1; i++ {

		if b == 2 {
			x[IX(i, 0)] = -x[IX(i, 1)]
			// x[IX(i, N-1)] = -x[IX(i, N-2)]

			//NEW IMP
			// x[IX(i, i)] = -x[IX(i, i-2)]
		} else {
			x[IX(i, 0)] = x[IX(i, 1)]
			// x[IX(i, N-1)] = x[IX(i, N-2)]

			//NEW IMP
			// x[IX(i, (i)-1)] = -x[IX(i, (i))]

		}
	}

	for j := 1; j < N-1; j++ {

		if b == 1 {
			x[IX(0, j)] = -x[IX(1, j)]
			// x[IX(N-1, j)] = -x[IX(N-2, j)]
		} else {
			x[IX(0, j)] = x[IX(1, j)]
			// x[IX(N-1, j)] = x[IX(N-2, j)]
		}
	}

	x[IX(0, 0)] = 0.5 * (x[IX(1, 0)] + x[IX(0, 1)])
	x[IX(0, N-1)] = 0.5 * (x[IX(1, N-1)] + x[IX(0, N-2)])
	x[IX(N-1, 0)] = 0.5 * (x[IX(N-2, 0)] + x[IX(N-1, 1)])
	x[IX(N-1, N-1)] = 0.5 * (x[IX(N-2, N-1)] + x[IX(N-1, N-2)])
}
