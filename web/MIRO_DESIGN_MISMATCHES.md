# Miro 设计不一致清单

生成日期：2026-05-12

修补状态：已按本清单完成第一轮 Miro 对齐修补。

- 已修补：全局字体 fallback、浅色主题语义、primary CTA / dark mode 映射、Miro 半径与阴影 token、营销容器、display typography、按钮、输入框、select、badge、表格、dialog/dropdown/sheet/popover/tooltip 等基础组件。
- 已修补：营销首页 nav、hero、白板 mockup、logo wall、feature cards、pricing monthly/yearly toggle、dark CTA banner、6 列 footer、app-store / review badge 文案及中英文 i18n。
- 已修补：Auth / Console / Project 顶层壳、分享页、邀请页等明显的 shadow、scale、过重字重与非 Miro 圆角残留。
- 保留说明：项目工作台内部仍有若干 `rounded-2xl` / `font-semibold` 等业务面板局部样式。这些属于复杂产品操作界面，`DESIGN.md` 主要覆盖 Miro marketing/pricing/footer/基础组件规范；本轮只收敛明显冲突，不做无设计依据的大规模工作台重排。
- 验证：`pnpm type-check`、`pnpm test:i18n`、`pnpm test -- src/test/marketing-home-page.test.tsx` 均已通过；本地首页已在 `http://localhost:3001` 打开，导航、hero、pricing、footer DOM 检查通过且控制台无错误。

范围：`web/src/app`、`web/src/components`、`web/src/themes`、`web/src/providers`、营销文案结构。已忽略构建产物 `web/.next`。

依据：`web/DESIGN.md` 中的 Miro 规范。重点规则包括：Roobert PRO 全局字体、白色 canvas、黑色胶囊 CTA、Miro 黄仅用于 wordmark / promo / yellow tag、1280px 营销容器、真实白板 mockup、扁平卡片、特定 pricing / footer / nav / form / tabs / table 组件规范。

## 总览

- [ ] 全局字体仍是 Inter + JetBrains Mono，不是 Miro 的 Roobert PRO。
- [ ] 全局 typography helper 仍叫 `figma-*`，并且多处使用 mono / uppercase / 0 letter-spacing，和 Miro 字体规范不一致。
- [ ] 主题 token 有 Miro 色板，但部分语义映射、dark mode、半径和阴影与 `DESIGN.md` 不一致。
- [ ] 营销首页接近 Miro 色系，但 hero、mockup、pricing、footer、logo wall、CTA banner 仍偏 SaaS dashboard，而不是 Miro marketing surface。
- [ ] 基础 UI 组件保留大量 shadcn 默认行为：hover、spring、scale、shadow、destructive red、非 pill filter/select 等。
- [ ] Console / Project / Auth 页面大量局部 Tailwind class 覆盖 token，导致圆角、阴影、字体权重和黄色使用不统一。

## 1. 全局 Token 与基础样式

- [ ] **字体栈不符合 Miro：应使用 Roobert PRO / Noto Sans fallback，当前使用 Inter + JetBrains Mono。**
  - 证据：`src/app/layout.tsx:2`、`src/app/layout.tsx:15`、`src/app/layout.tsx:19`
  - 证据：`src/app/globals.css:143`、`src/app/globals.css:146`
  - 影响：所有页面默认 `font-sans` 不是 Miro 的 Roobert PRO；`figma-caption` / `figma-eyebrow` 还切到了 mono。
  - 修复方向：定义 `--font-miro-sans: "Roobert PRO", "Noto Sans", -apple-system, BlinkMacSystemFont, sans-serif`，并让 `--font-sans`、caption、eyebrow 全部走该字体。代码块如确需 mono，应明确作为产品例外。

