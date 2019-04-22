package main

import (
		"fmt"
//		"math/rand"
//		"time" )
				)

const(	META = 100      // estado terminal
 		STATE_INIT = 1  // estado inicial 
 		EPSILON	= 0.09	// UN NUMERO MUY PEQUE
 		PH = 0.4 	    // probabilidad de la moneda
 		ALPHA = 0.5     // factor de aprendizaje
 		R = 1			// recompensa
 		GAMMA = 1 )     // factor de olvido     



type MDP_Gambler struct { 
	politica map[int]int
	V [META] float32
}

func init_policy()  map[int]int {
	politica := make( map[int]int )

	for i := 0; i < META; i++ {
		politica[i] = 0
	}

	return politica
}

func int_V() [META]float32 {
	var aux [META] float32
	return aux  
}

func ( g *MDP_Gambler ) max( num1, num2 float32 ) float32 {
	if num1 > num2 {
		return num1
	}else{
		return num2
	}
}

func ( g *MDP_Gambler ) min( num1, num2 int ) int {
	if num1 > num2 {
		return num2
	}else{
		return num1
	}
}

func ( g *MDP_Gambler ) abs( num1 float32 ) float32 {
	if num1 < 0.0 {
		return -1.0*num1
	}else{
		return num1
	}
}

func ( g *MDP_Gambler ) stake( s int ) (float32, int ){ // las acciones son las apuestas
	
	var top_apuesta int = g.min( s, META - s )
	var resultado float32 = 0.0
	var a_aux, r int = 1, 0
	var max float32 = PH * g.V[s] + ( 1.0 - PH ) * g.V[s] 
	
	for i := 0; i < top_apuesta; i++ { // el jugador puede apostar todo su capital hasta su tope
		
		if i+s == META-1 {
			r = R
		}else{
			r = 0 
		}

		resultado = PH * ( g.V[s+i] + float32(r) ) + ( 1.0 - PH ) * g.V[s-i]

		if resultado > max {
			max = resultado
			a_aux = i
		}
	}
	return max, a_aux
}

func ( g *MDP_Gambler ) run( teta float32 ) {
	
	var delta float32
	var v, n_v float32 = 0.0, 0.0
	
	for {
		
		delta = 0.0
		
		for state := 0; state < META ; state++ {
			v = g.V[state]

			g.V[state], g.politica[state] = g.stake( state )
			
			n_v = g.V[state]

			delta = g.max( delta, g.abs( v - n_v ) )
		}

		if delta < teta {
			break
		}
	}
}

func main() {
	
	jugador := MDP_Gambler{ politica : init_policy(), V: int_V() }
	
	jugador.run( 0.001 )

	for k, v := range jugador.politica { 
	    fmt.Printf("[%d]-> %d  \n", k, v)
	}

	fmt.Println()
}