/**
 * EJS 模板引擎 - 兼容 SillyTavern ST-Prompt-Template 插件
 *
 * 支持：
 * - EJS 标签：<% %> (执行), <%= %> (转义输出), <%- %> (原始输出)
 * - ST 原生宏：{{char}}, {{user}}, {{random:a,b,c}}, {{roll:2d6}}, {{time}}, {{date}}, {{idle_duration}}
 * - 内置函数：getvar/setvar/addvar, getchar/getuser/getwi/getpreset, messages/lastMessage 等
 * - 变量作用域：global (跨会话), local (当前聊天), message (单次处理)
 */

// ============ 变量存储 ============

const GLOBAL_VARS_KEY = 'ejs-vars-global'

function getGlobalVars(): Record<string, string> {
  try {
    const s = localStorage.getItem(GLOBAL_VARS_KEY)
    return s ? JSON.parse(s) : {}
  } catch {
    return {}
  }
}

function setGlobalVars(vars: Record<string, string>) {
  localStorage.setItem(GLOBAL_VARS_KEY, JSON.stringify(vars))
}

function getLocalVarsKey(chatId: string): string {
  return `ejs-vars-chat-${chatId}`
}

function getLocalVars(chatId: string): Record<string, string> {
  try {
    const s = localStorage.getItem(getLocalVarsKey(chatId))
    return s ? JSON.parse(s) : {}
  } catch {
    return {}
  }
}

function setLocalVars(chatId: string, vars: Record<string, string>) {
  localStorage.setItem(getLocalVarsKey(chatId), JSON.stringify(vars))
}

// ============ 上下文接口 ============

export interface EjsContext {
  charName: string
  userName: string
  chatId: string
  messages: Array<{ role: string; content: string }>
  characterData?: any
  presetData?: any
  worldBookEntries?: any[]
}

// ============ ST 原生宏处理 ============

/**
 * 处理 SillyTavern 原生双花括号宏
 */
