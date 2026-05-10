# i18n Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add on-demand file-based i18n (English + Simplified Chinese) to the Svelte 5 frontend.

**Architecture:** Custom Svelte 5 runes module (`i18n.svelte.ts`) using `$state` for reactive translations. JSON locale files fetched on demand — only one language file in memory at any time. Language switcher in App.svelte nav bar.

**Tech Stack:** Svelte 5 (runes), TypeScript, Vitest (TDD), Tailwind CSS, Python + Playwright

**Serving path:** Built frontend at `view/Panel/` served by Go backend at `/view/Panel/`.

---

### Task 1: Set Up Vitest + Write i18n Store Tests (TDD)

**Files:**
- Modify: `svelte/panel/vite.config.ts` (add `/// <reference types="vitest" />` and empty test config)
- Create: `svelte/panel/src/lib/i18n/i18n.test.ts`
- Modify: `svelte/panel/package.json` (add vitest devDependency)

- [ ] **Step 1: Add vitest devDependency**

Run in `svelte/panel/`:
```
npm install -D vitest
```

- [ ] **Step 2: Update vite.config.ts to include test config**

```ts
/// <reference types="vitest" />
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'path'

export default defineConfig({
  plugins: [
    svelte(),
    tailwindcss()
  ],
  resolve: {
    alias: {
      $lib: resolve('./src/lib')
    }
  },
  test: {
    environment: 'node',
    include: ['src/**/*.test.ts'],
  }
})
```

- [ ] **Step 3: Add test script to package.json**

Add to `scripts`:
```json
"test": "vitest run",
"test:watch": "vitest"
```

- [ ] **Step 4: Write the failing tests**

Create `svelte/panel/src/lib/i18n/i18n.test.ts`:

```ts
import { describe, it, expect, vi, beforeEach } from 'vitest'

const enMessages = {
  'app.title': 'File Manager',
  'common.download': 'Download',
  'common.cancel': 'Cancel',
  'common.back': 'Back',
  'common.view': 'View',
  'common.copy': 'Copy',
  'common.reset': 'Reset',
  'common.apply': 'Apply',
  'common.applied': 'Applied',
  'common.kb': ' KB',
  'common.mb': ' MB',
  'common.gb': ' GB',
  'app.hidePanel': 'Hide Panel',
  'app.showPanel': 'Show Panel',
  'app.pageMode': 'Page Mode:',
  'app.listView': '📁 List View',
  'app.photoGallery': '🖼️ Photo Gallery',
  'app.chatBoard': '💬 Chat Board',
  'app.uploadFiles': '📤 Upload Files',
  'app.sortBy': 'Sort By',
  'app.sortNameAZ': 'Name (A-Z)',
  'app.sortNameZA': 'Name (Z-A)',
  'app.sortDateOldest': 'Date (Oldest)',
  'app.sortDateNewest': 'Date (Newest)',
  'app.photoMode': 'Photo Mode',
  'app.showAllImages': 'Show All Images',
  'app.includeSubfolders': 'Include Subfolders',
  'app.imageList': 'Image List',
  'app.imageSizeControls': 'Image Size Controls',
  'app.width': 'Width:',
  'app.height': 'Height:',
  'app.language': 'Language',
  'fileList.filterByName': 'Filter by name',
  'fileList.searchPlaceholder': 'Search files and folders...',
  'fileList.columns': 'Columns',
  'fileList.oneColumn': '1 Column',
  'fileList.folders': 'Folders',
  'fileList.files': 'Files',
  'photo.noImagesFound': 'No Images Found',
  'photo.noImagesDescription': 'No images were found in this directory.',
  'photo.possibleReasons': 'Possible reasons:',
  'photo.reasonLsDisabled': 'The ls feature is disabled in the backend',
  'photo.reasonNoImages': 'No image files exist in this directory',
  'chat.title': 'Chat Board',
  'chat.status': 'Status:',
  'chat.connected': 'Connected',
  'chat.disconnected': 'Disconnected',
  'chat.username': 'Username',
  'chat.usernamePlaceholder': 'Enter your username',
  'chat.message': 'Message',
  'chat.messagePlaceholder': 'Type your message here...',
  'chat.send': 'Send',
  'chat.clear': 'Clear',
  'chat.messages': 'Messages ({count})',
  'upload.fileUpload': 'File Upload',
  'upload.chooseFile': 'Choose File to Upload',
  'upload.selected': 'Selected:',
  'upload.progress': 'Upload Progress',
  'upload.uploadFile': 'Upload File',
  'upload.textUpload': 'Text Upload',
  'upload.enterText': 'Enter your text content:',
  'upload.textPlaceholder': 'Type or paste your text here...',
  'upload.uploadText': 'Upload Text',
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn((url: string) => {
    if (url.endsWith('/en.json')) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve(enMessages) })
    }
    if (url.endsWith('/zh-CN.json')) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve({
        'app.title': '文件管理器',
        'common.download': '下载',
      }) })
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
})

describe('i18n', () => {
  it('t() returns key when messages are empty', async () => {
    const { t } = await import('./i18n.svelte')
    // before setLocale, messages should be empty
    expect(t('app.title')).toBe('app.title')
  })

  it('t() returns translated string after setLocale', async () => {
    const { t, setLocale } = await import('./i18n.svelte')
    await setLocale('en')
    expect(t('app.title')).toBe('File Manager')
    expect(t('common.download')).toBe('Download')
  })

  it('t() returns key for missing translations (graceful fallback)', async () => {
    const { t, setLocale } = await import('./i18n.svelte')
    await setLocale('en')
    expect(t('nonexistent.key')).toBe('nonexistent.key')
  })

  it('t() interpolates params', async () => {
    const { t, setLocale } = await import('./i18n.svelte')
    await setLocale('en')
    expect(t('chat.messages', { count: 5 })).toBe('Messages (5)')
  })

  it('setLocale switches locale', async () => {
    const { t, setLocale, getCurrentLocale } = await import('./i18n.svelte')
    await setLocale('en')
    expect(getCurrentLocale()).toBe('en')
    expect(t('app.title')).toBe('File Manager')
  })

  it('init() reads from localStorage first', async () => {
    const { t, init } = await import('./i18n.svelte')
    localStorage.setItem('locale', 'en')
    await init()
    expect(t('app.title')).toBe('File Manager')
  })

  it('init() detects zh from navigator.language', async () => {
    navigator = { language: 'zh-CN' } as any
    const { t, init } = await import('./i18n.svelte')
    await init()
    expect(t('app.title')).toBe('文件管理器')
  })
})
```

