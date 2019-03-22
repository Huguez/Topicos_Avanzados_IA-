package main

import (
		"fmt"
		"math/rand"
		"time"
		)

const( 
		INIT_X = 10
		INIT_Y = 2

		GOAL_X = 15
		GOAL_Y = 30

		EPSILON = 0.1
		ALPHA = 0.5
		GAMMA = 1.0

		ACT = 8
		ANCHO = 46
		ALTO = 24

		ARR = 0
		AB = 1
		I = 2
		D = 3
		AI = 4
		AD = 5
		ABI = 6
		ABD = 7
		STAY = 8
)

type A int

type S struct{
	x int 
	y int
}

type MDP struct {
	wind [ANCHO]int
	a [ACT]A
	St S
	R int
	G float32
}

func ( m *MDP ) is_terminal( s S ) bool { return s.x == m.St.x && s.y == m.St.y }

func ( m *MDP ) get_action( s S, Q [ALTO][ANCHO][ACT]float32  ) A {
	
	random := rand.Float64()
	if random < 1.0 - EPSILON {
		a , i := float32( -100.0 ), -1
		
		for i_act , act := range Q[s.x][s.y] {
			if float32(act) > float32(a) {
				a = float32(act)
				i = i_act		
			}
		}
		
		return A( i )
	}else{
		act_rand := rand.Intn( ACT )
			
		return m.a[act_rand]
	}
}

func ( m *MDP ) take_action( a A, s S ) (S, int) {
	var _s S
	var r int = m.R
	
	var aux int = m.wind[s.x]

	switch a {
	case ARR:
		if s.x - 1 - aux < 0 { ///// aqui
			_s.x = s.x
		}else{
			_s.x = s.x - 1 - aux 
		}

		_s.y = s.y
	case AB:

		if s.x + 1 - aux >= ALTO { ///// aqui
			_s.x = s.x
		}else{
			_s.x = s.x + 1 - aux	
		}

		_s.y = s.y
	case I:
		if  s.x - aux >= 0{
			_s.x = s.x - aux 
		}else{
			_s.x = 0
		}

		if  s.y - 1 < 0 {
			_s.y = s.y
		}else{
			_s.y = s.y - 1
		}
	case D:
		if  s.x - aux >= 0{
			_s.x = s.x - aux 
		}else{
			_s.x = 0
		}

		if s.y + 1 >= ALTO {
			_s.y = s.y
		}else{
			_s.y = s.y + 1
		}
	case AI:
		if s.x - 1 - aux < 0 {
			_s.x = s.x
		}else{
			_s.x = s.x - 1 - aux
		}

		if s.y - 1 - aux < 0 {
			_s.y = s.y
		}else{
			_s.y = s.y - 1 - aux
		
		}
	case AD:
		if s.x - 1 - aux < 0 {
			_s.x = s.x
		}else{
			_s.x = s.x - 1 - aux
		}

		if s.y + 1 >= ANCHO {
			_s.y = s.y
		}else{
			_s.y = s.y + 1
		}
	case ABI:
		if s.y - 1 < 0 {
			_s.y = s.y
		}else{
			_s.y = s.y - 1
		}

		if s.x + 1 - aux >= ALTO {
			_s.x = s.x
		}else{
			_s.x = s.x + 1  - aux
		}
	case ABD:
		if s.x + 1  - aux >= ALTO {
			_s.x = s.x
		}else{
			_s.x = s.x + 1 - aux
		}
		
		if s.y + 1 >= ANCHO {
			_s.y = s.y
		}else{
			_s.y = s.y + 1
		}
	}

	if m.is_terminal( _s ){
		r = 0
	}

	return _s, r
}

func ( mdp *MDP ) init_Q() [ALTO][ANCHO][ACT]float32 {

	var aux_Q [ALTO][ANCHO][ACT]float32	
	
	for i := 0; i < ALTO; i++ {
		for j := 0; j < ANCHO; j++ {
			for k := 0; k < ACT; k++ {
				aux_Q[i][j][k] = 0.0
				//aux_Q[i][j][k] = -float32(rand.Float64())
			}
		}
	}
	
	return aux_Q
}

func ( mdp *MDP ) init_s() (int, int) {	return INIT_X, INIT_Y }

func sarsa( mdp MDP, alpha float32, max_episodios int, max_pasos int ) [ALTO][ANCHO][ACT]float32 {
	Q := mdp.init_Q()

	var r int
	var s, n_s S 
	var a, n_a A

	for i := 0; i < max_episodios; i++ {
		
		s.x, s.y = mdp.init_s()
		a = mdp.get_action( s, Q )

		for j := 0; j < max_pasos; j++ {
			
			// hacer accion, tener nuevo estado y recompensa
			n_s, r = mdp.take_action( a, s )
			
			// obtener nueva accion apartir del nuevo estado con recompensa
			n_a = mdp.get_action( n_s, Q )
			
			Q[s.x][s.y][a] = Q[s.x][s.y][a] + float32(alpha)*( float32(r) + float32(mdp.G)*Q[n_s.x][n_s.y][n_a] - Q[s.x][s.y][a] )
			
			s.x, s.y = n_s.x, n_s.y 
			a = n_a

			if mdp.is_terminal( s ){
				break
			}
		}
	}

	return Q
}