- [ ] **Display typography 的 letter-spacing 不符合规范。**
  - 规范：hero `80px / 500 / 1.05 / -2px`，display-lg `60px / 500 / 1.10 / -1.5px`，heading-1 `48px / -1px`。
  - 证据：`src/app/globals.css:344`、`src/app/globals.css:352` 目前 `letter-spacing: 0`。
  - 影响：营销大标题缺少 Miro 视觉里的紧凑 display 字距。
  - 修复方向：新增 Miro 命名 typography utilities，按 `DESIGN.md` 映射负字距。

- [ ] **Hero 响应式字号断点不符合 Miro。**
  - 规范：mobile small 36px，mobile large 48px，tablet 60px，desktop 64px，wide 80px。
  - 证据：`src/app/globals.css:344`、`src/app/globals.css:361`、`src/app/globals.css:371` 当前为 48px -> 62.4px -> 80px，1024px 就进入 80px。
  - 影响：中等桌面上 hero 过大，手机小屏上又没有降到 36px。
  - 修复方向：按 `DESIGN.md` breakpoint 重写 `.figma-display-xl` 或替换为 `.miro-hero-display`。

- [ ] **`figma-*` 命名和语义没有迁移到 Miro design system。**
  - 证据：`src/app/globals.css:344`、`src/app/globals.css:380`、`src/app/globals.css:410`、`src/app/globals.css:419`
  - 证据：营销和项目页大量调用，如 `src/components/features/site/home/hero-section.tsx:25`、`src/components/features/project/project-dashboard-page.tsx:342`
  - 影响：实现层仍是旧 Figma 命名，后续维护时很难直接对照 `DESIGN.md` 的 Miro token。
  - 修复方向：把 `figma-display-xl/body/caption/color-block` 迁移为 `miro-hero-display/body-md/caption/card-feature-*` 等命名。

- [ ] **Caption / eyebrow 使用 mono 字体，与 Miro “Roobert PRO across every UI surface” 冲突。**
  - 证据：`src/app/globals.css:410`、`src/app/globals.css:419`
  - 证据：`src/components/features/site/home/logo-cloud.tsx:22`、`src/components/features/site/home/marketing-footer.tsx:68`
  - 影响：小标签、footer heading、badge 文本不像 Miro 的 Roobert PRO。
  - 修复方向：caption、micro、micro-uppercase 全部改为 Roobert PRO；只给 API path/code/editor 区域保留 mono 例外。

- [ ] **浅色主题中 `--bg-surface` 映射成了 canvas 白，而 Miro surface 应是 `#f7f8fa`。**
  - 规范：`colors.surface = #f7f8fa`，`colors.canvas = #ffffff`。
  - 证据：`src/themes/light.css:34` 当前 `--bg-surface: var(--miro-canvas)`，`--bg-subtle` 才是 `--miro-surface`。
  - 影响：使用 `bg-bg-surface` 的控件视觉上变成白底，和 search-pill / subtle section 的 Miro 灰底语义错位。
  - 修复方向：让 `--bg-surface` 指向 `--miro-surface`，必要时新增 `--bg-canvas` 专门表达白底。

- [ ] **Dark mode 没有 Miro 规范支撑，且会把 primary CTA 改成黄色。**
  - 规范：`DESIGN.md` 明确 dark-mode token values 未提取；Miro 主要只有 dark footer / dark CTA banner，不是全站 dark app。
  - 证据：`src/themes/dark.css:6`、`src/themes/dark.css:18`、`src/themes/dark.css:69`
  - 证据：`src/app/layout.tsx:57` 启用 `defaultTheme="system"` 和 `enableSystem`。
  - 影响：系统暗色模式下 primary button 变为 Miro 黄，违反“黑色胶囊 CTA 为主”的品牌规则。
  - 修复方向：营销面强制 light；如果 app 需要 dark，单独补一套非 Miro marketing 的产品设计规范。