export function processSTMacros(text: string, ctx: EjsContext): string {
  let result = text

  // {{char}} / {{user}}
  result = result.replace(/\{\{char\}\}/gi, ctx.charName)
  result = result.replace(/\{\{user\}\}/gi, ctx.userName)

  // {{time}} - 当前时间 HH:MM
  result = result.replace(/\{\{time\}\}/gi, () => {
    const now = new Date()
    return `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
  })

  // {{date}} - 当前日期
  result = result.replace(/\{\{date\}\}/gi, () => {
    const now = new Date()
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
  })

  // {{weekday}} - 星期几
  result = result.replace(/\{\{weekday\}\}/gi, () => {
    const days = ['日', '一', '二', '三', '四', '五', '六']
    return `星期${days[new Date().getDay()]}`
  })

  // {{random:a,b,c}} - 随机选择
  result = result.replace(/\{\{random:([^}]+)\}\}/gi, (_, items: string) => {
    const options = items.split(',').map(s => s.trim())
    return options[Math.floor(Math.random() * options.length)] || ''
  })

  // {{roll:NdM}} - 掷骰子，如 {{roll:2d6}}
  result = result.replace(/\{\{roll:(\d+)d(\d+)\}\}/gi, (_, countStr: string, sidesStr: string) => {
    const count = parseInt(countStr) || 1
    const sides = parseInt(sidesStr) || 6
    let total = 0
    for (let i = 0; i < Math.min(count, 100); i++) {
      total += Math.floor(Math.random() * sides) + 1
    }
    return String(total)
  })

  // {{roll:M}} - 简单掷骰 1dM
  result = result.replace(/\{\{roll:(\d+)\}\}/gi, (_, sidesStr: string) => {
    const sides = parseInt(sidesStr) || 20
    return String(Math.floor(Math.random() * sides) + 1)
  })

  // {{idle_duration}} - 距上次消息的时间（简化实现）
  result = result.replace(/\{\{idle_duration\}\}/gi, () => {
    return '刚刚'
  })

  // {{lastMessage}} - 最后一条消息
  result = result.replace(/\{\{lastMessage\}\}/gi, () => {
    const msgs = ctx.messages
    if (msgs.length === 0) return ''
    return msgs[msgs.length - 1].content
  })

  // {{lastUserMessage}} - 最后一条用户消息
  result = result.replace(/\{\{lastUserMessage\}\}/gi, () => {
    for (let i = ctx.messages.length - 1; i >= 0; i--) {
      if (ctx.messages[i].role === 'user') return ctx.messages[i].content
    }
    return ''
  })

  // {{lastCharMessage}} - 最后一条角色消息
  result = result.replace(/\{\{lastCharMessage\}\}/gi, () => {
    for (let i = ctx.messages.length - 1; i >= 0; i--) {
      if (ctx.messages[i].role === 'assistant') return ctx.messages[i].content
    }
    return ''
  })

  // {{messageCount}} - 消息数量（不含 system）
  result = result.replace(/\{\{messageCount\}\}/gi, () => {
    return String(ctx.messages.filter(m => m.role !== 'system').length)
  })

  // {{chatId}}
  result = result.replace(/\{\{chatId\}\}/gi, ctx.chatId)

  // {{// comment}} - 注释，移除
  result = result.replace(/\{\{\/\/[^}]*\}\}/g, '')

  return result
}

// ============ EJS 引擎 ============

/**
 * 构建 EJS 沙箱中可用的内置函数
 */
function buildEjsBuiltins(ctx: EjsContext, messageVars: Record<string, string>) {
  const globalVars = getGlobalVars()
  const localVars = getLocalVars(ctx.chatId)

  return {
    // ---- 变量操作 ----
    getvar(name: string, scope?: string): string {
      if (scope === 'global') return globalVars[name] ?? ''
      if (scope === 'local') return localVars[name] ?? ''
      // 默认优先级: message > local > global
      return messageVars[name] ?? localVars[name] ?? globalVars[name] ?? ''
    },

    setvar(name: string, value: string, scope?: string): string {
      const strVal = String(value ?? '')
      if (scope === 'global') {
        globalVars[name] = strVal
        setGlobalVars(globalVars)
      } else if (scope === 'local') {
        localVars[name] = strVal
        setLocalVars(ctx.chatId, localVars)
      } else {
        // 默认设置到 local
        localVars[name] = strVal
        setLocalVars(ctx.chatId, localVars)
      }
      return ''
    },

    addvar(name: string, increment: number, scope?: string): string {
      const current = scope === 'global'
        ? Number(globalVars[name] || 0)
        : Number(localVars[name] || 0)
      const newVal = String(current + (Number(increment) || 0))
      if (scope === 'global') {
        globalVars[name] = newVal
        setGlobalVars(globalVars)
      } else {
        localVars[name] = newVal
        setLocalVars(ctx.chatId, localVars)
      }
      return newVal
    },

    // 获取所有变量（某个作用域）
    variables(scope?: string): Record<string, string> {
      if (scope === 'global') return { ...globalVars }
      if (scope === 'local') return { ...localVars }
      return { ...globalVars, ...localVars, ...messageVars }
    },

    // ---- 角色/用户信息 ----
    getchar(field?: string): string {
      if (!ctx.characterData) return ''
      if (!field) return ctx.charName
      return String(ctx.characterData[field] ?? '')
    },

    getuser(field?: string): string {
      if (!field) return ctx.userName
      return ''
    },

    // ---- 世界书 ----
    getwi(titleOrIndex?: string | number): string {
      if (!ctx.worldBookEntries?.length) return ''
      if (titleOrIndex === undefined) {
        return ctx.worldBookEntries.map(e => e.content || '').join('\n')
      }
      if (typeof titleOrIndex === 'number') {
        return ctx.worldBookEntries[titleOrIndex]?.content || ''
      }
      const entry = ctx.worldBookEntries.find(e =>
        e.title === titleOrIndex || e.keywords?.includes(titleOrIndex)
      )
      return entry?.content || ''
    },

    // ---- 预设信息 ----
    getpreset(field?: string): string {
      if (!ctx.presetData) return ''
      if (!field) return ctx.presetData.name || ''
      return String(ctx.presetData[field] ?? '')
    },

    // ---- 消息相关 ----
    get messages() {
      return ctx.messages
    },

    get lastMessage(): string {
      const msgs = ctx.messages
      return msgs.length > 0 ? msgs[msgs.length - 1].content : ''
    },

    get messageCount(): number {
      return ctx.messages.filter(m => m.role !== 'system').length
    },

    get chatId(): string {
      return ctx.chatId
    },

    // ---- 实用函数 ----
    random(...args: any[]): any {
      if (args.length === 0) return Math.random()
      if (args.length === 1 && typeof args[0] === 'number') {
        return Math.floor(Math.random() * args[0])
      }
      // random(a, b, c) - 随机选择一个
      return args[Math.floor(Math.random() * args.length)]
    },

    pick(...items: any[]): any {
      if (items.length === 1 && Array.isArray(items[0])) {
        const arr = items[0]
        return arr[Math.floor(Math.random() * arr.length)]
      }
      return items[Math.floor(Math.random() * items.length)]
    },

    roll(spec: string | number = 6): number {
      if (typeof spec === 'number') {
        return Math.floor(Math.random() * spec) + 1
      }
      const match = String(spec).match(/^(\d+)d(\d+)$/i)
      if (match) {
        const count = parseInt(match[1]) || 1
        const sides = parseInt(match[2]) || 6
        let total = 0
        for (let i = 0; i < Math.min(count, 100); i++) {
          total += Math.floor(Math.random() * sides) + 1
        }
        return total
      }
      return Math.floor(Math.random() * (parseInt(spec) || 6)) + 1
    },

    time(): string {
      const now = new Date()
      return `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
    },

    date(): string {
      const now = new Date()
      return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
    },

    trim(s: string): string {
      return String(s || '').trim()
    },

    // ---- 数学 ----
    abs: Math.abs,
    floor: Math.floor,
    ceil: Math.ceil,
    round: Math.round,
    min: Math.min,
    max: Math.max,
  }
}

