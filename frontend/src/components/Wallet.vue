<template>
  <div class="wallet-container">
    <h2>ShareToken Wallet</h2>

    <!-- First Time / Backup Reminder -->
    <div v-if="showBackupReminder" class="backup-reminder">
      <h3>🔐 Backup Your Wallet</h3>
      <p>Your wallet has been created automatically. Please export and backup your mnemonic phrase to ensure you can recover your wallet.</p>
      <div class="backup-actions">
        <button @click="showExportDialog = true" class="btn btn-warning">Export Mnemonic</button>
        <button @click="dismissBackup" class="btn btn-secondary">Remind Me Later</button>
      </div>
    </div>

    <!-- Export Mnemonic Dialog -->
    <div v-if="showExportDialog" class="modal-overlay">
      <div class="modal">
        <h3>Export Mnemonic</h3>
        <p>Enter your password to reveal your mnemonic phrase:</p>
        <div class="form-group">
          <input v-model="exportPassword" type="password" placeholder="Password" class="input" />
        </div>
        <div v-if="exportedMnemonic" class="mnemonic-box">
          <p class="mnemonic-warning">⚠️ Write down these words in order and keep them safe!</p>
          <code class="mnemonic-words">{{ exportedMnemonic }}</code>
        </div>
        <div class="modal-actions">
          <button @click="doExport" class="btn btn-primary" :disabled="exporting || !exportPassword">
            {{ exporting ? 'Decrypting...' : 'Reveal' }}
          </button>
          <button @click="closeExportDialog" class="btn btn-secondary">Close</button>
        </div>
      </div>
    </div>

    <!-- Wallet Selection / Connection Status -->
    <div class="wallet-status">
      <div v-if="!connected" class="wallet-selection">
        <h3>Connect Wallet</h3>
        <div class="wallet-options">
          <button @click="connectLocalWallet" class="wallet-option local" :disabled="connecting || !hasLocalWallet">
            <span class="wallet-icon">💼</span>
            <span class="wallet-name">Local Wallet</span>
            <span v-if="hasLocalWallet" class="wallet-status-badge">Ready</span>
            <span v-else class="wallet-status-badge unavailable">Not Available</span>
          </button>
          <button @click="connectKeplr" class="wallet-option keplr" :disabled="connecting">
            <span class="wallet-icon">🌐</span>
            <span class="wallet-name">Keplr Browser Extension</span>
          </button>
          <button @click="connectWalletConnect" class="wallet-option walletconnect" :disabled="connecting">
            <span class="wallet-icon">📱</span>
            <span class="wallet-name">Mobile Wallet (WalletConnect)</span>
          </button>
        </div>
      </div>
      <div v-else class="wallet-info">
        <div class="wallet-header">
          <span class="wallet-type-badge" :class="walletType">
            {{ walletTypeLabel }}
          </span>
          <button @click="disconnect" class="btn btn-small btn-danger">Disconnect</button>
        </div>
        <div class="wallet-address">
          <strong>Address:</strong>
          <code>{{ address }}</code>
          <button @click="copyAddress" class="btn btn-small">Copy</button>
        </div>
        <div v-if="walletType === 'local'" class="wallet-actions">
          <button @click="showExportDialog = true" class="btn btn-small btn-warning">Export Mnemonic</button>
        </div>
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
        <label>Amount:</label>
        <div class="amount-input">
          <input v-model="amount" type="number" placeholder="0.00" class="input" step="0.000001" />
          <select v-model="selectedDenom" class="denom-select">
            <option value="stake">STAKE</option>
            <option value="stt">STT</option>
          </select>
        </div>
      </div>
      <div class="form-group">
        <label>Memo (optional):</label>
        <input v-model="memo" type="text" placeholder="Transaction memo" class="input" />
      </div>
      <button @click="sendTokens" class="btn btn-primary" :disabled="sending || !canSend">
        {{ sending ? 'Sending...' : 'Send' }}
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

    <!-- Restore Wallet Section (when no wallet exists) -->
    <div v-if="!hasLocalWallet && !connected" class="restore-section">
      <h3>Restore Existing Wallet</h3>
      <div class="form-group">
        <textarea v-model="restoreMnemonic" placeholder="Enter your 12/24 word mnemonic phrase..." class="input" rows="3"></textarea>
      </div>
      <button @click="restoreWallet" class="btn btn-secondary" :disabled="restoring || !isValidMnemonic">
        {{ restoring ? 'Restoring...' : 'Restore Wallet' }}
      </button>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- Success Message -->
    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
    </div>
  </div>