- [ ] **圆角 token 有冲突。**
  - 规范：xs 4px，sm 6px，md 8px，lg 12px，xl 16px，xxl 20px，xxxl 28px，feature 32px，full 9999px。
  - 证据：`src/themes/primitives.css:116` 定义 xs 为 4px，但 `src/app/globals.css:112` 又把 `--radius-xs` 改成 2px。
  - 证据：`src/app/globals.css:118` 把 `--radius-pill` 设为 `3.125rem`，不是 full 9999px。
  - 影响：Tailwind radius token 与 Miro radius scale 不完全一致。
  - 修复方向：保留一处权威 radius 定义，确保 Tailwind `rounded-*` 与 `DESIGN.md` 数值一致。

- [ ] **阴影 token 超出 Miro 规范，且标准卡片大量使用阴影。**
  - 规范：普通 card / documentation card flat，无阴影；白板 mockup 才用 `rgba(5,0,56,.08) 0 12px 32px -4px`；modal 用 `0 16px 48px -8px`。
  - 证据：`src/app/globals.css:129`、`src/app/globals.css:130`、`src/app/globals.css:131`、`src/app/globals.css:132`
  - 统计：`shadow-sm` / shadow token 在 app/components 中约 82 处；代表文件 `src/components/features/project/api-request-workbench.tsx`、`src/components/features/console/console-shell.tsx`。
  - 修复方向：标准 Card 只用 hairline border；把 `shadow-soft` 限定给 whiteboard mockup；dropdown/modal 使用对应 elevation token。

- [ ] **交互动效偏离 Miro 的 150-200ms ease。**
  - 证据：`src/app/globals.css:292`、`src/app/globals.css:309`、`src/app/globals.css:317`
  - 影响：按钮和可点击元素使用 `duration-300 ease-spring`、`active:scale-*`、brightness，这不是 `DESIGN.md` 的 pressed-state 风格。
  - 修复方向：按钮保留 pressed color state，去掉全局 spring/scale；把 transition 收敛到 150-200ms ease。

- [ ] **营销容器最大宽度和 gutter 不符合 Miro。**
  - 规范：marketing max-width 1280px，32px gutters。
  - 证据：`src/app/globals.css:208` 到 `src/app/globals.css:240` 中 `.container` 在 1536px 屏幕会扩到 1536px，默认 gutter 16px。
  - 影响：宽屏营销页比 Miro 规范更散，左右留白不足。
  - 修复方向：marketing shell 下的 container 固定 `max-width: 1280px; padding-inline: 32px`，移动端再缩小。

## 2. 基础 UI 组件

- [ ] **Button 尺寸和状态偏离 Miro button spec。**
  - 规范：button-primary padding `12px 24px`，`rounded.full`，typography `14px / 500 / 1.30`；pressed 改 `charcoal`。
  - 证据：`src/components/ui/button.tsx:14`、`src/components/ui/button.tsx:22`、`src/components/ui/button.tsx:32`
  - 当前问题：存在 `xl` / `2xl`，`2xl` 为 `h-13 px-10 text-base`；hover 用 `bg-primary/95`，active 用全局 scale/brightness。
  - 修复方向：为 Miro marketing CTA 固定 `h-auto px-6 py-3 text-sm`；删除 marketing 按钮上的 scale 和 hover alpha。

- [ ] **Button variant 数量和颜色超过 Miro spec。**
  - 规范：primary black、yellow、blue、secondary outline、on-dark、ghost、link、icon-circular。
  - 证据：`src/components/ui/button.tsx:19` 到 `src/components/ui/button.tsx:30`
  - 当前问题：`destructive`、`secondary` 等沿用 shadcn 语义；dark mode 下 default primary 会变黄。
  - 修复方向：拆出 `miroButtonVariants`，按 `DESIGN.md` 命名，不把 app destructive/action state 混进 marketing CTA。

