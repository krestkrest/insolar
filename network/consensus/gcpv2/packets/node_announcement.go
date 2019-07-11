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

package packets

import (
	"fmt"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"

	"github.com/insolar/insolar/network/consensus/common"
)

func NewNodeAnnouncement(np api.NodeProfile, ma api.MembershipAnnouncement, nodeCount int,
	pn pulse_data.PulseNumber) *NodeAnnouncementProfile {
	return &NodeAnnouncementProfile{
		nodeID:    np.GetShortNodeID(),
		nodeCount: uint16(nodeCount),
		ma:        ma,
		pn:        pn,
	}
}

//func NewNodeAnnouncementOf(na MembershipAnnouncementReader, pn common.PulseNumber) *NodeAnnouncementProfile {
//	nr := na.GetNodeRank()
//	return &NodeAnnouncementProfile{
//		nodeID:    na.GetNodeID(),
//		nodeCount: nr.GetTotalCount(),
//		pn:        pn,
//		isLeaving:  na.IsLeaving(),
//		leaveReason: na.GetLeaveReason(),
//		membership: common2.NewMembershipProfile(
//			nr.GetMode(),
//			nr.GetPower(),
//			nr.GetIndex(),
//			na.GetNodeStateHashEvidence(),
//			na.GetAnnouncementSignature(),
//			na.GetRequestedPower(),
//		),
//	}
//}

var _ MembershipAnnouncementReader = &NodeAnnouncementProfile{}

type NodeAnnouncementProfile struct {
	ma        api.MembershipAnnouncement
	nodeID    common.ShortNodeID
	pn        pulse_data.PulseNumber
	nodeCount uint16
}

func (c *NodeAnnouncementProfile) GetRequestedPower() api.MemberPower {
	return c.ma.Membership.RequestedPower
}

func (c *NodeAnnouncementProfile) IsLeaving() bool {
	return c.ma.IsLeaving
}

func (c *NodeAnnouncementProfile) GetLeaveReason() uint32 {
	return c.ma.LeaveReason
}

func (c *NodeAnnouncementProfile) GetJoinerID() common.ShortNodeID {
	if c.ma.Joiner == nil {
		return common.AbsentShortNodeID
	}
	return c.ma.Joiner.GetShortNodeID()
}

func (c *NodeAnnouncementProfile) GetJoinerAnnouncement() JoinerAnnouncementReader {
	panic("unsupported")
}

func (c *NodeAnnouncementProfile) GetNodeRank() MembershipRank {
	return NewMembershipRank(c.ma.Membership.Mode, c.ma.Membership.Power, c.ma.Membership.Index, c.nodeCount)
}

func (c *NodeAnnouncementProfile) GetAnnouncementSignature() api.MemberAnnouncementSignature {
	return c.ma.Membership.AnnounceSignature
}

func (c *NodeAnnouncementProfile) GetNodeID() common.ShortNodeID {
	return c.nodeID
}

func (c *NodeAnnouncementProfile) GetNodeCount() uint16 {
	return c.nodeCount
}

func (c *NodeAnnouncementProfile) GetNodeStateHashEvidence() api.NodeStateHashEvidence {
	return c.ma.Membership.StateEvidence
}

func (c NodeAnnouncementProfile) String() string {
	return fmt.Sprintf("{id:%d %03d/%d %s}", c.nodeID, c.ma.Membership.Index, c.nodeCount, c.ma.Membership.StringParts())
}

func (c *NodeAnnouncementProfile) GetMembershipProfile() api.MembershipProfile {
	return c.ma.Membership
}

func (c *NodeAnnouncementProfile) GetPulseNumber() pulse_data.PulseNumber {
	return c.pn
}