/**
 * 简易 EJS 模板编译和执行
 * 支持 <% %> (执行), <%= %> (输出并转义), <%- %> (原始输出)
 *
 * 实现原理：将 EJS 模板转换为 JS 函数字符串，然后用 new Function 执行
 */
function compileEjs(template: string): (data: Record<string, any>) => string {
  let code = "let __output = '';\n"
  let cursor = 0

  // 匹配所有 EJS 标签
  const tagRegex = /<%([=-]?)([\s\S]*?)%>/g
  let match: RegExpExecArray | null

  while ((match = tagRegex.exec(template)) !== null) {
    const prefix = match[1]   // '', '=', '-'
    const content = match[2]  // 标签内容
    const beforeTag = template.slice(cursor, match.index)

    // 添加纯文本部分
    if (beforeTag) {
      code += `__output += ${JSON.stringify(beforeTag)};\n`
    }

    if (prefix === '=' || prefix === '-') {
      // 输出标签（= 转义, - 原始）—— 在这个场景下都不做 HTML 转义
      code += `__output += String(${content.trim()});\n`
    } else {
      // 执行标签
      code += content + '\n'
    }

    cursor = match.index + match[0].length
  }

  // 添加剩余文本
  const remaining = template.slice(cursor)
  if (remaining) {
    code += `__output += ${JSON.stringify(remaining)};\n`
  }

  code += 'return __output;'

  return (data: Record<string, any>) => {
    try {
      // 将 data 的所有属性注入为局部变量
      const keys = Object.keys(data)
      const values = keys.map(k => data[k])
      const fn = new Function(...keys, code)
      return fn(...values)
    } catch (err) {
      console.warn('[EJS] Template execution error:', err)
      // 出错时返回原模板，避免丢失内容
      return template
    }
  }
}