- [ ] **Input focus state 没有 2px brand-blue border。**
  - 规范：text-input height 44px，border `1px hairline-strong`；focused border `2px solid brand-blue`。
  - 证据：`src/components/ui/input.tsx:19`、`src/components/ui/input.tsx:23`
  - 当前问题：focus 只是 `focus-visible:border-primary`，没有 2px 宽度；padding 是 `px-4 py-2`，不是 `12px 16px` 的明确 token。
  - 修复方向：focus-visible 改为 2px blue border，必要时用 `outline-offset` 避免布局跳动。

- [ ] **SearchInput 不是 Miro search-pill。**
  - 规范：search-pill `rounded.md`，height 40，surface background，steel text，hairline border。
  - 证据：`src/components/ui/input.tsx:164` 当前 `rounded-full`。
  - 修复方向：SearchInput 改成 rounded-md / h-10 / bg-surface / border-hairline。

- [ ] **Select / filter dropdown 不是 pill filter-dropdown。**
  - 规范：filter dropdown `rounded.full`，padding `8px 16px`，border hairline-strong。
  - 证据：`src/components/ui/select.tsx:35` 当前 `rounded-md`。
  - 影响：所有筛选器 / select 更像 shadcn form control，而不是 Miro pill dropdown。
  - 修复方向：提供 `variant="filter"` 或默认把 marketing filter select 改成 pill。

- [ ] **Dropdown / Popover / Dialog 阴影没有使用 Miro modal elevation。**
  - 规范：modal/dropdown Level 4 为 `rgba(5,0,56,.12) 0 16px 48px -8px`。
  - 证据：`src/components/ui/dropdown-menu.tsx:51`、`src/components/ui/select.tsx:74`、`src/components/ui/dialog.tsx:46` 都使用 `shadow-soft`。
  - 修复方向：增加 `--shadow-modal`，dropdown/dialog/sheet 用 Level 4，whiteboard mockup 才用 `shadow-soft`。

- [ ] **Table header 使用 uppercase semibold，不是 Miro comparison table row typography。**
  - 规范：comparison table body-sm，micro-uppercase 只用于 section dividers。
  - 证据：`src/components/ui/table.tsx:93`
  - 影响：普通数据表头都被强制 uppercase，和 Miro 的 dense comparison table 语义不一致。
  - 修复方向：普通表头用 `body-sm-medium`，只有 comparison section divider 用 micro-uppercase。

- [ ] **Badge default 是黑底，不是 Miro badge taxonomy。**
  - 规范：promo yellow、tag-yellow、tag-purple、tag-coral、success green、discount yellow-rect。
  - 证据：`src/components/ui/badge.tsx:15` 到 `src/components/ui/badge.tsx:29`
  - 当前问题：default 黑底 primary，secondary 使用 generic secondary，destructive 使用 solid red。
  - 修复方向：补齐 `promo/tagYellow/tagPurple/tagCoral/success/discount`，减少 generic default 用法。

- [ ] **Destructive / error color 没有使用 Miro red tokens。**
  - 规范：`brand-red = #fbd4d4`、`brand-red-dark = #e3c5c5`，更偏 soft red。
  - 证据：`src/themes/light.css:24`、`src/components/ui/badge.tsx:24`、`src/components/ui/input.tsx:27`
  - 当前问题：使用 `--red-500` solid red 和白字，视觉上偏 shadcn danger。
  - 修复方向：定义 Miro error background/border/text token，表单错误使用 soft red surface + dark readable foreground。

## 3. 营销首页

- [ ] **Marketing nav 结构不完全符合 Miro Top Navigation。**
  - 规范：sticky white bar，高约 64px；左侧 yellow Miro wordmark + Product / Solutions / Resources；右侧 Login / Pricing / Contact sales + black pill Get started free。
  - 证据：`src/components/features/site/home/marketing-navbar.tsx:103`、`src/components/features/site/home/marketing-navbar.tsx:118`、`src/components/features/site/home/marketing-navbar.tsx:128`
  - 当前问题：高度 68px；右侧是 LanguageSwitcher + Login + Sign Up，没有 Pricing / Contact sales；移动菜单弹层用了 `rounded-2xl` + `shadow-soft`。
  - 修复方向：按 Miro nav IA 和 CTA hierarchy 调整，语言切换降级到次要区域或 footer。

