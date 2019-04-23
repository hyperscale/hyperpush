// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package packets

import (
	"fmt"
	"io"
)

//DisconnectPacket is an internal representation of the fields of the
//Disconnect MQTT packet
type DisconnectPacket struct {
	FixedHeader
}

// NewDisconnectPacket return DisconnectPacket
func NewDisconnectPacket() *DisconnectPacket {
	return NewControlPacket(Disconnect).(*DisconnectPacket)
}

func (d *DisconnectPacket) String() string {
	str := fmt.Sprintf("%s", d.FixedHeader)
	return str
}

func (d *DisconnectPacket) Write(w io.Writer) error {
	packet := d.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (d *DisconnectPacket) Unpack(b io.Reader) error {
	return nil
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (d *DisconnectPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}
