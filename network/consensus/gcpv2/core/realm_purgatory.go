///
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
///

package core

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/consensuskit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"sync"
)

type NodePhantasm struct {
	mutex   sync.Mutex
	limiter phases.PacketLimiter

	//delayedPackets
	visions map[string]*nodeVision

	//callback *nodeContext
}

type nodeVision struct {
	visionOf          *NodeAppearance
	directAnnouncer   *NodeAppearance
	indirectAnnouncer *NodeAppearance

	profile profiles.ActiveNode // set by construction

	announceSignature proofs.MemberAnnouncementSignature // one-time set
	stateEvidence     proofs.NodeStateHashEvidence       // one-time set
	requestedPower    member.Power                       // one-time set

	firstFraudDetails *misbehavior.FraudError

	neighborReports uint8
}

func NewRealmPurgatory(baselineWeight uint32, local profiles.ActiveNode, nodeCountHint int, phase2ExtLimit uint8,
	fn NodeInitFunc) *RealmPurgatory {

	r := &RealmPurgatory{
		nodeInit:       fn,
		baselineWeight: baselineWeight,
		phase2ExtLimit: phase2ExtLimit,
	}
	r.initPopulation(local, nodeCountHint)

	return r
}

type RealmPurgatory struct {
	nodeInit       NodeInitFunc
	baselineWeight uint32
	phase2ExtLimit uint8
	self           *NodeAppearance

	rw sync.RWMutex

	joinerCount  int
	indexedCount int

	nodeIndex    []*NodeAppearance
	nodeShuffle  []*NodeAppearance // excluding self
	dynamicNodes map[insolar.ShortNodeID]*NodeAppearance
	//	purgatoryByEP map[string]*NodeAppearance
	purgatoryByPK map[string]*NodeAppearance
	purgatoryByID map[insolar.ShortNodeID]*[]*NodeAppearance
	purgatoryOuts map[insolar.ShortNodeID]*NodeAppearance
}

func (r *RealmPurgatory) initPopulation(local profiles.ActiveNode, nodeCountHint int) {
	r.self = r.CreateNodeAppearance(context.Background(), local)
	r.dynamicNodes = make(map[insolar.ShortNodeID]*NodeAppearance, nodeCountHint)
}

func (r *RealmPurgatory) GetSelf() *NodeAppearance {
	return r.self
}

func (r *RealmPurgatory) GetNodeCount() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return len(r.nodeIndex)
}

func (r *RealmPurgatory) GetJoinersCount() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.joinerCount
}

func (r *RealmPurgatory) GetOthersCount() int {
	return r.GetNodeCount() - 1
}

func (r *RealmPurgatory) GetBftMajorityCount() int {
	return consensuskit.BftMajority(r.GetNodeCount())
}

func (r *RealmPurgatory) GetNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.dynamicNodes[id]
}

func (r *RealmPurgatory) GetActiveNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := r.GetNodeAppearance(id)
	if !na.GetProfile().IsJoiner() {
		return na
	}
	return nil
}

func (r *RealmPurgatory) GetJoinerNodeAppearance(id insolar.ShortNodeID) *NodeAppearance {
	na := r.GetNodeAppearance(id)
	if !na.GetProfile().IsJoiner() {
		return nil
	}
	return na
}

func (r *RealmPurgatory) GetNodeAppearanceByIndex(idx int) *NodeAppearance {
	if idx < 0 {
		panic("illegal value")
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	if idx >= len(r.nodeIndex) {
		return nil
	}
	return r.nodeIndex[idx]
}

func (r *RealmPurgatory) GetShuffledOtherNodes() []*NodeAppearance {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.nodeShuffle
}

func (r *RealmPurgatory) IsComplete() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return len(r.nodeIndex) == r.indexedCount
}

func (r *RealmPurgatory) GetIndexedNodes() []*NodeAppearance {
	cp, _ := r.GetIndexedNodesWithCheck()
	// if !ok {
	//	panic("node set is incomplete")
	// }
	return cp
}

