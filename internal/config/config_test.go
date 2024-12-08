package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	envoy "gitlab.com/example/gophers/libs/route-registrator"
)

func TestEnvironment_IsProduction(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		env  string
		exp  bool
	}{
		{
			name: "production",
			env:  productionEnvironment,
			exp:  true,
		},
		{
			name: "empty",
			env:  "",
			exp:  false,
		},
		{
			name: "unknown",
			env:  "unknown",
			exp:  false,
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			e := Environment{Name: s.env}
			assert.Equal(t, s.exp, e.IsProduction())
		})
	}
}

func TestTracing_defineState(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		url  string
		exp  bool
	}{
		{
			name: "enabled",
			url:  "localhost:8029",
			exp:  true,
		},
		{
			name: "disabled",
			url:  "",
			exp:  false,
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			tr := &Tracing{JaegerURL: s.url}
			tr.defineState()
			assert.Equal(t, s.exp, tr.JaegerEnabled)
		})
	}
}

func TestServer_GetHTTPDomain(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		host string
		port int
		exp  string
	}{
		{
			name: "default",
			host: "localhost",
			port: 8080,
			exp:  "localhost:8080",
		},
		{
			name: "empty",
			host: "",
			port: 0,
			exp:  ":0",
		},
		{
			name: "empty host",
			host: "",
			port: 8080,
			exp:  ":8080",
		},
		{
			name: "empty port",
			host: "localhost",
			port: 0,
			exp:  "localhost:0",
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			srv := Server{Host: s.host, HTTPListenAddr: s.port}
			assert.Equal(t, s.exp, srv.GetHTTPDomain())
		})
	}
}

func TestConfig_ToEnvoyConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		serviceMesh ServiceMesh
		exp         *envoy.Config
	}{
		{
			name: "Все данные заполнены",
			serviceMesh: ServiceMesh{
				ServiceMeshHost:                "ServiceMeshHost",
				ServiceMeshPort:                9200,
				IsSendEnabledToServiceMesh:     true,
				SendPeriodToServiceMeshSeconds: 7,
				InternalTmpContainerHost:       "127.0.0.1",
				InternalTmpContainerPort:       8080,
			},
			exp: &envoy.Config{
				ServiceMeshHost:   "ServiceMeshHost",
				ServiceMeshPort:   9200,
				AdvEnabled:        true,
				AdvPeriodSecond:   7,
				AdvHost:           "127.0.0.1",
				AdvPort:           8080,
				ServiceRoutesYAML: serviceRoutesYaml,
			},
		},
		{
			name:        "Пустые значения",
			serviceMesh: ServiceMesh{},
			exp: &envoy.Config{
				ServiceRoutesYAML: serviceRoutesYaml,
			},
		},
	}

	for _, s := range cases {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			c := Config{
				Server:   Server{},
				Log:      Log{},
				Database: Database{},
				Kafka:    Kafka{},
				Cache:    Cache{},
				Tracing:  Tracing{},
				ServiceMesh: ServiceMesh{
					ServiceMeshHost:                s.serviceMesh.ServiceMeshHost,
					ServiceMeshPort:                s.serviceMesh.ServiceMeshPort,
					IsSendEnabledToServiceMesh:     s.serviceMesh.IsSendEnabledToServiceMesh,
					SendPeriodToServiceMeshSeconds: s.serviceMesh.SendPeriodToServiceMeshSeconds,
					InternalTmpContainerHost:       s.serviceMesh.InternalTmpContainerHost,
					InternalTmpContainerPort:       s.serviceMesh.InternalTmpContainerPort,
				},
				Environment: Environment{},
				ServiceName: "",
				Version:     "",
			}
			assert.Equal(t, s.exp, c.ToEnvoyConfig())
		})
	}
}
