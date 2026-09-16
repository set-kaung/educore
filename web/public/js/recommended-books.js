import { initPage, apiGet, showToast, esc, updateCount } from "./app.js";

await initPage({ require: ["professor", "admin"] });

const courseSelect = document.getElementById("course");
const bookSection = document.getElementById("book-section");
const bookRows = document.getElementById("book-rows");
const bookCount = document.getElementById("book-count");
const searchForm = document.getElementById("search-form");
const searchInput = document.getElementById("search");
const searchError = document.getElementById("search-error");
const searchResults = document.getElementById("search-results");

let courseOfferingId = null;

const coverUrl = (coverI, size) =>
    coverI ? `https://covers.openlibrary.org/b/id/${coverI}-${size}.jpg` : null;

async function loadCourses() {
    try {
        const courses = await apiGet("api/my/taught-courses");
        if (!courses.length) {
            courseSelect.innerHTML = '<option value="">You are not teaching any courses yet.</option>';
            return;
        }
        courseSelect.innerHTML = "";
        for (const course of courses) {
            courseSelect.append(new Option(
                `${course.name} (${course.course_code}) · ${course.section} · ${course.semester}`,
                course.semester_course_id,
            ));
        }
        await selectCourse(courses[0].semester_course_id);
    } catch {
        courseSelect.innerHTML = '<option value="">Failed to load your courses.</option>';
    }
}

async function selectCourse(id) {
    courseOfferingId = Number(id);
    bookSection.hidden = false;
    searchResults.innerHTML = "";
    searchError.hidden = true;
    await loadBooks();
}

async function loadBooks() {
    try {
        const books = await apiGet(`api/semester-courses/${courseOfferingId}/books`);
        if (!books.length) {
            bookRows.innerHTML = `<tr><td colspan="4" class="empty-row">No recommended books yet.</td></tr>`;
        } else {
            bookRows.innerHTML = books.map((book) => {
                const cover = coverUrl(book.cover_i, "S");
                const image = cover
                    ? `<img class="book-cover" src="${esc(cover)}" alt="" width="30">`
                    : "";
                return `<tr>
                    <td>${image} <span>${esc(book.title)}</span></td>
                    <td>${esc(book.author)}</td>
                    <td>${esc(book.isbn)}</td>
                    <td><button type="button" class="btn btn-link" data-remove="${book.id}">Remove</button></td>
                </tr>`;
            }).join("");
        }
        updateCount(bookCount, bookRows, "books");
    } catch {}
}

function bookCard(doc) {
    const cover = coverUrl(doc.cover_i, "M");
    const image = cover
        ? `<img class="book-cover" src="${esc(cover)}" alt="" loading="lazy">`
        : '<div class="book-cover-placeholder">No cover</div>';
    const authors = doc.author_name?.length ? doc.author_name.join(", ") : "Unknown author";
    const isbn = doc.isbn?.length ? doc.isbn[0] : "";
    const year = doc.first_publish_year ? ` · ${doc.first_publish_year}` : "";
    return `<article class="book-card">
        ${image}
        <div class="book-card-body">
            <h3 class="book-card-title">${esc(doc.title)}</h3>
            <p class="book-card-meta">${esc(authors)}${year}</p>
            <p class="book-card-meta">${esc(isbn)}</p>
        </div>
        <button type="button" class="btn btn-primary"
                data-add='${esc(JSON.stringify({
                    title: doc.title,
                    author: authors,
                    isbn,
                    cover_i: doc.cover_i,
                    open_library_key: doc.key,
                }))}'>Add</button>
    </article>`;
}

async function runSearch(query) {
    searchError.hidden = true;
    searchResults.innerHTML = '<p class="text-muted">Searching OpenLibrary…</p>';
    try {
        const data = await apiGet(`api/textbooks?q=${encodeURIComponent(query)}`);
        const docs = (data.docs || []).slice(0, 12);
        if (!docs.length) {
            searchResults.innerHTML = '<p class="text-muted">No books found. Try a different search.</p>';
            return;
        }
        searchResults.innerHTML = docs.map(bookCard).join("");
    } catch {
        searchResults.innerHTML = "";
        searchError.textContent = "Search failed. Please try again.";
        searchError.hidden = false;
    }
}

async function addBook(payload) {
    try {
        const res = await fetch(`api/semester-courses/${courseOfferingId}/books`, {
            method: "POST",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok) throw new Error(body.message || "Could not add the book.");
        showToast(`"${payload.title}" added to the recommended reading list.`);
        await loadBooks();
    } catch (err) {
        showToast(err.message, true);
    }
}

async function removeBook(id) {
    try {
        const res = await fetch(`api/books/${id}`, {
            method: "DELETE",
            credentials: "same-origin",
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok) throw new Error(body.message || "Could not remove the book.");
        showToast("Book removed.");
        await loadBooks();
    } catch (err) {
        showToast(err.message, true);
    }
}

courseSelect.addEventListener("change", () => selectCourse(courseSelect.value));

searchForm.addEventListener("submit", (event) => {
    event.preventDefault();
    const query = searchInput.value.trim();
    if (!query) return;
    runSearch(query);
});

bookRows.addEventListener("click", (event) => {
    const btn = event.target.closest("[data-remove]");
    if (btn && confirm("Remove this book from the recommended reading list?")) {
        removeBook(btn.dataset.remove);
    }
});

searchResults.addEventListener("click", (event) => {
    const btn = event.target.closest("[data-add]");
    if (!btn) return;
    try {
        addBook(JSON.parse(btn.dataset.add));
    } catch {
        showToast("Could not read book data.", true);
    }
});

loadCourses();