func (r *RealmPurgatory) GetIndexedNodesWithCheck() ([]*NodeAppearance, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	cp := make([]*NodeAppearance, len(r.nodeIndex))
	copy(cp, r.nodeIndex)

	return cp, len(r.nodeIndex) == r.indexedCount
}

func (r *RealmPurgatory) CreateNodeAppearance(ctx context.Context, np profiles.ActiveNode) *NodeAppearance {

	n := &NodeAppearance{}
	n.init(np, nil, r.baselineWeight, r.phase2ExtLimit)
	r.nodeInit(ctx, n)

	return n
}

func (r *RealmPurgatory) AddToPurgatory(n *NodeAppearance) (*NodeAppearance, PurgatoryNodeState) {
	nip := n.profile.GetStatic()
	if nip.GetIntroduction() != nil {
		panic("illegal value")
	}

	id := nip.GetStaticNodeID()
	na := r.GetActiveNodeAppearance(id)
	if na != nil {
		return na, PurgatoryExistingMember
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	nn := r.dynamicNodes[id]
	if nn != nil {
		return nn, PurgatoryExistingMember
	}

	if r.purgatoryByPK == nil {
		r.purgatoryByPK = make(map[string]*NodeAppearance)
		r.purgatoryByID = make(map[insolar.ShortNodeID]*[]*NodeAppearance)

		r.purgatoryByPK[nip.GetNodePublicKey().AsByteString()] = n
		r.purgatoryByID[nip.GetStaticNodeID()] = &[]*NodeAppearance{n}
		return n, 0
	}

	pk := nip.GetNodePublicKey().AsByteString()
	nn = r.purgatoryByPK[pk]
	if nn != nil {
		return nn, PurgatoryDuplicatePK
	}

	nodes := r.purgatoryByID[id]

	if nodes == nil {
		nodes = &[]*NodeAppearance{n}
		r.purgatoryByID[id] = nodes
		return n, 0
	}
	*nodes = append(*nodes, n)
	return n, PurgatoryNodeState(len(*nodes) - 1)
}

func (r *RealmPurgatory) AddToDynamics(n *NodeAppearance) (*NodeAppearance, []*NodeAppearance) {
	nip := n.profile.GetStatic()

	if nip.GetIntroduction() == nil {
		panic("illegal value")
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	id := nip.GetStaticNodeID()

	delete(r.purgatoryByPK, nip.GetNodePublicKey().AsByteString())
	nodes := r.purgatoryByID[id]
	if nodes != nil {
		delete(r.purgatoryByID, id)
	} else {
		nodes = &[]*NodeAppearance{}
	}

	na := r.GetActiveNodeAppearance(id)
	if na != nil {
		return na, *nodes
	}

	na = r.dynamicNodes[id]
	if na != nil {
		return na, *nodes
	}

	if n.profile.IsJoiner() {
		r.joinerCount++
	} else {
		ni := n.profile.GetIndex()
		switch {
		case ni.AsInt() == len(r.nodeIndex):
			r.nodeIndex = append(r.nodeIndex, n)
		case ni.AsInt() > len(r.nodeIndex):
			r.nodeIndex = append(r.nodeIndex, make([]*NodeAppearance, 1+ni.AsInt()-len(r.nodeIndex))...)
			r.nodeIndex[ni] = n
		default:
			if r.nodeIndex[ni] != nil {
				panic(fmt.Sprintf("duplicate node id(%v)", ni))
			}
			r.nodeIndex[ni] = n
		}
		r.indexedCount++
		r.nodeShuffle = append(r.nodeShuffle, n)
	}
	return n, *nodes
}

////
//func (r *RealmPurgatory) CreateVectorHelper() *RealmVectorHelper {
//	r.rw.RLock()
//	defer r.rw.RUnlock()
//
//	v := &RealmVectorHelper{realmPopulation: r}
//	v.setArrayNodes(r.nodeIndex, r.dynamicNodes, r.self.callback.GetPopulationVersion())
//	v.realmPopulation = r
//	return v
//}
//
