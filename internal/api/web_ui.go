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
      --panel: rgba(10, 15, 23, 0.92);
      --panel-2: #0d141e;
      --card: #111923;
      --card-2: #151f2b;
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
    code, .mono { font-family: var(--mono); }
    .app {
      min-height: 100vh;
      display: grid;
      grid-template-columns: 232px minmax(0, 1fr);
    }
    .rail {
      position: sticky;
      top: 0;
      height: 100vh;
      padding: 20px 14px;
      border-right: 1px solid var(--line);
      background: rgba(4, 8, 13, 0.88);
      backdrop-filter: blur(18px);
      display: flex;
      flex-direction: column;
      gap: 18px;
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
    .nav { display: grid; gap: 8px; }
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
    .nav-button.active { border-color: rgba(37, 215, 207, 0.6); background: rgba(37, 215, 207, 0.1); color: var(--text); }
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
    .rail-card {
      margin-top: auto;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(10, 16, 24, 0.74);
    }
    .rail-card-title { display: flex; justify-content: space-between; gap: 8px; margin-bottom: 9px; color: var(--muted); font-size: 12px; font-weight: 850; text-transform: uppercase; }
    .workspace { min-width: 0; padding: 16px; }
    .topbar {
      display: grid;
      grid-template-columns: minmax(240px, 1fr) minmax(520px, 780px);
      gap: 16px;
      align-items: center;
      padding: 12px 14px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(8, 13, 20, 0.88);
      backdrop-filter: blur(18px);
      box-shadow: 0 20px 70px rgba(0, 0, 0, 0.3);
    }
    .page-title h1 { margin: 0; font-size: 20px; line-height: 1.1; }
    .page-title p { margin: 5px 0 0; color: var(--muted); font-size: 13px; }
    .top-actions { display: grid; grid-template-columns: minmax(260px, 1fr) 78px 78px 92px 120px; gap: 8px; align-items: center; }
    .login-fields { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
    .login-field { display: grid; grid-template-columns: 82px minmax(0, 1fr); gap: 8px; align-items: center; }
    .login-field span { color: var(--muted); font-size: 12px; font-weight: 850; text-transform: uppercase; }
    .lang-toggle { display: grid; grid-template-columns: 1fr 1fr; gap: 4px; padding: 4px; border: 1px solid var(--line); border-radius: 8px; background: rgba(3, 7, 12, 0.72); }
    .lang-toggle button { min-height: 28px; border: 0; background: transparent; color: var(--muted); font-size: 12px; }
    .lang-toggle button.active { background: rgba(37, 215, 207, 0.15); color: var(--cyan); }
    .content { margin-top: 16px; }
    .view { display: none; }
    .view.active { display: block; }
    .status-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 10px; margin-bottom: 16px; }
    .metric { min-height: 66px; padding: 11px 12px; border: 1px solid var(--line); border-radius: 8px; background: rgba(10, 17, 26, 0.74); }
    .metric-label { color: var(--muted); font-size: 11px; font-weight: 850; text-transform: uppercase; }
    .metric-value { margin-top: 8px; display: flex; align-items: center; gap: 8px; font-size: 15px; font-weight: 850; overflow-wrap: anywhere; }
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
    .panel { border: 1px solid var(--line); border-radius: 8px; background: var(--panel); box-shadow: 0 24px 90px rgba(0, 0, 0, 0.28); overflow: hidden; }
    .panel-header { min-height: 58px; display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 13px 14px; border-bottom: 1px solid var(--line); }
    .panel-title { display: flex; align-items: center; gap: 10px; min-width: 0; }
    .panel-title h2 { margin: 0; font-size: 15px; }
    .panel-title p { margin: 2px 0 0; color: var(--muted); font-size: 12px; overflow-wrap: anywhere; }
    .panel-actions { display: flex; gap: 8px; flex-wrap: wrap; justify-content: flex-end; }
    .terminal-grid { display: grid; grid-template-columns: minmax(0, 1fr) 360px; gap: 16px; align-items: start; }
    .context-bar { display: grid; grid-template-columns: 150px minmax(160px, 1fr) minmax(160px, 1fr); gap: 10px; padding: 14px; border-bottom: 1px solid var(--line); background: rgba(5, 9, 14, 0.35); align-items: end; }
    .terminal-body { height: clamp(520px, calc(100vh - 336px), 780px); min-height: 520px; padding: 12px; }
    .xterm-frame { height: 100%; overflow: hidden; border: 1px solid #172332; border-radius: 8px; background: #05070a; box-shadow: inset 0 0 42px rgba(37, 215, 207, 0.05); }
    #terminal-container { height: 100%; width: 100%; padding: 8px; }
    .xterm { height: 100%; padding: 2px; }
    .manager-body { padding: 12px; display: grid; gap: 10px; }
    .cli-card {
      min-height: 116px;
      display: grid;
      grid-template-columns: 58px minmax(0, 1fr) auto;
      gap: 12px;
      align-items: center;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: linear-gradient(180deg, rgba(21, 31, 43, 0.98), rgba(13, 20, 30, 0.98));
    }
    .cli-logo { width: 46px; height: 46px; display: grid; place-items: center; overflow: hidden; border: 1px solid rgba(44, 66, 96, 0.78); border-radius: 8px; background: rgba(3, 7, 12, 0.62); color: var(--cyan); font-size: 13px; font-weight: 900; }
    .cli-logo img { max-width: 36px; max-height: 36px; object-fit: contain; }
    .cli-logo span { width: 100%; height: 100%; display: none; place-items: center; }
    .cli-name { font-size: 16px; font-weight: 850; }
    .cli-version { margin-top: 2px; color: var(--muted); font-family: var(--mono); font-size: 13px; }
    .cli-health { margin-top: 10px; }
    .cli-actions { display: grid; gap: 8px; min-width: 118px; }
    .cli-actions a { min-height: 36px; display: grid; place-items: center; padding: 8px 10px; border: 1px solid var(--line-2); border-radius: 8px; background: rgba(6, 10, 15, 0.7); color: var(--muted); font-size: 13px; font-weight: 760; text-decoration: none; }
    .cli-actions a:hover { border-color: rgba(37, 215, 207, 0.72); }
    .tenant-layout { display: grid; grid-template-columns: 1fr; gap: 16px; align-items: start; }
    .form-stack, .tenant-list, .file-body { display: grid; gap: 10px; padding: 14px; }
    .field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
    .tenant-card, .user-card { padding: 12px; border: 1px solid var(--line); border-radius: 8px; background: rgba(7, 12, 18, 0.72); }
    .tenant-card-head, .user-card-head { display: flex; justify-content: space-between; gap: 10px; align-items: center; margin-bottom: 10px; font-weight: 850; }
    .tenant-card-actions, .user-card-actions { margin-top: 10px; display: flex; gap: 8px; flex-wrap: wrap; }
    .kv { display: grid; gap: 8px; color: var(--muted); font-size: 12px; }
    .kv div { display: grid; grid-template-columns: 118px minmax(0, 1fr); gap: 8px; }
    .kv strong { color: var(--faint); font-size: 11px; font-weight: 850; text-transform: uppercase; }
    .file-controls { display: grid; grid-template-columns: 150px 150px minmax(180px, 1fr) auto auto; gap: 10px; align-items: end; }
    .file-path { padding: 9px 10px; border: 1px solid var(--line); border-radius: 8px; background: rgba(3, 7, 12, 0.78); color: var(--muted); font-family: var(--mono); font-size: 12px; overflow-wrap: anywhere; }
    .file-browser { height: clamp(420px, calc(100vh - 430px), 680px); min-height: 420px; display: grid; grid-template-columns: minmax(240px, 42%) minmax(0, 1fr); overflow: hidden; border: 1px solid var(--line); border-radius: 8px; background: rgba(4, 8, 13, 0.72); }
    .file-list { overflow: auto; border-right: 1px solid var(--line); }
    .file-row { display: grid; grid-template-columns: minmax(0, 1fr) 86px; gap: 10px; align-items: center; min-height: 36px; padding: 8px 10px; border-bottom: 1px solid rgba(30, 43, 57, 0.62); background: transparent; cursor: pointer; }
    .file-row:hover, .file-row.active { background: rgba(37, 215, 207, 0.08); }
    .file-name { display: flex; gap: 8px; align-items: center; min-width: 0; }
    .file-name button { min-height: 26px; padding: 0; border: 0; background: transparent; color: var(--text); text-align: left; font-weight: 760; overflow-wrap: anywhere; cursor: pointer; }
    .file-meta { color: var(--faint); font-family: var(--mono); font-size: 11px; }
    .file-viewer { min-width: 0; display: flex; flex-direction: column; }
    .file-viewbar { min-height: 38px; padding: 10px 12px; border-bottom: 1px solid var(--line); color: var(--muted); font-family: var(--mono); font-size: 12px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    .file-preview { flex: 1; margin: 0; padding: 12px; overflow: auto; white-space: pre; color: #d9e8f8; font-family: var(--mono); font-size: 12px; line-height: 1.45; }
    .empty { padding: 18px; color: var(--muted); border: 1px dashed var(--line-2); border-radius: 8px; background: rgba(255, 255, 255, 0.025); }
    .hidden { display: none !important; }
    .toast { position: fixed; right: 18px; bottom: 18px; z-index: 10; display: none; max-width: 420px; padding: 12px 14px; border: 1px solid rgba(37, 215, 207, 0.42); border-radius: 8px; background: rgba(5, 10, 16, 0.96); color: var(--text); box-shadow: 0 20px 70px rgba(0, 0, 0, 0.45); }
    .toast.show { display: block; }
    @media (max-width: 1180px) {
      .app { grid-template-columns: 1fr; }
      .rail { position: relative; height: auto; flex-direction: row; align-items: center; overflow-x: auto; }
      .nav { grid-auto-flow: column; grid-auto-columns: max-content; }
      .rail-card { display: none; }
      .topbar, .terminal-grid, .tenant-layout { grid-template-columns: 1fr; }
    }
    @media (max-width: 760px) {
      .workspace { padding: 10px; }
      .top-actions, .status-grid, .context-bar, .field-row, .file-controls, .login-fields, .file-browser { grid-template-columns: 1fr; }
      .login-field { grid-template-columns: 1fr; }
      .terminal-body { height: 460px; min-height: 460px; }
      .file-browser { height: 620px; }
      .file-list { border-right: 0; border-bottom: 1px solid var(--line); }
      .cli-card, .file-row { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <div class="app">
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
          <span class="nav-icon">&gt;_</span><span data-i18n="navTerminal">终端</span><span>01</span>
        </button>
        <button class="nav-button" type="button" data-view="tenants-view">
          <span class="nav-icon">TN</span><span data-i18n="navTenants">租户</span><span>02</span>
        </button>
      </nav>
      <div class="rail-card">
        <div class="rail-card-title"><span data-i18n="session">会话</span><span id="role-label">-</span></div>
        <div class="badge"><span class="led" id="session-led"></span><span id="session-label">not logged in</span></div>
      </div>
    </aside>

    <div class="workspace">
      <header class="topbar">
        <div class="page-title">
          <h1 id="page-heading">终端</h1>
          <p id="page-subtitle">登录认证、安装 CLI、检查租户隔离都从这里开始。</p>
        </div>
        <div class="top-actions">
          <div class="login-fields">
            <div class="login-field">
              <span data-i18n="username">用户名</span>
              <input id="username" autocomplete="username" placeholder="admin">
            </div>
            <div class="login-field">
              <span data-i18n="password">密码</span>
              <input id="password" type="password" autocomplete="current-password" placeholder="admin">
            </div>
          </div>
          <button class="primary" id="login" type="button" data-i18n="login">登录</button>
          <button class="ghost" id="logout" type="button" data-i18n="logout">退出</button>
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
            <div class="metric-label" data-i18n="availableCli">可用 CLI</div>
            <div class="metric-value"><span id="available-cli">0 / 0</span></div>
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
                    <p id="terminal-context">- / - / -</p>
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
                <div class="xterm-frame"><div id="terminal-container"></div></div>
              </div>
            </section>

            <aside class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">CL</span>
                  <div>
                    <h2>CLI Manager</h2>
                    <p data-i18n="cliManagerDesc">安装状态来自真实 PATH 探测。</p>
                  </div>
                </div>
                <button class="ghost" id="refresh-tools" type="button" title="Refresh">↻</button>
              </div>
              <div class="manager-body">
                <div id="installed-panel"></div>
              </div>
            </aside>
          </div>
        </section>

        <section class="view" id="tenants-view">
          <div class="tenant-layout">
            <section class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">TN</span>
                  <div>
                    <h2 data-i18n="tenantManager">用户与租户</h2>
                    <p data-i18n="tenantManagerDesc">admin 管理用户和租户边界，普通用户只看到自己的租户。</p>
                  </div>
                </div>
              </div>
              <div class="form-stack hidden" id="admin-user-form">
                <div class="field-row">
                  <div><label for="username-new" data-i18n="username">用户名</label><input id="username-new" placeholder="team-b"></div>
                  <div><label for="password-new" data-i18n="password">密码</label><input id="password-new" type="password" autocomplete="new-password" placeholder="••••••••"></div>
                </div>
                <div class="field-row">
                  <div><label for="subject-new" data-i18n="subject">主体</label><input id="subject-new" placeholder="tenant-user:team-b"></div>
                  <div><label for="tenant-new" data-i18n="tenant">租户</label><input id="tenant-new" placeholder="team-b"></div>
                </div>
                <div class="field-row">
                  <div><label for="role-new" data-i18n="role">角色</label><select id="role-new"><option value="tenant">tenant</option><option value="admin">admin</option></select></div>
                  <div><label for="duration-new" data-i18n="maxDuration">最长任务秒数</label><input id="duration-new" type="number" value="900"></div>
                </div>
                <div><label for="tools-new" data-i18n="allowedTools">允许工具</label><input id="tools-new" value="codex,claude,gemini,opencode,iflow,kimi,qoder"></div>
                <div><label for="workspaces-new" data-i18n="allowedWorkspaces">允许工作区</label><input id="workspaces-new" value="repo-*"></div>
                <div><label for="profiles-new" data-i18n="allowedProfiles">允许凭据配置</label><input id="profiles-new" value="team-default"></div>
                <div><label for="terminal-new" data-i18n="terminalAccess">终端权限</label><select id="terminal-new"><option value="true">allow</option><option value="false">block</option></select></div>
                <button class="primary" id="save-user" type="button" data-i18n="saveUser">保存用户</button>
              </div>
              <div class="tenant-list" id="tenant-list"></div>
              <div class="tenant-list hidden" id="user-list"></div>
            </section>

            <section class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">FS</span>
                  <div>
                    <h2 data-i18n="fileExplorer">文件浏览器</h2>
                    <p data-i18n="fileExplorerDesc">admin 可切换全部租户，普通租户只能访问自己的目录。</p>
                  </div>
                </div>
              </div>
              <div class="file-body">
                <div class="file-controls">
                  <div><label for="file-tenant" data-i18n="tenant">租户</label><select id="file-tenant"></select></div>
                  <div><label for="file-space" data-i18n="folder">目录</label><select id="file-space"><option value="workspaces">workspaces</option><option value="homes">homes</option></select></div>
                  <div><label for="file-path-input" data-i18n="path">路径</label><input id="file-path-input" value="/"></div>
                  <button class="ghost" id="file-parent" type="button" data-i18n="up">上级</button>
                  <button class="primary" id="file-refresh" type="button" data-i18n="open">打开</button>
                </div>
                <div class="file-path" id="file-abs-path">-</div>
                <div class="file-browser">
                  <div class="file-list" id="file-list"></div>
                  <div class="file-viewer">
                    <div class="file-viewbar" id="file-view-path" data-i18n="noFileSelected">未选择文件</div>
                    <pre class="file-preview" id="file-preview"></pre>
                  </div>
                </div>
              </div>
            </section>
          </div>
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
        navTenants: '租户',
        session: '会话',
        username: '用户名',
        password: '密码',
        login: '登录',
        logout: '退出',
        refresh: '刷新',
        health: '健康状态',
        ready: '就绪状态',
        availableCli: '可用 CLI',
        tenants: '租户',
        terminal: '终端',
        terminalSubtitle: '登录认证、安装 CLI、检查租户隔离都从这里开始。',
        tenantsSubtitle: '管理用户、查看租户边界和隔离文件系统。',
        connect: '连接',
        clear: '清屏',
        disconnect: '断开',
        tenant: '租户',
        workspace: '工作区',
        credentialProfile: '凭据配置',
        cliManagerDesc: '安装状态来自真实 PATH 探测。',
        installCli: '安装 CLI',
        quickLogin: 'Quick Login',
        verify: '验证',
        officialSource: '官方来源',
        healthOK: 'Health OK',
        notInstalled: '未安装',
        loginToCheck: '登录后检查',
        checking: '检查中',
        registeredOnly: '仅注册',
        delete: '删除',
        tenantManager: '用户与租户',
        tenantManagerDesc: 'admin 管理用户和租户边界，普通用户只看到自己的租户。',
        subject: '主体',
        role: '角色',
        allowedTools: '允许工具',
        allowedWorkspaces: '允许工作区',
        allowedProfiles: '允许凭据配置',
        terminalAccess: '终端权限',
        maxDuration: '最长任务秒数',
        saveUser: '保存用户',
        terminalAllowed: '允许终端',
        terminalBlocked: '禁止终端',
        subjects: '主体',
        tools: '工具',
        workspaces: '工作区',
        profiles: '凭据配置',
        tokenCount: '凭据数',
        dataFolders: '数据目录',
        browseFiles: '打开文件',
        fileExplorer: '文件浏览器',
        fileExplorerDesc: 'admin 可切换全部租户，普通租户只能访问自己的目录。',
        folder: '目录',
        path: '路径',
        up: '上级',
        open: '打开',
        noFiles: '目录为空',
        noFileSelected: '未选择文件',
        noTenants: '暂无可访问租户',
        noUsers: '暂无用户',
        userSaved: '用户已保存',
        userDeleted: '用户已删除',
        connected: '已连接',
        disconnected: '未连接',
        connecting: '连接中',
        connectionError: '连接错误',
        exited: '已退出',
        terminalWelcome: 'Agent Runtime 终端已就绪。先登录用户，再连接 shell 或执行 CLI 安装/登录。',
        terminalConnecting: '正在连接终端...',
        loginOK: '登录成功',
        refreshed: '状态已刷新',
        loggedOut: '已退出登录'
      },
      en: {
        brandSubtitle: 'CLI control plane',
        navTerminal: 'Terminal',
        navTenants: 'Tenants',
        session: 'Session',
        username: 'Username',
        password: 'Password',
        login: 'Login',
        logout: 'Logout',
        refresh: 'Refresh',
        health: 'Health',
        ready: 'Ready',
        availableCli: 'Available CLI',
        tenants: 'Tenants',
        terminal: 'Terminal',
        terminalSubtitle: 'Login, CLI installation, and tenant isolation checks start here.',
        tenantsSubtitle: 'Manage users, inspect tenant boundaries, and browse isolated filesystems.',
        connect: 'Connect',
        clear: 'Clear',
        disconnect: 'Disconnect',
        tenant: 'Tenant',
        workspace: 'Workspace',
        credentialProfile: 'Credential Profile',
        cliManagerDesc: 'Install state is checked against the real PATH.',
        installCli: 'Install CLI',
        quickLogin: 'Quick Login',
        verify: 'Verify',
        officialSource: 'Official Source',
        healthOK: 'Health OK',
        notInstalled: 'Not installed',
        loginToCheck: 'Login to check',
        checking: 'Checking',
        registeredOnly: 'Registered only',
        delete: 'Delete',
        tenantManager: 'Users & Tenants',
        tenantManagerDesc: 'Admins manage users and tenant boundaries. Users only see their own tenant.',
        subject: 'Subject',
        role: 'Role',
        allowedTools: 'Allowed Tools',
        allowedWorkspaces: 'Allowed Workspaces',
        allowedProfiles: 'Allowed Profiles',
        terminalAccess: 'Terminal Access',
        maxDuration: 'Max Job Seconds',
        saveUser: 'Save User',
        terminalAllowed: 'Terminal allowed',
        terminalBlocked: 'Terminal blocked',
        subjects: 'Subjects',
        tools: 'Tools',
        workspaces: 'Workspaces',
        profiles: 'Profiles',
        tokenCount: 'Credential Count',
        dataFolders: 'Data Folders',
        browseFiles: 'Open Files',
        fileExplorer: 'File Explorer',
        fileExplorerDesc: 'Admins can switch tenants. Tenant users can only access their own folders.',
        folder: 'Folder',
        path: 'Path',
        up: 'Up',
        open: 'Open',
        noFiles: 'Folder is empty',
        noFileSelected: 'No file selected',
        noTenants: 'No accessible tenants',
        noUsers: 'No users',
        userSaved: 'User saved',
        userDeleted: 'User deleted',
        connected: 'Connected',
        disconnected: 'Disconnected',
        connecting: 'Connecting',
        connectionError: 'Connection error',
        exited: 'Exited',
        terminalWelcome: 'Agent Runtime terminal is ready. Log in as a user, then connect a shell or run CLI install/login.',
        terminalConnecting: 'Connecting terminal...',
        loginOK: 'Login succeeded',
        refreshed: 'Runtime status refreshed',
        loggedOut: 'Logged out'
      }
    };

    const savedLanguage = localStorage.getItem('agent-runtime-lang');
    const state = {
      lang: savedLanguage || 'zh',
      view: 'terminal-view',
      sessionToken: localStorage.getItem('agent-runtime-session-token') || localStorage.getItem('agent-runtime-token') || '',
      session: null,
      tools: [],
      toolsLoaded: false,
      toolsContextReady: false,
      tenants: [],
      users: [],
      ws: null,
      connected: false,
      statusKey: 'disconnected',
      term: null,
      pendingCommand: '',
      installPollTimer: null,
      userDefaults: { tenant: '', subject: '' }
    };
    const usernameInput = document.getElementById('username');
    const passwordInput = document.getElementById('password');
    usernameInput.value = localStorage.getItem('agent-runtime-username') || 'admin';

    const installSources = [
      { label: 'Claude Code', fallback: 'CC', logo: 'https://claude.ai/favicon.svg', tool: 'claude', command: 'curl -fsSL https://claude.ai/install.sh | bash', verify: 'claude --version', login: 'claude', docs: 'https://docs.anthropic.com/en/docs/claude-code/quickstart', provider: 'Anthropic' },
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
    function join(values) { return values && values.length ? values.join(', ') : '-'; }
    function splitCSV(value) { return String(value || '').split(',').map(function(item) { return item.trim(); }).filter(Boolean); }
    function slugFromUsername(value) {
      const slug = String(value || '').toLowerCase().trim()
        .replace(/[^a-z0-9._-]+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^[-._]+|[-._]+$/g, '');
      return slug || 'tenant';
    }
    function syncUserDefaults(force) {
      const username = $('username-new').value.trim();
      const tenant = slugFromUsername(username);
      const subject = 'tenant-user:' + tenant;
      if (force || !$('tenant-new').value.trim() || $('tenant-new').value === state.userDefaults.tenant) {
        $('tenant-new').value = username ? tenant : '';
      }
      if (force || !$('subject-new').value.trim() || $('subject-new').value === state.userDefaults.subject) {
        $('subject-new').value = username ? subject : '';
      }
      state.userDefaults = {
        tenant: username ? tenant : '',
        subject: username ? subject : ''
      };
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
      updatePageTitle();
      renderSession();
      renderTerminalOptions();
      renderCliManager();
      renderTenants();
      renderUsers();
      setConnected(state.connected, state.statusKey);
      if (state.term && state.term.buffer.active.length <= 2) {
        state.term.clear();
        state.term.writeln(t('terminalWelcome'));
      }
    }

    function updatePageTitle() {
      const titles = {
        'terminal-view': [t('terminal'), t('terminalSubtitle')],
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
          red: '#ff5c73',
          green: '#37e681',
          yellow: '#f5bc4f',
          blue: '#4c8dff',
          magenta: '#b68cff',
          cyan: '#28e0d4'
        }
      });
      state.term.open($('terminal-container'));
      state.term.writeln(t('terminalWelcome'));
      state.term.onData(function(data) {
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
      const rect = $('terminal-container').getBoundingClientRect();
      const cols = Math.max(40, Math.floor(rect.width / 9));
      const rows = Math.max(12, Math.floor(rect.height / 18));
      state.term.resize(cols, rows);
      if (state.ws && state.ws.readyState === WebSocket.OPEN) {
        state.ws.send(JSON.stringify({ type: 'resize', cols: cols, rows: rows }));
      }
      return { cols: cols, rows: rows };
    }

    async function login(showMessage) {
      const username = usernameInput.value.trim();
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: username, password: passwordInput.value })
      });
      if (!response.ok) {
        let message = response.statusText;
        try {
          const body = await response.json();
          message = body.error || message;
        } catch (err) {}
        throw new Error(message);
      }
      const body = await response.json();
      state.sessionToken = body.token || '';
      state.session = body.session || null;
      localStorage.setItem('agent-runtime-session-token', state.sessionToken);
      localStorage.setItem('agent-runtime-username', username);
      passwordInput.value = '';
      renderSession();
      await refreshSecure();
      if (showMessage) showToast(t('loginOK') + ': ' + state.session.subject);
    }

    async function restoreSession() {
      if (!state.sessionToken) throw new Error('not logged in');
      state.session = await api('/api/session');
      localStorage.setItem('agent-runtime-session-token', state.sessionToken);
      localStorage.removeItem('agent-runtime-token');
      renderSession();
      await refreshSecure();
    }

    function logout() {
      state.sessionToken = '';
      state.session = null;
      state.tenants = [];
      state.users = [];
      state.tools = [];
      state.toolsLoaded = false;
      state.toolsContextReady = false;
      localStorage.removeItem('agent-runtime-session-token');
      localStorage.removeItem('agent-runtime-token');
      disconnectTerminal();
      renderSession();
      renderTerminalOptions();
      renderCliManager();
      renderTenants();
      renderUsers();
      showToast(t('loggedOut'));
    }

    async function refresh() {
      const results = await Promise.allSettled([api('/api/health'), api('/api/ready'), api('/api/status')]);
      setMetric('health', 'health-led', results[0].status === 'fulfilled' ? results[0].value.status : 'error');
      setMetric('ready', 'ready-led', results[1].status === 'fulfilled' ? results[1].value.status : 'error');
      try {
        await restoreSession();
      } catch (err) {
        state.session = null;
        state.tenants = [];
        state.users = [];
        state.toolsLoaded = false;
        state.toolsContextReady = false;
        renderSession();
        renderCliManager();
        renderTenants();
        renderUsers();
      }
    }

    async function refreshSecure() {
      await refreshTenants();
      renderTerminalOptions();
      await refreshTools();
      if (state.session && state.session.admin) {
        await refreshUsers();
      } else {
        state.users = [];
      }
      renderSession();
      renderCliManager();
      renderTenants();
      renderUsers();
      await refreshFiles().catch(function() {});
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
      $('tenant-count').textContent = String(state.tenants.length);
    }

    async function refreshUsers() {
      const body = await api('/api/users');
      state.users = body.users || [];
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
      $('session-label').textContent = loggedIn ? state.session.subject : 'not logged in';
      $('role-label').textContent = loggedIn ? state.session.role : '-';
      $('admin-user-form').classList.toggle('hidden', !(state.session && state.session.admin));
      $('user-list').classList.toggle('hidden', !(state.session && state.session.admin));
      $('login').disabled = loggedIn;
      $('logout').disabled = !loggedIn;
    }

    function renderTerminalOptions() {
      const tenantSelect = $('tenant');
      const fileTenantSelect = $('file-tenant');
      const currentTenant = tenantSelect.value;
      const currentFileTenant = fileTenantSelect.value;
      const tenantOptions = state.tenants.map(function(tenant) {
        return '<option value="' + escapeHTML(tenant.id) + '">' + escapeHTML(tenant.id) + '</option>';
      }).join('');
      tenantSelect.innerHTML = tenantOptions;
      fileTenantSelect.innerHTML = tenantOptions;
      if (currentTenant) tenantSelect.value = currentTenant;
      if (!tenantSelect.value && state.tenants[0]) tenantSelect.value = state.tenants[0].id;
      if (currentFileTenant) fileTenantSelect.value = currentFileTenant;
      if (!fileTenantSelect.value && tenantSelect.value) fileTenantSelect.value = tenantSelect.value;
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
      $('terminal-context').textContent = ($('tenant').value || '-') + ' / ' + ($('workspace').value || '-') + ' / ' + ($('profile').value || '-');
    }

    function terminalURL() {
      const size = resizeTerminal();
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const params = new URLSearchParams({
        token: state.sessionToken,
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

    function renderCliManager() {
      renderInstalledClis();
    }

    function renderCliLogo(source) {
      const fallback = source.fallback || source.label.slice(0, 2).toUpperCase();
      return '<div class="cli-logo"><img src="' + escapeHTML(source.logo || '') + '" alt="' + escapeHTML(source.label) + '" loading="lazy" referrerpolicy="no-referrer" onerror="this.style.display=&quot;none&quot;;this.nextElementSibling.style.display=&quot;grid&quot;"><span>' + escapeHTML(fallback) + '</span></div>';
    }

    function renderInstalledClis() {
      const known = new Map(state.tools.map(function(tool) { return [tool.name, tool]; }));
      $('installed-panel').innerHTML = installSources.map(function(source) {
        const tool = known.get(source.tool) || { name: source.tool, available: false, health: 'missing' };
        const canCheck = state.toolsLoaded && state.toolsContextReady;
        const available = canCheck && !!tool.available;
        const version = available ? (tool.detected_version || tool.version || '-') : (canCheck ? t('notInstalled') : t('loginToCheck'));
        const healthClass = available ? 'ok' : (canCheck ? 'bad' : 'warn');
        const healthText = available ? t('healthOK') : (canCheck ? t('notInstalled') : t('loginToCheck'));
        const actions = !canCheck
          ? '<button class="ghost" type="button" disabled>' + escapeHTML(t('checking')) + '</button>'
          : available
          ? '<button class="ghost" type="button" data-login-command="' + escapeHTML(source.login) + '">' + escapeHTML(t('quickLogin')) + ' ↪</button>' +
            '<button class="ghost" type="button" data-install-command="' + escapeHTML(source.verify) + '">' + escapeHTML(t('verify')) + '</button>' +
            (state.session && state.session.admin && known.has(source.tool) ? '<button class="danger" type="button" data-delete-tool="' + escapeHTML(source.tool) + '">' + escapeHTML(t('delete')) + '</button>' : '')
          : '<button class="primary" type="button" data-install-command="' + escapeHTML(source.command) + '">' + escapeHTML(t('installCli')) + '</button>' +
            '<a href="' + escapeHTML(source.docs) + '" target="_blank" rel="noopener noreferrer">' + escapeHTML(t('officialSource')) + '</a>';
        return '<article class="cli-card">' +
          renderCliLogo(source) +
          '<div>' +
            '<div class="cli-name">' + escapeHTML(source.label) + '</div>' +
            '<div class="cli-version">' + escapeHTML(version) + '</div>' +
            '<div class="cli-health"><span class="badge ' + healthClass + '"><span class="led ' + (available ? 'ok' : 'bad') + '"></span>' + escapeHTML(healthText) + '</span></div>' +
          '</div>' +
          '<div class="cli-actions">' + actions + '</div>' +
        '</article>';
      }).join('');
      bindManagerButtons($('installed-panel'));
    }

    function bindManagerButtons(root) {
      root.querySelectorAll('[data-install-command]').forEach(function(button) {
        button.addEventListener('click', function() {
          runCommand(button.dataset.installCommand);
          startInstallPolling();
        });
      });
      root.querySelectorAll('[data-login-command]').forEach(function(button) {
        button.addEventListener('click', function() { runCommand(button.dataset.loginCommand); });
      });
      root.querySelectorAll('[data-delete-tool]').forEach(function(button) {
        button.addEventListener('click', async function() {
          await api('/api/tools/' + encodeURIComponent(button.dataset.deleteTool), { method: 'DELETE' });
          await refreshTools();
          renderCliManager();
        });
      });
    }

    function renderTenants() {
      const container = $('tenant-list');
      if (!state.tenants.length) {
        container.innerHTML = '<div class="empty">' + escapeHTML(t('noTenants')) + '</div>';
        return;
      }
      container.innerHTML = state.tenants.map(function(tenant) {
        return '<article class="tenant-card">' +
          '<div class="tenant-card-head"><span>' + escapeHTML(tenant.id) + '</span><span class="badge ' + (tenant.allow_terminal ? 'ok' : 'warn') + '">' + escapeHTML(tenant.allow_terminal ? t('terminalAllowed') : t('terminalBlocked')) + '</span></div>' +
          '<div class="kv">' +
            '<div><strong>' + escapeHTML(t('subjects')) + '</strong><span>' + escapeHTML(join(tenant.subjects)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('tools')) + '</strong><span>' + escapeHTML(join(tenant.allowed_tools)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('workspaces')) + '</strong><span>' + escapeHTML(join(tenant.workspace_patterns)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('profiles')) + '</strong><span>' + escapeHTML(join(tenant.credential_profiles)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('tokenCount')) + '</strong><span>' + escapeHTML(tenant.token_count || 0) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('dataFolders')) + '</strong><span class="mono">tenants/' + escapeHTML(tenant.id) + '/workspaces<br>tenants/' + escapeHTML(tenant.id) + '/homes</span></div>' +
          '</div>' +
          '<div class="tenant-card-actions"><button class="ghost" type="button" data-browse-tenant="' + escapeHTML(tenant.id) + '">' + escapeHTML(t('browseFiles')) + '</button></div>' +
        '</article>';
      }).join('');
      container.querySelectorAll('[data-browse-tenant]').forEach(function(button) {
        button.addEventListener('click', function() {
          $('file-tenant').value = button.dataset.browseTenant;
          $('file-space').value = 'workspaces';
          $('file-path-input').value = '/';
          refreshFiles().catch(function(err) { showToast(err.message); });
        });
      });
    }

    function renderUsers() {
      if (!(state.session && state.session.admin)) return;
      const container = $('user-list');
      if (!state.users.length) {
        container.innerHTML = '<div class="empty">' + escapeHTML(t('noUsers')) + '</div>';
        return;
      }
      container.innerHTML = state.users.map(function(item) {
        return '<article class="user-card">' +
          '<div class="user-card-head"><span>' + escapeHTML(item.username) + '</span><span class="badge">' + escapeHTML(item.role) + '</span></div>' +
          '<div class="kv">' +
            '<div><strong>' + escapeHTML(t('subject')) + '</strong><span>' + escapeHTML(item.subject) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('tenant')) + '</strong><span>' + escapeHTML(item.tenant) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('tools')) + '</strong><span>' + escapeHTML(join(item.allowed_tools)) + '</span></div>' +
            '<div><strong>' + escapeHTML(t('profiles')) + '</strong><span>' + escapeHTML(join(item.allowed_credential_profiles)) + '</span></div>' +
          '</div>' +
          '<div class="user-card-actions"><button class="ghost" type="button" data-browse-tenant="' + escapeHTML(item.tenant) + '">' + escapeHTML(t('browseFiles')) + '</button><button class="danger" type="button" data-delete-user="' + escapeHTML(item.id) + '">' + escapeHTML(t('delete')) + '</button></div>' +
        '</article>';
      }).join('');
      container.querySelectorAll('[data-browse-tenant]').forEach(function(button) {
        button.addEventListener('click', function() {
          $('file-tenant').value = button.dataset.browseTenant;
          $('file-space').value = 'workspaces';
          $('file-path-input').value = '/';
          refreshFiles().catch(function(err) { showToast(err.message); });
        });
      });
      container.querySelectorAll('[data-delete-user]').forEach(function(button) {
        button.addEventListener('click', async function() {
          await api('/api/users/' + encodeURIComponent(button.dataset.deleteUser), { method: 'DELETE' });
          showToast(t('userDeleted'));
          await refreshSecure();
        });
      });
    }

    async function saveUser() {
      syncUserDefaults(false);
      const payload = {
        username: $('username-new').value.trim(),
        password: $('password-new').value,
        subject: $('subject-new').value.trim(),
        tenant: $('tenant-new').value.trim(),
        role: $('role-new').value,
        allowed_tools: splitCSV($('tools-new').value),
        allowed_workspaces: splitCSV($('workspaces-new').value),
        allowed_credential_profiles: splitCSV($('profiles-new').value),
        allow_terminal: $('terminal-new').value === 'true',
        max_job_seconds: Number($('duration-new').value || 0)
      };
      const user = await api('/api/users', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
      $('password-new').value = '';
      showToast(t('userSaved'));
      await refreshSecure();
      if (user && user.tenant) {
        $('file-tenant').value = user.tenant;
        $('file-space').value = 'workspaces';
        $('file-path-input').value = '/';
        await refreshFiles().catch(function(err) { showToast(err.message); });
      }
    }

    async function refreshFiles() {
      if (!$('file-tenant').value) {
        $('file-list').innerHTML = '<div class="empty">' + escapeHTML(t('noTenants')) + '</div>';
        $('file-view-path').textContent = t('noFileSelected');
        $('file-preview').textContent = '';
        return;
      }
      const params = new URLSearchParams({
        tenant: $('file-tenant').value,
        space: $('file-space').value,
        path: $('file-path-input').value.trim()
      });
      const body = await api('/api/files?' + params.toString());
      $('file-abs-path').textContent = body.abs_path || '-';
      $('file-path-input').value = body.path && body.path !== '.' ? '/' + body.path : '/';
      $('file-view-path').textContent = t('noFileSelected');
      $('file-preview').textContent = '';
      if (!body.entries || !body.entries.length) {
        $('file-list').innerHTML = '<div class="empty">' + escapeHTML(t('noFiles')) + '</div>';
        return;
      }
      $('file-list').innerHTML = body.entries.map(function(item) {
        const isDir = item.kind === 'directory';
        const icon = isDir ? '▸' : '·';
        const size = isDir ? item.kind : formatBytes(item.size || 0);
        return '<div class="file-row" data-entry-kind="' + escapeHTML(item.kind) + '" data-entry-path="/' + escapeHTML(item.path) + '">' +
          '<div class="file-name"><span>' + escapeHTML(icon) + '</span>' +
          '<button type="button">' + escapeHTML(item.name + (isDir ? '/' : '')) + '</button>' +
          '</div>' +
          '<div class="file-meta">' + escapeHTML(size) + '</div>' +
        '</div>';
      }).join('');
      $('file-list').querySelectorAll('[data-entry-path]').forEach(function(row) {
        row.addEventListener('click', async function() {
          $('file-list').querySelectorAll('.file-row').forEach(function(item) { item.classList.remove('active'); });
          row.classList.add('active');
          if (row.dataset.entryKind === 'directory') {
            $('file-path-input').value = row.dataset.entryPath;
            await refreshFiles().catch(function(err) { showToast(err.message); });
            return;
          }
          await previewFile(row.dataset.entryPath).catch(function(err) { showToast(err.message); });
        });
      });
    }

    async function previewFile(path) {
      const params = new URLSearchParams({
        tenant: $('file-tenant').value,
        space: $('file-space').value,
        path: path
      });
      $('file-view-path').textContent = path;
      $('file-preview').textContent = 'Loading...';
      const body = await api('/api/files/raw?' + params.toString());
      $('file-view-path').textContent = (body.path && body.path !== '.' ? '/' + body.path : path) + ' · ' + formatBytes(body.size || 0);
      $('file-preview').textContent = body.content + (body.truncated ? '\n\n[truncated]' : '');
    }

    function formatBytes(value) {
      value = Number(value || 0);
      if (value < 1024) return value + ' B';
      if (value < 1024 * 1024) return (value / 1024).toFixed(1) + ' KB';
      return (value / 1024 / 1024).toFixed(1) + ' MB';
    }

    function parentPath(path) {
      path = String(path || '/').replace(/\/+$/, '');
      if (!path || path === '/') return '/';
      const index = path.lastIndexOf('/');
      return index <= 0 ? '/' : path.slice(0, index);
    }

    function switchView(viewID) {
      state.view = viewID;
      document.querySelectorAll('.nav-button').forEach(function(item) { item.classList.toggle('active', item.dataset.view === viewID); });
      document.querySelectorAll('.view').forEach(function(item) { item.classList.toggle('active', item.id === viewID); });
      updatePageTitle();
      if (viewID === 'terminal-view') window.setTimeout(resizeTerminal, 0);
      if (viewID === 'tenants-view') refreshFiles().catch(function(err) { showToast(err.message); });
    }

    document.querySelectorAll('[data-view]').forEach(function(button) { button.addEventListener('click', function() { switchView(button.dataset.view); }); });
    document.querySelectorAll('[data-lang]').forEach(function(button) {
      button.addEventListener('click', function() {
        state.lang = button.dataset.lang;
        localStorage.setItem('agent-runtime-lang', state.lang);
        applyLanguage();
      });
    });
    $('login').addEventListener('click', function() { login(true).catch(function(err) { showToast(err.message); }); });
    $('logout').addEventListener('click', logout);
    $('password').addEventListener('keydown', function(event) {
      if (event.key === 'Enter') login(true).catch(function(err) { showToast(err.message); });
    });
    $('refresh').addEventListener('click', function() { refresh().then(function() { showToast(t('refreshed')); }).catch(function(err) { showToast(err.message); }); });
    $('refresh-tools').addEventListener('click', function() { refreshTools().then(renderCliManager).catch(function(err) { showToast(err.message); }); });
    $('tenant').addEventListener('change', function() { updateProfileOptions(); refreshTools().then(renderCliManager).catch(function() {}); });
    $('workspace').addEventListener('input', updateContextLabels);
    $('profile').addEventListener('input', function() { updateContextLabels(); refreshTools().then(renderCliManager).catch(function() {}); });
    $('connect-terminal').addEventListener('click', connectTerminal);
    $('disconnect-terminal').addEventListener('click', disconnectTerminal);
    $('clear-terminal').addEventListener('click', function() { if (state.term) state.term.clear(); });
    $('username-new').addEventListener('input', function() { syncUserDefaults(false); });
    $('save-user').addEventListener('click', function() { saveUser().catch(function(err) { showToast(err.message); }); });
    $('file-refresh').addEventListener('click', function() { refreshFiles().catch(function(err) { showToast(err.message); }); });
    $('file-parent').addEventListener('click', function() { $('file-path-input').value = parentPath($('file-path-input').value); refreshFiles().catch(function(err) { showToast(err.message); }); });
    $('file-tenant').addEventListener('change', function() { refreshFiles().catch(function(err) { showToast(err.message); }); });
    $('file-space').addEventListener('change', function() { $('file-path-input').value = '/'; refreshFiles().catch(function(err) { showToast(err.message); }); });

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
