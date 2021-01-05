package wasmer

import (
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/crypto/ed25519"
	"github.com/ChainSafe/gossamer/lib/genesis"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/scale"
	"github.com/ChainSafe/gossamer/lib/trie"

	log "github.com/ChainSafe/log15"
	"github.com/stretchr/testify/require"
)

func TestInstance_Version_PolkadotRuntime(t *testing.T) {
	expected := &runtime.Version{
		Spec_name:         []byte("polkadot"),
		Impl_name:         []byte("parity-polkadot"),
		Authoring_version: 0,
		Spec_version:      25,
		Impl_version:      0,
	}

	instance := NewTestInstance(t, runtime.POLKADOT_RUNTIME)

	ret, err := instance.exec(runtime.CoreVersion, []byte{})
	require.Nil(t, err)

	version := &runtime.VersionAPI{
		RuntimeVersion: &runtime.Version{},
		API:            nil,
	}
	version.Decode(ret)
	require.Nil(t, err)

	t.Logf("Spec_name: %s\n", version.RuntimeVersion.Spec_name)
	t.Logf("Impl_name: %s\n", version.RuntimeVersion.Impl_name)
	t.Logf("Authoring_version: %d\n", version.RuntimeVersion.Authoring_version)
	t.Logf("Spec_version: %d\n", version.RuntimeVersion.Spec_version)
	t.Logf("Impl_version: %d\n", version.RuntimeVersion.Impl_version)

	require.Equal(t, expected, version.RuntimeVersion)
}

func TestInstance_Version_NodeRuntime(t *testing.T) {
	expected := &runtime.Version{
		Spec_name:         []byte("node"),
		Impl_name:         []byte("substrate-node"),
		Authoring_version: 10,
		Spec_version:      260,
		Impl_version:      0,
	}

	instance := NewTestInstance(t, runtime.NODE_RUNTIME)

	ret, err := instance.exec(runtime.CoreVersion, []byte{})
	require.Nil(t, err)

	version := &runtime.VersionAPI{
		RuntimeVersion: &runtime.Version{},
		API:            nil,
	}
	version.Decode(ret)
	require.Nil(t, err)

	t.Logf("Spec_name: %s\n", version.RuntimeVersion.Spec_name)
	t.Logf("Impl_name: %s\n", version.RuntimeVersion.Impl_name)
	t.Logf("Authoring_version: %d\n", version.RuntimeVersion.Authoring_version)
	t.Logf("Spec_version: %d\n", version.RuntimeVersion.Spec_version)
	t.Logf("Impl_version: %d\n", version.RuntimeVersion.Impl_version)

	require.Equal(t, expected, version.RuntimeVersion)
}

func TestInstance_GrandpaAuthorities_NodeRuntime(t *testing.T) {
	tt := trie.NewEmptyTrie()

	value, err := common.HexToBytes("0x0108eea1eabcac7d2c8a6459b7322cf997874482bfc3d2ec7a80888a3a7d714103640100000000000000b64994460e59b30364cad3c92e3df6052f9b0ebbb8f88460c194dc5794d6d7170100000000000000")
	require.NoError(t, err)

	err = tt.Put(runtime.GrandpaAuthoritiesKey, value)
	require.NoError(t, err)

	rt := NewTestInstanceWithTrie(t, runtime.NODE_RUNTIME, tt, log.LvlTrace)

	auths, err := rt.GrandpaAuthorities()
	require.NoError(t, err)

	authABytes, _ := common.HexToBytes("0xeea1eabcac7d2c8a6459b7322cf997874482bfc3d2ec7a80888a3a7d71410364")
	authBBytes, _ := common.HexToBytes("0xb64994460e59b30364cad3c92e3df6052f9b0ebbb8f88460c194dc5794d6d717")

	authA, _ := ed25519.NewPublicKey(authABytes)
	authB, _ := ed25519.NewPublicKey(authBBytes)

	expected := []*types.Authority{
		{Key: authA, Weight: 1},
		{Key: authB, Weight: 1},
	}

	require.Equal(t, expected, auths)
}

func TestInstance_GrandpaAuthorities_PolkadotRuntime(t *testing.T) {
	tt := trie.NewEmptyTrie()

	value, err := common.HexToBytes("0x0108eea1eabcac7d2c8a6459b7322cf997874482bfc3d2ec7a80888a3a7d714103640100000000000000b64994460e59b30364cad3c92e3df6052f9b0ebbb8f88460c194dc5794d6d7170100000000000000")
	require.NoError(t, err)

	err = tt.Put(runtime.GrandpaAuthoritiesKey, value)
	require.NoError(t, err)

	rt := NewTestInstanceWithTrie(t, runtime.POLKADOT_RUNTIME, tt, log.LvlTrace)

	auths, err := rt.GrandpaAuthorities()
	require.NoError(t, err)

	authABytes, _ := common.HexToBytes("0xeea1eabcac7d2c8a6459b7322cf997874482bfc3d2ec7a80888a3a7d71410364")
	authBBytes, _ := common.HexToBytes("0xb64994460e59b30364cad3c92e3df6052f9b0ebbb8f88460c194dc5794d6d717")

	authA, _ := ed25519.NewPublicKey(authABytes)
	authB, _ := ed25519.NewPublicKey(authBBytes)

	expected := []*types.Authority{
		{Key: authA, Weight: 1},
		{Key: authB, Weight: 1},
	}

	require.Equal(t, expected, auths)
}

