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

package smachine

import (
	"context"
)

type slotContextMode uint8

const (
	inactiveContext slotContextMode = iota
	discardedContext
	constructContext
	initContext
	execContext
	migrateContext
)

type contextTemplate struct {
	marker struct{}
	mode   slotContextMode
}

func (p *contextTemplate) ensureExactState(expected slotContextMode) {
	if p.mode != expected {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureAtLeastState(s slotContextMode) {
	if p.mode >= s {
		return
	}
	panic("illegal state")
}

func (p *contextTemplate) setState(expected, updated slotContextMode) {
	p.ensureExactState(expected)
	p.mode = updated
}

type slotContext struct {
	contextTemplate
	s *Slot
}

func (p *slotContext) GetSlotID() SlotID {
	return p.s.slotID
}

func (p *slotContext) GetSelf() SlotLink {
	return p.s.NewLink()
}

func (p *slotContext) GetParent() SlotLink {
	return p.s.parent
}

//func (p *slotContext) SetMigrationForStep(fn MigrateFunc) {
//	p.ensureAtLeastState(initContext)
//	p.migrateState = fn
//}

func (p *slotContext) SetMigration(fn MigrateFunc) {
	p.ensureAtLeastState(initContext)
	p.s.migrateSlot = fn
}

func (p *slotContext) Next(fn StateFunc) StateUpdate {
	p.ensureAtLeastState(initContext)
	if fn == nil {
		panic("illegal value")
	}
	return StateUpdate{marker: &p.marker, flags: stateUpdateNext, nextStep: SlotStep{transition: fn}}
}

func (p *slotContext) Stop() StateUpdate {
	p.ensureAtLeastState(initContext)
	return StateUpdate{marker: &p.marker, flags: stateUpdateStop}
}

func (p *slotContext) Replace(fn CreateFunc) StateUpdate {
	p.ensureAtLeastState(migrateContext)
	if fn == nil {
		panic("illegal value")
	}
	return StateUpdate{marker: &p.marker, flags: stateUpdateReplace, param: fn}
}

func (p *slotContext) WaitAny() StateUpdate {
	p.ensureExactState(execContext)
	return StateUpdate{marker: &p.marker, flags: stateUpdateHotWait}
}

func (p *slotContext) WaitAdapter() StateUpdate {
	p.ensureExactState(execContext)
	return StateUpdate{marker: &p.marker, flags: stateUpdateColdWait}
}

func (p *slotContext) Repeat(limit int) StateUpdate {
	p.ensureExactState(execContext)
	return StateUpdate{marker: &p.marker, flags: stateUpdateRepeat, param: limit}
}

func (p *slotContext) Yield() StateUpdate {
	p.ensureExactState(execContext)
	return StateUpdate{marker: &p.marker, flags: stateUpdateRepeat | stateUpdateYield, param: nil}
}

var _ ConstructionContext = &constructionContext{}

type constructionContext struct {
	contextTemplate
	slotID SlotID
	parent SlotLink
}

func (p *constructionContext) GetSlotID() SlotID {
	if p.slotID == 0 {
		panic("illegal state")
	}
	return p.slotID
}

func (p *constructionContext) GetParent() SlotLink {
	return p.parent
}

func (p *constructionContext) executeCreate(nextCreate CreateFunc) StateMachine {
	p.setState(inactiveContext, constructContext)
	defer p.setState(constructContext, discardedContext)

	return nextCreate(p)
}

var _ MigrationContext = &migrationContext{}

type migrationContext struct {
	slotContext
}

func (p *migrationContext) Same() StateUpdate {
	return StateUpdate{marker: &p.marker, flags: stateUpdateNoChange}
}

func (p *migrationContext) executeMigrate(fn MigrateFunc) StateUpdate {
	p.setState(inactiveContext, migrateContext)
	defer p.setState(migrateContext, discardedContext)

	return fn(p).ensureContext(&p.marker)
}

var _ InitializationContext = &initializationContext{}

type initializationContext struct {
	slotContext
}

func (p *initializationContext) executeInitialization(fn InitFunc) StateUpdate {
	p.setState(inactiveContext, initContext)
	defer p.setState(initContext, discardedContext)

	return fn(p).ensureContext(&p.marker)
}

var _ ExecutionContext = &executionContext{}

type executionContext struct {
	slotContext
	worker          *SlotWorker
	countAsyncCalls uint32
}

func (p *executionContext) NewChild(fn CreateFunc) SlotLink {
	p.ensureExactState(execContext)
	if fn == nil {
		panic("illegal value")
	}
	_, link := p.worker.machine.applySlotCreate(nil, p.s.NewLink(), fn)
	return link
}

func (p *executionContext) NextAdapterCall(a ExecutionAdapter, fn AdapterCallFunc, resultState StateFunc) (StateUpdate, context.CancelFunc) {
	p.ensureExactState(execContext)
	if resultState == nil {
		panic("illegal value")
	}
	aq := p.worker.machine.GetAdapterQueue(a)

	cf := &indirectCancel{}

	slotLink := NewSlotLink(p.s)
	return StateUpdate{marker: &p.marker,
		flags:    stateUpdateColdWait | stateUpdateHasAsync,
		nextStep: SlotStep{transition: resultState},

		param: func() {
			cf.set(aq.CallAsync(slotLink, fn, func(fn AsyncResultFunc) {
				p.worker.machine.applyAsyncStateUpdate(slotLink, fn)
			}))
		}}, cf.cancel
}

type indirectCancel struct {
	cancelled bool
	cancelFn  context.CancelFunc
}

func (p *indirectCancel) cancel() {
	p.cancelled = true
	if p.cancelFn != nil {
		p.cancelFn()
	}
}

func (p *indirectCancel) set(cancel context.CancelFunc) {
	if p.cancelFn != nil {
		panic("illegal state")
	}
	if cancel == nil {
		return
	}
	if p.cancelled {
		p.cancel()
	}
}

func (p *executionContext) AdapterSyncCall(a ExecutionAdapter, fn AdapterCallFunc) bool {
	p.ensureExactState(execContext)
	aq := p.worker.machine.GetAdapterQueue(a)

	wc := p.worker.getCond()

	var resultFn AsyncResultFunc
	hasResult := false
	cancelFn := aq.CallAsync(p.s.NewLink(), fn, func(fn AsyncResultFunc) {
		hasResult = true
		resultFn = fn
		wc.L.Lock()
		wc.Broadcast()
		wc.L.Unlock()
	})

	wc.L.Lock()
	wc.Wait()
	wc.L.Unlock()

	if !hasResult && cancelFn != nil {
		// get a signal of something else
		cancelFn()
		return false
	}
	if resultFn == nil {
		return false
	}

	rc := asyncResultContext{slot: p.s}
	rc.executeResult(resultFn)
	return true
}

func (p *executionContext) AdapterAsyncCall(a ExecutionAdapter, fn AdapterCallFunc) context.CancelFunc {
	p.ensureExactState(execContext)
	aq := p.worker.machine.GetAdapterQueue(a)

	slotLink := p.s.NewLink()
	p.countAsyncCalls++

	return aq.CallAsync(slotLink, fn, func(fn AsyncResultFunc) {
		p.worker.machine.applyAsyncStateUpdate(slotLink, fn)
	})
}

func (p *executionContext) executeNextStep() (stopNow bool, stateUpdate StateUpdate, asyncCount uint32) {
	p.setState(inactiveContext, execContext)
	defer p.setState(execContext, discardedContext)

	loopLimit := p.worker.GetLoopLimit()

	for loopCount := 0; loopCount < loopLimit; loopCount++ {
		if p.worker.HasSignal() {
			return true, stateUpdate, p.countAsyncCalls
		}

		current := p.s.nextState
		stateUpdate = current.transition(p)
		stateUpdate.ensureContext(&p.marker)

		switch stateUpdate.getMode() { // fast path(s)
		case stateUpdateRepeat:
			if loopCount < stateUpdate.getInt() {
				continue
			}
		case stateUpdateNext:
			if !stateUpdate.hasAny(stateUpdateHasAsync|stateUpdateYield) &&
				p.s.machine.IsConsecutive(current.transition, stateUpdate.nextStep.transition) {

				p.s.nextState = stateUpdate.nextStep
				continue
			}
		}
		break
	}

	return false, stateUpdate, p.countAsyncCalls
}

var _ AsyncResultContext = &asyncResultContext{}

type asyncResultContext struct {
	slot   *Slot
	wakeup bool
}

func (p *asyncResultContext) GetSlotID() SlotID {
	return p.slot.slotID
}

func (p *asyncResultContext) GetParent() SlotLink {
	return p.slot.parent
}

func (p *asyncResultContext) WakeUp() {
	p.wakeup = true
}

func (p *asyncResultContext) executeResult(fn AsyncResultFunc) bool {
	fn(p)
	return p.wakeup
}
