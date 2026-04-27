package api

const webUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Agent Runtime</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #05070a;
      --bg-2: #081019;
      --panel: rgba(12, 18, 27, 0.88);
      --panel-solid: #0b111a;
      --panel-soft: #0f1722;
      --line: #1e2a36;
      --line-strong: #2a3b4d;
      --text: #e8f2ff;
      --muted: #8ca0b7;
      --faint: #536579;
      --cyan: #28e0d4;
      --blue: #4c8dff;
      --green: #3ee486;
      --amber: #f5b84b;
      --red: #ff5c73;
      --terminal: #020407;
      --mono: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    html, body { min-height: 100%; }
    body {
      margin: 0;
      background:
        linear-gradient(rgba(40, 224, 212, 0.035) 1px, transparent 1px),
        linear-gradient(90deg, rgba(76, 141, 255, 0.03) 1px, transparent 1px),
        radial-gradient(circle at 30% -10%, rgba(40, 224, 212, 0.14), transparent 34%),
        linear-gradient(135deg, var(--bg), var(--bg-2) 62%, #030507);
      background-size: 36px 36px, 36px 36px, auto, auto;
      color: var(--text);
      line-height: 1.45;
    }
    button, input, select {
      font: inherit;
    }
    button {
      min-height: 36px;
      border: 1px solid var(--line-strong);
      border-radius: 8px;
      background: #111b28;
      color: var(--text);
      cursor: pointer;
      font-weight: 700;
      letter-spacing: 0;
    }
    button:hover { border-color: rgba(40, 224, 212, 0.65); }
    button:disabled { opacity: 0.46; cursor: not-allowed; }
    button.primary {
      border-color: rgba(40, 224, 212, 0.78);
      background: linear-gradient(135deg, rgba(40, 224, 212, 0.96), rgba(76, 141, 255, 0.9));
      color: #021014;
      box-shadow: 0 0 22px rgba(40, 224, 212, 0.18);
    }
    button.ghost {
      background: rgba(13, 20, 31, 0.68);
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
      background: rgba(3, 7, 12, 0.68);
      color: var(--text);
      outline: none;
    }
    input:focus, select:focus {
      border-color: rgba(40, 224, 212, 0.8);
      box-shadow: 0 0 0 3px rgba(40, 224, 212, 0.11);
    }
    label {
      display: block;
      margin-bottom: 6px;
      color: var(--muted);
      font-size: 11px;
      font-weight: 800;
      text-transform: uppercase;
    }
    code, pre, .mono { font-family: var(--mono); }
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
      background: rgba(4, 8, 13, 0.82);
      backdrop-filter: blur(18px);
      display: flex;
      flex-direction: column;
      gap: 18px;
    }
    .brand {
      display: grid;
      grid-template-columns: 38px 1fr;
      gap: 10px;
      align-items: center;
      min-height: 48px;
    }
    .brand-mark {
      width: 38px;
      height: 38px;
      border: 1px solid rgba(40, 224, 212, 0.56);
      border-radius: 8px;
      background:
        linear-gradient(135deg, rgba(40, 224, 212, 0.18), rgba(76, 141, 255, 0.16)),
        #071018;
      display: grid;
      place-items: center;
      color: var(--cyan);
      font-family: var(--mono);
      font-weight: 900;
      box-shadow: inset 0 0 24px rgba(40, 224, 212, 0.1);
    }
    .brand-title { font-size: 17px; font-weight: 800; }
    .brand-subtitle { color: var(--faint); font-size: 12px; }
    .nav {
      display: grid;
      gap: 8px;
    }
    .nav-button {
      display: grid;
      grid-template-columns: 22px 1fr auto;
      gap: 10px;
      align-items: center;
      width: 100%;
      padding: 9px 10px;
      text-align: left;
      background: transparent;
      color: var(--muted);
    }
    .nav-button.active {
      border-color: rgba(40, 224, 212, 0.58);
      background: rgba(40, 224, 212, 0.1);
      color: var(--text);
    }
    .nav-icon {
      width: 22px;
      height: 22px;
      border-radius: 6px;
      display: grid;
      place-items: center;
      background: rgba(255, 255, 255, 0.04);
      color: var(--cyan);
      font-family: var(--mono);
      font-size: 12px;
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
      font-size: 12px;
      color: var(--muted);
      text-transform: uppercase;
      font-weight: 800;
    }
    .workspace {
      min-width: 0;
      padding: 16px;
    }
    .topbar {
      min-height: 64px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(8, 13, 20, 0.84);
      backdrop-filter: blur(18px);
      display: grid;
      grid-template-columns: minmax(260px, 1fr) minmax(380px, 560px);
      gap: 16px;
      align-items: center;
      padding: 12px 14px;
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
      grid-template-columns: 1fr 92px;
      gap: 8px;
      align-items: center;
    }
    .token-field {
      display: grid;
      grid-template-columns: 86px minmax(0, 1fr);
      gap: 8px;
      align-items: center;
    }
    .token-field span {
      color: var(--muted);
      font-size: 12px;
      font-weight: 800;
      text-transform: uppercase;
    }
    .content {
      margin-top: 16px;
    }
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
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(10, 17, 26, 0.74);
      padding: 11px 12px;
    }
    .metric-label {
      color: var(--muted);
      font-size: 11px;
      font-weight: 800;
      text-transform: uppercase;
    }
    .metric-value {
      margin-top: 8px;
      display: flex;
      align-items: center;
      gap: 8px;
      color: var(--text);
      font-size: 15px;
      font-weight: 800;
      overflow-wrap: anywhere;
    }
    .led {
      width: 8px;
      height: 8px;
      border-radius: 999px;
      background: var(--amber);
      box-shadow: 0 0 14px rgba(245, 184, 75, 0.65);
      flex: 0 0 auto;
    }
    .led.ok {
      background: var(--green);
      box-shadow: 0 0 14px rgba(62, 228, 134, 0.72);
    }
    .led.bad {
      background: var(--red);
      box-shadow: 0 0 14px rgba(255, 92, 115, 0.72);
    }
    .terminal-grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 340px;
      gap: 16px;
      align-items: start;
    }
    .panel {
      border: 1px solid var(--line);
      border-radius: 8px;
      background: var(--panel);
      box-shadow: 0 24px 90px rgba(0, 0, 0, 0.28);
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
      font-weight: 800;
      white-space: nowrap;
    }
    .badge.ok { color: #b7ffd4; border-color: rgba(62, 228, 134, 0.38); }
    .badge.warn { color: #ffe0a3; border-color: rgba(245, 184, 75, 0.42); }
    .context-bar {
      display: grid;
      grid-template-columns: 150px minmax(160px, 1fr) minmax(160px, 1fr) 116px;
      gap: 10px;
      padding: 14px;
      border-bottom: 1px solid var(--line);
      background: rgba(5, 9, 14, 0.35);
      align-items: end;
    }
    .terminal-shell {
      padding: 12px;
    }
    .terminal-screen {
      position: relative;
      height: clamp(360px, calc(100vh - 500px), 560px);
      min-height: 360px;
      border: 1px solid #172332;
      border-radius: 8px;
      background:
        linear-gradient(rgba(40, 224, 212, 0.028) 1px, transparent 1px),
        linear-gradient(90deg, rgba(40, 224, 212, 0.018) 1px, transparent 1px),
        var(--terminal);
      background-size: 22px 22px, 22px 22px, auto;
      overflow: auto;
      outline: none;
      box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.02), inset 0 0 42px rgba(40, 224, 212, 0.04);
    }
    .terminal-screen:focus {
      border-color: rgba(40, 224, 212, 0.72);
      box-shadow: 0 0 0 3px rgba(40, 224, 212, 0.09), inset 0 0 42px rgba(40, 224, 212, 0.04);
    }
    .terminal-output {
      margin: 0;
      min-height: 100%;
      padding: 15px;
      color: #d7fbe9;
      font-size: 13px;
      line-height: 1.5;
      white-space: pre-wrap;
      overflow-wrap: anywhere;
    }
    .screen-hint {
      position: sticky;
      bottom: 0;
      display: flex;
      justify-content: space-between;
      gap: 10px;
      padding: 8px 12px;
      border-top: 1px solid rgba(30, 42, 54, 0.68);
      background: rgba(3, 6, 10, 0.82);
      color: var(--faint);
      font-size: 12px;
    }
    .command-dock {
      display: grid;
      grid-template-columns: auto minmax(0, 1fr) 80px;
      gap: 8px;
      align-items: center;
      margin-top: 10px;
      padding: 9px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(6, 11, 17, 0.82);
    }
    .prompt-label {
      color: var(--cyan);
      font-family: var(--mono);
      font-size: 13px;
      font-weight: 800;
      white-space: nowrap;
    }
    .terminal-actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      padding: 0 12px 12px;
    }
    .side-stack {
      display: grid;
      gap: 16px;
    }
    .quick-grid {
      display: grid;
      gap: 8px;
      padding: 12px;
    }
    .quick-button {
      display: grid;
      grid-template-columns: 1fr auto;
      gap: 8px;
      align-items: center;
      width: 100%;
      min-height: 46px;
      padding: 10px 11px;
      text-align: left;
      background: rgba(9, 16, 24, 0.72);
    }
    .quick-button strong {
      display: block;
      font-size: 13px;
    }
    .quick-button span {
      display: block;
      margin-top: 2px;
      color: var(--faint);
      font-family: var(--mono);
      font-size: 11px;
      font-weight: 500;
    }
    .tool-cards {
      display: grid;
      gap: 10px;
      padding: 12px;
    }
    .tool-card {
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(5, 10, 16, 0.62);
      padding: 11px;
    }
    .tool-card-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
    }
    .tool-name {
      font-family: var(--mono);
      font-size: 13px;
      font-weight: 900;
      color: var(--text);
    }
    .tool-path {
      margin-top: 8px;
      color: var(--muted);
      font-family: var(--mono);
      font-size: 11px;
      overflow-wrap: anywhere;
    }
    .empty {
      padding: 18px;
      color: var(--muted);
      border: 1px dashed var(--line-strong);
      border-radius: 8px;
      background: rgba(255, 255, 255, 0.025);
    }
    .manager-grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 360px;
      gap: 16px;
      align-items: start;
    }
    .table-wrap {
      overflow: auto;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 13px;
    }
    th, td {
      padding: 12px 12px;
      border-bottom: 1px solid var(--line);
      text-align: left;
      vertical-align: top;
    }
    th {
      color: var(--muted);
      font-size: 11px;
      text-transform: uppercase;
      font-weight: 900;
    }
    td code {
      color: #d5e8ff;
      font-size: 12px;
      overflow-wrap: anywhere;
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
      border: 1px solid var(--line);
      border-radius: 8px;
      background: rgba(6, 12, 18, 0.7);
      padding: 14px;
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
      font-weight: 800;
      text-transform: uppercase;
      font-size: 11px;
    }
    .toast {
      position: fixed;
      right: 18px;
      bottom: 18px;
      z-index: 10;
      max-width: 420px;
      padding: 12px 14px;
      border: 1px solid rgba(40, 224, 212, 0.42);
      border-radius: 8px;
      background: rgba(5, 10, 16, 0.94);
      color: var(--text);
      box-shadow: 0 20px 70px rgba(0, 0, 0, 0.45);
      display: none;
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
      .terminal-grid, .manager-grid { grid-template-columns: 1fr; }
      .topbar { grid-template-columns: 1fr; }
    }
    @media (max-width: 760px) {
      .workspace { padding: 10px; }
      .status-grid, .context-bar, .field-row { grid-template-columns: 1fr; }
      .terminal-screen { min-height: 360px; }
      .command-dock { grid-template-columns: 1fr; }
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
          <div class="brand-subtitle">CLI control plane</div>
        </div>
      </div>
      <nav class="nav" aria-label="Main navigation">
        <button class="nav-button active" type="button" data-view="terminal-view">
          <span class="nav-icon">&gt;_</span><span>Terminal</span><span class="nav-kbd">01</span>
        </button>
        <button class="nav-button" type="button" data-view="tools-view">
          <span class="nav-icon">CL</span><span>CLI Manager</span><span class="nav-kbd">02</span>
        </button>
        <button class="nav-button" type="button" data-view="tenants-view">
          <span class="nav-icon">TN</span><span>Tenants</span><span class="nav-kbd">03</span>
        </button>
      </nav>
      <div class="rail-card">
        <div class="rail-card-title"><span>Runtime</span><span id="rail-health">checking</span></div>
        <div class="badge ok"><span class="led ok"></span><span id="rail-context">team-a / repo-main</span></div>
      </div>
    </aside>

    <div class="workspace">
      <header class="topbar">
        <div class="page-title">
          <h1 id="page-heading">Terminal</h1>
          <p id="page-subtitle">Interactive shell for login, auth flows, and CLI diagnostics.</p>
        </div>
        <div class="top-actions">
          <div class="token-field">
            <span>Token</span>
            <input id="token" type="password" autocomplete="off" placeholder="dev-token">
          </div>
          <button class="ghost" id="refresh" type="button">Refresh</button>
        </div>
      </header>

      <main class="content">
        <div class="status-grid">
          <div class="metric">
            <div class="metric-label">Health</div>
            <div class="metric-value"><span class="led" id="health-led"></span><span id="health">loading</span></div>
          </div>
          <div class="metric">
            <div class="metric-label">Ready</div>
            <div class="metric-value"><span class="led" id="ready-led"></span><span id="ready">loading</span></div>
          </div>
          <div class="metric">
            <div class="metric-label">Registered CLI</div>
            <div class="metric-value"><span id="tool-count">0</span></div>
          </div>
          <div class="metric">
            <div class="metric-label">Tenants</div>
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
                    <h2>Terminal</h2>
                    <p id="terminal-context">Select tenant, workspace, and credential profile.</p>
                  </div>
                </div>
                <span class="badge warn" id="connection-badge"><span class="led" id="terminal-led"></span><span id="terminal-state">disconnected</span></span>
              </div>

              <div class="context-bar">
                <div>
                  <label for="tenant">Tenant</label>
                  <select id="tenant"></select>
                </div>
                <div>
                  <label for="workspace">Workspace</label>
                  <input id="workspace" value="repo-main">
                </div>
                <div>
                  <label for="profile">Credential Profile</label>
                  <input id="profile" value="team-default" list="profile-options">
                  <datalist id="profile-options"></datalist>
                </div>
                <button class="primary" id="connect-terminal" type="button">Connect</button>
              </div>

              <div class="terminal-shell">
                <div class="terminal-screen" id="terminal-screen" tabindex="0" spellcheck="false" aria-label="Interactive terminal">
                  <pre class="terminal-output" id="terminal-output">Terminal is offline.

Connect to start an interactive shell. After connecting, click this terminal and type directly.
Credential state is isolated by tenant and profile.</pre>
                  <div class="screen-hint">
                    <span id="screen-hint-left">Click terminal to focus. Ctrl-C, arrows, tab, paste are supported.</span>
                    <span id="screen-hint-right" class="mono">120x32</span>
                  </div>
                </div>
                <div class="command-dock">
                  <span class="prompt-label">agent-runtime $</span>
                  <input id="command-input" class="mono" autocomplete="off" placeholder="Type a command, for example codex login" disabled>
                  <button id="send-command" type="button" disabled>Send</button>
                </div>
              </div>

              <div class="terminal-actions">
                <button class="ghost" id="disconnect-terminal" type="button" disabled>Disconnect</button>
                <button class="ghost" id="clear-terminal" type="button">Clear</button>
                <button class="ghost" id="ctrl-c" type="button" disabled>Ctrl-C</button>
                <button class="ghost" id="ctrl-l" type="button" disabled>Ctrl-L</button>
              </div>
            </section>

            <aside class="side-stack">
              <section class="panel">
                <div class="panel-header">
                  <div class="panel-title">
                    <span class="nav-icon">QL</span>
                    <div>
                      <h2>Quick Login</h2>
                      <p>Starts a terminal when needed, then runs the login command.</p>
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
                      <h2>CLI Manager</h2>
                      <p>Installed command entries and credential homes.</p>
                    </div>
                  </div>
                  <button class="ghost" type="button" data-view-jump="tools-view">Manage</button>
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
                    <h2>CLI Manager</h2>
                    <p>Add, update, or remove CLI definitions used by jobs and terminal credential envs.</p>
                  </div>
                </div>
              </div>
              <div class="table-wrap">
                <table>
                  <thead><tr><th>Name</th><th>Version</th><th>Path</th><th>Credential Env</th><th>Credential Subdir</th><th></th></tr></thead>
                  <tbody id="tools"><tr><td colspan="6" class="empty">Loading CLI tools</td></tr></tbody>
                </table>
              </div>
            </section>

            <aside class="panel">
              <div class="panel-header">
                <div class="panel-title">
                  <span class="nav-icon">IN</span>
                  <div>
                    <h2>Install CLI</h2>
                    <p>Register an existing binary path.</p>
                  </div>
                </div>
              </div>
              <div class="form-stack">
                <div>
                  <label for="tool-name">Name</label>
                  <input id="tool-name" placeholder="codex">
                </div>
                <div>
                  <label for="tool-path">Path</label>
                  <input id="tool-path" placeholder="/usr/local/bin/codex">
                </div>
                <div class="field-row">
                  <div>
                    <label for="tool-version">Version</label>
                    <input id="tool-version" placeholder="latest">
                  </div>
                  <div>
                    <label for="tool-env">Credential Env</label>
                    <input id="tool-env" placeholder="CODEX_HOME">
                  </div>
                </div>
                <div>
                  <label for="tool-subdir">Credential Subdir</label>
                  <input id="tool-subdir" placeholder=".codex">
                </div>
                <button class="primary" id="save-tool" type="button">Save CLI</button>
              </div>
            </aside>
          </div>
        </section>

        <section class="view" id="tenants-view">
          <section class="panel">
            <div class="panel-header">
              <div class="panel-title">
                <span class="nav-icon">TN</span>
                <div>
                  <h2>Tenants</h2>
                  <p>Token-derived tenant boundaries for tools, workspaces, credential profiles, and terminal access.</p>
                </div>
              </div>
            </div>
            <div class="tenant-grid" id="tenants"></div>
          </section>
        </section>
      </main>
    </div>
  </div>

  <div class="toast" id="toast"></div>

  <script>
    const state = {
      tools: [],
      tenants: [],
      ws: null,
      connected: false,
      pendingCommand: '',
      cols: 120,
      rows: 32
    };
    const tokenInput = document.getElementById('token');
    tokenInput.value = localStorage.getItem('agent-runtime-token') || 'dev-token';

    const pageTitles = {
      'terminal-view': ['Terminal', 'Interactive shell for login, auth flows, and CLI diagnostics.'],
      'tools-view': ['CLI Manager', 'Install, update, and remove CLI entries used by shared runtimes.'],
      'tenants-view': ['Tenants', 'Inspect tenant-scoped access boundaries and credential profiles.']
    };
    const loginCommands = [
      { label: 'Claude Code', command: 'claude login', tool: 'claude', env: 'CLAUDE_CONFIG_DIR' },
      { label: 'Codex', command: 'codex login', tool: 'codex', env: 'CODEX_HOME' },
      { label: 'Gemini', command: 'gemini', tool: 'gemini', env: 'GEMINI_HOME' },
      { label: 'OpenCode', command: 'opencode auth login', tool: 'opencode', env: 'OPENCODE_HOME' },
      { label: 'iFlow', command: 'iflow login', tool: 'iflow', env: 'IFLOW_HOME' },
      { label: 'Kimi', command: 'kimi login', tool: 'kimi', env: 'KIMI_HOME' },
      { label: 'Qoder', command: 'qodercli login', tool: 'qoder', env: 'QODER_HOME' }
    ];

    function $(id) { return document.getElementById(id); }
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
      $('rail-health').textContent = health.status === 'fulfilled' ? 'Health OK' : 'Offline';
      renderTools();
      renderToolCards();
      renderTenants();
      renderTerminalOptions();
      renderLoginShortcuts();
      updateContextLabels();
    }

    function setMetric(textID, ledID, value) {
      $(textID).textContent = value;
      const led = $(ledID);
      led.classList.remove('ok', 'bad');
      if (value === 'ok' || value === 'ready') led.classList.add('ok');
      if (value === 'error') led.classList.add('bad');
    }

    function renderTools() {
      const body = $('tools');
      if (!state.tools.length) {
        body.innerHTML = '<tr><td colspan="6"><div class="empty">No CLI tools registered</div></td></tr>';
        return;
      }
      body.innerHTML = state.tools.map(function(tool) {
        return '<tr>' +
          '<td><code>' + escapeHTML(tool.name) + '</code></td>' +
          '<td>' + escapeHTML(tool.version || '-') + '</td>' +
          '<td><code>' + escapeHTML(tool.path) + '</code></td>' +
          '<td><code>' + escapeHTML(tool.credential_env || '-') + '</code></td>' +
          '<td><code>' + escapeHTML(tool.credential_subdir || '-') + '</code></td>' +
          '<td><button class="danger" type="button" data-delete-tool="' + escapeHTML(tool.name) + '">Delete</button></td>' +
        '</tr>';
      }).join('');
      document.querySelectorAll('[data-delete-tool]').forEach(function(button) {
        button.addEventListener('click', async function() {
          await api('/api/tools/' + encodeURIComponent(button.dataset.deleteTool), { method: 'DELETE' });
          showToast('Deleted CLI ' + button.dataset.deleteTool);
          await refresh();
        });
      });
    }

    function renderToolCards() {
      const container = $('cli-cards');
      if (!state.tools.length) {
        container.innerHTML = '<div class="empty">No CLI tools registered</div>';
        return;
      }
      container.innerHTML = state.tools.slice(0, 6).map(function(tool) {
        return '<div class="tool-card">' +
          '<div class="tool-card-head">' +
            '<span class="tool-name">' + escapeHTML(tool.name) + '</span>' +
            '<span class="badge ok"><span class="led ok"></span>Registered</span>' +
          '</div>' +
          '<div class="tool-path">' + escapeHTML(tool.path) + '</div>' +
          '<div class="tool-path">' + escapeHTML(tool.credential_env || 'HOME') + ' -> ' + escapeHTML(tool.credential_subdir || '.') + '</div>' +
        '</div>';
      }).join('');
    }

    function renderTenants() {
      const container = $('tenants');
      if (!state.tenants.length) {
        container.innerHTML = '<div class="empty">No tenants configured</div>';
        return;
      }
      container.innerHTML = state.tenants.map(function(tenant) {
        return '<article class="tenant-card">' +
          '<div class="tenant-name"><span>' + escapeHTML(tenant.id) + '</span>' +
          '<span class="badge ' + (tenant.allow_terminal ? 'ok' : 'warn') + '">' + (tenant.allow_terminal ? 'Terminal allowed' : 'Terminal blocked') + '</span></div>' +
          '<div class="kv">' +
            '<div><strong>Subjects</strong><span>' + escapeHTML(join(tenant.subjects)) + '</span></div>' +
            '<div><strong>Tools</strong><span>' + escapeHTML(join(tenant.allowed_tools)) + '</span></div>' +
            '<div><strong>Workspaces</strong><span>' + escapeHTML(join(tenant.workspace_patterns)) + '</span></div>' +
            '<div><strong>Profiles</strong><span>' + escapeHTML(join(tenant.credential_profiles)) + '</span></div>' +
          '</div>' +
        '</article>';
      }).join('');
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

    function renderLoginShortcuts() {
      const knownTools = new Set(state.tools.map(function(tool) { return tool.name; }));
      $('login-shortcuts').innerHTML = loginCommands.map(function(item) {
        const registered = knownTools.has(item.tool);
        return '<button class="quick-button" type="button" data-login-command="' + escapeHTML(item.command) + '">' +
          '<span><strong>' + escapeHTML(item.label) + '</strong><span>' + escapeHTML(item.command) + '</span></span>' +
          '<span class="badge ' + (registered ? 'ok' : 'warn') + '">' + (registered ? 'Ready' : 'Add') + '</span>' +
        '</button>';
      }).join('');
      document.querySelectorAll('[data-login-command]').forEach(function(button) {
        button.addEventListener('click', function() { runCommand(button.dataset.loginCommand); });
      });
    }

    function updateContextLabels() {
      const context = ($('tenant').value || '-') + ' / ' + ($('workspace').value || '-');
      $('terminal-context').textContent = context + ' / ' + ($('profile').value || '-');
      $('rail-context').textContent = context;
    }

    function calculateTerminalSize() {
      const screen = $('terminal-screen');
      const width = Math.max(360, screen.clientWidth - 30);
      const height = Math.max(240, screen.clientHeight - 58);
      state.cols = Math.max(40, Math.min(240, Math.floor(width / 8.4)));
      state.rows = Math.max(12, Math.min(80, Math.floor(height / 18)));
      $('screen-hint-right').textContent = state.cols + 'x' + state.rows;
    }

    function terminalURL() {
      calculateTerminalSize();
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const params = new URLSearchParams({
        token: tokenInput.value.trim(),
        tenant: $('tenant').value,
        workspace: $('workspace').value.trim(),
        credential_profile: $('profile').value.trim(),
        cols: String(state.cols),
        rows: String(state.rows)
      });
      return protocol + '//' + window.location.host + '/api/terminal?' + params.toString();
    }

    function connectTerminal() {
      if (state.ws) state.ws.close();
      $('terminal-output').textContent = '';
      setConnected(false, 'connecting');
      localStorage.setItem('agent-runtime-token', tokenInput.value.trim());
      const ws = new WebSocket(terminalURL());
      state.ws = ws;
      ws.onopen = function() {
        setConnected(true, 'Connected');
        $('terminal-screen').focus();
        sendResize();
        if (state.pendingCommand) {
          const command = state.pendingCommand;
          state.pendingCommand = '';
          window.setTimeout(function() { sendTerminal(command); }, 180);
        }
      };
      ws.onmessage = function(event) {
        try {
          const payload = JSON.parse(event.data);
          if (payload.type === 'output') appendTerminal(payload.data || '');
          if (payload.type === 'error') appendTerminal('\n[terminal error] ' + (payload.data || 'unknown error') + '\n');
          if (payload.type === 'exit') setConnected(false, 'exited');
        } catch (err) {
          appendTerminal(String(event.data));
        }
      };
      ws.onclose = function() {
        if (state.ws === ws) state.ws = null;
        setConnected(false, 'disconnected');
      };
      ws.onerror = function() { setConnected(false, 'connection error'); };
    }

    function setConnected(connected, label) {
      state.connected = connected;
      const terminalState = label || (connected ? 'Connected' : 'disconnected');
      $('terminal-state').textContent = terminalState;
      $('connect-terminal').disabled = connected;
      $('disconnect-terminal').disabled = !connected;
      $('ctrl-c').disabled = !connected;
      $('ctrl-l').disabled = !connected;
      $('command-input').disabled = !connected;
      $('send-command').disabled = !connected;
      $('terminal-led').classList.remove('ok', 'bad');
      if (connected) $('terminal-led').classList.add('ok');
      if (terminalState === 'connection error') $('terminal-led').classList.add('bad');
      $('connection-badge').classList.toggle('ok', connected);
      $('connection-badge').classList.toggle('warn', !connected);
      $('screen-hint-left').textContent = connected ? 'Focused terminal accepts direct keyboard input and paste.' : 'Connect to start an interactive shell.';
    }

    function cleanTerminalText(data) {
      let text = String(data || '');
      if (text.indexOf('\x1b[2J') !== -1 || text.indexOf('\x1bc') !== -1) {
        $('terminal-output').textContent = '';
      }
      text = text.replace(/\x1B\][^\x07]*(\x07|\x1B\\)/g, '');
      text = text.replace(/\x1B\[[0-?]*[ -/]*[@-~]/g, '');
      text = text.replace(/\r\n/g, '\n').replace(/\r/g, '');
      return text;
    }

    function appendTerminal(text) {
      const output = $('terminal-output');
      const cleaned = cleanTerminalText(text);
      if (output.textContent.indexOf('Terminal is offline.') === 0) output.textContent = '';
      output.textContent += cleaned;
      const screen = $('terminal-screen');
      screen.scrollTop = screen.scrollHeight;
    }

    function sendTerminal(data) {
      if (!state.ws || state.ws.readyState !== WebSocket.OPEN) return;
      state.ws.send(JSON.stringify({ type: 'input', data: data }));
    }

    function sendResize() {
      if (!state.ws || state.ws.readyState !== WebSocket.OPEN) return;
      calculateTerminalSize();
      state.ws.send(JSON.stringify({ type: 'resize', cols: state.cols, rows: state.rows }));
    }

    function runCommand(command) {
      const payload = command + '\r';
      if (state.connected) {
        sendTerminal(payload);
        $('terminal-screen').focus();
        return;
      }
      state.pendingCommand = payload;
      connectTerminal();
    }

    function handleTerminalKey(event) {
      if (!state.connected) return;
      let data = '';
      if (event.ctrlKey && !event.metaKey) {
        const key = event.key.toLowerCase();
        if (key === 'c') data = '\x03';
        if (key === 'd') data = '\x04';
        if (key === 'l') data = '\x0c';
      } else if (event.key === 'Enter') data = '\r';
      else if (event.key === 'Backspace') data = '\x7f';
      else if (event.key === 'Tab') data = '\t';
      else if (event.key === 'Escape') data = '\x1b';
      else if (event.key === 'ArrowUp') data = '\x1b[A';
      else if (event.key === 'ArrowDown') data = '\x1b[B';
      else if (event.key === 'ArrowRight') data = '\x1b[C';
      else if (event.key === 'ArrowLeft') data = '\x1b[D';
      else if (event.key === 'Home') data = '\x1b[H';
      else if (event.key === 'End') data = '\x1b[F';
      else if (event.key === 'Delete') data = '\x1b[3~';
      else if (!event.metaKey && !event.altKey && event.key.length === 1) data = event.key;
      if (data) {
        event.preventDefault();
        sendTerminal(data);
      }
    }

    function switchView(viewID) {
      document.querySelectorAll('.nav-button').forEach(function(item) {
        item.classList.toggle('active', item.dataset.view === viewID);
      });
      document.querySelectorAll('.view').forEach(function(item) {
        item.classList.toggle('active', item.id === viewID);
      });
      const title = pageTitles[viewID] || pageTitles['terminal-view'];
      $('page-heading').textContent = title[0];
      $('page-subtitle').textContent = title[1];
    }

    document.querySelectorAll('[data-view]').forEach(function(button) {
      button.addEventListener('click', function() { switchView(button.dataset.view); });
    });
    document.querySelectorAll('[data-view-jump]').forEach(function(button) {
      button.addEventListener('click', function() { switchView(button.dataset.viewJump); });
    });
    $('refresh').addEventListener('click', function() {
      refresh().then(function() { showToast('Runtime status refreshed'); }).catch(function(err) { showToast(err.message); });
    });
    $('tenant').addEventListener('change', updateProfileOptions);
    $('workspace').addEventListener('input', updateContextLabels);
    $('profile').addEventListener('input', updateContextLabels);
    $('connect-terminal').addEventListener('click', connectTerminal);
    $('disconnect-terminal').addEventListener('click', function() { if (state.ws) state.ws.close(); });
    $('clear-terminal').addEventListener('click', function() { $('terminal-output').textContent = ''; });
    $('ctrl-c').addEventListener('click', function() { sendTerminal('\x03'); $('terminal-screen').focus(); });
    $('ctrl-l').addEventListener('click', function() { sendTerminal('\x0c'); $('terminal-screen').focus(); });
    $('terminal-screen').addEventListener('keydown', handleTerminalKey);
    $('terminal-screen').addEventListener('paste', function(event) {
      if (!state.connected) return;
      event.preventDefault();
      sendTerminal((event.clipboardData || window.clipboardData).getData('text'));
    });
    $('command-input').addEventListener('keydown', function(event) {
      if (event.key === 'Enter') {
        event.preventDefault();
        $('send-command').click();
      }
    });
    $('send-command').addEventListener('click', function() {
      const input = $('command-input');
      if (!input.value.trim()) return;
      sendTerminal(input.value + '\r');
      input.value = '';
      $('terminal-screen').focus();
    });
    $('save-tool').addEventListener('click', async function() {
      await api('/api/tools', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: $('tool-name').value.trim(),
          path: $('tool-path').value.trim(),
          version: $('tool-version').value.trim(),
          credential_env: $('tool-env').value.trim(),
          credential_subdir: $('tool-subdir').value.trim()
        })
      });
      ['tool-name', 'tool-path', 'tool-version', 'tool-env', 'tool-subdir'].forEach(function(id) { $(id).value = ''; });
      showToast('CLI saved');
      await refresh();
    });

    calculateTerminalSize();
    if (window.ResizeObserver) {
      new ResizeObserver(function() { sendResize(); }).observe($('terminal-screen'));
    } else {
      window.addEventListener('resize', sendResize);
    }
    refresh().catch(function(err) {
      setMetric('health', 'health-led', 'error');
      setMetric('ready', 'ready-led', 'error');
      appendTerminal('[ui error] ' + err.message + '\n');
    });
  </script>
</body>
</html>`
