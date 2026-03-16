<template>
  <div class="credits-page">
    <NavBar title="我的额度" :show-back="true" back-to="/profile" />

    <div class="credits-content">
      <!-- 额度卡片 -->
      <div class="credits-card">
        <div class="credits-icon">◈</div>
        <div class="credits-info">
          <div class="credits-amount">{{ settings?.credits || 0 }}</div>
          <div class="credits-label">剩余额度</div>
        </div>
      </div>

      <!-- 统计数据 -->
      <div class="stats-grid">
        <div class="stat-box">
          <div class="stat-num">{{ settings?.total_tokens || 0 }}</div>
          <div class="stat-label">累计Token</div>
        </div>
        <div class="stat-box">
          <div class="stat-num">{{ settings?.streak || 0 }}</div>
          <div class="stat-label">连续签到</div>
        </div>
      </div>

      <!-- 签到区域 -->
      <div class="action-section">
        <div class="section-title">每日签到</div>
        <div class="signin-card">
          <div class="signin-info">
            <span class="signin-bonus">+{{ signinDaily }} / 天</span>
            <span class="streak-bonus">连续+{{ streakBonus }}</span>
          </div>
          <button
            class="signin-btn"
            :disabled="!canSignIn || signingIn"
            @click="handleSignIn"
          >
            {{ signingIn ? '签到中...' : (canSignIn ? '立即签到' : '已签到') }}
          </button>
        </div>
      </div>

      <!-- 邀请好友 -->
      <div class="action-section" v-if="settings?.settings?.invite_enabled">
        <div class="section-title">邀请好友</div>
        <div class="invite-card">
          <div class="invite-code-box">
            <span class="invite-label">邀请码</span>
            <div class="invite-code" @click="copyInviteCode">
              {{ inviteCode || '生成中...' }}
              <span class="copy-icon">📋</span>
            </div>
          </div>
          <div class="invite-stats">
            <span>已邀请 {{ inviteCount }} 人</span>
            <span>已获得 {{ rewardsClaimed }} 额度</span>
          </div>
        </div>
      </div>

      <!-- 兑换码 -->
      <div class="action-section">
        <div class="section-title">兑换码</div>
        <div class="redeem-card">
          <input
            v-model="redeemCode"
            type="text"
            placeholder="请输入兑换码"
            class="redeem-input"
          />
          <button class="redeem-btn" @click="handleRedeem" :disabled="!redeemCode || redeeming">
            {{ redeeming ? '兑换中...' : '兑换' }}
          </button>
        </div>
      </div>

      <!-- 绑定邀请码 -->
      <div class="action-section" v-if="!hasInvited">
        <div class="section-title">绑定邀请码</div>
        <div class="redeem-card">
          <input
            v-model="bindCode"
            type="text"
            placeholder="输入邀请码"
            class="redeem-input"
          />
          <button class="redeem-btn" @click="handleBindInvite" :disabled="!bindCode || binding">
            {{ binding ? '绑定中...' : '绑定' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCreditsStore } from '@/stores/credits'
import NavBar from '@/components/common/NavBar.vue'
import { showToast } from '@/utils/toast'

const creditsStore = useCreditsStore()

const settings = computed(() => creditsStore.settings)
const canSignIn = computed(() => settings.value?.can_signin ?? false)
const signinDaily = computed(() => settings.value?.settings?.signin_daily_credits ?? 10)
const streakBonus = computed(() => settings.value?.settings?.signin_streak_bonus ?? 5)

const inviteCode = ref('')
const inviteCount = ref(0)
const rewardsClaimed = ref(0)
const hasInvited = ref(false)

const redeemCode = ref('')
const bindCode = ref('')

const signingIn = ref(false)
const redeeming = ref(false)
const binding = ref(false)

async function loadData() {
  await creditsStore.fetchSettings()
  
  try {
    const inviteData = await creditsStore.getInviteCode()
    inviteCode.value = inviteData.invite_code
    inviteCount.value = inviteData.invite_count
    rewardsClaimed.value = inviteData.rewards_claimed
    
    const stored = localStorage.getItem('used_invite_code')
    hasInvited.value = !!stored
  } catch (e) {
    console.error('Failed to load invite code:', e)
  }
}

async function handleSignIn() {
  if (!canSignIn.value || signingIn.value) return
  
  signingIn.value = true
  try {
    const result = await creditsStore.signIn()
    showToast(`签到成功！获得 ${result?.credits_earned} 额度`)
  } catch (e: any) {
    showToast(e.message || '签到失败')
  } finally {
    signingIn.value = false
  }
}

async function copyInviteCode() {
  if (!inviteCode.value) return
  await navigator.clipboard.writeText(inviteCode.value)
  showToast('邀请码已复制')
}

async function handleRedeem() {
  if (!redeemCode.value || redeeming.value) return
  
  redeeming.value = true
  try {
    const result = await creditsStore.redeemCoupon(redeemCode.value)
    showToast(`兑换成功！获得 ${result.credits_earned} 额度`)
    redeemCode.value = ''
    await creditsStore.fetchSettings()
  } catch (e: any) {
    showToast(e.message || '兑换失败')
  } finally {
    redeeming.value = false
  }
}

async function handleBindInvite() {
  if (!bindCode.value || binding.value) return
  
  binding.value = true
  try {
    const result = await creditsStore.useInviteCode(bindCode.value)
    showToast(`绑定成功！获得 ${result.credits_earned} 额度`)
    localStorage.setItem('used_invite_code', 'true')
    hasInvited.value = true
    bindCode.value = ''
    await creditsStore.fetchSettings()
    await loadData()
  } catch (e: any) {
    showToast(e.message || '绑定失败')
  } finally {
    binding.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.credits-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

.credits-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.credits-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.credits-icon {
  width: 56px;
  height: 56px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: #fff;
}

.credits-info {
  flex: 1;
}

.credits-amount {
  font-size: 36px;
  font-weight: 700;
  color: #fff;
}

.credits-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.stat-box {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
}

.stat-num {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-tertiary);
}

.action-section {
  margin-bottom: 16px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 4px 8px;
}

.signin-card,
.invite-card,
.redeem-card {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 16px;
}

.signin-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.signin-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.signin-bonus {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.streak-bonus {
  font-size: 12px;
  color: var(--color-green);
}

.signin-btn {
  padding: 10px 20px;
  background: var(--color-green);
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.signin-btn:disabled {
  background: var(--bg-quaternary);
  color: var(--text-tertiary);
}

.invite-code-box {
  margin-bottom: 12px;
}

.invite-label {
  font-size: 12px;
  color: var(--text-tertiary);
  display: block;
  margin-bottom: 6px;
}

.invite-code {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-secondary);
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 2px;
}

.copy-icon {
  font-size: 16px;
  cursor: pointer;
}

.invite-stats {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-tertiary);
}

.redeem-card {
  display: flex;
  gap: 12px;
}

.redeem-input {
  flex: 1;
  padding: 12px;
  background: var(--bg-secondary);
  border: none;
  border-radius: 8px;
  font-size: 14px;
  color: var(--text-primary);
}

.redeem-input::placeholder {
  color: var(--text-quaternary);
}

.redeem-btn {
  padding: 12px 20px;
  background: var(--color-blue);
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.redeem-btn:disabled {
  background: var(--bg-quaternary);
  color: var(--text-tertiary);
}
</style>