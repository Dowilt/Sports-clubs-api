console.log("script.js loaded at", new Date().toLocaleTimeString());
// Функция для обработки нажатия Enter в поиске
function handleEnter(event) {
    if (event.key === "Enter") searchClubs();
}

// Дебаунс для фильтров и поиска
let debounceTimer;
function debounceSearch() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
        searchClubs();
    }, 500); // Задержка 500 мс
}

// Показать/скрыть панель фильтров
function toggleFilters() {
    const panel = document.getElementById("filtersPanel");
    panel.classList.toggle("active");
    panel.style.display = panel.style.display === "block" ? "none" : "block";
}

// Сбросить фильтры
function resetFilters() {
    document.getElementById("titlesMin").value = "";
    document.getElementById("titlesMax").value = "";
    document.getElementById("ageMin").value = "";
    document.getElementById("ageMax").value = "";
    searchClubs();
}

// Переключение темы
function toggleDarkMode() {
    document.body.classList.toggle("dark-mode");
}

// Получить список клубов
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

    url = url.slice(0, -1); // убрать последний &

    try {
        const response = await fetch(url);
        const data = await response.json();
        const tbody = document.querySelector("#clubsTable tbody");
        tbody.innerHTML = "";

        if (!data || data.length === 0) {
            tbody.innerHTML = "<tr><td colspan='6' class='no-results'>Клубы не найдены</td></tr>";
            return;
        }

        data.forEach(club => {
            const row = document.createElement("tr");

            row.innerHTML = `
                <td>${club.id}</td>
                <td contenteditable="true" onblur="updateCell(this, ${club.id}, 'name')">${club.name}</td>
                <td contenteditable="true" onblur="updateCell(this, ${club.id}, 'city')">${club.city}</td>
                <td class="center-cell" contenteditable="true" onblur="updateCell(this, ${club.id}, 'titles')">${club.titles ?? 0}</td>
                <td class="center-cell" contenteditable="true" onblur="updateCell(this, ${club.id}, 'avgAge')">${club.avgAge}</td>
                <td class="delete-cell" onclick="deleteClub(${club.id})">×</td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error("Ошибка загрузки данных:", error);
        alert("Ошибка подключения к серверу.");
    }
}

// Переключение панели добавления
function toggleAddPanel() {
    const panel = document.getElementById("addPanel");
    panel.classList.toggle("active");
    panel.style.display = panel.style.display === "block" ? "none" : "block";
}

// Добавление нового клуба
async function addNewClub() {
    const name = document.getElementById("newClubName").value.trim();
    const city = document.getElementById("newClubCity").value.trim();
    const titles = parseInt(document.getElementById("newClubTitles").value);
    const avgAge = parseInt(document.getElementById("newClubAvgAge").value);

    if (!name || !city || isNaN(titles) || isNaN(avgAge)) {
        alert("Пожалуйста, заполните все поля корректно.");
        return;
    }

    const newClub = {
        name,
        city,
        titles,
        avgAge
    };

    try {
        const res = await fetch("http://localhost:8080/clubs", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(newClub)
        });

        if (!res.ok) throw new Error(await res.text());

        alert("Клуб успешно добавлен");
        document.getElementById("newClubName").value = "";
        document.getElementById("newClubCity").value = "";
        document.getElementById("newClubTitles").value = "";
        document.getElementById("newClubAvgAge").value = "";
        searchClubs(); // Обновляем таблицу
    } catch (err) {
        console.error(err);
        alert("Ошибка при добавлении клуба: " + err.message);
    }
}

// Редактирование ячейки
async function updateCell(cell, clubId, field) {
    const value = cell.innerText.trim();
    const payload = {};

    if (field === "titles") {
        payload["titles_count"] = parseInt(value);
    } else if (field === "avgAge") {
        payload["avg_age"] = parseInt(value);
    } else {
        payload[field] = value;
    }

    try {
        const res = await fetch(`http://localhost:8080/clubs/${clubId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) throw new Error(await res.text());
    } catch (err) {
        console.error(err);
        alert("Ошибка при обновлении данных: " + err.message);
        searchClubs(); // Откат к актуальным данным
    }
}

// Удаление клуба
async function deleteClub(id) {
    if (!confirm("Вы уверены, что хотите удалить этот клуб?")) return;

    const res = await fetch(`http://localhost:8080/clubs/${id}`, {
        method: "DELETE"
    });

    if (!res.ok) {
        alert("Ошибка удаления");
        return;
    }

    searchClubs(); // Обновляем таблицу
}

// Автозагрузка данных при открытии страницы
window.onload = () => {
    searchClubs();
};
