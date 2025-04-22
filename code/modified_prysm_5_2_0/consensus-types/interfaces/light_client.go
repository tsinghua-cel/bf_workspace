package interfaces

import (
	ssz "github.com/prysmaticlabs/fastssz"
	fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	pb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/proto"
)

type LightClientExecutionBranch = [fieldparams.ExecutionBranchDepth][fieldparams.RootLength]byte
type LightClientSyncCommitteeBranch = [fieldparams.SyncCommitteeBranchDepth][fieldparams.RootLength]byte
type LightClientSyncCommitteeBranchElectra = [fieldparams.SyncCommitteeBranchDepthElectra][fieldparams.RootLength]byte
type LightClientFinalityBranch = [fieldparams.FinalityBranchDepth][fieldparams.RootLength]byte
type LightClientFinalityBranchElectra = [fieldparams.FinalityBranchDepthElectra][fieldparams.RootLength]byte

type LightClientHeader interface {
	ssz.Marshaler
	Proto() proto.Message
	Version() int
	Beacon() *pb.BeaconBlockHeader
	Execution() (ExecutionData, error)
	ExecutionBranch() (LightClientExecutionBranch, error)
}

type LightClientBootstrap interface {
	ssz.Marshaler
	Version() int
	Header() LightClientHeader
	SetHeader(header LightClientHeader) error
	CurrentSyncCommittee() *pb.SyncCommittee
	SetCurrentSyncCommittee(sc *pb.SyncCommittee) error
	CurrentSyncCommitteeBranch() (LightClientSyncCommitteeBranch, error)
	CurrentSyncCommitteeBranchElectra() (LightClientSyncCommitteeBranchElectra, error)
	SetCurrentSyncCommitteeBranch(branch [][]byte) error
}

type LightClientUpdate interface {
	ssz.Marshaler
	Proto() proto.Message
	Version() int
	AttestedHeader() LightClientHeader
	SetAttestedHeader(header LightClientHeader) error
	NextSyncCommittee() *pb.SyncCommittee
	SetNextSyncCommittee(sc *pb.SyncCommittee)
	NextSyncCommitteeBranch() (LightClientSyncCommitteeBranch, error)
	SetNextSyncCommitteeBranch(branch [][]byte) error
	NextSyncCommitteeBranchElectra() (LightClientSyncCommitteeBranchElectra, error)
	FinalizedHeader() LightClientHeader
	SetFinalizedHeader(header LightClientHeader) error
	FinalityBranch() (LightClientFinalityBranch, error)
	FinalityBranchElectra() (LightClientFinalityBranchElectra, error)
	SetFinalityBranch(branch [][]byte) error
	SyncAggregate() *pb.SyncAggregate
	SetSyncAggregate(sa *pb.SyncAggregate)
	SignatureSlot() primitives.Slot
	SetSignatureSlot(slot primitives.Slot)
}

type LightClientFinalityUpdate interface {
	ssz.Marshaler
	ssz.Unmarshaler
	Proto() proto.Message
	Version() int
	AttestedHeader() LightClientHeader
	FinalizedHeader() LightClientHeader
	FinalityBranch() (LightClientFinalityBranch, error)
	FinalityBranchElectra() (LightClientFinalityBranchElectra, error)
	SyncAggregate() *pb.SyncAggregate
	SignatureSlot() primitives.Slot
}

type LightClientOptimisticUpdate interface {
	ssz.Marshaler
	ssz.Unmarshaler
	Proto() proto.Message
	Version() int
	AttestedHeader() LightClientHeader
	SyncAggregate() *pb.SyncAggregate
	SignatureSlot() primitives.Slot
}
