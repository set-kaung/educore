import { initPage, flash, apiGet, esc, showToast, updateCount } from "./app.js";

await initPage({ require: ["professor", "admin"] });

const form = document.getElementById("course-form");
const errorBox = document.getElementById("course-error");
const tbody = document.getElementById("course-rows");
const count = document.getElementById("course-count");

async function loadCourses() {
    try {
        const data = await apiGet("api/courses");
        if (!data.length) {
            tbody.innerHTML = `<tr><td colspan="3" class="empty-row">No courses yet.</td></tr>`;
        } else {
            tbody.innerHTML = data.map((course) => `<tr>
                <td>${esc(course.name)}</td>
                <td>${esc(course.course_code)}</td>
                <td><button type="button" class="btn btn-danger btn-sm" data-delete-course="${course.id}">Delete</button></td>
            </tr>`).join("");
        }
        updateCount(count, tbody, "courses");
    } catch {}
}

tbody.addEventListener("click", async (event) => {
    const btn = event.target.closest("[data-delete-course]");
    if (!btn) return;

    const name = btn.closest("tr").cells[0].textContent;
    if (!confirm(`Delete course "${name}"? This cannot be undone.`)) return;

    btn.disabled = true;
    try {
        const res = await fetch(`api/courses/${btn.dataset.deleteCourse}`, {
            method: "DELETE",
            credentials: "same-origin",
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok) {
            showToast(body.message || "Could not delete the course.", true);
            btn.disabled = false;
            return;
        }
        showToast(`Course "${name}" deleted.`);
        loadCourses();
    } catch {
        showToast("Network error. Please try again.", true);
        btn.disabled = false;
    }
});

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    errorBox.hidden = true;

    try {
        const res = await fetch("api/courses", {
            method: "POST",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                name: form.name.value.trim(),
                course_code: form.course_code.value.trim(),
            }),
        });
        const body = await res.json();
        if (!res.ok) throw new Error(body.message || "Could not create the course.");

        flash(`Course "${body.data.name}" created successfully.`);
        location.reload();
    } catch (err) {
        errorBox.textContent = err.message;
        errorBox.hidden = false;
    }
});

loadCourses();
