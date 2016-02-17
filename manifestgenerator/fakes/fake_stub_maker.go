// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/pivotal-cf-experimental/mkman/manifestgenerator"
)

type FakeStubMaker struct {
	MakeStubStub        func() (string, error)
	makeStubMutex       sync.RWMutex
	makeStubArgsForCall []struct{}
	makeStubReturns struct {
		result1 string
		result2 error
	}
}

func (fake *FakeStubMaker) MakeStub() (string, error) {
	fake.makeStubMutex.Lock()
	fake.makeStubArgsForCall = append(fake.makeStubArgsForCall, struct{}{})
	fake.makeStubMutex.Unlock()
	if fake.MakeStubStub != nil {
		return fake.MakeStubStub()
	} else {
		return fake.makeStubReturns.result1, fake.makeStubReturns.result2
	}
}

func (fake *FakeStubMaker) MakeStubCallCount() int {
	fake.makeStubMutex.RLock()
	defer fake.makeStubMutex.RUnlock()
	return len(fake.makeStubArgsForCall)
}

func (fake *FakeStubMaker) MakeStubReturns(result1 string, result2 error) {
	fake.MakeStubStub = nil
	fake.makeStubReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

var _ manifestgenerator.StubMaker = new(FakeStubMaker)