package nsnitro

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	jp "github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(mux *http.ServeMux) (*httptest.Server, *Client) {
	server := httptest.NewServer(mux)
	config := &Config{
		URL: server.URL,
	}

	client := NewClient(config)
	return server, client
}

func Test_client_get_version(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/nitro/v1/config/nsversion",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"errorcode": 0,
				"nsversion": {
					"version": "foo"
				}
			}`)
		})

	server, client := setup(mux)
	defer server.Close()

	version, err := client.Version()
	require.Nil(t, err)
	assert.Equal(t, "foo", version)
}

func Test_client_fetch_objects(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/nitro/v1/config/server",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"server": [{
					"name": "server1",
					"ipaddress": "1.1.1.1"
				}, {
					"name": "server2",
					"ipaddress": "2.2.2.2"
				}]
			}`)
		})
	server, client := setup(mux)
	defer server.Close()

	servers := []Server{}
	err := client.fetch(nsrequest{
		Type: "server",
	}, &servers)

	require.Nil(t, err)
	require.Equal(t, 2, len(servers))
	assert.Equal(t, "server1", servers[0].Name)
	assert.Equal(t, "1.1.1.1", servers[0].IP)
	assert.Equal(t, "server2", servers[1].Name)
	assert.Equal(t, "2.2.2.2", servers[1].IP)
}

func Test_client_post_object(t *testing.T) {
	var objectName string

	mux := http.NewServeMux()
	mux.HandleFunc("/nitro/v1/config/server",
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			objectName, _ = jp.GetString(body, "server", "name")
			w.WriteHeader(201)
		})

	server, client := setup(mux)
	defer server.Close()

	object := NewServer("server1", "1.1.1.1")
	err := client.do("POST",
		nsrequest{Type: "server"},
		&nsresource{Server: &object})
	require.Nil(t, err)
	assert.Equal(t, "server1", objectName)
}

func Test_client_handle_api_error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/nitro/v1/config/server",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{
				"errorcode": 1335,
				"message": "Server already exists [server1]",
				"severity": "ERROR"
			}`))
		})

	server, client := setup(mux)
	defer server.Close()

	object := NewServer("server1", "1.1.1.1")
	err := client.do("POST",
		nsrequest{Type: "server"},
		&nsresource{Server: &object})

	require.NotNil(t, err)
	assert.Equal(t, err.Error(), "Server already exists [server1]")
}
