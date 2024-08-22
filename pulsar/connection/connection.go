package connection

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logger = log.ChildLogger(log.RootLogger(), "pulsar.connection")
var engineLogLevel string

func init() {
	_ = connection.RegisterManagerFactory(&Factory{})
}

type Settings struct {
	URL                  string            `md:"url,required"`
	CaCert               string            `md:"caCert"`
	CertFile             string            `md:"certFile"`
	KeyFile              string            `md:"keyFile"`
	Auth                 string            `md:"auth"`
	AthenzAuthentication map[string]string `md:"athenzAuth"`
	JWT                  string            `md:"jwt"`
	AllowInsecure        bool              `md:"allowInsecure"`
	ConnectionTimeout    int               `md:"connTimeout"`
	OperationTimeout     int               `md:"opTimeout"`
	Audience             string            `md:"audience"`
	PrivateKey           string            `md:"privateKey"`
	Scope                string            `md:"scope"`
	IssuerUrl            string            `md:"issuerUrl"`
}

type PulsarConnection struct {
	client      pulsar.Client
	keystoreDir string
	clientOpts  pulsar.ClientOptions
	connected   bool
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

	var auth pulsar.Authentication
	keystoreDir, err := createTempKeystoreDir(s)
	if err != nil {
		return nil, err
	}
	if s.Auth == "" {

	} else {
		if s.Auth == "TLS" {
			auth, err = getTLSAuthentication(keystoreDir, s)
			if err != nil {
				return nil, err
			}
		} else if s.Auth == "JWT" {
			auth, err = getJWTAuthentication(s)
			if err != nil {
				return nil, err
			}
		} else if s.Auth == "Athenz" {
			auth = getAthenzAuthentication(s)
		} else if s.Auth == "OAuth2" {
			auth = getOAuth2Authentication(s, keystoreDir)
			if auth == nil {
				return nil, fmt.Errorf("Authentication error")
			}
		}
	}

	engineLogLevel = os.Getenv(log.EnvKeyLogLevel)

	customLogger := zapLoggerWrapper{logger: logger}

	connTimeout := s.ConnectionTimeout

	if connTimeout <= 0 {
		connTimeout = 30
	}

	opTimeout := s.OperationTimeout

	if opTimeout <= 0 {
		opTimeout = 30
	}

	clientOpts := pulsar.ClientOptions{
		URL:                        s.URL,
		Authentication:             auth,
		TLSValidateHostname:        false,
		TLSAllowInsecureConnection: s.AllowInsecure,
		Logger:                     &customLogger,
		ConnectionTimeout:          time.Duration(connTimeout) * time.Second,
		OperationTimeout:           time.Duration(opTimeout) * time.Second,
	}

	if strings.Index(s.URL, "pulsar+ssl") >= 0 {
		if keystoreDir == "" {
			clientOpts.TLSTrustCertsFilePath = s.CaCert
		} else if !s.AllowInsecure {
			clientOpts.TLSTrustCertsFilePath = keystoreDir + string(os.PathSeparator) + "cacert.pem"
		}
	}
	logger.Debugf("pulsar.ClientOptions: %v", clientOpts)

	pulsarCnn := &PulsarConnection{keystoreDir: keystoreDir, clientOpts: clientOpts}

	return pulsarCnn, nil

}

func (p *PulsarConnection) Type() string {

	return "pulsar"
}

func (p *PulsarConnection) GetConnection() interface{} {
	return PulsarConnManager{
		Client:     p.client,
		ClientOpts: p.clientOpts,
		Connected:  p.connected,
		Lock:       &sync.RWMutex{}}
}

func (p *PulsarConnection) Stop() error {
	logger.Debug("Stop Pulsar Connection")
	if p.keystoreDir != "" {
		os.RemoveAll(p.keystoreDir)
	}
	return nil
}

