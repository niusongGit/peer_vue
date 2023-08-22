package peer

import (
	pcli "changeme/apiSdk/clents/jsonrpc2"
	"fmt"
	"time"
)

const UpdateDur = 200 * time.Millisecond //更新节点数据间隔时间

func (this *Api) run() {
	doneChannel := make(map[string]chan int)
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-time.After(UpdateDur):
			for _, peer := range this.getActivePeer() {
				if _, ok := doneChannel[peer.ID]; !ok {
					doneChannel[peer.ID] = make(chan int, 1)
					doneChannel[peer.ID] <- 0
				}
				select {
				case rc := <-doneChannel[peer.ID]:
					if rc > 0 {
						rc--
						doneChannel[peer.ID] <- rc
						continue
					}

					go func(peer *Peer) {
						err := update(peer)
						if err != nil {
							rc = 10
							doneChannel[peer.ID] <- rc
							return
						}
						doneChannel[peer.ID] <- rc
					}(peer)
				default:
					continue
				}

			}
		}
	}

}

func update(peer *Peer) error {
	cli := pcli.NewClient(fmt.Sprintf("%s:%d", peer.Ip, peer.Port), peer.Username, peer.Password)
	starAt := time.Now()
	listAccounts, err := cli.ListAccounts()
	if err != nil {
		peer.Error = err.Error()
		peer.setUsedTime(starAt)
		return err
	}

	if len(listAccounts) <= 0 {
		peer.Error = "缺省地址为空"
		peer.setUsedTime(starAt)
		return nil
	}
	addresses := []*AddressInfo{}
	totalBalance := float64(0)
	types := make(map[int]uint, 0)
	for i, v := range listAccounts {
		balance := pcli.AmountDivCV(v.Value)
		addresses = append(addresses, &AddressInfo{
			Index:         i,
			Address:       v.AddrCoin,
			Balance:       balance,
			BalanceFrozen: pcli.AmountDivCV(v.ValueFrozen),
			BalanceLockup: pcli.AmountDivCV(v.ValueLockip),
			PeerType:      v.Type,
		})
		totalBalance += balance
		if v.Type == 1 || v.Type == 2 || v.Type == 3 {
			types[v.Type] = types[v.Type] + 1
		}
	}

	peer.Addresses = addresses
	peer.TotalBalance = totalBalance
	peer.Types = types

	if peer.DefaultAddress == nil {
		peer.DefaultAddress = addresses[0]
	}

	accountInfo, err := cli.GetInfo()
	if err != nil {
		peer.Error = err.Error()
		peer.setUsedTime(starAt)
		return err
	}

	endAt := time.Now().UnixMilli()
	peer.CurrentBlock = accountInfo.CurrentBlock
	peer.HighestBlock = accountInfo.HighestBlock
	peer.SnapshotHeight = accountInfo.SnapshotHeight
	peer.UpdatedAt = endAt
	peer.setUsedTime(starAt)
	return nil
}
