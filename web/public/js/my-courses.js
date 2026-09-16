import { initPage, apiGet, renderRows, updateCount } from "./app.js";

await initPage({});

const tbody = document.getElementById("course-rows");
const count = document.getElementById("course-count");

try {
    const data = await apiGet("api/my/courses");
    renderRows(tbody, data, {
        columns: ["name", "course_code", "section", "semester", "professor_name", "schedule"],
        empty: "You are not enrolled in any courses.",
        colspan: 6,
    });
    updateCount(count, tbody, "enrolled");
} catch {}
