package sparkplug

import "fmt"

const (
	SparkplugNamespace = "spBv1.0"
)

type TopicBuilder struct {
	groupID     string
	edgeNodeID  string
	deviceID    string
	scadaHostID string
}

func NewTopicBuilder(groupID, edgeNodeID, deviceID, scadaHostID string) *TopicBuilder {
	return &TopicBuilder{
		groupID:     groupID,
		edgeNodeID:  edgeNodeID,
		deviceID:    deviceID,
		scadaHostID: scadaHostID,
	}
}

func (tb *TopicBuilder) NBirthTopic() string {
	return fmt.Sprintf("%s/%s/NBIRTH/%s", SparkplugNamespace, tb.groupID, tb.edgeNodeID)
}

func (tb *TopicBuilder) NDataTopic() string {
	return fmt.Sprintf("%s/%s/NDATA/%s", SparkplugNamespace, tb.groupID, tb.edgeNodeID)
}

func (tb *TopicBuilder) NDeathTopic() string {
	return fmt.Sprintf("%s/%s/NDEATH/%s", SparkplugNamespace, tb.groupID, tb.edgeNodeID)
}

func (tb *TopicBuilder) DBirthTopic() string {
	return fmt.Sprintf("%s/%s/DBIRTH/%s/%s", SparkplugNamespace, tb.groupID, tb.edgeNodeID, tb.deviceID)
}

func (tb *TopicBuilder) DDataTopic() string {
	return fmt.Sprintf("%s/%s/DDATA/%s/%s", SparkplugNamespace, tb.groupID, tb.edgeNodeID, tb.deviceID)
}

func (tb *TopicBuilder) DDeathTopic() string {
	return fmt.Sprintf("%s/%s/DDEATH/%s/%s", SparkplugNamespace, tb.groupID, tb.edgeNodeID, tb.deviceID)
}

func (tb *TopicBuilder) StateTopic() string {
	return fmt.Sprintf("STATE/%s", tb.scadaHostID)
}
