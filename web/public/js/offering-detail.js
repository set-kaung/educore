import { initPage, apiGet, esc } from "./app.js";

const params = new URLSearchParams(location.search);
const id = params.get("id");

const titleEl = document.getElementById("offering-title");
const metaEl = document.getElementById("offering-meta");
const bookList = document.getElementById("book-list");
const bookCount = document.getElementById("book-count");

const coverUrl = (coverI, size) =>
    coverI ? `https://covers.openlibrary.org/b/id/${coverI}-${size}.jpg` : null;

if (!id) {
    location.replace("course-offerings");
} else {
    await initPage({});
    await loadOffering(id);
}

async function loadOffering(offeringId) {
    try {
        const offering = await apiGet(`api/semester-courses/${offeringId}`);
        document.title = `${offering.name} · EduCore`;
        titleEl.textContent = `${offering.name} (${offering.course_code})`;
        metaEl.textContent = `Section ${offering.section} · ${offering.semester} · ${offering.professor_name} · ${offering.schedule}`;
    } catch {
        titleEl.textContent = "Course offering not found.";
        metaEl.textContent = "";
        return;
    }

    try {
        const books = await apiGet(`api/semester-courses/${offeringId}/books`);
        console.log("Recommended books response:", books);
        if (!books.length) {
            bookList.innerHTML = '<p class="empty-row">No recommended books for this course yet.</p>';
        } else {
            bookCount.textContent = `${books.length} ${books.length === 1 ? "book" : "books"}`;
            bookList.innerHTML = books.map((book) => {
                const cover = coverUrl(book.cover_i, "M");
                const image = cover
                    ? `<img class="book-cover" src="${esc(cover)}" alt="" loading="lazy" width="60" height="90">`
                    : '<div class="book-cover-placeholder">No cover</div>';
                return `<article class="book-card">
                    ${image}
                    <div class="book-card-body">
                        <h3 class="book-card-title">${esc(book.title)}</h3>
                        <p class="book-card-meta">${esc(book.author || "Unknown author")}</p>
                        ${book.isbn ? `<p class="book-card-meta">ISBN: ${esc(book.isbn)}</p>` : ""}
                    </div>
                </article>`;
            }).join("");
        }
    } catch (err) {
        console.error("Failed to load recommended books:", err);
        bookList.innerHTML = `<p class="empty-row">Failed to load recommended books: ${esc(err.message)}</p>`;
    }
}
