import { initPage, flash, apiGet, renderRows, updateCount } from "/js/app.js";

await initPage({ require: ["professor", "admin"] });

const form = document.getElementById("course-form");
const errorBox = document.getElementById("course-error");
const tbody = document.getElementById("course-rows");
const count = document.getElementById("course-count");

async function loadCourses() {
    try {
        const data = await apiGet("/api/courses");
        renderRows(tbody, data, {
            columns: ["name", "course_code"],
            empty: "No courses yet.",
            colspan: 2,
        });
        updateCount(count, tbody, "courses");
    } catch {}
}

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    errorBox.hidden = true;

    try {
        const res = await fetch("/api/courses", {
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
