package main

import (
		"fmt"
 		"math/rand"
 		"math"
 		"time"
 		)

const(
		WIDTH = 46
		HEIGHT = 24

		α = 0.5  // factor de aprendizaje
		ε = 0.1  // numero muy peque
		γ = 1.0  // factor de olvido
		RECOMPENSA = -1

		INIT_X = 1	
		INIT_Y = HEIGHT - 2

		GOAL_X = WIDTH - 2
		GOAL_Y = HEIGHT - 2

		ACT = 4  // numero a de acciones
		
		ARR = 0
		AB = 1
		IZQ = 2
		DER = 3
	)

type A uint8

type S struct{
	x uint16
	y uint16
}

type MDP struct{
	a [ACT]A
	St S
	R int8
	γ float32
}

func ( mdp *MDP ) init_S() S {
	var s S

	s.x = INIT_X
	s.y = INIT_Y

	return s
}

func ( mdp *MDP ) is_termnal( s S ) bool {
	return s.x == mdp.St.x && s.y == mdp.St.y
}

func ( mdp *MDP ) init_Q() [HEIGHT][WIDTH][ACT]float32 {
	var Q[HEIGHT][WIDTH][ACT]float32

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			for k := 0; k < ACT; k++ {
				Q[i][j][k] = 0
			}
		}
	}
	return Q
}

func ( mdp *MDP ) get_action( s S, Q[HEIGHT][WIDTH][ACT]float32 ) A {
	if rand.Float64() < 1 - ε {
		
		return mdp.max_a( Q[s.y][s.x] )
	} else {
		
		return A( rand.Intn( ACT ) )
	}
}

func ( mdp *MDP ) take_action( a A, s S) (S, int8){

	var n_s S
	n_s.x, n_s.y = s.x, s.y

	switch a {
		case ARR: // 0
			n_s.x = s.x
			
			if s.y - 1 > 0 || s.y - 1 == 1{
				n_s.y = s.y - 1
			}else{
				n_s.y = s.y
			}

		case AB: // 1
			n_s.x = s.x
			if s.y + 1 < HEIGHT {
				n_s.y = s.y + 1
			}else{
				n_s.y = s.y
			}

		case IZQ: // 2
			
			if s.x - 1 > 0 || s.x - 1 == 1{
				n_s.x = s.x - 1
			}else{
				n_s.x = s.x  
			}
			n_s.y = s.y
			
		case DER: // 3
			//fmt.Println( s.x + 1, " < WIDTH" )
			if s.x + 1 < WIDTH {
				n_s.x = s.x + 1
			}else{
				n_s.x = s.x 
			}
			n_s.y = s.y

			//fmt.Println( s.y , s.x )
	}

	var r int8 = 0
	
	if mdp.in_cliff( n_s )  {
		r = 100*mdp.R
		n_s.x = INIT_X
		n_s.y = INIT_Y

	}else{
		r = mdp.R
	}
	return n_s, r
}

func ( mdp *MDP ) in_cliff( s S ) bool {
	return ( s.y >= HEIGHT-3 && s.y < HEIGHT ) && ( s.x >= 3 && s.x < WIDTH - 3 )
}

func ( mdp *MDP ) max_a( act [ACT]float32  ) A {
	var (
		a A
		mayor float32 = float32( math.Inf( -1 ) )
	)

	for index, value := range act {
		if mayor < value {
			mayor = value
			a = A( index )
		}
	}
	return a
}

func Q_learning( mdp MDP, α float32, max_episode int, max_step int ) [HEIGHT][WIDTH][ACT]float32 {
	Q := mdp.init_Q()
	var (
		a A
		s, n_s S
		r int8
	)

	for i := 0; i < max_episode; i++ {

		s.x = INIT_X
		s.y = INIT_Y

		for j := 0; j < max_step; j++ {
			
			a = mdp.get_action( s, Q )

			n_s, r = mdp.take_action( a, s )
			
			Q[s.y][s.x][a] += α*( float32(r) + mdp.γ * Q[n_s.y][n_s.x][ int(mdp.max_a( Q[n_s.y][n_s.x] )) ] - Q[s.y][s.x][a] ) 
			
			if mdp.in_cliff( n_s ){
				s.x = INIT_X
				s.y = INIT_Y
			}else{
				s.y = n_s.y
				s.x = n_s.x
			}

			if mdp.is_termnal( s ){
				break
			}
		}
	}
	return Q
}

func action( s S, a A ) S {
	
	var n_s S

	switch a {
		case ARR: // 0
			n_s.x = s.x
			
			if s.y - 1 > 0 {
				n_s.y = s.y - 1
			}else{
				n_s.y = s.y
			}

		case AB: // 1
			n_s.x = s.x
			if s.y + 1 < HEIGHT {
				n_s.y = s.y + 1
			}else{
				n_s.y = s.y
			}

		case IZQ: // 2
			
			if s.x - 1 > 0 {
				n_s.x = s.x - 1
			}else{
				n_s.x = s.x  
			}
			n_s.y = s.y
			
		case DER: // 3
			//fmt.Println( s.x + 1, " < WIDTH" )
			if s.x + 1 < WIDTH {
				n_s.x = s.x + 1
			}else{
				n_s.x = s.x 
			}
			n_s.y = s.y

			//fmt.Println( s.y , s.x )
	}

	return n_s 
}

func optimal_policy( Q [HEIGHT][WIDTH][ACT]float32 ) map[S]A {
	policy := make( map[S]A )
	s := S{ x: INIT_X, y: INIT_Y  }
	var (
		a A
		mayor float32 = float32( math.Inf( -1 ) )
	)

 	for {
 		for act, value := range Q[s.y][s.x] {
			
			if mayor < value{
				mayor = value
				a = A(act)
			}	
		}		
		policy[s] = a
		
		s = action( s , a )

		if s.x == GOAL_X && s.y == GOAL_Y{
			break
		}
 	}
 	
 	return policy
}

func get_move( a A ) rune {
	var char rune
	
	switch a {
		case ARR:
			char = '↑'
		case AB:
			char = '↓'
		case IZQ:
			char = '←'
		case DER:
			char = '→'
	}
	return char
}

func print_policy( policy map[S]A ){
	var mapa[HEIGHT][WIDTH] rune	

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			mapa[i][j] = '_'
		}
	}

	for i := HEIGHT-3; i < HEIGHT ; i++ {
		for j := 3; j < WIDTH - 3; j++ {
			mapa[i][j] = ' '
		}
	}
	
	for s, a := range policy{
		mapa[s.y][s.x] = get_move( a )
	}
	
	mapa[INIT_Y][INIT_X] = 'S'
	mapa[GOAL_Y][GOAL_X] = 'G'

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			fmt.Printf( " %c ", mapa[i][j] )
		}
		fmt.Println()
	}
	fmt.Println("\nPolitica: ", policy )
}

func main() {
	rand.Seed( time.Now().Unix() )

	_a := [ACT]A{ ARR, AB, IZQ, DER }
	_st := S{ x: GOAL_X, y: GOAL_Y }
	mdp := MDP{ a :_a, St: _st, R: RECOMPENSA, γ: γ  }

	Q := Q_learning( mdp, α, 1100, 500 )
	fmt.Println( "Cargando Politica..." )
	policy := optimal_policy( Q )

	print_policy( policy )
}