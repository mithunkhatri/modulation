package audio

import (
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type Player struct {
	mu          sync.Mutex
	ctrl        *beep.Ctrl
	streamer    beep.StreamCloser
	format      beep.Format
	done        chan bool
	isPlaying   bool
	currentURL  string
	stationName string
	cmd         *exec.Cmd
}

func NewPlayer() (*Player, error) {
	// Initialize speaker with a default sample rate.
	// We'll update it when we actually start a stream.
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return nil, err
	}

	return &Player{
		done: make(chan bool),
	}, nil
}

func (p *Player) Play(url string, name string) error {
	// Stop existing command if any
	p.mu.Lock()
	p.stop()
	p.mu.Unlock()

	// Use ffmpeg to decode stream to PCM
	cmd := exec.Command("ffmpeg",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-i", url,
		"-f", "s16le",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		"-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	streamer := &pcmStreamer{
		reader: stdout,
		cmd:    cmd,
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.streamer = streamer
	p.format = format
	p.cmd = cmd
	p.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	p.currentURL = url
	p.stationName = name
	p.isPlaying = true

	speaker.Play(p.ctrl)
	return nil
}

type pcmStreamer struct {
	reader io.ReadCloser
	cmd    *exec.Cmd
}

func (s *pcmStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	buf := make([]byte, len(samples)*4) // 2 channels * 2 bytes per sample
	nRead, err := io.ReadFull(s.reader, buf)
	if nRead > 0 {
		for i := 0; i < nRead/4; i++ {
			lowL := float64(int16(uint16(buf[i*4]) | uint16(buf[i*4+1])<<8))
			lowR := float64(int16(uint16(buf[i*4+2]) | uint16(buf[i*4+3])<<8))
			samples[i][0] = lowL / 32768.0
			samples[i][1] = lowR / 32768.0
		}
		return nRead / 4, true
	}
	if err != nil {
		return 0, false
	}
	return 0, true
}

func (s *pcmStreamer) Err() error {
	return nil
}

func (s *pcmStreamer) Close() error {
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	return s.reader.Close()
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ctrl != nil {
		p.ctrl.Paused = !p.ctrl.Paused
	}
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stop()
}

func (p *Player) stop() {
	if p.streamer != nil {
		p.streamer.Close()
		p.streamer = nil
	}
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd = nil
	}
	p.isPlaying = false
	p.stationName = ""
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isPlaying && p.ctrl != nil && !p.ctrl.Paused
}

func (p *Player) CurrentStation() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stationName
}

func (p *Player) IsCurrentStation(url string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentURL == url
}
