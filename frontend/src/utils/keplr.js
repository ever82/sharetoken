/**
 * Keplr Wallet Integration for ShareToken
 * Provides wallet connection, balance query, and transaction signing
 */

import { SigningStargateClient } from "@cosmjs/stargate";

// ShareToken chain configuration for Keplr
const STT_CHAIN_ID = "sharetoken-devnet";
const STT_RPC_ENDPOINT = "http://localhost:26657";
const STT_REST_ENDPOINT = "http://localhost:1317";

// Chain configuration for Keplr
export const getKeplrChainConfig = () => ({
  chainId: STT_CHAIN_ID,
  chainName: "ShareToken",
  rpc: STT_RPC_ENDPOINT,
  rest: STT_REST_ENDPOINT,
  bip44: {
    coinType: 118,
  },
  bech32Config: {
    bech32PrefixAccAddr: "sharetoken",
    bech32PrefixAccPub: "sharetokenpub",
    bech32PrefixValAddr: "sharetokenvaloper",
    bech32PrefixValPub: "sharetokenvaloperpub",
    bech32PrefixConsAddr: "sharetokenvalcons",
    bech32PrefixConsPub: "sharetokenvalconspub",
  },
  currencies: [
    {
      coinDenom: "STT",
      coinMinimalDenom: "stt",
      coinDecimals: 6,
      coinGeckoId: "sharetoken",
    },
    {
      coinDenom: "STAKE",
      coinMinimalDenom: "stake",
      coinDecimals: 6,
    },
  ],
  feeCurrencies: [
    {
      coinDenom: "STT",
      coinMinimalDenom: "stt",
      coinDecimals: 6,
      coinGeckoId: "sharetoken",
      gasPriceStep: {
        low: 0.01,
        average: 0.025,
        high: 0.04,
      },
    },
  ],
  stakeCurrency: {
    coinDenom: "STAKE",
    coinMinimalDenom: "stake",
    coinDecimals: 6,
  },
  features: ["stargate", "ibc-transfer"],
});

/**
 * KeplrWallet class for managing wallet connection
 */
export class KeplrWallet {
  constructor() {
    this.client = null;
    this.account = null;
    this.chainId = STT_CHAIN_ID;
  }

  /**
   * Check if Keplr is installed
   */
  static isInstalled() {
    return typeof window !== "undefined" && window.keplr !== undefined;
  }

  /**
   * Suggest ShareToken chain to Keplr
   */
  async suggestChain() {
    if (!KeplrWallet.isInstalled()) {
      throw new Error("Keplr wallet not found. Please install Keplr extension.");
    }

    try {
      await window.keplr.experimentalSuggestChain(getKeplrChainConfig());
      await window.keplr.enable(this.chainId);
      return true;
    } catch (error) {
      console.error("Failed to suggest chain:", error);
      throw error;
    }
  }

  /**
   * Connect to Keplr wallet
   */
  async connect() {
    if (!KeplrWallet.isInstalled()) {
      throw new Error("Keplr wallet not found. Please install Keplr extension.");
    }

    try {
      // Suggest chain if not already added
      await this.suggestChain();

      // Get signer
      const offlineSigner = window.getOfflineSigner(this.chainId);

      // Create signing client
      this.client = await SigningStargateClient.connectWithSigner(
        STT_RPC_ENDPOINT,
        offlineSigner
      );

      // Get account
      const accounts = await offlineSigner.getAccounts();
      this.account = accounts[0];

      return {
        address: this.account.address,
        pubKey: this.account.pubkey,
      };
    } catch (error) {
      console.error("Failed to connect to Keplr:", error);
      throw error;
    }
  }

  /**
   * Disconnect wallet
   */
  disconnect() {
    this.client = null;
    this.account = null;
  }

  /**
   * Get wallet address
   */
  getAddress() {
    return this.account ? this.account.address : null;
  }

  /**
   * Query STT balance
   */
  async getBalance() {
    if (!this.client || !this.account) {
      throw new Error("Wallet not connected");
    }

    try {
      const balance = await this.client.getBalance(
        this.account.address,
        "stt"
      );
      return balance;
    } catch (error) {
      console.error("Failed to get balance:", error);
      throw error;
    }
  }

  /**
   * Query all balances
   */
  async getAllBalances() {
    if (!this.client || !this.account) {
      throw new Error("Wallet not connected");
    }

    try {
      const balances = await this.client.getAllBalances(this.account.address);
      return balances;
    } catch (error) {
      console.error("Failed to get all balances:", error);
      throw error;
    }
  }

  /**
   * Send STT tokens
   */
  async sendSTT(recipientAddress, amount, memo = "") {
    if (!this.client || !this.account) {
      throw new Error("Wallet not connected");
    }

    try {
      const amountFinal = {
        denom: "stt",
        amount: amount.toString(),
      };

      const fee = {
        amount: [
          {
            denom: "stt",
            amount: "5000",
          },
        ],
        gas: "200000",
      };

      const result = await this.client.sendTokens(
        this.account.address,
        recipientAddress,
        [amountFinal],
        fee,
        memo
      );

      return {
        transactionHash: result.transactionHash,
        height: result.height,
        gasUsed: result.gasUsed,
        gasWanted: result.gasWanted,
      };
    } catch (error) {
      console.error("Failed to send STT:", error);
      throw error;
    }
  }

  /**
   * Get transaction history
   */
  async getTransactionHistory() {
    if (!this.account) {
      throw new Error("Wallet not connected");
    }

    try {
      const response = await fetch(
        `${STT_REST_ENDPOINT}/cosmos/tx/v1beta1/txs?events=message.sender='${this.account.address}'&order_by=ORDER_BY_DESC`
      );
      const data = await response.json();
      return data.tx_responses || [];
    } catch (error) {
      console.error("Failed to get transaction history:", error);
      throw error;
    }
  }

  /**
   * Check if wallet is connected
   */
  isConnected() {
    return this.client !== null && this.account !== null;
  }
}

// Export singleton instance
export const keplrWallet = new KeplrWallet();
export default keplrWallet;
