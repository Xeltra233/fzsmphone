<template>
  <div class="system-settings-page">
    <NavBar title="系统设置" :show-back="true" back-to="/profile" />

    <div class="page-content" v-if="authStore.isSuperAdmin || authStore.isAdmin">
      <!-- 标签切换 -->
      <div class="settings-tabs">
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'credits' }"
          @click="activeTab = 'credits'"
        >
          额度设置
        </button>
        <button
          class="tab-btn"
          :class="{ active: activeTab === 'system' }"
          @click="activeTab = 'system'"
        >
          系统配置
        </button>
<button
  class="tab-btn"
  :class="{ active: activeTab === 'api' }"
  @click="activeTab = 'api'"
>
  API设置
</button>
<button
  class="tab-btn"
  :class="{ active: activeTab === 'oauth' }"
  @click="activeTab = 'oauth'"
>
  OAuth配置
</button>
<button
  class="tab-btn"
  :class="{ active: activeTab === 'contact' }"
  @click="activeTab = 'contact'"
>
  联系方式
</button>
</div>

      <!-- 额度设置 -->
      <template v-if="activeTab === 'credits'">
        <div v-if="!authStore.isSuperAdmin" class="permission-notice">
          <span class="notice-icon">ℹ️</span>
          您的账户为管理员，只能查看额度设置。如需修改，请联系超级管理员。
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载中...</p>
        </div>

        <template v-else>
          <div class="settings-section">
            <div class="section-header">
              <h3>签到设置</h3>
            </div>
            <div class="settings-list">
              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">每日签到奖励</span>
                  <span class="label-desc">用户每日签到可获得的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.signin_daily_credits"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度</span>
                </div>
              </div>

              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">连续签到加成</span>
                  <span class="label-desc">连续签到时每日额外增加的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.signin_streak_bonus"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度/天</span>
                </div>
              </div>

              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">签到功能</span>
                  <span class="label-desc">是否启用每日签到功能</span>
                </div>
                <label class="toggle">
                  <input
                    type="checkbox"
                    v-model="creditSettings.signin_enabled"
                    :disabled="!authStore.isSuperAdmin"
                  />
                  <span class="toggle-slider"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="settings-section">
            <div class="section-header">
              <h3>邀请设置</h3>
            </div>
            <div class="settings-list">
              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">邀请奖励额度</span>
                  <span class="label-desc">成功邀请好友后奖励的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.invite_reward_credits"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度</span>
                </div>
              </div>

              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">邀请功能</span>
                  <span class="label-desc">是否启用邀请好友功能</span>
                </div>
                <label class="toggle">
                  <input
                    type="checkbox"
                    v-model="creditSettings.invite_enabled"
                    :disabled="!authStore.isSuperAdmin"
                  />
                  <span class="toggle-slider"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="settings-section">
            <div class="section-header">
              <h3>新用户设置</h3>
            </div>
            <div class="settings-list">
              <div class="setting-item">
                <div class="setting-label">
                  <span class="label-text">新用户初始额度</span>
                  <span class="label-desc">新注册用户自动获得的额度</span>
                </div>
                <div class="setting-input-wrap">
                  <input
                    type="number"
                    v-model.number="creditSettings.default_credits"
                    class="setting-input"
                    :disabled="!authStore.isSuperAdmin"
                    min="0"
                  />
                  <span class="input-suffix">额度</span>
                </div>
              </div>
            </div>
          </div>

          <div class="settings-section">
            <div class="section-header">
              <h3>兑换码管理</h3>
            </div>
            <div class="settings-list coupon-manager">
              <div class="setting-item vertical">
                <div class="setting-label">
                  <span class="label-text">生成兑换码</span>
                  <span class="label-desc">支持自定义或随机生成，生成后用户可在额度页直接兑换</span>
                </div>
                <div class="coupon-form-grid">
                  <input v-model.trim="couponForm.code" class="setting-input" :disabled="!authStore.isSuperAdmin || couponCreating" placeholder="留空则随机生成，例如 8 位兑换码" />
                  <input v-model.number="couponForm.credits" type="number" min="1" class="setting-input" :disabled="!authStore.isSuperAdmin || couponCreating" placeholder="奖励额度" />
                  <input v-model.number="couponForm.maxUses" type="number" min="1" class="setting-input" :disabled="!authStore.isSuperAdmin || couponCreating" placeholder="可使用次数" />
                  <input v-model.number="couponForm.expiresIn" type="number" min="0" class="setting-input" :disabled="!authStore.isSuperAdmin || couponCreating" placeholder="有效天数，0 为不过期" />
                  <input v-model.trim="couponForm.note" class="setting-input coupon-note-input" :disabled="!authStore.isSuperAdmin || couponCreating" placeholder="备注/用途，例如 开服补偿、活动礼包" />
                  <input v-model.number="couponForm.batchSize" type="number" min="1" max="100" class="setting-input" :disabled="!authStore.isSuperAdmin || couponCreating || !!couponForm.code.trim()" placeholder="批量数量，1-100" />
                </div>
                <div class="coupon-actions-row">
                  <button class="add-model-btn" type="button" @click="handleCreateCoupon" :disabled="!authStore.isSuperAdmin || couponCreating">
                    {{ couponCreating ? '生成中...' : '生成兑换码' }}
                  </button>
                  <button class="add-model-btn secondary" type="button" @click="fetchCoupons" :disabled="couponListLoading">
                    {{ couponListLoading ? '刷新中...' : '刷新列表' }}
                  </button>
                </div>
                <div v-if="latestCoupon" class="coupon-latest-card">
                  <span class="coupon-latest-label">最近生成</span>
                  <div v-if="latestCoupon.codes.length === 1" class="coupon-latest-code">{{ latestCoupon.codes[0] }}</div>
                  <div v-else class="coupon-latest-code coupon-latest-multi">共生成 {{ latestCoupon.codes.length }} 个兑换码</div>
                  <div class="coupon-latest-meta">额度 {{ latestCoupon.credits }} / 可用 {{ latestCoupon.max_uses }} 次 / {{ latestCoupon.note || '无备注' }}</div>
                  <div v-if="latestCoupon.codes.length > 1" class="coupon-latest-list">
                    <span v-for="code in latestCoupon.codes" :key="code" class="coupon-latest-chip">{{ code }}</span>
                  </div>
                </div>
              </div>

              <div class="setting-item vertical">
                <div class="setting-label">
                  <span class="label-text">最近兑换码</span>
                  <span class="label-desc">展示最近 50 条，方便核对是否过期或已用完</span>
                </div>
                <div class="coupon-toolbar">
                  <input v-model.trim="couponSearch" class="setting-input" placeholder="搜索兑换码或备注" />
                  <select v-model="couponStatusFilter" class="setting-input coupon-filter-select">
                    <option value="all">全部状态</option>
                    <option value="active">仅启用</option>
                    <option value="inactive">仅停用</option>
                    <option value="available">仅未用完</option>
                    <option value="exhausted">仅已用完</option>
                    <option value="expired">仅已过期</option>
                  </select>
                </div>
                <div v-if="couponList.length > 0" class="coupon-toolbar-meta">
                  <span>共 {{ couponList.length }} 条，当前显示 {{ filteredCouponList.length }} 条</span>
                  <div class="coupon-export-actions">
                    <button class="batch-action-btn" type="button" @click="exportCoupons('txt')">导出 TXT</button>
                    <button class="batch-action-btn" type="button" @click="exportCoupons('csv')">导出 CSV</button>
                  </div>
                </div>
                <div v-if="couponListLoading" class="empty-models">加载兑换码中...</div>
                <div v-else-if="couponList.length === 0" class="empty-models">暂无兑换码</div>
                <div v-else-if="filteredCouponList.length === 0" class="empty-models">没有符合筛选条件的兑换码</div>
                <div v-else class="coupon-list">
                  <div v-for="coupon in filteredCouponList" :key="coupon.code" class="coupon-list-item">
                    <div class="coupon-list-main">
                      <div class="coupon-list-code">{{ coupon.code }}</div>
                      <div class="coupon-list-meta">
                        <span>额度 {{ coupon.credits }}</span>
                        <span>已用 {{ coupon.current_uses }}/{{ coupon.max_uses }}</span>
                        <span>{{ formatCouponExpiry(coupon.expires_at) }}</span>
                        <span>{{ coupon.is_active ? '启用中' : '已停用' }}</span>
                        <span v-if="coupon.note">备注：{{ coupon.note }}</span>
                      </div>
                    </div>
                    <div class="coupon-item-actions">
                      <button class="batch-action-btn" type="button" @click="copyCouponCode(coupon.code)">复制</button>
                      <button class="batch-action-btn" type="button" @click="toggleCouponStatus(coupon)">{{ coupon.is_active ? '停用' : '启用' }}</button>
                      <button class="batch-action-btn danger" type="button" @click="deleteCoupon(coupon.code)">删除</button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <button
            class="save-btn"
            @click="saveCreditSettings"
            :disabled="saving || !authStore.isSuperAdmin"
          >
            <template v-if="saving">
              <div class="btn-spinner"></div>
              保存中...
            </template>
            <template v-else>
              ▣ 保存设置
            </template>
          </button>
        </template>
      </template>

