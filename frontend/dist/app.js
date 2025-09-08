// API базовый URL
const API_BASE = '/api';

// Состояние приложения
let currentBoard = null;

// Инициализация приложения
document.addEventListener('DOMContentLoaded', function() {
    checkAuth();
});

// Проверка авторизации
async function checkAuth() {
    try {
        const response = await fetch(`${API_BASE}/board`, {
            credentials: 'include'
        });

        if (response.ok) {
            const board = await response.json();
            currentBoard = board;
            showApp();
            renderBoard();
        } else {
            showAuth();
        }
    } catch (error) {
        console.error('Ошибка проверки авторизации:', error);
        showAuth();
    }
}

// Показать форму авторизации
function showAuth() {
    document.getElementById('authSection').classList.remove('hidden');
    document.getElementById('app').classList.add('hidden');
}

// Показать приложение
function showApp() {
    document.getElementById('authSection').classList.add('hidden');
    document.getElementById('app').classList.remove('hidden');
    if (currentBoard) {
        document.getElementById('boardTitle').textContent = currentBoard.name;
    }
}

// Генерация пароля для новой доски
async function generatePassword() {
    const boardName = document.getElementById('boardName').value.trim();
    if (!boardName) {
        alert('Введите название доски');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/generate-password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ board_name: boardName }),
            credentials: 'include'
        });

        const data = await response.json();

        if (response.ok) {
            document.getElementById('passwordDisplay').textContent = data.password;
            document.getElementById('generatedPassword').classList.remove('hidden');
        } else {
            alert('Ошибка создания доски: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка генерации пароля:', error);
        alert('Ошибка сети');
    }
}

// Вход в систему
async function login() {
    const password = document.getElementById('password').value.trim();
    if (!password) {
        alert('Введите пароль');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ password: password }),
            credentials: 'include'
        });

        const data = await response.json();

        if (response.ok) {
            await checkAuth(); // Перезагружаем данные доски
        } else {
            alert('Ошибка входа: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка входа:', error);
        alert('Ошибка сети');
    }
}

// Выход из системы
async function logout() {
    try {
        await fetch(`${API_BASE}/logout`, {
            method: 'POST',
            credentials: 'include'
        });

        currentBoard = null;
        showAuth();
        // Очищаем форму
        document.getElementById('password').value = '';
        document.getElementById('boardName').value = '';
        document.getElementById('generatedPassword').classList.add('hidden');
    } catch (error) {
        console.error('Ошибка выхода:', error);
    }
}

// Отрисовка доски
function renderBoard() {
    const boardElement = document.getElementById('board');
    boardElement.innerHTML = '';

    // Отрисовка колонок
    currentBoard.columns.forEach(column => {
        const columnElement = createColumnElement(column);
        boardElement.appendChild(columnElement);
    });

    // Добавляем кнопку для создания новой колонки
    const addColumnElement = document.createElement('div');
    addColumnElement.className = 'add-column';
    addColumnElement.innerHTML = '<h3>+ Добавить колонку</h3>';
    addColumnElement.onclick = showAddColumnModal;
    boardElement.appendChild(addColumnElement);
}

// Создание элемента колонки
function createColumnElement(column) {
    const columnDiv = document.createElement('div');
    columnDiv.className = 'column';
    columnDiv.dataset.columnId = column.id;

    // Настройка drag and drop
    columnDiv.addEventListener('dragover', handleDragOver);
    columnDiv.addEventListener('drop', handleDrop);

    columnDiv.innerHTML = `
        <div class="column-header">
            <h3 class="column-title">${column.name}</h3>
            <div>
                <button class="btn btn-small" onclick="showAddCardModal('${column.id}')">+ Карточка</button>
                <button class="btn btn-danger btn-small" onclick="deleteColumn('${column.id}')">×</button>
            </div>
        </div>
        <div class="cards" id="cards-${column.id}">
            ${column.cards.map(card => createCardHTML(card)).join('')}
        </div>
    `;

    return columnDiv;
}

// Создание HTML для карточки
function createCardHTML(card) {
    return `
        <div class="card" draggable="true" data-card-id="${card.id}" 
             ondragstart="handleDragStart(event)">
            <div class="card-title">${card.title}</div>
            ${card.description ? `<div class="card-description">${card.description}</div>` : ''}
            ${card.assignee ? `<div class="card-assignee">${card.assignee}</div>` : ''}
            <div class="card-actions">
                <button class="btn btn-danger btn-small" onclick="deleteCard('${card.id}')">Удалить</button>
            </div>
        </div>
    `;
}

