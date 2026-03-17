/**
 * AI 生图配置系统 — 对齐酒馆脚本 chunk 920
 * localStorage key: phone_novelai_config
 */

// ==================== 常量 ====================

export const IMAGE_ASPECT_RATIO_OPTIONS = [
    'auto', '1:1', '2:3', '3:2', '3:4', '4:3', '4:5', '5:4', '9:16', '16:9', '21:9',
] as const

export const GEMINI_IMAGE_SIZE_OPTIONS = ['auto', '1K', '2K', '4K'] as const

export type AspectRatio = (typeof IMAGE_ASPECT_RATIO_OPTIONS)[number]
export type GeminiImageSize = (typeof GEMINI_IMAGE_SIZE_OPTIONS)[number]
export type ApiFormat = 'novelai' | 'openai' | 'gemini'
export type PromptComposeMode = 'tags' | 'natural' | 'none'

// ==================== 接口 ====================

export interface NovelAiProviderConfig {
    url: string
    key: string
    model: string
    promptPrefix: string
    promptSuffix: string
    negativePrompt: string
    width: number
    height: number
    steps: number
    cfg: number
    sampler: string
    scheduler: string
    seed: number
    ucPreset: number
    qualityToggle: boolean
    autoSmea: boolean
    decrisper: boolean
    sm: boolean
    smDyn: boolean
    cfgRescale: number
    varietyBoost: boolean
    characterReferenceEnabled: boolean
    characterReferenceStyleAware: boolean
    characterReferenceFidelity: number
}

export interface OpenAiProviderConfig {
  url: string
  key: string
  model: string
  promptPrefix: string
  promptSuffix: string
  aspectRatio: AspectRatio
  requestFormat: 'chat' | 'direct'  // chat: OpenAI messages format, direct: simple prompt format
}

export interface GeminiProviderConfig {
    url: string
    key: string
    model: string
    promptPrefix: string
    promptSuffix: string
    aspectRatio: AspectRatio
    imageSize: GeminiImageSize
}

export interface ImageGenConfig {
    apiFormat: ApiFormat
    promptComposeMode: PromptComposeMode
    novelai: NovelAiProviderConfig
    openai: OpenAiProviderConfig
    gemini: GeminiProviderConfig
}

// ==================== 默认值 ====================

const STORAGE_KEY = 'phone_novelai_config'

const DEFAULT_NOVELAI: NovelAiProviderConfig = {
  url: '',
  key: '',
  model: 'nai-diffusion-4-5-full',
  promptPrefix: '',
  promptSuffix: '',
  negativePrompt: '',
  width: 1024,
  height: 1024,
  steps: 28,
  cfg: 5,
  sampler: 'k_euler',
  scheduler: 'karras',
  seed: -1,
  ucPreset: 0,
  qualityToggle: true,
  autoSmea: false,
  decrisper: false,
  sm: false,
  smDyn: false,
  cfgRescale: 0,
  varietyBoost: false,
  characterReferenceEnabled: false,
  characterReferenceStyleAware: true,
  characterReferenceFidelity: 0.8,
}

const DEFAULT_OPENAI: OpenAiProviderConfig = {
  url: '',
  key: '',
  model: '',
  promptPrefix: '',
  promptSuffix: '',
  aspectRatio: 'auto',
  requestFormat: 'chat',
}

// Nano/Banano 模型列表（用于参考）
export const NANO_BANANO_MODELS = [
  'nano-banano',
  'nano-banano-v2', 
  'banano',
  'banano-diffusion',
  'nai-diffusion-4-5-full',
  'nai-diffusion-4-standard',
] as const

// Grok 模型列表（用于参考）
export const GROK_MODELS = [
  'grok-2-vision-1212',
  'grok-2-vision',
  'grok-2',
  'grok-beta',
  'grok-vision-beta',
] as const

const DEFAULT_GEMINI: GeminiProviderConfig = {
    url: '',
    key: '',
    model: '',
    promptPrefix: '',
    promptSuffix: '',
    aspectRatio: 'auto',
    imageSize: 'auto',
}

// ==================== 工具函数 ====================

function clamp(val: number, min: number, max: number): number {
    return Number.isFinite(val) ? Math.min(max, Math.max(min, val)) : min
}

function toNum(val: unknown, fallback: number): number {
    if (typeof val === 'number' && Number.isFinite(val)) return val
    if (typeof val === 'string' && val.trim() !== '') {
        const n = Number(val)
        if (Number.isFinite(n)) return n
    }
    return fallback
}

