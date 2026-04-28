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
    .api-shell { display: grid; grid-template-columns: 280px minmax(0, 1fr); gap: 16px; align-items: start; }
    .api-index { position: sticky; top: 104px; padding: 10px; display: grid; gap: 8px; }
    .api-index button { text-align: left; padding: 9px 10px; background: transparent; color: var(--muted); }
    .api-index button.active { border-color: rgba(37, 215, 207, 0.6); background: rgba(37, 215, 207, 0.1); color: var(--text); }
    .api-docs { display: grid; gap: 12px; }
    .api-card { border: 1px solid var(--line); border-radius: 8px; background: rgba(7, 12, 18, 0.72); overflow: hidden; }
    .api-card-head { display: flex; gap: 10px; justify-content: space-between; align-items: center; padding: 13px 14px; border-bottom: 1px solid var(--line); }
    .api-card-head h3 { margin: 0; font-size: 15px; }
    .api-card-body { padding: 14px; display: grid; gap: 10px; color: var(--muted); }
    .method { min-width: 58px; justify-content: center; font-family: var(--mono); }
    .method.get { color: #b7ffd4; border-color: rgba(57, 217, 122, 0.38); }
    .method.post { color: #c8dcff; border-color: rgba(76, 134, 255, 0.42); }
    .method.delete { color: #ffb4c0; border-color: rgba(255, 92, 116, 0.42); }
    .method.ws { color: #c5f9ff; border-color: rgba(37, 215, 207, 0.42); }
    .code-block { position: relative; }
    .code-block pre { margin: 0; padding: 12px; overflow: auto; border: 1px solid var(--line); border-radius: 8px; background: rgba(3, 7, 12, 0.84); color: #d9e8f8; font-size: 12px; line-height: 1.45; }
    .copy { position: absolute; right: 8px; top: 8px; min-height: 28px; font-size: 11px; }
    .empty { padding: 18px; color: var(--muted); border: 1px dashed var(--line-2); border-radius: 8px; background: rgba(255, 255, 255, 0.025); }
    .toast { position: fixed; right: 18px; bottom: 18px; z-index: 10; display: none; max-width: 420px; padding: 12px 14px; border: 1px solid rgba(37, 215, 207, 0.42); border-radius: 8px; background: rgba(5, 10, 16, 0.96); color: var(--text); box-shadow: 0 20px 70px rgba(0, 0, 0, 0.45); }
    .toast.show { display: block; }
    @media (max-width: 1180px) {
      .topbar, .hero, .home-grid, .api-shell { grid-template-columns: 1fr; }
      .api-index { position: static; }
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
              <pre class="activity-log" id="activity-log"></pre>
            </div>
          </aside>
        </div>
      </section>

      <section class="view" id="api-view">
        <div class="api-shell">
          <aside class="panel api-index">
            <button class="active" type="button" data-api-target="overview">Overview</button>
            <button type="button" data-api-target="status">GET /api/status</button>
            <button type="button" data-api-target="tools">GET /api/tools</button>
            <button type="button" data-api-target="tool-update">POST /api/tools</button>
            <button type="button" data-api-target="jobs">POST /api/jobs</button>
            <button type="button" data-api-target="job-events">GET /api/jobs/{id}/events</button>
            <button type="button" data-api-target="terminal-api">WS /api/terminal</button>
          </aside>
          <section class="api-docs">
            <article class="api-card" id="overview">
              <div class="api-card-head">
                <h3 data-i18n="apiOverview">Agent Runtime API</h3>
                <span class="badge">OpenAPI</span>
              </div>
              <div class="api-card-body">
                <p data-i18n="apiOverviewDesc">服务调用默认使用 Bearer Token。人用入口走首页按钮，服务间调用走 Job API。</p>
                <div class="code-block"><button class="copy" data-copy="#base-url" type="button">Copy</button><pre id="base-url">Base URL: BASE_URL
Authorization: Bearer &lt;token&gt;</pre></div>
                <p><a href="/openapi.json" target="_blank" rel="noopener noreferrer">/openapi.json</a></p>
              </div>
            </article>

            <article class="api-card" id="status">
              <div class="api-card-head"><h3><span class="badge method get">GET</span> /api/status</h3></div>
              <div class="api-card-body">
                <p data-i18n="statusDesc">查询运行时状态、CLI 数量和使用者数量。</p>
                <div class="code-block"><button class="copy" data-copy="#status-code" type="button">Copy</button><pre id="status-code">curl -s BASE_URL/api/status</pre></div>
              </div>
            </article>

            <article class="api-card" id="tools">
              <div class="api-card-head"><h3><span class="badge method get">GET</span> /api/tools</h3></div>
              <div class="api-card-body">
                <p data-i18n="toolsDesc">列出已注册 CLI，并可按租户和凭据配置探测真实 PATH 状态。</p>
                <div class="code-block"><button class="copy" data-copy="#tools-code" type="button">Copy</button><pre id="tools-code">curl -s "BASE_URL/api/tools?tenant=team-a&credential_profile=team-default" \
  -H "Authorization: Bearer &lt;token&gt;"</pre></div>
              </div>
            </article>

            <article class="api-card" id="tool-update">
              <div class="api-card-head"><h3><span class="badge method post">POST</span> /api/tools</h3></div>
              <div class="api-card-body">
                <p data-i18n="toolUpdateDesc">管理员注册或更新 CLI wrapper。普通安装建议使用首页 CLI Manager。</p>
                <div class="code-block"><button class="copy" data-copy="#tool-code" type="button">Copy</button><pre id="tool-code">curl -s -X POST BASE_URL/api/tools \
  -H "Authorization: Bearer &lt;admin-token&gt;" \
  -H "Content-Type: application/json" \
  -d '{"name":"codex","path":"codex","version":"official","credential_env":"CODEX_HOME","credential_subdir":".codex"}'</pre></div>
              </div>
            </article>

            <article class="api-card" id="jobs">
              <div class="api-card-head"><h3><span class="badge method post">POST</span> /api/jobs</h3></div>
              <div class="api-card-body">
                <p data-i18n="jobsDesc">服务间调用入口。调用方只能使用 token 策略允许的 tool、workspace 和 credential profile。</p>
                <div class="code-block"><button class="copy" data-copy="#job-code" type="button">Copy</button><pre id="job-code">curl -s -X POST BASE_URL/api/jobs \
  -H "Authorization: Bearer &lt;token&gt;" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant": "team-a",
    "tool": "codex",
    "args": ["exec", "fix tests"],
    "workspace": "repo-main",
    "credential_profile": "team-default",
    "timeout_seconds": 900
  }'</pre></div>
              </div>
            </article>

            <article class="api-card" id="job-events">
              <div class="api-card-head"><h3><span class="badge method get">GET</span> /api/jobs/{id}/events</h3></div>
              <div class="api-card-body">
                <p data-i18n="eventsDesc">读取 job 事件流。当前实现返回 Server-Sent Events。</p>
                <div class="code-block"><button class="copy" data-copy="#events-code" type="button">Copy</button><pre id="events-code">curl -N BASE_URL/api/jobs/&lt;job-id&gt;/events</pre></div>
              </div>
            </article>

            <article class="api-card" id="terminal-api">
              <div class="api-card-head"><h3><span class="badge method ws">WS</span> /api/terminal</h3></div>
              <div class="api-card-body">
                <p data-i18n="terminalApiDesc">交互式 PTY API 仍保留给集成方；首页不会显示终端。</p>
                <div class="code-block"><button class="copy" data-copy="#terminal-code" type="button">Copy</button><pre id="terminal-code">WSS_BASE/api/terminal?token=&lt;token&gt;&tenant=team-a&workspace=repo-main&credential_profile=team-default</pre></div>
              </div>
            </article>
          </section>
        </div>
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
        idle: '空闲',
        running: '运行中',
        disconnected: '未连接',
        connected: '已连接',
        loginOK: '会话已就绪',
        loginFailed: '默认会话不可用',
        refreshed: '状态已刷新',
        commandStarted: '已启动',
        commandStopped: '已停止',
        copied: '已复制',
        apiOverview: 'Agent Runtime API',
        apiOverviewDesc: '服务调用默认使用 Bearer Token。人用入口走首页按钮，服务间调用走 Job API。',
        statusDesc: '查询运行时状态、CLI 数量和使用者数量。',
        toolsDesc: '列出已注册 CLI，并可按租户和凭据配置探测真实 PATH 状态。',
        toolUpdateDesc: '管理员注册或更新 CLI wrapper。普通安装建议使用首页 CLI Manager。',
        jobsDesc: '服务间调用入口。调用方只能使用 token 策略允许的 tool、workspace 和 credential profile。',
        eventsDesc: '读取 job 事件流。当前实现返回 Server-Sent Events。',
        terminalApiDesc: '交互式 PTY API 仍保留给集成方；首页不会显示终端。'
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
        idle: 'Idle',
        running: 'Running',
        disconnected: 'Disconnected',
        connected: 'Connected',
        loginOK: 'Session ready',
        loginFailed: 'Default session unavailable',
        refreshed: 'Runtime status refreshed',
        commandStarted: 'Started',
        commandStopped: 'Stopped',
        copied: 'Copied',
        apiOverview: 'Agent Runtime API',
        apiOverviewDesc: 'Service calls use Bearer tokens by default. Human workflows use the home screen; service workflows use the Job API.',
        statusDesc: 'Inspect runtime health, CLI counts, and user counts.',
        toolsDesc: 'List registered CLIs and probe tenant/profile-specific PATH health.',
        toolUpdateDesc: 'Admins can register or update CLI wrappers. Normal installs should use the home CLI Manager.',
        jobsDesc: 'Service-to-service execution entrypoint constrained by token policy.',
        eventsDesc: 'Read job event output. Current implementation returns Server-Sent Events.',
        terminalApiDesc: 'Interactive PTY API remains available for integrations; the home screen does not show a terminal.'
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
      ws: null,
      connected: false,
      actionBusy: false,
      activityText: '',
      installPollTimer: null
    };

    const installSources = [
      { label: 'Claude Code', fallback: 'CC', logo: '/assets/logos/claude.svg', tool: 'claude', command: 'curl -fsSL https://claude.ai/install.sh | bash', verify: 'claude --version', login: 'claude', docs: 'https://docs.anthropic.com/en/docs/claude-code/quickstart', provider: 'Anthropic' },
      { label: 'Codex', fallback: 'CX', logo: 'https://avatars.githubusercontent.com/u/14957082?s=96&v=4', tool: 'codex', command: 'npm install -g @openai/codex', verify: 'codex --version', login: 'codex login', docs: 'https://github.com/openai/codex', provider: 'OpenAI' },
      { label: 'Gemini', fallback: 'GM', logo: 'https://avatars.githubusercontent.com/u/161781182?s=96&v=4', tool: 'gemini', command: 'npm install -g @google/gemini-cli', verify: 'gemini --version', login: 'gemini', docs: 'https://github.com/google-gemini/gemini-cli', provider: 'Google' },
      { label: 'OpenCode', fallback: 'OC', logo: 'https://opencode.ai/favicon-96x96-v3.png', tool: 'opencode', command: 'curl -fsSL https://opencode.ai/install | bash', verify: 'opencode --version', login: 'opencode auth login', docs: 'https://opencode.ai/download', provider: 'SST' },
      { label: 'iFlow', fallback: 'IF', logo: 'https://img.alicdn.com/imgextra/i1/O1CN01jgdyc81WIsdSepA4X_!!6000000002766-55-tps-162-162.svg', tool: 'iflow', command: 'bash -c "$(curl -fsSL https://gitee.com/iflow-ai/iflow-cli/raw/main/install.sh)"', verify: 'iflow --version', login: 'iflow', docs: 'https://platform.iflow.cn/cli/quickstart', provider: 'iFlow' },
      { label: 'Kimi', fallback: 'KM', logo: 'https://www.kimi.com/favicon.ico', tool: 'kimi', command: 'curl -LsSf https://code.kimi.com/install.sh | bash', verify: 'kimi --version', login: 'kimi', docs: 'https://www.kimi.com/code/docs/en/kimi-code-cli/getting-started.html', provider: 'Moonshot AI' },
      { label: 'Qoder', fallback: 'QD', logo: 'https://docs.qoder.com/mintlify-assets/_mintlify/favicons/qoder/-6DIoH8zsEnnm9G9/_generated/favicon-dark/favicon-32x32.png', tool: 'qoder', command: 'curl -fsSL https://qoder.com/install | bash', verify: 'qodercli --version', login: 'qodercli', docs: 'https://docs.qoder.com/cli/quick-start', provider: 'Qoder' }
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
      $('lang-zh').classList.toggle('active', state.lang === 'zh');
      $('lang-en').classList.toggle('active', state.lang === 'en');
      applyApiExamples();
      updatePageTitle();
      renderSession();
      renderCliManager();
      renderActionState();
    }

    function applyApiExamples() {
      const base = window.location.origin;
      const wsBase = (window.location.protocol === 'https:' ? 'wss://' : 'ws://') + window.location.host;
      $('base-url').textContent = 'Base URL: ' + base + '\nAuthorization: Bearer <token>';
      $('status-code').textContent = 'curl -s ' + base + '/api/status';
      $('tools-code').textContent = 'curl -s "' + base + '/api/tools?tenant=team-a&credential_profile=team-default" \\\n  -H "Authorization: Bearer <token>"';
      $('tool-code').textContent = 'curl -s -X POST ' + base + '/api/tools \\\n  -H "Authorization: Bearer <admin-token>" \\\n  -H "Content-Type: application/json" \\\n  -d ' + "'{\"name\":\"codex\",\"path\":\"codex\",\"version\":\"official\",\"credential_env\":\"CODEX_HOME\",\"credential_subdir\":\".codex\"}'";
      $('job-code').textContent = 'curl -s -X POST ' + base + '/api/jobs \\\n  -H "Authorization: Bearer <token>" \\\n  -H "Content-Type: application/json" \\\n  -d ' + "'{\\n    \"tenant\": \"team-a\",\\n    \"tool\": \"codex\",\\n    \"args\": [\"exec\", \"fix tests\"],\\n    \"workspace\": \"repo-main\",\\n    \"credential_profile\": \"team-default\",\\n    \"timeout_seconds\": 900\\n  }'";
      $('events-code').textContent = 'curl -N ' + base + '/api/jobs/<job-id>/events';
      $('terminal-code').textContent = wsBase + '/api/terminal?token=<token>&tenant=team-a&workspace=repo-main&credential_profile=team-default';
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
      $('session-label').textContent = loggedIn ? state.session.subject : t('disconnected');
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

    function terminalURL() {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const params = new URLSearchParams({
        token: state.sessionToken,
        tenant: $('tenant').value,
        workspace: $('workspace').value.trim(),
        credential_profile: $('profile').value.trim(),
        cols: '120',
        rows: '32'
      });
      return protocol + '//' + window.location.host + '/api/v1/terminal/ws?' + params.toString();
    }

    function connectActionSocket() {
      if (state.ws && state.ws.readyState === WebSocket.OPEN) return Promise.resolve(state.ws);
      return new Promise(function(resolve, reject) {
        if (!state.session) {
          reject(new Error(t('loginFailed')));
          return;
        }
        const ws = new WebSocket(terminalURL());
        state.ws = ws;
        const timer = window.setTimeout(function() { reject(new Error('connection timeout')); }, 8000);
        ws.onopen = function() {
          window.clearTimeout(timer);
          state.connected = true;
          renderActionState();
          resolve(ws);
        };
        ws.onmessage = function(event) {
          try {
            const payload = JSON.parse(event.data);
            if (payload.type === 'output') appendActivity(payload.data || '');
            if (payload.type === 'error') appendActivity('\n[error] ' + (payload.data || 'unknown error') + '\n');
            if (payload.type === 'exit') {
              state.connected = false;
              state.actionBusy = false;
              renderActionState();
            }
          } catch (err) {
            appendActivity(String(event.data));
          }
        };
        ws.onclose = function() {
          state.connected = false;
          state.actionBusy = false;
          renderActionState();
        };
        ws.onerror = function() {
          state.connected = false;
          state.actionBusy = false;
          renderActionState();
        };
      });
    }

    function disconnectActionSocket() {
      if (state.ws) state.ws.close();
      state.ws = null;
      state.connected = false;
      state.actionBusy = false;
      renderActionState();
    }

    async function runCommand(command, title) {
      if (!state.session) {
        showToast(t('loginFailed'));
        return;
      }
      $('action-title').textContent = title;
      appendActivity('\n$ ' + command + '\n');
      state.actionBusy = true;
      renderActionState();
      const ws = await connectActionSocket();
      window.setTimeout(function() {
        ws.send(JSON.stringify({ type: 'input', data: command + '\r' }));
      }, 180);
      showToast(t('commandStarted') + ': ' + title);
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
      $('activity-links').innerHTML = urls.map(function(url) {
        return '<a href="' + escapeHTML(url) + '" target="_blank" rel="noopener noreferrer">' + escapeHTML(url) + '</a>';
      }).join('');
    }

    function renderActionState() {
      const key = state.actionBusy ? 'running' : (state.connected ? 'connected' : 'idle');
      $('action-state').textContent = t(key);
      $('action-led').classList.toggle('ok', state.connected || state.actionBusy);
      $('action-led').classList.toggle('bad', false);
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
      $('installed-panel').innerHTML = installSources.map(function(source) {
        const tool = known.get(source.tool) || { name: source.tool, available: false, health: 'missing' };
        const canCheck = state.toolsLoaded && state.toolsContextReady;
        const available = canCheck && !!tool.available;
        const version = available ? (tool.detected_version || tool.version || '-') : (canCheck ? t('notInstalled') : t('initializing'));
        const healthClass = available ? 'ok' : (canCheck ? 'bad' : 'warn');
        const healthText = available ? t('healthOK') : (canCheck ? t('notInstalled') : t('checking'));
        const disabled = canCheck ? '' : ' disabled';
        const actions = available
          ? '<button class="ghost" type="button" data-login-command="' + escapeHTML(source.login) + '"' + disabled + '>' + escapeHTML(t('authorize')) + '</button>' +
            '<button class="ghost" type="button" data-install-command="' + escapeHTML(source.verify) + '"' + disabled + '>' + escapeHTML(t('verify')) + '</button>' +
            (state.session && state.session.admin && known.has(source.tool) ? '<button class="danger" type="button" data-delete-tool="' + escapeHTML(source.tool) + '">' + escapeHTML(t('remove')) + '</button>' : '')
          : '<button class="primary" type="button" data-install-command="' + escapeHTML(source.command) + '"' + disabled + '>' + escapeHTML(t('installCli')) + '</button>' +
            '<button class="ghost" type="button" data-login-command="' + escapeHTML(source.login) + '"' + disabled + '>' + escapeHTML(t('authorize')) + '</button>' +
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
      bindManagerButtons($('installed-panel'));
    }

    function bindManagerButtons(root) {
      root.querySelectorAll('[data-install-command]').forEach(function(button) {
        button.addEventListener('click', function() {
          runCommand(button.dataset.installCommand, button.textContent.trim()).then(startInstallPolling).catch(function(err) { showToast(err.message); });
        });
      });
      root.querySelectorAll('[data-login-command]').forEach(function(button) {
        button.addEventListener('click', function() {
          runCommand(button.dataset.loginCommand, button.textContent.trim()).catch(function(err) { showToast(err.message); });
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

    function bindApiIndex() {
      document.querySelectorAll('[data-api-target]').forEach(function(button) {
        button.addEventListener('click', function() {
          document.querySelectorAll('[data-api-target]').forEach(function(item) { item.classList.toggle('active', item === button); });
          const target = document.getElementById(button.dataset.apiTarget);
          if (target) target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        });
      });
      document.querySelectorAll('[data-copy]').forEach(function(button) {
        button.addEventListener('click', async function() {
          const target = document.querySelector(button.dataset.copy);
          if (!target) return;
          await navigator.clipboard.writeText(target.textContent);
          showToast(t('copied'));
        });
      });
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
    $('stop-action').addEventListener('click', function() { disconnectActionSocket(); showToast(t('commandStopped')); });

    bindApiIndex();
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
