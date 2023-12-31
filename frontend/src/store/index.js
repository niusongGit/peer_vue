import { createStore } from 'vuex'
export default createStore({
    state: {
        peerInfoIndex:-1,//选中的节点的index
        selectedAddress:"",//选中的地址
        peerList:[],//节点列表数据
        addressesOptions:[],//级联地址选择器数据
        drawer:false, //详情抽屉页是否展示
    },
    getters: {
        getPeerInfo(state) { //详情抽屉页的节点信息
            return state.peerList[state.peerInfoIndex]
        },
        getSelectedAddressInfo(state,getters){//详情抽屉页选择的地址信息
            if (getters.getPeerInfo.addresses && state.selectedAddress) {
                return getters.getPeerInfo.addresses.find(item => item.address === state.selectedAddress)
            }
            return {
                index: -1,
                address: '',
                balance: 0,
                balance_frozen: 0,
                balance_lockup: 0,
                peer_type: 0,
            }
        },
    },
    mutations: {
        setAddressesOptions(state) { //获取级联地址选择器数据
            const options = [];

            // 创建一个 Map 对象，用于存储分组及其对应的节点
            const groupMap = new Map();

            state.peerList.forEach(peer => {
                if (peer.addresses && peer.addresses.length > 0) {
                    const group = peer.group;

                    // 如果当前分组不存在于 groupMap 中，则新建一个分组
                    if (!groupMap.has(group)) {
                        groupMap.set(group, {
                            value: group,
                            label: group,
                            children: []
                        });
                    }

                    const node = {
                        value: peer.name,
                        label: peer.name,
                        children: []
                    };

                    const addresses = peer.addresses.map(address => {
                        return {
                            value: address.address,
                            label: address.address,
                            children: []
                        };
                    });

                    node.children = addresses;
                    groupMap.get(group).children.push(node);
                }
            });

            // 将 groupMap 中的分组添加到 options 数组中
            groupMap.forEach(group => {
                options.push(group);
            });

            state.addressesOptions = options;
        }
    }
})