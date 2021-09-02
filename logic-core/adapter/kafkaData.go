package adapter

import (
	"log"
	"strconv"
	"time"

	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

var (
	loc     *time.Location
	timeFmt string
)

func init() {
	loc, _ = time.LoadLocation("Asia/Seoul")
	timeFmt = "2006-01-02 15:04:05"
}

type KafkaData struct {
	// NodeID    int       `json:"node_id"`
	NodeID    string    `json:"node_id"`
	Values    []float64 `json:"values"`
	Timestamp string    `json:"timestamp"`
}

func KafkaToModel(d *KafkaData) (model.KafkaData, error) {
	t, err := time.ParseInLocation(timeFmt, d.Timestamp, loc)
	if err != nil {
		log.Println("Error in KafkaToModel from ParseInLocation")
		return model.KafkaData{}, err
	}
	log.Println("Success KafkaToModel")

	nodeType := d.NodeID[:3]
	log.Println("nodeType = ", nodeType)

	var nodeID int
	nodeID, err = strconv.Atoi(d.NodeID[3:])
	if err != nil {
		log.Println("Error in KafkaToModel from strconv.Atoi")
		return model.KafkaData{}, err
	}
	log.Println("nodeID = ", nodeID)

	return model.KafkaData{
		NodeID:    nodeID,
		Values:    d.Values,
		Timestamp: t,
	}, nil
}
