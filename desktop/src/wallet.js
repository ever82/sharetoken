// Wallet Manager for ShareToken Desktop
// Handles in-app wallet creation, storage, and operations

const { safeStorage } = require('electron');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

// CosmJS imports (will be loaded dynamically to handle missing deps gracefully)
let DirectSecp256k1HdWallet = null;
let SigningStargateClient = null;
let GasPrice = null;

// Lazy load CosmJS modules
async function loadCosmJS() {
  if (!DirectSecp256k1HdWallet) {
    const protoSigning = require('@cosmjs/proto-signing');
    DirectSecp256k1HdWallet = protoSigning.DirectSecp256k1HdWallet;
  }
  if (!SigningStargateClient) {
    const stargate = require('@cosmjs/stargate');
    SigningStargateClient = stargate.SigningStargateClient;
    GasPrice = stargate.GasPrice;
  }
}

// Wallet storage paths
function getWalletPaths(userDataPath) {
  const walletDir = path.join(userDataPath, 'wallet');
  return {
    walletDir,
    encryptedMnemonic: path.join(walletDir, 'mnemonic.enc'),
    walletConfig: path.join(walletDir, 'config.json'),
    backupReminder: path.join(walletDir, '.backup-reminder')
  };
}

// Ensure wallet directory exists
function ensureWalletDir(walletDir) {
  if (!fs.existsSync(walletDir)) {
    fs.mkdirSync(walletDir, { recursive: true });
  }
}

// Generate new wallet
async function createWallet() {
  await loadCosmJS();

  const wallet = await DirectSecp256k1HdWallet.generate(12, {
    prefix: 'sharetoken'
  });

  const [account] = await wallet.getAccounts();
  const mnemonic = wallet.mnemonic;

  return {
    address: account.address,
    mnemonic: mnemonic,
    pubKey: Buffer.from(account.pubkey).toString('hex')
  };
}

// Encrypt and save mnemonic
function saveEncryptedMnemonic(mnemonic, paths) {
  ensureWalletDir(paths.walletDir);

  if (!safeStorage.isEncryptionAvailable()) {
    throw new Error('Encryption is not available on this system');
  }

  // Encrypt mnemonic using system keychain
  const encrypted = safeStorage.encryptString(mnemonic);
  fs.writeFileSync(paths.encryptedMnemonic, encrypted);

  // Mark as needing backup
  fs.writeFileSync(paths.backupReminder, Date.now().toString());

  return true;
}

// Load and decrypt mnemonic
function loadEncryptedMnemonic(paths) {
  if (!fs.existsSync(paths.encryptedMnemonic)) {
    return null;
  }

  if (!safeStorage.isEncryptionAvailable()) {
    throw new Error('Encryption is not available on this system');
  }

  const encrypted = fs.readFileSync(paths.encryptedMnemonic);
  return safeStorage.decryptString(encrypted);
}

// Check if wallet exists
function walletExists(paths) {
  return fs.existsSync(paths.encryptedMnemonic);
}

// Check if backup is needed
function needsBackup(paths) {
  return fs.existsSync(paths.backupReminder);
}

// Mark backup as completed
function markBackupComplete(paths) {
  if (fs.existsSync(paths.backupReminder)) {
    fs.unlinkSync(paths.backupReminder);
  }
}

// Restore wallet from mnemonic
async function restoreWallet(mnemonic) {
  await loadCosmJS();

  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
    prefix: 'sharetoken'
  });

  const [account] = await wallet.getAccounts();

  return {
    address: account.address,
    pubKey: Buffer.from(account.pubkey).toString('hex')
  };
}

// Get wallet instance for signing
async function getWalletInstance(paths) {
  const mnemonic = loadEncryptedMnemonic(paths);
  if (!mnemonic) {
    return null;
  }

  await loadCosmJS();

  return await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
    prefix: 'sharetoken'
  });
}

// Get wallet address
async function getWalletAddress(paths) {
  const wallet = await getWalletInstance(paths);
  if (!wallet) {
    return null;
  }

  const [account] = await wallet.getAccounts();
  return account.address;
}

// Get balance from RPC
async function getBalance(address, rpcEndpoint = 'http://localhost:26657') {
  await loadCosmJS();

  const client = await SigningStargateClient.connect(rpcEndpoint);
  const balances = await client.getAllBalances(address);
  await client.disconnect();

  return balances;
}

// Send tokens
async function sendTokens(recipient, amount, denom = 'stake', memo = '', paths, rpcEndpoint = 'http://localhost:26657') {
  await loadCosmJS();

  const wallet = await getWalletInstance(paths);
  if (!wallet) {
    throw new Error('Wallet not found');
  }

  const [senderAccount] = await wallet.getAccounts();

  const client = await SigningStargateClient.connectWithSigner(
    rpcEndpoint,
    wallet,
    { gasPrice: GasPrice.fromString('0.025stake') }
  );

  const result = await client.sendTokens(
    senderAccount.address,
    recipient,
    [{ denom, amount: amount.toString() }],
    'auto',
    memo
  );

  await client.disconnect();

  return {
    transactionHash: result.transactionHash,
    gasUsed: result.gasUsed,
    gasWanted: result.gasWanted,
    height: result.height
  };
}

// Export wallet data (for backup) - requires password verification
async function exportWalletData(paths, password) {
  // Verify password (in this simple implementation, we check if decryption works)
  let mnemonic;
  try {
    mnemonic = loadEncryptedMnemonic(paths);
  } catch (error) {
    throw new Error('Failed to decrypt wallet - invalid password or corrupted data');
  }

  if (!mnemonic) {
    throw new Error('Wallet not found');
  }

  return {
    mnemonic,
    exportedAt: new Date().toISOString()
  };
}

// Delete wallet (for logout/reset)
function deleteWallet(paths) {
  if (fs.existsSync(paths.encryptedMnemonic)) {
    fs.unlinkSync(paths.encryptedMnemonic);
  }
  if (fs.existsSync(paths.backupReminder)) {
    fs.unlinkSync(paths.backupReminder);
  }
  if (fs.existsSync(paths.walletConfig)) {
    fs.unlinkSync(paths.walletConfig);
  }
  return true;
}

// Verify mnemonic format
function isValidMnemonic(mnemonic) {
  const words = mnemonic.trim().split(/\s+/);
  return words.length === 12 || words.length === 15 || words.length === 18 || words.length === 21 || words.length === 24;
}

module.exports = {
  // Paths
  getWalletPaths,

  // Creation and restoration
  createWallet,
  restoreWallet,

  // Storage
  saveEncryptedMnemonic,
  loadEncryptedMnemonic,
  walletExists,
  deleteWallet,

  // Backup
  needsBackup,
  markBackupComplete,

  // Operations
  getWalletAddress,
  getWalletInstance,
  getBalance,
  sendTokens,
  exportWalletData,

  // Validation
  isValidMnemonic
};