<!-- 系统配置 -->
  <template v-if="activeTab === 'system'">
    <div v-if="!authStore.isSuperAdmin" class="permission-notice">
      <span class="notice-icon">ℹ️</span>
      您的账户为管理员，只能查看系统配置。如需修改，请联系超级管理员。
    </div>

    <div class="settings-section">
      <div class="section-header">
        <h3>基础设置</h3>
      </div>
      <div class="settings-list">
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">应用名称</span>
            <span class="label-desc">显示在登录页和状态栏</span>
          </div>
          <input
            v-model="systemSettings.app_name"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="贩子死妈小手机"
          />
        </div>
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">网页标题</span>
            <span class="label-desc">浏览器标签页显示的标题</span>
          </div>
          <input
            v-model="systemSettings.app_title"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="贩子死妈小手机"
          />
        </div>
      </div>
    </div>

    <div class="settings-section">
      <div class="section-header">
        <h3>公告提示</h3>
      </div>
      <div class="settings-list">
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">公告内容</span>
            <span class="label-desc">显示在登录页顶部，支持换行</span>
          </div>
          <textarea
            v-model="systemSettings.announcement"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="欢迎使用..."
            rows="3"
          />
        </div>
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">每日提示</span>
            <span class="label-desc">显示在额度页面的小技巧</span>
          </div>
          <textarea
            v-model="systemSettings.tips"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="每日签到可获得额外额度..."
            rows="2"
          />
        </div>
      </div>
    </div>

    <div class="settings-section">
      <div class="section-header">
        <h3>版权声明</h3>
      </div>
      <div class="settings-list">
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">免责声明</span>
            <span class="label-desc">作为附加声明并行显示，官方版权声明始终保留</span>
          </div>
          <textarea
            v-model="systemSettings.disclaimer"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
            placeholder="贩子死妈保留所有权益..."
            rows="3"
          />
        </div>
        <div class="setting-item vertical">
          <div class="setting-label">
            <span class="label-text">附加群二维码路径</span>
            <span class="label-desc">与官方群二维码并行显示，不会覆盖官方入口</span>
          </div>
          <input
            v-model="systemSettings.qun_qrcode"
            class="setting-input full-width"
            :disabled="!authStore.isSuperAdmin"
placeholder="qun_qrcode.jpg"
      />
      <div class="file-upload-area" v-if="authStore.isSuperAdmin">
        <input type="file" id="qrcode-upload" accept="image/*" @change="handleQrcodeUpload" hidden />
        <label for="qrcode-upload" class="upload-btn">上传二维码图片</label>
        <span v-if="uploadingQrcode" class="upload-status">上传中...</span>
      </div>
    </div>
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">网站图标</span>
        <span class="label-desc">浏览器标签页图标，显示在浏览器标签上</span>
      </div>
      <div class="favicon-preview">
        <img v-if="systemSettings.favicon" :src="systemSettings.favicon" alt="favicon" class="favicon-img" @error="faviconError = true" />
        <span v-else class="favicon-placeholder">未设置</span>
      </div>
      <div class="file-upload-area" v-if="authStore.isSuperAdmin">
        <input type="file" id="favicon-upload" accept="image/*,.ico" @change="handleFaviconUpload" hidden />
        <label for="favicon-upload" class="upload-btn">上传图标</label>
        <span v-if="uploadingFavicon" class="upload-status">上传中...</span>
      </div>
      <input
        v-model="systemSettings.favicon"
        class="setting-input full-width"
        :disabled="!authStore.isSuperAdmin"
        placeholder="/icon.png"
      />
    </div>
  </div>
</div>

    <button
      class="save-btn"
      @click="saveSystemSettings"
      :disabled="saving || !authStore.isSuperAdmin"
    >
      <template v-if="saving">
        <div class="btn-spinner"></div>
        保存中...
      </template>
      <template v-else>
        ▣ 保存设置
      </template>
    </button>
  </template>

<!-- API设置 -->
<template v-if="activeTab === 'api'">
<div class="api-settings-section">
  <div v-if="apiSettingsLoading" class="loading-state">
    <div class="spinner"></div>
    <p>加载API设置...</p>
  </div>
  <template v-else>
    <!-- 全局API设置 -->
    <div class="settings-card">
      <div class="card-header">
        <h3>全局API设置</h3>
        <span class="badge">全局</span>
      </div>
      <div class="card-body">
        <div class="input-group">
          <label>模型提供商</label>
          <select v-model="apiForm.globalApiUrl" class="setting-input">
            <option value="https://generativelanguage.googleapis.com/v1beta/openai/chat/completions">Gemini AI Studio</option>
            <option value="https://api.openai.com/v1/chat/completions">OpenAI</option>
            <option value="https://openrouter.ai/api/v1/chat/completions">OpenRouter</option>
            <option value="https://api.deepseek.com/chat/completions">DeepSeek</option>
            <option value="custom">自定义提供商</option>
          </select>
        </div>
        <div v-if="apiForm.globalApiUrl === 'custom'" class="input-group">
          <label>自定义接口地址</label>
          <input v-model="apiForm.globalCustomUrl" placeholder="https://api.example.com/v1/chat/completions" class="setting-input" />
        </div>
        <div class="input-group">
          <label>API Key <span class="input-tip">全局用户共用</span></label>
          <input v-model="apiForm.globalApiKey" type="password" placeholder="留空则强制用户填写个人API" class="setting-input" />
        </div>
        <div class="input-group model-manager-group">
          <label>默认模型 <span class="input-tip">用户未选择时的默认</span></label>
          <select v-model="apiForm.globalModel" class="setting-select">
            <option value="">请选择默认模型</option>
            <option v-for="model in enabledModels" :key="model.id" :value="model.id">
              {{ model.id }}
            </option>
          </select>
          <span class="input-desc">默认模型必须来自已启用模型列表</span>
        </div>
        <div class="input-group model-manager-group">
          <label>可用模型列表 <span class="input-tip">可添加、开关与删除</span></label>
          <div class="model-add-row">
            <input
              v-model="newModelId"
              placeholder="输入模型 ID，如 gpt-4o-mini"
              class="setting-input"
              @keyup.enter="addModel"
            />
            <button class="add-model-btn" type="button" @click="addModel">添加</button>
            <button class="add-model-btn secondary" type="button" @click="pullModels" :disabled="pullingModels || showModelPullModal">
              {{ pullingModels ? '拉取中...' : '拉取模型' }}
            </button>
          </div>
          <div v-if="modelPullError" class="model-pull-error">{{ modelPullError }}</div>
          <div v-if="apiForm.globalModels.length > 0" class="model-list-editor">
            <div v-for="model in apiForm.globalModels" :key="model.id" class="model-chip-row">
              <label class="model-toggle">
                <input v-model="model.enabled" type="checkbox" />
                <span class="model-toggle-slider"></span>
              </label>
              <div class="model-chip-main">
                <span class="model-chip-id">{{ model.id }}</span>
                <input
                  v-model.trim="model.displayName"
                  class="setting-input model-display-input"
                  placeholder="显示名称，可留空"
                />
              </div>
              <button class="model-delete-btn" type="button" @click="removeModel(model.id)">删除</button>
            </div>
          </div>
          <div v-else class="empty-models">暂未添加模型</div>
          <span class="input-desc">支持从当前 API 拉取模型；显示名称仅用于后台展示，留空则直接显示模型 ID</span>
        </div>
          <div class="input-row">
            <div class="input-group">
              <label>Temperature <span class="input-tip">创造性</span></label>
              <input v-model.number="apiForm.globalTemperature" type="number" step="0.1" min="0" max="2" class="setting-input" />
            </div>
            <div class="input-group">
              <label>Max Length <span class="input-tip">最大回复长度</span></label>
              <input v-model.number="apiForm.globalMaxLength" type="number" min="100" max="32000" class="setting-input" />
            </div>
          </div>
          <div class="input-row">
            <div class="input-group">
              <label>Context Size <span class="input-tip">历史消息数</span></label>
              <input v-model.number="apiForm.globalContextSize" type="number" min="1" max="100" class="setting-input" />
            </div>
            <div class="input-group">
