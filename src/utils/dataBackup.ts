import api from '@/api/client'
import {
  callApi,
  characterApi,
  diaryApi,
  gameApi,
  personaApi,
  presetApi,
  smsApi,
  walletApi,
} from '@/api/services'
import { loadImageGenConfig } from '@/utils/imageGenConfig'
import type {
  Character,
  CharacterInput,
  Diary,
  DiaryInput,
  GameRecordItem,
  Persona,
  PersonaInput,
  Preset,
  PresetInput,
  SmsMessage,
  SmsThread,
  WalletTransaction,
  WalletInfo,
} from '@/api/types'

const BACKUP_VERSION = 1
const SETTINGS_KEY = 'fzsm-settings'
const PRESETS_KEY = 'aiPresets'
const ACTIVE_PRESET_KEY = 'activePresetId'
const WORLDBOOKS_KEY = 'worldBooks'
const CURRENT_PERSONA_KEY = 'currentPersonaId'
const CURRENT_CHARACTER_KEY = 'currentUserCharId'
const CHAT_CONVERSATIONS_KEY = 'chat-conversations'
const SMS_STORAGE_KEY = 'fzsm-sms-conversations'
const CALL_HISTORY_KEY = 'fzsm-call-history'
const WALLET_KEY = 'wallet-data'
const TX_KEY = 'wallet-transactions'
const REDPACKET_KEY = 'wallet-redpackets'
const TRANSFER_KEY = 'wallet-transfers'
const DIARY_LOCAL_KEY = 'diary-local-data'
const SECURITY_LOG_KEY = 'security_logs'
const GLOBAL_EJS_VARS_KEY = 'ejs-vars-global'
const CHAT_AUTO_IMAGE_KEY = 'phone_auto_gen_chat_images'
const SOCIAL_AUTO_IMAGE_KEY = 'phone_auto_gen_social_images'
const IMAGE_GEN_KEY = 'phone_novelai_config'

const SOCIAL_STORAGE_KEYS = {
  forum: 'social-data-forum',
  weibo: 'social-data-weibo',
  moments: 'social-data-moments',
  zhihu: 'social-data-zhihu',
  xiaohongshu: 'social-data-xiaohongshu',
  douyin: 'social-data-douyin',
  music: 'social-data-music',
  live: 'social-data-live',
  theater: 'social-data-theater',
  takeaway: 'social-data-takeaway',
  shopping: 'social-data-shopping',
  couple: 'social-data-couple',
  stock: 'social-data-stock',
  email: 'social-data-email',
  browser: 'social-data-browser',
  map: 'social-data-map',
  calendar: 'social-data-calendar',
} as const

type JsonObject = Record<string, any>

async function safeCall<T>(action: () => Promise<T>, fallback: T): Promise<T> {
  try {
    return await action()
  } catch {
    return fallback
  }
}

export interface BackupFile {
  format: 'fzsmphone-backup'
  version: number
  exported_at: string
  app_version: string
  includes: string[]
  data: {
    settings: JsonObject
    image_gen_config: JsonObject
    local_state: JsonObject
    social_data: JsonObject
    characters: CharacterInput[]
    presets: PresetInput[]
    personas: PersonaInput[]
    diaries: DiaryInput[]
    sms: Array<{ thread: Omit<SmsThread, 'id' | 'user_id' | 'created_at'>; messages: Array<Pick<SmsMessage, 'role' | 'content'>> }>
    calls: Array<{ name: string; number: string; avatar?: string; character_id: string; type: string; call_type: string; duration: string }>
    wallet: {
      balance: number | null
      transactions: Array<{ type: string; amount: number; description: string; target?: string }>
      redpackets: any[]
      transfers: any[]
    }
    games: Array<{ game: string; detail: string; amount: number; win: boolean }>
  }
}

function safeParse<T>(raw: string | null, fallback: T): T {
  if (!raw) return fallback
  try {
    return JSON.parse(raw) as T
  } catch {
    return fallback
  }
}

function clone<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

function sanitizeSettings(settings: JsonObject): JsonObject {
  const next = clone(settings || {})
  delete next.apiKey
  delete next.socialApiKey
  return next
}

function sanitizeImageGenConfig(config: JsonObject): JsonObject {
  const next = clone(config || {})
  if (next.novelai) delete next.novelai.key
  if (next.openai) delete next.openai.key
  if (next.gemini) delete next.gemini.key
  return next
}

function getLocalEjsVars(): JsonObject {
  const data: JsonObject = {}
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i)
    if (!key || !key.startsWith('ejs-vars-chat-')) continue
    data[key] = safeParse(localStorage.getItem(key), {})
  }
  return data
}