func (p *PulsarConnection) Start() error {
	logger.Info("attempting to create client")
	type ClientInfo struct {
		client pulsar.Client
		err    error
	}
	infoChan := make(chan ClientInfo)
	go func() {
		client, err := pulsar.NewClient(p.clientOpts)
		infoChan <- ClientInfo{client: client, err: err}
	}()

	var data ClientInfo
	select {
	case data = <-infoChan:
		p.client = data.client
	case <-time.After(30 * time.Second):
		data.client = nil
		data.err = fmt.Errorf("client creation has timedout after 30 seconds")
	}
	if data.err != nil {
		if strings.Contains(strings.ToLower(data.err.Error()), "authentication error") || strings.Contains(strings.ToLower(data.err.Error()), "empty token credentials") || strings.Contains(strings.ToLower(data.err.Error()), "missing configuration for token auth") || strings.Contains(strings.ToLower(data.err.Error()), "unsupported authentication type") {
			return data.err
		} else {
			logger.Warnf("%v", data.err)
		}
	} else {
		p.connected = true
		logger.Info("new client created")
	}

	return nil
}

// ReleaseConnection clean up connection resources
func (p *PulsarConnection) ReleaseConnection(connection interface{}) {
	logger.Debug("ReleaseConnection")
	if p.keystoreDir != "" {
		os.RemoveAll(p.keystoreDir)
	}
}

func getAthenzAuthentication(s *Settings) pulsar.Authentication {
	if len(s.AthenzAuthentication) != 0 {
		return pulsar.NewAuthenticationAthenz(s.AthenzAuthentication)
	}
	return nil
}

func getOAuth2Authentication(s *Settings, keystoreDir string) pulsar.Authentication {
	return pulsar.NewAuthenticationOAuth2(map[string]string{
		"type":       "client_credentials",
		"privateKey": keystoreDir + string(os.PathSeparator) + "privateKey.json",
		"issuerUrl":  s.IssuerUrl,
		"audience":   s.Audience,
		"scope":      s.Scope,
	})
}

func getTLSAuthentication(keystoreDir string, s *Settings) (auth pulsar.Authentication, err error) {
	if keystoreDir == "" {
		auth = pulsar.NewAuthenticationTLS(s.CertFile, s.KeyFile)
	} else {
		auth = pulsar.NewAuthenticationTLS(keystoreDir+string(os.PathSeparator)+"certfile.pem",
			keystoreDir+string(os.PathSeparator)+"keyfile.pem")
	}
	return
}
func getJWTAuthentication(s *Settings) (auth pulsar.Authentication, err error) {
	auth = pulsar.NewAuthenticationToken(s.JWT)
	return
}

func createTempKeystoreDir(s *Settings) (keystoreDir string, err error) {
	var prikeyObj map[string]interface{}
	var flogoFileValue = true
	logger.Debugf("createTempCertificateDir:  %v", *s)
	if s.CertFile != "" || s.KeyFile != "" || s.CaCert != "" || s.PrivateKey != "" {
		keystoreDir, err = ioutil.TempDir(os.TempDir(), "pulsar")
		if err != nil {
			return
		}
	} else {
		return "", nil
	}
	if s.CaCert != "" {
		flogoFileValue, err = createFile(keystoreDir, &s.CaCert, "cacert.pem")
		if err != nil {
			return
		}
	}
	if s.CertFile != "" {
		flogoFileValue, err = createFile(keystoreDir, &s.CertFile, "certfile.pem")
		if err != nil {
			return
		}
	}
	if s.KeyFile != "" {
		flogoFileValue, err = createFile(keystoreDir, &s.KeyFile, "keyfile.pem")
		if err != nil {
			return
		}
	}
	if s.PrivateKey != "" {
		err = json.Unmarshal([]byte(s.PrivateKey), &prikeyObj)
		if err != nil { //if its not a json string, then its an OSS file spec
			flogoFileValue = false
		} else {
			var keyBytes []byte
			keyBytes, err = getBytesFromFileSetting(prikeyObj)
			if err != nil {
				return
			}
			if keyBytes != nil {
				err = ioutil.WriteFile(keystoreDir+string(os.PathSeparator)+"privateKey.json", keyBytes, 0644)
				if err != nil {
					return
				}
			}
		}
	}
	if !flogoFileValue {
		os.RemoveAll(keystoreDir)
		return "", nil
	}
	return
}

func getBytesFromFileSetting(fileSetting map[string]interface{}) (destArray []byte, err error) {
	var header = "base64,"
	value := fileSetting["content"].(string)
	if value == "" {
		return nil, nil
	}
	if strings.Index(value, header) >= 0 {
		value = value[strings.Index(value, header)+len(header):]
		return decodeValue(value)
	}
	return nil, fmt.Errorf("internal error; file based setting not formatted correctly")
}

