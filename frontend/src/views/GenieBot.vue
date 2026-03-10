<template>
  <div class="geniebot-container">
    <div class="chat-header">
      <h2>🧞 GenieBot</h2>
      <p class="subtitle">AI Assistant - Describe what you need, I'll find the best service</p>
    </div>

    <div class="model-selector">
      <label>AI Model:</label>
      <select v-model="selectedModel" class="model-select">
        <option value="openai/gpt-4">OpenAI GPT-4</option>
        <option value="openai/gpt-3.5-turbo">OpenAI GPT-3.5</option>
        <option value="anthropic/claude-3">Anthropic Claude 3</option>
        <option value="anthropic/claude-instant">Claude Instant</option>
        <option value="google/gemini-pro">Google Gemini Pro</option>
      </select>
      <span class="cost-estimate">Est. cost: {{ estimatedCost }} STT/token</span>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div v-for="(message, index) in messages" :key="index" class="message" :class="message.role">
        <div class="message-avatar">
          {{ message.role === 'user' ? '👤' : '🧞' }}
        </div>
        <div class="message-content">
          <div class="message-text">{{ message.content }}</div>
          <div v-if="message.serviceRecommendations" class="service-recommendations">
            <h4>Recommended Services:</h4>
            <div v-for="service in message.serviceRecommendations" :key="service.id" class="service-card" @click="selectService(service)">
              <div class="service-header">
                <span class="service-name">{{ service.name }}</span>
                <span class="service-price">{{ service.price }} STT</span>
              </div>
              <div class="service-description">{{ service.description }}</div>
              <div class="service-meta">
                <span class="provider">by {{ service.provider }}</span>
                <span class="rating">⭐ {{ service.rating }}</span>
              </div>
            </div>
          </div>
          <div v-if="message.intent" class="intent-tag">
            Detected: {{ message.intent }}
          </div>
          <div class="message-meta">
            <span class="timestamp">{{ formatTime(message.timestamp) }}</span>
            <span v-if="message.tokens" class="tokens">{{ message.tokens }} tokens</span>
            <span v-if="message.cost" class="cost">{{ message.cost }} STT</span>
          </div>
        </div>
      </div>

      <div v-if="isTyping" class="message assistant typing">
        <div class="message-avatar">🧞</div>
        <div class="message-content">
          <div class="typing-indicator">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>
    </div>

    <div class="chat-input-area">
      <div class="quick-actions">
        <button @click="sendQuick('Write a blog post about blockchain')" class="quick-btn">📝 Write blog</button>
        <button @click="sendQuick('Code review for my smart contract')" class="quick-btn">💻 Code review</button>
        <button @click="sendQuick('Research on DeFi protocols')" class="quick-btn">🔍 Research</button>
        <button @click="sendQuick('Translate this to Chinese')" class="quick-btn">🌐 Translate</button>
      </div>

      <div class="input-container">
        <textarea
          v-model="inputMessage"
          @keydown.enter.prevent="sendMessage"
          placeholder="Describe what you need... (e.g., 'Write a Python script to analyze blockchain data')"
          rows="3"
          class="chat-input"
        ></textarea>
        <button @click="sendMessage" :disabled="!inputMessage.trim() || isTyping" class="send-btn">
          {{ isTyping ? '...' : '➤' }}
        </button>
      </div>

      <div class="input-footer">
        <span class="hint">Press Enter to send, Shift+Enter for new line</span>
        <span class="balance" v-if="userBalance !== null">Balance: {{ userBalance }} STT</span>
      </div>
    </div>

    <!-- Service Details Modal -->
    <div v-if="selectedService" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h3>{{ selectedService.name }}</h3>
        <p class="modal-description">{{ selectedService.description }}</p>
        <div class="modal-details">
          <div class="detail-row">
            <span>Price:</span>
            <strong>{{ selectedService.price }} STT</strong>
          </div>
          <div class="detail-row">
            <span>Provider:</span>
            <span>{{ selectedService.provider }}</span>
          </div>
          <div class="detail-row">
            <span>Rating:</span>
            <span>⭐ {{ selectedService.rating }}/5</span>
          </div>
          <div class="detail-row">
            <span>Est. Time:</span>
            <span>{{ selectedService.estimatedTime }}</span>
          </div>
        </div>
        <div class="modal-actions">
          <button @click="confirmService" class="btn-primary">Confirm & Pay</button>
          <button @click="closeModal" class="btn-secondary">Cancel</button>
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
          content: 'Hello! I\'m GenieBot, your AI assistant. I can help you with:\n\n• Writing and content creation\n• Code review and development\n• Research and analysis\n• Translation and summarization\n• And much more!\n\nJust describe what you need, and I\'ll find the best service for you.',
          timestamp: Date.now(),
          intent: 'greeting'
        }
      ],
      inputMessage: '',
      selectedModel: 'openai/gpt-4',
      isTyping: false,
      userBalance: null,
      selectedService: null,
      serviceRecommendations: []
    };
  },
  computed: {
    estimatedCost() {
      const costs = {
        'openai/gpt-4': '0.03',
        'openai/gpt-3.5-turbo': '0.002',
        'anthropic/claude-3': '0.008',
        'anthropic/claude-instant': '0.0016',
        'google/gemini-pro': '0.001'
      };
      return costs[this.selectedModel] || '0.01';
    }
  },
  mounted() {
    this.scrollToBottom();
    this.fetchUserBalance();
  },
  methods: {
    async sendMessage() {
      if (!this.inputMessage.trim() || this.isTyping) return;

      const userMessage = {
        role: 'user',
        content: this.inputMessage.trim(),
        timestamp: Date.now()
      };

      this.messages.push(userMessage);
      this.inputMessage = '';
      this.isTyping = true;
      this.scrollToBottom();

      // Simulate AI response with intent recognition
      setTimeout(() => {
        this.generateResponse(userMessage.content);
      }, 1500);
    },

    sendQuick(text) {
      this.inputMessage = text;
      this.sendMessage();
    },

    generateResponse(userInput) {
      // Simple intent recognition simulation
      const intents = this.recognizeIntent(userInput);
      const response = this.createResponse(intents, userInput);

      this.messages.push({
        role: 'assistant',
        content: response.text,
        serviceRecommendations: response.recommendations,
        intent: intents.primary,
        tokens: Math.floor(userInput.length / 4) + 50,
        cost: (Math.floor(userInput.length / 4) + 50) * parseFloat(this.estimatedCost),
        timestamp: Date.now()
      });

      this.isTyping = false;
      this.scrollToBottom();
    },

    recognizeIntent(text) {
      const lowerText = text.toLowerCase();
      const intents = [];

      // Writing intent
      if (lowerText.includes('write') || lowerText.includes('blog') || lowerText.includes('article') || lowerText.includes('content')) {
        intents.push({ type: 'writing', confidence: 0.92 });
      }

      // Code intent
      if (lowerText.includes('code') || lowerText.includes('program') || lowerText.includes('script') || lowerText.includes('review')) {
        intents.push({ type: 'coding', confidence: 0.88 });
      }

      // Research intent
      if (lowerText.includes('research') || lowerText.includes('analyze') || lowerText.includes('study') || lowerText.includes('investigate')) {
        intents.push({ type: 'research', confidence: 0.85 });
      }

      // Translation intent
      if (lowerText.includes('translate') || lowerText.includes('translation') || lowerText.includes('language')) {
        intents.push({ type: 'translation', confidence: 0.90 });
      }

      // Default intent
      if (intents.length === 0) {
        intents.push({ type: 'general', confidence: 0.70 });
      }

      return {
        primary: intents[0].type,
        confidence: intents[0].confidence,
        all: intents
      };
    },

    createResponse(intents, userInput) {
      const responses = {
        writing: {
          text: 'I can help you with writing! Here are some professional writing services available:',
          recommendations: [
            { id: '1', name: 'Pro Content Writer', description: 'High-quality blog posts, articles, and marketing copy', price: '50', provider: 'AliceWriter', rating: 4.8, estimatedTime: '2 hours' },
            { id: '2', name: 'Technical Writer', description: 'Documentation, whitepapers, and technical guides', price: '80', provider: 'TechWrite Pro', rating: 4.9, estimatedTime: '4 hours' },
            { id: '3', name: 'Creative Writing', description: 'Stories, scripts, and creative content', price: '60', provider: 'CreativeAI', rating: 4.7, estimatedTime: '3 hours' }
          ]
        },
        coding: {
          text: 'I found several coding experts who can help you! Check out these services:',
          recommendations: [
            { id: '4', name: 'Smart Contract Auditor', description: 'Security audit for blockchain smart contracts', price: '200', provider: 'SecureCode', rating: 4.9, estimatedTime: '24 hours' },
            { id: '5', name: 'Python Developer', description: 'Data analysis scripts and automation tools', price: '100', provider: 'PyExpert', rating: 4.8, estimatedTime: '6 hours' },
            { id: '6', name: 'Code Reviewer', description: 'Professional code review and optimization', price: '75', provider: 'CodeQuality', rating: 4.7, estimatedTime: '4 hours' }
          ]
        },
        research: {
          text: 'For research tasks, here are specialized services:',
          recommendations: [
            { id: '7', name: 'DeFi Analyst', description: 'Deep dive into DeFi protocols and yield strategies', price: '150', provider: 'DeFiResearch', rating: 4.8, estimatedTime: '12 hours' },
            { id: '8', name: 'Market Researcher', description: 'Market analysis and competitive intelligence', price: '120', provider: 'MarketIntel', rating: 4.6, estimatedTime: '8 hours' }
          ]
        },
        translation: {
          text: 'Translation services available in 50+ languages:',
          recommendations: [
            { id: '9', name: 'Professional Translator', description: 'Native-level translation with cultural adaptation', price: '30', provider: 'LinguaPro', rating: 4.9, estimatedTime: '2 hours' },
            { id: '10', name: 'Technical Translator', description: 'Technical documents and manuals', price: '50', provider: 'TechLingua', rating: 4.7, estimatedTime: '4 hours' }
          ]
        },
        general: {
          text: 'Here are some popular AI services that might help you:',
          recommendations: [
            { id: '11', name: 'AI Assistant', description: 'General-purpose AI assistance', price: '20', provider: 'GenieBot', rating: 4.5, estimatedTime: '1 hour' },
            { id: '12', name: 'Data Analyst', description: 'Data processing and visualization', price: '100', provider: 'DataWiz', rating: 4.8, estimatedTime: '6 hours' }
          ]
        }
      };

      return responses[intents.primary] || responses.general;
    },

    selectService(service) {
      this.selectedService = service;
    },

    closeModal() {
      this.selectedService = null;
    },

    confirmService() {
      // Add confirmation message
      this.messages.push({
        role: 'assistant',
        content: `Service "${this.selectedService.name}" selected! The task has been submitted to the provider. You'll be notified when it's complete.`,
        timestamp: Date.now()
      });
      this.closeModal();
      this.scrollToBottom();
    },

    async fetchUserBalance() {
      // TODO: Fetch actual balance from wallet
      this.userBalance = 1000; // Placeholder
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

.model-selector {
  padding: 1rem 1.5rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.model-select {
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

.intent-tag {
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: #17a2b8;
  font-weight: bold;
}

.service-recommendations {
  margin-top: 1rem;
}

.service-recommendations h4 {
  margin: 0 0 0.75rem;
  color: #495057;
  font-size: 0.9rem;
}

.service-card {
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 0.75rem;
  cursor: pointer;
  transition: all 0.2s;
}

.service-card:hover {
  border-color: #667eea;
  box-shadow: 0 4px 8px rgba(102, 126, 234, 0.15);
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.service-name {
  font-weight: bold;
  color: #2c3e50;
}

.service-price {
  color: #28a745;
  font-weight: bold;
}

.service-description {
  color: #6c757d;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

.service-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.8rem;
  color: #adb5bd;
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

/* Typing indicator */
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

.modal-description {
  color: #6c757d;
  margin-bottom: 1.5rem;
}

.modal-details {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #e9ecef;
}

.detail-row:last-child {
  border-bottom: none;
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
