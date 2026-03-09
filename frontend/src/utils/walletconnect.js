/**
 * WalletConnect Integration for ShareToken
 * Provides mobile wallet support via WalletConnect protocol
 */

import WalletConnect from "@walletconnect/client";
import QRCodeModal from "@walletconnect/qrcode-modal";
import { SigningStargateClient } from "@cosmjs/stargate";
import { fromBase64, toBase64 } from "@cosmjs/encoding";

// ShareToken chain configuration
const STT_CHAIN_ID = "sharetoken-devnet";
const STT_RPC_ENDPOINT = "http://localhost:26657";

/**
 * WalletConnectWallet class for mobile wallet connection
 */
export class WalletConnectWallet {
  constructor() {
    this.connector = null;
    this.client = null;
    this.address = null;
    this.pubKey = null;
  }

  /**
   * Initialize WalletConnect session
   */
  async init() {
    // Create connector
    this.connector = new WalletConnect({
      bridge: "https://bridge.walletconnect.org",
      qrcodeModal: QRCodeModal,
    });

    // Check if already connected
    if (!this.connector.connected) {
      // Create new session
      await this.connector.createSession();
    } else {
      // Already connected
      this.onConnect();
    }

    // Subscribe to events
    this.connector.on("connect", (error, payload) => {
      if (error) {
        console.error("WalletConnect connect error:", error);
        return;
      }
      this.onConnect(payload);
    });

    this.connector.on("session_update", (error, payload) => {
      if (error) {
        console.error("WalletConnect session_update error:", error);
        return;
      }
      this.onSessionUpdate(payload);
    });

    this.connector.on("disconnect", (error, payload) => {
      if (error) {
        console.error("WalletConnect disconnect error:", error);
        return;
      }
      this.onDisconnect();
    });

    return {
      uri: this.connector.uri,
    };
  }

  /**
   * Handle connect event
   */
  onConnect(payload) {
    const { accounts, chainId } = payload.params[0];
    this.address = accounts[0];

    console.log("WalletConnect connected:", {
      address: this.address,
      chainId,
    });
  }

  /**
   * Handle session update
   */
  onSessionUpdate(payload) {
    const { accounts, chainId } = payload.params[0];
    this.address = accounts[0];

    console.log("WalletConnect session updated:", {
      address: this.address,
      chainId,
    });
  }

  /**
   * Handle disconnect
   */
  onDisconnect() {
    this.connector = null;
    this.client = null;
    this.address = null;
    this.pubKey = null;

    console.log("WalletConnect disconnected");
  }

  /**
   * Connect to mobile wallet
   */
  async connect() {
    if (!this.connector) {
      await this.init();
    }

    // Wait for connection
    return new Promise((resolve, reject) => {
      const checkConnection = setInterval(() => {
        if (this.address) {
          clearInterval(checkConnection);
          resolve({
            address: this.address,
            pubKey: this.pubKey,
          });
        }
      }, 1000);

      // Timeout after 5 minutes
      setTimeout(() => {
        clearInterval(checkConnection);
        reject(new Error("Connection timeout"));
      }, 300000);
    });
  }

  /**
   * Disconnect wallet
   */
  async disconnect() {
    if (this.connector) {
      await this.connector.killSession();
    }
    this.onDisconnect();
  }

  /**
   * Get wallet address
   */
  getAddress() {
    return this.address;
  }

  /**
   * Query STT balance
   */
  async getBalance() {
    if (!this.address) {
      throw new Error("Wallet not connected");
    }

    try {
      const response = await fetch(
        `${STT_RPC_ENDPOINT}/cosmos/bank/v1beta1/balances/${this.address}/by_denom?denom=stt`
      );
      const data = await response.json();
      return data.balance || { denom: "stt", amount: "0" };
    } catch (error) {
      console.error("Failed to get balance:", error);
      throw error;
    }
  }

  /**
   * Query all balances
   */
  async getAllBalances() {
    if (!this.address) {
      throw new Error("Wallet not connected");
    }

    try {
      const response = await fetch(
        `${STT_RPC_ENDPOINT}/cosmos/bank/v1beta1/balances/${this.address}`
      );
      const data = await response.json();
      return data.balances || [];
    } catch (error) {
      console.error("Failed to get all balances:", error);
      throw error;
    }
  }

  /**
   * Send STT tokens via WalletConnect
   */
  async sendSTT(recipientAddress, amount, memo = "") {
    if (!this.connector || !this.address) {
      throw new Error("Wallet not connected");
    }

    // Build transaction
    const tx = {
      msgs: [
        {
          type: "cosmos-sdk/MsgSend",
          value: {
            from_address: this.address,
            to_address: recipientAddress,
            amount: [
              {
                denom: "stt",
                amount: amount.toString(),
              },
            ],
          },
        },
      ],
      fee: {
        amount: [
          {
            denom: "stt",
            amount: "5000",
          },
        ],
        gas: "200000",
      },
      chain_id: STT_CHAIN_ID,
      memo,
    };

    try {
      // Request signature via WalletConnect
      const result = await this.connector.sendCustomRequest({
        jsonrpc: "2.0",
        method: "cosmos_signAmino",
        params: [this.address, tx],
      });

      return {
        transactionHash: result.hash,
        signature: result.signature,
      };
    } catch (error) {
      console.error("Failed to send STT via WalletConnect:", error);
      throw error;
    }
  }

  /**
   * Get transaction history
   */
  async getTransactionHistory() {
    if (!this.address) {
      throw new Error("Wallet not connected");
    }

    try {
      const response = await fetch(
        `https://rest.sharetoken.network/cosmos/tx/v1beta1/txs?events=message.sender='${this.address}'&order_by=ORDER_BY_DESC`
      );
      const data = await response.json();
      return data.tx_responses || [];
    } catch (error) {
      console.error("Failed to get transaction history:", error);
      return [];
    }
  }

  /**
   * Check if wallet is connected
   */
  isConnected() {
    return this.connector && this.connector.connected && this.address !== null;
  }

  /**
   * Check if WalletConnect is available
   */
  static isAvailable() {
    return typeof window !== "undefined";
  }
}

// Export singleton instance
export const walletConnectWallet = new WalletConnectWallet();
export default walletConnectWallet;