<label>Timeout <span class="input-tip">超时秒数</span></label>
<input v-model.number="apiForm.globalTimeout" type="number" min="10" max="300" class="setting-input" />
</div>
</div>
</div>
</div>

<!-- 社交内容全局API -->
    <div class="settings-card">
      <div class="card-header">
        <h3>社交内容全局API</h3>
        <span class="badge">全局</span>
      </div>
      <div class="card-body">
        <div class="input-group">
          <label>社交内容 API 地址</label>
          <input v-model="apiForm.globalSocialApiUrl" placeholder="留空则使用主接口地址" class="setting-input" />
        </div>
        <div class="input-group">
          <label>社交内容 API Key <span class="input-tip">用于微博、邮箱等社交内容生成</span></label>
          <input v-model="apiForm.globalSocialApiKey" type="password" placeholder="留空则使用主API" class="setting-input" />
        </div>
<div class="input-group">
<label>社交内容默认模型</label>
<select v-model="apiForm.globalSocialModel" class="setting-select">
  <option value="">跟随主默认模型</option>
  <option v-for="model in enabledModels" :key="`social-${model.id}`" :value="model.id">
    {{ model.displayName || model.id }}
  </option>
</select>
<span class="input-desc">默认从已启用模型中选择，留空则跟随主默认模型</span>
</div>
</div>
</div>

<!-- 全局生图配置 -->
<div class="settings-card">
<div class="card-header">
<h3>全局生图配置</h3>
<span class="badge">全局</span>
</div>
<div class="card-body">
<div class="input-group">
<label>生图 API 格式</label>
<select v-model="imgGenForm.apiFormat" class="setting-select">
<option value="novelai">NovelAI</option>
<option value="openai">OpenAI/自定义</option>
<option value="gemini">Gemini</option>
</select>
</div>
<template v-if="imgGenForm.apiFormat === 'openai'">
<div class="input-group">
<label>生图 Base URL</label>
<input v-model="imgGenForm.openaiUrl" placeholder="https://api.openai.com" class="setting-input" />
</div>
<div class="input-group">
<label>生图 API Key</label>
<input v-model="imgGenForm.openaiKey" type="password" placeholder="sk-..." class="setting-input" />
</div>
<div class="input-group">
<label>生图模型</label>
<input v-model="imgGenForm.openaiModel" placeholder="dall-e-3 / flux" class="setting-input" />
</div>
</template>
<template v-else-if="imgGenForm.apiFormat === 'novelai'">
<div class="input-group">
<label>NovelAI 地址</label>
<input v-model="imgGenForm.novelaiUrl" placeholder="https://image.novelai.run" class="setting-input" />
</div>
<div class="input-group">
<label>NovelAI Key</label>
<input v-model="imgGenForm.novelaiKey" type="password" placeholder="user-key" class="setting-input" />
</div>
<div class="input-group">
<label>模型</label>
<input v-model="imgGenForm.novelaiModel" placeholder="nai-diffusion-4-5-full" class="setting-input" />
</div>
</template>
<template v-else-if="imgGenForm.apiFormat === 'gemini'">
<div class="input-group">
<label>Gemini 地址</label>
<input v-model="imgGenForm.geminiUrl" placeholder="https://generativelanguage.googleapis.com" class="setting-input" />
</div>
<div class="input-group">
<label>Gemini Key</label>
<input v-model="imgGenForm.geminiKey" type="password" placeholder="AIza..." class="setting-input" />
</div>
<div class="input-group">
<label>模型</label>
<input v-model="imgGenForm.geminiModel" placeholder="gemini-2.0-flash-exp" class="setting-input" />
</div>
</template>
</div>
</div>

<button class="save-btn" @click="saveApiSettings" :disabled="apiSaving">
              <template v-if="apiSaving">
                <div class="btn-spinner"></div>
                保存中...
              </template>
              <template v-else>
                ▣ 保存API设置
              </template>
</button>
</template>
</div>
</template>

<!-- OAuth配置 -->
<div v-show="activeTab === 'oauth'" class="settings-section">
  <div class="section-header">
    <h3>Discord OAuth设置</h3>
  </div>
  <div class="settings-list">
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">Client ID</span>
        <span class="label-desc">Discord应用程序的Client ID</span>
      </div>
      <input v-model="oauthSettings.discord_client_id" class="setting-input full-width" :disabled="!authStore.isSuperAdmin" placeholder="your-discord-client-id" />
    </div>
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">Client Secret</span>
        <span class="label-desc">Discord应用程序的Client Secret</span>
      </div>
      <input v-model="oauthSettings.discord_client_secret" type="password" class="setting-input full-width" :disabled="!authStore.isSuperAdmin" placeholder="your-discord-client-secret" />
    </div>
    <div class="setting-item vertical">
      <div class="setting-label">
        <span class="label-text">Redirect URI</span>
        <span class="label-desc">OAuth回调地址，需在Discord开发者门户配置</span>
      </div>
      <input v-model="oauthSettings.discord_redirect_uri" class="setting-input full-width" :disabled="!authStore.isSuperAdmin" placeholder="https://yourdomain.com/api/auth/discord/callback" />
    </div>
  </div>
</div>
<button v-show="activeTab === 'oauth'" class="save-btn" @click="saveOauthSettings" :disabled="saving || !authStore.isSuperAdmin">
  <template v-if="saving">
    <div class="btn-spinner"></div>
    保存中...
  </template>
  <template v-else>
    ▣ 保存OAuth设置
  </template>
</button>
</div>

<!-- 联系方式 -->
  <div v-show="activeTab === 'contact'" class="settings-section">
    <div class="section-header">
      <h3>官方群聊</h3>
    </div>
    <div class="settings-list">
      <div class="contact-qrcode">
        <img :src="'https://raw.githubusercontent.com/Xeltra233/fzsmphone/master/public/qun_qrcode.jpg'" alt="官方群二维码" class="qrcode-image" @error="qrcodeError = true" />
        <div v-if="qrcodeError" class="qrcode-placeholder">群二维码加载失败，请联系管理员</div>
        <div class="contact-password">
          <span class="password-label">进群密码：</span>
          <span class="password-value">贩子死妈</span>
        </div>
      </div>
    </div>
  </div>

<!-- 提示信息 -->
<div class="toast" v-if="toast.show" :class="toast.type">
      {{ toast.message }}
    </div>
  </div>

  <Teleport to="body">
    <div v-if="showModelPullModal" class="modal-overlay" @click.self="closeModelPullModal">
      <div class="model-pull-modal">
        <div class="model-pull-header">
          <div>
            <h3>选择要添加的模型</h3>
            <p>仅在你点击“添加所选模型”后，模型才会写入列表。</p>
          </div>
          <button class="model-pull-close" type="button" @click="closeModelPullModal">×</button>
        </div>

        <div class="model-pull-toolbar">
          <input v-model.trim="modelPullSearch" class="setting-input" placeholder="搜索显示名或模型 ID" />
          <div class="model-pull-tabs" v-if="pulledModels.length > 0">
            <button
              class="model-pull-tab-btn"
              :class="{ active: modelPullTab === 'new' }"
              type="button"
              @click="modelPullTab = 'new'"
            >
              新增模型 ({{ newPulledModels.length }})
            </button>
            <button
              class="model-pull-tab-btn"
              :class="{ active: modelPullTab === 'existing' }"
              type="button"
              @click="modelPullTab = 'existing'"
            >
              已有模型 ({{ existingPulledModels.length }})
            </button>
          </div>
        </div>

        <div v-if="pulledModels.length > 0" class="model-pull-stats">
          <span>共拉取 {{ modelPullStats.total }} 个</span>
          <span>新增 {{ modelPullStats.newCount }} 个</span>
          <span>已存在 {{ modelPullStats.existingCount }} 个</span>
          <span>当前显示 {{ filteredPulledModels.length }} 个</span>
          <span>当前已选 {{ selectedVisiblePulledModelCount }} 个</span>
        </div>

        <div class="model-pull-batch-actions" v-if="pulledModels.length > 0">
          <button class="batch-action-btn" type="button" @click="selectAllPulledModels">全选当前</button>
          <button class="batch-action-btn" type="button" @click="clearPulledModelSelection">清空选择</button>
        </div>

        <div v-if="filteredPulledModels.length === 0" class="empty-models">
          {{ pulledModels.length === 0 ? '未拉取到可添加模型' : '没有匹配的模型' }}
        </div>
        <div v-else class="model-pull-list">
          <label v-for="model in filteredPulledModels" :key="model.id" class="model-pull-item" :class="{ disabled: model.alreadyExists, selected: isPulledModelSelected(model.id) }">
            <input
              type="checkbox"
              v-model="selectedPulledModelIds"
              :value="model.id"
              :checked="isPulledModelSelected(model.id)"
              :disabled="model.alreadyExists"
            />
            <div class="model-pull-item-main">
              <span class="model-pull-display">{{ getPulledModelPreviewName(model) }}</span>
              <span class="model-pull-id">ID: {{ model.id }}</span>
              <span class="model-pull-meta">{{ model.alreadyExists ? getPulledModelExistsReason(model) : '可添加到可用模型列表' }}</span>
            </div>
          </label>
        </div>

        <div class="model-pull-actions">
          <button class="add-model-btn secondary" type="button" @click="closeModelPullModal">取消</button>
          <button class="add-model-btn" type="button" @click="confirmAddPulledModels" :disabled="selectedPulledModelIds.length === 0">添加所选模型</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import apiClient from '@/api/client'
