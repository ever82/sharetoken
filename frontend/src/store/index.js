import { createStore } from 'vuex'

export default createStore({
  state: {
    wallet: {
      address: '',
      balance: 0
    },
    nodeStatus: 'disconnected'
  },
  getters: {
    isWalletConnected: state => !!state.wallet.address
  },
  mutations: {
    SET_WALLET(state, wallet) {
      state.wallet = wallet
    },
    SET_NODE_STATUS(state, status) {
      state.nodeStatus = status
    }
  },
  actions: {
    connectWallet({ commit }, address) {
      commit('SET_WALLET', { address, balance: 0 })
    }
  }
})
