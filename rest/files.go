package rest

import (
	"context"
	"io"
	"net/http"

	"github.com/bdragon300/terminusgo/srverror"
	"github.com/bdragon300/tusgo"
)

type FilesIntroducer struct {
	BaseIntroducer
	tusClient *tusgo.Client
	ctx       context.Context
}

func (fi *FilesIntroducer) WithContext(ctx context.Context) *FilesIntroducer {
	r := *fi
	r.ctx = ctx
	return &r
}

func (fi *FilesIntroducer) WithRemoteBaseURI(baseURI string) *FilesIntroducer {
	fi.tusClient.GetRequest = func(method, url string, body io.Reader, _ *tusgo.Client, _ *http.Client, _ *tusgo.ServerCapabilities) (req *http.Request, err error) {
		if req, err = http.NewRequest(method, url, body); err != nil {
			return
		}
		req.Header.Set(srverror.XTerminusDBApiBaseHeader, baseURI)
		return
	}
	return fi
}

func (fi *FilesIntroducer) TusClient() *tusgo.Client {
	res := fi.tusClient
	if fi.ctx != nil {
		res = res.WithContext(fi.ctx)
	}
	return res
}

func NewTusFile(filename string, size int64, metadata map[string]string) *tusgo.File {
	meta := make(map[string]string)
	for k, v := range metadata {
		meta[k] = v
	}
	meta["filename"] = filename
	return &tusgo.File{
		Metadata:   meta,
		RemoteSize: size,
	}
}