function toInt(val: unknown, fallback: number): number {
    const n = toNum(val, fallback)
    return Number.isFinite(n) ? Math.trunc(n) : fallback
}

function snapTo64(val: unknown, fallback: number): number {
    const n = clamp(toInt(val, fallback), 64, 4096)
    return 64 * Math.round(n / 64)
}

function toSeed(val: unknown, fallback: number): number {
    const n = toInt(val, fallback)
    return n < 0 ? -1 : clamp(n, 0, 4294967295)
}

function toApiFormat(val: unknown, fallback: ApiFormat): ApiFormat {
    const s = String(val ?? '').trim()
    return s === 'novelai' || s === 'openai' || s === 'gemini' ? s : fallback
}

function toPromptMode(val: unknown, fallback: PromptComposeMode): PromptComposeMode {
    const s = String(val ?? '').trim()
    return s === 'tags' || s === 'natural' ? s : fallback
}

function toAspectRatio(val: unknown, fallback: AspectRatio): AspectRatio {
    const s = String(val ?? '').trim()
    return (IMAGE_ASPECT_RATIO_OPTIONS as readonly string[]).includes(s) ? (s as AspectRatio) : fallback
}

function toImageSize(val: unknown, fallback: GeminiImageSize): GeminiImageSize {
    const s = String(val ?? '').trim()
    return s === '1K' || s === '2K' || s === '4K' ? s : fallback === 'auto' ? 'auto' : fallback
}

// ==================== 解析 ====================

function parseNovelAi(raw: Record<string, unknown> | undefined, base: NovelAiProviderConfig): NovelAiProviderConfig {
    const o = (raw && typeof raw === 'object' ? raw : {}) as Record<string, unknown>

    // 角色参考图模式兼容
    const refMode = String(o.referenceImageMode ?? o.reference_image_mode ?? o.referenceMode ?? o.reference_mode ?? '').trim()
    const charRefEnabled = typeof o.characterReferenceEnabled === 'boolean'
        ? o.characterReferenceEnabled
        : refMode === 'character_reference' || refMode === 'auto' || (refMode !== 'vibe' && base.characterReferenceEnabled)

    return {
        url: typeof o.url === 'string' ? o.url : base.url,
        key: typeof o.key === 'string' ? o.key : base.key,
        model: typeof o.model === 'string' ? o.model : base.model,
        promptPrefix: typeof o.promptPrefix === 'string' ? o.promptPrefix : base.promptPrefix,
        promptSuffix: typeof o.promptSuffix === 'string' ? o.promptSuffix : base.promptSuffix,
        negativePrompt: typeof o.negativePrompt === 'string' ? o.negativePrompt : base.negativePrompt,
        width: snapTo64(o.width, base.width),
        height: snapTo64(o.height, base.height),
        steps: clamp(toInt(o.steps, base.steps), 1, 60),
        cfg: clamp(toNum(o.cfg, base.cfg), 0, 30),
        sampler: typeof o.sampler === 'string' ? o.sampler : base.sampler,
        scheduler: typeof o.scheduler === 'string' ? o.scheduler : base.scheduler,
        seed: toSeed(o.seed, base.seed),
        ucPreset: clamp(toInt(o.ucPreset, base.ucPreset), 0, 3),
        qualityToggle: typeof o.qualityToggle === 'boolean' ? o.qualityToggle : base.qualityToggle,
        autoSmea: typeof o.autoSmea === 'boolean' ? o.autoSmea : base.autoSmea,
        decrisper: typeof o.decrisper === 'boolean' ? o.decrisper : base.decrisper,
        sm: typeof o.sm === 'boolean' ? o.sm : base.sm,
        smDyn: typeof o.smDyn === 'boolean' ? o.smDyn : base.smDyn,
        cfgRescale: clamp(toNum(o.cfgRescale, base.cfgRescale), 0, 1),
        varietyBoost: typeof o.varietyBoost === 'boolean' ? o.varietyBoost : base.varietyBoost,
        characterReferenceEnabled: charRefEnabled,
        characterReferenceStyleAware: typeof o.characterReferenceStyleAware === 'boolean' ? o.characterReferenceStyleAware : base.characterReferenceStyleAware,
        characterReferenceFidelity: clamp(toNum(o.characterReferenceFidelity, base.characterReferenceFidelity), 0, 1),
    }
}

