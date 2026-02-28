export interface GradientIconDef {
    char: string
    bg: string
}

export const contentIcons: Record<string, GradientIconDef> = {
    // ===== Shopping / Digital =====
    headphones: { char: '♪', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    phone: { char: '▢', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    laptop: { char: '⌘', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    watch: { char: '◔', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    digital: { char: '▢', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },

    // ===== Fashion =====
    dress: { char: '✾', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    couple_wear: { char: '⚯', bg: 'linear-gradient(135deg, #a18cd1, #fbc2eb)' },
    fashion: { char: '✾', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },

    // ===== Beauty =====
    cosmetics: { char: '✦', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    lipstick: { char: '♥', bg: 'linear-gradient(135deg, #ff9a9e, #fad0c4)' },
    beauty: { char: '✦', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },

    // ===== Food =====
    chocolate: { char: '◆', bg: 'linear-gradient(135deg, #795548, #D7CCC8)' },
    matcha: { char: '❦', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    food: { char: '◈', bg: 'linear-gradient(135deg, #ff9a9e, #fad0c4)' },

    // ===== Home =====
    candle: { char: '△', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    pillow: { char: '▢', bg: 'linear-gradient(135deg, #a1c4fd, #c2e9fb)' },
    home: { char: '⌂', bg: 'linear-gradient(135deg, #a1c4fd, #c2e9fb)' },

    // ===== Gifts =====
    necklace: { char: '○', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    flower: { char: '✿', bg: 'linear-gradient(135deg, #ff9a9e, #fecfef)' },
    gift: { char: '★', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },

    // ===== Navigation / Status =====
    all: { char: '⌂', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    search: { char: '◎', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    cart: { char: '▢', bg: 'linear-gradient(135deg, #ff9a9e, #fad0c4)' },
    heart: { char: '♥', bg: 'linear-gradient(135deg, #ff9a9e, #fad0c4)' },
    fire: { char: '▲', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    checkmark: { char: '✓', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    star: { char: '★', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },

    // ===== Mood =====
    mood_happy: { char: '◠', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    mood_excited: { char: '✦', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    mood_calm: { char: '○', bg: 'linear-gradient(135deg, #a1c4fd, #c2e9fb)' },
    mood_sad: { char: '◡', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    mood_angry: { char: '✕', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    mood_tired: { char: '—', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },

    // ===== Weather =====
    weather_sunny: { char: '☀', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    weather_cloudy: { char: '☁', bg: 'linear-gradient(135deg, #bdc3c7, #95a5a6)' },
    weather_rainy: { char: '☂', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    weather_snowy: { char: '❄', bg: 'linear-gradient(135deg, #e0e5ec, #d2d6dc)' },
    weather_windy: { char: '≋', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },

    // ===== Diary =====
    diary: { char: '▤', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    empty_diary: { char: '▤', bg: 'linear-gradient(135deg, #bdc3c7, #95a5a6)' },

    // ===== Social =====
    avatar_male: { char: '○', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    avatar_female: { char: '○', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    photo: { char: '▣', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    camera: { char: '▣', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    music: { char: '♪', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    location: { char: '◆', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    calendar: { char: '▦', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    letter: { char: '✉', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    ring: { char: '○', bg: 'linear-gradient(135deg, #ffd93d, #f6d365)' },
    cake: { char: '◈', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    plane: { char: '➤', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    tree: { char: '✦', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    fireworks: { char: '✸', bg: 'linear-gradient(135deg, #ff6b6b, #fda085)' },
    couple: { char: '♥', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },

    // ===== Entertainment =====
    gamepad: { char: '▶', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    puzzle: { char: '⊞', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    target: { char: '◎', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    chess: { char: '♟', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },
    cards: { char: '♠', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    dice: { char: '⚄', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    slots: { char: '▣', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    wheel: { char: '◎', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    trophy: { char: '★', bg: 'linear-gradient(135deg, #ffd93d, #f6d365)' },
    movie: { char: '▷', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    live: { char: '●', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    mic: { char: '◉', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    radio: { char: '♪', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },

    // ===== Takeaway =====
    rice: { char: '◈', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    noodle: { char: '≈', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    burger: { char: '◆', bg: 'linear-gradient(135deg, #ff9a9e, #fad0c4)' },
    drink: { char: '◇', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    cake_food: { char: '●', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    hotpot: { char: '◎', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    chicken: { char: '▲', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    salad: { char: '✿', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    delivery: { char: '➤', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    clock: { char: '◔', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },

    // ===== Tools =====
    settings: { char: '⚙', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },
    shield: { char: '■', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    key: { char: '■', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    lock: { char: '■', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },
    warning: { char: '△', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    edit: { char: '✎', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    copy: { char: '▤', bg: 'linear-gradient(135deg, #a1c4fd, #c2e9fb)' },
    upload: { char: '↑', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    download: { char: '↓', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    trash: { char: '✕', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    save: { char: '▣', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    note: { char: '▤', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    book: { char: '▤', bg: 'linear-gradient(135deg, #a18cd1, #fbc2eb)' },
    link: { char: '↗', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    pin: { char: '◆', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    box: { char: '☐', bg: 'linear-gradient(135deg, #bdc3c7, #95a5a6)' },
    sparkle: { char: '✦', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    rocket: { char: '➤', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },
    moon: { char: '☽', bg: 'linear-gradient(135deg, #667eea, #764ba2)' },

    // ===== Preset / Chat =====
    roleplay: { char: '◈', bg: 'linear-gradient(135deg, #f093fb, #f5576c)' },
    robot: { char: '◎', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    writing: { char: '✎', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    learning: { char: '▥', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    fun: { char: '▶', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    crystal: { char: '◇', bg: 'linear-gradient(135deg, #a18cd1, #fbc2eb)' },
    art: { char: '✧', bg: 'linear-gradient(135deg, #f6d365, #fda085)' },
    cat: { char: '◈', bg: 'linear-gradient(135deg, #ffecd2, #fcb69f)' },
    ghost: { char: '◈', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },
    sword: { char: '✕', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
    globe: { char: '◎', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    briefcase: { char: '▢', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },
    brain: { char: '◎', bg: 'linear-gradient(135deg, #a18cd1, #fbc2eb)' },

    // ===== Currency =====
    money: { char: '¤', bg: 'linear-gradient(135deg, #ffd93d, #f6d365)' },
    diamond: { char: '◇', bg: 'linear-gradient(135deg, #89f7fe, #66a6ff)' },
    crown: { char: '♛', bg: 'linear-gradient(135deg, #ffd93d, #f6d365)' },

    // ===== Video / Voice Call =====
    video_call: { char: '▷', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    voice_call: { char: '☎', bg: 'linear-gradient(135deg, #43e97b, #38f9d7)' },
    mute: { char: '◇', bg: 'linear-gradient(135deg, #636e72, #b2bec3)' },
    speaker: { char: '◉', bg: 'linear-gradient(135deg, #4facfe, #00f2fe)' },
    hangup: { char: '✕', bg: 'linear-gradient(135deg, #ff6b6b, #ee5a24)' },
}

export function getContentIcon(key: string): GradientIconDef {
    return contentIcons[key] || { char: '?', bg: 'linear-gradient(135deg, #bdc3c7, #95a5a6)' }
}
