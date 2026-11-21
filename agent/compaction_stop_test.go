package agent

import (
	"context"
	"iter"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/norma/llm"
)

// truncatedProvider answers every call as if the reply had hit the token limit.
type truncatedProvider struct{}

func (truncatedProvider) Stream(context.Context, llm.CompletionRequest) iter.Seq2[llm.StreamEvent, error] {
	return func(func(llm.StreamEvent, error) bool) {}
}

func (truncatedProvider) Complete(context.Context, llm.CompletionRequest) (llm.Message, string, llm.Usage, error) {
	return llm.Message{}, "max_tokens", llm.Usage{}, nil
}

// A summary cut at the token limit must be refused, so the nodes stay visible to the planner.
func TestCompressRejectsTruncatedReply(t *testing.T) {
	c := &Compactor{prov: truncatedProvider{}}
	_, err := c.compress(context.Background(), &coldGraph{}, block{}, map[int64]*db.Node{})
	if err == nil || !strings.Contains(err.Error(), "token limit") {
		t.Fatalf("compress accepted a reply cut at the token limit: err=%v", err)
	}
}
