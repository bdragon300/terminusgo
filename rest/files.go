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
}

func (fr *FilesIntroducer) WithRemoteBaseURI(baseURI string) *FilesIntroducer {
	fr.tusClient.GetRequest = func(method, url string, body io.Reader, _ *tusgo.Client, _ *http.Client, _ *tusgo.ServerCapabilities) (req *http.Request, err error) {
		if req, err = http.NewRequest(method, url, body); err != nil {
			return
		}
		req.Header.Set(srverror.XTerminusDBApiBaseHeader, baseURI)
		return
	}
	return fr
}

func (fr *FilesIntroducer) TusClient(ctx context.Context) *tusgo.Client {
	return fr.tusClient.WithContext(ctx)
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
