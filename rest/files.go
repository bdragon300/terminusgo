package rest

import (
	"context"
	"io"
	"net/http"

	"github.com/bdragon300/terminusgo/tusc"

	"github.com/bdragon300/terminusgo/srverror"
	"github.com/bdragon300/tusgo"
)

type FilesIntroducer struct {
	BaseIntroducer
	tusClient *tusc.Client
	ctx       context.Context
}

func (fi *FilesIntroducer) WithContext(ctx context.Context) *FilesIntroducer {
	r := *fi
	r.ctx = ctx
	return &r
}

func (fi *FilesIntroducer) WithRemoteBaseURI(baseURI string) *FilesIntroducer {
	fi.tusClient.GetRequest = func(method, url string, body io.Reader, _ *tusgo.Client, _ *http.Client) (req *http.Request, err error) {
		if req, err = http.NewRequest(method, url, body); err != nil {
			return
		}
		req.Header.Set(srverror.XTerminusDBApiBaseHeader, baseURI)
		return
	}
	return fi
}

func (fi *FilesIntroducer) GetClient() *tusc.Client {
	res := fi.tusClient
	if fi.ctx != nil {
		res = res.WithContext(fi.ctx)
	}
	return res
}
