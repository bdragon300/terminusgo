package objects

import "context"

type Remote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type RemoteIntroducer BaseIntroducer

func (ri *RemoteIntroducer) OnDatabase(path DatabasePath) *RemoteRequester {
	return &RemoteRequester{Client: ri.client, path: path}
}

type RemoteRequester BaseRequester

func (rr *RemoteRequester) ListAll(ctx context.Context, buf *[]Remote) error {
	var httpResponse struct {
		RemoteNames []Remote `json:"api:remote_names"` // FIXME: figure out the real field signature
	}
	sl := rr.Client.C.Get(rr.path.GetPath("remote"))
	if _, err := doRequest(ctx, sl, &httpResponse); err != nil {
		return err
	}

	copy(*buf, httpResponse.RemoteNames)
	return nil
}

// TODO: test on localhost
func (rr *RemoteRequester) Get(ctx context.Context, buf *Remote, name string) error {
	query := struct {
		RemoteName string `url:"remote_name"`
	}{name}
	var httpResponse struct {
		Type       string `json:"@type"`
		RemoteName string `json:"remote_name"`
		RemoteURL  string `json:"remote_url"`
	}
	sl := rr.Client.C.QueryStruct(query).Get(rr.path.GetPath("remote"))
	if _, err := doRequest(ctx, sl, &httpResponse); err != nil {
		return err
	}

	if buf == nil {
		buf = new(Remote)
	}
	*buf = Remote{Name: httpResponse.RemoteName, URL: httpResponse.RemoteURL}

	return nil
}

type RemoteCreateOptions struct {
	RemoteLocation string `json:"remote_location" validate:"required,uri" default:"http://example.com/user/test_db"`
}

func (rr *RemoteRequester) Create(ctx context.Context, name string, options *RemoteCreateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		RemoteCreateOptions
		RemoteLocation string `json:"remote_location"`
	}{*options, name}
	sl := rr.Client.C.BodyJSON(body).Post(rr.path.GetPath("remote")) // FIXME: figure out if such URL is enough
	if _, err = doRequest(ctx, sl, nil); err != nil {
		return err
	}
	return
}

type RemoteUpdateOptions struct {
	RemoteLocation string `json:"remote_location" validate:"required,uri" default:"http://example.com/user/test_db"`
}

func (rr *RemoteRequester) Update(ctx context.Context, name string, options *RemoteUpdateOptions) (err error) {
	if options, err = prepareOptions(options); err != nil {
		return err
	}
	body := struct {
		RemoteUpdateOptions
		RemoteLocation string `json:"remote_location"`
	}{*options, name}
	sl := rr.Client.C.BodyJSON(body).Put(rr.path.GetPath("remote")) // FIXME: figure out if such URL is enough
	if _, err = doRequest(ctx, sl, nil); err != nil {
		return err
	}
	return
}

func (rr *RemoteRequester) Delete(ctx context.Context, name string) error {
	query := struct {
		RemoteName string `url:"remote_name"`
	}{name}
	sl := rr.Client.C.QueryStruct(query).Delete(rr.path.GetPath("remote")) // FIXME: figure out if such URL is enough
	if _, err := doRequest(ctx, sl, nil); err != nil {
		return err
	}

	return nil
}
