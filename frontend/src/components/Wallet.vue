<template>
  <div class="wallet-container">
    <h2>ShareToken Wallet</h2>

    <!-- Wallet Connection Status -->
    <div class="wallet-status">
      <div v-if="!connected" class="connect-buttons">
        <button @click="connectKeplr" class="btn btn-primary" :disabled="connecting">
          {{ connecting ? 'Connecting...' : 'Connect Keplr' }}
        </button>
        <button @click="connectWalletConnect" class="btn btn-secondary" :disabled="connecting">
          {{ connecting ? 'Connecting...' : 'Connect Mobile Wallet' }}
        </button>
      </div>
      <div v-else class="wallet-info">
        <p><strong>Address:</strong> {{ address }}</p>
        <button @click="disconnect" class="btn btn-danger">Disconnect</button>
      </div>
    </div>

    <!-- Balance Display -->
    <div v-if="connected" class="balance-section">
      <h3>Your Balances</h3>
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else class="balances">
        <div v-for="balance in balances" :key="balance.denom" class="balance-item">
          <span class="denom">{{ formatDenom(balance.denom) }}</span>
          <span class="amount">{{ formatAmount(balance.amount) }}</span>
        </div>
        <div v-if="balances.length === 0" class="no-balance">
          No balances found
        </div>
      </div>
      <button @click="refreshBalance" class="btn btn-small" :disabled="loading">
        Refresh
      </button>
    </div>

    <!-- Transfer Section -->
    <div v-if="connected" class="transfer-section">
      <h3>Send Tokens</h3>
      <div class="form-group">
        <label>Recipient Address:</label>
        <input v-model="recipient" type="text" placeholder="sharetoken1..." class="input" />
      </div>
      <div class="form-group">
        <label>Amount (STT):</label>
        <input v-model="amount" type="number" placeholder="0.00" class="input" step="0.000001" />
      </div>
      <div class="form-group">
        <label>Memo (optional):</label>
        <input v-model="memo" type="text" placeholder="Transaction memo" class="input" />
      </div>
      <button @click="sendTokens" class="btn btn-primary" :disabled="sending || !canSend">
        {{ sending ? 'Sending...' : 'Send STT' }}
      </button>
    </div>

    <!-- Transaction History -->
    <div v-if="connected" class="history-section">
      <h3>Transaction History</h3>
      <div v-if="loadingHistory" class="loading">Loading...</div>
      <div v-else-if="transactions.length > 0" class="transactions">
        <div v-for="tx in transactions" :key="tx.txhash" class="transaction-item">
          <div class="tx-hash">{{ truncateHash(tx.txhash) }}</div>
          <div class="tx-height">Height: {{ tx.height }}</div>
          <div class="tx-status" :class="tx.code === 0 ? 'success' : 'failed'">
            {{ tx.code === 0 ? 'Success' : 'Failed' }}
          </div>
        </div>
      </div>
      <div v-else class="no-transactions">
        No transactions found
      </div>
      <button @click="refreshHistory" class="btn btn-small" :disabled="loadingHistory">
        Refresh
      </button>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { keplrWallet } from '../utils/keplr.js';
import { walletConnectWallet } from '../utils/walletconnect.js';