// ============ 公开 API ============

/**
 * 处理文本中的 EJS 模板和 ST 宏
 * @param text 待处理的文本
 * @param ctx 上下文信息
 * @param stage 'prompt' (发送前处理) | 'render' (渲染时处理)
 * @returns 处理后的文本
 */
export function processTemplate(
  text: string,
  ctx: EjsContext,
  stage: 'prompt' | 'render' = 'prompt',
): string {
  if (!text) return text

  // 1. 先处理 ST 原生宏（{{...}}）
  let result = processSTMacros(text, ctx)

  // 2. 检查是否包含 EJS 标签
  if (!result.includes('<%')) {
    return result
  }

  // 3. 编译并执行 EJS 模板
  const messageVars: Record<string, string> = {}
  const builtins = buildEjsBuiltins(ctx, messageVars)

  try {
    const compiled = compileEjs(result)
    result = compiled({
      ...builtins,
      // 额外注入常用变量
      char: ctx.charName,
      user: ctx.userName,
      stage,
    })
  } catch (err) {
    console.warn('[EJS] Processing error:', err)
    // 出错时返回宏处理后的文本
  }

  return result
}

/**
 * 批量处理消息列表中的模板
 */
export function processMessagesTemplate(
  messages: Array<{ role: string; content: string }>,
  ctx: EjsContext,
  stage: 'prompt' | 'render' = 'prompt',
): Array<{ role: string; content: string }> {
  return messages.map(msg => ({
    role: msg.role,
    content: processTemplate(msg.content, ctx, stage),
  }))
}

/**
 * 处理 @INJECT 语法
 * 格式: @INJECT position=<pos> depth=<depth> role=<role>
 * 在文本中查找 @INJECT 指令并提取注入配置
 */
export interface InjectDirective {
  content: string
  position: number  // 0-4, 与世界书 position 一致
  depth: number
  role: 'system' | 'user' | 'assistant'
}

export function extractInjectDirectives(text: string): {
  cleanText: string
  directives: InjectDirective[]
} {
  const directives: InjectDirective[] = []
  const lines = text.split('\n')
  const cleanLines: string[] = []
  let currentDirective: Partial<InjectDirective> | null = null
  let directiveContent: string[] = []

  for (const line of lines) {
    const injectMatch = line.match(
      /^@INJECT\s+(?:position=(\d+))?\s*(?:depth=(\d+))?\s*(?:role=(system|user|assistant))?/i
    )

    if (injectMatch) {
      // 保存之前的 directive
      if (currentDirective && directiveContent.length > 0) {
        directives.push({
          content: directiveContent.join('\n').trim(),
          position: currentDirective.position ?? 0,
          depth: currentDirective.depth ?? 4,
          role: currentDirective.role ?? 'system',
        })
      }

      currentDirective = {
        position: injectMatch[1] ? parseInt(injectMatch[1]) : 0,
        depth: injectMatch[2] ? parseInt(injectMatch[2]) : 4,
        role: (injectMatch[3] as any) || 'system',
      }
      directiveContent = []
    } else if (currentDirective) {
      if (line.trim() === '@END_INJECT') {
        directives.push({
          content: directiveContent.join('\n').trim(),
          position: currentDirective.position ?? 0,
          depth: currentDirective.depth ?? 4,
          role: currentDirective.role ?? 'system',
        })
        currentDirective = null
        directiveContent = []
      } else {
        directiveContent.push(line)
      }
    } else {
      cleanLines.push(line)
    }
  }

  // 处理未闭合的 directive
  if (currentDirective && directiveContent.length > 0) {
    directives.push({
      content: directiveContent.join('\n').trim(),
      position: currentDirective.position ?? 0,
      depth: currentDirective.depth ?? 4,
      role: currentDirective.role ?? 'system',
    })
  }

  return {
    cleanText: cleanLines.join('\n'),
    directives,
  }
}
