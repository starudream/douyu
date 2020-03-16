package douyu

import (
	"io/ioutil"
	"net/http"

	"github.com/go-sdk/utilx/json"
)

func httpGet(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return err
	}
	return nil
}
