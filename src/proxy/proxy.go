package proxy

import (
	"github.com/andrecronje/babble/src/hashgraph"
)

type AppProxy interface {
	SubmitCh() chan []byte
	CommitBlock(block hashgraph.Block) (CommitResponse, error)
	GetSnapshot(blockIndex int) ([]byte, error)
	Restore(snapshot []byte) error
}
