package main

import (
	"bytes"
	"encoding/binary"
	//"fmt"
	"io"
)

/*var events []string = []string{
	"request",
	"impression",
	"creativeView",
	"start",
	"firstQuartile",
	"midpoint",
	"thirdQuartile",
		"complete",
}
var elen int = len(events)

func mkstring (s string) string {
	l := len(s)
	return string(byte(l)) + s
}

// generate n metric names
func genMetricNames(prefix string, id, n int) []string {
	names := make([]string, n * n * n * elen)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				for h := 0; h < elen; h++ {
					names[i * n * n * elen + j * n * elen + k * elen + h] =
						fmt.Sprintf("domain.domain%5.5d.aid%5.5d.cmp%5.5d.event.%s", i, j, k, events[h])
				}
			}
		}
	}

	return names
}*/

// generate n metric names
/* func genMetricNames(prefix string, id, n int) []string {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("\x05agent\x06%6d\x07metrics\x06%6d", id, i)
	}

	return names
} */

// actually write the data in carbon line format
func carbonate(w io.ReadWriteCloser, epoch int64, agent string, domain string, aid string, cmp string, event string) error {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(5))
	binary.Write(buf, binary.BigEndian, uint64(epoch))

	slen := len(agent) + 1 + /*len(domain) + 1*/ + len(aid) + 1 + /*len(cmp) + 1*/ + len(event) + 1
	binary.Write(buf, binary.BigEndian, uint16(slen))
	w.Write(buf.Bytes())

	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(len(agent)))
	w.Write(buf.Bytes())
	w.Write([]byte(agent))

	/*buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(len(domain)))
	w.Write(buf.Bytes())
	w.Write([]byte(domain))*/

	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(len(aid)))
	w.Write(buf.Bytes())
	w.Write([]byte(aid))

	/*buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(len(cmp)))
	w.Write(buf.Bytes())
	w.Write([]byte(cmp))*/

	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint8(len(event)))
	w.Write(buf.Bytes())
	w.Write([]byte(event))

	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(8))
	binary.Write(buf, binary.BigEndian, uint8(1))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	binary.Write(buf, binary.BigEndian, uint8(0))
	value := 1
	if epoch % 10 == 0 {
		value = 0
	}
	binary.Write(buf, binary.BigEndian, uint32(value)) // value
	w.Write(buf.Bytes())
	return nil
}
