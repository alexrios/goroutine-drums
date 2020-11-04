package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/fatih/color"
)

var (
	hihatSound = sound("drums/hihat.wav")
	kickSound  = sound("drums/kick.wav")
	snareSound = sound("drums/snare.wav")
)

type Drum struct {
	Tempo int
}

func (d Drum) hihat(rhythm string, wg *sync.WaitGroup) {
	defer wg.Done()
	d.playbeat("prato", rhythm, hihatSound, blue)
}

func (d Drum) kick(rhythm string, wg *sync.WaitGroup) {
	defer wg.Done()
	d.playbeat("bumbo", rhythm, kickSound, green)
}

func (d Drum) snare(rhythm string, wg *sync.WaitGroup) {
	defer wg.Done()
	d.playbeat("caixa", rhythm, snareSound, yellow)
}

func (d Drum) playbeat(name string, beats string, sound beep.StreamSeeker, color *color.Color) {
	ticker := time.NewTicker(time.Duration(d.Tempo) * time.Millisecond)
	defer ticker.Stop()

	runes := []rune(beats)
	count := 0
	loop := 1
	for count < len(runes) {
		select {
		case _ = <-ticker.C:
			r := runes[count]
			if r == 'x' {
				go play(sound)
				color.Printf("[T%d]%s\n", loop, name)
			}
			count++
			loop++
		}
	}
	//Give a little time to the last note to echo
	time.Sleep(time.Duration(d.Tempo) * time.Millisecond)
}

func play(sound beep.StreamSeeker) {
	playbackDone := make(chan bool)
	speaker.Play(beep.Seq(sound, beep.Callback(func() {
		playbackDone <- true
		sound.Seek(0)
	})))
	<-playbackDone
}

func sound(filename string) beep.StreamSeeker {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	sound := buffer.Streamer(0, buffer.Len())

	return sound
}
