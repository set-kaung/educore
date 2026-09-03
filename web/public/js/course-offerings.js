import { initPage, apiGet, showToast, esc, updateCount } from "/js/app.js";

const session = await initPage({});
const isStudent = session?.role === "student";

const select = document.getElementById("semester");
const tbody = document.getElementById("offering-rows");
const count = document.getElementById("offering-count");
const enrollCol = document.getElementById("enroll-col");

if (isStudent) enrollCol.hidden = false;

const enrolledIds = new Set();

async function loadEnrolled() {
    if (!isStudent) return;
    try {
        const data = await apiGet("/api/my/courses");
        for (const row of data) {
            enrolledIds.add(row.semester_course_id);
        }
    } catch {}
}

function rowHtml(row) {
    const action = enrolledIds.has(row.semester_course_id)
        ? '<span class="enrolled-tag">Enrolled</span>'
        : `<button type="button" class="btn btn-primary" data-enroll="${row.semester_course_id}">Enroll</button>`;
    return `<tr>
        <td>${esc(row.name)}</td>
        <td>${esc(row.course_code)}</td>
        <td>${esc(row.section)}</td>
        <td>${esc(row.professor_name)}</td>
        <td>${esc(row.schedule)}</td>
        ${isStudent ? `<td>${action}</td>` : ""}
    </tr>`;
}

async function loadOfferings() {
    try {
        const data = await apiGet(`/api/semester-courses?semester_id=${encodeURIComponent(select.value)}`);
        if (!data.length) {
            tbody.innerHTML = `<tr><td colspan="6" class="empty-row">No offerings for this semester yet.</td></tr>`;
        } else {
            tbody.innerHTML = data.map(rowHtml).join("");
        }
        updateCount(count, tbody, "offerings");
    } catch {}
}

async function loadSemesters() {
    try {
        const semesters = await apiGet("/api/semesters");
        if (semesters.length) {
            select.innerHTML = "";
            let current = semesters[0];
            for (const semester of semesters) {
                select.append(new Option(semester.name, semester.id));
                if (semester.is_current) current = semester;
            }
            select.value = current.id;
            await loadOfferings();
        } else {
            select.innerHTML = '<option value="">No semesters yet</option>';
        }
    } catch {
        select.innerHTML = '<option value="">Failed to load semesters</option>';
    }
}

tbody.addEventListener("click", async (event) => {
    const btn = event.target.closest("[data-enroll]");
    if (!btn) return;
    btn.disabled = true;

    try {
        const res = await fetch(`/api/semester-courses/${btn.dataset.enroll}/enroll`, {
            method: "POST",
            credentials: "same-origin",
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok) {
            showToast(body.message || "Enrollment failed.", true);
            btn.disabled = false;
            return;
        }
        enrolledIds.add(Number(btn.dataset.enroll));
        showToast(`Enrolled in ${btn.closest("tr").cells[0].textContent}.`);
        await loadOfferings();
    } catch {
        showToast("Network error. Please try again.", true);
        btn.disabled = false;
    }
});

select.addEventListener("change", loadOfferings);
await loadEnrolled();
loadSemesters();
