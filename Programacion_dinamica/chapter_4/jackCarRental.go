package main 

import (
		"math"
		"fmt"
		)

const( 
		MAX_AUTO = 20
		MAX_MOVE_AUTO = 5
		BANKRUPTCY = 0

		COST_MOVE_AUTO = -2
		CREDITED = 10

		REQUEST_A = 3
		RETURN_A = 3
		
		REQUEST_B = 4
		RETURN_B = 2

		γ = 0.9
	)

type A int8

type Brand_Office_jack struct {
	return_cars [MAX_AUTO+1]float32
	request_cars [MAX_AUTO+1]float32
}

type MDP_jackCarRental struct{
	s [MAX_AUTO+1][MAX_AUTO+1] float32
	a [1+MAX_MOVE_AUTO*2]A
	R int8
	γ float32
	p [2]Brand_Office_jack
}

func ( j *MDP_jackCarRental ) Min( num1 , num2 float32 ) float32 {
	if num1 < num2 {
		return num1	
	}else{
 		return num2
	}
}

func ( j *MDP_jackCarRental ) Max( num1 , num2 float32 ) float32 {
	if num1 < num2 {
		return num2	
	}else{
 		return num1
	}
}

func ( j *MDP_jackCarRental ) Abs_i( num int ) int {
	if num < 0{
		return -1*num
	}else{
		return num
	}
}

func ( j *MDP_jackCarRental ) Abs( num float32 ) float32 {
	if num < 0{
		return -1*num
	}else{
		return num
	}
}

func ( j *MDP_jackCarRental ) Factorial( n float32 ) float32 {
	if n == 0.0 {
		return 1
	}else{
		return n * j.Factorial( n - 1 )
	}
}

func ( j *MDP_jackCarRental ) Poisson( x , λ float32 ) float32 {
	return float32( math.Pow( float64( λ ), float64( x ) ) * math.Exp( -1*float64( λ ) ) ) /  j.Factorial( x )
}

func ( j *MDP_jackCarRental ) Init_Policy(  ) [MAX_AUTO+1][MAX_AUTO+1]int {
	return [MAX_AUTO+1][MAX_AUTO+1]int {}
}

func ( j *MDP_jackCarRental ) Init_A() {
	act := -5
	for i := 0; i <= MAX_MOVE_AUTO * 2; i++ {
		j.a[i] = A(act)
		act++
	}
}

func ( j *MDP_jackCarRental ) Init_Brand_Office() ( Brand_Office_jack, Brand_Office_jack ) {
	var brand_office_A Brand_Office_jack
	var brand_office_B Brand_Office_jack
	
	for i := 0; i <= MAX_AUTO; i++ {
	
		brand_office_A.request_cars[i] = j.Poisson( float32( i ) , REQUEST_A )
		brand_office_A.return_cars[i] = j.Poisson( float32( i ), RETURN_A )
	
		brand_office_B.request_cars[i] = j.Poisson( float32( i ) , REQUEST_B )
		brand_office_B.return_cars[i] = j.Poisson( float32( i ), RETURN_B )
	}

	return brand_office_A, brand_office_B
}

func ( jack *MDP_jackCarRental ) value_function( num_car_1, num_car_2 int, a A ) float32 {
	act := int( a )
	total_acts := jack.Abs( float32(act) ) * COST_MOVE_AUTO

	// BANKRUPTCY == 0

	// sucursal A 
	//j.p[0].return_cars[i]
	//jack.p[0].request_cars[i]

	// sucursal B
	//j.p[1].return_cars[i]
	//jack.p[1].request_cars[i]
	//num_car_1 num_car_2

	value := float32(0.0)
	recompensa_total := float32(0.0)

	for i := 0; i <= MAX_AUTO; i++ {
		for j := 0; j <= MAX_AUTO; j++ { // iteraciones que son los dias 
			
			tope_A := jack.Min( float32( num_car_1 - act ), MAX_AUTO )

			if tope_A > MAX_AUTO {
				tope_A = MAX_AUTO
			}			
			
			tope_B := jack.Min( float32(num_car_2 + act), MAX_AUTO )

			if tope_B > MAX_AUTO {
				tope_B = MAX_AUTO
			}			
			
			renta_A := jack.Min( float32( tope_A ), float32( i ) )
			
			renta_B := jack.Min( float32( tope_B ), float32( j ) )			

			recompensa_total = float32( ( renta_A + renta_B ) * CREDITED )

			aux := jack.p[0].request_cars[i] * jack.p[1].request_cars[i]

			for k := 0; k <= MAX_AUTO; k++ {
				for l := 0; l <= MAX_AUTO; l++ { // iteraciones que son las noches
			
					aux_car_A := jack.Abs_i( int( jack.Min( float32( tope_A + float32( k ) ), MAX_AUTO ) ) )
					
					aux_car_B := jack.Abs_i( int( jack.Min( float32( tope_B + float32( l ) ), MAX_AUTO ) ) )

					aux = aux * jack.p[0].return_cars[i] * jack.p[1].return_cars[j]

					if aux_car_A < 0 || aux_car_B < 0 || aux_car_A > 21 || aux_car_B > 21 {
						//fmt.Println(aux_car_A, aux_car_B )
					}

					value += aux * total_acts * ( recompensa_total + jack.γ * jack.s[aux_car_A][aux_car_B] )
				}
			}
		}
	}

	return value
}

func ( jack *MDP_jackCarRental ) Init_S() {
	for i := 0; i <= MAX_AUTO; i++ {
		for j := 0; j <= MAX_AUTO; j++ {
			jack.s[i][j] = float32( 0.0 )
		}
	}
}

func ( jack *MDP_jackCarRental ) run( θ float32 ) [MAX_AUTO+1][MAX_AUTO+1]int {
	
	jack.p[0], jack.p[1] = jack.Init_Brand_Office()
	π := jack.Init_Policy() 
	jack.Init_S()
	jack.Init_A()
	var ν float32 = 0.0
	var Δ float32 = 0.0
	var policy_stable bool = true
	
	for policy_stable {
		// Policy Evaluation 
		for {
			Δ = 0.0
			for i := 0; i <= MAX_AUTO; i++ { 
				for j :=  0; j <= MAX_AUTO; j++ {
					ν = jack.s[i][j]
					
					jack.s[i][j] = jack.value_function( i, j, A(π[i][j]) )
		 			
					Δ = jack.Max( Δ, jack.Abs( ν - jack.s[i][j] ) )
				}
			}

			if Δ < θ {
				break
			}
		}

		// Policy Improvement
		for i := 0; i <= MAX_AUTO; i++ {
			for j := 0; j <= MAX_AUTO; j++ {
		
				old_action := π[i][j]
		
				π[i][j] = jack.max_a( i, j )
		
				if old_action != π[i][j] {
					policy_stable = false
				}
			}
		}

		if policy_stable {
			return π
		}
	}
	return π
}

func ( jack *MDP_jackCarRental ) max_a( i , j int ) int {
	v := jack.s[i][j]
	var accion int = 0.0

	for i := 0; i <= MAX_MOVE_AUTO*2; i++ {
				
		aux := jack.value_function( i, j, jack.a[i] )
		if aux > v {
			v = aux
			accion = int( jack.a[i] )
		}
	}
	return accion
}

func main() {
	
	//var jack MDP_jackCarRental

	jack := MDP_jackCarRental{ s: [MAX_AUTO+1][MAX_AUTO+1]float32{}, γ: γ, R:CREDITED, a: [1+MAX_MOVE_AUTO * 2]A{} }
	
	poli := jack.run( 0.000001 ) 

	fmt.Println( poli )

}