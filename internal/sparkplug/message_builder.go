package sparkplug

import (
	"time"

	pb "sparkplug-go/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type MessageBuilder struct {
	seqNum uint64
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		seqNum: 0,
	}
}

func (mb *MessageBuilder) nextSeq() uint64 {
	mb.seqNum++
	return mb.seqNum
}

func (mb *MessageBuilder) CreateNBirth(metrics []*pb.Metric) ([]byte, error) {
	payload := &pb.Payload{
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Metrics:   metrics,
		Seq:       mb.nextSeq(),
		Uuid:      uuid.New().String(),
	}

	return proto.Marshal(payload)
}

func (mb *MessageBuilder) CreateNData(metrics []*pb.Metric) ([]byte, error) {
	payload := &pb.Payload{
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Metrics:   metrics,
		Seq:       mb.nextSeq(),
		Uuid:      uuid.New().String(),
	}

	return proto.Marshal(payload)
}

func (mb *MessageBuilder) CreateDBirth(metrics []*pb.Metric) ([]byte, error) {
	payload := &pb.Payload{
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Metrics:   metrics,
		Seq:       mb.nextSeq(),
		Uuid:      uuid.New().String(),
	}

	return proto.Marshal(payload)
}

func (mb *MessageBuilder) CreateDData(metrics []*pb.Metric) ([]byte, error) {
	payload := &pb.Payload{
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Metrics:   metrics,
		Seq:       mb.nextSeq(),
		Uuid:      uuid.New().String(),
	}

	return proto.Marshal(payload)
}

func (mb *MessageBuilder) CreateNDeath() ([]byte, error) {
	payload := &pb.Payload{
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Metrics:   nil,
		Seq:       mb.nextSeq(),
		Uuid:      uuid.New().String(),
	}

	return proto.Marshal(payload)
}

// CreateMetric creates a new metric with the specified parameters
func CreateMetric(name string, dataType pb.DataType, value interface{}) *pb.Metric {
	metric := &pb.Metric{
		Name:      name,
		Timestamp: uint64(time.Now().UnixNano() / 1000000),
		Datatype:  dataType,
		IsNull:    false,
	}

	switch v := value.(type) {
	case float64:
		metric.Value = &pb.Metric_DoubleValue{DoubleValue: v}
		metric.Datatype = pb.DataType_Double
	case float32:
		metric.Value = &pb.Metric_FloatValue{FloatValue: v}
		metric.Datatype = pb.DataType_Float
	case int64:
		metric.Value = &pb.Metric_LongValue{LongValue: v}
		metric.Datatype = pb.DataType_Int64
	case int32:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_Int32
	case int16:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_Int16
	case int8:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_Int8
	case uint64:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_UInt64
	case uint32:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_UInt32
	case uint16:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_UInt16
	case uint8:
		metric.Value = &pb.Metric_LongValue{LongValue: int64(v)}
		metric.Datatype = pb.DataType_UInt8
	case bool:
		metric.Value = &pb.Metric_BooleanValue{BooleanValue: v}
		metric.Datatype = pb.DataType_Boolean
	case string:
		metric.Value = &pb.Metric_StringValue{StringValue: v}
		metric.Datatype = pb.DataType_String
	case []byte:
		metric.Value = &pb.Metric_BytesValue{BytesValue: v}
		metric.Datatype = pb.DataType_Bytes
	case *pb.DataSet:
		metric.Value = &pb.Metric_DatasetValue{DatasetValue: v}
		metric.Datatype = pb.DataType_DataSetType
	case *pb.Template:
		metric.Value = &pb.Metric_TemplateValue{TemplateValue: v}
		metric.Datatype = pb.DataType_TemplateType
	default:
		metric.IsNull = true
	}

	return metric
}

// CreateDataSet creates a new DataSet with the specified columns and types
func CreateDataSet(columns []string, types []string) *pb.DataSet {
	return &pb.DataSet{
		NumOfColumns: uint64(len(columns)),
		Columns:      columns,
		Types:        types,
		Rows:         make([]*pb.DataSet_Row, 0),
	}
}

// CreatePropertySet creates a new PropertySet with the specified properties
func CreatePropertySet(properties map[string]interface{}) *pb.PropertySet {
	propertySet := &pb.PropertySet{
		Properties: make(map[string]*pb.PropertyValue),
	}

	for key, value := range properties {
		propValue := &pb.PropertyValue{
			IsNull: value == nil,
		}

		if value != nil {
			switch v := value.(type) {
			case float64:
				propValue.Type = pb.DataType_Double
				propValue.Value = &pb.PropertyValue_DoubleValue{DoubleValue: v}
			case float32:
				propValue.Type = pb.DataType_Float
				propValue.Value = &pb.PropertyValue_FloatValue{FloatValue: v}
			case int64:
				propValue.Type = pb.DataType_Int64
				propValue.Value = &pb.PropertyValue_LongValue{LongValue: v}
			case bool:
				propValue.Type = pb.DataType_Boolean
				propValue.Value = &pb.PropertyValue_BooleanValue{BooleanValue: v}
			case string:
				propValue.Type = pb.DataType_String
				propValue.Value = &pb.PropertyValue_StringValue{StringValue: v}
			case *pb.DataSet:
				propValue.Type = pb.DataType_DataSetType
				propValue.Value = &pb.PropertyValue_DatasetValue{DatasetValue: v}
			case *pb.Template:
				propValue.Type = pb.DataType_TemplateType
				propValue.Value = &pb.PropertyValue_TemplateValue{TemplateValue: v}
			}
		}

		propertySet.Properties[key] = propValue
	}

	return propertySet
}

// CreateTemplate creates a new Template with the specified version and reference
func CreateTemplate(version string, templateRef string, isDefinition bool) *pb.Template {
	return &pb.Template{
		Version:      version,
		TemplateRef:  templateRef,
		IsDefinition: isDefinition,
		Parameters:   make([]*pb.Template_Parameter, 0),
	}
}
