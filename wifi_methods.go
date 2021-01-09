// Custom methods not part of the generated set

package hitron

import (
	"context"
	"strconv"
)

// WiFiRadioDetails - get details from /WiFi/Radios/<n>
func (c *CableModem) WiFiRadioDetails(ctx context.Context, radio int) (out WiFiRadio, err error) {
	err = c.getJSON(ctx, "/WiFi/Radios/"+strconv.Itoa(radio), &out)

	return out, err
}