import NavBar from '@/components/common/NavBar.vue'

const authStore = useAuthStore()

const loading = ref(true)
const saving = ref(false)
const qrcodeError = ref(false)
const faviconError = ref(false)
const uploadingQrcode = ref(false)
const uploadingFavicon = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const activeTab = ref<'credits' | 'system' | 'api' | 'oauth' | 'contact'>('credits')

const creditSettings = reactive({
  default_credits: 1000,
  signin_daily_credits: 10,
  signin_streak_bonus: 5,
  invite_reward_credits: 100,
  signin_enabled: true,
  invite_enabled: true,
})

const couponForm = reactive({
  code: '',
  credits: 100,
  maxUses: 1,
  expiresIn: 0,
  note: '',
  batchSize: 1,
})
const couponList = ref<Array<{
  code: string
  credits: number
  max_uses: number
  current_uses: number
  expires_at: string | null
  created_at: string
  note: string
  is_active: boolean
}>>([])
const latestCoupon = ref<null | { codes: string[]; credits: number; max_uses: number; note: string }>(null)
const couponCreating = ref(false)
const couponListLoading = ref(false)
const couponSearch = ref('')
const couponStatusFilter = ref<'all' | 'active' | 'inactive' | 'available' | 'exhausted' | 'expired'>('all')

const systemSettings = reactive({
  app_name: '',
  app_title: '',
  announcement: '',
  tips: '',
  disclaimer: '',
  qun_qrcode: '',
  favicon: '/icon.png',
})

const oauthSettings = reactive({
  discord_client_id: '',
  discord_client_secret: '',
  discord_redirect_uri: '',
})

const apiSettingsLoading = ref(false)
const apiSaving = ref(false)
const apiSettings = ref<Record<string, any> | null>(null)
type ManagedModel = { id: string; enabled: boolean; displayName: string }

const apiForm = ref({
globalApiKey: '',
globalApiUrl: '',
globalCustomUrl: '',
globalModel: '',
globalModels: [] as ManagedModel[],
globalTemperature: 0.9,
globalMaxLength: 4000,
globalContextSize: 20,
globalTimeout: 60,
globalSocialApiKey: '',
globalSocialApiUrl: '',
globalSocialModel: '',
})
const newModelId = ref('')
const pullingModels = ref(false)
const modelPullError = ref('')
const showModelPullModal = ref(false)
const pulledModels = ref<Array<{ id: string; displayName: string; alreadyExists: boolean; existsReason: string }>>([])
const selectedPulledModelIds = ref<string[]>([])
const modelPullSearch = ref('')
const modelPullTab = ref<'new' | 'existing'>('new')

const imgGenForm = ref({
apiFormat: 'openai',
novelaiUrl: '',
novelaiKey: '',
novelaiModel: 'nai-diffusion-4-5-full',
openaiUrl: '',
openaiKey: '',
openaiModel: '',
geminiUrl: '',
geminiKey: '',
geminiModel: '',
})