func TestInstance_BabeConfiguration_NodeRuntime_NoAuthorities(t *testing.T) {
	rt := NewTestInstance(t, runtime.NODE_RUNTIME)
	cfg, err := rt.BabeConfiguration()
	require.NoError(t, err)

	expected := &types.BabeConfiguration{
		SlotDuration:       3000,
		EpochLength:        200,
		C1:                 1,
		C2:                 4,
		GenesisAuthorities: nil,
		Randomness:         [32]byte{},
		SecondarySlots:     true,
	}

	require.Equal(t, expected, cfg)
}

func TestInstance_BabeConfiguration_NodeRuntime_WithAuthorities(t *testing.T) {
	tt := trie.NewEmptyTrie()

	rvalue, err := common.HexToHash("0x01")
	require.NoError(t, err)
	err = tt.Put(runtime.BABERandomnessKey(), rvalue[:])
	require.NoError(t, err)

	avalue, err := common.HexToBytes("0x08eea1eabcac7d2c8a6459b7322cf997874482bfc3d2ec7a80888a3a7d714103640100000000000000b64994460e59b30364cad3c92e3df6052f9b0ebbb8f88460c194dc5794d6d7170100000000000000")
	require.NoError(t, err)

	err = tt.Put(runtime.BABEAuthoritiesKey(), avalue)
	require.NoError(t, err)

	rt := NewTestInstanceWithTrie(t, runtime.NODE_RUNTIME, tt, log.LvlTrace)

	cfg, err := rt.BabeConfiguration()
	require.NoError(t, err)

	authA, _ := common.HexToHash("0xeea1eabcac7d2c8a6459b7322cf997874482bfc3d2ec7a80888a3a7d71410364")
	authB, _ := common.HexToHash("0xb64994460e59b30364cad3c92e3df6052f9b0ebbb8f88460c194dc5794d6d717")

	expectedAuthData := []*types.AuthorityRaw{
		{Key: authA, Weight: 1},
		{Key: authB, Weight: 1},
	}

	expected := &types.BabeConfiguration{
		SlotDuration:       3000,
		EpochLength:        200,
		C1:                 1,
		C2:                 4,
		GenesisAuthorities: expectedAuthData,
		Randomness:         [32]byte{1},
		SecondarySlots:     true,
	}

	require.Equal(t, expected, cfg)
}

func TestInstance_InitializeBlock_NodeRuntime(t *testing.T) {
	rt := NewTestInstance(t, runtime.NODE_RUNTIME)

	header := &types.Header{
		Number: big.NewInt(1),
		Digest: [][]byte{},
	}

	err := rt.InitializeBlock(header)
	require.NoError(t, err)
}

func TestInstance_InitializeBlock_PolkadotRuntime(t *testing.T) {
	rt := NewTestInstance(t, runtime.POLKADOT_RUNTIME)

	header := &types.Header{
		Number: big.NewInt(1),
		Digest: [][]byte{},
	}

	err := rt.InitializeBlock(header)
	require.NoError(t, err)
}

func buildBlock(t *testing.T, instance runtime.Instance) *types.Block {
	header := &types.Header{
		ParentHash: trie.EmptyHash,
		Number:     big.NewInt(77),
		Digest:     [][]byte{},
	}

	err := instance.InitializeBlock(header)
	require.NoError(t, err)

	idata := types.NewInherentsData()
	err = idata.SetInt64Inherent(types.Timstap0, uint64(time.Now().Unix()))
	require.NoError(t, err)

	err = idata.SetInt64Inherent(types.Babeslot, 1)
	require.NoError(t, err)

	err = idata.SetBigIntInherent(types.Finalnum, big.NewInt(0))
	require.NoError(t, err)

	ienc, err := idata.Encode()
	require.NoError(t, err)

	// Call BlockBuilder_inherent_extrinsics which returns the inherents as extrinsics
	inherentExts, err := instance.InherentExtrinsics(ienc)
	require.NoError(t, err)

	// decode inherent extrinsics
	exts, err := scale.Decode(inherentExts, [][]byte{})
	require.NoError(t, err)

	// apply each inherent extrinsic
	for _, ext := range exts.([][]byte) {
		in, err := scale.Encode(ext) //nolint
		require.NoError(t, err)

		ret, err := instance.ApplyExtrinsic(append([]byte{1}, in...))
		require.NoError(t, err, in)
		require.Equal(t, ret, []byte{0, 0})
	}

	res, err := instance.FinalizeBlock()
	require.NoError(t, err)

	res.Number = header.Number

	expected := &types.Header{
		ParentHash: header.ParentHash,
		Number:     big.NewInt(77),
		Digest:     [][]byte{},
	}

	require.Equal(t, expected.ParentHash, res.ParentHash)
	require.Equal(t, expected.Number, res.Number)
	require.Equal(t, expected.Digest, res.Digest)
	require.NotEqual(t, common.Hash{}, res.StateRoot)
	require.NotEqual(t, common.Hash{}, res.ExtrinsicsRoot)
	require.NotEqual(t, trie.EmptyHash, res.StateRoot)

	return &types.Block{
		Header: res,
		Body:   types.NewBody(inherentExts),
	}
}

