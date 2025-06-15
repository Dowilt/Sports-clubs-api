function handleEnter(event) {
    if (event.key === "Enter") searchClubs();
}

let debounceTimer;

function debounceSearch() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
        searchClubs();
    }, 500); // Задержка 
}

function toggleFilters() {
    const panel = document.getElementById("filtersPanel");
    panel.classList.toggle("active");
}

function resetFilters() {
    document.getElementById("titlesMin").value = "";
    document.getElementById("titlesMax").value = "";
    document.getElementById("ageMin").value = "";
    document.getElementById("ageMax").value = "";
    searchClubs();
}

function toggleDarkMode() {
    document.body.classList.toggle("dark-mode");
}

async function searchClubs() {
    const query = document.getElementById("searchQuery").value.trim();
    const titlesMin = document.getElementById("titlesMin").value.trim();
    const titlesMax = document.getElementById("titlesMax").value.trim();
    const ageMin = document.getElementById("ageMin").value.trim();
    const ageMax = document.getElementById("ageMax").value.trim();

    const sortField = document.getElementById("sortField").value;
    const sortOrder = document.getElementById("sortOrder").value;

    let url = "http://localhost:8080/clubs?";

    if (query) url += `q=${encodeURIComponent(query)}&`;
    if (titlesMin !== "") url += `titles_min=${encodeURIComponent(titlesMin)}&`;
    if (titlesMax !== "") url += `titles_max=${encodeURIComponent(titlesMax)}&`;
    if (ageMin !== "") url += `age_min=${encodeURIComponent(ageMin)}&`;
    if (ageMax !== "") url += `age_max=${encodeURIComponent(ageMax)}&`;
    if (sortField && sortOrder) url += `sort_by=${encodeURIComponent(sortField)}&sort_order=${encodeURIComponent(sortOrder)}&`;

    url = url.slice(0, -1);

    try {
        const response = await fetch(url);
        const data = await response.json();

        const tbody = document.querySelector("#clubsTable tbody");
        tbody.innerHTML = "";

        if (!data || data.length === 0) {
            tbody.innerHTML = "<tr><td colspan='5' class='no-results'>Клубы не найдены</td></tr>";
            return;
        }

        data.forEach(club => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${club.id}</td>
                <td>${club.name}</td>
                <td>${club.city}</td>
                <td>${club.titles_count ?? 0}</td>
                <td>${club.avg_age}</td>
            `;
            tbody.appendChild(row);
        });

    } catch (error) {
        console.error("Ошибка загрузки данных:", error);
        alert("Ошибка подключения к серверу.");
    }
}

window.onload = () => {
    searchClubs();
};