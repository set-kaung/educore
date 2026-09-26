import { initPage, apiGet, showToast, esc, updateCount } from "./app.js";

const session = await initPage({});
const isStudent = session?.role === "student";
const isStaff = session?.role === "professor" || session?.role === "admin";

const select = document.getElementById("semester");
const tbody = document.getElementById("offering-rows");
const count = document.getElementById("offering-count");
const enrollCol = document.getElementById("enroll-col");
const actionCol = document.getElementById("action-col");

if (isStudent) enrollCol.hidden = false;
if (isStaff) actionCol.hidden = false;

const enrolledIds = new Set();

async function loadEnrolled() {
    if (!isStudent) return;
    try {
        const data = await apiGet("api/my/courses");
        for (const row of data) {
            enrolledIds.add(row.semester_course_id);
        }
    } catch {}
}

function rowHtml(row) {
    const action = enrolledIds.has(row.semester_course_id)
        ? '<span class="enrolled-tag">Enrolled</span>'
        : `<button type="button" class="btn btn-primary" data-enroll="${row.semester_course_id}">Enroll</button>`;
    const actions = isStaff
        ? `<button type="button" class="btn btn-danger btn-sm" data-delete-offering="${row.semester_course_id}">Delete</button>`
        : "";
    const lastCol = isStudent ? `<td>${action}</td>` : (isStaff ? `<td>${actions}</td>` : "");
    return `<tr>
        <td><a class="nav-link" href="offering-detail?id=${row.semester_course_id}">${esc(row.name)}</a></td>
        <td>${esc(row.course_code)}</td>
        <td>${esc(row.section)}</td>
        <td>${esc(row.professor_name)}</td>
        <td>${esc(row.schedule)}</td>
        ${lastCol}
    </tr>`;
}

async function loadOfferings() {
    try {
        const data = await apiGet(`api/semester-courses?semester_id=${encodeURIComponent(select.value)}`);
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
        const semesters = await apiGet("api/semesters");
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
    const enrollBtn = event.target.closest("[data-enroll]");
    if (enrollBtn) {
        enrollBtn.disabled = true;

        try {
            const res = await fetch(`api/semester-courses/${enrollBtn.dataset.enroll}/enroll`, {
                method: "POST",
                credentials: "same-origin",
            });
            const body = await res.json().catch(() => ({}));
            if (!res.ok) {
                showToast(body.message || "Enrollment failed.", true);
                enrollBtn.disabled = false;
                return;
            }
            enrolledIds.add(Number(enrollBtn.dataset.enroll));
            showToast(`Enrolled in ${enrollBtn.closest("tr").cells[0].textContent}.`);
            await loadOfferings();
        } catch {
            showToast("Network error. Please try again.", true);
            enrollBtn.disabled = false;
        }
        return;
    }

    const deleteBtn = event.target.closest("[data-delete-offering]");
    if (!deleteBtn) return;

    const name = deleteBtn.closest("tr").cells[0].textContent;
    if (!confirm(`Delete offering "${name}"? This cannot be undone.`)) return;

    deleteBtn.disabled = true;
    try {
        const res = await fetch(`api/semester-courses/${deleteBtn.dataset.deleteOffering}`, {
            method: "DELETE",
            credentials: "same-origin",
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok) {
            showToast(body.message || "Could not delete the offering.", true);
            deleteBtn.disabled = false;
            return;
        }
        showToast(`Offering "${name}" deleted.`);
        loadOfferings();
    } catch {
        showToast("Network error. Please try again.", true);
        deleteBtn.disabled = false;
    }
});

select.addEventListener("change", loadOfferings);
await loadEnrolled();
loadSemesters();
