//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package common

import (
	"bytes"
	"io"
)

type HostAddress string

func (addr *HostAddress) IsLocalHost() bool {
	return addr != nil && len(*addr) == 0
}

func (addr *HostAddress) Equals(o HostAddress) bool {
	return addr != nil && *addr == o
}

func (addr *HostAddress) EqualsToString(o string) bool {
	return addr.Equals(HostAddress(o))
}

type HostIdentityHolder interface {
	GetHostAddress() HostAddress
	GetTransportKey() SignatureKeyHolder
	GetTransportCert() CertificateHolder
}

var _ HostIdentityHolder = &HostIdentity{}

func NewHostIdentityFromHolder(h HostIdentityHolder) HostIdentity {
	return HostIdentity{
		Addr: h.GetHostAddress(),
		Key:  h.GetTransportKey(),
		Cert: h.GetTransportCert(),
	}
}

type HostIdentity struct {
	Addr HostAddress
	Key  SignatureKeyHolder
	Cert CertificateHolder
}

func (v *HostIdentity) GetHostAddress() HostAddress {
	return v.Addr
}

func (v *HostIdentity) GetTransportKey() SignatureKeyHolder {
	return v.Key
}

func (v *HostIdentity) GetTransportCert() CertificateHolder {
	return v.Cert
}

type Foldable interface {
	FoldToUint64() uint64
	FixedByteSize() int
}

type FoldableReader interface {
	io.WriterTo
	io.Reader
	Foldable
}

func FoldUint64(v uint64) uint32 {
	return uint32(v) ^ uint32(v>>32)
}

/*
This function guarantees that
	(float(bftMajorityCount)/nodeCount > 2.0/3.0)	AND	(float(bftMajorityCount - 1)/nodeCount <= 2.0/3.0)
*/
func BftMajority(nodeCount int) int {
	return nodeCount - BftMinority(nodeCount)
}

func BftMinority(nodeCount int) int {
	return (nodeCount - 1) / 3
}

//TODO ?NeedFix - current implementation can only work for limited cases
func EqualFixedLenWriterTo(t, o io.WriterTo) bool {
	if t == nil || o == nil {
		return false
	}
	return (&writerToComparer{}).compare(t, o)
}

type writerToComparer struct {
	thisValue *[]byte
	other     io.WriterTo
}

func (c *writerToComparer) compare(this io.WriterTo, other io.WriterTo) bool {
	c.thisValue = nil
	c.other = other
	_, _ = this.WriteTo(c)
	return c.other == nil
}

func (c *writerToComparer) Write(otherValue []byte) (int, error) {
	if c.thisValue == nil {
		c.thisValue = &otherValue //result of &var is never nil
		_, err := c.other.WriteTo(c)
		if err != nil {
			return 0, err
		}
	} else {
		if bytes.Equal(*c.thisValue, otherValue) {
			c.other = nil //mark "true"
		}
	}
	return len(otherValue), nil
}

type fixedSize struct {
	data []byte
}

func (c *fixedSize) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, c)
}

func (c *fixedSize) Read(p []byte) (n int, err error) {
	return copy(p, c.data), nil
}

func (c *fixedSize) FoldToUint64() uint64 {
	return FoldToUint64(c.data)
}

func (c *fixedSize) FixedByteSize() int {
	return len(c.data)
}

func CopyFixedSize(v FoldableReader) FoldableReader {
	r := fixedSize{}
	r.data = make([]byte, v.FixedByteSize())
	n, err := v.Read(r.data)
	if err != nil {
		panic(err)
	}
	if n != len(r.data) {
		panic("unexpected")
	}
	return &r
}

type ShortNodeId uint32 // ZERO is RESERVED
const AbsentShortNodeId ShortNodeId = 0