- [ ] **Step 5: Run tests to verify they fail**

Run: `npm test`
Expected: FAIL with import errors (module doesn't exist yet)

### Task 2: Implement i18n Store + Locale Files

**Files:**
- Create: `svelte/panel/src/lib/i18n/i18n.svelte.ts`
- Create: `svelte/panel/public/locales/en.json`
- Create: `svelte/panel/public/locales/zh-CN.json`

- [ ] **Step 1: Create i18n.svelte.ts**

```ts
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
  const resp = await fetch(`/locales/${locale}.json`)
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
```

- [ ] **Step 2: Create en.json**

`public/locales/en.json`:
```json
{
  "app.title": "File Manager",
  "app.hidePanel": "Hide Panel",
  "app.showPanel": "Show Panel",
  "app.pageMode": "Page Mode:",
  "app.listView": "📁 List View",
  "app.photoGallery": "🖼️ Photo Gallery",
  "app.chatBoard": "💬 Chat Board",
  "app.uploadFiles": "📤 Upload Files",
  "app.sortBy": "Sort By",
  "app.sortNameAZ": "Name (A-Z)",
  "app.sortNameZA": "Name (Z-A)",
  "app.sortDateOldest": "Date (Oldest)",
  "app.sortDateNewest": "Date (Newest)",
  "app.photoMode": "Photo Mode",
  "app.showAllImages": "Show All Images",
  "app.includeSubfolders": "Include Subfolders",
  "app.imageList": "Image List",
  "app.imageSizeControls": "Image Size Controls",
  "app.width": "Width:",
  "app.height": "Height:",
  "common.back": "Back",
  "common.download": "Download",
  "common.view": "View",
  "common.copy": "Copy",
  "common.reset": "Reset",
  "common.apply": "Apply",
  "common.applied": "Applied",
  "common.kb": " KB",
  "common.mb": " MB",
  "common.gb": " GB",
  "fileList.filterByName": "Filter by name",
  "fileList.searchPlaceholder": "Search files and folders...",
  "fileList.columns": "Columns",
  "fileList.oneColumn": "1 Column",
  "fileList.folders": "Folders",
  "fileList.files": "Files",
  "photo.noImagesFound": "No Images Found",
  "photo.noImagesDescription": "No images were found in this directory.",
  "photo.possibleReasons": "Possible reasons:",
  "photo.reasonLsDisabled": "The ls feature is disabled in the backend",
  "photo.reasonNoImages": "No image files exist in this directory",
  "chat.title": "Chat Board",
  "chat.status": "Status:",
  "chat.connected": "Connected",
  "chat.disconnected": "Disconnected",
  "chat.username": "Username",
  "chat.usernamePlaceholder": "Enter your username",
  "chat.message": "Message",
  "chat.messagePlaceholder": "Type your message here...",
  "chat.send": "Send",
  "chat.clear": "Clear",
  "chat.messages": "Messages ({count})",
  "upload.fileUpload": "File Upload",
  "upload.chooseFile": "Choose File to Upload",
  "upload.selected": "Selected:",
  "upload.progress": "Upload Progress",
  "upload.uploadFile": "Upload File",
  "upload.textUpload": "Text Upload",
  "upload.enterText": "Enter your text content:",
  "upload.textPlaceholder": "Type or paste your text here...",
  "upload.uploadText": "Upload Text",
  "app.language": "Language"
}
```

- [ ] **Step 3: Create zh-CN.json**

`public/locales/zh-CN.json`:
```json
{
  "app.title": "文件管理器",
  "app.hidePanel": "隐藏面板",
  "app.showPanel": "显示面板",
  "app.pageMode": "页面模式：",
  "app.listView": "📁 列表视图",
  "app.photoGallery": "🖼️ 图片浏览",
  "app.chatBoard": "💬 聊天室",
  "app.uploadFiles": "📤 上传文件",
  "app.sortBy": "排序方式",
  "app.sortNameAZ": "名称 (A-Z)",
  "app.sortNameZA": "名称 (Z-A)",
  "app.sortDateOldest": "日期 (从旧到新)",
  "app.sortDateNewest": "日期 (从新到旧)",
  "app.photoMode": "图片模式",
  "app.showAllImages": "显示全部图片",
  "app.includeSubfolders": "包含子文件夹",
  "app.imageList": "图片列表",
  "app.imageSizeControls": "图片尺寸控制",
  "app.width": "宽度：",
  "app.height": "高度：",
  "common.back": "返回",
  "common.download": "下载",
  "common.view": "查看",
  "common.copy": "复制",
  "common.reset": "重置",
  "common.apply": "应用",
  "common.applied": "已应用",
  "common.kb": " KB",
  "common.mb": " MB",
  "common.gb": " GB",
  "fileList.filterByName": "按名称筛选",
  "fileList.searchPlaceholder": "搜索文件和文件夹...",
  "fileList.columns": "列数",
  "fileList.oneColumn": "1 列",
  "fileList.folders": "文件夹",
  "fileList.files": "文件",
  "photo.noImagesFound": "未找到图片",
  "photo.noImagesDescription": "此目录中未找到任何图片。",
  "photo.possibleReasons": "可能的原因：",
  "photo.reasonLsDisabled": "后端未启用 ls 功能",
  "photo.reasonNoImages": "此目录中没有图片文件",
  "chat.title": "聊天室",
  "chat.status": "状态：",
  "chat.connected": "已连接",
  "chat.disconnected": "未连接",
  "chat.username": "用户名",
  "chat.usernamePlaceholder": "请输入用户名",
  "chat.message": "消息",
  "chat.messagePlaceholder": "在此输入消息...",
  "chat.send": "发送",
  "chat.clear": "清空",
  "chat.messages": "消息 ({count})",
  "upload.fileUpload": "文件上传",
  "upload.chooseFile": "选择要上传的文件",
  "upload.selected": "已选择：",
  "upload.progress": "上传进度",
  "upload.uploadFile": "上传文件",
  "upload.textUpload": "文本上传",
  "upload.enterText": "输入文本内容：",
  "upload.textPlaceholder": "在此输入或粘贴文本...",
  "upload.uploadText": "上传文本",
  "app.language": "语言"
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npm test`
Expected: All tests PASS

### Task 3: Apply i18n to All Svelte Components

**Files:**
- Modify: `svelte/panel/src/App.svelte` (import i18n, replace strings, add language switcher)
- Modify: `svelte/panel/src/lib/FileListPanel.svelte`
- Modify: `svelte/panel/src/lib/PhotoPanel.svelte`
- Modify: `svelte/panel/src/lib/ChatBoard.svelte`
- Modify: `svelte/panel/src/lib/Upload.svelte`
- Modify: `svelte/panel/src/lib/FileCardList.svelte`
- Modify: `svelte/panel/src/main.ts` (call `init()` on startup)

- [ ] **Step 1: Update main.ts to initialize i18n**

```ts
import "./app.pcss";
import "./app.css";
import App from "./App.svelte";
import { mount } from "svelte";
import { init } from "$lib/i18n/i18n.svelte";

init();

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;
```

- [ ] **Step 2: Apply i18n to App.svelte**

Add import at top:
```ts
import { t, setLocale, getCurrentLocale, init } from '$lib/i18n/i18n.svelte'
```

Replace all hardcoded English strings with `t('key')` calls:
- `"File Manager"` → `{t('app.title')}`
- `"Hide Panel"` / `"Show Panel"` → `{visible ? t('app.hidePanel') : t('app.showPanel')}`
- etc.

Add language switcher button near the navigation area:
```svelte
<button onclick={() => setLocale(getCurrentLocale() === 'en' ? 'zh-CN' : 'en')}>
  {getCurrentLocale() === 'en' ? '中' : 'EN'}
</button>
```

- [ ] **Step 3: Apply i18n to FileListPanel.svelte**

Replace all hardcoded strings with `t()` calls.

- [ ] **Step 4: Apply i18n to FileCardList.svelte**

Replace all hardcoded strings with `t()` calls.

- [ ] **Step 5: Apply i18n to PhotoPanel.svelte**

Replace all hardcoded strings with `t()` calls.

- [ ] **Step 6: Apply i18n to ChatBoard.svelte**

Replace all hardcoded strings with `t()` calls. Handle `{count}` param in messages count.

- [ ] **Step 7: Apply i18n to Upload.svelte**

Replace all hardcoded strings with `t()` calls.

### Task 4: Build Verification

**Files:**
- Build output: `svelte/panel/dist/`
- Copy to: `view/Panel/`

- [ ] **Step 1: Build the Svelte frontend**

```bash
cd svelte/panel && npm run build
```
Expected: Build succeeds, output in `svelte/panel/dist/`

- [ ] **Step 2: Copy dist to view/Panel/**

```bash
cp -r svelte/panel/dist/* view/Panel/
```
Expected: `view/Panel/` now contains updated `index.html`, `index.js`, `index.css` plus `locales/` directory

- [ ] **Step 3: Build the Go binary**

```bash
go build -o simserver.exe
```
Expected: Build succeeds, simserver.exe created

### Task 5: Python Playwright Verification

**Files:**
- Create: `scripts/check_i18n.py`

- [ ] **Step 1: Install Playwright for Python**

```bash
pip install playwright
python -m playwright install chromium
```

- [ ] **Step 2: Create check_i18n.py**

The script should:
1. Start the Go server (simserver.exe)
2. Wait for it to be ready on port 5000
3. Open the page at `http://localhost:5000/view/Panel/`
4. Check that the page loads with Chinese locale
5. Verify key translated elements are present
6. Switch locale and verify English
7. Take screenshots
8. Navigate to each page/panel view and check for untranslated strings
9. Report any hardcoded English strings that were missed

```python
import subprocess
import time
import sys
from pathlib import Path
from playwright.sync_api import sync_playwright

def main():
    server_dir = Path(__file__).parent.parent
    
    # Modify config to enable view
    import json
    config_path = server_dir / 'config.json'
    with open(config_path) as f:
        config = json.load(f)
    config['view']['enable'] = True
    config['port'] = '5000'
    with open(config_path, 'w') as f:
        json.dump(config, f, indent=4)
    
    # Start server
    server = subprocess.Popen(
        [str(server_dir / 'simserver.exe')],
        cwd=server_dir,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    time.sleep(2)
    
    try:
        with sync_playwright() as p:
            browser = p.chromium.launch(headless=True)
            page = browser.new_page()
            
            # Set Chinese locale
            page.context.add_init_script("""
                Object.defineProperty(navigator, 'language', {
                    get: () => 'zh-CN'
                });
                localStorage.setItem('locale', 'zh-CN');
            """)
            
            page.goto('http://localhost:5000/view/Panel/', wait_until='networkidle')
            page.wait_for_timeout(1000)
            
            # Take screenshot
            page.screenshot(path='zh-CN-page.png', full_page=True)
            
            # Check for common English strings that should NOT appear
            english_phrases = [
                'File Manager', 'Hide Panel', 'Show Panel', 'Download',
                'Upload Files', 'Chat Board', 'Photo Gallery', 'List View',
                'Sort By', 'Filter by name', 'No Images Found',
                'File Upload', 'Text Upload', 'Username', 'Message',
            ]
            
            body_text = page.inner_text('body')
            found_english = [p for p in english_phrases if p in body_text]
            
            if found_english:
                print(f"WARNING: Found untranslated English strings: {found_english}")
            else:
                print("OK: No untranslated English strings found in zh-CN mode")
            
            # Switch to English
            page.evaluate("localStorage.setItem('locale', 'en')")
            page.reload(wait_until='networkidle')
            page.wait_for_timeout(1000)
            page.screenshot(path='en-page.png', full_page=True)
            
            # Verify English text is visible
            en_body = page.inner_text('body')
            if 'File Manager' in en_body:
                print("OK: English mode shows English text")
            else:
                print("WARNING: English mode may not have loaded correctly")
            
            browser.close()
            print("Playwright verification complete")
            
    finally:
        server.terminate()
        server.wait()
        # Restore config
        config['view']['enable'] = False
        with open(config_path, 'w') as f:
            json.dump(config, f, indent=4)

if __name__ == '__main__':
    main()
```

- [ ] **Step 3: Run the verification script**

```bash
python scripts/check_i18n.py
```
Expected: All checks pass, no untranslated English strings when zh-CN is active

- [ ] **Step 4: Review screenshots for missed translations**

Check `zh-CN-page.png` and `en-page.png` for any remaining hardcoded English strings that weren't caught by the automated checks.
