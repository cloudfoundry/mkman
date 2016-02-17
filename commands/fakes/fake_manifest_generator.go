// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/pivotal-cf-experimental/mkman/commands"
)

type FakeManifestGenerator struct {
	GenerateManifestStub        func() error
	generateManifestMutex       sync.RWMutex
	generateManifestArgsForCall []struct{}
	generateManifestReturns struct {
		result1 error
	}
}

func (fake *FakeManifestGenerator) GenerateManifest() error {
	fake.generateManifestMutex.Lock()
	fake.generateManifestArgsForCall = append(fake.generateManifestArgsForCall, struct{}{})
	fake.generateManifestMutex.Unlock()
	if fake.GenerateManifestStub != nil {
		return fake.GenerateManifestStub()
	} else {
		return fake.generateManifestReturns.result1
	}
}

func (fake *FakeManifestGenerator) GenerateManifestCallCount() int {
	fake.generateManifestMutex.RLock()
	defer fake.generateManifestMutex.RUnlock()
	return len(fake.generateManifestArgsForCall)
}

func (fake *FakeManifestGenerator) GenerateManifestReturns(result1 error) {
	fake.GenerateManifestStub = nil
	fake.generateManifestReturns = struct {
		result1 error
	}{result1}
}

var _ commands.ManifestGenerator = new(FakeManifestGenerator)