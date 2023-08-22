package peer

import (
	"changeme/apiSdk/clents"
	pcli "changeme/apiSdk/clents/jsonrpc2"
	"changeme/apiSdk/utils"
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

type Api struct {
	ctx       context.Context
	PeerGroup map[string]string `json:"peer_group" yaml:"peer_group"`
	PeerList  []*Peer           `json:"peer_list" yaml:"peer_list"`
	peerMap   *sync.Map
}

//func (a *Api) startup(ctx context.Context) {
//	a.ctx = ctx
//}

func NewApi() *Api {
	return &Api{
		PeerGroup: make(map[string]string),
		peerMap:   new(sync.Map),
	}
}

func Load(this *Api) error {
	ok, err := utils.PathExists(PeerConfigFile)
	if err != nil {
		return err
	}
	if ok {
		bs, err := os.ReadFile(PeerConfigFile)
		if err != nil {
			return err
		}

		yaml.Unmarshal(bs, &this)
		for k, v := range this.PeerList {
			this.PeerList[k].ID = fmt.Sprintf("%s:%s", v.Group, v.Name)
			this.peerMap.Store(v.ID, v)
		}
	}
	this.ctx = context.Background()
	go this.run()
	return nil
}

func (this *Api) GetPeerList() []*Peer {
	return this.PeerList
}

func (this *Api) GetPeerGroup() map[string]string {
	return this.PeerGroup
}

func (this *Api) SetStatus(id string, status bool) error {
	v, ok := this.peerMap.Load(id)
	if !ok {
		return errors.New("节点不存在！")
	}
	peer, ok := v.(*Peer)
	if !ok {
		return errors.New("节点不存在！")
	}
	peer.Status = status
	return nil
}

func (this *Api) savePeerConfig() error {
	bs, err := yaml.Marshal(this)
	if err != nil {
		return err
	}
	return utils.SaveFile(PeerConfigFile, &bs)
}

func (this *Api) AddGroup(group string) error {
	_, ok := this.PeerGroup[group]
	if ok {
		return errors.New("分组已存在！")
	}
	this.PeerGroup[group] = group
	return this.savePeerConfig()
}

func (this *Api) EditGroup(group string, old string) error {
	_, ok := this.PeerGroup[old]
	if !ok {
		return errors.New("分组不存在！")
	}
	_, ok = this.PeerGroup[group]
	if ok {
		return errors.New("分组已存在！")
	}
	//修改数据
	this.PeerGroup[group] = group
	delete(this.PeerGroup, old)
	this.peerMap.Range(func(key, value any) bool {
		peer, ok := value.(*Peer)
		if ok {
			if peer.Group == old {
				peer.Group = group
				peer.ID = fmt.Sprintf("%s:%s", peer.Group, peer.Name)
			}
		}
		return true
	})
	return this.savePeerConfig()
}

func (this *Api) DelGroup(group string) error {
	_, ok := this.PeerGroup[group]
	if !ok {
		return errors.New("分组不存在！")
	}
	//修改数据
	delete(this.PeerGroup, group)
	ids := make([]string, 0, 0)
	this.peerMap.Range(func(key, value any) bool {
		peer, ok := value.(*Peer)
		if ok {
			if peer.Group == group {
				peer.Status = false
				ids = append(ids, peer.ID)
				for k, p := range this.PeerList {
					if peer.ID == p.ID {
						this.PeerList = append(this.PeerList[:k], this.PeerList[k+1:]...)
						break
					}
				}
			}
		}
		return true
	})
	this.savePeerConfig()
	if len(ids) > 0 {
		time.Sleep(UpdateDur)
		for _, v := range ids {
			this.peerMap.Delete(v)
		}

	}
	return nil
}

func (this *Api) getClentsAndPeer(id string) (c clents.PeerClient, p *Peer, e error) {
	v, ok := this.peerMap.Load(id)
	if !ok {
		return c, p, errors.New("节点不存在！")
	}
	p, ok = v.(*Peer)
	if !ok {
		return c, p, errors.New("节点不存在！！")
	}

	c = pcli.NewClient(fmt.Sprintf("%s:%d", p.Ip, p.Port), p.Username, p.Password)
	return
}

func (this *Api) GetNewAddress(id, pwd string) error {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return err
	}
	err = cli.GetNewAddress(pwd)
	return err
}

func (this *Api) SendToAddress(id, srcAddr string, toAddr string, pwd string, amount float64, gas float64) (*clents.RespSendToAddress, error) {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return nil, err
	}
	return cli.SendToAddress(srcAddr, toAddr, pwd, amount, gas)
}

func (this *Api) Depositin(id string, amount float64, gas float64, rate float64, pwd string, payload string) error {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return err
	}
	return cli.Depositin(amount, gas, rate, pwd, payload)
}

func (this *Api) Depositout(id string, witnessAddr string, amount float64, gas float64, pwd string) error {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return err
	}
	return cli.Depositout(witnessAddr, amount, gas, pwd)
}

func (this *Api) VoteIn(id string, votetype uint16, address string, witness string, amount float64, gas float64, rate float64, pwd string, payload string) error {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return err
	}
	return cli.VoteIn(votetype, address, witness, amount, gas, rate, pwd, payload)
}

func (this *Api) VoteOut(id string, votetype uint16, address string, amount float64, gas float64, rate float64, pwd string, payload string) error {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return err
	}
	return cli.VoteOut(votetype, address, amount, gas, rate, pwd, payload)
}

func (this *Api) CommunityDistribute(id, srcaddress string, gas float64, pwd string) error {
	cli, _, err := this.getClentsAndPeer(id)
	if err != nil {
		return err
	}
	return cli.CommunityDistribute(srcaddress, gas, pwd)
}
