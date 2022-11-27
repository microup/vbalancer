package proxy

import (
	"fmt"
	"io"
	"net/http"
	"vbalancer/internal/types"
	"vbalancer/internal/vlog"
)

func (p *Proxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI))

	if len(p.peers) == 0 {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), "Peers not found")
		return
	}

	peer, err := p.getNextPeer()
	if err != nil || peer == nil {
		http.Error(w, "Proxy error", http.StatusServiceUnavailable)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), err.Error())
		return
	}

	newProxyURI := fmt.Sprintf("%s://%s%s", peer.Proto, peer.Uri, r.RequestURI)

	newRequest, err := http.NewRequest(r.Method, newProxyURI, r.Body)
	if err != nil {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), err.Error())
		return
	}
	p.copyHeader(r.Header, &newRequest.Header)

	var transport http.Transport
	resp, err := transport.RoundTrip(newRequest)
	if err != nil {
		http.Error(w, "Proxy error", http.StatusInternalServerError)
		if resp != nil {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), vlog.ProxyHost(resp.Request.Host), vlog.ProxyMethod(resp.Request.Method), vlog.ProxyProto(resp.Request.Proto), vlog.ProxyURI(resp.Request.URL.Path), err.Error())
		} else {
			p.logger.Add(vlog.Debug, types.ProxyError, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), err.Error())
		}
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)

	p.logger.Add(vlog.Debug, types.ResultOK, vlog.RemoteAddr(r.RemoteAddr), vlog.ClientHost(r.Host), vlog.ClientMethod(r.Method), vlog.ClientProto(r.Proto), vlog.ClientURI(r.RequestURI), vlog.ProxyHost(resp.Request.Host), vlog.ProxyMethod(resp.Request.Method), vlog.ProxyProto(resp.Request.Proto), vlog.ProxyURI(resp.Request.URL.Path))

}
