package rest

import "context"

type Remote struct {
	Name     string `json:"remote_name"`
	Location string `json:"remote_location"`
}

type RemoteIntroducer BaseIntroducer

func (ri *RemoteIntroducer) OnDatabase(path DatabasePath) *RemoteRequester {
	return &RemoteRequester{Client: ri.client, path: path}
}

type RemoteRequester BaseRequester

func (rr *RemoteRequester) WithContext(ctx context.Context) *RemoteRequester {
	r := *rr
	r.ctx = ctx
	return &r
}

func (rr *RemoteRequester) ListAllNames(buf *[]string) (response TerminusResponse, err error) {
	var httpResponse struct {
		RemoteNames []string `json:"api:remote_names"`
	}
	sl := rr.Client.C.Get(rr.path.GetURL("remote"))
	response, err = doRequest(rr.ctx, sl, &httpResponse)
	if err != nil {
		return
	}

	copy(*buf, httpResponse.RemoteNames)
	return
}

func (rr *RemoteRequester) Get(buf *Remote, name string) (response TerminusResponse, err error) {
	query := struct {
		RemoteName string `url:"remote_name"`
	}{name}
	var httpResponse struct {
		RemoteName string `json:"api:remote_name"`
		RemoteURL  string `json:"api:remote_url"`
	}
	sl := rr.Client.C.QueryStruct(query).Get(rr.path.GetURL("remote"))
	response, err = doRequest(rr.ctx, sl, &httpResponse)
	if err != nil {
		return
	}

	*buf = Remote{Name: httpResponse.RemoteName, Location: httpResponse.RemoteURL}
	return
}

func (rr *RemoteRequester) Create(name, uri string) (response TerminusResponse, err error) {
	body := struct {
		RemoteName     string `json:"remote_name"`
		RemoteLocation string `json:"remote_location"`
	}{name, uri}
	sl := rr.Client.C.BodyJSON(body).Post(rr.path.GetURL("remote"))
	return doRequest(rr.ctx, sl, nil)
}

func (rr *RemoteRequester) Update(name, uri string) (response TerminusResponse, err error) {
	body := struct {
		RemoteName     string `json:"remote_name"`
		RemoteLocation string `json:"remote_location"`
	}{name, uri}
	sl := rr.Client.C.BodyJSON(body).Put(rr.path.GetURL("remote"))
	return doRequest(rr.ctx, sl, nil)
}

func (rr *RemoteRequester) Delete(name string) (response TerminusResponse, err error) {
	query := struct {
		RemoteName string `url:"remote_name"`
	}{name}
	sl := rr.Client.C.QueryStruct(query).Delete(rr.path.GetURL("remote"))
	return doRequest(rr.ctx, sl, nil)
}
