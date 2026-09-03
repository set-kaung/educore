// Shared glue between the static HTML pages and the JSON API.
// All business logic lives behind /api/*; this file only handles
// sessions, navigation, redirects, error toasts and rendering
// API data into tables and selects.

const sessionCache = { value: undefined };

export async function getSession(force = false) {
  if (!force && sessionCache.value !== undefined) {
    return sessionCache.value;
  }
  try {
    const res = await fetch("/api/session", { credentials: "same-origin" });
    sessionCache.value = res.ok ? (await res.json()).data : null;
  } catch {
    sessionCache.value = null;
  }
  return sessionCache.value;
}

export async function logout() {
  await fetch("/api/logout", { method: "POST", credentials: "same-origin" });
  sessionCache.value = null;
}

export function homeFor(session) {
  if (!session) return "/login";
  return session.role === "student" ? "/my-courses" : "/students";
}

function applyNav(session) {
  const nav = document.getElementById("nav");
  if (!nav) return;

  for (const link of nav.querySelectorAll("[data-role]")) {
    const roles = link.dataset.role.split(" ");
    link.hidden = !session || !roles.includes(session.role);
  }
  const guestLink = nav.querySelector("[data-guest]");
  if (guestLink) guestLink.hidden = !!session;

  const badge = document.getElementById("role-badge");
  if (badge && session) {
    badge.textContent = session.role;
    badge.hidden = false;
  }

  const logoutBtn = document.getElementById("logout-btn");
  if (logoutBtn) {
    logoutBtn.hidden = !session;
    logoutBtn.onclick = async () => {
      await logout();
      location.assign("/login");
    };
  }

  const current = document.body.dataset.page;
  for (const link of nav.querySelectorAll("[data-page]")) {
    if (link.dataset.page === current) link.classList.add("nav-link-active");
  }

  nav.hidden = false;
}

export function showToast(message, isError = false) {
  let host = document.getElementById("toast-host");
  if (!host) {
    host = document.createElement("div");
    host.id = "toast-host";
    host.className = "toast-host";
    document.body.appendChild(host);
  }
  const toast = document.createElement("div");
  toast.textContent = message;
  toast.className = isError ? "toast toast-error" : "toast";
  host.appendChild(toast);
  setTimeout(() => toast.remove(), 4000);
}

function popFlash() {
  const message = sessionStorage.getItem("edu_flash");
  if (message) {
    sessionStorage.removeItem("edu_flash");
    showToast(message);
  }
}

export function flash(message) {
  sessionStorage.setItem("edu_flash", message);
}

function guardPage(session, options) {
  if (options.redirectHome) {
    location.replace(homeFor(session));
    return false;
  }
  if (options.guestOnly) {
    if (session) location.replace(homeFor(session));
    return !session;
  }
  if (!session) {
    location.replace("/login");
    return false;
  }
  const required = options.require;
  if (required && !required.includes(session.role)) {
    location.replace(homeFor(session));
    return false;
  }
  return true;
}

const escapeMap = {
  "&": "&amp;",
  "<": "&lt;",
  ">": "&gt;",
  '"': "&quot;",
  "'": "&#39;",
  "`": "&#x60;",
  "=": "&#x3D;",
  "/": "&#x2F;",
};

export function esc(value) {
  return String(value ?? "").replace(/[&<>"'`=\/]/g, (c) => escapeMap[c]);
}

export async function apiGet(url) {
  let res;
  try {
    res = await fetch(url, { credentials: "same-origin" });
  } catch {
    showToast("Network error. Please try again.", true);
    throw new Error("Network error.");
  }
  if (res.status === 401) {
    location.assign("/login");
    throw new Error("Unauthorized.");
  }
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    const message = body.message || `Request failed (${res.status}).`;
    showToast(message, true);
    throw new Error(message);
  }
  return body.data;
}

export function renderRows(tbody, rows, { columns, empty, colspan }) {
  if (!rows.length) {
    tbody.innerHTML = `<tr><td colspan="${colspan}" class="empty-row">${esc(empty)}</td></tr>`;
    return;
  }
  tbody.innerHTML = rows
    .map((row) => `<tr>${columns.map((key) => `<td>${esc(row[key])}</td>`).join("")}</tr>`)
    .join("");
}

export function renderOptions(select, items, { value, label, placeholder }) {
  select.innerHTML =
    `<option value="" disabled selected>${esc(placeholder)}</option>` +
    items.map((item) => `<option value="${esc(value(item))}">${esc(label(item))}</option>`).join("");
}

export function updateCount(el, tbody, label) {
  el.textContent = tbody.querySelector(".empty-row") ? "" : `${tbody.rows.length} ${label}`;
}

export async function initPage(options = {}) {
  const session = await getSession();
  applyNav(session);

  if (!guardPage(session, options)) return;

  popFlash();
  return session;
}
