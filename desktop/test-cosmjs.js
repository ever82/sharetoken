const fs = require('fs');

// The encrypted mnemonic is stored using safeStorage
// Since we can't easily decrypt it outside Electron, let's create a test wallet
// and verify the CosmJS integration works

async function testCosmJS() {
  console.log('=== Testing CosmJS Integration ===\n');
  
  try {
    const { DirectSecp256k1HdWallet } = require('@cosmjs/proto-signing');
    const { SigningStargateClient } = require('@cosmjs/stargate');
    
    // Test 1: Create a new wallet
    console.log('1. Creating test wallet...');
    const wallet = await DirectSecp256k1HdWallet.generate(12, {
      prefix: 'sharetoken'
    });
    const [account] = await wallet.getAccounts();
    console.log('   Address:', account.address);
    console.log('   Mnemonic:', wallet.mnemonic);
    
    // Test 2: Connect to node
    console.log('\n2. Connecting to ShareToken node...');
    const client = await SigningStargateClient.connect('http://localhost:26657');
    console.log('   Connected successfully');
    
    // Test 3: Query validator balance
    const validatorAddr = 'sharetoken1z843agmw5nvfrn46v3rkz3ehrphjpmy4m3vsy4';
    console.log('\n3. Querying validator balance...');
    const balances = await client.getAllBalances(validatorAddr);
    console.log('   Validator balances:', balances);
    
    // Test 4: Query node status
    console.log('\n4. Querying node status...');
    const status = await client.getChainId();
    console.log('   Chain ID:', status);
    
    await client.disconnect();
    
    console.log('\n=== CosmJS Integration Test Complete ===');
    console.log('\n✅ All tests passed!');
    
  } catch (err) {
    console.error('Test failed:', err.message);
    console.error(err.stack);
    process.exit(1);
  }
}

testCosmJS();