function parseOpenAi(raw: Record<string, unknown> | undefined, base: OpenAiProviderConfig): OpenAiProviderConfig {
  const o = (raw && typeof raw === 'object' ? raw : {}) as Record<string, unknown>
  const rf = String(o.requestFormat ?? '').trim()
  return {
    url: typeof o.url === 'string' ? o.url : base.url,
    key: typeof o.key === 'string' ? o.key : base.key,
    model: typeof o.model === 'string' ? o.model : base.model,
    promptPrefix: typeof o.promptPrefix === 'string' ? o.promptPrefix : base.promptPrefix,
    promptSuffix: typeof o.promptSuffix === 'string' ? o.promptSuffix : base.promptSuffix,
    aspectRatio: toAspectRatio(o.aspectRatio, base.aspectRatio),
    requestFormat: rf === 'direct' ? 'direct' : 'chat',
  }
}

function parseGemini(raw: Record<string, unknown> | undefined, base: GeminiProviderConfig): GeminiProviderConfig {
    const o = (raw && typeof raw === 'object' ? raw : {}) as Record<string, unknown>
    const imgSize = o.imageSize ?? o.geminiImageSize
    return {
        url: typeof o.url === 'string' ? o.url : base.url,
        key: typeof o.key === 'string' ? o.key : base.key,
        model: typeof o.model === 'string' ? o.model : base.model,
        promptPrefix: typeof o.promptPrefix === 'string' ? o.promptPrefix : base.promptPrefix,
        promptSuffix: typeof o.promptSuffix === 'string' ? o.promptSuffix : base.promptSuffix,
        aspectRatio: toAspectRatio(o.aspectRatio, base.aspectRatio),
        imageSize: toImageSize(imgSize, base.imageSize),
    }
}

// ==================== 公开 API ====================

export function getDefaultImageGenConfig(): ImageGenConfig {
    return {
        apiFormat: 'novelai',
        promptComposeMode: 'none',
        novelai: { ...DEFAULT_NOVELAI },
        openai: { ...DEFAULT_OPENAI },
        gemini: { ...DEFAULT_GEMINI },
    }
}

export function loadImageGenConfig(): ImageGenConfig {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        if (!raw) return getDefaultImageGenConfig()

        const data = JSON.parse(raw)
        if (!data || typeof data !== 'object') return getDefaultImageGenConfig()

        const apiFormat = toApiFormat(data.apiFormat, 'novelai')
        const promptComposeMode = toPromptMode(data.promptComposeMode, 'none')

        // 新版格式（有 novelai/openai/gemini 子对象）
        if ((data.novelai && typeof data.novelai === 'object') ||
            (data.openai && typeof data.openai === 'object') ||
            (data.gemini && typeof data.gemini === 'object')) {
            return {
                apiFormat,
                promptComposeMode,
                novelai: parseNovelAi(data.novelai, DEFAULT_NOVELAI),
                openai: parseOpenAi(data.openai, DEFAULT_OPENAI),
                gemini: parseGemini(data.gemini, DEFAULT_GEMINI),
            }
        }

        // 旧版格式（扁平结构，只用于一个 provider）
        const nai = parseNovelAi(data as Record<string, unknown>, DEFAULT_NOVELAI)
        const oai = parseOpenAi(data as Record<string, unknown>, DEFAULT_OPENAI)
        const gem = parseGemini(data as Record<string, unknown>, DEFAULT_GEMINI)

        // 只保留当前 format 对应的数据
        if (apiFormat !== 'novelai') {
            Object.assign(nai, { url: DEFAULT_NOVELAI.url, key: DEFAULT_NOVELAI.key, model: DEFAULT_NOVELAI.model })
        }
        if (apiFormat !== 'openai') {
            Object.assign(oai, { url: DEFAULT_OPENAI.url, key: DEFAULT_OPENAI.key, model: DEFAULT_OPENAI.model })
        }
        if (apiFormat !== 'gemini') {
            Object.assign(gem, { url: DEFAULT_GEMINI.url, key: DEFAULT_GEMINI.key, model: DEFAULT_GEMINI.model })
        }

        return { apiFormat, promptComposeMode, novelai: nai, openai: oai, gemini: gem }
    } catch {
        return getDefaultImageGenConfig()
    }
}

export function saveImageGenConfig(config: ImageGenConfig): void {
    const data: ImageGenConfig = {
        apiFormat: toApiFormat(config?.apiFormat, 'novelai'),
        promptComposeMode: toPromptMode(config?.promptComposeMode, 'none'),
        novelai: parseNovelAi(config?.novelai as unknown as Record<string, unknown>, DEFAULT_NOVELAI),
        openai: parseOpenAi(config?.openai as unknown as Record<string, unknown>, DEFAULT_OPENAI),
        gemini: parseGemini(config?.gemini as unknown as Record<string, unknown>, DEFAULT_GEMINI),
    }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
}