export default {
  name: 'Wallet',
  data() {
    return {
      connected: false,
      connecting: false,
      walletType: null, // 'keplr' or 'walletconnect'
      address: '',
      balances: [],
      transactions: [],
      recipient: '',
      amount: '',
      memo: '',
      loading: false,
      loadingHistory: false,
      sending: false,
      error: null,
    };
  },
  computed: {
    canSend() {
      return this.recipient && this.amount && parseFloat(this.amount) > 0;
    },
  },
  mounted() {
    // Check if already connected
    this.checkConnection();
  },
  methods: {
    async checkConnection() {
      if (keplrWallet.isConnected()) {
        this.walletType = 'keplr';
        this.address = keplrWallet.getAddress();
        this.connected = true;
        await this.refreshBalance();
      } else if (walletConnectWallet.isConnected()) {
        this.walletType = 'walletconnect';
        this.address = walletConnectWallet.getAddress();
        this.connected = true;
        await this.refreshBalance();
      }
    },

    async connectKeplr() {
      this.connecting = true;
      this.error = null;

      try {
        const result = await keplrWallet.connect();
        this.walletType = 'keplr';
        this.address = result.address;
        this.connected = true;
        await this.refreshBalance();
      } catch (err) {
        this.error = err.message || 'Failed to connect to Keplr';
        console.error('Keplr connection error:', err);
      } finally {
        this.connecting = false;
      }
    },

    async connectWalletConnect() {
      this.connecting = true;
      this.error = null;

      try {
        const result = await walletConnectWallet.connect();
        this.walletType = 'walletconnect';
        this.address = result.address;
        this.connected = true;
        await this.refreshBalance();
      } catch (err) {
        this.error = err.message || 'Failed to connect mobile wallet';
        console.error('WalletConnect error:', err);
      } finally {
        this.connecting = false;
      }
    },

    async disconnect() {
      if (this.walletType === 'keplr') {
        keplrWallet.disconnect();
      } else if (this.walletType === 'walletconnect') {
        await walletConnectWallet.disconnect();
      }

      this.connected = false;
      this.walletType = null;
      this.address = '';
      this.balances = [];
      this.transactions = [];
    },

    async refreshBalance() {
      if (!this.connected) return;

      this.loading = true;
      this.error = null;

      try {
        if (this.walletType === 'keplr') {
          this.balances = await keplrWallet.getAllBalances();
        } else if (this.walletType === 'walletconnect') {
          this.balances = await walletConnectWallet.getAllBalances();
        }
      } catch (err) {
        this.error = err.message || 'Failed to load balances';
        console.error('Balance error:', err);
      } finally {
        this.loading = false;
      }
    },

    async sendTokens() {
      if (!this.canSend) return;

      this.sending = true;
      this.error = null;

      try {
        const amountInUstt = Math.floor(parseFloat(this.amount) * 1000000);

        let result;
        if (this.walletType === 'keplr') {
          result = await keplrWallet.sendSTT(this.recipient, amountInUstt, this.memo);
        } else if (this.walletType === 'walletconnect') {
          result = await walletConnectWallet.sendSTT(this.recipient, amountInUstt, this.memo);
        }

        // Clear form
        this.recipient = '';
        this.amount = '';
        this.memo = '';

        // Refresh balance and history
        await this.refreshBalance();
        await this.refreshHistory();

        alert(`Transaction sent! Hash: ${result.transactionHash}`);
      } catch (err) {
        this.error = err.message || 'Failed to send tokens';
        console.error('Send error:', err);
      } finally {
        this.sending = false;
      }
    },

    async refreshHistory() {
      if (!this.connected) return;

      this.loadingHistory = true;
      this.error = null;

      try {
        if (this.walletType === 'keplr') {
          this.transactions = await keplrWallet.getTransactionHistory();
        } else if (this.walletType === 'walletconnect') {
          this.transactions = await walletConnectWallet.getTransactionHistory();
        }
      } catch (err) {
        this.error = err.message || 'Failed to load transaction history';
        console.error('History error:', err);
      } finally {
        this.loadingHistory = false;
      }
    },

    formatDenom(denom) {
      if (denom === 'stt') return 'STT';
      if (denom === 'stake') return 'STAKE';
      return denom.toUpperCase();
    },

    formatAmount(amount) {
      const num = parseInt(amount);
      return (num / 1000000).toFixed(6);
    },

    truncateHash(hash) {
      if (!hash) return '';
      return `${hash.slice(0, 8)}...${hash.slice(-8)}`;
    },
  },
};
</script>

<style scoped>
.wallet-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.wallet-status {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.connect-buttons {
  display: flex;
  gap: 10px;
}

.wallet-info {
  word-break: break-all;
}

.balance-section,
.transfer-section,
.history-section {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.balance-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.denom {
  font-weight: bold;
  color: #333;
}

.amount {
  color: #666;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

.input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: opacity 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-small {
  padding: 5px 10px;
  font-size: 12px;
  background: #28a745;
  color: white;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
}

.tx-status.success {
  color: #28a745;
}

.tx-status.failed {
  color: #dc3545;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  padding: 10px;
  border-radius: 4px;
  margin-top: 20px;
}

.loading {
  color: #666;
  font-style: italic;
}

.no-balance,
.no-transactions {
  color: #999;
  text-align: center;
  padding: 20px;
}
</style>