- [ ] **Hero 没有 Miro 的 centered hero-band 构图。**
  - 规范：centered headline、centered subtitle、centered button row、whiteboard mockup below。
  - 证据：`src/components/features/site/home/hero-section.tsx:19`、`src/components/features/site/home/hero-section.tsx:21`、`src/components/features/site/home/hero-section.tsx:28`
  - 当前问题：内容是左对齐 `max-w-5xl`，不是 `mx-auto text-center`；section padding `lg:py-28` 不等于 desktop hero 120px；CTA 使用 `size="2xl"` 过大。
  - 修复方向：hero 文案居中，desktop padding 使用 `spacing.hero`，CTA 回到 `12px 24px` pill。

- [ ] **Hero badge 使用 Miro 黄作为通用 highlight，容易违反黄色使用规则。**
  - 规范：brand yellow 保留给 wordmark、top promo banner、yellow tag chips；不是通用强调色。
  - 证据：`src/components/features/site/home/hero-section.tsx:25`
  - 证据：类似用法还出现在 `src/components/features/site/home/marketing-navbar.tsx:77`、`src/components/features/site/home/pricing-section.tsx:31`
  - 修复方向：如果是 tag chip，用 `badge-tag-yellow` 的 `surface-yellow + yellow-dark`；只有 promo banner badge 用 full brand yellow。

- [ ] **Hero mockup 不是 Miro whiteboard / board mockup。**
  - 规范：真实 Miro-board mockup， sticky notes、kanban、mind map 等作为视觉主体。
  - 证据：`src/components/features/site/home/product-preview-mockup.tsx:24` 到 `src/components/features/site/home/product-preview-mockup.tsx:158`
  - 当前问题：是三列 API dashboard panel + timeline + result card；更像 SaaS admin screenshot，不像白板画布。
  - 修复方向：重做为白板画布式 product mockup：白底 board frame、sticky-note palette、flow arrows、kanban/cards、subtle Level 3 mockup shadow。

- [ ] **Marketing mockup 使用 grid background，不是 Miro 的真实产品图 / board mockup。**
  - 证据：`src/app/globals.css:503` 的 `.marketing-grid`
  - 影响：视觉更像工程画布背景，而不是 `DESIGN.md` 要求的 real board mockup imagery。
  - 修复方向：用真实产品 UI 截图或手工构建 Miro-like board frame，避免纯 CSS 网格当主视觉。

- [ ] **Feature card hover 不符合 no-hover policy，且 padding 不是 token。**
  - 规范：hover states 不记录；card-feature-* padding `spacing.xxl` 32px，rounded.xxxl 28px。
  - 证据：`src/components/features/site/home/feature-cards.tsx:60`、`src/components/features/site/home/feature-cards.tsx:63`
  - 当前问题：`p-7` 是 28px，不是 32px；有 `hover:-translate-y-0.5 hover:shadow-sm`。
  - 修复方向：使用 `p-8` / `rounded-[28px]`，去掉 hover lift。

- [ ] **Logo cloud 不符合 Miro logo-wall-item。**
  - 规范：customer logo wordmarks inline，consistent 100px height，transparent background。
  - 证据：`src/components/features/site/home/logo-cloud.tsx:26` 到 `src/components/features/site/home/logo-cloud.tsx:31`
  - 当前问题：logo 被渲染成圆角灰色文字 pill，没有真实 wordmark，也没有 100px 统一高度。
  - 修复方向：使用真实或模拟 wordmark asset，透明背景，按固定高度排列。