function showToast(message: string, type: 'success' | 'error' = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

function normalizeModelId(value: string) {
  return value.trim()
}

function modelIdKey(value: string) {
  return normalizeModelId(value).toLowerCase()
}

function sanitizeDisplayName(value: unknown) {
  return typeof value === 'string' ? value.trim() : ''
}

function pickModelCollection(payload: any): any[] {
  if (Array.isArray(payload)) return payload
  if (!payload || typeof payload !== 'object') return []

  if (Array.isArray(payload.data)) return payload.data
  if (payload.data && typeof payload.data === 'object' && Array.isArray(payload.data.data)) return payload.data.data
  if (Array.isArray(payload.models)) return payload.models
  if (payload.models && typeof payload.models === 'object' && Array.isArray(payload.models.data)) return payload.models.data
  if (Array.isArray(payload.list)) return payload.list
  if (payload.list && typeof payload.list === 'object' && Array.isArray(payload.list.data)) return payload.list.data

  return []
}

function normalizeSearchText(value: string) {
  return value.trim().toLowerCase()
}

function isCouponExpired(value: string | null) {
  if (!value) return false
  const time = new Date(value).getTime()
  return !Number.isNaN(time) && time < Date.now()
}

function parseManagedModels(value: unknown): ManagedModel[] {
  if (Array.isArray(value)) {
    return value
      .map((item) => {
        if (typeof item === 'string') {
          const id = normalizeModelId(item)
          if (!id) return null
          return { id, enabled: true, displayName: '' }
        }
        if (!item || typeof item !== 'object') return null
        const model = item as Record<string, unknown>
        const id = normalizeModelId(String(model.id || ''))
        if (!id) return null
        return {
          id,
          enabled: model.enabled !== false,
          displayName: sanitizeDisplayName(model.display_name ?? model.displayName),
        }
      })
      .filter((item): item is ManagedModel => Boolean(item))
  }
  if (typeof value !== 'string') return []
  return value
    .split(',')
    .map(normalizeModelId)
    .filter(Boolean)
    .map((id) => ({ id, enabled: true, displayName: '' }))
}

function getEnabledModelIds() {
  return apiForm.value.globalModels
    .filter((model) => model.enabled)
    .map((model) => model.id)
}

const enabledModels = computed(() => apiForm.value.globalModels.filter((model) => model.enabled))
const couponDisplayRows = computed(() => {
  return couponList.value.map((coupon) => ({
    ...coupon,
    expired: isCouponExpired(coupon.expires_at),
    exhausted: coupon.current_uses >= coupon.max_uses,
  }))
})
const filteredCouponList = computed(() => {
  const keyword = couponSearch.value.trim().toLowerCase()
  return couponDisplayRows.value.filter((coupon) => {
    if (couponStatusFilter.value === 'active' && !coupon.is_active) return false
    if (couponStatusFilter.value === 'inactive' && coupon.is_active) return false
    if (couponStatusFilter.value === 'available' && (coupon.exhausted || coupon.expired)) return false
    if (couponStatusFilter.value === 'exhausted' && !coupon.exhausted) return false
    if (couponStatusFilter.value === 'expired' && !coupon.expired) return false
    if (!keyword) return true
    return coupon.code.toLowerCase().includes(keyword) || coupon.note.toLowerCase().includes(keyword)
  })
})
const managedModelDisplayMap = computed(() => {
  const map = new Map<string, string>()
  for (const model of apiForm.value.globalModels) {
    map.set(modelIdKey(model.id), model.displayName.trim() || model.id)
  }
  return map
})
const newPulledModels = computed(() => pulledModels.value.filter((model) => !model.alreadyExists))
const existingPulledModels = computed(() => pulledModels.value.filter((model) => model.alreadyExists))
const filteredPulledModels = computed(() => {
  const keyword = normalizeSearchText(modelPullSearch.value)
  const source = modelPullTab.value === 'existing' ? existingPulledModels.value : newPulledModels.value
  return source.filter((model) => {
    if (!keyword) return true
    const previewName = normalizeSearchText(getPulledModelPreviewName(model))
    return modelIdKey(model.id).includes(keyword) || previewName.includes(keyword)
  })
})
const modelPullStats = computed(() => {
  const total = pulledModels.value.length
  const existingCount = pulledModels.value.filter((model) => model.alreadyExists).length
  return {
    total,
    existingCount,
    newCount: total - existingCount,
  }
})

function ensureDefaultModelValid() {
  const enabledIds = getEnabledModelIds()
  if (!enabledIds.length) {
    apiForm.value.globalModel = ''
    apiForm.value.globalSocialModel = ''
    return
  }
  if (!enabledIds.includes(apiForm.value.globalModel)) {
    apiForm.value.globalModel = enabledIds[0]
  }
  if (apiForm.value.globalSocialModel && !enabledIds.includes(apiForm.value.globalSocialModel)) {
    apiForm.value.globalSocialModel = ''
  }
}

function addModel() {
  const id = normalizeModelId(newModelId.value)
  if (!id) return
  const exists = apiForm.value.globalModels.some((model) => modelIdKey(model.id) === modelIdKey(id))
  if (exists) {
    showToast('该模型已存在', 'error')
    return
  }
  apiForm.value.globalModels.push({ id, enabled: true, displayName: '' })
  newModelId.value = ''
  ensureDefaultModelValid()
}

function mergeModels(models: ManagedModel[]) {
  const merged = [...apiForm.value.globalModels]
  for (const model of models) {
    const index = merged.findIndex((item) => modelIdKey(item.id) === modelIdKey(model.id))
    if (index >= 0) {
      const current = merged[index]
      merged[index] = {
        ...current,
        displayName: current.displayName || model.displayName || '',
      }
      continue
    }
    merged.push(model)
  }
  apiForm.value.globalModels = merged
  ensureDefaultModelValid()
}

function resolveApiUrl() {
  if (apiForm.value.globalApiUrl === 'custom') {
    return apiForm.value.globalCustomUrl.trim()
  }
  return apiForm.value.globalApiUrl.trim()
}

function extractModelItems(payload: any) {
  const list = pickModelCollection(payload)

  return list
    .map((item) => {
      if (typeof item === 'string') {
        const id = normalizeModelId(item)
        return id ? { id, displayName: '' } : null
      }
      if (!item || typeof item !== 'object') return null
      const id = normalizeModelId(item.id || item.name || item.model || '')
      if (!id) return null
      const displayName = sanitizeDisplayName(
        item.display_name
        ?? item.displayName
        ?? item.name_zh
        ?? item.cn
        ?? item.label
        ?? item.title
        ?? item.meta?.display_name
        ?? item.meta?.displayName
        ?? item.meta?.name
      )
      return {
        id,
        displayName,
      }
    })
    .filter((item): item is { id: string; displayName: string } => Boolean(item))
    .filter((item, index, arr) => arr.findIndex((current) => modelIdKey(current.id) === modelIdKey(item.id)) === index)
}

function closeModelPullModal() {
  showModelPullModal.value = false
  pulledModels.value = []
  selectedPulledModelIds.value = []
  modelPullSearch.value = ''
  modelPullTab.value = 'new'
}

function isPulledModelSelected(id: string) {
  return selectedPulledModelIds.value.some((item) => modelIdKey(item) === modelIdKey(id))
}

function selectAllPulledModels() {
  selectedPulledModelIds.value = filteredPulledModels.value.filter((model) => !model.alreadyExists).map((model) => model.id)
}

function clearPulledModelSelection() {
  selectedPulledModelIds.value = []
}

const selectedVisiblePulledModelCount = computed(() => {
  return filteredPulledModels.value.filter((model) => isPulledModelSelected(model.id)).length
})

function confirmAddPulledModels() {
  const selected = pulledModels.value
    .filter((model) => !model.alreadyExists && isPulledModelSelected(model.id))
    .map((model) => ({ id: model.id, enabled: true, displayName: model.displayName }))

  if (!selected.length) {
    showToast('没有可添加的新模型', 'error')
    return
  }

  const before = apiForm.value.globalModels.length
  mergeModels(selected)
  const addedCount = apiForm.value.globalModels.length - before
  closeModelPullModal()
  showToast(addedCount > 0 ? `已添加 ${addedCount} 个模型` : '模型列表已是最新')
}

async function pullModels() {
  pullingModels.value = true
  modelPullError.value = ''
  try {
    const response = await apiClient.post<any>('/api/ai/models', {
      apiUrl: resolveApiUrl(),
      apiKey: apiForm.value.globalApiKey,
    })
    const pulled = extractModelItems(response).map((model) => ({ id: model.id, enabled: true, displayName: model.displayName }))
    if (!pulled.length) {
      throw new Error('接口已响应，但没有返回可用模型')
    }
    const existingModels = new Map(apiForm.value.globalModels.map((model) => [modelIdKey(model.id), model]))
    pulledModels.value = pulled.map((model) => ({
      id: model.id,
      displayName: model.displayName,
      alreadyExists: existingModels.has(modelIdKey(model.id)),
      existsReason: existingModels.get(modelIdKey(model.id))?.displayName?.trim() ? '已存在，当前已配置自定义显示名' : '已存在于可用模型列表',
    }))
    selectedPulledModelIds.value = []
    modelPullTab.value = newPulledModels.value.length > 0 ? 'new' : 'existing'
    showModelPullModal.value = true
  } catch (err: any) {
    modelPullError.value = err.message || '拉取模型失败'
  } finally {
    pullingModels.value = false
  }
}

function removeModel(id: string) {
  apiForm.value.globalModels = apiForm.value.globalModels.filter((model) => model.id !== id)
  ensureDefaultModelValid()
}

async function fetchCreditSettings() {
  loading.value = true
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    
    if (data.default_credits !== undefined) creditSettings.default_credits = Number(data.default_credits)
    if (data.signin_daily_credits !== undefined) creditSettings.signin_daily_credits = Number(data.signin_daily_credits)
    if (data.signin_streak_bonus !== undefined) creditSettings.signin_streak_bonus = Number(data.signin_streak_bonus)
    if (data.invite_reward_credits !== undefined) creditSettings.invite_reward_credits = Number(data.invite_reward_credits)
    if (data.signin_enabled !== undefined) creditSettings.signin_enabled = data.signin_enabled === true || data.signin_enabled === 'true'
    if (data.invite_enabled !== undefined) creditSettings.invite_enabled = data.invite_enabled === true || data.invite_enabled === 'true'
  } catch (err: any) {
    showToast('加载设置失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    loading.value = false
  }
}

async function fetchCoupons() {
  if (!authStore.isSuperAdmin) return
  couponListLoading.value = true
  try {
    const res = await apiClient.get<{ coupons?: any[] }>('/api/credits/coupons')
    couponList.value = Array.isArray(res.coupons) ? res.coupons : []
  } catch (err: any) {
    showToast('加载兑换码失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    couponListLoading.value = false
  }
}

async function handleCreateCoupon() {
  if (!authStore.isSuperAdmin || couponCreating.value) return
  if (couponForm.credits <= 0) {
    showToast('奖励额度必须大于 0', 'error')
    return
  }
  if (couponForm.maxUses <= 0) {
    showToast('可使用次数必须大于 0', 'error')
    return
  }
  if (couponForm.batchSize <= 0 || couponForm.batchSize > 100) {
    showToast('批量数量必须在 1 到 100 之间', 'error')
    return
  }
  if (couponForm.code.trim() && couponForm.batchSize > 1) {
    showToast('自定义兑换码时不能批量生成', 'error')
    return
  }

  couponCreating.value = true
  try {
    const res = await apiClient.post<{ code?: string; codes?: string[]; credits: number; max_uses: number; note?: string }>('/api/credits/coupon', {
      code: couponForm.code.trim(),
      credits: couponForm.credits,
      max_uses: couponForm.maxUses,
      expires_in: couponForm.expiresIn,
      note: couponForm.note.trim(),
      batch_size: couponForm.batchSize,
    })
    latestCoupon.value = {
      codes: Array.isArray(res.codes) && res.codes.length ? res.codes : (res.code ? [res.code] : []),
      credits: res.credits,
      max_uses: res.max_uses,
      note: res.note || couponForm.note.trim(),
    }
    couponForm.code = ''
    couponForm.note = ''
    couponForm.batchSize = 1
    if (latestCoupon.value.codes.length === 1) {
      await navigator.clipboard.writeText(latestCoupon.value.codes[0])
      showToast(`兑换码已生成并复制：${latestCoupon.value.codes[0]}`)
    } else if (latestCoupon.value.codes.length > 1) {
      await navigator.clipboard.writeText(latestCoupon.value.codes.join('\n'))
      showToast(`已批量生成 ${latestCoupon.value.codes.length} 个兑换码，并复制到剪贴板`)
    }
    await fetchCoupons()
  } catch (err: any) {
    showToast('生成兑换码失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    couponCreating.value = false
  }
}

function formatCouponExpiry(value: string | null) {
  if (!value) return '永久有效'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '过期时间未知'
  return `截止 ${date.toLocaleString()}`
}

function getPulledModelPreviewName(model: { id: string; displayName: string }) {
  return managedModelDisplayMap.value.get(modelIdKey(model.id)) || model.displayName || model.id
}

function getPulledModelExistsReason(model: { existsReason: string }) {
  return model.existsReason || '已存在，已自动跳过'
}

function escapeCsvCell(value: string | number | boolean) {
  const text = String(value ?? '')
  return /[",\n]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text
}

function downloadTextFile(filename: string, content: string, mimeType: string) {
  const blob = new Blob([content], { type: mimeType })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

function exportCoupons(format: 'txt' | 'csv') {
  if (!filteredCouponList.value.length) {
    showToast('当前没有可导出的兑换码', 'error')
    return
  }

  const rows = filteredCouponList.value.map((coupon) => ({
    code: coupon.code,
    credits: coupon.credits,
    currentUses: coupon.current_uses,
    maxUses: coupon.max_uses,
    expiresAt: coupon.expires_at ? formatCouponExpiry(coupon.expires_at) : '永久有效',
    status: coupon.is_active ? '启用中' : '已停用',
    note: coupon.note || '',
    createdAt: new Date(coupon.created_at).toLocaleString(),
  }))

  if (format === 'txt') {
    const content = rows
      .map((row, index) => `${index + 1}. ${row.code}\n额度: ${row.credits}\n已用: ${row.currentUses}/${row.maxUses}\n状态: ${row.status}\n过期: ${row.expiresAt}\n备注: ${row.note || '无'}\n创建时间: ${row.createdAt}`)
      .join('\n\n')
    downloadTextFile(`coupons-${Date.now()}.txt`, content, 'text/plain;charset=utf-8')
    showToast(`已导出 ${rows.length} 条兑换码 TXT`)
    return
  }

  const header = ['code', 'credits', 'current_uses', 'max_uses', 'status', 'expires_at', 'note', 'created_at']
  const content = [
    header.join(','),
    ...rows.map((row) => [row.code, row.credits, row.currentUses, row.maxUses, row.status, row.expiresAt, row.note, row.createdAt].map(escapeCsvCell).join(',')),
  ].join('\n')
  downloadTextFile(`coupons-${Date.now()}.csv`, content, 'text/csv;charset=utf-8')
  showToast(`已导出 ${rows.length} 条兑换码 CSV`)
}

async function copyCouponCode(code: string) {
  try {
    await navigator.clipboard.writeText(code)
    showToast('兑换码已复制')
  } catch {
    showToast('复制失败', 'error')
  }
}

async function toggleCouponStatus(coupon: { code: string; is_active: boolean }) {
  try {
    await apiClient.patch(`/api/credits/coupon/${encodeURIComponent(coupon.code)}`, {
      is_active: !coupon.is_active,
    })
    showToast(coupon.is_active ? '兑换码已停用' : '兑换码已启用')
    await fetchCoupons()
  } catch (err: any) {
    showToast('更新兑换码失败: ' + (err.message || '未知错误'), 'error')
  }
}

async function deleteCoupon(code: string) {
  const confirmed = confirm(`确定删除兑换码 ${code} 吗？删除后无法恢复。`)
  if (!confirmed) return
  try {
    await apiClient.delete(`/api/credits/coupon/${encodeURIComponent(code)}`)
    showToast('兑换码已删除')
    await fetchCoupons()
  } catch (err: any) {
    showToast('删除兑换码失败: ' + (err.message || '未知错误'), 'error')
  }
}

async function fetchSystemSettings() {
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    systemSettings.app_name = data.app_name || ''
    systemSettings.app_title = data.app_title || ''
    systemSettings.announcement = data.announcement || ''
    systemSettings.tips = data.tips || ''
    systemSettings.disclaimer = data.disclaimer || ''
    systemSettings.qun_qrcode = data.qun_qrcode || ''
    systemSettings.favicon = data.favicon || '/icon.png'
  } catch (err: any) {
    console.error('Failed to load system settings:', err)
  }
}

async function saveCreditSettings() {
  saving.value = true
  try {
    await apiClient.put('/api/settings', {
      default_credits: creditSettings.default_credits,
      signin_daily_credits: creditSettings.signin_daily_credits,
      signin_streak_bonus: creditSettings.signin_streak_bonus,
      invite_reward_credits: creditSettings.invite_reward_credits,
      signin_enabled: creditSettings.signin_enabled,
      invite_enabled: creditSettings.invite_enabled,
    })
    showToast('额度设置已保存')
  } catch (err: any) {
    showToast('保存失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    saving.value = false
  }
}

async function handleQrcodeUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  uploadingQrcode.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res: any = await apiClient.upload('/api/settings/upload', formData)
    systemSettings.qun_qrcode = res.path
    showToast('二维码上传成功')
  } catch (err: any) {
    showToast('上传失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    uploadingQrcode.value = false
  }
}

async function handleFaviconUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  uploadingFavicon.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res: any = await apiClient.upload('/api/settings/upload', formData)
    systemSettings.favicon = res.path
    faviconError.value = false
    showToast('图标上传成功')
  } catch (err: any) {
    showToast('上传失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    uploadingFavicon.value = false
  }
}

async function saveSystemSettings() {
  saving.value = true
  try {
    await apiClient.put('/api/settings', {
      app_name: systemSettings.app_name,
      app_title: systemSettings.app_title,
      announcement: systemSettings.announcement,
      tips: systemSettings.tips,
      disclaimer: systemSettings.disclaimer,
      qun_qrcode: systemSettings.qun_qrcode,
      favicon: systemSettings.favicon,
    })
    showToast('系统设置已保存')
  } catch (err: any) {
    showToast('保存失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    saving.value = false
  }
}

async function fetchOauthSettings() {
  try {
    const res: any = await apiClient.get('/api/settings')
    const data = res.data || res || {}
    oauthSettings.discord_client_id = data.discord_client_id || ''
    oauthSettings.discord_client_secret = data.discord_client_secret || ''
    oauthSettings.discord_redirect_uri = data.discord_redirect_uri || ''
  } catch (err: any) {
    console.error('Failed to load OAuth settings:', err)
  }
}

async function saveOauthSettings() {
  saving.value = true
  try {
    await apiClient.put('/api/settings', {
      discord_client_id: oauthSettings.discord_client_id,
      discord_client_secret: oauthSettings.discord_client_secret,
      discord_redirect_uri: oauthSettings.discord_redirect_uri,
    })
    showToast('OAuth设置已保存')
  } catch (err: any) {
    showToast('保存失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    saving.value = false
  }
}

async function fetchApiSettings() {
  apiSettingsLoading.value = true
  try {
    const res = await apiClient.get<Record<string, any>>('/api/settings/api') as Record<string, any>
    apiSettings.value = res
    if (authStore.isSuperAdmin && res.global) {
      apiForm.value.globalApiKey = res.global.api_key || ''
      apiForm.value.globalApiUrl = res.global.api_url || ''
      apiForm.value.globalCustomUrl = res.global.custom_url || ''
      apiForm.value.globalModel = res.global.model || ''
      apiForm.value.globalModels = parseManagedModels(res.global.model_list || [])
      apiForm.value.globalTemperature = res.global.temperature ?? 0.9
      apiForm.value.globalMaxLength = res.global.max_length ?? 4000
      apiForm.value.globalContextSize = res.global.context_size ?? 20
      apiForm.value.globalTimeout = res.global.timeout ?? 60
      apiForm.value.globalSocialApiKey = res.global.social_api_key || ''
      apiForm.value.globalSocialApiUrl = res.global.social_api_url || ''
      apiForm.value.globalSocialModel = res.global.social_model || ''
      ensureDefaultModelValid()
    }

    try {
      const settingsData = await apiClient.get<Record<string, any>>('/api/settings')
      if (settingsData.img_gen_config) {
        const imgCfg = JSON.parse(settingsData.img_gen_config)
        imgGenForm.value.apiFormat = imgCfg.apiFormat || 'openai'
        if (imgCfg.novelai) {
          imgGenForm.value.novelaiUrl = imgCfg.novelai.url || ''
          imgGenForm.value.novelaiKey = imgCfg.novelai.key || ''
          imgGenForm.value.novelaiModel = imgCfg.novelai.model || 'nai-diffusion-4-5-full'
        }
        if (imgCfg.openai) {
          imgGenForm.value.openaiUrl = imgCfg.openai.url || ''
          imgGenForm.value.openaiKey = imgCfg.openai.key || ''
          imgGenForm.value.openaiModel = imgCfg.openai.model || ''
        }
        if (imgCfg.gemini) {
          imgGenForm.value.geminiUrl = imgCfg.gemini.url || ''
          imgGenForm.value.geminiKey = imgCfg.gemini.key || ''
          imgGenForm.value.geminiModel = imgCfg.gemini.model || ''
        }
      }
    } catch (e) {
      console.error('Failed to load img gen config', e)
    }
  } catch (err: any) {
    showToast('加载API设置失败', 'error')
  } finally {
    apiSettingsLoading.value = false
  }
}

async function saveApiSettings() {
  apiSaving.value = true
  try {
    if (authStore.isSuperAdmin) {
      await apiClient.put('/api/settings/api', {
        api_key: apiForm.value.globalApiKey,
        api_url: apiForm.value.globalApiUrl,
        custom_url: apiForm.value.globalCustomUrl,
        model: apiForm.value.globalModel,
        model_list: apiForm.value.globalModels.map((model) => ({
          id: model.id,
          enabled: model.enabled,
          display_name: model.displayName.trim(),
        })),
        temperature: apiForm.value.globalTemperature,
        max_length: apiForm.value.globalMaxLength,
        context_size: apiForm.value.globalContextSize,
        timeout: apiForm.value.globalTimeout,
        social_api_key: apiForm.value.globalSocialApiKey,
        social_api_url: apiForm.value.globalSocialApiUrl,
        social_model: apiForm.value.globalSocialModel,
        is_global: true,
      })

      await apiClient.put('/api/settings', {
        img_gen_config: JSON.stringify({
          apiFormat: imgGenForm.value.apiFormat,
          novelai: { url: imgGenForm.value.novelaiUrl, key: imgGenForm.value.novelaiKey, model: imgGenForm.value.novelaiModel },
          openai: { url: imgGenForm.value.openaiUrl, key: imgGenForm.value.openaiKey, model: imgGenForm.value.openaiModel },
          gemini: { url: imgGenForm.value.geminiUrl, key: imgGenForm.value.geminiKey, model: imgGenForm.value.geminiModel },
        }),
      })
    }
    showToast('API设置已保存')
  } catch (err: any) {
    showToast('保存API设置失败: ' + (err.message || '未知错误'), 'error')
  } finally {
    apiSaving.value = false
  }
}

onMounted(() => {
  if (authStore.isAdmin) {
    fetchCreditSettings()
    fetchCoupons()
    fetchSystemSettings()
    fetchApiSettings()
    fetchOauthSettings()
  }
})
</script>

<style scoped>
.system-settings-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary, #f5f5f7);
  color: var(--text-primary, #1d1d1f);
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  position: relative;
}

.page-content::-webkit-scrollbar {
  display: none;
}

/* 加载 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  gap: 16px;
  color: var(--text-tertiary, #86868b);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #5B6EF5;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Section */
.settings-section {
  margin-bottom: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-tertiary, 0.5);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0;
}

/* 设置项 */
.settings-list {
  background: var(--bg-tertiary, 0.06);
  border: 1px solid var(--separator, 0.08);
  border-radius: 12px;
  overflow: hidden;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--separator, 0.06);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label-text {
  font-size: 15px;
  color: var(--text-primary);
}

.label-desc {
  font-size: 12px;
  color: var(--text-tertiary, 0.4);
}

.setting-input-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.setting-input {
  width: 100%;
  min-width: 0;
  padding: 8px 12px;
  background: var(--bg-secondary, 0.3);
  border: 1px solid var(--separator, 0.15);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  text-align: left;
  outline: none;
  box-sizing: border-box;
}

.setting-input:focus {
  border-color: #5B6EF5;
}

.setting-input.full-width {
  width: 100%;
  text-align: left;
}

.setting-item.vertical {
  flex-direction: column;
  align-items: stretch;
  gap: 8px;
}

.setting-item.vertical .setting-label {
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.setting-item.vertical .setting-input {
  width: 100%;
  box-sizing: border-box;
}

.setting-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.input-suffix {
  font-size: 13px;
  color: var(--text-tertiary, 0.5);
}

/* Toggle */
.toggle {
  position: relative;
  width: 51px;
  height: 31px;
  cursor: pointer;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
  position: absolute;
}

.toggle-slider {
  position: absolute;
  inset: 0;
  background: var(--bg-quaternary, rgba(255, 255, 255, 0.1));
  border-radius: 31px;
  transition: background 0.3s;
}

.toggle-slider::after {
  content: '';
  position: absolute;
  left: 2px;
  top: 2px;
  width: 27px;
  height: 27px;
  background: #fff;
  border-radius: 50%;
  transition: transform 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.toggle input:checked + .toggle-slider {
  background: #34C759;
}

.toggle input:checked + .toggle-slider::after {
  transform: translateX(20px);
}

/* 操作按钮 */
.save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  width: 100%;
  background: linear-gradient(135deg, #5B6EF5, #8B5CF6);
  border: none;
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 20px;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.save-btn:not(:disabled):hover {
  filter: brightness(1.1);
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: var(--text-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Toast */
.toast {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  z-index: 1000;
  animation: fadeInUp 0.3s ease;
}

.toast.success {
  background: rgba(52, 199, 89, 0.9);
  color: var(--text-primary);
}

.toast.error {
  background: rgba(255, 59, 48, 0.9);
  color: var(--text-primary);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateX(-50%) translateY(10px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}

/* 无权限 */
.no-access {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary, 0.4);
  font-size: 16px;
}

/* 标签切换 */
.settings-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 8px 16px;
  background: var(--bg-tertiary, 0.08);
  border: 1px solid var(--separator, 0.15);
  border-radius: 20px;
  color: var(--text-tertiary, 0.6);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn.active {
  background: rgba(91, 110, 245, 0.3);
  border-color: #5B6EF5;
  color: var(--text-primary);
}

/* 权限提示 */
.permission-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(91, 110, 245, 0.15);
  border: 1px solid rgba(91, 110, 245, 0.3);
  border-radius: 12px;
  margin-bottom: 16px;
  font-size: 13px;
  color: var(--text-tertiary, 0.7);
}

.notice-icon {
  font-size: 16px;
}

/* API设置区域 */
.api-settings-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.settings-card {
  background: var(--bg-tertiary, 0.06);
  border: 1px solid var(--separator, 0.08);
  border-radius: 12px;
  overflow: hidden;
}

.coupon-manager {
  gap: 14px;
}

.coupon-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.coupon-actions-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.coupon-toolbar {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.coupon-toolbar .setting-input {
  flex: 1;
  min-width: 180px;
}

.coupon-filter-select {
  max-width: 180px;
}

.coupon-toolbar-meta {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  font-size: 12px;
  color: var(--text-secondary);
}

.coupon-export-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.coupon-latest-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 14px;
  border-radius: 12px;
  background: rgba(52, 199, 89, 0.1);
  border: 1px solid rgba(52, 199, 89, 0.18);
}

.coupon-latest-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.coupon-latest-code {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--text-primary);
}

.coupon-latest-multi {
  font-size: 16px;
  letter-spacing: normal;
}

.coupon-latest-meta {
  font-size: 12px;
  color: var(--text-secondary);
}

.coupon-latest-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.coupon-latest-chip {
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(52, 199, 89, 0.12);
  color: var(--text-primary);
  font-size: 12px;
  font-weight: 600;
}

.coupon-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.coupon-list-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid var(--separator, 0.08);
  background: var(--bg-primary, rgba(255, 255, 255, 0.6));
}

.coupon-list-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.coupon-list-code {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  word-break: break-all;
}

.coupon-list-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 12px;
  font-size: 12px;
  color: var(--text-secondary);
}

.coupon-item-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-end;
}

.batch-action-btn.danger {
  background: rgba(255, 59, 48, 0.12);
  color: #ff3b30;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid var(--separator, 0.08);
}

.card-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.badge {
  padding: 4px 10px;
  background: rgba(91, 110, 245, 0.2);
  border-radius: 12px;
  font-size: 11px;
  color: #5B6EF5;
}

.card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.input-group label {
  font-size: 13px;
  color: var(--text-tertiary, 0.6);
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.input-tip {
  font-size: 11px;
  color: rgba(91, 110, 245, 0.8);
  margin-left: 6px;
}

.input-desc {
  font-size: 11px;
  color: var(--text-tertiary, 0.4);
  margin-top: 4px;
  display: block;
}

.input-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.input-row .input-group {
  flex: 1;
  min-width: 140px;
}

.api-settings-section .setting-input,
.api-settings-section .setting-select {
  width: 100%;
  min-height: 42px;
}

.api-settings-section .card-body {
  gap: 14px;
}

.model-manager-group {
  gap: 10px;
}

.model-add-row {
  display: flex;
  gap: 10px;
  align-items: stretch;
}

.add-model-btn {
  flex: 0 0 auto;
  min-width: 72px;
  padding: 0 14px;
  border: 1px solid rgba(91, 110, 245, 0.35);
  border-radius: 10px;
  background: rgba(91, 110, 245, 0.14);
  color: #5B6EF5;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.add-model-btn.secondary {
  border-color: var(--separator, rgba(120, 120, 128, 0.18));
  background: var(--bg-secondary, rgba(120, 120, 128, 0.08));
  color: var(--text-primary);
}

.add-model-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.model-pull-error {
  font-size: 12px;
  color: #ff3b30;
}

.model-pull-modal {
  width: min(560px, calc(100vw - 32px));
  max-height: min(78vh, 720px);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border-radius: 18px;
  background: var(--bg-primary, #fff);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
}

.model-pull-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 18px 18px 12px;
  border-bottom: 1px solid var(--separator, rgba(120, 120, 128, 0.14));
}

.model-pull-header h3 {
  margin: 0;
  font-size: 17px;
  color: var(--text-primary);
}

.model-pull-header p {
  margin: 6px 0 0;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.model-pull-close {
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
}

.model-pull-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px 18px;
  overflow-y: auto;
}

.model-pull-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: flex-start;
  padding: 14px 18px 0;
}

.model-pull-toolbar .setting-input {
  flex: 1;
  min-width: 220px;
}

.model-pull-tabs {
  display: inline-flex;
  gap: 8px;
  flex-wrap: wrap;
}

.model-pull-tab-btn {
  border: 1px solid var(--separator, rgba(120, 120, 128, 0.14));
  background: var(--bg-secondary, rgba(120, 120, 128, 0.06));
  color: var(--text-secondary);
  border-radius: 999px;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.model-pull-tab-btn.active {
  color: #5B6EF5;
  border-color: rgba(91, 110, 245, 0.32);
  background: rgba(91, 110, 245, 0.1);
}

.model-pull-batch-actions {
  display: flex;
  gap: 10px;
  padding: 12px 18px 0;
}

.model-pull-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 12px;
  padding: 12px 18px 0;
  font-size: 12px;
  color: var(--text-secondary);
}

.batch-action-btn {
  border: none;
  background: rgba(91, 110, 245, 0.1);
  color: #5B6EF5;
  border-radius: 999px;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.model-pull-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid var(--separator, rgba(120, 120, 128, 0.14));
  background: var(--bg-secondary, rgba(120, 120, 128, 0.06));
  cursor: pointer;
}

.model-pull-item.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.model-pull-item.selected {
  border-color: rgba(91, 110, 245, 0.48);
  background: rgba(91, 110, 245, 0.12);
  box-shadow: 0 0 0 1px rgba(91, 110, 245, 0.12);
}

.model-pull-item input {
  margin-top: 2px;
}

.model-pull-item-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.model-pull-display {
  color: var(--text-primary);
  font-size: 15px;
  font-weight: 600;
  word-break: break-word;
}

.model-pull-id {
  color: var(--text-secondary);
  font-size: 12px;
  word-break: break-all;
}

.model-pull-meta {
  color: var(--text-secondary);
  font-size: 12px;
}

.model-pull-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  padding: 14px 18px 18px;
  border-top: 1px solid var(--separator, rgba(120, 120, 128, 0.14));
}

.model-list-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.model-chip-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid var(--separator, 0.12);
  border-radius: 12px;
  background: var(--bg-secondary, 0.35);
}

.model-chip-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.model-chip-id {
  min-width: 0;
  color: var(--text-primary);
  font-size: 13px;
  word-break: break-all;
}

.model-display-input {
  min-height: 36px !important;
}

.model-delete-btn {
  flex: 0 0 auto;
  border: none;
  background: rgba(255, 59, 48, 0.14);
  color: #ff3b30;
  border-radius: 999px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
}

.empty-models {
  padding: 12px;
  border: 1px dashed var(--separator, 0.16);
  border-radius: 12px;
  text-align: center;
  font-size: 13px;
  color: var(--text-tertiary, 0.5);
}

.model-toggle {
  position: relative;
  width: 42px;
  height: 24px;
  flex: 0 0 auto;
}

.model-toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.model-toggle-slider {
  position: absolute;
  inset: 0;
  border-radius: 999px;
  background: var(--fill-secondary, rgba(120, 120, 128, 0.2));
  transition: background 0.2s ease;
}

.model-toggle-slider::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.2s ease;
}

.model-toggle input:checked + .model-toggle-slider {
  background: #34c759;
}

.model-toggle input:checked + .model-toggle-slider::after {
  transform: translateX(18px);
}

@media (max-width: 430px) {
  .input-row {
    flex-direction: column;
    gap: 10px;
  }

  .coupon-form-grid {
    grid-template-columns: 1fr;
  }

  .coupon-list-item {
    flex-direction: column;
    align-items: stretch;
  }

  .coupon-filter-select {
    max-width: none;
  }

  .input-row .input-group {
    min-width: 0;
  }

  .model-add-row,
  .model-chip-row {
    flex-wrap: wrap;
  }

  .add-model-btn,
  .model-delete-btn {
    width: 100%;
  }
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: var(--text-primary);
  font-size: 14px;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  appearance: none;
  -webkit-appearance: none;
  border: 2px solid var(--text-tertiary, #86868b);
  border-radius: 4px;
  background: var(--bg-tertiary, #e5e5ea);
  cursor: pointer;
  transition: all 0.2s ease;
}

.checkbox-label input[type="checkbox"]:checked {
  background: #34c759;
  border-color: #34c759;
}

.checkbox-label input[type="checkbox"]:checked::after {
  content: '✓';
  display: block;
  text-align: center;
  color: white;
  font-size: 12px;
  line-height: 14px;
}

.qrcode-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  width: 100%;
}

.contact-qrcode {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  gap: 12px;
}

.contact-password {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
}

.password-label {
  color: var(--text-tertiary, 0.6);
}

.password-value {
  color: #5B6EF5;
  font-weight: 600;
}

.qrcode-image {
  max-width: 200px;
  max-height: 200px;
  border-radius: 12px;
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.qrcode-placeholder {
  color: var(--text-tertiary, 0.4);
  font-size: 14px;
  padding: 20px;
}

.password-display {
  font-size: 18px;
  font-weight: 600;
  color: #5B6EF5;
  padding: 12px 20px;
  background: rgba(91, 110, 245, 0.1);
  border-radius: 8px;
  letter-spacing: 2px;
}

.file-upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 8px 0;
}

.upload-btn {
  display: inline-block;
  padding: 8px 16px;
  background: rgba(91, 110, 245, 0.2);
  border: 1px solid rgba(91, 110, 245, 0.4);
  border-radius: 6px;
  color: #5B6EF5;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-btn:hover {
  background: rgba(91, 110, 245, 0.3);
  border-color: #5B6EF5;
}

.upload-status {
  font-size: 12px;
  color: var(--text-tertiary, 0.5);
}

.favicon-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.favicon-img {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  object-fit: contain;
  background: var(--bg-tertiary, 0.1);
}

.favicon-placeholder {
  font-size: 13px;
  color: var(--text-tertiary, 0.4);
}
</style>
