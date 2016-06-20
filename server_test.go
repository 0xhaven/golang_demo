package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	wh, rh := startServer()
	mux := http.NewServeMux()
	mux.Handle("/write", wh)
	mux.Handle("/read", rh)

	srv := httptest.NewServer(mux)
	defer srv.Close()
	if err := doWrite(srv.URL, "foo", "bar"); err != nil {
		t.Fatal(err)
	}
	if value, err := doRead(srv.URL, "foo"); err != nil {
		t.Fatal(err)
	} else if value != "bar" {
		t.Fatalf("%s != %s", value, "bar") // TODO: use testify assert package
	}
}

func doWrite(addr, name, value string) error {
	b, err := json.Marshal(writeReq{Name: name, Value: value})
	if err != nil {
		return err
	}

	_, err = http.Post(addr+"/write", "encoding/json", bytes.NewReader(b))
	return err
}

func doRead(addr, name string) (string, error) {
	b, err := json.Marshal(readReq{Name: name})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(addr+"/read", "encoding/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
