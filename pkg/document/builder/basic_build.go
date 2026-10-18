package builder

import (
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
)

type (
	pulseCommand     = schema.PulseCommand
	beepCommand      = schema.BeepCommand
	feedCommand      = schema.FeedCommand
	cutCommand       = schema.CutCommand
	separatorCommand = schema.SeparatorCommand
)

func intPtr(v int) *int {
	return &v
}

// Feed adds paper feed lines
func (b *DocumentBuilder) Feed(lines int) *DocumentBuilder {
	return b.addCommand("feed", feedCommand{Lines: lines})
}

// Cut adds a paper cut command
func (b *DocumentBuilder) Cut() *DocumentBuilder {
	return b.addCommand("cut", cutCommand{Mode: "partial", Feed: intPtr(2)})
}

// FullCut adds a full paper cut
func (b *DocumentBuilder) FullCut() *DocumentBuilder {
	return b.addCommand("cut", cutCommand{Mode: "full", Feed: intPtr(2)})
}

// CutWithFeed adds a cut with custom feed
func (b *DocumentBuilder) CutWithFeed(mode string, feed int) *DocumentBuilder {
	return b.addCommand("cut", cutCommand{Mode: mode, Feed: intPtr(feed)})
}

// Separator adds a separator line
func (b *DocumentBuilder) Separator(char string) *DocumentBuilder {
	return b.addCommand("separator", separatorCommand{Char: char, Length: 48})
}

// SeparatorWithLength adds a separator with custom length
func (b *DocumentBuilder) SeparatorWithLength(char string, length int) *DocumentBuilder {
	return b.addCommand("separator", separatorCommand{Char: char, Length: length})
}

// Pulse opens the cash drawer
func (b *DocumentBuilder) Pulse() *DocumentBuilder {
	return b.addCommand("pulse", pulseCommand{
		Pin:     constants.DefaultPulsePin,
		OnTime:  constants.DefaultPulseOnTime,
		OffTime: constants.DefaultPulseOffTime,
	})
}

// PulseWithOptions opens the cash drawer with custom timing
func (b *DocumentBuilder) PulseWithOptions(pin, onTime, offTime int) *DocumentBuilder {
	return b.addCommand("pulse", pulseCommand{Pin: pin, OnTime: onTime, OffTime: offTime})
}

// Beep emits beep sounds
func (b *DocumentBuilder) Beep(times, lapse int) *DocumentBuilder {
	return b.addCommand("beep", beepCommand{Times: times, Lapse: lapse})
}

// BeepOnce emits a single beep sound
func (b *DocumentBuilder) BeepOnce() *DocumentBuilder {
	return b.addCommand("beep", beepCommand{
		Times: constants.DefaultBeepTimes,
		Lapse: constants.DefaultBeepLapse,
	})
}