func wind() [ANCHO] int {
	var wind [ANCHO] int
	for i := ANCHO/2; i < ANCHO; i++ {
		wind[i] = rand.Intn( 3 ) + 1
	}
	return wind
}

func get_move( a A ) rune {
	var char rune
	
	switch a {
	case ARR:
		char = '↑'
	case AB:
		char = '↓'
	case I:
		char = '←'
	case D:
		char = '→'
	case AI:
		char = '↖'
	case AD:
		char = '↗'
	case ABI:
		char = '↙'
	case ABD:
		char = '↘'
	case STAY:
		 	char = ' '
	}
	return char
}	

func print_gridworld( policy map[S]A, wind [ANCHO]int ) {
	var gridworld [ALTO][ANCHO] rune
	
	for i := 0; i < ALTO; i++ {
		for j := 0; j < ANCHO/2; j++ {
			gridworld[i][j] = '-'
		}
		for j := ANCHO/2; j < ANCHO; j++ {
			gridworld[i][j] = '𐋇'
		}
	}

	for s, a := range policy {
		gridworld[s.x][s.y] = get_move( a )
	}

	gridworld[GOAL_X][GOAL_Y] = 'F'
	gridworld[INIT_X][INIT_Y] = 'I'

	for i := 0; i < ALTO; i++ {
		for j := 0; j < ANCHO; j++ {
			fmt.Printf( " %c ",gridworld[i][j] )
		}
		fmt.Println()
	}
	for i := 0; i < ANCHO; i++ {
		fmt.Printf( " %d ", wind[i] )
	}
	fmt.Println()
}

func get_state( s S, a A ) (S, int) {
	var _s S
	var aux int = 0

	switch a {
	case ARR:
		if s.x - 1 - aux < 0 { ///// aqui
			_s.x = s.x
		}else{
			_s.x = s.x - 1 - aux 
		}

		_s.y = s.y
	case AB:

		if s.x + 1 - aux >= ALTO { ///// aqui
			_s.x = s.x
		}else{
			_s.x = s.x + 1 - aux	
		}

		_s.y = s.y
	case I:
		if  s.x - aux >= 0{
			_s.x = s.x - aux 
		}else{
			_s.x = 0
		}

		if  s.y - 1 < 0 {
			_s.y = s.y
		}else{
			_s.y = s.y - 1
		}
	case D:
		if  s.x - aux >= 0{
			_s.x = s.x - aux 
		}else{
			_s.x = 0
		}

		if s.y + 1 >= ALTO {
			_s.y = s.y
		}else{
			_s.y = s.y + 1
		}
	case AI:
		if s.x - 1 - aux < 0 {
			_s.x = s.x
		}else{
			_s.x = s.x - 1 - aux
		}

		if s.y - 1 - aux < 0 {
			_s.y = s.y
		}else{
			_s.y = s.y - 1 - aux
		
		}
	case AD:
		if s.x - 1 - aux < 0 {
			_s.x = s.x
		}else{
			_s.x = s.x - 1 - aux
		}

		if s.y + 1 >= ANCHO {
			_s.y = s.y
		}else{
			_s.y = s.y + 1
		}
	case ABI:
		if s.y - 1 < 0 {
			_s.y = s.y
		}else{
			_s.y = s.y - 1
		}

		if s.x + 1 - aux >= ALTO {
			_s.x = s.x
		}else{
			_s.x = s.x + 1  - aux
		}
	case ABD:
		if s.x + 1  - aux >= ALTO {
			_s.x = s.x
		}else{
			_s.x = s.x + 1 - aux
		}
		
		if s.y + 1 >= ANCHO {
			_s.y = s.y
		}else{
			_s.y = s.y + 1
		}
	}
	return _s, 0
}

func politica_optima( Q [ALTO][ANCHO][ACT]float32 ) map[S]A {
	
	policy := make( map[S]A )
	var (
		a A = -1
		mayor float32 = -1000.00
		meta bool = false
	)

	s := S{ x: INIT_X, y: INIT_Y  }

 	for !meta {
 		for act, value := range Q[s.x][s.y] {
			if mayor < value{
				mayor = value
				a = A(act)
			}
		}
		policy[s] = a
		s, _ = get_state( s , a )
		if s.x == GOAL_X && s.y == GOAL_Y{
			meta = true
		} 
 	}
 	fmt.Println( "GOAL_X", s.x, "GOAL_Y", s.y )
 	return policy
}

func main() {
	rand.Seed( time.Now().Unix() )

	_a := [ACT]A { ARR, AB, I, D, AI, AD, ABI, ABD }
	_st :=  S{ x: GOAL_X, y : GOAL_Y }
	_wind := wind()

	mdp := MDP{ wind: _wind, a: _a, St: _st , R: -1 , G: GAMMA }

	Q := sarsa( mdp, ALPHA, 17000, 8000 )
	
	policy := politica_optima( Q )

	print_gridworld( policy, _wind )
}