- [ ] **Product story blocks 使用整块 pastel / navy panel，不完全符合 Miro feature card + mockup 组合。**
  - 证据：`src/components/features/site/home/product-story-section.tsx:34`
  - 当前问题：`figma-color-block` 是 28px / 32-40px padding 的大包裹容器；Miro 更强调白板 mockup + pastel feature cards 的组合。
  - 修复方向：把大块背景拆成 Miro feature cards / whiteboard mockup frame，减少整区块彩色容器。

- [ ] **Dark CTA banner 使用 ink-deep navy 且尺寸不符合 cta-banner-dark。**
  - 规范：`cta-banner-dark` background `primary #1c1c1e`，rounded.feature 32px，padding section 64px，centered。
  - 证据：`src/app/globals.css:558`、`src/components/features/site/home/final-cta.tsx:18`
  - 当前问题：`figma-color-block-navy` 背景是 `#050038`；radius 28px；padding 32/40px；内容左对齐，并有两个 CTA。
  - 修复方向：改为 `#1c1c1e`、32px radius、64px padding、centered headline + subtitle + white pill primary CTA。

- [ ] **Stats token 不符合 Miro stat-display。**
  - 规范：stat-display 64px / 500 / 1.10 / -1.5px。
  - 证据：`src/components/features/site/home/stats-section.tsx:34`
  - 当前问题：使用 `text-5xl` 48px，缺少负字距。
  - 修复方向：新增 `.miro-stat-display` 并应用到 stats value。

- [ ] **Pricing section 缺少 monthly/yearly pill toggle。**
  - 规范：`toggle-monthly-yearly` 是 pricing 关键组件。
  - 证据：`src/components/features/site/home/pricing-section.tsx:79` 到 `src/components/features/site/home/pricing-section.tsx:118` 没有 toggle。
  - 修复方向：在 pricing heading 和 cards 之间加 Monthly / Annual pill toggle。

- [ ] **Pricing card padding / radius / CTA hierarchy 不完全符合规范。**
  - 规范：pricing-card rounded.xl 16px，padding xxl 32px；featured lavender + 2px brand-blue；enterprise primary dark。
  - 证据：`src/components/features/site/home/pricing-section.tsx:13`、`src/components/features/site/home/pricing-section.tsx:16`
  - 当前问题：`rounded-2xl` 20px，`p-6` 24px；free/starter CTA 使用 outline，featured 使用 default，但没有统一对照 Miro pricing CTA。
  - 修复方向：改为 `rounded-xl p-8`，并按 tier 定义 button-primary / button-secondary / button-on-dark。

- [ ] **Comparison table 不够 Miro 式 dense pricing comparison。**
  - 规范：4-tier grid 后跟约 80 行 dense comparison table，有 section dividers。
  - 证据：`src/components/features/site/home/pricing-section.tsx:106`、`src/components/features/site/home/home-content.ts` pricing rows 仅少量核心项。
  - 当前问题：comparison table 太短，缺少 section divider / dense row rhythm。
  - 修复方向：扩展 comparisonRows 结构，加 category divider 和更多功能项；移动端保留 horizontal scroll。

- [ ] **Footer 不符合 Miro massive footer。**
  - 规范：6-column link grid，Product / Solutions / Tools / Resources / Company / Plans & Pricing；含 app-store badge / Capterra badge。
  - 证据：`src/components/features/site/home/marketing-footer.tsx:85`、`src/components/features/site/home/home-content.ts:38`
  - 当前问题：只有 5 列，列名是 Product / API Docs / Resources / Company / Legal；没有 app-store badge、Capterra badge，也没有 Plans & Pricing 列。
  - 修复方向：按 Miro footer IA 重组列，补充 badge 组件或明确删除相关 requirement。

## 4. Auth 页面

- [ ] **Auth 页面使用大面积黄色系 side panel，不是 Miro 规范里的 auth surface。**
  - 证据：`src/app/(auth)/layout.tsx:25`
  - 当前问题：`figma-color-block-lime` 覆盖整个左侧，且布局是登录页装饰面板；`DESIGN.md` 没有定义 auth page，若严格复用 Miro marketing，应避免把 yellow/pastel 当全页背景。
  - 修复方向：要么新增 auth page 设计规范，要么把它改为白 canvas + black CTA + 小面积 pastel feature cards。

