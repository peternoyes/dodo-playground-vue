package main

import (
	"syscall/js"
)

type WebSpeaker struct {
	Context    js.Value
	Oscillator js.Value
	Gain       js.Value
	Mute       bool
	Enabled    bool
}

func (s *WebSpeaker) New() {
	v := js.Global().Get("window").Get("AudioContext")
	if v.IsUndefined() {
		v = js.Global().Get("window").Get("webkitAudioContext")
	}

	s.Context = v.New()
	s.Gain = s.Context.Call("createGain")
	s.Mute = false
	s.Enabled = false
}

func (s *WebSpeaker) SetFrequency(freq int) {
	if !s.Enabled || s.Mute {
		return
	}

	if freq == 0 {
		s.Gain.Get("gain").Set("value", 0)
		return
	}

	if !s.Oscillator.IsNull() {
		s.Oscillator.Call("stop")
	}

	s.Oscillator = s.Context.Call("createOscillator")
	s.Oscillator.Set("type", "square")
	s.Oscillator.Get("frequency").Set("value", freq)
	s.Gain.Get("gain").Set("value", 0.05)

	s.Oscillator.Call("connect", s.Gain)
	s.Gain.Call("connect", s.Context.Get("destination"))

	s.Oscillator.Call("start")
}

func (s *WebSpeaker) Enable() {
	s.Enabled = true
}

func (s *WebSpeaker) Disable() {
	s.Gain.Get("gain").Set("value", 0)
	s.Enabled = false
}

func (s *WebSpeaker) ToggleMute() {
	s.Mute = !s.Mute

	if s.Mute {
		s.Gain.Get("gain").Set("value", 0)
	}
}
