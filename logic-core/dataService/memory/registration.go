package memory

import (
	"errors"
	"log"
	"sync"

	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

var regist *registRepo

func NewRegistRepo() (*registRepo, map[int]model.Sink) {
	if regist != nil {
		return regist, nil
	}

	regist := &registRepo{
		nodeRepo{
			nmu:   &sync.RWMutex{},
			ninfo: make(map[int]model.Node),
		},
		sinkAddrRepo{
			samu:  &sync.RWMutex{},
			addrs: make(map[int]model.Sink),
		},
		nodeInfoRepo{
			nmu:   &sync.RWMutex{},
			ninfo: make(map[int]model.Nodeinfo),
		},
		pathRepo{
			pmu:   &sync.RWMutex{},
			pinfo: make(map[int]model.Path),
		},
		deliveryRepo{
			dmu:   &sync.RWMutex{},
			dinfo: make(map[int]model.Delivery),
		},
	}

	return regist, regist.addrs
}

type registRepo struct {
	nodeRepo
	sinkAddrRepo
	nodeInfoRepo // not used
	pathRepo
	deliveryRepo
}

var (
	Pid int = 0
)

/**************************************************************/
/* Node Repo                                                  */
/**************************************************************/
type nodeRepo struct {
	nmu   *sync.RWMutex
	ninfo map[int]model.Node
}

func (nr *nodeRepo) FindNode(key int) (*model.Node, error) {
	nr.nmu.RLock()
	defer nr.nmu.RUnlock()

	n, ok := nr.ninfo[key]

	if !ok {
		return nil, errors.New("nodeRepo: cannot find node")
	}
	return &n, nil
}

func (nr *nodeRepo) FindNodesBySinkID(sid int) ([]model.Node, error) {
	nr.nmu.RLock()
	defer nr.nmu.RUnlock()

	if len(nr.ninfo) == 0 {
		return nil, errors.New("nodeRepo: cannot find node")
	}

	res := []model.Node{}
	for _, node := range(nr.ninfo) {
		if node.Sid == sid {
			res = append(res, node)
		}
	}
	return res, nil
}

func (nr *nodeRepo) CreateNode(key int, n *model.Node) error {
	_, ok := nr.ninfo[key]
	if ok {
		return errors.New("nodeRepo: already exist node")
	}
	nr.ninfo[key] = *n
	return nil
}

func (nr *nodeRepo) DeleteNode(key int) error {
	_, ok := nr.ninfo[key]
	if !ok {
		return errors.New("nodeRepo: cannot find node")
	}
	delete(nr.ninfo, key)
	return nil
}

/**************************************************************/
/* NodeInfo Repo                                              */
/**************************************************************/
type nodeInfoRepo struct {
	nmu   *sync.RWMutex
	ninfo map[int]model.Nodeinfo
}

func (nir *nodeInfoRepo) AppendNodeMap(nid int, sid int) error {
	nir.nmu.RLock()
	defer nir.nmu.RUnlock()

	_, ok := nir.ninfo[nid]

	if ok {
		return errors.New("nodeInfoRepo: already exist nid")
	}
	ni := model.Nodeinfo{SinkID: sid}

	nir.ninfo[nid] = ni

	log.Println("test >>>>>> in memory/AppendNodeMap, sinkID : ", ni, "sinkADDR : ")
	return nil

}

func (nir *nodeInfoRepo) GetSid(nid int) (*model.Nodeinfo, error) {
	nir.nmu.RLock()
	defer nir.nmu.RUnlock()

	n, ok := nir.ninfo[nid]

	if !ok {
		return nil, errors.New("nodeRepo: cannot find node")
	}
	return &n, nil
}

/**************************************************************/
/* SinkAddr Repo                                              */
/**************************************************************/
type sinkAddrRepo struct {
	samu  *sync.RWMutex
	addrs map[int]model.Sink
}

func (sar *sinkAddrRepo) AppendSinkAddr(sid int, s *string) error {
	sar.samu.RLock()
	defer sar.samu.RUnlock()
	_, ok := sar.addrs[sid]
	if ok {
		return errors.New("sinkAddrRepo: already exist sink")
	}
	var sink model.Sink
	sink.Addr = *s
	sar.addrs[sid] = sink
	log.Println("test >>>>>> in memory/appendSinkAddr, sinkID : ", sid, "sinkADDR : ", *s)
	return nil
}

/**************************************************************/
/* Path Repo                                                  */
/**************************************************************/
type pathRepo struct {
	pmu   *sync.RWMutex
	pinfo map[int]model.Path
}

func (pr *pathRepo) FindPath(key int) (*model.Path, error) {
	pr.pmu.RLock()
	defer pr.pmu.RUnlock()

	p, ok := pr.pinfo[key]
	if !ok {
		return nil, errors.New("pathRepo: cannot find path")
	}
	return &p, nil
}

func (pr *pathRepo) FindShortestPathStationID(tagid int) (stationid int, err error) {
	pr.pmu.RLock()
	defer pr.pmu.RUnlock()

	if len(pr.pinfo) == 0 {
		return -1, errors.New("pathRepo: no path")
	}

	min := pr.pinfo[1].Distance
	stationid = 1
	for _, path := range(pr.pinfo) {
		if (path.Distance < min) {
			min = path.Distance
			stationid = path.StationID
		}
	}
	return stationid, nil
}

func (pr *pathRepo) CreatePath(p *model.Path) (int, error) {
	Pid += 1
	_, ok := pr.pinfo[Pid]
	if ok {
		return -1, errors.New("pathRepo: already exist path")
	}
	pr.pinfo[Pid] = *p
	return Pid, nil
}

func (pr *pathRepo) DeletePath(key int) error {
	_, ok := pr.pinfo[key]
	if !ok {
		return errors.New("pathRepo: cannot find path")
	}
	delete(pr.pinfo, key)
	return nil
}

/**************************************************************/
/* Delivery Repo                                              */
/**************************************************************/
type deliveryRepo struct {
	dmu   *sync.RWMutex
	dinfo map[int]model.Delivery
}

func (dr *deliveryRepo) FindDelivery(key int) (*model.Delivery, error) {
	dr.dmu.RLock()
	defer dr.dmu.RUnlock()

	d, ok := dr.dinfo[key]
	if !ok {
		return nil, errors.New("deliveryRepo: cannot find delivery")
	}
	return &d, nil
}

func (dr *deliveryRepo) CreateDelivery(key int, d *model.Delivery) error {
	_, ok := dr.dinfo[key]
	if ok {
		return errors.New("deliveryRepo: already exist delivery")
	}
	dr.dinfo[key] = *d
	return nil
}

func (dr *deliveryRepo) DeleteDelivery(key int) error {
	_, ok := dr.dinfo[key]
	if !ok {
		return errors.New("deliveryRepo: cannot find delivery")
	}
	delete(dr.dinfo, key)
	return nil
}
