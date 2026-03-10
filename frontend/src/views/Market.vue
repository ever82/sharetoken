<template>
  <div class="market-view">
    <h1>服务市场</h1>
    <p class="subtitle">发现和使用 AI 服务</p>

    <div class="services">
      <div class="service-card">
        <h3>GPT-4</h3>
        <p>OpenAI 最强大的语言模型</p>
        <span class="price">0.1 STT / 1K tokens</span>
        <button class="btn" disabled>即将上线</button>
      </div>

      <div class="service-card">
        <h3>Claude 3</h3>
        <p>Anthropic 的先进 AI 助手</p>
        <span class="price">0.08 STT / 1K tokens</span>
        <button class="btn" disabled>即将上线</button>
      </div>

      <div class="service-card">
        <h3>代码助手</h3>
        <p>专业的编程辅助 Agent</p>
        <span class="price">0.05 STT / 请求</span>
        <button class="btn" disabled>即将上线</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.market-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

h1 {
  text-align: center;
  margin-bottom: 0.5rem;
}

.subtitle {
  text-align: center;
  color: #666;
  margin-bottom: 3rem;
}

.services {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.service-card {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  text-align: center;
}

.service-card h3 {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.service-card p {
  color: #666;
  margin-bottom: 1rem;
}

.price {
  display: block;
  color: #667eea;
  font-weight: 600;
  margin-bottom: 1rem;
}

.btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}

.btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style>
