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
      --bg: #f7f8fa;
      --panel: #ffffff;
      --text: #18202a;
      --muted: #5b6675;
      --line: #d9dee7;
      --accent: #0f766e;
      --accent-strong: #115e59;
      --danger: #b42318;
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
      width: min(1180px, calc(100% - 32px));
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
    .subtle { color: var(--muted); }
    main {
      padding: 24px 0 40px;
    }
    .grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr) minmax(340px, 420px);
      gap: 16px;
      align-items: start;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 16px;
    }
    .panel + .panel { margin-top: 16px; }
    h2 {
      margin: 0 0 12px;
      font-size: 15px;
      font-weight: 700;
      letter-spacing: 0;
    }
    .metric-row {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 12px;
    }
    .metric {
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 12px;
      min-height: 76px;
    }
    .metric .label {
      font-size: 12px;
      color: var(--muted);
    }
    .metric .value {
      margin-top: 8px;
      font-size: 18px;
      font-weight: 700;
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
    code, pre, input, textarea, select {
      font-family: var(--mono);
      font-size: 13px;
    }
    label {
      display: block;
      margin: 12px 0 6px;
      font-size: 12px;
      font-weight: 700;
      color: var(--muted);
      text-transform: uppercase;
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
      min-height: 76px;
      resize: vertical;
    }
    .split {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 10px;
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
    .actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
      margin-top: 14px;
    }
    pre {
      margin: 0;
      min-height: 160px;
      max-height: 420px;
      overflow: auto;
      padding: 12px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: color-mix(in srgb, var(--panel), #000 4%);
      white-space: pre-wrap;
      overflow-wrap: anywhere;
    }
    .error { color: var(--danger); }
    @media (max-width: 860px) {
      .grid, .metric-row, .split { grid-template-columns: 1fr; }
      .topbar { align-items: flex-start; flex-direction: column; padding: 14px 0; }
    }
  </style>
</head>
<body>
  <header>
    <div class="shell topbar">
      <div>
        <h1>Agent Runtime</h1>
        <div class="subtle">Portable CLI agent runtime control surface</div>
      </div>
      <button class="secondary" id="refresh" type="button">Refresh</button>
    </div>
  </header>

  <main class="shell">
    <div class="grid">
      <section>
        <div class="panel">
          <h2>Status</h2>
          <div class="metric-row">
            <div class="metric"><div class="label">Health</div><div class="value" id="health">loading</div></div>
            <div class="metric"><div class="label">Ready</div><div class="value" id="ready">loading</div></div>
            <div class="metric"><div class="label">Tools</div><div class="value" id="tool-count">0</div></div>
          </div>
        </div>

        <div class="panel">
          <h2>Tools</h2>
          <table>
            <thead><tr><th>Name</th><th>Version</th><th>Path</th><th>Credential Home</th></tr></thead>
            <tbody id="tools"><tr><td colspan="4" class="subtle">Loading tools</td></tr></tbody>
          </table>
        </div>

        <div class="panel">
          <h2>Last Result</h2>
          <pre id="output">No job submitted yet.</pre>
        </div>
      </section>

      <aside class="panel">
        <h2>Create Job</h2>
        <label for="token">Bearer Token</label>
        <input id="token" type="password" autocomplete="off" placeholder="dev-token">

        <div class="split">
          <div>
            <label for="tenant">Tenant</label>
            <input id="tenant" value="team-a">
          </div>
          <div>
            <label for="profile">Credential Profile</label>
            <input id="profile" value="team-default">
          </div>
        </div>

        <div class="split">
          <div>
            <label for="workspace">Workspace</label>
            <input id="workspace" value="repo-main">
          </div>
          <div>
            <label for="tool">Tool</label>
            <select id="tool"></select>
          </div>
        </div>

        <label for="args">Args JSON Array</label>
        <textarea id="args">[]</textarea>

        <div class="split">
          <div>
            <label for="timeout">Timeout Seconds</label>
            <input id="timeout" type="number" min="1" value="60">
          </div>
          <div>
            <label for="job-id">Job ID</label>
            <input id="job-id" placeholder="created job id">
          </div>
        </div>

        <div class="actions">
          <button id="submit" type="button">Run Job</button>
          <button class="secondary" id="fetch-job" type="button">Fetch Job</button>
          <button class="secondary" id="events" type="button">Load Events</button>
        </div>
      </aside>
    </div>
  </main>

  <script>
    const state = { tools: [] };
    const $ = (id) => document.getElementById(id);

    function writeOutput(value) {
      $("output").textContent = typeof value === "string" ? value : JSON.stringify(value, null, 2);
    }

    async function request(path, options = {}) {
      const response = await fetch(path, options);
      const text = await response.text();
      let payload = text;
      try { payload = text ? JSON.parse(text) : null; } catch (_) {}
      if (!response.ok) {
        throw new Error(typeof payload === "object" && payload && payload.error ? payload.error : text || response.statusText);
      }
      return payload;
    }

    async function refresh() {
      $("health").textContent = "loading";
      $("ready").textContent = "loading";
      try {
        const [health, ready, tools] = await Promise.all([
          request("/api/health"),
          request("/api/ready"),
          request("/api/tools")
        ]);
        $("health").textContent = health.status || "ok";
        $("ready").textContent = ready.status || "ready";
        state.tools = tools.tools || [];
        $("tool-count").textContent = String(state.tools.length);
        renderTools();
      } catch (error) {
        $("health").innerHTML = '<span class="error">error</span>';
        $("ready").innerHTML = '<span class="error">error</span>';
        writeOutput(error.message);
      }
    }

    function renderTools() {
      const body = $("tools");
      const select = $("tool");
      body.innerHTML = "";
      select.innerHTML = "";
      if (state.tools.length === 0) {
        body.innerHTML = '<tr><td colspan="4" class="subtle">No tools registered</td></tr>';
        return;
      }
      for (const tool of state.tools) {
        const row = document.createElement("tr");
        row.innerHTML = '<td><code></code></td><td></td><td><code></code></td><td><code></code></td>';
        row.children[0].firstChild.textContent = tool.name || "";
        row.children[1].textContent = tool.version || "";
        row.children[2].firstChild.textContent = tool.path || "";
        row.children[3].firstChild.textContent = tool.credential_env || "";
        body.appendChild(row);

        const option = document.createElement("option");
        option.value = tool.name;
        option.textContent = tool.name;
        select.appendChild(option);
      }
    }

    async function submitJob() {
      const token = $("token").value.trim();
      localStorage.setItem("agentRuntimeToken", token);
      const payload = {
        tenant: $("tenant").value.trim(),
        tool: $("tool").value,
        args: JSON.parse($("args").value || "[]"),
        workspace: $("workspace").value.trim(),
        credential_profile: $("profile").value.trim(),
        timeout_seconds: Number($("timeout").value || 60)
      };
      const job = await request("/api/jobs", {
        method: "POST",
        headers: {
          "Authorization": "Bearer " + token,
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      });
      $("job-id").value = job.id;
      writeOutput(job);
    }

    async function fetchJob() {
      const id = $("job-id").value.trim();
      if (!id) throw new Error("Job ID is required");
      writeOutput(await request("/api/jobs/" + encodeURIComponent(id)));
    }

    async function loadEvents() {
      const id = $("job-id").value.trim();
      if (!id) throw new Error("Job ID is required");
      const text = await fetch("/api/jobs/" + encodeURIComponent(id) + "/events").then((r) => r.text());
      writeOutput(text || "No events");
    }

    $("refresh").addEventListener("click", () => refresh());
    $("submit").addEventListener("click", () => submitJob().catch((error) => writeOutput(error.message)));
    $("fetch-job").addEventListener("click", () => fetchJob().catch((error) => writeOutput(error.message)));
    $("events").addEventListener("click", () => loadEvents().catch((error) => writeOutput(error.message)));
    $("token").value = localStorage.getItem("agentRuntimeToken") || "";
    refresh();
  </script>
</body>
</html>`