function getChatMessageBuckets(): JsonObject {
  const data: JsonObject = {}
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i)
    if (!key || !key.startsWith('chat-messages-')) continue
    data[key] = safeParse(localStorage.getItem(key), [])
  }
  return data
}

function getSocialData(): JsonObject {
  const data: JsonObject = {}
  for (const [name, key] of Object.entries(SOCIAL_STORAGE_KEYS)) {
    data[name] = safeParse(localStorage.getItem(key), null)
  }
  return data
}

function normalizePresetInput(preset: Preset): PresetInput {
  return {
    name: preset.name,
    emoji: preset.emoji,
    category: preset.category,
    description: preset.description,
    content: preset.content,
    prefill: preset.prefill,
    enable_prefill: preset.enable_prefill ?? (preset as any).enablePrefill ?? false,
    gradient: preset.gradient,
    is_builtin: preset.is_builtin ?? (preset as any).isBuiltin ?? false,
    prompt_items: preset.prompt_items || preset.promptItems || [],
  }
}

function normalizeCharacterInput(character: Character): CharacterInput {
  return {
    name: character.name,
    avatar_url: character.avatar_url,
    description: character.description,
    personality: character.personality,
    system_prompt: character.system_prompt,
    greeting: character.greeting,
    is_public: character.is_public,
    tags: Array.isArray(character.tags) ? character.tags : [],
    extra: character.extra || {},
  }
}

function normalizePersonaInput(persona: Persona): PersonaInput {
  return {
    name: persona.name,
    description: persona.description,
    avatar_url: persona.avatar_url,
    is_default: persona.is_default,
  }
}

function normalizeDiaryInput(diary: Diary): DiaryInput {
  return {
    title: diary.title,
    content: diary.content,
    mood: diary.mood,
    weather: diary.weather,
    tags: Array.isArray(diary.tags) ? diary.tags : [],
    is_private: diary.is_private,
  }
}

export async function exportAppBackup(): Promise<BackupFile> {
  const [charactersRes, presetsRes, personasRes, diariesRes, smsThreadsRes, walletRes, walletTxRes, gamesRes, callsRes] = await Promise.all([
    safeCall(() => characterApi.listWithStorage(), { data: [] as Character[] }),
    safeCall(() => presetApi.list(), { data: [] as Preset[] }),
    safeCall(() => personaApi.list(), { data: [] as Persona[] }),
    safeCall(() => diaryApi.list(), { data: [] as Diary[] }),
    safeCall(() => smsApi.listThreads(), { data: [] as SmsThread[] }),
    safeCall<WalletInfo | { balance: null }>(() => walletApi.get(), { balance: null }),
    safeCall(() => walletApi.listTransactions(), { data: [] as WalletTransaction[] }),
    safeCall(() => gameApi.listRecords(), { data: [] as GameRecordItem[] }),
    safeCall(() => callApi.list(), { data: [] as any[] }),
  ])

  const smsThreads = Array.isArray(smsThreadsRes.data) ? smsThreadsRes.data : []
  const sms = await Promise.all(smsThreads.map(async (thread) => {
    const messagesRes = await safeCall(() => smsApi.listMessages(thread.id), { data: [] as SmsMessage[] })
    return {
      thread: {
        recipient: thread.recipient,
        number: thread.number,
        character_id: thread.character_id,
        avatar: thread.avatar,
        last_content: thread.last_content,
        last_at: thread.last_at,
      },
      messages: (messagesRes.data || []).map((message) => ({
        role: message.role,
        content: message.content,
      })),
    }
  }))

  return {
    format: 'fzsmphone-backup',
    version: BACKUP_VERSION,
    exported_at: new Date().toISOString(),
    app_version: import.meta.env.VITE_APP_VERSION || 'unknown',
    includes: [
      'settings',
      'characters',
      'presets',
      'personas',
      'diaries',
      'sms',
      'wallet',
      'games',
      'calls',
      'chat',
      'social-data',
      'worldbooks',
      'template-vars',
    ],
    data: {
      settings: sanitizeSettings(safeParse(localStorage.getItem(SETTINGS_KEY), {})),
      image_gen_config: sanitizeImageGenConfig(loadImageGenConfig() as unknown as JsonObject),
      local_state: {
        presets_cache: safeParse(localStorage.getItem(PRESETS_KEY), []),
        active_preset_id: localStorage.getItem(ACTIVE_PRESET_KEY),
        worldbooks: safeParse(localStorage.getItem(WORLDBOOKS_KEY), []),
        current_persona_id: localStorage.getItem(CURRENT_PERSONA_KEY),
        current_character_id: localStorage.getItem(CURRENT_CHARACTER_KEY),
        chat_conversations: safeParse(localStorage.getItem(CHAT_CONVERSATIONS_KEY), []),
        chat_messages: getChatMessageBuckets(),
        sms_local_cache: safeParse(localStorage.getItem(SMS_STORAGE_KEY), []),
        call_history_local: safeParse(localStorage.getItem(CALL_HISTORY_KEY), []),
        wallet_local: safeParse(localStorage.getItem(WALLET_KEY), {}),
        wallet_transactions_local: safeParse(localStorage.getItem(TX_KEY), []),
        wallet_redpackets_local: safeParse(localStorage.getItem(REDPACKET_KEY), []),
        wallet_transfers_local: safeParse(localStorage.getItem(TRANSFER_KEY), []),
        diary_local: safeParse(localStorage.getItem(DIARY_LOCAL_KEY), []),
        security_logs: safeParse(localStorage.getItem(SECURITY_LOG_KEY), []),
        chat_auto_image: localStorage.getItem(CHAT_AUTO_IMAGE_KEY),
        social_auto_image: localStorage.getItem(SOCIAL_AUTO_IMAGE_KEY),
        global_template_vars: safeParse(localStorage.getItem(GLOBAL_EJS_VARS_KEY), {}),
        local_template_vars: getLocalEjsVars(),
      },
      social_data: getSocialData(),
      characters: (charactersRes.data || []).map(normalizeCharacterInput),
      presets: (presetsRes.data || []).map(normalizePresetInput),
      personas: (personasRes.data || []).map(normalizePersonaInput),
      diaries: (diariesRes.data || []).map(normalizeDiaryInput),
      sms,
      calls: (callsRes.data || []).map((item: any) => ({
        name: item.name || '',
        number: item.number || '',
        avatar: item.avatar || '',
        character_id: item.character_id || '',
        type: item.type || 'outgoing',
        call_type: item.call_type || 'voice',
        duration: item.duration || '',
      })),
      wallet: {
        balance: typeof walletRes.balance === 'number' ? walletRes.balance : null,
        transactions: (walletTxRes.data || []).map((tx) => ({
          type: tx.type,
          amount: tx.amount,
          description: tx.description,
          target: tx.target,
        })),
        redpackets: safeParse(localStorage.getItem(REDPACKET_KEY), []),
        transfers: safeParse(localStorage.getItem(TRANSFER_KEY), []),
      },
      games: (gamesRes.data || []).map((record) => ({
        game: record.game,
        detail: record.detail,
        amount: record.amount,
        win: record.win,
      })),
    },
  }
}

