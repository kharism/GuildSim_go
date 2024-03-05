package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	_ "embed"
)

type musicType int

const (
	sampleRate = 48000
)
const (
	typeOgg musicType = iota
	typeMP3
)

func (t musicType) String() string {
	switch t {
	case typeOgg:
		return "Ogg"
	case typeMP3:
		return "MP3"
	default:
		panic("not reached")
	}
}

type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte

	volume128 int

	musicType musicType
}

//go:embed sound/638697__captainyulef__card.wav
var drawSound []byte

//go:embed sound/sword-hit-7160.mp3
var damageSound []byte

//go:embed sound/98491__tec_studio__walking1_gravel.mp3
var playSound []byte

//go:embed sound/breeze-of-blood-122253_short.mp3
var defeatSound []byte

// which SE we should play
var seSignal chan string

var SFX_DRAW = "DRAW"
var SFX_ATTACKED = "ATK"
var SFX_PLAY = "PLAY"
var SFX_DEFEAT = "DEFEAT"

func init() {
	seSignal = make(chan string, 7)
	go func() {
		for s := range seSignal {
			switch s {
			case SFX_DRAW:
				s, _ := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(drawSound))
				b, _ := io.ReadAll(s)
				sePlayer := audio.CurrentContext().NewPlayerFromBytes(b)
				sePlayer.Play()
				time.Sleep(time.Second / 7)
			case SFX_ATTACKED:
				s, _ := mp3.DecodeWithoutResampling(bytes.NewReader(damageSound))
				b, _ := io.ReadAll(s)
				sePlayer := audio.CurrentContext().NewPlayerFromBytes(b)
				sePlayer.Play()
				// time.Sleep(time.Second / 7)
			case SFX_PLAY:
				s, _ := mp3.DecodeWithoutResampling(bytes.NewReader(playSound))
				b, _ := io.ReadAll(s)
				sePlayer := audio.CurrentContext().NewPlayerFromBytes(b)
				sePlayer.Play()
			case SFX_DEFEAT:
				s, _ := mp3.DecodeWithoutResampling(bytes.NewReader(defeatSound))
				b, _ := io.ReadAll(s)
				sePlayer := audio.CurrentContext().NewPlayerFromBytes(b)
				sePlayer.Play()
			}
		}
	}()
}
func NewPlayer(audioContext *audio.Context, audioPath string, musicType musicType) (*Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream
	audio, err := os.Open(audioPath)
	if err != nil {
		return nil, err
	}
	switch musicType {

	case typeMP3:
		var err error
		s, err = mp3.DecodeWithoutResampling(audio)
		if err != nil {
			return nil, err
		}
	default:
		panic("not reached")
	}
	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		seCh:         make(chan []byte),

		musicType: musicType,
	}
	if player.total == 0 {
		player.total = 1
	}

	const buttonPadding = 16

	player.audioPlayer.Play()
	go func() {
		// s, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(raudio.Jab_wav))
		// if err != nil {
		// 	log.Fatal(err)
		// 	return
		// }
		s := bytes.NewReader(drawSound)
		b, err := io.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		player.seCh <- b
	}()
	return player, nil
}

func (p *Player) Close() error {
	return p.audioPlayer.Close()
}
func (p *Player) update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}
	p.playSEIfNeeded()
	return nil
}
func (p *Player) shouldPlaySE() bool {
	if p.seBytes == nil {
		// Bytes for the SE is not loaded yet.
		return false
	}
	return false
}

func (p *Player) playSEIfNeeded() {
	// if !p.shouldPlaySE() {
	// 	return
	// }
	// sePlayer := p.audioContext.NewPlayerFromBytes(p.seBytes)
	// sePlayer.Play()

}
