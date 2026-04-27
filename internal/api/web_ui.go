package api

const webUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Agent Runtime</title>
  <style>
    :root {
      color-scheme: light dark;
      --bg: #f6f7f9;
      --panel: #ffffff;
      --text: #17202a;
      --muted: #647184;
      --line: #d8dee8;
      --accent: #0f766e;
      --accent-strong: #115e59;
      --danger: #b42318;
      --terminal: #05070a;
      --terminal-text: #d8dee9;
      --mono: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg: #111418;
        --panel: #181d23;
        --text: #eef2f7;
        --muted: #a2adbd;
        --line: #2b333d;
        --accent: #2dd4bf;
        --accent-strong: #5eead4;
      }
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      line-height: 1.45;
    }
    header {
      border-bottom: 1px solid var(--line);
      background: var(--panel);
    }
    .shell {
      width: min(1240px, calc(100% - 32px));
      margin: 0 auto;
    }
    .topbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      min-height: 68px;
    }
    h1 {
      margin: 0;
      font-size: 20px;
      font-weight: 700;
      letter-spacing: 0;
    }
    h2 {
      margin: 0 0 12px;
      font-size: 15px;
      font-weight: 700;
      letter-spacing: 0;
    }
    .subtle { color: var(--muted); }
    .token-row {
      display: grid;
      grid-template-columns: 180px auto;
      gap: 8px;
      align-items: center;
    }
    main { padding: 18px 0 40px; }
    .tabs {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      margin-bottom: 14px;
    }
    .tab {
      border: 1px solid var(--line);
      background: var(--panel);
      color: var(--text);
      min-height: 36px;
      padding: 7px 12px;
      border-radius: 6px;
    }
    .tab.active {
      border-color: var(--accent);
      color: #ffffff;
      background: var(--accent);
    }
    .view { display: none; }
    .view.active { display: block; }
    .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 16px;
    }
    .panel + .panel { margin-top: 14px; }
    .toolbar {
      display: flex;
      align-items: end;
      gap: 10px;
      flex-wrap: wrap;
    }
    .toolbar > div { min-width: 170px; flex: 1; }
    .terminal-layout {
      display: grid;
      grid-template-columns: minmax(0, 1fr) 280px;
      gap: 14px;
      align-items: start;
    }
    label {
      display: block;
      margin: 0 0 6px;
      font-size: 12px;
      font-weight: 700;
      color: var(--muted);
      text-transform: uppercase;
    }
    code, pre, input, textarea, select {
      font-family: var(--mono);
      font-size: 13px;
    }
    input, textarea, select {
      width: 100%;
      min-height: 38px;
      padding: 8px 10px;
      border: 1px solid var(--line);
      border-radius: 6px;
      background: transparent;
      color: var(--text);
    }
    textarea {
      min-height: 58px;
      resize: vertical;
    }
    button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 38px;
      padding: 8px 12px;
      border: 1px solid var(--accent);
      border-radius: 6px;
      background: var(--accent);
      color: #ffffff;
      font-weight: 700;
      cursor: pointer;
    }
    button.secondary {
      background: transparent;
      color: var(--accent-strong);
    }
    button.danger {
      border-color: var(--danger);
      background: transparent;
      color: var(--danger);
    }
    button:disabled {
      cursor: not-allowed;
      opacity: 0.5;
    }
    .actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      margin-top: 12px;
    }
    .quick-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: 8px;
    }
    .terminal {
      margin: 0;
      min-height: 460px;
      max-height: 60vh;
      overflow: auto;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: var(--terminal);
      color: var(--terminal-text);
      white-space: pre-wrap;
      overflow-wrap: anywhere;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 14px;
    }
    th, td {
      text-align: left;
      padding: 10px 8px;
      border-bottom: 1px solid var(--line);
      vertical-align: top;
    }
    th {
      font-size: 12px;
      color: var(--muted);
      font-weight: 700;
      text-transform: uppercase;
    }
    .form-grid {
      display: grid;
      grid-template-columns: repeat(5, minmax(130px, 1fr));
      gap: 10px;
      align-items: end;
    }
    .status-strip {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
      margin-bottom: 14px;
    }
    .pill {
      border: 1px solid var(--line);
      border-radius: 6px;
      padding: 7px 10px;
      background: var(--panel);
      color: var(--muted);
      font-size: 13px;
    }
    .error { color: var(--danger); }
    @media (max-width: 900px) {
      .topbar, .terminal-layout, .form-grid, .token-row { grid-template-columns: 1fr; }
      .topbar { align-items: flex-start; flex-direction: column; padding: 14px 0; }
      .terminal { min-height: 360px; }
    }
  </style>
