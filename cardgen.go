package main

// func dealCards (players [3]*Player) {
// 	var deck [32]Card

// 	c := 0
// 	for i := Seven; i <= Ace; i++ {
// 		for j := Spades; j <= Clubs; j++ {
// 			deck[c] = Card{suit: j, value: i}
// 			c++
// 		}
// 	}

// 	for i := 31; i > 0; i-- {
// 		j := rand.IntN(i)

// 		deck[i], deck[j] = deck[j], deck[i]
// 	}

// 	copy(players[0].hand[:], deck[0:10])
// 	copy(players[1].hand[:], deck[10:20])
// 	copy(players[2].hand[:], deck[20:30])

// 	return
// }
