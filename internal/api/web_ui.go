package api

const webUIHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Agent Runtime</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #05070a;
      --panel: rgba(10, 15, 23, 0.94);
      --panel-2: #0d141e;
      --card: #111923;
      --line: #1e2b39;
      --line-2: #2c4260;
      --text: #e8f3ff;
      --muted: #91a4ba;
      --faint: #5f7288;
      --cyan: #25d7cf;
      --blue: #4c86ff;
      --green: #39d97a;
      --amber: #f4bd57;
      --red: #ff5c74;
      --mono: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      color: var(--text);
      background:
        linear-gradient(rgba(37, 215, 207, 0.035) 1px, transparent 1px),
        linear-gradient(90deg, rgba(76, 134, 255, 0.025) 1px, transparent 1px),
        radial-gradient(circle at 18% -8%, rgba(37, 215, 207, 0.16), transparent 30%),
        radial-gradient(circle at 88% 0%, rgba(76, 134, 255, 0.14), transparent 28%),
        #05070a;
      background-size: 36px 36px, 36px 36px, auto, auto, auto;
      line-height: 1.4;
    }
    button, input, select { font: inherit; }
    button {
      min-height: 36px;
      border: 1px solid var(--line-2);
      border-radius: 8px;
      background: rgba(16, 25, 37, 0.92);
      color: var(--text);
      cursor: pointer;
      font-weight: 760;
      letter-spacing: 0;
    }
    button:hover { border-color: rgba(37, 215, 207, 0.72); }
    button:disabled { cursor: not-allowed; opacity: 0.45; }
    .primary {
      border-color: rgba(76, 134, 255, 0.86);
      background: linear-gradient(135deg, #3178ff, #4c86ff);
      color: #eef5ff;
      box-shadow: 0 0 22px rgba(76, 134, 255, 0.2);
    }
    .ghost { background: rgba(6, 10, 15, 0.7); color: var(--muted); }
    .danger { border-color: rgba(255, 92, 116, 0.48); background: rgba(255, 92, 116, 0.08); color: #ff9cad; }
    input, select {
      width: 100%;
      min-height: 38px;
      padding: 8px 10px;
      border: 1px solid var(--line);
      border-radius: 8px;
      outline: none;
      background: rgba(4, 8, 13, 0.82);
      color: var(--text);
    }
    input:focus, select:focus { border-color: rgba(37, 215, 207, 0.72); box-shadow: 0 0 0 3px rgba(37, 215, 207, 0.1); }
    label {
      display: block;
      margin-bottom: 6px;
      color: var(--muted);
      font-size: 11px;
      font-weight: 850;
      text-transform: uppercase;
    }
    code, pre, .mono { font-family: var(--mono); }
    a { color: var(--cyan); text-decoration: none; }
    a:hover { text-decoration: underline; }
    .shell {
      width: min(1480px, calc(100vw - 32px));
      margin: 0 auto;
      padding: 16px 0 28px;
    }
    .topbar {
      position: sticky;
      top: 12px;
      z-index: 5;
      display: grid;
      grid-template-columns: minmax(260px, 1fr) auto auto;
      gap: 14px;
      align-items: center;
      padding: 12px 14px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(8, 13, 20, 0.9);
      backdrop-filter: blur(18px);
      box-shadow: 0 20px 70px rgba(0, 0, 0, 0.3);
    }
    .brand { display: grid; grid-template-columns: 38px minmax(0, 1fr); gap: 10px; align-items: center; }
    .brand-mark {
      width: 38px;
      height: 38px;
      display: grid;
      place-items: center;
      border: 1px solid rgba(37, 215, 207, 0.54);
      border-radius: 8px;
      background: linear-gradient(135deg, rgba(37, 215, 207, 0.17), rgba(76, 134, 255, 0.14)), #071018;
      color: var(--cyan);
      font-family: var(--mono);
      font-weight: 900;
    }
    .brand-title { font-size: 17px; font-weight: 850; }
    .brand-subtitle { color: var(--faint); font-size: 12px; }
    .view-switch, .lang-toggle {
      display: grid;
      grid-auto-flow: column;
      gap: 4px;
      padding: 4px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(3, 7, 12, 0.72);
    }
    .view-switch button, .lang-toggle button {
      min-height: 30px;
      border: 0;
      padding: 0 12px;
      background: transparent;
      color: var(--muted);
      font-size: 12px;
    }
    .view-switch button.active, .lang-toggle button.active { background: rgba(37, 215, 207, 0.15); color: var(--cyan); }
    .top-actions { display: flex; gap: 8px; align-items: center; justify-content: flex-end; }
    .hero {
      margin-top: 16px;
      display: grid;
      grid-template-columns: minmax(280px, 1fr) minmax(520px, 720px);
      gap: 16px;
      align-items: stretch;
    }
    .headline {
      min-height: 150px;
      padding: 18px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: linear-gradient(135deg, rgba(13, 20, 30, 0.96), rgba(8, 13, 20, 0.9));
      display: flex;
      flex-direction: column;
      justify-content: space-between;
    }
    .headline h1 { margin: 0; font-size: 26px; line-height: 1.12; }
    .headline p { margin: 8px 0 0; color: var(--muted); max-width: 720px; }
    .status-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 10px; }
    .metric { min-height: 72px; padding: 11px 12px; border: 1px solid var(--line); border-radius: 8px; background: rgba(10, 17, 26, 0.74); }
    .metric-label { color: var(--muted); font-size: 11px; font-weight: 850; text-transform: uppercase; }
    .metric-value { margin-top: 8px; display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 850; overflow-wrap: anywhere; }
    .led { width: 8px; height: 8px; border-radius: 999px; background: var(--amber); box-shadow: 0 0 14px rgba(244, 189, 87, 0.65); flex: 0 0 auto; }
    .led.ok { background: var(--green); box-shadow: 0 0 14px rgba(57, 217, 122, 0.72); }
    .led.bad { background: var(--red); box-shadow: 0 0 14px rgba(255, 92, 116, 0.72); }
    .badge {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      min-height: 26px;
      padding: 4px 8px;
      border: 1px solid var(--line);
      border-radius: 999px;
      background: rgba(255, 255, 255, 0.035);
      color: var(--muted);
      font-size: 12px;
      font-weight: 850;
      white-space: nowrap;
    }
    .badge.ok { color: #b7ffd4; border-color: rgba(57, 217, 122, 0.38); }
    .badge.warn { color: #ffe0a3; border-color: rgba(244, 189, 87, 0.42); }
    .badge.bad { color: #ffb4c0; border-color: rgba(255, 92, 116, 0.42); }
    .content { margin-top: 16px; }
    .view { display: none; }
    .view.active { display: block; }
    .home-grid { display: grid; grid-template-columns: minmax(0, 1fr) 390px; gap: 16px; align-items: start; }
    .panel { border: 1px solid var(--line); border-radius: 8px; background: var(--panel); box-shadow: 0 24px 90px rgba(0, 0, 0, 0.28); overflow: hidden; }
    .panel-header { min-height: 58px; display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 13px 14px; border-bottom: 1px solid var(--line); }
    .panel-title { display: flex; align-items: center; gap: 10px; min-width: 0; }
    .panel-title h2 { margin: 0; font-size: 15px; }
    .panel-title p { margin: 2px 0 0; color: var(--muted); font-size: 12px; overflow-wrap: anywhere; }
    .nav-icon {
      width: 24px;
      height: 24px;
      display: grid;
      place-items: center;
      border-radius: 7px;
      background: rgba(255, 255, 255, 0.04);
      color: var(--cyan);
      font-family: var(--mono);
      font-size: 12px;
      font-weight: 900;
    }
    .manager-body { padding: 14px; }
    .cli-grid { display: grid; grid-template-columns: repeat(2, minmax(280px, 1fr)); gap: 10px; }
    .cli-card {
      min-height: 128px;
      display: grid;
      grid-template-columns: 58px minmax(0, 1fr);
      gap: 12px;
      align-items: start;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: linear-gradient(180deg, rgba(21, 31, 43, 0.98), rgba(13, 20, 30, 0.98));
    }
    .cli-logo { width: 46px; height: 46px; display: grid; place-items: center; overflow: hidden; border: 1px solid rgba(44, 66, 96, 0.78); border-radius: 8px; background: rgba(3, 7, 12, 0.62); color: var(--cyan); font-size: 13px; font-weight: 900; }
    .cli-logo img { max-width: 36px; max-height: 36px; object-fit: contain; }
    .cli-logo span { width: 100%; height: 100%; display: none; place-items: center; }
    .cli-main { min-width: 0; }
    .cli-head { display: flex; gap: 8px; justify-content: space-between; align-items: start; }
    .cli-name { font-size: 16px; font-weight: 850; }
    .cli-provider { color: var(--faint); font-size: 12px; }
    .cli-version { margin-top: 4px; color: var(--muted); font-family: var(--mono); font-size: 13px; overflow-wrap: anywhere; }
    .cli-actions { margin-top: 12px; display: flex; flex-wrap: wrap; gap: 8px; }
    .cli-actions button, .cli-actions a { min-height: 34px; padding: 7px 10px; font-size: 12px; }
    .cli-actions a { display: inline-grid; place-items: center; border: 1px solid var(--line-2); border-radius: 8px; background: rgba(6, 10, 15, 0.7); color: var(--muted); font-weight: 760; }
    .activity-body { padding: 14px; display: grid; gap: 10px; }
    .activity-log {
      min-height: 360px;
      max-height: calc(100vh - 430px);
      margin: 0;
      padding: 12px;
      overflow: auto;
      white-space: pre-wrap;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(3, 7, 12, 0.84);
      color: #d9e8f8;
      font-size: 12px;
      line-height: 1.45;
    }
    .activity-links { display: grid; gap: 8px; }
    .activity-links a { display: block; padding: 8px 10px; border: 1px solid rgba(37, 215, 207, 0.3); border-radius: 8px; background: rgba(37, 215, 207, 0.08); overflow-wrap: anywhere; }
    .action-input { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 8px; }
    .action-input[hidden] { display: none; }
    .action-input button { padding: 0 12px; }
    .swagger-panel { min-height: calc(100vh - 262px); }
    .swagger-frame {
      width: 100%;
      height: calc(100vh - 274px);
      min-height: 720px;
      display: block;
      border: 0;
      background: #0b1118;
    }
    .swagger-actions { display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
    .swagger-actions a { display: inline-grid; place-items: center; min-height: 34px; padding: 7px 10px; border: 1px solid var(--line-2); border-radius: 8px; background: rgba(6, 10, 15, 0.7); color: var(--muted); font-size: 12px; font-weight: 760; }
    .empty { padding: 18px; color: var(--muted); border: 1px dashed var(--line-2); border-radius: 8px; background: rgba(255, 255, 255, 0.025); }
    .toast { position: fixed; right: 18px; bottom: 18px; z-index: 10; display: none; max-width: 420px; padding: 12px 14px; border: 1px solid rgba(37, 215, 207, 0.42); border-radius: 8px; background: rgba(5, 10, 16, 0.96); color: var(--text); box-shadow: 0 20px 70px rgba(0, 0, 0, 0.45); }
    .toast.show { display: block; }
    @media (max-width: 1180px) {
      .topbar, .hero, .home-grid { grid-template-columns: 1fr; }
      .cli-grid { grid-template-columns: 1fr; }
    }
    @media (max-width: 760px) {
      .shell { width: min(100vw - 20px, 1480px); padding-top: 10px; }
      .topbar { position: static; }
      .top-actions, .status-grid { grid-template-columns: 1fr; display: grid; }
      .view-switch, .lang-toggle { width: 100%; grid-auto-flow: column; }
      .cli-card { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <div class="shell">
    <header class="topbar">
      <div class="brand">
        <div class="brand-mark">AR</div>
        <div>
          <div class="brand-title">Agent Runtime</div>
          <div class="brand-subtitle" data-i18n="brandSubtitle">CLI 控制平面</div>
        </div>
      </div>
      <div class="view-switch" aria-label="View switch">
        <button class="active" type="button" data-view="home-view" data-i18n="home">首页</button>
        <button type="button" data-view="api-view" data-i18n="api">API</button>
      </div>
      <div class="top-actions">
        <button class="ghost" id="refresh" type="button" data-i18n="refresh">刷新</button>
        <div class="lang-toggle" aria-label="Language">
          <button id="lang-zh" type="button" data-lang="zh">中文</button>
          <button id="lang-en" type="button" data-lang="en">EN</button>
        </div>
      </div>
    </header>

    <section class="hero">
      <div class="headline">
        <div>
          <h1 id="page-heading">CLI Agent Runtime</h1>
          <p id="page-subtitle">统一安装、授权和探测 Claude Code、Codex、Gemini、OpenCode 等 CLI。</p>
        </div>
      </div>
      <div class="status-grid">
        <div class="metric">
          <div class="metric-label" data-i18n="health">健康状态</div>
          <div class="metric-value"><span class="led" id="health-led"></span><span id="health">loading</span></div>
        </div>
        <div class="metric">
          <div class="metric-label" data-i18n="ready">就绪状态</div>
          <div class="metric-value"><span class="led" id="ready-led"></span><span id="ready">loading</span></div>
        </div>
        <div class="metric">
          <div class="metric-label" data-i18n="availableCli">可用 CLI</div>
          <div class="metric-value"><span id="available-cli">0 / 0</span></div>
        </div>
        <div class="metric">
          <div class="metric-label" data-i18n="usersOnline">使用者</div>
          <div class="metric-value"><span id="user-count">0</span></div>
        </div>
      </div>
    </section>

    <main class="content">
      <section class="view active" id="home-view">
        <div class="home-grid">
          <section class="panel">
            <div class="panel-header">
              <div class="panel-title">
                <span class="nav-icon">CL</span>
                <div>
                  <h2>CLI Manager</h2>
                  <p data-i18n="cliManagerDesc">安装、授权和验证都通过按钮完成，不暴露 shell。</p>
                </div>
              </div>
              <button class="ghost" id="refresh-tools" type="button" data-i18n="refresh">刷新</button>
            </div>
            <div hidden>
              <span id="session-badge"><span id="session-led"></span><span id="session-label">initializing</span></span>
              <span id="context-label">- / - / -</span>
              <select id="tenant"></select>
              <input id="workspace" value="repo-main">
              <input id="profile" value="team-default" list="profile-options">
              <datalist id="profile-options"></datalist>
            </div>
            <div class="manager-body">
              <div class="cli-grid" id="installed-panel"></div>
            </div>
          </section>

          <aside class="panel">
            <div class="panel-header">
              <div class="panel-title">
                <span class="nav-icon">AC</span>
                <div>
                  <h2 data-i18n="activity">操作状态</h2>
                  <p id="action-title" data-i18n="activityDesc">安装和授权输出会整理在这里。</p>
                </div>
              </div>
              <button class="ghost" id="stop-action" type="button" data-i18n="stop">停止</button>
            </div>
            <div class="activity-body">
              <div class="badge warn"><span class="led" id="action-led"></span><span id="action-state" data-i18n="idle">空闲</span></div>
              <div class="activity-links" id="activity-links"></div>
              <div class="action-input" id="action-input-box" hidden>
                <input id="action-input" autocomplete="one-time-code" data-i18n-placeholder="actionInputPlaceholder" placeholder="输入授权码或令牌">
                <button class="primary" id="send-action-input" type="button" data-i18n="send">发送</button>
              </div>
              <pre class="activity-log" id="activity-log"></pre>
            </div>
          </aside>
        </div>
      </section>

      <section class="view" id="api-view">
        <section class="panel swagger-panel">
          <div class="panel-header">
            <div class="panel-title">
              <span class="nav-icon">API</span>
              <div>
                <h2>Swagger UI</h2>
                <p data-i18n="swaggerDesc">基于 /openapi.json 的交互式 API 文档。</p>
              </div>
            </div>
            <div class="swagger-actions">
              <a href="/docs" target="_blank" rel="noopener noreferrer">/docs</a>
              <a href="/openapi.json" target="_blank" rel="noopener noreferrer">/openapi.json</a>
            </div>
          </div>
          <iframe class="swagger-frame" id="swagger-frame" src="/docs" title="Swagger UI"></iframe>
        </section>
      </section>
    </main>
  </div>

  <div class="toast" id="toast"></div>

  <script>
    const messages = {
      zh: {
        brandSubtitle: 'CLI 控制平面',
        home: '首页',
        api: 'API',
        refresh: '刷新',
        health: '健康状态',
        ready: '就绪状态',
        availableCli: '可用 CLI',
        usersOnline: '使用者',
        cliManagerDesc: '安装、授权和验证都通过按钮完成，不暴露 shell。',
        installCli: '安装',
        authorize: '授权',
        verify: '验证',
        officialSource: '官方来源',
        remove: '移除',
        healthOK: 'Health OK',
        notInstalled: '未安装',
        initializing: '初始化中',
        checking: '检查中',
        unavailable: '不可用',
        activity: '操作状态',
        activityDesc: '安装和授权输出会整理在这里。',
        stop: '停止',
        send: '发送',
        actionInputPlaceholder: '输入授权码或令牌',
        idle: '空闲',
        running: '运行中',
        loginOK: '会话已就绪',
        loginFailed: '默认会话不可用',
        refreshed: '状态已刷新',
        commandStarted: '已启动',
        commandStopped: '已停止',
        copied: '已复制',
        swaggerDesc: '基于 /openapi.json 的交互式 API 文档。',
        apiOverview: 'Agent Runtime API',
        apiOverviewDesc: 'Swagger UI 会直接读取 /openapi.json。',
        statusDesc: '查询运行时状态、CLI 数量和使用者数量。',
        registerDesc: '开放注册租户用户，返回调用自己租户资源的 Bearer Token。',
        toolsDesc: '列出已注册 CLI，并可按租户和凭据配置探测真实 PATH 状态。',
        toolUpdateDesc: '管理员注册或更新 CLI wrapper。普通安装建议使用首页 CLI Manager。',
        cliActionsDesc: '启动后台安装、授权或验证动作；授权输出会提取 URL 和一次性 code。',
        jobsDesc: '服务间调用入口。调用方只能使用 token 策略允许的 tool、workspace 和 credential profile。',
        eventsDesc: '读取 job 事件流。返回 Server-Sent Events，适合 curl 和服务端消费。',
        eventsWSDesc: '实时 WebSocket job 事件流，stdout、stderr、status、exit 会逐条推送。',
        appServerDesc: '为外部客户端启动 tenant 隔离的 codex app-server，并透明转发 WebSocket 消息。'
      },
      en: {
        brandSubtitle: 'CLI control plane',
        home: 'Home',
        api: 'API',
        refresh: 'Refresh',
        health: 'Health',
        ready: 'Ready',
        availableCli: 'Available CLI',
        usersOnline: 'Users',
        cliManagerDesc: 'Install, authorize, and verify CLIs through UI actions without exposing a shell.',
        installCli: 'Install',
        authorize: 'Authorize',
        verify: 'Verify',
        officialSource: 'Official Source',
        remove: 'Remove',
        healthOK: 'Health OK',
        notInstalled: 'Not installed',
        initializing: 'Initializing',
        checking: 'Checking',
        unavailable: 'Unavailable',
        activity: 'Activity',
        activityDesc: 'Install and authorization output is summarized here.',
        stop: 'Stop',
        send: 'Send',
        actionInputPlaceholder: 'Enter auth code or token',
        idle: 'Idle',
        running: 'Running',
        loginOK: 'Session ready',
        loginFailed: 'Default session unavailable',
        refreshed: 'Runtime status refreshed',
        commandStarted: 'Started',
        commandStopped: 'Stopped',
        copied: 'Copied',
        swaggerDesc: 'Interactive API documentation powered by /openapi.json.',
        apiOverview: 'Agent Runtime API',
        apiOverviewDesc: 'Swagger UI reads /openapi.json directly.',
        statusDesc: 'Inspect runtime health, CLI counts, and user counts.',
        registerDesc: 'Register a tenant user and receive a Bearer token for that user’s own tenant resources.',
        toolsDesc: 'List registered CLIs and probe tenant/profile-specific PATH health.',
        toolUpdateDesc: 'Admins can register or update CLI wrappers. Normal installs should use the home CLI Manager.',
        cliActionsDesc: 'Start background install, authorization, or verification actions; auth output extracts URLs and one-time codes.',
        jobsDesc: 'Service-to-service execution entrypoint constrained by token policy.',
        eventsDesc: 'Read job event output as Server-Sent Events for curl and backend clients.',
        eventsWSDesc: 'Real-time WebSocket job event stream for stdout, stderr, status, and exit events.',
        appServerDesc: 'Start a tenant-isolated codex app-server for external clients and transparently proxy WebSocket messages.'
      }
    };

    const savedLanguage = localStorage.getItem('agent-runtime-lang');
    const state = {
      lang: savedLanguage || 'zh',
      view: 'home-view',
      sessionToken: localStorage.getItem('agent-runtime-session-token') || localStorage.getItem('agent-runtime-token') || '',
      session: null,
      tools: [],
      toolsLoaded: false,
      toolsContextReady: false,
      tenants: [],
      users: [],
      actionBusy: false,
      activeActionID: '',
      activeActionKind: '',
      actionEventCount: 0,
      activityText: '',
      actionPollTimer: null,
      installPollTimer: null,
      cliRenderSignature: ''
    };

    const installSources = [
      { label: 'Claude Code', fallback: 'CC', logo: '/assets/logos/claude.svg', tool: 'claude', docs: 'https://docs.anthropic.com/en/docs/claude-code/quickstart', provider: 'Anthropic' },
      { label: 'Codex', fallback: 'CX', logo: 'https://avatars.githubusercontent.com/u/14957082?s=96&v=4', tool: 'codex', docs: 'https://github.com/openai/codex', provider: 'OpenAI' },
      { label: 'Gemini', fallback: 'GM', logo: 'https://avatars.githubusercontent.com/u/161781182?s=96&v=4', tool: 'gemini', docs: 'https://github.com/google-gemini/gemini-cli', provider: 'Google' },
      { label: 'OpenCode', fallback: 'OC', logo: 'https://opencode.ai/favicon-96x96-v3.png', tool: 'opencode', docs: 'https://opencode.ai/download', provider: 'SST' },
      { label: 'iFlow', fallback: 'IF', logo: 'https://img.alicdn.com/imgextra/i1/O1CN01jgdyc81WIsdSepA4X_!!6000000002766-55-tps-162-162.svg', tool: 'iflow', docs: 'https://platform.iflow.cn/cli/quickstart', provider: 'iFlow' },
      { label: 'Kimi', fallback: 'KM', logo: 'https://www.kimi.com/favicon.ico', tool: 'kimi', docs: 'https://www.kimi.com/code/docs/en/kimi-code-cli/getting-started.html', provider: 'Moonshot AI' },
      { label: 'Qoder', fallback: 'QD', logo: '/assets/logos/qoder.svg', tool: 'qoder', docs: 'https://docs.qoder.com/cli/quick-start', provider: 'Qoder' }
    ];

    function $(id) { return document.getElementById(id); }
    function t(key) { return (messages[state.lang] && messages[state.lang][key]) || messages.en[key] || key; }
    function escapeHTML(value) {
      return String(value == null ? '' : value).replace(/[&<>"']/g, function(char) {
        return ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char];
      });
    }
    function authHeaders() {
      return state.sessionToken ? { Authorization: 'Bearer ' + state.sessionToken } : {};
    }
    function showToast(message) {
      const toast = $('toast');
      toast.textContent = message;
      toast.classList.add('show');
      window.clearTimeout(showToast.timer);
      showToast.timer = window.setTimeout(function() { toast.classList.remove('show'); }, 3200);
    }
    async function api(path, options) {
      options = options || {};
      options.headers = Object.assign({}, authHeaders(), options.headers || {});
      const response = await fetch(path, options);
      if (!response.ok) {
        let message = response.statusText;
        try {
          const body = await response.json();
          message = body.error || message;
        } catch (err) {}
        throw new Error(message);
      }
      if (response.status === 204) return null;
      return response.json();
    }

    function applyLanguage() {
      document.documentElement.lang = state.lang === 'zh' ? 'zh-CN' : 'en';
      document.querySelectorAll('[data-i18n]').forEach(function(el) { el.textContent = t(el.dataset.i18n); });
      document.querySelectorAll('[data-i18n-placeholder]').forEach(function(el) { el.placeholder = t(el.dataset.i18nPlaceholder); });
      $('lang-zh').classList.toggle('active', state.lang === 'zh');
      $('lang-en').classList.toggle('active', state.lang === 'en');
      updatePageTitle();
      renderSession();
      renderCliManager();
      renderActionState();
    }

    function updatePageTitle() {
      if (state.view === 'api-view') {
        $('page-heading').textContent = 'Agent Runtime API';
        $('page-subtitle').textContent = t('apiOverviewDesc');
      } else {
        $('page-heading').textContent = 'CLI Agent Runtime';
        $('page-subtitle').textContent = t('cliManagerDesc');
      }
    }

    async function loginDefault() {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: 'admin', password: 'admin' })
      });
      if (!response.ok) throw new Error(t('loginFailed'));
      const body = await response.json();
      state.sessionToken = body.token || '';
      state.session = body.session || null;
      localStorage.setItem('agent-runtime-session-token', state.sessionToken);
      localStorage.removeItem('agent-runtime-token');
    }

    async function restoreOrLogin() {
      if (state.session) return;
      if (state.sessionToken) {
        try {
          state.session = await api('/api/session');
          localStorage.setItem('agent-runtime-session-token', state.sessionToken);
          localStorage.removeItem('agent-runtime-token');
          return;
        } catch (err) {
          state.sessionToken = '';
          localStorage.removeItem('agent-runtime-session-token');
          localStorage.removeItem('agent-runtime-token');
        }
      }
      await loginDefault();
    }

    async function refresh() {
      const results = await Promise.allSettled([api('/api/health'), api('/api/ready'), api('/api/status')]);
      setMetric('health', 'health-led', results[0].status === 'fulfilled' ? results[0].value.status : 'error');
      setMetric('ready', 'ready-led', results[1].status === 'fulfilled' ? results[1].value.status : 'error');
      if (results[2].status === 'fulfilled') {
        $('user-count').textContent = String(results[2].value.users || 0);
      }
      try {
        await restoreOrLogin();
        await refreshSecure();
      } catch (err) {
        state.session = null;
        state.tenants = [];
        state.users = [];
        state.toolsContextReady = false;
        await refreshTools().catch(function() {});
        renderSession();
        renderContextOptions();
        renderCliManager();
        showToast(err.message);
      }
    }

    async function refreshSecure() {
      await refreshTenants();
      renderContextOptions();
      await refreshTools();
      if (state.session && state.session.admin) {
        await refreshUsers().catch(function() {});
      }
      renderSession();
      renderCliManager();
    }

    async function refreshTools() {
      let path = '/api/tools';
      state.toolsContextReady = !!(state.session && $('tenant').value && $('profile').value.trim());
      if (state.toolsContextReady) {
        const params = new URLSearchParams({ tenant: $('tenant').value, credential_profile: $('profile').value.trim() });
        path += '?' + params.toString();
      }
      const body = await api(path);
      state.tools = body.tools || [];
      state.toolsLoaded = true;
      const available = state.tools.filter(function(tool) { return tool.available; }).length;
      $('available-cli').textContent = available + ' / ' + state.tools.length;
    }

    async function refreshTenants() {
      const body = await api('/api/tenants');
      state.tenants = body.tenants || [];
    }

    async function refreshUsers() {
      const body = await api('/api/users');
      state.users = body.users || [];
      if (state.users.length) $('user-count').textContent = String(state.users.length);
    }

    function setMetric(textID, ledID, value) {
      $(textID).textContent = value;
      const led = $(ledID);
      led.classList.remove('ok', 'bad');
      if (value === 'ok' || value === 'ready') led.classList.add('ok');
      if (value === 'error') led.classList.add('bad');
    }

    function renderSession() {
      const loggedIn = !!state.session;
      $('session-led').classList.toggle('ok', loggedIn);
      $('session-led').classList.toggle('bad', !loggedIn);
      $('session-badge').classList.toggle('ok', loggedIn);
      $('session-badge').classList.toggle('bad', !loggedIn);
      $('session-label').textContent = loggedIn ? state.session.subject : t('idle');
    }

    function renderContextOptions() {
      const tenantSelect = $('tenant');
      const currentTenant = tenantSelect.value;
      const sourceTenants = state.tenants.length ? state.tenants : (state.session ? [{ id: state.session.tenant, credential_profiles: state.session.allowed_credential_profiles, workspace_patterns: state.session.allowed_workspaces }] : []);
      tenantSelect.innerHTML = sourceTenants.map(function(tenant) {
        return '<option value="' + escapeHTML(tenant.id) + '">' + escapeHTML(tenant.id) + '</option>';
      }).join('');
      if (currentTenant) tenantSelect.value = currentTenant;
      if (!tenantSelect.value && sourceTenants[0]) tenantSelect.value = sourceTenants[0].id;
      updateProfileOptions();
    }

    function updateProfileOptions() {
      const tenant = state.tenants.find(function(item) { return item.id === $('tenant').value; });
      const profiles = tenant && tenant.credential_profiles ? tenant.credential_profiles : (state.session ? state.session.allowed_credential_profiles || [] : []);
      $('profile-options').innerHTML = profiles.map(function(profile) {
        return '<option value="' + escapeHTML(profile) + '"></option>';
      }).join('');
      if (profiles.length && (!$('profile').value || $('profile').value === 'team-default')) $('profile').value = profiles[0];
      const workspaces = tenant && tenant.workspace_patterns ? tenant.workspace_patterns : (state.session ? state.session.allowed_workspaces || [] : []);
      if (workspaces[0] && (!$('workspace').value || $('workspace').value === 'repo-main')) $('workspace').value = workspaces[0].replace('*', 'main');
      updateContextLabels();
    }

    function updateContextLabels() {
      $('context-label').textContent = ($('tenant').value || '-') + ' / ' + ($('workspace').value || '-') + ' / ' + ($('profile').value || '-');
    }

    async function startCliAction(tool, action, title) {
      if (!state.session) {
        showToast(t('loginFailed'));
        return;
      }
      if (!state.toolsContextReady) {
        showToast(t('checking'));
        return;
      }
      stopActionPolling();
      state.activityText = '';
      state.actionBusy = true;
      state.activeActionID = '';
      state.activeActionKind = action;
      state.actionEventCount = 0;
      $('activity-log').textContent = '';
      $('activity-links').innerHTML = '';
      $('action-title').textContent = title;
      renderActionState();
      const body = await api('/api/cli-actions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tool: tool,
          action: action,
          tenant: $('tenant').value,
          workspace: $('workspace').value.trim(),
          credential_profile: $('profile').value.trim()
        })
      });
      state.activeActionID = body.id;
      appendActivity('$ ' + (body.command || (tool + ' ' + action)) + '\n');
      showToast(t('commandStarted') + ': ' + title);
      await pollCliAction(true);
      state.actionPollTimer = window.setInterval(function() { pollCliAction(false).catch(function(err) { showToast(err.message); }); }, 1400);
    }

    function stopActionPolling() {
      if (state.actionPollTimer) window.clearInterval(state.actionPollTimer);
      state.actionPollTimer = null;
    }

    async function pollCliAction(runOnce) {
      if (!state.activeActionID) return;
      const action = await api('/api/cli-actions/' + encodeURIComponent(state.activeActionID));
      const events = action.events || [];
      events.slice(state.actionEventCount).forEach(function(event) {
        if (event.type === 'status') appendActivity('[status] ' + event.message + '\n');
        if (event.type === 'stdout' || event.type === 'stderr') appendActivity(event.message || '');
        if (event.type === 'input') appendActivity('[input] ' + event.message);
        if (event.type === 'exit') appendActivity('[exit] ' + event.message + '\n');
      });
      state.actionEventCount = events.length;
      renderAuthHints(action);
      const done = ['succeeded', 'failed', 'timed_out', 'canceled'].indexOf(action.status) >= 0;
      state.actionBusy = !done;
      renderActionState();
      if (done && !runOnce) stopActionPolling();
      if (done && (action.action === 'install' || action.action === 'verify')) {
        await refreshTools().then(renderCliManager).catch(function() {});
      }
    }

    function renderAuthHints(action) {
      const urls = action.auth_urls || [];
      const codes = action.auth_codes || [];
      const html = []
        .concat(urls.slice(-5).map(function(url) {
          return '<a href="' + escapeHTML(url) + '" target="_blank" rel="noopener noreferrer">' + escapeHTML(url) + '</a>';
        }))
        .concat(codes.slice(-5).map(function(code) {
          return '<button class="ghost" type="button" data-copy-code="' + escapeHTML(code) + '">' + escapeHTML(code) + '</button>';
        }))
        .join('');
      if (html) $('activity-links').innerHTML = html;
      $('activity-links').querySelectorAll('[data-copy-code]').forEach(function(button) {
        button.addEventListener('click', function() {
          navigator.clipboard.writeText(button.dataset.copyCode || '').then(function() { showToast(t('copied')); }).catch(function() {});
        });
      });
    }

    async function sendActionInput() {
      const input = $('action-input');
      const data = input.value.trim();
      if (!data || !state.activeActionID) return;
      await api('/api/cli-actions/' + encodeURIComponent(state.activeActionID) + '/input', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ data: data })
      });
      input.value = '';
      showToast(t('commandStarted'));
      await pollCliAction(true);
    }

    async function cancelCliAction() {
      if (state.activeActionID && state.actionBusy) {
        await api('/api/cli-actions/' + encodeURIComponent(state.activeActionID), { method: 'DELETE' });
      }
      stopActionPolling();
      state.actionBusy = false;
      renderActionState();
    }

    function cleanOutput(value) {
      return String(value || '').replace(/\x1b\[[0-9;?]*[ -/]*[@-~]/g, '');
    }

    function appendActivity(value) {
      state.activityText += cleanOutput(value);
      if (state.activityText.length > 24000) state.activityText = state.activityText.slice(-24000);
      const log = $('activity-log');
      log.textContent = state.activityText.trimStart();
      log.scrollTop = log.scrollHeight;
      renderActivityLinks();
    }

    function renderActivityLinks() {
      const urls = Array.from(new Set((state.activityText.match(/https?:\/\/[^\s"'<>]+/g) || []).map(function(url) {
        return url.replace(/[),.]+$/, '');
      }))).slice(-5);
      if (!$('activity-links').innerHTML) $('activity-links').innerHTML = urls.map(function(url) {
        return '<a href="' + escapeHTML(url) + '" target="_blank" rel="noopener noreferrer">' + escapeHTML(url) + '</a>';
      }).join('');
    }

    function renderActionState() {
      const key = state.actionBusy ? 'running' : 'idle';
      $('action-state').textContent = t(key);
      $('action-led').classList.toggle('ok', state.actionBusy);
      $('action-led').classList.toggle('bad', false);
      $('action-input-box').hidden = !(state.actionBusy && state.activeActionKind === 'auth');
    }

    function startInstallPolling() {
      if (state.installPollTimer) window.clearInterval(state.installPollTimer);
      let remaining = 40;
      state.installPollTimer = window.setInterval(function() {
        if (remaining <= 0) {
          window.clearInterval(state.installPollTimer);
          state.installPollTimer = null;
          return;
        }
        remaining--;
        refreshTools().then(renderCliManager).catch(function() {});
      }, 3000);
    }

    function renderCliLogo(source) {
      const fallback = source.fallback || source.label.slice(0, 2).toUpperCase();
      return '<div class="cli-logo"><img src="' + escapeHTML(source.logo || '') + '" alt="' + escapeHTML(source.label) + '" loading="lazy" referrerpolicy="no-referrer" onerror="this.style.display=&quot;none&quot;;this.nextElementSibling.style.display=&quot;grid&quot;"><span>' + escapeHTML(fallback) + '</span></div>';
    }

    function renderCliManager() {
      const known = new Map(state.tools.map(function(tool) { return [tool.name, tool]; }));
      const cards = installSources.map(function(source) {
        const tool = known.get(source.tool) || { name: source.tool, available: false, health: 'missing' };
        const canCheck = state.toolsLoaded && state.toolsContextReady;
        const available = canCheck && !!tool.available;
        const version = available ? (tool.detected_version || tool.version || '-') : (canCheck ? t('notInstalled') : t('initializing'));
        const healthClass = available ? 'ok' : (canCheck ? 'bad' : 'warn');
        const healthText = available ? t('healthOK') : (canCheck ? t('notInstalled') : t('checking'));
        const disabled = canCheck ? '' : ' disabled';
        const actions = available
          ? '<button class="ghost" type="button" data-cli-action="auth" data-cli-tool="' + escapeHTML(source.tool) + '"' + disabled + '>' + escapeHTML(t('authorize')) + '</button>' +
            '<button class="ghost" type="button" data-cli-action="verify" data-cli-tool="' + escapeHTML(source.tool) + '"' + disabled + '>' + escapeHTML(t('verify')) + '</button>' +
            (state.session && state.session.admin && known.has(source.tool) ? '<button class="danger" type="button" data-delete-tool="' + escapeHTML(source.tool) + '">' + escapeHTML(t('remove')) + '</button>' : '')
          : '<button class="primary" type="button" data-cli-action="install" data-cli-tool="' + escapeHTML(source.tool) + '"' + disabled + '>' + escapeHTML(t('installCli')) + '</button>' +
            '<button class="ghost" type="button" data-cli-action="auth" data-cli-tool="' + escapeHTML(source.tool) + '"' + disabled + '>' + escapeHTML(t('authorize')) + '</button>' +
            '<a href="' + escapeHTML(source.docs) + '" target="_blank" rel="noopener noreferrer">' + escapeHTML(t('officialSource')) + '</a>';
        return '<article class="cli-card">' +
          renderCliLogo(source) +
          '<div class="cli-main">' +
            '<div class="cli-head"><div><div class="cli-name">' + escapeHTML(source.label) + '</div><div class="cli-provider">' + escapeHTML(source.provider) + '</div></div><span class="badge ' + healthClass + '"><span class="led ' + (available ? 'ok' : 'bad') + '"></span>' + escapeHTML(healthText) + '</span></div>' +
            '<div class="cli-version">' + escapeHTML(version) + '</div>' +
            '<div class="cli-actions">' + actions + '</div>' +
          '</div>' +
        '</article>';
      }).join('');
      const signature = state.lang + '|' + state.toolsContextReady + '|' + state.tools.map(function(tool) {
        return [tool.name, tool.available, tool.detected_version || tool.version || '', tool.health || '', tool.error || ''].join(':');
      }).join('|');
      if (signature === state.cliRenderSignature && $('installed-panel').innerHTML) return;
      state.cliRenderSignature = signature;
      $('installed-panel').innerHTML = cards;
      bindManagerButtons($('installed-panel'));
    }

    function bindManagerButtons(root) {
      root.querySelectorAll('[data-cli-action]').forEach(function(button) {
        button.addEventListener('click', function() {
          const action = button.dataset.cliAction;
          const tool = button.dataset.cliTool;
          startCliAction(tool, action, button.textContent.trim()).then(function() {
            if (action === 'install') startInstallPolling();
          }).catch(function(err) { showToast(err.message); });
        });
      });
      root.querySelectorAll('[data-delete-tool]').forEach(function(button) {
        button.addEventListener('click', async function() {
          await api('/api/tools/' + encodeURIComponent(button.dataset.deleteTool), { method: 'DELETE' });
          await refreshTools();
          renderCliManager();
        });
      });
    }

    function switchView(viewID) {
      state.view = viewID;
      document.querySelectorAll('[data-view]').forEach(function(item) { item.classList.toggle('active', item.dataset.view === viewID); });
      document.querySelectorAll('.view').forEach(function(item) { item.classList.toggle('active', item.id === viewID); });
      updatePageTitle();
    }

    document.querySelectorAll('[data-view]').forEach(function(button) { button.addEventListener('click', function() { switchView(button.dataset.view); }); });
    document.querySelectorAll('[data-lang]').forEach(function(button) {
      button.addEventListener('click', function() {
        state.lang = button.dataset.lang;
        localStorage.setItem('agent-runtime-lang', state.lang);
        applyLanguage();
      });
    });
    $('refresh').addEventListener('click', function() { refresh().then(function() { showToast(t('refreshed')); }).catch(function(err) { showToast(err.message); }); });
    $('refresh-tools').addEventListener('click', function() { refreshTools().then(renderCliManager).catch(function(err) { showToast(err.message); }); });
    $('tenant').addEventListener('change', function() { updateProfileOptions(); refreshTools().then(renderCliManager).catch(function() {}); });
    $('workspace').addEventListener('input', updateContextLabels);
    $('profile').addEventListener('input', function() { updateContextLabels(); refreshTools().then(renderCliManager).catch(function() {}); });
    $('stop-action').addEventListener('click', function() { cancelCliAction().then(function() { showToast(t('commandStopped')); }).catch(function(err) { showToast(err.message); }); });
    $('send-action-input').addEventListener('click', function() { sendActionInput().catch(function(err) { showToast(err.message); }); });
    $('action-input').addEventListener('keydown', function(event) {
      if (event.key === 'Enter') {
        event.preventDefault();
        sendActionInput().catch(function(err) { showToast(err.message); });
      }
    });

    applyLanguage();
    renderActionState();
    refresh().catch(function(err) {
      setMetric('health', 'health-led', 'error');
      setMetric('ready', 'ready-led', 'error');
      showToast(err.message);
    });
  </script>
</body>
</html>`
