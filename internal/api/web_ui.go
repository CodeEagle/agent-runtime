package api

const webUIHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Agent Runtime</title>
  <link rel="stylesheet" href="/assets/xterm/xterm.css">
  <style>
    :root {
      color-scheme: dark;
      --bg: #05070a;
      --bg-panel: rgba(9, 14, 22, 0.88);
      --bg-card: #0b111a;
      --bg-card-2: #0f1722;
      --line: #1d2a38;
      --line-strong: #2a3d51;
      --text: #e8f2ff;
      --muted: #94a7bd;
      --faint: #607286;
      --cyan: #28e0d4;
      --blue: #4c8dff;
      --green: #37e681;
      --amber: #f5bc4f;
      --red: #ff5c73;
      --terminal-bg: #05070a;
      --mono: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    html, body { min-height: 100%; }
    body {
      margin: 0;
      background:
        linear-gradient(rgba(40, 224, 212, 0.035) 1px, transparent 1px),
        linear-gradient(90deg, rgba(76, 141, 255, 0.025) 1px, transparent 1px),
        radial-gradient(circle at 25% -10%, rgba(40, 224, 212, 0.16), transparent 32%),
        linear-gradient(135deg, #040609, #09121c 58%, #030507);
      background-size: 36px 36px, 36px 36px, auto, auto;
      color: var(--text);
      line-height: 1.45;
    }
    button, input, select { font: inherit; }
    button {
      min-height: 36px;
      border: 1px solid var(--line-strong);
      border-radius: 8px;
      background: rgba(17, 27, 40, 0.9);
      color: var(--text);
      cursor: pointer;
      font-weight: 760;
      letter-spacing: 0;
    }
    button:hover { border-color: rgba(40, 224, 212, 0.66); }
    button:disabled { cursor: not-allowed; opacity: 0.45; }
    button.primary {
      border-color: rgba(40, 224, 212, 0.84);
      background: linear-gradient(135deg, rgba(40, 224, 212, 0.96), rgba(76, 141, 255, 0.9));
      color: #021014;
      box-shadow: 0 0 22px rgba(40, 224, 212, 0.18);
    }
    button.ghost {
      background: rgba(7, 12, 18, 0.68);
      color: var(--muted);
    }
    button.danger {
      border-color: rgba(255, 92, 115, 0.55);
      background: rgba(255, 92, 115, 0.08);
      color: #ff9aaa;
    }
    input, select {
      width: 100%;
      min-height: 38px;
      padding: 8px 10px;
      border: 1px solid var(--line);
      border-radius: 8px;
      outline: none;
      background: rgba(3, 7, 12, 0.76);
      color: var(--text);
    }
    input:focus, select:focus {
      border-color: rgba(40, 224, 212, 0.78);
      box-shadow: 0 0 0 3px rgba(40, 224, 212, 0.1);
    }
    label {
      display: block;
      margin-bottom: 6px;
      color: var(--muted);
      font-size: 11px;
      font-weight: 850;
      text-transform: uppercase;
    }
    code, .mono { font-family: var(--mono); }
    .app-shell {
      min-height: 100vh;
      display: grid;
      grid-template-columns: 248px minmax(0, 1fr);
    }
    .rail {
      position: sticky;
      top: 0;
      height: 100vh;
      padding: 20px 16px;
      border-right: 1px solid var(--line);
      background: rgba(4, 8, 13, 0.84);
      backdrop-filter: blur(18px);
      display: flex;
      flex-direction: column;
      gap: 18px;
    }
    .brand {
      display: grid;
      grid-template-columns: 38px minmax(0, 1fr);
      gap: 10px;
      align-items: center;
      min-height: 48px;
    }
    .brand-mark {
      width: 38px;
      height: 38px;
      display: grid;
      place-items: center;
      border: 1px solid rgba(40, 224, 212, 0.56);
      border-radius: 8px;
      background: linear-gradient(135deg, rgba(40, 224, 212, 0.18), rgba(76, 141, 255, 0.16)), #071018;
      color: var(--cyan);
      font-family: var(--mono);
      font-weight: 900;
      box-shadow: inset 0 0 24px rgba(40, 224, 212, 0.1);
    }
    .brand-title { font-size: 17px; font-weight: 850; }
    .brand-subtitle { color: var(--faint); font-size: 12px; }
    .nav {
      display: grid;
      gap: 8px;
    }
    .nav-button {
      width: 100%;
      display: grid;
      grid-template-columns: 24px minmax(0, 1fr) auto;
      gap: 10px;
      align-items: center;
      padding: 9px 10px;
      text-align: left;
      background: transparent;
      color: var(--muted);
    }
    .nav-button.active {
      border-color: rgba(40, 224, 212, 0.6);
      background: rgba(40, 224, 212, 0.1);
      color: var(--text);
    }
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
    .nav-kbd {
      color: var(--faint);
      font-family: var(--mono);
      font-size: 11px;
    }
    .rail-card {
      margin-top: auto;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(10, 16, 24, 0.74);
    }
    .rail-card-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 8px;
      margin-bottom: 9px;
      color: var(--muted);
      font-size: 12px;
      text-transform: uppercase;
      font-weight: 850;
    }
    .workspace {
      min-width: 0;
      padding: 16px;
    }
    .topbar {
      min-height: 64px;
      display: grid;
      grid-template-columns: minmax(280px, 1fr) minmax(480px, 700px);
      gap: 16px;
      align-items: center;
      padding: 12px 14px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(8, 13, 20, 0.86);
      backdrop-filter: blur(18px);
      box-shadow: 0 20px 70px rgba(0, 0, 0, 0.28);
    }
    .page-title h1 {
      margin: 0;
      font-size: 20px;
      line-height: 1.1;
      letter-spacing: 0;
    }
    .page-title p {
      margin: 5px 0 0;
      color: var(--muted);
      font-size: 13px;
    }
    .top-actions {
      display: grid;
      grid-template-columns: 1fr 94px 126px;
      gap: 8px;
      align-items: center;
    }
    .token-field {
      display: grid;
      grid-template-columns: 80px minmax(0, 1fr);
      gap: 8px;
      align-items: center;
    }
    .token-field span {
      color: var(--muted);
      font-size: 12px;
      font-weight: 850;
      text-transform: uppercase;
    }
    .lang-toggle {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 4px;
      padding: 4px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(3, 7, 12, 0.72);
    }
    .lang-toggle button {
      min-height: 28px;
      border: 0;
      background: transparent;
      color: var(--muted);
      font-size: 12px;
    }
    .lang-toggle button.active {
      background: rgba(40, 224, 212, 0.15);
      color: var(--cyan);
    }
    .content { margin-top: 16px; }
    .view { display: none; }
    .view.active { display: block; }
    .status-grid {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 10px;
      margin-bottom: 16px;
    }
    .metric {
      min-height: 66px;
      padding: 11px 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(10, 17, 26, 0.74);
    }
    .metric-label {
      color: var(--muted);
      font-size: 11px;
      font-weight: 850;
      text-transform: uppercase;
    }
    .metric-value {
      margin-top: 8px;
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      font-weight: 850;
      overflow-wrap: anywhere;
    }
    .led {
      width: 8px;
      height: 8px;
      border-radius: 999px;
      background: var(--amber);
      box-shadow: 0 0 14px rgba(245, 188, 79, 0.65);
      flex: 0 0 auto;
    }
    .led.ok { background: var(--green); box-shadow: 0 0 14px rgba(55, 230, 129, 0.72); }
    .led.bad { background: var(--red); box-shadow: 0 0 14px rgba(255, 92, 115, 0.72); }
    .terminal-grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 340px;
      gap: 16px;
      align-items: start;
    }
    .panel {
      border: 1px solid var(--line);
      border-radius: 8px;
      background: var(--bg-panel);
      box-shadow: 0 24px 90px rgba(0, 0, 0, 0.28);
      overflow: hidden;
    }
    .panel-header {
      min-height: 58px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 13px 14px;
      border-bottom: 1px solid var(--line);
    }
    .panel-title {
      display: flex;
      align-items: center;
      gap: 10px;
      min-width: 0;
    }
    .panel-title h2 {
      margin: 0;
      font-size: 15px;
      letter-spacing: 0;
    }
    .panel-title p {
      margin: 2px 0 0;
      color: var(--muted);
      font-size: 12px;
      overflow-wrap: anywhere;
    }
    .panel-actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      justify-content: flex-end;
    }
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
    .badge.ok { color: #b7ffd4; border-color: rgba(55, 230, 129, 0.38); }
    .badge.warn { color: #ffe0a3; border-color: rgba(245, 188, 79, 0.42); }
    .context-bar {
      display: grid;
      grid-template-columns: 150px minmax(160px, 1fr) minmax(160px, 1fr);
      gap: 10px;
      padding: 14px;
      border-bottom: 1px solid var(--line);
      background: rgba(5, 9, 14, 0.35);
      align-items: end;
    }
    .terminal-body {
      height: clamp(500px, calc(100vh - 336px), 760px);
      min-height: 500px;
      padding: 12px;
    }
    .xterm-frame {
      height: 100%;
      overflow: hidden;
      border: 1px solid #172332;
      border-radius: 8px;
      background: var(--terminal-bg);
      box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.02), inset 0 0 42px rgba(40, 224, 212, 0.04);
    }
    #terminal-container {
      height: 100%;
      width: 100%;
      padding: 8px;
    }
    .xterm {
      height: 100%;
      padding: 2px;
    }
    .xterm .xterm-viewport {
      background-color: transparent;
      scrollbar-color: #233447 transparent;
    }
    .quick-grid {
      display: grid;
      gap: 8px;
      padding: 12px;
    }
    .quick-button {
      width: 100%;
      min-height: 46px;
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 8px;
      align-items: center;
      padding: 10px 11px;
      text-align: left;
      background: rgba(9, 16, 24, 0.72);
    }
    .quick-button strong {
      display: block;
      font-size: 13px;
    }
    .quick-button span span {
      display: block;
      margin-top: 2px;
      color: var(--faint);
      font-family: var(--mono);
      font-size: 11px;
      font-weight: 500;
    }
    .side-stack {
      display: grid;
      gap: 16px;
    }
    .tool-cards {
      display: grid;
      gap: 10px;
      padding: 12px;
    }
    .tool-card {
      padding: 11px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(5, 10, 16, 0.62);
    }
    .tool-card-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
    }
    .tool-name {
      color: var(--text);
      font-family: var(--mono);
      font-size: 13px;
      font-weight: 900;
    }
    .tool-path {
      margin-top: 8px;
      color: var(--muted);
      font-family: var(--mono);
      font-size: 11px;
      overflow-wrap: anywhere;
    }
    .tool-meta {
      display: grid;
      gap: 7px;
      margin-top: 10px;
      color: var(--muted);
      font-size: 12px;
    }
    .tool-meta div {
      display: grid;
      grid-template-columns: 110px minmax(0, 1fr);
      gap: 8px;
      align-items: start;
    }
    .tool-meta strong {
      color: var(--faint);
      font-size: 11px;
      font-weight: 850;
      text-transform: uppercase;
    }
    .tool-actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      margin-top: 12px;
    }
    .manager-grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 360px;
      gap: 16px;
      align-items: start;
    }
    .manager-tools {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
      gap: 12px;
      padding: 14px;
    }
    .manager-tools .tool-card {
      min-height: 190px;
      display: flex;
      flex-direction: column;
    }
    .manager-tools .tool-actions {
      margin-top: auto;
      padding-top: 12px;
    }
    .install-grid {
      display: grid;
      gap: 10px;
      padding: 12px;
    }
    .install-card {
      display: grid;
      gap: 10px;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(5, 10, 16, 0.66);
    }
    .install-card-head {
      display: flex;
      justify-content: space-between;
      gap: 10px;
      align-items: flex-start;
    }
    .install-title {
      display: grid;
      gap: 2px;
      min-width: 0;
    }
    .install-title strong {
      color: var(--text);
      font-size: 13px;
    }
    .install-title span {
      color: var(--faint);
      font-size: 11px;
    }
    .install-command {
      padding: 8px;
      border: 1px solid rgba(42, 61, 81, 0.82);
      border-radius: 8px;
      background: rgba(3, 7, 12, 0.8);
      color: #d5e8ff;
      font-family: var(--mono);
      font-size: 11px;
      overflow-wrap: anywhere;
    }
    .install-actions {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 8px;
    }
    .install-actions a {
      min-height: 36px;
      display: grid;
      place-items: center;
      padding: 8px 10px;
      border: 1px solid var(--line-strong);
      border-radius: 8px;
      background: rgba(7, 12, 18, 0.68);
      color: var(--muted);
      font-size: 13px;
      font-weight: 760;
      text-decoration: none;
    }
    .install-actions button {
      width: 100%;
    }
    .form-stack {
      display: grid;
      gap: 12px;
      padding: 14px;
    }
    .field-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 10px;
    }
    .tenant-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
      gap: 14px;
      padding: 14px;
    }
    .tenant-card {
      padding: 14px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(6, 12, 18, 0.7);
    }
    .tenant-name {
      display: flex;
      justify-content: space-between;
      gap: 10px;
      align-items: center;
      margin-bottom: 12px;
      font-weight: 900;
    }
    .kv {
      display: grid;
      gap: 9px;
      color: var(--muted);
      font-size: 12px;
    }
    .kv div {
      display: grid;
      grid-template-columns: 116px minmax(0, 1fr);
      gap: 8px;
    }
    .kv strong {
      color: var(--faint);
      font-size: 11px;
      font-weight: 850;
      text-transform: uppercase;
    }
    .empty {
      padding: 18px;
      color: var(--muted);
      border: 1px dashed var(--line-strong);
      border-radius: 8px;
      background: rgba(255, 255, 255, 0.025);
    }
    .toast {
      position: fixed;
      right: 18px;
      bottom: 18px;
      z-index: 10;
      display: none;
      max-width: 420px;
      padding: 12px 14px;
      border: 1px solid rgba(40, 224, 212, 0.42);
      border-radius: 8px;
      background: rgba(5, 10, 16, 0.95);
      color: var(--text);
      box-shadow: 0 20px 70px rgba(0, 0, 0, 0.45);
    }
    .toast.show { display: block; }
    @media (max-width: 1120px) {
      .app-shell { grid-template-columns: 1fr; }
      .rail {
        position: relative;
        height: auto;
        flex-direction: row;
        align-items: center;
        overflow-x: auto;
      }
      .nav { grid-auto-flow: column; grid-auto-columns: max-content; }
      .rail-card { display: none; }
      .topbar { grid-template-columns: 1fr; }
      .terminal-grid, .manager-grid { grid-template-columns: 1fr; }
    }
    @media (max-width: 760px) {
      .workspace { padding: 10px; }
      .top-actions, .status-grid, .context-bar, .field-row { grid-template-columns: 1fr; }
      .token-field { grid-template-columns: 1fr; }
      .terminal-body { height: 460px; min-height: 460px; }
    }
  </style>
