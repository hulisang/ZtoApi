# SEO 和国际化 (i18n) 功能说明

本模板已集成 SEO 优化和多语言国际化支持。

## 功能特性

### 🔍 SEO 优化

- **Meta 标签优化**：自动生成 title、description、keywords
- **Open Graph**：支持社交媒体分享优化
- **Twitter Card**：Twitter 分享卡片优化
- **结构化数据**：JSON-LD 格式的 Schema.org 数据
- **Canonical URL**：规范链接支持
- **Robots 控制**：搜索引擎爬虫指令

### 🌍 国际化支持

- **多语言**：支持中文（zh-CN）、英文（en-US）、日文（ja-JP）
- **自动检测**：根据浏览器语言自动选择
- **手动切换**：页面右上角语言切换器
- **URL 参数**：支持 `?lang=zh-CN` 参数
- **全站翻译**：所有页面文本完全国际化

## 文件说明

### 核心文件

```
template/lib/
├── types.ts       # 类型定义（ProxyConfig、I18nTranslations）
├── i18n.ts        # 国际化翻译文件
├── seo.ts         # SEO 工具函数
└── pages.ts       # 页面生成函数（集成 SEO 和 i18n）
```

## 配置说明

### 1. SEO 配置

在 `main.ts` 中配置 SEO 参数：

```typescript
const CONFIG: ProxyConfig = {
  // ... 其他配置

  // SEO 配置（可选）
  seoTitle: "Your Service Name",
  seoDescription: "Your service description for search engines",
  seoKeywords: "keyword1,keyword2,keyword3",
  seoAuthor: "Your Name or Organization",
  seoOgImage: "https://your-domain.com/og-image.png", // Open Graph 图片
};
```

### 2. 使用 SEO Meta 标签

页面会自动包含以下 SEO 元素：

```html
<!-- 基础 SEO -->
<meta name="description" content="...">
<meta name="keywords" content="...">
<meta name="author" content="...">

<!-- Open Graph -->
<meta property="og:title" content="...">
<meta property="og:description" content="...">
<meta property="og:image" content="...">

<!-- Twitter Card -->
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="...">

<!-- 结构化数据 -->
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  ...
}
</script>
```

## 国际化使用

### 1. 支持的语言

- `zh-CN`：简体中文（默认）
- `en-US`：English
- `ja-JP`：日本語

### 2. 语言检测优先级

1. URL 参数：`?lang=en-US`
2. 浏览器 Accept-Language header
3. 默认语言：中文（zh-CN）

### 3. 添加新翻译

在 `lib/i18n.ts` 中添加新的翻译字段：

```typescript
export interface I18nTranslations {
  // ... 现有字段

  // 添加新字段
  newField: string;
}

// 然后在每个语言中添加翻译
const zhCN: I18nTranslations = {
  // ...
  newField: "新字段",
};

const enUS: I18nTranslations = {
  // ...
  newField: "New Field",
};

const jaJP: I18nTranslations = {
  // ...
  newField: "新しいフィールド",
};
```

### 4. 在页面中使用翻译

```typescript
import { getTranslations } from "./lib/i18n.ts";

const lang = detectLanguage(req);
const t = getTranslations(lang);

// 使用翻译
const html = `
  <h1>${t.homeTitle}</h1>
  <p>${t.homeSubtitle}</p>
  <button>${t.playgroundSendBtn}</button>
`;
```

## 语言切换器

每个页面右上角都有语言切换器：

```html
<div class="lang-switcher">
    <select id="langSelect" onchange="window.location.href='?lang='+this.value">
        <option value="zh-CN">🇨🇳 中文</option>
        <option value="en-US">🇺🇸 English</option>
        <option value="ja-JP">🇯🇵 日本語</option>
    </select>
</div>
```

## API 函数

### SEO 函数

```typescript
// 生成 SEO meta 标签
getSeoMeta(
  config: ProxyConfig,
  pageTitle: string,
  pageDescription?: string,
  currentUrl?: string
): string

// 生成结构化数据
getStructuredData(
  config: ProxyConfig,
  currentUrl?: string
): string
```

### i18n 函数

```typescript
// 获取翻译
getTranslations(lang: Language): I18nTranslations

// 检测语言
detectLanguage(req: Request): Language

// 从 URL 获取语言
getLanguageFromUrl(url: URL): Language | null
```

## 完整示例

创建带 SEO 和 i18n 的页面：

```typescript
async function handler(req: Request): Promise<Response> {
  const url = new URL(req.url);

  // 检测语言
  const urlLang = getLanguageFromUrl(url);
  const browserLang = detectLanguage(req);
  const lang = urlLang || browserLang;

  // 获取翻译
  const t = getTranslations(lang);

  if (url.pathname === "/") {
    return new Response(
      getHomePage(CONFIG, lang, url.toString()),
      { headers: { "Content-Type": "text/html" } }
    );
  }

  // ...
}
```

## 最佳实践

### SEO

1. **设置有意义的 title**：每个页面应有唯一的标题
2. **编写吸引人的 description**：150-160 字符，包含关键词
3. **选择相关的 keywords**：5-10 个相关关键词
4. **提供 OG 图片**：建议尺寸 1200x630px
5. **使用规范 URL**：避免重复内容问题

### i18n

1. **保持翻译一致性**：使用统一的术语
2. **考虑文化差异**：不同地区的表达习惯
3. **测试所有语言**：确保布局在各语言下正常
4. **提供语言切换**：让用户方便地切换语言
5. **默认语言选择**：根据目标用户选择合适的默认语言

## 扩展支持

### 添加新语言

1. 在 `types.ts` 中添加新的语言类型：
```typescript
export type Language = "zh-CN" | "en-US" | "ja-JP" | "fr-FR"; // 添加法语
```

2. 在 `i18n.ts` 中添加翻译：
```typescript
const frFR: I18nTranslations = {
  home: "Accueil",
  // ... 其他翻译
};

const translations: Record<Language, I18nTranslations> = {
  "zh-CN": zhCN,
  "en-US": enUS,
  "ja-JP": jaJP,
  "fr-FR": frFR, // 添加法语
};
```

3. 在语言切换器中添加选项：
```html
<option value="fr-FR">🇫🇷 Français</option>
```

## 注意事项

1. **SEO 配置是可选的**：如果不设置，会使用默认值
2. **语言切换会刷新页面**：保持 URL 参数状态
3. **翻译完整性**：确保所有语言的字段都有对应翻译
4. **性能影响**：SEO meta 和 i18n 对性能影响很小

## 测试

### 测试 SEO

1. 使用 Google Search Console 测试
2. 查看源代码确认 meta 标签
3. 使用 Facebook Debugger 测试 OG 标签
4. 使用 Twitter Card Validator 测试 Twitter Card

### 测试 i18n

1. 测试 URL 参数：`?lang=en-US`
2. 修改浏览器语言设置
3. 检查每个页面的翻译
4. 测试语言切换器

## 相关资源

- [Open Graph Protocol](https://ogp.me/)
- [Twitter Cards](https://developer.twitter.com/en/docs/twitter-for-websites/cards)
- [Schema.org](https://schema.org/)
- [BCP 47 语言标签](https://tools.ietf.org/html/bcp47)
