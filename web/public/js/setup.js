import { apiGet } from "./app.js";

const form = document.getElementById("setup-form");
const errorBox = document.getElementById("setup-error");
const emailLine = document.getElementById("setup-email");
const departmentSelect = document.getElementById("department_id");
const studentIdInput = document.getElementById("student_id");

let context = { email: "", name: "" };
try {
    context = await apiGet("api/setup/context");
} catch {}

if (context.email) {
    emailLine.textContent = `Signed in as ${context.email}`;
}
if (context.name) {
    form.name.value = context.name;
}

const localPart = (context.email || "").split("@")[0];
if (localPart && /^\d+$/.test(localPart)) {
    studentIdInput.value = localPart;
}

try {
    const departments = await apiGet("api/departments");
    if (departments.length) {
        departmentSelect.innerHTML = "";
        for (const department of departments) {
            departmentSelect.append(new Option(department.name, department.id));
        }
    } else {
        departmentSelect.innerHTML = '<option value="" disabled selected>No departments available</option>';
    }
} catch {
    departmentSelect.innerHTML = '<option value="" disabled selected>Failed to load departments</option>';
}

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    errorBox.hidden = true;

    try {
        const res = await fetch("api/setup", {
            method: "POST",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                name: form.name.value.trim(),
                department_id: Number(form.department_id.value),
                student_id: studentIdInput.value.trim(),
            }),
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok) {
            throw new Error(body.message || "Setup failed.");
        }
        location.assign("./");
    } catch (err) {
        errorBox.textContent = err.message;
        errorBox.hidden = false;
    }
});
