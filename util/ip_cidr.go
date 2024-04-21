package util

import (
	"github.com/pic4xiu/iprange"
)

func CheckCidrIPs(cidr string) <-chan *string {
	ch := make(chan *string)
	list, err := iprange.ParseList(cidr)
	if err != nil {
		return nil
	}
	go func() {
		for _, rng := range list.Expand() {
			ing := rng.String()
			ch <- &ing
		}
		close(ch)
	}()
	return ch
}
