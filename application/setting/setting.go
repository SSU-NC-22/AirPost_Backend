// AirPost
package setting

import (
	"log"
	"os"
	"strconv"
)

func GetenvInt(target *int, init int, env string) {
	var err error

	temp := os.Getenv(env)
	if temp == "" {
		*target = init
	} else {
		if *target, err = strconv.Atoi(temp); err != nil {
			*target = init
		}
	}
}

/* App setting */
type App struct {
	Server string
}

func (as *App) Getenv() {
	as.Server = os.Getenv("APP_SERVER")
	if as.Server == "" {
		as.Server = "221.140.150.7:8081"
	}
}

var Appsetting = &App{}

/* Database setting */
type Database struct {
	Driver   string `toml:"driver"`
	Server   string `toml:"tcp"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Database string `toml:"database"`
}

func (ds *Database) Getenv() {
	ds.Driver = os.Getenv("DB_DRIVER")
	if ds.Driver == "" {
		ds.Driver = "mysql"
	}
	ds.Server = os.Getenv("DB_SERVER")
	if ds.Server == "" {
		ds.Server = "localhost:3306"
	}
	ds.User = os.Getenv("DB_USER")
	if ds.User == "" {
		ds.User = "airpost"
	}
	ds.Pass = os.Getenv("DB_PASS")
	if ds.Pass == "" {
		ds.Pass = "airpost203"
	}
	ds.Database = os.Getenv("DB_DATABASE")
	if ds.Database == "" {
		ds.Database = "airpost"
	}
}

var Databasesetting = &Database{}

/* Topic setting */
type Topic struct {
	Name         string
	Partitions   int
	Replications int
}

func (ts *Topic) Getenv() {
	ts.Name = os.Getenv("TOPIC_NAME")
	if ts.Name == "" {
		ts.Name = "sensor-data"
	}
	GetenvInt(&ts.Partitions, 1, "TOPIC_PARTITIONS")
	GetenvInt(&ts.Replications, 1, "TOPIC_REPLICATIONS")
}

var Topicsetting = &Topic{}

/* Sink setting */
type DroneSink struct {
	Name		string
	Addr		string
	TopicID		int
}

func (ss *DroneSink) Getenv() {
	ss.Name = os.Getenv("SINK_NAME")
	if ss.Name == "" {
		ss.Name = "drone-sink"
	}
	ss.Addr = os.Getenv("SINK_ADDR")
	if ss.Addr == "" {
		ss.Addr = "111.111.111:8080"
	}
	GetenvInt(&ss.TopicID, 1, "SINK_TOPICID")
}

var DroneSinksetting = &DroneSink{}

type StationSink struct {
	Name		string
	Addr		string
	TopicID		int
}

func (ss *StationSink) Getenv() {
	ss.Name = os.Getenv("SINK_NAME")
	if ss.Name == "" {
		ss.Name = "station-sink"
	}
	ss.Addr = os.Getenv("SINK_ADDR")
	if ss.Addr == "" {
		ss.Addr = "222.222.222:8080"
	}
	GetenvInt(&ss.TopicID, 1, "SINK_TOPICID")
}

var StationSinksetting = &StationSink{}

type TagSink struct {
	Name		string
	Addr		string
	TopicID		int
}

func (ss *TagSink) Getenv() {
	ss.Name = os.Getenv("SINK_NAME")
	if ss.Name == "" {
		ss.Name = "tag-sink"
	}
	ss.Addr = os.Getenv("SINK_ADDR")
	if ss.Addr == "" {
		ss.Addr = "333.333.333:8080"
	}
	GetenvInt(&ss.TopicID, 1, "SINK_TOPICID")
}

var TagSinksetting = &TagSink{}


func init() {
	Appsetting.Getenv()
	Databasesetting.Getenv()
	Topicsetting.Getenv()
	DroneSinksetting.Getenv()
	StationSinksetting.Getenv()
	TagSinksetting.Getenv()

	log.Printf("app : %v\ndb : %v\n", Appsetting, Databasesetting)
}
