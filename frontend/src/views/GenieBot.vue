<template>
  <div class="geniebot-container">
    <div class="chat-header">
      <h2>🧞 GenieBot</h2>
      <p class="subtitle">AI Assistant - Connected to local Agent Gateway</p>
    </div>

    <div class="connection-status">
      <span :class="['status-dot', connectionStatus]"></span>
      <span class="status-text">{{ connectionStatusText }}</span>
      <span v-if="agentGatewayUrl" class="gateway-url">@ {{ agentGatewayUrl }}</span>
    </div>

    <div class="agent-selector">
      <label>AI Agent:</label>
      <select v-model="selectedAgent" class="agent-select">
        <option value="local">🧞 GenieBot (Local)</option>
        <option value="claude-code" :disabled="!hasClaudeCode">🤖 Claude Code {{ hasClaudeCode ? '' : '(Not Available)' }}</option>
        <option value="custom">⚙️ Custom MCP Agent</option>
      </select>
      <span class="cost-estimate" v-if="selectedAgent === 'local'">Cost: 0.001 STT/message</span>
      <span class="cost-estimate free" v-else>Free (uses your own agent)</span>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div v-for="(message, index) in messages" :key="index" class="message" :class="message.role">
        <div class="message-avatar">
          {{ message.role === 'user' ? '👤' : message.agent === 'claude-code' ? '🤖' : '🧞' }}
        </div>
        <div class="message-content">
          <div class="message-text">{{ message.content }}</div>
          <div v-if="message.toolCalls" class="tool-calls">
            <div v-for="(tool, idx) in message.toolCalls" :key="idx" class="tool-call">
              <span class="tool-icon">🔧</span>
              <span class="tool-name">{{ tool.name }}</span>
              <span class="tool-status">{{ tool.status }}</span>
            </div>
          </div>
          <div class="message-meta">
            <span class="timestamp">{{ formatTime(message.timestamp) }}</span>
            <span v-if="message.cost !== undefined" class="cost" :class="{ free: message.cost === 0 }">
              {{ message.cost === 0 ? 'Free' : message.cost + ' STT' }}
            </span>
            <span v-if="message.agent" class="agent-tag">{{ message.agent }}</span>
          </div>
        </div>
      </div>

      <div v-if="isTyping" class="message assistant typing">
        <div class="message-avatar">{{ selectedAgent === 'claude-code' ? '🤖' : '🧞' }}</div>
        <div class="message-content">
          <div class="typing-indicator">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>

      <div v-if="error" class="error-message">
        <span class="error-icon">⚠️</span>
        {{ error }}
      </div>
    </div>

    <div class="chat-input-area">
      <div class="quick-actions">
        <button @click="sendQuick('查询我的余额')" class="quick-btn">💰 查余额</button>
        <button @click="sendQuick('创建一个新任务')" class="quick-btn">📋 创建任务</button>
        <button @click="sendQuick('查看可用服务')" class="quick-btn">🔍 浏览服务</button>
        <button @click="sendQuick('我想做任务赚STT')" class="quick-btn">💎 赚STT</button>
      </div>

      <div class="input-container">
        <textarea
          v-model="inputMessage"
          @keydown.enter.prevent="sendMessage"
          placeholder="输入消息... (支持连接本地 Claude Code 或 OpenClaw)"
          rows="3"
          class="chat-input"
          :disabled="isTyping"
        ></textarea>
        <button @click="sendMessage" :disabled="!inputMessage.trim() || isTyping" class="send-btn">
          {{ isTyping ? '...' : '➤' }}
        </button>
      </div>

      <div class="input-footer">
        <span class="hint">Press Enter to send</span>
        <span class="balance" v-if="userBalance !== null">Balance: {{ userBalance }} STT</span>
      </div>
    </div>

    <!-- Agent Configuration Modal -->
    <div v-if="showAgentConfig" class="modal-overlay" @click.self="closeAgentConfig">
      <div class="modal">
        <h3>Configure Local Agent</h3>
        <div class="modal-content">
          <div class="form-group">
            <label>Agent Type:</label>
            <select v-model="customAgentConfig.type" class="form-select">
              <option value="stdio">Stdio (Claude Code, OpenClaw)</option>
              <option value="http">HTTP Endpoint</option>
            </select>
          </div>
          <div class="form-group">
            <label>Command/URL:</label>
            <input
              v-model="customAgentConfig.command"
              :placeholder="customAgentConfig.type === 'stdio' ? 'claude' : 'http://localhost:8080/mcp'"
              class="form-input"
            />
          </div>
          <div class="form-group" v-if="customAgentConfig.type === 'stdio'">
            <label>Arguments (comma separated):</label>
            <input
              v-model="customAgentConfig.args"
              placeholder="--transport,stdio"
              class="form-input"
            />
          </div>
        </div>
        <div class="modal-actions">
          <button @click="saveAgentConfig" class="btn-primary">Save</button>
          <button @click="closeAgentConfig" class="btn-secondary">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GenieBot',
  data() {
    return {
      messages: [
        {
          role: 'assistant',
          content: '你好！我是 GenieBot，已连接到本地 Agent Gateway。\n\n我可以帮你：\n• 查询余额、创建任务\n• 使用本地 Claude Code / OpenClaw (免费)\n• 浏览服务市场并接单赚 STT\n\n直接输入你的需求即可开始！',
          timestamp: Date.now(),
          cost: 0,
          agent: 'local'
        }
      ],
      inputMessage: '',
      selectedAgent: 'local',
      isTyping: false,
      userBalance: null,
      error: null,
      connectionStatus: 'connecting',
      agentGatewayUrl: '',
      hasClaudeCode: false,
      showAgentConfig: false,
      customAgentConfig: {
        type: 'stdio',
        command: '',
        args: ''
      },
      sessionId: 'geniebot-session-' + Date.now()
    };
  },
  computed: {
    connectionStatusText() {
      const texts = {
        connected: 'Connected',
        connecting: 'Connecting...',
        error: 'Connection Error'
      };
      return texts[this.connectionStatus] || 'Unknown';
    }
  },
  async mounted() {
    this.scrollToBottom();
    await this.checkAgentGateway();
    await this.fetchUserBalance();
    await this.checkLocalAgents();
  },
  methods: {
    async checkAgentGateway() {
      // Try different possible gateway URLs
      const urls = [
        'http://localhost:18080',
        'http://127.0.0.1:18080',
        window.location.origin + ':18080'
      ];

      for (const url of urls) {
        try {
          const response = await fetch(url + '/health', {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
          });

          if (response.ok) {
            const data = await response.json();
            this.connectionStatus = 'connected';
            this.agentGatewayUrl = url;
            console.log('Agent Gateway connected:', data);
            return;
          }
        } catch (e) {
          console.log('Failed to connect to:', url);
        }
      }

      this.connectionStatus = 'error';
      this.error = 'Agent Gateway not available. Please start it with: ./bin/agent-gateway -transport=http -port=18080';
    },

    async checkLocalAgents() {
      // Check if user has Claude Code configured
      const savedConfig = localStorage.getItem('geniebot_agent_config');
      if (savedConfig) {
        this.customAgentConfig = JSON.parse(savedConfig);
      }

      // Check for Claude Code MCP config
      const claudeConfig = localStorage.getItem('claude_mcp_config');
      if (claudeConfig) {
        this.hasClaudeCode = true;
      }

      // Also try to detect if claude command is available (in Electron)
      if (window.electronAPI && window.electronAPI.checkClaudeCode) {
        try {
          const result = await window.electronAPI.checkClaudeCode();
          this.hasClaudeCode = result.available;
        } catch (e) {
          console.log('Could not check Claude Code:', e);
        }
      }
    },

    async sendMessage() {
      if (!this.inputMessage.trim() || this.isTyping) return;

      const userMessage = {
        role: 'user',
        content: this.inputMessage.trim(),
        timestamp: Date.now()
      };

      this.messages.push(userMessage);
      const messageText = this.inputMessage.trim();
      this.inputMessage = '';
      this.isTyping = true;
      this.error = null;
      this.scrollToBottom();

      try {
        if (this.selectedAgent === 'local') {
          await this.callLocalGenieBot(messageText);
        } else if (this.selectedAgent === 'claude-code') {
          await this.callClaudeCode(messageText);
        } else {
          await this.callCustomAgent(messageText);
        }
      } catch (err) {
        this.error = 'Error: ' + err.message;
        this.messages.push({
          role: 'assistant',
          content: 'Sorry, I encountered an error: ' + err.message,
          timestamp: Date.now(),
          cost: 0,
          agent: 'error'
        });
      } finally {
        this.isTyping = false;
        this.scrollToBottom();
      }
    },

    async callLocalGenieBot(message) {
      // Call Agent Gateway MCP endpoint
      const response = await fetch(this.agentGatewayUrl + '/mcp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          jsonrpc: '2.0',
          method: 'tools/call',
          id: Date.now(),
          params: {
            name: 'chat_with_genie',
            arguments: {
              message: message,
              session_id: this.sessionId
            }
          }
        })
      });

      const data = await response.json();

      if (data.error) {
        throw new Error(data.error.message);
      }

      const resultText = data.result?.content?.[0]?.text;
      if (!resultText) {
        throw new Error('Invalid response format');
      }

      const result = JSON.parse(resultText);

      this.messages.push({
        role: 'assistant',
        content: result.content,
        timestamp: Date.now(),
        cost: result.cost,
        agent: 'local'
      });

      // Update balance if there's a cost
      if (result.cost && this.userBalance !== null) {
        this.userBalance = Math.max(0, this.userBalance - result.cost);
      }
    },

    async callClaudeCode(message) {
      // This would require a backend proxy to call Claude Code MCP
      // For now, show instructions
      this.messages.push({
        role: 'assistant',
        content: `🤖 **Claude Code Integration**\n\nTo use your local Claude Code:\n\n1. Configure Claude Code MCP:\n\`\`\`json\n{\n  "mcpServers": {\n    "sharetoken": {\n      "command": "/path/to/agent-gateway",\n      "args": ["-transport", "stdio"]\n    }\n  }\n}\n\`\`\`\n\n2. Or run Claude Code with:\n\`claude --mcp-server sharetoken\`\n\nThis lets you use Claude Code for free while still earning STT by completing tasks!`,
        timestamp: Date.now(),
        cost: 0,
        agent: 'claude-code'
      });
    },

    async callCustomAgent(message) {
      this.messages.push({
        role: 'assistant',
        content: `⚙️ **Custom Agent: ${this.customAgentConfig.command}**\n\nMessage: "${message}"\n\nCustom agent integration requires backend support. Please configure in settings.`,
        timestamp: Date.now(),
        cost: 0,
        agent: 'custom'
      });
    },

    sendQuick(text) {
      this.inputMessage = text;
      this.sendMessage();
    },

    async fetchUserBalance() {
      if (!this.agentGatewayUrl) return;

      try {
        // Get wallet address from electron API or use a default
        let address = 'cosmos1user';
        if (window.electronAPI && window.electronAPI.walletGetAddress) {
          address = await window.electronAPI.walletGetAddress();
        }

        const response = await fetch(this.agentGatewayUrl + '/mcp', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            jsonrpc: '2.0',
            method: 'tools/call',
            id: Date.now(),
            params: {
              name: 'query_balance',
              arguments: { address }
            }
          })
        });

        const data = await response.json();
        if (data.result?.content?.[0]?.text) {
          const result = JSON.parse(data.result.content[0].text);
          this.userBalance = result.balance / 1000000; // Convert to STT
        }
      } catch (e) {
        console.log('Could not fetch balance:', e);
        this.userBalance = 0;
      }
    },

    showAgentConfiguration() {
      this.showAgentConfig = true;
    },

    closeAgentConfig() {
      this.showAgentConfig = false;
    },

    saveAgentConfig() {
      localStorage.setItem('geniebot_agent_config', JSON.stringify(this.customAgentConfig));
      this.showAgentConfig = false;
      this.hasClaudeCode = true;
    },

    scrollToBottom() {
      const container = this.$refs.messagesContainer;
      if (container) {
        setTimeout(() => {
          container.scrollTop = container.scrollHeight;
        }, 100);
      }
    },

    formatTime(timestamp) {
      return new Date(timestamp).toLocaleTimeString();
    }
  }
};
</script>

