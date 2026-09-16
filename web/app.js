async function loadStats() {
  const res = await fetch('/api/stats/overview');
  const body = await res.json();
  const data = body.data || {};
  const grid = document.getElementById('stats');
  grid.innerHTML = `
    <div class="stat-card"><div class="stat-title">会话总数</div><div class="stat-value">${data.total_sessions || 0}</div></div>
    <div class="stat-card"><div class="stat-title">任务总数</div><div class="stat-value">${data.total_tasks || 0}</div></div>
    <div class="stat-card"><div class="stat-title">流量记录</div><div class="stat-value">${data.total_records || 0}</div></div>
    <div class="stat-card"><div class="stat-title">环境数</div><div class="stat-value">${data.total_envs || 0}</div></div>
    <div class="stat-card"><div class="stat-title">成功率(%)</div><div class="stat-value">${(data.success_rate || 0).toFixed(2)}</div></div>
    <div class="stat-card"><div class="stat-title">匹配率(%)</div><div class="stat-value">${(data.match_rate || 0).toFixed(2)}</div></div>
  `;
}

async function loadTasks() {
  const res = await fetch('/api/replay-tasks');
  const body = await res.json();
  const items = (body.data && body.data.items) || [];
  const tbody = document.querySelector('#task-list tbody');
  tbody.innerHTML = '';
  items.forEach(t => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${t.id}</td><td>${t.session_id}</td><td>${t.target_env_id}</td><td>${t.status}</td><td>${t.total_requests}</td><td>${t.success_count}</td><td>${t.failed_count}</td><td>${t.created_at}</td>`;
    tbody.appendChild(tr);
  });
}

async function loadRecords() {
  const res = await fetch('/api/traffic-records');
  const body = await res.json();
  const items = (body.data && body.data.items) || [];
  const tbody = document.querySelector('#record-list tbody');
  tbody.innerHTML = '';
  items.forEach(r => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${r.id}</td><td>${r.session_id}</td><td>${r.method}</td><td>${r.path}</td><td>${r.status_code}</td><td>${r.duration_ms}</td><td>${r.timestamp}</td>`;
    tbody.appendChild(tr);
  });
}

async function load() {
  await loadStats();
  await loadTasks();
  await loadRecords();
}
load();