function validateBackupFile(raw: any): BackupFile {
  if (!raw || typeof raw !== 'object' || raw.format !== 'fzsmphone-backup') {
    throw new Error('不是有效的 fzsmphone 备份文件')
  }
  if (!raw.data || typeof raw.data !== 'object') {
    throw new Error('备份文件缺少 data 字段')
  }
  return raw as BackupFile
}

function mergeObjectStorage(key: string, value: unknown) {
  localStorage.setItem(key, JSON.stringify(value ?? {}))
}

function mergeArrayStorage(key: string, value: unknown) {
  localStorage.setItem(key, JSON.stringify(Array.isArray(value) ? value : []))
}

function setScalarStorage(key: string, value: unknown) {
  if (value === null || value === undefined || value === '') {
    localStorage.removeItem(key)
    return
  }
  localStorage.setItem(key, String(value))
}

async function importCharacters(items: CharacterInput[]) {
  for (const item of items) {
    await safeCall(() => characterApi.create(item), null)
  }
}

async function importPresets(items: PresetInput[]) {
  for (const item of items) {
    await safeCall(() => presetApi.create(item), null)
  }
}

async function importPersonas(items: PersonaInput[]) {
  for (const item of items) {
    await safeCall(() => personaApi.create(item), null)
  }
}

async function importDiaries(items: DiaryInput[]) {
  for (const item of items) {
    await safeCall(() => diaryApi.create(item), null)
  }
}

async function importSms(items: BackupFile['data']['sms']) {
  for (const item of items) {
    const thread = await safeCall(() => smsApi.createThread({
      recipient: item.thread.recipient,
      number: item.thread.number,
      character_id: item.thread.character_id,
      avatar: item.thread.avatar,
    }), { id: 0 } as any)
    const threadId = thread.id
    if (!threadId) continue
    for (const message of item.messages || []) {
      await safeCall(() => smsApi.createMessage(threadId, {
        role: message.role,
        content: message.content,
      }), null)
    }
  }
}

async function importCalls(items: BackupFile['data']['calls']) {
  for (const item of items) {
    await safeCall(() => callApi.create(item), null)
  }
}