</head>
<body>
  <header>
    <div class="shell topbar">
      <div>
        <h1>Agent Runtime</h1>
        <div class="subtle">Shared CLI installs, login state, credential homes, and tenant workspaces</div>
      </div>
      <div class="token-row">
        <input id="token" type="password" autocomplete="off" placeholder="dev-token">
        <button class="secondary" id="refresh" type="button">Refresh</button>
      </div>
    </div>
  </header>

  <main class="shell">
    <nav class="tabs" aria-label="Agent Runtime sections">
      <button class="tab active" type="button" data-tab="terminal-view">Terminal</button>
      <button class="tab" type="button" data-tab="tools-view">CLI Manager</button>
      <button class="tab" type="button" data-tab="tenants-view">Tenants</button>
    </nav>

    <div class="status-strip">
      <div class="pill">Health: <strong id="health">loading</strong></div>
      <div class="pill">Ready: <strong id="ready">loading</strong></div>
      <div class="pill">Tools: <strong id="tool-count">0</strong></div>
      <div class="pill">Tenants: <strong id="tenant-count">0</strong></div>
      <div class="pill">Terminal: <strong id="terminal-state">disconnected</strong></div>
    </div>

    <section class="view active" id="terminal-view">
      <div class="terminal-layout">
        <div class="panel">
          <h2>Terminal</h2>
          <div class="toolbar">
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
            <button id="connect-terminal" type="button">Connect</button>
          </div>
          <div class="actions">
            <button class="secondary" id="disconnect-terminal" type="button" disabled>Disconnect</button>
            <button class="secondary" id="clear-terminal" type="button">Clear</button>
            <button class="secondary" id="ctrl-c" type="button" disabled>Ctrl-C</button>
          </div>
          <pre class="terminal" id="terminal-output">Terminal is disconnected.</pre>
          <label for="terminal-input" style="margin-top: 12px;">Input</label>
          <textarea id="terminal-input" placeholder="Type a command or response, then press Enter. Use Shift+Enter for a newline." disabled></textarea>
          <div class="actions">
            <button id="send-terminal" type="button" disabled>Send</button>
          </div>
        </div>

        <aside class="panel">
          <h2>Login Shortcuts</h2>
          <div class="quick-grid" id="login-shortcuts"></div>
        </aside>
      </div>
    </section>

    <section class="view" id="tools-view">
      <div class="panel">
        <h2>CLI Manager</h2>
        <table>
          <thead><tr><th>Name</th><th>Version</th><th>Path</th><th>Credential Env</th><th>Credential Subdir</th><th></th></tr></thead>
          <tbody id="tools"><tr><td colspan="6" class="subtle">Loading tools</td></tr></tbody>
        </table>
      </div>

      <div class="panel">
        <h2>Add Or Update CLI</h2>
        <div class="form-grid">
          <div>
            <label for="tool-name">Name</label>
            <input id="tool-name" placeholder="codex">
          </div>
          <div>
            <label for="tool-path">Path</label>
            <input id="tool-path" placeholder="/usr/local/bin/codex">
          </div>
          <div>
            <label for="tool-version">Version</label>
            <input id="tool-version" placeholder="latest">
          </div>
          <div>
            <label for="tool-env">Credential Env</label>
            <input id="tool-env" placeholder="CODEX_HOME">
          </div>
          <div>
            <label for="tool-subdir">Credential Subdir</label>
            <input id="tool-subdir" placeholder=".codex">
          </div>
        </div>
        <div class="actions">
          <button id="save-tool" type="button">Save CLI</button>
        </div>
      </div>
    </section>

    <section class="view" id="tenants-view">
      <div class="panel">
        <h2>Tenants</h2>
        <table>
          <thead><tr><th>Tenant</th><th>Subjects</th><th>Tools</th><th>Workspaces</th><th>Credential Profiles</th><th>Terminal</th></tr></thead>
          <tbody id="tenants"><tr><td colspan="6" class="subtle">Loading tenants</td></tr></tbody>
        </table>
      </div>
    </section>
  </main>

  <script>
    const state = { tools: [], tenants: [], ws: null, connected: false };
    const tokenInput = document.getElementById('token');
    tokenInput.value = localStorage.getItem('agent-runtime-token') || 'dev-token';

    const loginCommands = [
      { label: 'Claude Code', command: 'claude login', tool: 'claude' },
      { label: 'Codex', command: 'codex login', tool: 'codex' },
      { label: 'Gemini', command: 'gemini', tool: 'gemini' },
      { label: 'OpenCode', command: 'opencode auth login', tool: 'opencode' },
      { label: 'iFlow', command: 'iflow login', tool: 'iflow' },
      { label: 'Kimi', command: 'kimi login', tool: 'kimi' },
      { label: 'Qoder', command: 'qodercli login', tool: 'qoder' }
    ];

    function $(id) { return document.getElementById(id); }
    function escapeHTML(value) {
      return String(value ?? '').replace(/[&<>"']/g, (char) => ({
        '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'
      })[char]);
    }
    function join(values) { return values && values.length ? values.join(', ') : '-'; }

    async function api(path, options = {}) {
      const response = await fetch(path, options);
      if (!response.ok) {
        let message = response.statusText;
        try {
          const body = await response.json();
          message = body.error || message;
        } catch {}
        throw new Error(message);
      }
      if (response.status === 204) return null;
      return response.json();
    }

    async function refresh() {
      localStorage.setItem('agent-runtime-token', tokenInput.value.trim());
      const [health, ready, status, tools, tenants] = await Promise.allSettled([
        api('/api/health'),
        api('/api/ready'),
        api('/api/status'),
        api('/api/tools'),
        api('/api/tenants')
      ]);
      $('health').textContent = health.status === 'fulfilled' ? health.value.status : 'error';
      $('ready').textContent = ready.status === 'fulfilled' ? ready.value.status : 'error';
      if (status.status === 'fulfilled') {
        $('tool-count').textContent = status.value.tools ?? 0;
        $('tenant-count').textContent = status.value.tenants ?? 0;
      }
      state.tools = tools.status === 'fulfilled' ? tools.value.tools : [];
      state.tenants = tenants.status === 'fulfilled' ? tenants.value.tenants : [];
      renderTools();
      renderTenants();
      renderTerminalOptions();
      renderLoginShortcuts();
    }

    function renderTools() {
      const body = $('tools');
      if (!state.tools.length) {
        body.innerHTML = '<tr><td colspan="6" class="subtle">No CLI tools registered</td></tr>';
        return;
      }
      body.innerHTML = state.tools.map((tool) =>
        '<tr>' +
          '<td><code>' + escapeHTML(tool.name) + '</code></td>' +
          '<td>' + escapeHTML(tool.version || '-') + '</td>' +
          '<td><code>' + escapeHTML(tool.path) + '</code></td>' +
          '<td><code>' + escapeHTML(tool.credential_env || '-') + '</code></td>' +
          '<td><code>' + escapeHTML(tool.credential_subdir || '-') + '</code></td>' +
          '<td><button class="danger" type="button" data-delete-tool="' + escapeHTML(tool.name) + '">Delete</button></td>' +
        '</tr>'
      ).join('');
      document.querySelectorAll('[data-delete-tool]').forEach((button) => {
        button.addEventListener('click', async () => {
          await api('/api/tools/' + encodeURIComponent(button.dataset.deleteTool), { method: 'DELETE' });
          await refresh();
        });
      });
    }

    function renderTenants() {
      const body = $('tenants');
      if (!state.tenants.length) {
        body.innerHTML = '<tr><td colspan="6" class="subtle">No tenants configured</td></tr>';
        return;
      }
      body.innerHTML = state.tenants.map((tenant) =>
        '<tr>' +
          '<td><code>' + escapeHTML(tenant.id) + '</code></td>' +
          '<td>' + escapeHTML(join(tenant.subjects)) + '</td>' +
          '<td>' + escapeHTML(join(tenant.allowed_tools)) + '</td>' +
          '<td>' + escapeHTML(join(tenant.workspace_patterns)) + '</td>' +
          '<td>' + escapeHTML(join(tenant.credential_profiles)) + '</td>' +
          '<td>' + (tenant.allow_terminal ? 'allowed' : 'blocked') + '</td>' +
        '</tr>'
      ).join('');
    }

    function renderTerminalOptions() {
      const tenantSelect = $('tenant');
      const currentTenant = tenantSelect.value;
      tenantSelect.innerHTML = state.tenants.map((tenant) =>
        '<option value="' + escapeHTML(tenant.id) + '">' + escapeHTML(tenant.id) + '</option>'
      ).join('');
      if (currentTenant) tenantSelect.value = currentTenant;
      if (!tenantSelect.value && state.tenants[0]) tenantSelect.value = state.tenants[0].id;
      updateProfileOptions();
    }

    function updateProfileOptions() {
      const tenant = state.tenants.find((item) => item.id === $('tenant').value);
      const profiles = tenant?.credential_profiles || [];
      $('profile-options').innerHTML = profiles.map((profile) =>
        '<option value="' + escapeHTML(profile) + '"></option>'
      ).join('');
      if (profiles.length && !$('profile').value) $('profile').value = profiles[0];
      const workspaces = tenant?.workspace_patterns || [];
      if (!$('workspace').value && workspaces[0]) $('workspace').value = workspaces[0].replace('*', 'main');
    }

    function renderLoginShortcuts() {
      const knownTools = new Set(state.tools.map((tool) => tool.name));
      $('login-shortcuts').innerHTML = loginCommands.map((item) => {
        const registered = knownTools.has(item.tool);
        return '<button class="secondary" type="button" data-login-command="' + escapeHTML(item.command) + '" ' +
          (state.connected ? '' : 'disabled') + '>' +
          escapeHTML(item.label) + (registered ? '' : ' (not registered)') +
        '</button>';
      }).join('');
      document.querySelectorAll('[data-login-command]').forEach((button) => {
        button.addEventListener('click', () => sendTerminal(button.dataset.loginCommand + '\r'));
      });
    }

    function appendTerminal(text) {
      const output = $('terminal-output');
      if (output.textContent === 'Terminal is disconnected.') output.textContent = '';
      output.textContent += text;
      output.scrollTop = output.scrollHeight;
    }

    function setConnected(connected, label) {
      state.connected = connected;
      $('terminal-state').textContent = label || (connected ? 'connected' : 'disconnected');
      $('connect-terminal').disabled = connected;
      $('disconnect-terminal').disabled = !connected;
      $('ctrl-c').disabled = !connected;
      $('terminal-input').disabled = !connected;
      $('send-terminal').disabled = !connected;
      renderLoginShortcuts();
    }

    function terminalURL() {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const params = new URLSearchParams({
        token: tokenInput.value.trim(),
        tenant: $('tenant').value,
        workspace: $('workspace').value.trim(),
        credential_profile: $('profile').value.trim(),
        cols: '120',
        rows: '32'
      });
      return protocol + '//' + window.location.host + '/api/terminal?' + params.toString();
    }

    function connectTerminal() {
      if (state.ws) state.ws.close();
      $('terminal-output').textContent = '';
      setConnected(false, 'connecting');
      const ws = new WebSocket(terminalURL());
      state.ws = ws;
      ws.onopen = () => setConnected(true, 'connected');
      ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          if (payload.type === 'output') appendTerminal(payload.data || '');
          if (payload.type === 'error') appendTerminal('\r\n[terminal error] ' + (payload.data || 'unknown error') + '\r\n');
          if (payload.type === 'exit') setConnected(false, 'exited');
        } catch {
          appendTerminal(String(event.data));
        }
      };
      ws.onclose = () => {
        if (state.ws === ws) state.ws = null;
        setConnected(false, 'disconnected');
      };
      ws.onerror = () => setConnected(false, 'connection error');
    }

    function sendTerminal(data) {
      if (!state.ws || state.ws.readyState !== WebSocket.OPEN) return;
      state.ws.send(JSON.stringify({ type: 'input', data }));
    }

    document.querySelectorAll('.tab').forEach((button) => {
      button.addEventListener('click', () => {
        document.querySelectorAll('.tab').forEach((item) => item.classList.remove('active'));
        document.querySelectorAll('.view').forEach((item) => item.classList.remove('active'));
        button.classList.add('active');
        $(button.dataset.tab).classList.add('active');
      });
    });
    $('refresh').addEventListener('click', () => refresh().catch((err) => alert(err.message)));
    $('tenant').addEventListener('change', updateProfileOptions);
    $('connect-terminal').addEventListener('click', connectTerminal);
    $('disconnect-terminal').addEventListener('click', () => state.ws?.close());
    $('clear-terminal').addEventListener('click', () => $('terminal-output').textContent = '');
    $('ctrl-c').addEventListener('click', () => sendTerminal('\x03'));
    $('send-terminal').addEventListener('click', () => {
      const input = $('terminal-input');
      sendTerminal(input.value + '\r');
      input.value = '';
    });
    $('terminal-input').addEventListener('keydown', (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        $('send-terminal').click();
      }
    });
    $('save-tool').addEventListener('click', async () => {
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
      ['tool-name', 'tool-path', 'tool-version', 'tool-env', 'tool-subdir'].forEach((id) => $(id).value = '');
      await refresh();
    });

    refresh().catch((err) => {
      $('health').textContent = 'error';
      $('ready').textContent = 'error';
      appendTerminal('[ui error] ' + err.message + '\r\n');
    });
  </script>
</body>
</html>`
