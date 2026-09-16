import { initPage, flash, apiGet, renderOptions } from "./app.js";

await initPage({ require: ["professor", "admin"] });

const form = document.getElementById("offering-form");
const errorBox = document.getElementById("offering-error");
const scheduleRows = document.getElementById("schedule-rows");
const rowTemplate = document.getElementById("tpl-schedule-row");
const courseSelect = document.getElementById("course_id");
const professorSelect = document.getElementById("taught_by");
const semesterSelect = document.getElementById("semester_id");

async function loadSelects() {
    try {
        const semesters = await apiGet("api/semesters");
        renderOptions(semesterSelect, semesters, {
            value: (semester) => semester.id,
            label: (semester) => semester.name,
            placeholder: "Select a semester…",
        });
    } catch {}
    try {
        const courses = await apiGet("api/courses");
        renderOptions(courseSelect, courses, {
            value: (course) => course.id,
            label: (course) => `${course.name} (${course.course_code})`,
            placeholder: "Select a course…",
        });
    } catch {}
    try {
        const professors = await apiGet("api/professors");
        renderOptions(professorSelect, professors, {
            value: (professor) => professor.id,
            label: (professor) => professor.name,
            placeholder: "Select a professor…",
        });
    } catch {}
}

function addScheduleRow() {
    scheduleRows.append(rowTemplate.content.cloneNode(true));
}

scheduleRows.addEventListener("click", (event) => {
    if (event.target.closest("[data-remove-row]")) {
        event.target.closest("[data-schedule-row]").remove();
    }
});

document.getElementById("add-schedule").addEventListener("click", () => {
    if (!scheduleRows.children.length) {
        addScheduleRow();
        return;
    }
    scheduleRows.append(
        scheduleRows.lastElementChild.cloneNode(true)
    );
    for (const input of scheduleRows.lastElementChild.querySelectorAll("input, select")) {
        input.value = "";
    }
});

addScheduleRow();

loadSelects();

function collectSchedules() {
    return Array.from(form.querySelectorAll("[data-schedule-row]"), (row) => ({
        weekday: row.querySelector("[name=schedule_weekday]").value,
        start_time: row.querySelector("[name=schedule_start_time]").value,
        end_time: row.querySelector("[name=schedule_end_time]").value,
    }));
}

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    errorBox.hidden = true;

    try {
        const res = await fetch("api/semester-courses", {
            method: "POST",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                semester_id: Number(form.semester_id.value),
                course_id: Number(form.course_id.value),
                section: form.section.value.trim(),
                taught_by: Number(form.taught_by.value),
                schedules: collectSchedules(),
            }),
        });
        const body = await res.json();
        if (!res.ok) throw new Error(body.message || "Could not create the course offering.");

        flash(`Course offering created for ${body.data.semester}.`);
        location.assign("course-offerings");
    } catch (err) {
        errorBox.textContent = err.message;
        errorBox.hidden = false;
    }
});
