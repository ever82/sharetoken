<template>
  <div class="home-view">
    <div class="hero">
      <h1>欢迎来到 ShareToken</h1>
      <p>去中心化的 AI 服务市场</p>
      <div class="actions">
        <router-link to="/wallet" class="btn btn-primary">打开钱包</router-link>
        <router-link to="/market" class="btn btn-secondary">浏览市场</router-link>
      </div>
    </div>
    <div class="features">
      <div class="feature">
        <h3>🔐 安全钱包</h3>
        <p>私钥自持，完全掌控您的资产</p>
      </div>
      <div class="feature">
        <h3>🤖 AI 服务</h3>
        <p>一键调用主流 AI 模型</p>
      </div>
      <div class="feature">
        <h3>💰 安全托管</h3>
        <p>资金托管，确认后付款</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.hero {
  text-align: center;
  padding: 4rem 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 16px;
  margin-bottom: 3rem;
}

.hero h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.hero p {
  font-size: 1.25rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.btn {
  padding: 0.875rem 2rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: transform 0.2s;
}

.btn:hover {
  transform: translateY(-2px);
}

.btn-primary {
  background: white;
  color: #667eea;
}

.btn-secondary {
  background: transparent;
  color: white;
  border: 2px solid white;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.feature {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.feature h3 {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.feature p {
  color: #666;
}
</style>