<style scoped>
.geniebot-container {
  max-width: 900px;
  margin: 0 auto;
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  background: white;
}

.chat-header {
  padding: 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-align: center;
}

.chat-header h2 {
  margin: 0;
  font-size: 1.8rem;
}

.subtitle {
  margin: 0.5rem 0 0;
  opacity: 0.9;
  font-size: 0.95rem;
}

.connection-status {
  padding: 0.5rem 1rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-dot.connected {
  background: #28a745;
}

.status-dot.connecting {
  background: #ffc107;
  animation: pulse 1.5s infinite;
}

.status-dot.error {
  background: #dc3545;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  color: #6c757d;
}

.gateway-url {
  color: #17a2b8;
  margin-left: auto;
  font-family: monospace;
}

.agent-selector {
  padding: 1rem 1.5rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.agent-select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
}

.cost-estimate {
  color: #28a745;
  font-weight: bold;
  font-size: 0.9rem;
}

.cost-estimate.free {
  color: #17a2b8;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  background: #f8f9fa;
}

.message {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.message-content {
  flex: 1;
  max-width: 80%;
}

.message.user .message-content {
  margin-left: auto;
}

.message-text {
  background: white;
  padding: 1rem;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  line-height: 1.5;
  white-space: pre-wrap;
}

.message.user .message-text {
  background: #667eea;
  color: white;
}

.message-meta {
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: #6c757d;
  display: flex;
  gap: 1rem;
}

.message.user .message-meta {
  justify-content: flex-end;
}

.cost {
  color: #28a745;
  font-weight: bold;
}

.cost.free {
  color: #17a2b8;
}

.agent-tag {
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.7rem;
}

.tool-calls {
  margin-top: 0.5rem;
}

.tool-call {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 0.25rem;
  font-size: 0.85rem;
}

.tool-icon {
  font-size: 1rem;
}

.tool-name {
  font-weight: bold;
  color: #495057;
}

.tool-status {
  margin-left: auto;
  color: #6c757d;
  font-size: 0.75rem;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  padding: 1rem;
  border-radius: 8px;
  margin: 1rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.error-icon {
  font-size: 1.2rem;
}

.chat-input-area {
  padding: 1rem 1.5rem;
  background: white;
  border-top: 1px solid #e9ecef;
}

.quick-actions {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.quick-btn {
  padding: 0.5rem 1rem;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

.input-container {
  display: flex;
  gap: 0.75rem;
  align-items: flex-end;
}

.chat-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid #dee2e6;
  border-radius: 12px;
  resize: none;
  font-family: inherit;
  font-size: 0.95rem;
}

.chat-input:focus {
  outline: none;
  border-color: #667eea;
}

.chat-input:disabled {
  background: #f8f9fa;
  cursor: not-allowed;
}

.send-btn {
  padding: 0.75rem 1.25rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: #5a6fd6;
}

.send-btn:disabled {
  background: #adb5bd;
  cursor: not-allowed;
}

.input-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: #adb5bd;
}

.balance {
  color: #28a745;
  font-weight: bold;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 1rem;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #adb5bd;
  border-radius: 50%;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-10px); }
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  max-width: 500px;
  width: 90%;
}

.modal h3 {
  margin: 0 0 1rem;
  color: #2c3e50;
}

.modal-content {
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #495057;
  font-weight: 500;
}

.form-select,
.form-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  font-size: 0.95rem;
}

.form-select:focus,
.form-input:focus {
  outline: none;
  border-color: #667eea;
}

.modal-actions {
  display: flex;
  gap: 1rem;
}

.btn-primary {
  flex: 1;
  padding: 0.75rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
}

.btn-secondary {
  flex: 1;
  padding: 0.75rem;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
}
</style>
