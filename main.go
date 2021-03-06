package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

type snake struct {
	x             int
	y             int
	snakeNitroLvl int //default/initially at 0, will be altered randomly
	snakeMoved    int
	targetStruck  int
}

//Implement one snake first and then follow on....
type snake2 struct {
	x           int
	y           int
	snacksEaten int
}

func drawSnakes(snk snake) {

	//For one snake.for both, initiate with the loops
	termbox.SetCell(snk.x, snk.y+1, '*', termbox.Attribute(3), termbox.Attribute(8))
	termbox.SetCell(snk.x, snk.y+2, '*', termbox.Attribute(3), termbox.Attribute(8))
	termbox.SetCell(snk.x, snk.y+3, '*', termbox.Attribute(3), termbox.Attribute(8))
	//termbox.SetCell(snk.x, snk.y+7, '/', termbox.Attribute(3), termbox.Attribute(8))
	//termbox.SetCell(snk.x, snk.y+8, '+', termbox.Attribute(3), termbox.Attribute(8))
}

//creates the entire map environment
func drawWorld(mapEnvDup [25]string) {
	getColour := func(x int, y int, mapEnvDup [25]string) int {
		switch mapEnvDup[y][x] {
		case ' ':
			return 7
		case '0':
			return 3
		case '=':
			return 4
		case '|':
			return 14
		case '-':
			return 13
		}
		return -1
	}
	getColour(0, 0, mapEnvDup)
	wid, hei := 70, 25
	for y := 0; y < hei; y++ {
		for x := 0; x < wid; x++ {
			termbox.SetCell(x, y, rune(mapEnvDup[y][x]), termbox.Attribute(getColour(x, y, mapEnvDup)), termbox.Attribute(getColour(x, y, mapEnvDup))) //sets the colour

		}
	}
}

func drawEnv(mapEnvDup [25]string, snk snake) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawWorld(mapEnvDup)
	drawSnakes(snk)

	termbox.Flush()

}

func main() {
	fmt.Println("making")

	mapEnv := [25]string{}

	mapEnv[0] = "|-------------------------------------------0-------------==----------------|"
	mapEnv[1] = "|                                                                             |"
	mapEnv[2] = "|                                                                             |"
	mapEnv[3] = "|                                                                             |"
	mapEnv[4] = "|                                                                             |"
	mapEnv[5] = "|                                                                             |"
	mapEnv[6] = "|                                                                             |"
	mapEnv[7] = "|                                                                             |"
	mapEnv[8] = "|                                                                             |"
	mapEnv[9] = "|                                  -------------                              |"
	mapEnv[10] = "|                                                                             |"
	mapEnv[11] = "|                                                                             |"
	mapEnv[12] = "|                                                                             |"
	mapEnv[13] = "|                                                                             |"
	mapEnv[14] = "|                                                                             |"
	mapEnv[15] = "|                                                                             |"
	mapEnv[16] = "|                                                                             |"
	mapEnv[17] = "|                                        -------------                        |"
	mapEnv[18] = "|                                                                             |"
	mapEnv[19] = "|                                                                             |"
	mapEnv[20] = "|                                                                             |"
	mapEnv[21] = "|                                                                             |"
	mapEnv[22] = "|                                                                             |"
	mapEnv[23] = "|                                                                             |"
	mapEnv[24] = "------------------------------------------------------------------------------"
	/*ASCII values
	- 45
	| 124
	0 48
	= 61
	*/
	err := termbox.Init()

	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(5000 * time.Millisecond)
	quit := make(chan string)

	//Go routine for even handler
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent() //polls event
		}
	}()

	redrawProcess := make(chan snake)
	go func(snkArg chan snake, mapEnvArg [25]string) {
		for {
			select {
			case val := <-snkArg:
				drawEnv(mapEnvArg, val)
			}
		}
	}(redrawProcess, mapEnv)

	//using map[y][x] method
	/*
	 */

	go func(t *time.Ticker, mainSnk chan snake, worldMap [25]string) {
		snk := snake{
			x:             3,
			y:             0,
			snakeNitroLvl: 0,
			snakeMoved:    0, //TBW

		}
		for {
			select {
			case <-t.C: //TBC
				if snk.snakeNitroLvl > 0 {
					if snk.y < 22 {
						if worldMap[snk.y][snk.x] == 32 { //triggered by space bar
							snk.y += 1
							mainSnk <- snk
						}
					}
					snk.snakeNitroLvl--
				} else if snk.snakeNitroLvl == 0 {
					if snk.y <= 22 {
						if worldMap[snk.y+3][snk.x] == 32 && snk.snakeMoved == 0 {
							snk.y++
							mainSnk <- snk
						}
					}
				}
			case event := <-eventQueue:
				if event.Type == termbox.EventKey {
					switch event.Key { // 70x25
					case termbox.KeyArrowUp:
						if snk.y > 0 {
							if worldMap[snk.y-1][snk.x] != 45 && snk.snakeNitroLvl == 0 {
								snk.y--
								mainSnk <- snk
							}
						}
					case termbox.KeyArrowDown:
						if snk.y < 20 {
							if worldMap[snk.y+1][snk.x] != 45 || worldMap[snk.y+2][snk.x] != 45 {
								snk.y++
								mainSnk <- snk
							}
						}
					case termbox.KeyArrowLeft: //for left check snks left move -
						if snk.x > 0 && snk.y > 0 && snk.y < 23 {
							c1 := worldMap[snk.y][snk.x-1]
							//c2 := worldMap[snk.y][snk.x-2] (c2 == 32 && c2 != 124) ||
							//c3 := worldMap[snk.y][snk.x-3]
							if c1 == 32 && c1 != 124 {
								snk.x--
								mainSnk <- snk
							}
						}
					case termbox.KeyArrowRight:
						if snk.x < 67 && snk.y > 0 && snk.y < 22 {
							if worldMap[snk.y][snk.x+1] == 32 && worldMap[snk.y][snk.x+2] == 32 && worldMap[snk.y][snk.x+3] == 32 {
								snk.x++
								mainSnk <- snk
							}
						}
					case termbox.KeySpace:
						snk.snakeNitroLvl++

					case termbox.KeyEsc:
						quit <- "Game has ended.... Thanks for playing."
					}
					snk.snakeMoved = 1
					mainSnk <- snk
				}
			}
			//checks here for end bit, hitting a pre recorded zone.
			if (snk.y == 0 && snk.x == 44) || (snk.y == 70 && snk.x == 44) {
				snk.targetStruck = 1
				mainSnk <- snk
				//fmt.Println("ad")
			}
			//fmt.Print( snk.x, "---",snk.y,"---",snk.targetStruck,"|||| ")
			if snk.targetStruck == 1 && snk.y == 0 && snk.x == 58 {
				mainSnk <- snk
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				t.Stop()
				quit <- "You win the game!!!"
			}
		}

	}(ticker, redrawProcess, mapEnv)
	msg := <-quit
	termbox.Close()
	fmt.Println(msg)

}