func TestInstance_FinalizeBlock_NodeRuntime(t *testing.T) {
	instance := NewTestInstance(t, runtime.NODE_RUNTIME)
	buildBlock(t, instance)
}

func TestInstance_ExecuteBlock_NodeRuntime(t *testing.T) {
	instance := NewTestInstance(t, runtime.NODE_RUNTIME)
	block := buildBlock(t, instance)

	// reset state back to parent state before executing
	parentState := runtime.NewTestRuntimeStorage(t, nil)
	instance.SetContext(parentState)

	_, err := instance.ExecuteBlock(block)
	require.NoError(t, err)
}

func TestInstance_ExecuteBlock_PolkadotRuntime(t *testing.T) {
	instance := NewTestInstance(t, runtime.POLKADOT_RUNTIME)
	block := buildBlock(t, instance)

	// reset state back to parent state before executing
	parentState := runtime.NewTestRuntimeStorage(t, nil)
	instance.SetContext(parentState)

	_, err := instance.ExecuteBlock(block)
	require.NoError(t, err)
}

func TestInstance_ExecuteBlock_PolkadotRuntime_PolkadotBlock1(t *testing.T) {
	t.Skip() // TODO: complete this
	instance := NewTestInstance(t, runtime.POLKADOT_RUNTIME)

	gen, err := genesis.NewGenesisFromJSONRaw("../../../chain/polkadot/genesis.json")
	require.NoError(t, err)

	genTrie, err := genesis.NewTrieFromGenesis(gen)
	require.NoError(t, err)

	// set state to genesis state
	genState := runtime.NewTestRuntimeStorage(t, genTrie)
	instance.SetContext(genState)

	digestBytes := common.MustHexToBytes("0x0c0642414245b501010000000093decc0f00000000362ed8d6055645487fe42e9c8640be651f70a3a2a03658046b2b43f021665704501af9b1ca6e974c257e3d26609b5f68b5b0a1da53f7f252bbe5d94948c39705c98ffa4b869dd44ac29528e3723d619cc7edf1d3f7b7a57a957f6a7e9bdb270a044241424549040118fa3437b10f6e7af8f31362df3a179b991a8c56313d1bcd6307a4d0c734c1ae310100000000000000d2419bc8835493ac89eb09d5985281f5dff4bc6c7a7ea988fd23af05f301580a0100000000000000ccb6bef60defc30724545d57440394ed1c71ea7ee6d880ed0e79871a05b5e40601000000000000005e67b64cf07d4d258a47df63835121423551712844f5b67de68e36bb9a21e12701000000000000006236877b05370265640c133fec07e64d7ca823db1dc56f2d3584b3d7c0f1615801000000000000006c52d02d95c30aa567fda284acf25025ca7470f0b0c516ddf94475a1807c4d250100000000000000000000000000000000000000000000000000000000000000000000000000000005424142450101d468680c844b19194d4dfbdc6697a35bf2b494bda2c5a6961d4d4eacfbf74574379ba0d97b5bb650c2e8670a63791a727943bcb699dc7a228bdb9e0a98c9d089")
	digest, err := scale.Decode(digestBytes, [][]byte{})
	require.NoError(t, err)

	// polkadot block 1, from polkadot.js
	block := &types.Block{
		Header: &types.Header{
			ParentHash:     common.MustHexToHash("0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3"),
			Number:         big.NewInt(1),
			StateRoot:      common.MustHexToHash("0xc56fcd6e7a757926ace3e1ecff9b4010fc78b90d459202a339266a7f6360002f"),
			ExtrinsicsRoot: common.MustHexToHash("0x9a87f6af64ef97aff2d31bebfdd59f8fe2ef6019278b634b2515a38f1c4c2420"),
			Digest:         digest.([][]byte),
		},
		Body: types.NewBody([]byte{8, 40, 4, 3, 0, 11, 80, 149, 160, 81, 114, 1, 16, 4, 20, 0, 0}),
	}

	_, err = instance.ExecuteBlock(block)
	require.NoError(t, err)
}