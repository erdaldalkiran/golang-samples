package invalidator

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//Invalidator ...
type Invalidator struct {
	client *http.Client
	hosts  []string
}

//New ...
func New() *Invalidator {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          300,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c := &http.Client{Transport: t}

	return &Invalidator{c, []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}}
}

//Invalidate varnish cache
func (i *Invalidator) Invalidate(id string) error {
	for _, host := range i.hosts {
		url := fmt.Sprintf("http://%s/%s", host, "entity/id/"+id)
		r, err := http.NewRequest("PURGE", url, nil)
		if err != nil {
			return fmt.Errorf("invalidator request error: %s", err.Error())
		}
		r.Header.Add("X-Host", "host-name")

		res, err := i.client.Do(r)
		if err != nil {
			return fmt.Errorf("invalidator response url: %s error: %s", url, err.Error())
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			return fmt.Errorf("invalidator unexpected status code url: %s status code: %d body: %s ", url, res.StatusCode, string(body))
		}

		io.Copy(ioutil.Discard, res.Body)
	}

	return nil
}
