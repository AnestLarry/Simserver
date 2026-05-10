import { describe, it, expect, vi, beforeEach } from 'vitest'

const enMessages = {
  'app.title': 'File Manager',
  'common.download': 'Download',
  'chat.messages': 'Messages ({count})',
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn((url: string) => {
    if (url.endsWith('/en.json')) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve(enMessages) })
    }
    if (url.endsWith('/zh-CN.json')) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve({ 'app.title': '\u6587\u4EF6\u7BA1\u7406\u5668' }) })
    }
    return Promise.reject(new Error('not found'))
  }))
  const store: Record<string, string> = {}
  vi.stubGlobal('localStorage', {
    getItem: vi.fn((k: string) => store[k] ?? null),
    setItem: vi.fn((k: string, v: string) => { store[k] = v }),
    removeItem: vi.fn((k: string) => { delete store[k] }),
    clear: vi.fn(() => { for (const k in store) delete store[k] }),
    get length() { return Object.keys(store).length },
    key: vi.fn((i: number) => Object.keys(store)[i] ?? null),
  })
  vi.stubGlobal('navigator', { language: 'en-US' })
  vi.stubGlobal('window', { location: { pathname: '/view/Panel/' } })
})

it('t() returns key when no locale loaded', async () => {
  const { t } = await import('./i18n.svelte')
  expect(t('app.title')).toBe('app.title')
})

it('t() returns translated value after setLocale', async () => {
  const { t, setLocale, getCurrentLocale } = await import('./i18n.svelte')
  await setLocale('en')
  expect(t('app.title')).toBe('File Manager')
  expect(getCurrentLocale()).toBe('en')
  await setLocale('zh-CN')
  expect(t('app.title')).toBe('\u6587\u4EF6\u7BA1\u7406\u5668')
  expect(getCurrentLocale()).toBe('zh-CN')
})

it('t() falls back to key for missing translations', async () => {
  const { t, setLocale } = await import('./i18n.svelte')
  await setLocale('en')
  expect(t('nonexistent.key')).toBe('nonexistent.key')
})

it('t() interpolates {count} params', async () => {
  const { t, setLocale } = await import('./i18n.svelte')
  await setLocale('en')
  expect(t('chat.messages', { count: 3 })).toBe('Messages (3)')
})

it('setLocale saves to localStorage', async () => {
  const { setLocale } = await import('./i18n.svelte')
  await setLocale('en')
  expect(localStorage.setItem).toHaveBeenCalledWith('locale', 'en')
})

it('init reads from localStorage', async () => {
  localStorage.setItem('locale', 'en')
  const { t, init } = await import('./i18n.svelte')
  await init()
  expect(t('app.title')).toBe('File Manager')
})
