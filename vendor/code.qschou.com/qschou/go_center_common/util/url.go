package util

import "net/url"

func AddQueryParams(rawurl string, values url.Values) (ret string, err error) {
	data, err := url.Parse(rawurl)
	if err != nil {
		return
	}
	q := data.Query()
	for k, _ := range values {
		q.Add(k, values.Get(k))
	}
	data.RawQuery = q.Encode()
	ret = data.String()
	return
}
func AddQueryParam(rawurl string, key, val string) (ret string, err error) {
	data := url.Values{}
	data.Add(key, val)
	return AddQueryParams(rawurl, data)
}
