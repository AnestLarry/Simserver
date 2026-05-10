let _currentLocale = $state('en')
let _messages = $state<Record<string, string>>({})

export function t(key: string, params?: Record<string, string | number>): string {
  let msg = _messages[key] ?? key
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      msg = msg.replace(`{${k}}`, String(v))
    }
  }
  return msg
}

export async function setLocale(locale: string): Promise<void> {
  _currentLocale = locale
  const basePath = window.location.pathname.replace(/\/[^/]*$/, '/')
  const resp = await fetch(`${basePath}locales/${locale}.json`)
  _messages = await resp.json()
  localStorage.setItem('locale', locale)
}

export function getCurrentLocale(): string {
  return _currentLocale
}

export async function init(): Promise<void> {
  const saved = localStorage.getItem('locale')
  if (saved) {
    await setLocale(saved)
    return
  }
  const detected = navigator.language.startsWith('zh') ? 'zh-CN' : 'en'
  await setLocale(detected)
}