async function importWallet(data: BackupFile['data']['wallet']) {
  if (typeof data.balance === 'number') {
    await safeCall(() => walletApi.setBalance(data.balance as number), null)
  }
  for (const tx of data.transactions || []) {
    await safeCall(() => walletApi.createTransaction(tx), null)
  }
  mergeArrayStorage(REDPACKET_KEY, data.redpackets)
  mergeArrayStorage(TRANSFER_KEY, data.transfers)
}

async function importGames(items: BackupFile['data']['games']) {
  for (const item of items) {
    await safeCall(() => gameApi.createRecord(item), null)
  }
}

function importLocalState(file: BackupFile) {
  const settings = sanitizeSettings(file.data.settings || {})
  const currentSettings = safeParse(localStorage.getItem(SETTINGS_KEY), {})
  localStorage.setItem(SETTINGS_KEY, JSON.stringify({ ...currentSettings, ...settings }))

  mergeObjectStorage(IMAGE_GEN_KEY, sanitizeImageGenConfig(file.data.image_gen_config || {}))

  const localState = file.data.local_state || {}
  mergeArrayStorage(PRESETS_KEY, localState.presets_cache)
  mergeArrayStorage(WORLDBOOKS_KEY, localState.worldbooks)
  mergeArrayStorage(CHAT_CONVERSATIONS_KEY, localState.chat_conversations)
  mergeArrayStorage(SMS_STORAGE_KEY, localState.sms_local_cache)
  mergeArrayStorage(CALL_HISTORY_KEY, localState.call_history_local)
  mergeObjectStorage(WALLET_KEY, localState.wallet_local)
  mergeArrayStorage(TX_KEY, localState.wallet_transactions_local)
  mergeArrayStorage(REDPACKET_KEY, localState.wallet_redpackets_local)
  mergeArrayStorage(TRANSFER_KEY, localState.wallet_transfers_local)
  mergeArrayStorage(DIARY_LOCAL_KEY, localState.diary_local)
  mergeArrayStorage(SECURITY_LOG_KEY, localState.security_logs)
  mergeObjectStorage(GLOBAL_EJS_VARS_KEY, localState.global_template_vars)

  setScalarStorage(ACTIVE_PRESET_KEY, localState.active_preset_id)
  setScalarStorage(CURRENT_PERSONA_KEY, localState.current_persona_id)
  setScalarStorage(CURRENT_CHARACTER_KEY, localState.current_character_id)
  setScalarStorage(CHAT_AUTO_IMAGE_KEY, localState.chat_auto_image)
  setScalarStorage(SOCIAL_AUTO_IMAGE_KEY, localState.social_auto_image)

  const chatMessages = localState.chat_messages || {}
  for (const [key, value] of Object.entries(chatMessages)) {
    if (key.startsWith('chat-messages-')) {
      mergeArrayStorage(key, value)
    }
  }

  const localTemplateVars = localState.local_template_vars || {}
  for (const [key, value] of Object.entries(localTemplateVars)) {
    if (key.startsWith('ejs-vars-chat-')) {
      mergeObjectStorage(key, value)
    }
  }

  for (const [name, storageKey] of Object.entries(SOCIAL_STORAGE_KEYS)) {
    if (Object.prototype.hasOwnProperty.call(file.data.social_data || {}, name)) {
      localStorage.setItem(storageKey, JSON.stringify(file.data.social_data[name]))
    }
  }
}

export async function importAppBackup(raw: unknown) {
  const file = validateBackupFile(raw)
  importLocalState(file)

  await importCharacters(file.data.characters || [])
  await importPresets(file.data.presets || [])
  await importPersonas(file.data.personas || [])
  await importDiaries(file.data.diaries || [])
  await importSms(file.data.sms || [])
  await importCalls(file.data.calls || [])
  await importWallet(file.data.wallet || { balance: null, transactions: [], redpackets: [], transfers: [] })
  await importGames(file.data.games || [])

  return {
    characters: (file.data.characters || []).length,
    presets: (file.data.presets || []).length,
    personas: (file.data.personas || []).length,
    diaries: (file.data.diaries || []).length,
    sms_threads: (file.data.sms || []).length,
    calls: (file.data.calls || []).length,
    games: (file.data.games || []).length,
  }
}

export async function downloadAppBackup() {
  const backup = await exportAppBackup()
  const blob = new Blob([JSON.stringify(backup, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `fzsmphone-backup-${new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')}.json`
  a.click()
  URL.revokeObjectURL(url)
}

export async function readBackupFile(file: File): Promise<BackupFile> {
  const text = await file.text()
  const parsed = JSON.parse(text)
  return validateBackupFile(parsed)
}

export async function pingAuth() {
  await api.get('/api/users/me')
}