</head>
<body>
  <div class="app-shell">
    <aside class="rail">
      <div class="brand">
        <div class="brand-mark">AR</div>
        <div>
          <div class="brand-title">Agent Runtime</div>
          <div class="brand-subtitle" data-i18n="brandSubtitle">CLI 控制平面</div>
        </div>
      </div>
      <nav class="nav" aria-label="Main navigation">
        <button class="nav-button active" type="button" data-view="terminal-view">
          <span class="nav-icon">&gt;_</span><span data-i18n="navTerminal">终端</span><span class="nav-kbd">01</span>
        </button>
        <button class="nav-button" type="button" data-view="tools-view">
          <span class="nav-icon">CL</span><span data-i18n="navTools">CLI 管理</span><span class="nav-kbd">02</span>
        </button>
        <button class="nav-button" type="button" data-view="tenants-view">
          <span class="nav-icon">TN</span><span data-i18n="navTenants">租户</span><span class="nav-kbd">03</span>
        </button>
      </nav>
      <div class="rail-card">
        <div class="rail-card-title"><span data-i18n="runtime">运行时</span><span id="rail-health">Health OK</span></div>
        <div class="badge ok"><span class="led ok"></span><span id="rail-context">team-a / repo-main</span></div>
      </div>
    </aside>

    <div class="workspace">
      <header class="topbar">
        <div class="page-title">
          <h1 id="page-heading">终端</h1>
          <p id="page-subtitle">用于登录认证、CLI 调试和手动维护的完整 shell。</p>
        </div>
        <div class="top-actions">
          <div class="token-field">
            <span data-i18n="token">Token</span>
            <input id="token" type="password" autocomplete="off" placeholder="dev-token">
          </div>
          <button class="ghost" id="refresh" type="button" data-i18n="refresh">刷新</button>
          <div class="lang-toggle" aria-label="Language">
            <button id="lang-zh" type="button" data-lang="zh">中文</button>
            <button id="lang-en" type="button" data-lang="en">EN</button>
          </div>
        </div>
      </header>

      <main class="content">
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
            <div class="metric-label" data-i18n="registeredCli">已注册 CLI</div>
            <div class="metric-value"><span id="tool-count">0</span></div>
          </div>
          <div class="metric">
            <div class="metric-label" data-i18n="tenants">租户</div>
            <div class="metric-value"><span id="tenant-count">0</span></div>
          </div>
        </div>

        <section class="view active" id="terminal-view">
          <div class="terminal-grid">
            <section class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">&gt;_</span>
                  <div>
                    <h2 data-i18n="terminal">终端</h2>
                    <p id="terminal-context">team-a / repo-main / team-default</p>
                  </div>
                </div>
                <div class="panel-actions">
                  <span class="badge warn" id="connection-badge"><span class="led" id="terminal-led"></span><span id="terminal-state">disconnected</span></span>
                  <button class="primary" id="connect-terminal" type="button" data-i18n="connect">连接</button>
                  <button class="ghost" id="clear-terminal" type="button" data-i18n="clear">清屏</button>
                  <button class="ghost" id="disconnect-terminal" type="button" data-i18n="disconnect">断开</button>
                </div>
              </div>

              <div class="context-bar">
                <div>
                  <label for="tenant" data-i18n="tenant">租户</label>
                  <select id="tenant"></select>
                </div>
                <div>
                  <label for="workspace" data-i18n="workspace">工作区</label>
                  <input id="workspace" value="repo-main">
                </div>
                <div>
                  <label for="profile" data-i18n="credentialProfile">凭据配置</label>
                  <input id="profile" value="team-default" list="profile-options">
                  <datalist id="profile-options"></datalist>
                </div>
              </div>

              <div class="terminal-body">
                <div class="xterm-frame">
                  <div id="terminal-container"></div>
                </div>
              </div>
            </section>

            <aside class="side-stack">
              <section class="panel">
                <div class="panel-header">
                  <div class="panel-title">
                    <span class="nav-icon">QL</span>
                    <div>
                      <h2 data-i18n="quickLogin">快捷登录</h2>
                      <p data-i18n="quickLoginDesc">必要时自动连接终端，然后执行登录命令。</p>
                    </div>
                  </div>
                </div>
                <div class="quick-grid" id="login-shortcuts"></div>
              </section>

              <section class="panel">
                <div class="panel-header">
                  <div class="panel-title">
                    <span class="nav-icon">CL</span>
                    <div>
                      <h2 data-i18n="cliManager">CLI 管理</h2>
                      <p data-i18n="cliManagerDesc">已注册命令和凭据目录。</p>
                    </div>
                  </div>
                  <button class="ghost" type="button" data-view-jump="tools-view" data-i18n="manage">管理</button>
                </div>
                <div class="tool-cards" id="cli-cards"></div>
              </section>
            </aside>
          </div>
        </section>

        <section class="view" id="tools-view">
          <div class="manager-grid">
            <section class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">CL</span>
                  <div>
                    <h2 data-i18n="cliManager">CLI 管理</h2>
                    <p data-i18n="toolsSubtitle">添加、更新或删除共享运行时使用的 CLI 定义。</p>
                  </div>
                </div>
              </div>
              <div class="manager-tools" id="tools">
                <div class="empty">Loading CLI tools</div>
              </div>
            </section>

            <aside class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">IN</span>
                  <div>
                    <h2 data-i18n="installCli">安装 CLI</h2>
                    <p data-i18n="installCliDesc">从官方安装源执行推荐命令。</p>
                  </div>
                </div>
              </div>
              <div class="install-grid" id="install-options"></div>
            </aside>
          </div>
        </section>

        <section class="view" id="tenants-view">
          <section class="panel">
            <div class="panel-header">
              <div class="panel-title">
                <span class="nav-icon">TN</span>
                <div>
                  <h2 data-i18n="tenants">租户</h2>
                  <p data-i18n="tenantSubtitle">查看 token 推导出的工具、工作区、凭据配置和终端访问边界。</p>
                </div>
              </div>
            </div>
            <div class="tenant-grid" id="tenants-list"></div>
          </section>
        </section>
      </main>
    </div>
  </div>

  <div class="toast" id="toast"></div>

  <script src="/assets/xterm/xterm.js"></script>
  <script>
    const messages = {
      zh: {
        brandSubtitle: 'CLI 控制平面',
        navTerminal: '终端',
        navTools: 'CLI 管理',
        navTenants: '租户',
        runtime: '运行时',
        token: 'Token',
        refresh: '刷新',
        health: '健康状态',
        ready: '就绪状态',
        registeredCli: '已注册 CLI',
        tenants: '租户',
        terminal: '终端',
        terminalSubtitle: '用于登录认证、CLI 调试和手动维护的完整 shell。',
        cliManager: 'CLI 管理',
        cliSubtitle: '安装、更新和删除共享运行时使用的 CLI。',
        tenantsSubtitle: '查看租户边界、凭据配置和终端权限。',
        connect: '连接',
        clear: '清屏',
        disconnect: '断开',
        tenant: '租户',
        workspace: '工作区',
        credentialProfile: '凭据配置',
        quickLogin: '快捷登录',
        quickLoginDesc: '必要时自动连接终端，然后执行登录命令。',
        cliManagerDesc: '已注册命令和凭据目录。',
        toolsSubtitle: '添加、更新或删除共享运行时使用的 CLI 定义。',
        manage: '管理',
        name: '名称',
        version: '版本',
        path: '路径',
        credentialEnv: '凭据环境变量',
        credentialSubdir: '凭据子目录',
        credentialHome: '凭据目录',
        commandPath: '命令路径',
        loginCommand: '登录命令',
        installCli: '安装 CLI',
        installCliDesc: '从官方安装源执行推荐命令。',
        officialSource: '官方来源',
        installCommand: '安装命令',
        runInstall: '终端安装',
        verifyInstall: '验证',
        tenantSubtitle: '查看 token 推导出的工具、工作区、凭据配置和终端访问边界。',
        noTools: '暂无已注册 CLI',
        noTenants: '暂无租户配置',
        delete: '删除',
        registered: '已注册',
        readyLabel: '可用',
        addLabel: '待添加',
        terminalAllowed: '允许终端',
        terminalBlocked: '禁止终端',
        subjects: '主体',
        tools: '工具',
        workspaces: '工作区',
        profiles: '凭据配置',
        connected: '已连接',
        disconnected: '未连接',
        connecting: '连接中',
        connectionError: '连接错误',
        exited: '已退出',
        terminalWelcome: 'Agent Runtime 终端已就绪。点击连接打开 shell，或使用右侧快捷登录。',
        terminalConnecting: '正在连接终端...',
        runtimeReady: 'Health OK',
        runtimeOffline: 'Offline',
        refreshed: '状态已刷新',
        cliSaved: 'CLI 已保存',
        cliDeleted: 'CLI 已删除'
      },
      en: {
        brandSubtitle: 'CLI control plane',
        navTerminal: 'Terminal',
        navTools: 'CLI Manager',
        navTenants: 'Tenants',
        runtime: 'Runtime',
        token: 'Token',
        refresh: 'Refresh',
        health: 'Health',
        ready: 'Ready',
        registeredCli: 'Registered CLI',
        tenants: 'Tenants',
        terminal: 'Terminal',
        terminalSubtitle: 'Full shell for login, CLI diagnostics, and manual maintenance.',
        cliManager: 'CLI Manager',
        cliSubtitle: 'Install, update, and remove CLI entries used by shared runtimes.',
        tenantsSubtitle: 'Inspect tenant boundaries, credential profiles, and terminal permissions.',
        connect: 'Connect',
        clear: 'Clear',
        disconnect: 'Disconnect',
        tenant: 'Tenant',
        workspace: 'Workspace',
        credentialProfile: 'Credential Profile',
        quickLogin: 'Quick Login',
        quickLoginDesc: 'Connects the terminal when needed, then runs the login command.',
        cliManagerDesc: 'Registered commands and credential homes.',
        toolsSubtitle: 'Add, update, or remove CLI definitions used by the shared runtime.',
        manage: 'Manage',
        name: 'Name',
        version: 'Version',
        path: 'Path',
        credentialEnv: 'Credential Env',
        credentialSubdir: 'Credential Subdir',
        credentialHome: 'Credential Home',
        commandPath: 'Command Path',
        loginCommand: 'Login Command',
        installCli: 'Install CLI',
        installCliDesc: 'Run the recommended command from the official install source.',
        officialSource: 'Official Source',
        installCommand: 'Install Command',
        runInstall: 'Install in Terminal',
        verifyInstall: 'Verify',
        tenantSubtitle: 'Inspect token-derived tool, workspace, credential profile, and terminal boundaries.',
        noTools: 'No CLI tools registered',
        noTenants: 'No tenants configured',
        delete: 'Delete',
        registered: 'Registered',
        readyLabel: 'Ready',
        addLabel: 'Add',
        terminalAllowed: 'Terminal allowed',
        terminalBlocked: 'Terminal blocked',
        subjects: 'Subjects',
        tools: 'Tools',
        workspaces: 'Workspaces',
        profiles: 'Profiles',
        connected: 'Connected',
        disconnected: 'Disconnected',
        connecting: 'Connecting',
        connectionError: 'Connection error',
        exited: 'Exited',
        terminalWelcome: 'Agent Runtime terminal is ready. Click Connect to open a shell, or use Quick Login.',
        terminalConnecting: 'Connecting terminal...',
        runtimeReady: 'Health OK',
        runtimeOffline: 'Offline',
        refreshed: 'Runtime status refreshed',
        cliSaved: 'CLI saved',
        cliDeleted: 'CLI deleted'
      }
    };

    const savedLanguage = localStorage.getItem('agent-runtime-lang');
    const state = {
      lang: savedLanguage || 'zh',
      view: 'terminal-view',
      tools: [],
      tenants: [],
      ws: null,
      connected: false,
      statusKey: 'disconnected',
      term: null,
      dataDisposable: null,
      pendingCommand: ''
    };
    const tokenInput = document.getElementById('token');
    tokenInput.value = localStorage.getItem('agent-runtime-token') || 'dev-token';

    const loginCommands = [
      { label: 'Claude Code', command: 'claude', tool: 'claude' },
      { label: 'Codex', command: 'codex login', tool: 'codex' },
      { label: 'Gemini', command: 'gemini', tool: 'gemini' },
      { label: 'iFlow', command: 'iflow', tool: 'iflow' },
      { label: 'OpenCode', command: 'opencode auth login', tool: 'opencode' },
      { label: 'Kimi', command: 'kimi', tool: 'kimi' },
      { label: 'Qoder', command: 'qodercli', tool: 'qoder' }
    ];

    const installSources = [
      {
        label: 'Claude Code',
        tool: 'claude',
        provider: 'Anthropic',
        docs: 'https://docs.anthropic.com/en/docs/claude-code/quickstart',
        command: 'curl -fsSL https://claude.ai/install.sh | bash',
        verify: 'claude --version'
      },
      {
        label: 'Codex',
        tool: 'codex',
        provider: 'OpenAI',
        docs: 'https://github.com/openai/codex',
        command: 'npm install -g @openai/codex',
        verify: 'codex --version'
      },
      {
        label: 'Gemini CLI',
        tool: 'gemini',
        provider: 'Google',
        docs: 'https://github.com/google-gemini/gemini-cli',
        command: 'npm install -g @google/gemini-cli',
        verify: 'gemini --version'
      },
      {
        label: 'OpenCode',
        tool: 'opencode',
        provider: 'SST',
        docs: 'https://opencode.ai/download',
        command: 'curl -fsSL https://opencode.ai/install | bash',
        verify: 'opencode --version'
      },
      {
        label: 'iFlow',
        tool: 'iflow',
        provider: 'iFlow',
        docs: 'https://platform.iflow.cn/cli/quickstart',
        command: 'bash -c "$(curl -fsSL https://gitee.com/iflow-ai/iflow-cli/raw/main/install.sh)"',
        verify: 'iflow --version'
      },
      {
        label: 'Kimi',
        tool: 'kimi',
        provider: 'Moonshot AI',
        docs: 'https://www.kimi.com/code/docs/en/kimi-code-cli/getting-started.html',
        command: 'curl -LsSf https://code.kimi.com/install.sh | bash',
        verify: 'kimi --version'
      },
      {
        label: 'Qoder',
        tool: 'qoder',
        provider: 'Qoder',
        docs: 'https://docs.qoder.com/cli/quick-start',
        command: 'curl -fsSL https://qoder.com/install | bash',
        verify: 'qodercli --version'
      }
    ];

    function $(id) { return document.getElementById(id); }
    function t(key) { return (messages[state.lang] && messages[state.lang][key]) || messages.en[key] || key; }
    function escapeHTML(value) {
      return String(value == null ? '' : value).replace(/[&<>"']/g, function(char) {
        return ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char];
      });
    }
    function join(values) { return values && values.length ? values.join(', ') : '-'; }
    function showToast(message) {
      const toast = $('toast');
      toast.textContent = message;
      toast.classList.add('show');
      window.clearTimeout(showToast.timer);
      showToast.timer = window.setTimeout(function() { toast.classList.remove('show'); }, 3200);
    }

    async function api(path, options) {
      const response = await fetch(path, options || {});
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
      document.querySelectorAll('[data-i18n]').forEach(function(el) {
        el.textContent = t(el.dataset.i18n);
      });
      $('lang-zh').classList.toggle('active', state.lang === 'zh');
      $('lang-en').classList.toggle('active', state.lang === 'en');
      updatePageTitle();
      renderLoginShortcuts();
      renderTools();
      renderToolCards();
      renderInstallOptions();
      renderTenants();
      setConnected(state.connected, state.statusKey);
      $('rail-health').textContent = $('health').textContent === 'ok' ? t('runtimeReady') : t('runtimeOffline');
      if (state.term && state.term.buffer.active.length <= 2) {
        state.term.clear();
        state.term.writeln(t('terminalWelcome'));
      }
    }

    function updatePageTitle() {
      const titles = {
        'terminal-view': [t('terminal'), t('terminalSubtitle')],
        'tools-view': [t('cliManager'), t('cliSubtitle')],
        'tenants-view': [t('tenants'), t('tenantsSubtitle')]
      };
      const value = titles[state.view] || titles['terminal-view'];
      $('page-heading').textContent = value[0];
      $('page-subtitle').textContent = value[1];
    }

    function initTerminal() {
      if (!window.Terminal) {
        showToast('xterm.js not loaded');
        return;
      }
      state.term = new window.Terminal({
        cursorBlink: true,
        convertEol: true,
        fontFamily: 'JetBrains Mono, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
        fontSize: 13,
        lineHeight: 1.25,
        scrollback: 5000,
        theme: {
          background: '#05070a',
          foreground: '#d8dee9',
          cursor: '#42ff9c',
          selectionBackground: '#284b3c',
          black: '#05070a',
          red: '#ff5c73',
          green: '#37e681',
          yellow: '#f5bc4f',
          blue: '#4c8dff',
          magenta: '#b68cff',
          cyan: '#28e0d4',
          white: '#d8dee9'
        }
      });
      state.term.open($('terminal-container'));
      state.term.writeln(t('terminalWelcome'));
      state.dataDisposable = state.term.onData(function(data) {
        if (state.ws && state.ws.readyState === WebSocket.OPEN) {
          state.ws.send(JSON.stringify({ type: 'input', data: data }));
        }
      });
      if (window.ResizeObserver) {
        new ResizeObserver(function() { resizeTerminal(); }).observe($('terminal-container'));
      } else {
        window.addEventListener('resize', resizeTerminal);
      }
      window.setTimeout(resizeTerminal, 0);
    }

    function resizeTerminal() {
      if (!state.term) return { cols: 120, rows: 32 };
      const container = $('terminal-container');
      const rect = container.getBoundingClientRect();
      const cols = Math.max(40, Math.floor(rect.width / 9));
      const rows = Math.max(12, Math.floor(rect.height / 18));
      state.term.resize(cols, rows);
      if (state.ws && state.ws.readyState === WebSocket.OPEN) {
        state.ws.send(JSON.stringify({ type: 'resize', cols: cols, rows: rows }));
      }
      return { cols: cols, rows: rows };
    }

    async function refresh() {
      localStorage.setItem('agent-runtime-token', tokenInput.value.trim());
      const results = await Promise.allSettled([
        api('/api/health'),
        api('/api/ready'),
        api('/api/status'),
        api('/api/tools'),
        api('/api/tenants')
      ]);
      const health = results[0];
      const ready = results[1];
      const status = results[2];
      const tools = results[3];
      const tenants = results[4];
      setMetric('health', 'health-led', health.status === 'fulfilled' ? health.value.status : 'error');
      setMetric('ready', 'ready-led', ready.status === 'fulfilled' ? ready.value.status : 'error');
      if (status.status === 'fulfilled') {
        $('tool-count').textContent = status.value.tools == null ? '0' : String(status.value.tools);
        $('tenant-count').textContent = status.value.tenants == null ? '0' : String(status.value.tenants);
      }
      state.tools = tools.status === 'fulfilled' ? tools.value.tools : [];
      state.tenants = tenants.status === 'fulfilled' ? tenants.value.tenants : [];
      $('rail-health').textContent = health.status === 'fulfilled' ? t('runtimeReady') : t('runtimeOffline');
      renderTerminalOptions();
      renderLoginShortcuts();
      renderTools();
      renderToolCards();
      renderInstallOptions();
      renderTenants();
      updateContextLabels();
    }

    function setMetric(textID, ledID, value) {
      $(textID).textContent = value;
      const led = $(ledID);
      led.classList.remove('ok', 'bad');
      if (value === 'ok' || value === 'ready') led.classList.add('ok');
      if (value === 'error') led.classList.add('bad');
    }

    function renderTerminalOptions() {
      const tenantSelect = $('tenant');
      const currentTenant = tenantSelect.value;
      tenantSelect.innerHTML = state.tenants.map(function(tenant) {
        return '<option value="' + escapeHTML(tenant.id) + '">' + escapeHTML(tenant.id) + '</option>';
      }).join('');
      if (currentTenant) tenantSelect.value = currentTenant;
      if (!tenantSelect.value && state.tenants[0]) tenantSelect.value = state.tenants[0].id;
      updateProfileOptions();
    }

    function updateProfileOptions() {
      const tenant = state.tenants.find(function(item) { return item.id === $('tenant').value; });
      const profiles = tenant && tenant.credential_profiles ? tenant.credential_profiles : [];
      $('profile-options').innerHTML = profiles.map(function(profile) {
        return '<option value="' + escapeHTML(profile) + '"></option>';
      }).join('');
      if (profiles.length && !$('profile').value) $('profile').value = profiles[0];
      const workspaces = tenant && tenant.workspace_patterns ? tenant.workspace_patterns : [];
      if (!$('workspace').value && workspaces[0]) $('workspace').value = workspaces[0].replace('*', 'main');
      updateContextLabels();
    }

    function updateContextLabels() {
      const context = ($('tenant').value || '-') + ' / ' + ($('workspace').value || '-');
      $('terminal-context').textContent = context + ' / ' + ($('profile').value || '-');
      $('rail-context').textContent = context;
    }

    function terminalURL() {
      const size = resizeTerminal();
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const params = new URLSearchParams({
        token: tokenInput.value.trim(),
        tenant: $('tenant').value,
        workspace: $('workspace').value.trim(),
        credential_profile: $('profile').value.trim(),
        cols: String(size.cols),
        rows: String(size.rows)
      });
      return protocol + '//' + window.location.host + '/api/v1/terminal/ws?' + params.toString();
    }

    function connectTerminal() {
      if (!state.term) initTerminal();
      if (!state.term) return;
      if (state.ws) state.ws.close();
      state.term.clear();
      state.term.writeln(t('terminalConnecting'));
      setConnected(false, 'connecting');
      localStorage.setItem('agent-runtime-token', tokenInput.value.trim());
      const ws = new WebSocket(terminalURL());
      state.ws = ws;
      ws.onopen = function() {
        setConnected(true, 'connected');
        state.term.clear();
        state.term.focus();
        resizeTerminal();
        if (state.pendingCommand) {
          const command = state.pendingCommand;
          state.pendingCommand = '';
          window.setTimeout(function() { sendTerminal(command); }, 180);
        }
      };
      ws.onmessage = function(event) {
        try {
          const payload = JSON.parse(event.data);
          if (payload.type === 'output') state.term.write(payload.data || '');
          if (payload.type === 'error') state.term.writeln('\r\n[terminal error] ' + (payload.data || 'unknown error'));
          if (payload.type === 'exit') setConnected(false, 'exited');
        } catch (err) {
          state.term.write(String(event.data));
        }
      };
      ws.onclose = function() {
        if (state.ws === ws) state.ws = null;
        setConnected(false, 'disconnected');
      };
      ws.onerror = function() { setConnected(false, 'connectionError'); };
    }

    function disconnectTerminal() {
      if (state.ws) state.ws.close();
      state.ws = null;
      setConnected(false, 'disconnected');
    }

    function setConnected(connected, statusKey) {
      state.connected = connected;
      state.statusKey = statusKey || (connected ? 'connected' : 'disconnected');
      $('terminal-state').textContent = t(state.statusKey);
      $('connect-terminal').disabled = connected;
      $('disconnect-terminal').disabled = !connected;
      $('terminal-led').classList.remove('ok', 'bad');
      if (connected) $('terminal-led').classList.add('ok');
      if (state.statusKey === 'connectionError') $('terminal-led').classList.add('bad');
      $('connection-badge').classList.toggle('ok', connected);
      $('connection-badge').classList.toggle('warn', !connected);
      renderLoginShortcuts();
    }

    function sendTerminal(data) {
      if (!state.ws || state.ws.readyState !== WebSocket.OPEN) return;
      state.ws.send(JSON.stringify({ type: 'input', data: data }));
    }

    function runCommand(command) {
      const payload = command + '\r';
      if (state.connected) {
        sendTerminal(payload);
        state.term.focus();
        return;
      }
      state.pendingCommand = payload;
      connectTerminal();
    }

    function renderLoginShortcuts() {
      const knownTools = new Set(state.tools.map(function(tool) { return tool.name; }));
      const container = $('login-shortcuts');
      container.innerHTML = loginCommands.map(function(item) {
        const registered = knownTools.has(item.tool);
        return '<button class="quick-button" type="button" data-login-command="' + escapeHTML(item.command) + '">' +
          '<span><strong>' + escapeHTML(item.label) + '</strong><span>' + escapeHTML(item.command) + '</span></span>' +
          '<span class="badge ' + (registered ? 'ok' : 'warn') + '">' + (registered ? t('readyLabel') : t('addLabel')) + '</span>' +
        '</button>';
      }).join('');
      container.querySelectorAll('[data-login-command]').forEach(function(button) {
        button.addEventListener('click', function() { runCommand(button.dataset.loginCommand); });
      });
    }

    function renderToolCards() {
      const container = $('cli-cards');
      if (!state.tools.length) {
        container.innerHTML = '<div class="empty">' + escapeHTML(t('noTools')) + '</div>';
        return;
      }
      container.innerHTML = state.tools.slice(0, 6).map(function(tool) {
        return toolCardHTML(tool, false);
      }).join('');
    }

    function renderTools() {
      const container = $('tools');
      if (!state.tools.length) {
        container.innerHTML = '<div class="empty">' + escapeHTML(t('noTools')) + '</div>';
        return;
      }
      container.innerHTML = state.tools.map(function(tool) {
        return toolCardHTML(tool, true);
      }).join('');
      container.querySelectorAll('[data-login-command]').forEach(function(button) {
        button.addEventListener('click', function() { runCommand(button.dataset.loginCommand); });
      });
      container.querySelectorAll('[data-delete-tool]').forEach(function(button) {
        button.addEventListener('click', async function() {
          await api('/api/tools/' + encodeURIComponent(button.dataset.deleteTool), { method: 'DELETE' });
          showToast(t('cliDeleted') + ': ' + button.dataset.deleteTool);
          await refresh();
        });
      });
    }

    function renderInstallOptions() {
      const container = $('install-options');
      if (!container) return;
      const knownTools = new Set(state.tools.map(function(tool) { return tool.name; }));
      container.innerHTML = installSources.map(function(source) {
        const registered = knownTools.has(source.tool);
        return '<article class="install-card">' +
          '<div class="install-card-head">' +
            '<div class="install-title">' +
              '<strong>' + escapeHTML(source.label) + '</strong>' +
              '<span>' + escapeHTML(source.provider) + ' · ' + escapeHTML(t('officialSource')) + '</span>' +
            '</div>' +
            '<span class="badge ' + (registered ? 'ok' : 'warn') + '">' + escapeHTML(registered ? t('registered') : t('addLabel')) + '</span>' +
          '</div>' +
          '<div class="install-command">' + escapeHTML(source.command) + '</div>' +
          '<div class="install-actions">' +
            '<a href="' + escapeHTML(source.docs) + '" target="_blank" rel="noopener noreferrer">' + escapeHTML(t('officialSource')) + '</a>' +
            '<button class="primary" type="button" data-install-command="' + escapeHTML(source.command) + '">' + escapeHTML(t('runInstall')) + '</button>' +
            '<button class="ghost" type="button" data-install-command="' + escapeHTML(source.verify) + '">' + escapeHTML(t('verifyInstall')) + '</button>' +
            '<button class="ghost" type="button" data-login-command="' + escapeHTML((loginCommands.find(function(item) { return item.tool === source.tool; }) || {}).command || source.tool) + '">' + escapeHTML(t('quickLogin')) + '</button>' +
          '</div>' +
        '</article>';
      }).join('');
      container.querySelectorAll('[data-install-command]').forEach(function(button) {
        button.addEventListener('click', function() { runCommand(button.dataset.installCommand); });
      });
      container.querySelectorAll('[data-login-command]').forEach(function(button) {
        button.addEventListener('click', function() { runCommand(button.dataset.loginCommand); });
      });
    }

    function toolCardHTML(tool, manageable) {
      const login = loginCommands.find(function(item) { return item.tool === tool.name; });
      const credentialHome = (tool.credential_env || 'HOME') + ' -> ' + (tool.credential_subdir || '.');
      let html = '<div class="tool-card">' +
        '<div class="tool-card-head">' +
          '<span class="tool-name">' + escapeHTML(tool.name) + '</span>' +
          '<span class="badge ok"><span class="led ok"></span>' + escapeHTML(t('registered')) + '</span>' +
        '</div>' +
        '<div class="tool-meta">' +
          '<div><strong>' + escapeHTML(t('version')) + '</strong><span>' + escapeHTML(tool.version || '-') + '</span></div>' +
          '<div><strong>' + escapeHTML(t('commandPath')) + '</strong><span class="mono">' + escapeHTML(tool.path) + '</span></div>' +
          '<div><strong>' + escapeHTML(t('credentialHome')) + '</strong><span class="mono">' + escapeHTML(credentialHome) + '</span></div>';
      if (login) {
        html += '<div><strong>' + escapeHTML(t('loginCommand')) + '</strong><span class="mono">' + escapeHTML(login.command) + '</span></div>';
      }
      html += '</div>';
      if (manageable) {
        html += '<div class="tool-actions">' +
          (login ? '<button class="ghost" type="button" data-login-command="' + escapeHTML(login.command) + '">' + escapeHTML(t('quickLogin')) + '</button>' : '') +
          '<button class="danger" type="button" data-delete-tool="' + escapeHTML(tool.name) + '">' + escapeHTML(t('delete')) + '</button>' +
        '</div>';
      }
      html += '</div>';
      return html;
    }

    function renderTenants() {
      const container = $('tenants-list');
      if (!state.tenants.length) {
        container.innerHTML = '<div class="empty">' + escapeHTML(t('noTenants')) + '</div>';
        return;
      }
      container.innerHTML = state.tenants.map(function(tenant) {
        return '<article class="tenant-card">' +
          '<div class="tenant-name"><span>' + escapeHTML(tenant.id) + '</span>' +
          '<span class="badge ' + (tenant.allow_terminal ? 'ok' : 'warn') + '">' + escapeHTML(tenant.allow_terminal ? t('terminalAllowed') : t('terminalBlocked')) + '</span></div>' +
          '<div class="kv">' +
            '<div><strong>' + escapeHTML(t('subjects')) + '</strong><span>' + escapeHTML(join(tenant.subjects)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('tools')) + '</strong><span>' + escapeHTML(join(tenant.allowed_tools)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('workspaces')) + '</strong><span>' + escapeHTML(join(tenant.workspace_patterns)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('profiles')) + '</strong><span>' + escapeHTML(join(tenant.credential_profiles)) + '</span></div>' +
          '</div>' +
        '</article>';
      }).join('');
    }

    function switchView(viewID) {
      state.view = viewID;
      document.querySelectorAll('.nav-button').forEach(function(item) {
        item.classList.toggle('active', item.dataset.view === viewID);
      });
      document.querySelectorAll('.view').forEach(function(item) {
        item.classList.toggle('active', item.id === viewID);
      });
      updatePageTitle();
      if (viewID === 'terminal-view') window.setTimeout(resizeTerminal, 0);
    }

    document.querySelectorAll('[data-view]').forEach(function(button) {
      button.addEventListener('click', function() { switchView(button.dataset.view); });
    });
    document.querySelectorAll('[data-view-jump]').forEach(function(button) {
      button.addEventListener('click', function() { switchView(button.dataset.viewJump); });
    });
    document.querySelectorAll('[data-lang]').forEach(function(button) {
      button.addEventListener('click', function() {
        state.lang = button.dataset.lang;
        localStorage.setItem('agent-runtime-lang', state.lang);
        applyLanguage();
      });
    });
    $('refresh').addEventListener('click', function() {
      refresh().then(function() { showToast(t('refreshed')); }).catch(function(err) { showToast(err.message); });
    });
    $('tenant').addEventListener('change', updateProfileOptions);
    $('workspace').addEventListener('input', updateContextLabels);
    $('profile').addEventListener('input', updateContextLabels);
    $('connect-terminal').addEventListener('click', connectTerminal);
    $('disconnect-terminal').addEventListener('click', disconnectTerminal);
    $('clear-terminal').addEventListener('click', function() { if (state.term) state.term.clear(); });

    initTerminal();
    applyLanguage();
    refresh().catch(function(err) {
      setMetric('health', 'health-led', 'error');
      setMetric('ready', 'ready-led', 'error');
      showToast(err.message);
    });
  </script>
</body>
</html>`
