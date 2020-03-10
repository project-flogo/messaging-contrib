package connection

import (
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
)

func init() {
	connection.RegisterManager("pulsarConnection", &PulsarConnection{})
	connection.RegisterManagerFactory(&Factory{})
}

type Settings struct {
	URL                  string `md:"url,required"`
	AthenzAuthentication string `md:"athenzauth"`
	CertFile             string `md:"certFile"`
	KeyFile              string `md:"keyFile"`
}

type PulsarConnection struct {
	client pulsar.Client
}

type Factory struct {
}

func (*Factory) Type() string {

	return "pulsar"
}

func (*Factory) NewManager(settings map[string]interface{}) (connection.Manager, error) {
	s := &Settings{}
	err := metadata.MapToStruct(settings, s, true)
	if err != nil {
		return nil, err
	}

	auth := getAuthentication(s)

	clientOps := pulsar.ClientOptions{
		URL:            s.URL,
		Authentication: auth,
	}
	client, err := pulsar.NewClient(clientOps)

	if err != nil {
		return nil, err
	}

	return &PulsarConnection{client: client}, nil
}

func (p *PulsarConnection) Type() string {

	return "pulsar"
}

func (p *PulsarConnection) GetConnection() interface{} {

	return p.client
}

func (p *PulsarConnection) Stop() error {
	return nil
}

func (p *PulsarConnection) Start() error {
	return nil
}

func (p *PulsarConnection) ReleaseConnection(connection interface{}) {

}

func getAuthentication(s *Settings) pulsar.Authentication {
	if s.AthenzAuthentication != "" {
		return pulsar.NewAuthenticationAthenz(s.AthenzAuthentication)

	}
	if s.CertFile != "" && s.KeyFile != "" {
		return pulsar.NewAuthenticationTLS(s.CertFile, s.KeyFile)
	}
	return nil
}