- [ ] **Auth 卡片和 insight card 使用非规范阴影。**
  - 证据：`src/app/(auth)/layout.tsx:43`、`src/app/(auth)/layout.tsx:55`、`src/app/(auth)/layout.tsx:86`
  - 当前问题：多个普通卡片使用 `shadow-sm` / `shadow-soft`；Miro 普通 cards 应 flat，只用 hairline border。
  - 修复方向：普通 auth card 去阴影；如需 mockup 视觉，再使用 whiteboard mockup shadow。

- [ ] **Auth 表单 badge 使用 full brand yellow。**
  - 证据：`src/components/features/auth/login-form.tsx:56`、`src/components/features/auth/register-form.tsx:90`、`src/components/features/auth/forgot-password-form.tsx:35`
  - 当前问题：登录/注册小 badge 使用 `bg-highlight`，属于 generic emphasis。
  - 修复方向：改为 `badge-tag-yellow` 的 pale yellow surface，或去掉 badge。

## 5. Console / Project 应用页面

- [ ] **应用页大量使用 heading `font-semibold`，不符合 Miro heading 500 weight。**
  - 统计：app/components 中 `font-semibold` 约 98 处。
  - 代表证据：`src/components/features/project/project-management-page.tsx:189`、`src/components/features/project/category-management-page.tsx:470`、`src/components/features/project/api-spec-management-page.tsx:1997`、`src/components/features/console/account-settings.tsx:138`
  - 影响：页面标题和卡片标题偏 shadcn dashboard weight，不是 Miro 500 display/heading scale。
  - 修复方向：页面标题改用 Miro heading utilities，badge / micro label 才保留 600。

- [ ] **应用页大量使用普通卡片阴影。**
  - 代表证据：`src/components/features/project/flow-management-page.tsx:1288`、`src/components/features/project/api-request-workbench.tsx:4970`、`src/components/features/project/project-detail-page.tsx:780`、`src/components/features/console/dashboard-stats.tsx:95`
  - 影响：Miro 的 flat documentation card 规则被破坏；页面看起来更像 generic SaaS dashboard。
  - 修复方向：默认 Card 去 `shadow-sm`，只保留 border；hover 不加 lift。

- [ ] **应用页圆角使用非常不统一。**
  - 统计：`rounded-md/lg/xl/2xl/[1.75rem]` 在 app/components 中约 442 处。
  - 代表证据：`src/components/features/project/project-dashboard-page.tsx:352`、`src/components/features/project/project-workspace-page.tsx:1588`、`src/components/features/project/flow-management-page.tsx:3631`
  - 当前问题：同一类内容容器混用 8/12/16/20/28px；Miro 对 input/card/pricing/feature/cta 有明确 radius。
  - 修复方向：建立组件级 token，不在业务页直接写 `rounded-*`；按 input 8、card 16、feature 28、cta 32、pill full 区分。

- [ ] **Miro 黄 / highlight 在项目页被当成通用状态色。**
  - 证据：`src/components/features/project/project-topbar.tsx:97`、`src/components/features/project/flow-management-page.tsx:238`、`src/components/features/project/project-home-status.tsx:36`
  - 影响：黄色不再是 wordmark / promo / yellow tag 的品牌信号，而变成项目状态色。
  - 修复方向：状态类改用 neutral / success / warning semantic；yellow 仅保留给 Miro tag chip 或 promo。

- [ ] **Console / Project topbar 使用阴影和非 Miro nav 结构。**
  - 证据：`src/components/features/console/console-shell.tsx:91`、`src/components/features/project/project-topbar.tsx:74`
  - 当前问题：header 使用 `shadow-sm`，avatar/icon buttons 多处也有 shadow。
  - 修复方向：topbar 保持 flat border-bottom；icon button 使用 36px circular + hairline border，无阴影。

