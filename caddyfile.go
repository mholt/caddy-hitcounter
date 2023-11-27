package hitcounter

import (
	"strconv"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// UnmarshalCaddyfile deserializes Caddyfile tokens into h.
// The style names are subject to change!
//
//	hitCounter {
//	    style green|bright_green|odometer|yellow
//	    pad_digits <num_digits>
//	}
func (hc *HitCounter) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "style":
				if !d.NextArg() {
					return d.ArgErr()
				}
				hc.Style = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "pad_digits":
				if !d.NextArg() {
					return d.ArgErr()
				}
				num, err := strconv.Atoi(d.Val())
				if err != nil {
					return d.Errf("invalid digit padding number '%s': %v", d.Val(), err)
				}
				hc.PadDigits = num
				if d.NextArg() {
					return d.ArgErr()
				}

			default:
				return d.Errf("unrecognized hitCounter config property: %s", d.Val())
			}
		}
	}
	return nil
}
