package node

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hienduyph/genesis/database"
	"github.com/hienduyph/genesis/utils/coders"
	"github.com/hienduyph/goss/errorx"
	"github.com/hienduyph/goss/logger"
)

func NewSyncHandler(
	db *database.State,
) *SyncHandler {
	return &SyncHandler{db: db}
}

type SyncHandler struct {
	db        *database.State
	peerState *PeerState
}

type FromBlockReq struct {
	FromBlock string `json:"fromBlock"`
}

func (r FromBlockReq) AsReqURI(path string) string {
	u := make(url.Values)
	if err := coders.EncodeQuery(r, u); err != nil {
		logger.Error(err, "encode from block req")
	}
	return fmt.Sprintf("%s?%s", path, u.Encode())
}

type SyncRes struct {
	Blocks []database.Block `json:"blocks"`
}

func (s *SyncHandler) FromBlockHandler(r *http.Request) (interface{}, error) {
	d := new(FromBlockReq)
	if e := coders.DecodeQuery(d, r.URL.Query()); e != nil {

		return nil, fmt.Errorf("decode params error: %s, %w", e.Error(), errorx.ErrBadInput)
	}
	hash := database.Hash{}
	if e := hash.UnmarshalText([]byte(d.FromBlock)); e != nil {
		return nil, fmt.Errorf("invalid hash: %s, %w", e.Error(), errorx.ErrBadInput)
	}
	blocks, err := s.db.GetBlockAfter(r.Context(), hash)
	if err != nil {
		return nil, fmt.Errorf("read blocks failed: %w", err)
	}
	return &SyncRes{Blocks: blocks}, nil

}