- [ ] **项目工作区导航和请求编辑器是产品 UI，但没有独立于 Miro marketing 的规范。**
  - 证据：`src/components/features/project/project-workspace-layout.tsx:87`、`src/components/features/project/api-request-workbench.tsx:5477`
  - 当前问题：它混合 Miro pastel、shadcn controls、API workbench dense UI；`DESIGN.md` 主要描述 marketing/pricing/customer surfaces，不足以约束这个复杂应用界面。
  - 修复方向：为 app/workbench 补一份 Miro-compatible product UI spec，明确 dense table、editor、sidebar、request method badge、code block 的例外规则。

- [ ] **代码/路径/ID 大量使用 mono，与 Miro 全局字体冲突但可能是产品必要例外。**
  - 代表证据：`src/components/features/project/api-spec-management-page.tsx:420`、`src/components/features/project/test-case-management-page.tsx:342`、`src/components/features/project/environment-management-page.tsx:477`
  - 判断：如果严格套 `DESIGN.md`，这些是 off-spec；如果 API 工具需要代码可读性，应在设计规范里新增 code typography token，而不是混用 JetBrains Mono。
  - 修复方向：补 `code-sm/code-md` token，并限定使用场景。

## 6. Share / 测试页面

- [ ] **Share API Spec 页面使用 mono 大标题，不符合 Miro heading。**
  - 证据：`src/app/share/api-spec/[slug]/page.tsx:150`
  - 当前问题：`h1` 使用 `font-mono text-2xl font-semibold`；Miro heading 应是 Roobert PRO、500 weight。
  - 修复方向：标题用 Miro heading，path/method/code 细节保留 code token。

- [ ] **Share 页面普通容器使用 28px 大卡和 whiteboard shadow。**
  - 证据：`src/app/share/api-spec/[slug]/page.tsx:136`
  - 当前问题：`rounded-[1.75rem]` + `shadow-soft` 用在普通 API spec shell 上；Miro 的 `shadow-soft` 应留给 whiteboard mockup。
  - 修复方向：改为 standard card：16px radius + hairline border + no shadow。

- [ ] **`i18n-test` 页面是开发测试页，不符合 Miro visual surface。**
  - 证据：`src/app/(normal)/i18n-test/page.tsx:36`、`src/app/(normal)/i18n-test/page.tsx:63`、`src/app/(normal)/i18n-test/page.tsx:93`
  - 当前问题：大量 test card、uppercase labels、generic shadcn interactions。
  - 修复方向：若仍对用户可见，应隐藏或重做；若只是开发页，应从正常导航和生产构建中隔离。

## 7. 需要优先处理的修复顺序

1. [ ] 先修全局 design tokens：字体、typography utilities、surface mapping、radius、shadow、dark mode 策略。
2. [ ] 再修基础 UI primitives：Button、Input、SearchInput、Select/filter、Badge、Dropdown/Dialog elevation。
3. [ ] 然后修营销首页关键 surface：nav、hero centered layout、whiteboard mockup、pricing toggle/table、footer。
4. [ ] 最后清理业务页局部 class：移除普通卡片 shadow、统一 radius、减少 generic yellow、把 headings 接到 Miro typography token。
5. [ ] 为 Console / Project workbench 单独补充 app UI 规范，避免用 marketing-only `DESIGN.md` 硬套复杂 API 编辑器。

## 8. 快速复核命令

```bash
rg -n "figma-|Inter|JetBrains|font-semibold|shadow-sm|shadow-soft|bg-highlight|rounded-2xl|rounded-\\[1\\.75rem\\]|font-mono" web/src --glob '*.{tsx,ts,css}'
rg -n "defaultTheme=\"system\"|enableSystem|\\.dark|--bg-surface|--radius-xs|--shadow-md" web/src --glob '*.{tsx,ts,css}'
```
