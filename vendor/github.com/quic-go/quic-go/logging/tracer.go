package logging

import "net"

// A Tracer traces events.
type Tracer struct {
	SentPacket                   func(net.Addr, *Header, ByteCount, []Frame)
	SentVersionNegotiationPacket func(_ net.Addr, dest, src ArbitraryLenConnectionID, _ []Version)
	DroppedPacket                func(net.Addr, PacketType, ByteCount, PacketDropReason)
	Debug                        func(name, msg string)
	Close                        func()
}

// NewMultiplexedTracer creates a new tracer that multiplexes events to multiple tracers.
func NewMultiplexedTracer(tracers ...*Tracer) *Tracer {
	if len(tracers) == 0 {
		return nil
	}
	if len(tracers) == 1 {
		return tracers[0]
	}
	return &Tracer{
		SentPacket: func(remote net.Addr, hdr *Header, size ByteCount, frames []Frame) {
			for _, t := range tracers {
				if t.SentPacket != nil {
					t.SentPacket(remote, hdr, size, frames)
				}
			}
		},
		SentVersionNegotiationPacket: func(remote net.Addr, dest, src ArbitraryLenConnectionID, versions []Version) {
			for _, t := range tracers {
				if t.SentVersionNegotiationPacket != nil {
					t.SentVersionNegotiationPacket(remote, dest, src, versions)
				}
			}
		},
		DroppedPacket: func(remote net.Addr, typ PacketType, size ByteCount, reason PacketDropReason) {
			for _, t := range tracers {
				if t.DroppedPacket != nil {
					t.DroppedPacket(remote, typ, size, reason)
				}
			}
		},
		Debug: func(name, msg string) {
			for _, t := range tracers {
				if t.Debug != nil {
					t.Debug(name, msg)
				}
			}
		},
		Close: func() {
			for _, t := range tracers {
				if t.Close != nil {
					t.Close()
				}
			}
		},
	}
}
