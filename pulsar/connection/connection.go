package connection

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logger = log.ChildLogger(log.RootLogger(), "pulsar-connection")
var plogger = log.ChildLogger(log.RootLogger(), "pulsarCustomLogger")
var engineLogLevel string

func init() {
	_ = connection.RegisterManager("pulsarConnection", &PulsarConnection{})
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
	ConnectionTimeout    int8              `md:connTimeout`
	OperationTimeout     int8              `md:opTimeout`
}

type PulsarConnection struct {
	client      pulsar.Client
	keystoreDir string
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
		}
	}

	engineLogLevel = os.Getenv(log.EnvKeyLogLevel)

	customLogger := zapLoggerWrapper{logger: plogger}

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

	client, err := pulsar.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}
	return &PulsarConnection{client: client, keystoreDir: keystoreDir}, nil
}

func (p *PulsarConnection) Type() string {

	return "pulsar"
}

func (p *PulsarConnection) GetConnection() interface{} {
	return p.client
}

func (p *PulsarConnection) Stop() error {
	logger.Debug("Stop Pulsar Connection")
	if p.keystoreDir != "" {
		os.RemoveAll(p.keystoreDir)
	}
	return nil
}

func (p *PulsarConnection) Start() error {
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
	var certObj, keyObj, cacertObj map[string]interface{}
	var flogoFileValue = true
	logger.Debugf("createTempCertificateDir:  %v", *s)
	if s.CertFile != "" || s.KeyFile != "" || s.CaCert != "" {
		keystoreDir, err = ioutil.TempDir(os.TempDir(), "pulsar")
		if err != nil {
			return
		}
	} else {
		return "", nil
	}
	if s.CaCert != "" {
		err = json.Unmarshal([]byte(s.CaCert), &cacertObj)
		if err != nil { //if its not a json string, then its an OSS file spec
			flogoFileValue = false
		} else {
			var certBytes []byte
			certBytes, err = getBytesFromFileSetting(cacertObj)
			if err != nil {
				return
			}
			if certBytes != nil {
				err = ioutil.WriteFile(keystoreDir+string(os.PathSeparator)+"cacert.pem", certBytes, 0644)
				if err != nil {
					return
				}
			}
		}
	}
	if s.CertFile != "" {
		err = json.Unmarshal([]byte(s.CertFile), &certObj)
		if err != nil { //if its not a json string, then its an OSS file spec
			flogoFileValue = false
		} else {
			var certBytes []byte
			certBytes, err = getBytesFromFileSetting(certObj)
			if err != nil {
				return
			}
			if certBytes != nil {
				err = ioutil.WriteFile(keystoreDir+string(os.PathSeparator)+"certfile.pem", certBytes, 0644)
				if err != nil {
					return
				}
			}
		}
	}
	if s.KeyFile != "" {
		err = json.Unmarshal([]byte(s.KeyFile), &keyObj)
		if err != nil { //if its not a json string, then its an OSS file spec
			flogoFileValue = false
		} else {
			var keyBytes []byte
			keyBytes, err = getBytesFromFileSetting(keyObj)
			if err != nil {
				return
			}
			if keyBytes != nil {
				err = ioutil.WriteFile(keystoreDir+string(os.PathSeparator)+"keyfile.pem", keyBytes, 0644)
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
	return nil, fmt.Errorf("internal error; file based setting not formatted correctly")
}
