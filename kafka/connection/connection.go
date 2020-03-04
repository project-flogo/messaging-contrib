package connection

import (
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

func init() {
	connection.RegisterManager("kafkaconnection", &KafkaSharedConn{})
	connection.RegisterManagerFactory(&Factory{})
}

type KafkaSharedConn struct {
	conn KafkaConnection
}
type Factory struct {
}

func (*Factory) Type() string {

	return "kafka:sarama"
}

func (*Factory) NewManager(settings map[string]interface{}) (connection.Manager, error) {
	settingStruct := &Settings{}
	err := metadata.MapToStruct(settings, settingStruct, true)
	if err != nil {
		return nil, err
	}

	conn, err := getKafkaConnection(log.ChildLogger(log.RootLogger(), "kafka-shared-conn"), settingStruct)
	if err != nil {
		//ctx.Logger().Errorf("Kafka parameters initialization got error: [%s]", err.Error())
		return nil, err
	}

	sharedConn := &KafkaSharedConn{conn: conn}

	return sharedConn, nil
}
func (h *KafkaSharedConn) Type() string {

	return "kafka:sarama"
}

func (h *KafkaSharedConn) GetConnection() interface{} {

	return h.conn
}

func (h *KafkaSharedConn) ReleaseConnection(connection interface{}) {

}

func (h *KafkaSharedConn) Stop() error {
	return h.conn.Stop()
}

func (h *KafkaSharedConn) Start() error {
	return nil
}