</template>

<script>
import { keplrWallet } from '../utils/keplr.js';
import { walletConnectWallet } from '../utils/walletconnect.js';

export default {
  name: 'WalletView',
  data() {
    return {
      // Connection state
      connected: false,
      connecting: false,
      walletType: null, // 'local', 'keplr', or 'walletconnect'
      address: '',

      // Local wallet state
      hasLocalWallet: false,
      needsBackup: false,
      showBackupReminder: false,

      // Balance and transactions
      balances: [],
      transactions: [],
      loading: false,
      loadingHistory: false,

      // Transfer form
      recipient: '',
      amount: '',
      selectedDenom: 'stake',
      memo: '',
      sending: false,

      // Export dialog
      showExportDialog: false,
      exportPassword: '',
      exporting: false,
      exportedMnemonic: '',

      // Restore
      restoreMnemonic: '',
      restoring: false,

      // Messages
      error: null,
      successMessage: null,

      // Electron API
      electron: null,
    };
  },
  computed: {
    canSend() {
      return this.recipient && this.amount && parseFloat(this.amount) > 0;
    },
    walletTypeLabel() {
      const labels = {
        local: '💼 Local Wallet',
        keplr: '🌐 Keplr',
        walletconnect: '📱 Mobile Wallet'
      };
      return labels[this.walletType] || 'Unknown';
    },
    isValidMnemonic() {
      const words = this.restoreMnemonic.trim().split(/\s+/);
      return [12, 15, 18, 21, 24].includes(words.length);
    },
  },
  mounted() {
    // Check for Electron API
    if (window.electronAPI) {
      this.electron = window.electronAPI;
      this.checkLocalWallet();

      // Listen for wallet initialization
      if (this.electron.onWalletInitialized) {
        this.electron.onWalletInitialized((data) => {
          console.log('Wallet initialized:', data);
          if (data.created) {
            this.hasLocalWallet = true;
            this.needsBackup = data.needsBackup;
            this.showBackupReminder = data.needsBackup;
          }
        });
      }
    }

    // Check if already connected
    this.checkConnection();
  },
  methods: {
    async checkLocalWallet() {
      if (!this.electron) return;

      try {
        const status = await this.electron.walletGetStatus();
        this.hasLocalWallet = status.exists;
        this.needsBackup = status.needsBackup;

        // Show backup reminder if wallet exists but hasn't been backed up
        if (status.exists && status.needsBackup) {
          this.showBackupReminder = true;
        }
      } catch (err) {
        console.error('Failed to check local wallet:', err);
      }
    },

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

    async connectLocalWallet() {
      if (!this.electron || !this.hasLocalWallet) {
        this.error = 'Local wallet is not available';
        return;
      }

      this.connecting = true;
      this.error = null;

      try {
        const address = await this.electron.walletGetAddress();
        if (address) {
          this.walletType = 'local';
          this.address = address;
          this.connected = true;
          await this.refreshBalance();
        } else {
          throw new Error('Failed to get wallet address');
        }
      } catch (err) {
        this.error = err.message || 'Failed to connect local wallet';
        console.error('Local wallet connection error:', err);
      } finally {
        this.connecting = false;
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
        if (this.walletType === 'local' && this.electron) {
          const result = await this.electron.walletGetBalance(this.address);
          if (result.success) {
            this.balances = result.balances;
          } else {
            throw new Error(result.error);
          }
        } else if (this.walletType === 'keplr') {
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
        const amountInMicro = Math.floor(parseFloat(this.amount) * 1000000);

        let result;
        if (this.walletType === 'local' && this.electron) {
          result = await this.electron.walletSend({
            recipient: this.recipient,
            amount: amountInMicro.toString(),
            denom: this.selectedDenom,
            memo: this.memo
          });
          if (!result.success) {
            throw new Error(result.error);
          }
          result = result.result;
        } else if (this.walletType === 'keplr') {
          result = await keplrWallet.sendSTT(this.recipient, amountInMicro, this.memo);
        } else if (this.walletType === 'walletconnect') {
          result = await walletConnectWallet.sendSTT(this.recipient, amountInMicro, this.memo);
        }

        // Clear form
        this.recipient = '';
        this.amount = '';
        this.memo = '';

        // Refresh balance and history
        await this.refreshBalance();
        await this.refreshHistory();

        this.showSuccess(`Transaction sent! Hash: ${result.transactionHash || result.txHash || 'N/A'}`);
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
        } else {
          // Local wallet - we can query via REST API
          this.transactions = [];
        }
      } catch (err) {
        this.error = err.message || 'Failed to load transaction history';
        console.error('History error:', err);
      } finally {
        this.loadingHistory = false;
      }
    },

    async doExport() {
      if (!this.electron || !this.exportPassword) return;

      this.exporting = true;
      this.error = null;

      try {
        const result = await this.electron.walletExport({
          password: this.exportPassword
        });

        if (result.success) {
          this.exportedMnemonic = result.mnemonic;
          this.needsBackup = false;
        } else {
          throw new Error(result.error);
        }
      } catch (err) {
        this.error = err.message || 'Failed to export wallet';
      } finally {
        this.exporting = false;
      }
    },

    closeExportDialog() {
      this.showExportDialog = false;
      this.exportPassword = '';
      this.exportedMnemonic = '';
    },

    dismissBackup() {
      this.showBackupReminder = false;
    },

    async restoreWallet() {
      if (!this.electron || !this.isValidMnemonic) return;

      this.restoring = true;
      this.error = null;

      try {
        const result = await this.electron.walletRestore({
          mnemonic: this.restoreMnemonic.trim()
        });

        if (result.success) {
          this.hasLocalWallet = true;
          this.needsBackup = false;
          this.restoreMnemonic = '';
          this.showSuccess(`Wallet restored! Address: ${result.address}`);
          // Auto-connect
          await this.connectLocalWallet();
        } else {
          throw new Error(result.error);
        }
      } catch (err) {
        this.error = err.message || 'Failed to restore wallet';
      } finally {
        this.restoring = false;
      }
    },

    copyAddress() {
      navigator.clipboard.writeText(this.address).then(() => {
        this.showSuccess('Address copied to clipboard');
      });
    },

    showSuccess(message) {
      this.successMessage = message;
      setTimeout(() => {
        this.successMessage = null;
      }, 5000);
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

/* Backup Reminder */
.backup-reminder {
  background: #fff3cd;
  border: 2px solid #ffc107;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.backup-reminder h3 {
  margin-top: 0;
  color: #856404;
}

.backup-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 30px;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
}

.modal h3 {
  margin-top: 0;
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.mnemonic-box {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  margin: 15px 0;
}

.mnemonic-warning {
  color: #dc3545;
  font-weight: bold;
  margin-bottom: 10px;
}

.mnemonic-words {
  display: block;
  padding: 10px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  word-break: break-all;
}

/* Wallet Selection */
.wallet-selection h3 {
  margin-top: 0;
}

.wallet-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.wallet-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px;
  border: 2px solid #ddd;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.wallet-option:hover:not(:disabled) {
  border-color: #007bff;
  background: #f8f9fa;
}

.wallet-option:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.wallet-option.local {
  border-color: #28a745;
}

.wallet-option.keplr {
  border-color: #007bff;
}

.wallet-option.walletconnect {
  border-color: #6f42c1;
}

.wallet-icon {
  font-size: 24px;
}

.wallet-name {
  flex: 1;
  font-weight: 500;
}

.wallet-status-badge {
  padding: 2px 8px;
  background: #28a745;
  color: white;
  border-radius: 4px;
  font-size: 12px;
}

.wallet-status-badge.unavailable {
  background: #6c757d;
}

/* Wallet Info */
.wallet-status {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.wallet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.wallet-type-badge {
  padding: 5px 10px;
  border-radius: 4px;
  font-weight: bold;
}

.wallet-type-badge.local {
  background: #d4edda;
  color: #155724;
}

.wallet-type-badge.keplr {
  background: #cce5ff;
  color: #004085;
}

.wallet-type-badge.walletconnect {
  background: #e2e3f3;
  color: #383d41;
}

.wallet-address {
  display: flex;
  align-items: center;
  gap: 10px;
  word-break: break-all;
}

.wallet-address code {
  flex: 1;
  padding: 5px 10px;
  background: #e9ecef;
  border-radius: 4px;
  font-size: 14px;
}

.wallet-actions {
  margin-top: 15px;
  display: flex;
  gap: 10px;
}

/* Sections */
.balance-section,
.transfer-section,
.history-section,
.restore-section {
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
  box-sizing: border-box;
}

.amount-input {
  display: flex;
  gap: 10px;
}

.amount-input .input {
  flex: 1;
}

.denom-select {
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

.btn-warning {
  background: #ffc107;
  color: #212529;
}

.btn-small {
  padding: 5px 10px;
  font-size: 12px;
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

.success-message {
  background: #d4edda;
  color: #155724;
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
