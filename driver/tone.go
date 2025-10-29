package driver

import (
	"fmt"
	"math"
)

const (
	SampleRate = 48000
)

type Tone struct {
	Freq    float64
	Samples int
}

var ch = make(chan []Tone)
var currentTone Tone
var currentMalody []Tone
var phase, step float64

func bin2freq(b byte) float64 {
	n := float64(b-0x7B) / 12
	return 440 * math.Pow(2, n)
}

func bin2samp(b byte) int {
	return (SampleRate / 1000) * 8 * int(b) // 1 unit means 8 milisecond
}

func PlayTone(malodyBin []byte) {
	malody := []Tone{}
	loopLoc := 0
	loopCount := 0

	for i := 0; i < len(malodyBin); i += 2 {
		if malodyBin[i] == 0x05 {
			loopLoc = i
			loopCount = int(malodyBin[i+1])
		} else if malodyBin[i] == 0x06 {
			if loopCount > 0 {
				loopCount -= 1
				i = loopLoc
			} else {
				i -= 1 // 0x06 is only a single byte
			}
		} else if malodyBin[i] == 0x40 {
			var freq float64 = 0
			samp := bin2samp(malodyBin[i+1])
			malody = append(malody, Tone{freq, samp})
		} else {
			freq := bin2freq(malodyBin[i])
			samp := bin2samp(malodyBin[i+1])
			malody = append(malody, Tone{freq, samp})
		}
	}
	fmt.Printf("PLAYTONE DRIVER: Malody: %v\n", malody)
	go func() {
		ch <- malody
	}()
}

func updateStep() {
	step = 2 * math.Pi * currentTone.Freq / SampleRate
}

func getNextTone() {
	defer updateStep()           // after getting a new tone update step
	if currentTone.Samples > 0 { // keep playing current tone
		return
	}
	if len(currentMalody) != 0 { // get next tone from the malody
		currentTone = currentMalody[0]
		currentMalody = currentMalody[1:]
		return
	}
	select {
	case m := <-ch:
		// new malody recieved
		currentMalody = m
	default:
		// nothing to play, just fill 0s
		currentTone.Freq = -1
	}
}

func binarify(n float32) float32 {
	if n < 0 {
		return -1
	} else if n == 0 {
		return 0
	} else {
		return 1
	}
}

func GenerateWave(buffer []float32, frames int) {
	getNextTone()
	for i := range frames {
		if currentTone.Freq == -1 {
			buffer[i] = 0
			continue
		}
		if currentTone.Samples != 0 {
			if currentTone.Freq == 0 {
				buffer[i] = 0
			} else {
				buffer[i] = binarify(float32(math.Sin(phase))) // osilator frequency goes here
				phase += step
				if phase >= 2*math.Pi {
					phase -= 2 * math.Pi
				}
			}
			currentTone.Samples--
		}

		if currentTone.Samples == 0 {
			getNextTone()
		}
	}
}
