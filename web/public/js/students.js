import { initPage, apiGet, renderRows, updateCount } from "./app.js";

await initPage({ require: ["professor", "admin"] });

const tbody = document.getElementById("student-rows");
const count = document.getElementById("student-count");
const search = document.querySelector('input[name="q"]');

async function loadStudents() {
    try {
        const q = search.value.trim();
        const data = await apiGet(q ? `api/students?q=${encodeURIComponent(q)}` : "api/students");
        const rows = Array.isArray(data) ? data : [];
        renderRows(tbody, rows, {
            columns: ["name", "student_id", "department_name"],
            empty: rows.length ? "No students found." : "No students yet.",
            colspan: 3,
        });
        updateCount(count, tbody, "shown");
    } catch {}
}

let debounce;
search.addEventListener("input", () => {
    clearTimeout(debounce);
    debounce = setTimeout(loadStudents, 300);
});

loadStudents();
