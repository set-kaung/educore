// Shared glue between the static HTML pages and the JSON API.
// All business logic lives behind /api/*; this file only handles
// sessions, navigation, redirects and error toasts.

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

export async function login(username, password) {
  const res = await fetch("/api/login", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(body.message || "Login failed.");
  }
  sessionCache.value = { user_id: body.data.user_id, role: body.data.role };
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

function showToast(message, isError = false) {
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

function wireHtmxErrors() {
  document.body.addEventListener("htmx:responseError", (evt) => {
    const status = evt.detail.xhr.status;
    if (status === 401) {
      location.assign("/login");
      return;
    }
    let message = `Request failed (${status}).`;
    try {
      message = JSON.parse(evt.detail.xhr.responseText).message || message;
    } catch {}
    showToast(message, true);
  });
  document.body.addEventListener("htmx:sendError", () => {
    showToast("Network error. Please try again.", true);
  });
}

export async function initPage(options = {}) {
  wireHtmxErrors();

  const session = await getSession();
  applyNav(session);

  if (!guardPage(session, options)) return;

  popFlash();
  return session;
}