type PulsarConnManager struct {
	Client     pulsar.Client
	ClientOpts pulsar.ClientOptions
	Connected  bool
	Lock       *sync.RWMutex
}

func (p *PulsarConnManager) Connect() error {

	logger.Debugf("Acquiring lock for client creation")
	p.Lock.Lock()
	logger.Debugf("lock acquired for creating client")
	defer p.Lock.Unlock()

	if p.Connected {
		return nil
	}

	logger.Info("attempting to create client")
	type ClientInfo struct {
		client pulsar.Client
		err    error
	}
	infoChan := make(chan ClientInfo)

	go func() {
		client, err := pulsar.NewClient(p.ClientOpts)
		infoChan <- ClientInfo{client: client, err: err}
	}()

	select {
	case data := <-infoChan:
		if data.err != nil {
			return data.err
		}
		logger.Info("new client created")
		p.Client = data.client
		p.Connected = true
		return nil
	case <-time.After(30 * time.Second):
		return fmt.Errorf("client creation has timedout after 30 seconds")
	}

}

func (p *PulsarConnManager) GetProducer(producerOptions pulsar.ProducerOptions) (pulsar.Producer, error) {

	err := p.Connect()
	if err != nil {
		return nil, err
	}

	logger.Info("attempting to create producer")
	type ProducerInfo struct {
		producer pulsar.Producer
		err      error
	}
	infoChan := make(chan ProducerInfo)

	go func() {
		producer, err := p.Client.CreateProducer(producerOptions)
		infoChan <- ProducerInfo{producer: producer, err: err}
	}()

	select {
	case data := <-infoChan:
		if data.err != nil {
			return nil, data.err
		}
		logger.Info("producer created")
		return data.producer, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("producer creation has timedout after 30 seconds")
	}
}

func (p *PulsarConnManager) GetSubscriber(consumerOptions pulsar.ConsumerOptions) (pulsar.Consumer, error) {

	if !p.Connected {
		err := p.Connect()
		if err != nil {
			return nil, err
		}
	}

	logger.Info("attempting to create subscriber")
	type ConsumerInfo struct {
		consumer pulsar.Consumer
		err      error
	}
	infoChan := make(chan ConsumerInfo)

	go func() {
		consumer, err := p.Client.Subscribe(consumerOptions)
		infoChan <- ConsumerInfo{consumer: consumer, err: err}
	}()

	select {
	case data := <-infoChan:
		if data.err != nil {
			return nil, data.err
		}
		logger.Info("subscriber created")
		return data.consumer, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("subscriber creation has timedout after 30 seconds")
	}

}

// createFile is a helper function to create a file in the specified keystore directory.
func createFile(keystoreDir string, certVal *string, fileName string) (bool, error) {
	var certObj map[string]interface{}
	if strings.HasPrefix(*certVal, "file://") {
		val := *certVal
		*certVal = val[7:]
		return false, nil
	} else {
		var certBytes []byte
		err := json.Unmarshal([]byte(*certVal), &certObj)
		if err == nil { //if its not a json string, then its an OSS file spec
			certBytes, err = getBytesFromFileSetting(certObj)
			if err != nil {
				return true, err
			}
		} else {
			value := *certVal

			certBytes, err = decodeValue(value)
			if err != nil {
				return true, err
			}

		}
		if certBytes != nil {
			err = ioutil.WriteFile(keystoreDir+string(os.PathSeparator)+fileName, certBytes, 0644)
			if err != nil {
				return true, err
			}
		}
	}
	return true, nil
}

// DecodeValue decodes a base64-encoded string value and returns the corresponding byte array.
func decodeValue(value string) ([]byte, error) {
	decodedLen := base64.StdEncoding.DecodedLen(len(value))
	destArray := make([]byte, decodedLen)
	actualen, err := base64.StdEncoding.Decode(destArray, []byte(value))
	if err != nil {
		return nil, fmt.Errorf("file based setting not base64 encoded: [%s]", err)
	}
	if decodedLen != actualen {
		newDestArray := make([]byte, actualen)
		copy(newDestArray, destArray)
		destArray = newDestArray
		return newDestArray, nil
	}
	return destArray, nil
}
