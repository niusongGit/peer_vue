package peer

import (
	"errors"
	"fmt"
	"time"
)

const PeerConfigFile = "peerconfig.yaml"

type Peer struct {
	ID             string         `json:"id" yaml:"-"`
	Group          string         `json:"group" yaml:"group"`
	Name           string         `json:"name"  yaml:"name"`
	Ip             string         `json:"ip"   yaml:"ip"`
	Port           int            `json:"port" yaml:"port"`
	Username       string         `json:"username" yaml:"username"`
	Password       string         `json:"password" yaml:"password"`
	Status         bool           `json:"status" yaml:"-"`          //是否开启
	HighestBlock   uint64         `json:"highest_block" yaml:"-"`   //所链接的节点的最高高度
	CurrentBlock   uint64         `json:"current_block" yaml:"-"`   // 块高度
	SnapshotHeight uint64         `json:"snapshot_height" yaml:"-"` //快照高度
	TotalBalance   float64        `json:"total_balance" yaml:"-"`   //总余额
	Types          map[int]uint   `json:"types" yaml:"-"`           //角色
	UsedTime       int64          `json:"used_time" yaml:"-"`       //耗时
	IsDel          bool           `json:"is_del" yaml:"-"`
	Error          string         `json:"error" yaml:"-"`
	UpdatedAt      int64          `json:"updated_at" yaml:"-"`
	DefaultAddress *AddressInfo   `json:"default_address" yaml:"-"`
	Addresses      []*AddressInfo `json:"addresses" yaml:"-"`
	usedTimeArr    []int64        //用于计算UsedTime的时间差统计
}

type AddressInfo struct {
	Index         int     `json:"index"`
	Address       string  `json:"address"`
	Balance       float64 `json:"balance"`        //可用余额
	BalanceFrozen float64 `json:"balance_frozen"` //冻结余额
	BalanceLockup float64 `json:"balance_lockup"` //锁仓余额
	PeerType      int     `json:"peer_type"`
}

func (this *Api) AddPeer(p *Peer) (*Peer, error) {
	p.ID = fmt.Sprintf("%s:%s", p.Group, p.Name)
	//是否有重复
	_, ok := this.peerMap.Load(p.ID)
	if ok {
		return nil, errors.New("节点已存在请不要重复添加")
	}
	this.PeerList = append(this.PeerList, p)
	this.peerMap.Store(p.ID, p)
	this.PeerGroup[p.Group] = p.Group
	err := this.savePeerConfig()
	return p, err
}

func (this *Api) EditPeer(p *Peer) (*Peer, error) {
	v, ok := this.peerMap.Load(p.ID)
	if !ok {
		return nil, errors.New("节点不存在！")
	}
	peer, ok := v.(*Peer)
	if !ok {
		return nil, errors.New("节点不存在！！")
	}

	id := fmt.Sprintf("%s:%s", p.Group, p.Name)
	if id != peer.ID {
		_, ok = this.peerMap.Load(id)
		if ok {
			return nil, errors.New("节点不能重复！")
		}
		this.peerMap.Store(id, peer)
		this.peerMap.Delete(p.ID)
	}
	peer.ID = id
	peer.Ip = p.Ip
	peer.Group = p.Group
	peer.Name = p.Name
	peer.Port = p.Port
	peer.Username = p.Username
	peer.Password = p.Password
	this.PeerGroup[p.Group] = p.Group

	err := this.savePeerConfig()
	return peer, err
}

func (this *Api) DelPeer(id string) error {
	v, ok := this.peerMap.Load(id)
	if !ok {
		return errors.New("节点不存在！")
	}
	peer, ok := v.(*Peer)
	if !ok {
		return errors.New("节点不存在！")
	}
	peer.Status = false
	for k, p := range this.PeerList {
		if id == p.ID {
			this.PeerList = append(this.PeerList[:k], this.PeerList[k+1:]...)
			break
		}
	}
	this.savePeerConfig()
	time.Sleep(UpdateDur)
	this.peerMap.Delete(peer.ID)
	return nil
}

func (this *Api) getActivePeer() (list []*Peer) {
	this.peerMap.Range(func(key, value any) bool {
		peer, ok := value.(*Peer)
		if ok {
			if peer.Status {
				list = append(list, peer)
			}
		}
		return true
	})
	return
}

func (p *Peer) setUsedTime(starAt time.Time) {
	dur := time.Since(starAt).Milliseconds()
	if len(p.usedTimeArr) < 10 { //最多统计10次
		p.usedTimeArr = append(p.usedTimeArr, dur)
	} else {
		p.usedTimeArr = p.usedTimeArr[1:]
		p.usedTimeArr = append(p.usedTimeArr, dur)
	}

	sum := int64(0)
	for _, v := range p.usedTimeArr {
		sum += v
	}
	if len(p.usedTimeArr) > 0 {
		p.UsedTime = int64(float64(sum) / float64(len(p.usedTimeArr)))
	} else {
		p.UsedTime = 0
	}
}