// Показать модальное окно добавления колонки
function showAddColumnModal() {
    document.getElementById('columnModal').style.display = 'block';
}

// Показать модальное окно добавления карточки
function showAddCardModal(columnId) {
    document.getElementById('cardColumnId').value = columnId;
    document.getElementById('cardModal').style.display = 'block';
}

// Закрыть модальное окно
function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
    // Очищаем формы
    if (modalId === 'columnModal') {
        document.getElementById('columnName').value = '';
    } else if (modalId === 'cardModal') {
        document.getElementById('cardTitle').value = '';
        document.getElementById('cardDescription').value = '';
        document.getElementById('cardAssignee').value = '';
    }
}

// Добавление колонки
async function addColumn() {
    const name = document.getElementById('columnName').value.trim();
    if (!name) {
        alert('Введите название колонки');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/columns`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: name }),
            credentials: 'include'
        });

        if (response.ok) {
            await refreshBoard();
            closeModal('columnModal');
        } else {
            const data = await response.json();
            alert('Ошибка добавления колонки: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка добавления колонки:', error);
        alert('Ошибка сети');
    }
}

// Удаление колонки
async function deleteColumn(columnId) {
    if (!confirm('Вы уверены, что хотите удалить эту колонку?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/columns/${columnId}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        if (response.ok) {
            await refreshBoard();
        } else {
            const data = await response.json();
            alert('Ошибка удаления колонки: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка удаления колонки:', error);
        alert('Ошибка сети');
    }
}

// Добавление карточки
async function addCard() {
    const title = document.getElementById('cardTitle').value.trim();
    const description = document.getElementById('cardDescription').value.trim();
    const assignee = document.getElementById('cardAssignee').value.trim();
    const columnId = document.getElementById('cardColumnId').value;

    if (!title) {
        alert('Введите заголовок карточки');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/cards`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title,
                description: description,
                assignee: assignee,
                column_id: columnId
            }),
            credentials: 'include'
        });

        if (response.ok) {
            await refreshBoard();
            closeModal('cardModal');
        } else {
            const data = await response.json();
            alert('Ошибка добавления карточки: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка добавления карточки:', error);
        alert('Ошибка сети');
    }
}

// Удаление карточки
async function deleteCard(cardId) {
    if (!confirm('Вы уверены, что хотите удалить эту карточку?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/cards/${cardId}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        if (response.ok) {
            await refreshBoard();
        } else {
            const data = await response.json();
            alert('Ошибка удаления карточки: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка удаления карточки:', error);
        alert('Ошибка сети');
    }
}

// Обновление данных доски
async function refreshBoard() {
    try {
        const response = await fetch(`${API_BASE}/board`, {
            credentials: 'include'
        });

        if (response.ok) {
            currentBoard = await response.json();
            renderBoard();
        }
    } catch (error) {
        console.error('Ошибка обновления доски:', error);
    }
}

// Drag and Drop функционал
let draggedCardId = null;

function handleDragStart(event) {
    draggedCardId = event.target.dataset.cardId;
    event.dataTransfer.effectAllowed = 'move';
}

function handleDragOver(event) {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    event.currentTarget.classList.add('drag-over');
}

function handleDragLeave(event) {
    event.currentTarget.classList.remove('drag-over');
}

async function handleDrop(event) {
    event.preventDefault();
    event.currentTarget.classList.remove('drag-over');

    const targetColumnId = event.currentTarget.dataset.columnId;

    if (draggedCardId && targetColumnId) {
        try {
            const response = await fetch(`${API_BASE}/cards/move`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    card_id: draggedCardId,
                    column_id: targetColumnId,
                    order: 0 // Можно улучшить для точного позиционирования
                }),
                credentials: 'include'
            });

            if (response.ok) {
                await refreshBoard();
            } else {
                const data = await response.json();
                alert('Ошибка перемещения карточки: ' + data.error);
            }
        } catch (error) {
            console.error('Ошибка перемещения карточки:', error);
            alert('Ошибка сети');
        }
    }

    draggedCardId = null;
}

// Закрытие модальных окон при клике вне их
window.onclick = function(event) {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
}
