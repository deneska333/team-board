// Глобальные переменные
let currentBoardId = null;
let currentColumnId = null;
let editingCardId = null;
let draggedCard = null;

// API базовый URL
const API_BASE = '/api';

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    setupEventListeners();
    checkAuthStatus();
});

// Настройка обработчиков событий
function setupEventListeners() {
    // Форма создания доски
    document.getElementById('create-form').addEventListener('submit', createBoard);

    // Форма входа
    document.getElementById('login-form-element').addEventListener('submit', loginToBoard);

    // Форма карт��чки
    document.getElementById('card-form').addEventListener('submit', saveCard);

    // Закрытие модального окна по клику вне его
    document.getElementById('card-modal').addEventListener('click', function(e) {
        if (e.target === this) {
            closeCardModal();
        }
    });
}

// Проверка статуса авторизации
async function checkAuthStatus() {
    try {
        const response = await fetch(`${API_BASE}/board`, {
            credentials: 'include'
        });

        if (response.ok) {
            const board = await response.json();
            showBoard(board);
        }
    } catch (error) {
        console.log('Пользователь не авторизован');
    }
}

// Создание новой доски
async function createBoard(e) {
    e.preventDefault();

    const name = document.getElementById('board-name').value;
    const password = document.getElementById('board-password').value;

    const errorDiv = document.getElementById('create-error');
    errorDiv.classList.add('hidden');

    try {
        const response = await fetch(`${API_BASE}/boards`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ name, password }),
        });

        if (response.ok) {
            const board = await response.json();
            showBoard(board);
        } else {
            const error = await response.json();
            showError('create-error', error.error);
        }
    } catch (error) {
        showError('create-error', 'Ошибка соединения с сервером');
    }
}

// Вход в существующую доску
async function loginToBoard(e) {
    e.preventDefault();

    const boardId = document.getElementById('board-id').value;
    const password = document.getElementById('login-password').value;

    const errorDiv = document.getElementById('login-error');
    errorDiv.classList.add('hidden');

    try {
        const response = await fetch(`${API_BASE}/boards/${boardId}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ password }),
        });

        if (response.ok) {
            // Получаем данные доски после успешного входа
            const boardResponse = await fetch(`${API_BASE}/board`, {
                credentials: 'include'
            });

            if (boardResponse.ok) {
                const board = await boardResponse.json();
                showBoard(board);
            }
        } else {
            const error = await response.json();
            showError('login-error', error.error);
        }
    } catch (error) {
        showError('login-error', 'Ошибка соединения с сервером');
    }
}

// Отображение доски
function showBoard(board) {
    currentBoardId = board.id;

    // Скрываем формы входа и показываем доску
    document.getElementById('create-board-form').classList.add('hidden');
    document.getElementById('login-form').classList.add('hidden');
    document.getElementById('board-container').classList.remove('hidden');

    // Обновляем заголовок и ID доски
    document.getElementById('board-title').textContent = board.name;
    document.getElementById('board-id-display').textContent = `ID: ${board.id}`;

    // Отображаем колонки
    renderBoard(board);
}

// Отрисовка доски с колонками
function renderBoard(board) {
    const boardElement = document.getElementById('board');
    boardElement.innerHTML = '';

    // Создаем колонки из данных доски
    board.columns.forEach(column => {
        const columnElement = createColumnElement(column);
        boardElement.appendChild(columnElement);
    });

    // Добавляем кнопку создания новой колонки
    const addColumnButton = document.createElement('div');
    addColumnButton.className = 'add-column-container';
    addColumnButton.innerHTML = `
        <button class="btn add-column-btn" onclick="openColumnModal()">
            + Добавить колонку
        </button>
    `;
    boardElement.appendChild(addColumnButton);
}

// Создание элемента колонки
function createColumnElement(column) {
    const columnDiv = document.createElement('div');
    columnDiv.className = 'column';
    columnDiv.dataset.columnId = column.id;

    columnDiv.innerHTML = `
        <div class="column-header">
            <h3 class="column-title">${column.name}</h3>
            <div class="column-actions">
                <button class="btn add-card-btn" onclick="openCardModal('${column.id}')">+ Карточка</button>
                <button class="btn-icon" onclick="editColumn('${column.id}', '${column.name}')" title="Редактировать колонку">✎</button>
                <button class="btn-icon btn-danger-icon" onclick="deleteColumn('${column.id}')" title="Удалить колонку">×</button>
            </div>
        </div>
        <div class="cards" id="cards-${column.id}">
            ${column.cards.map(card => createCardHTML(card)).join('')}
        </div>
        <div class="drop-zone" ondrop="dropCard(event, '${column.id}')" ondragover="allowDrop(event)">
            Перетащите карточку сюда
        </div>
    `;

    return columnDiv;
}

// Создание HTML карточки
function createCardHTML(card) {
    return `
        <div class="card" draggable="true" data-card-id="${card.id}" 
             ondragstart="dragStart(event)" ondragend="dragEnd(event)">
            <div class="card-actions">
                <button onclick="editCard('${card.id}')" title="Редактировать">✎</button>
                <button onclick="deleteCard('${card.id}')" title="Удалить">×</button>
            </div>
            <div class="card-title">${card.title}</div>
            ${card.description ? `<div class="card-description">${card.description}</div>` : ''}
            ${card.assignee ? `<div class="card-assignee">${card.assignee}</div>` : ''}
        </div>
    `;
}

// Drag and Drop функции
function dragStart(e) {
    draggedCard = e.target;
    e.target.classList.add('dragging');
    e.dataTransfer.setData('text/plain', e.target.dataset.cardId);
}

function dragEnd(e) {
    e.target.classList.remove('dragging');
    draggedCard = null;
}

function allowDrop(e) {
    e.preventDefault();
    e.currentTarget.classList.add('drag-over');
}

function dropCard(e, columnId) {
    e.preventDefault();
    e.currentTarget.classList.remove('drag-over');

    const cardId = e.dataTransfer.getData('text/plain');
    moveCardToColumn(cardId, columnId);
}

// Перемещение карточки между колонками
async function moveCardToColumn(cardId, columnId) {
    try {
        const response = await fetch(`${API_BASE}/cards/${cardId}/move`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ column_id: columnId }),
        });

        if (response.ok) {
            await refreshBoard();
        } else {
            const error = await response.json();
            alert(`Ошибка перемещения: ${error.error}`);
        }
    } catch (error) {
        alert('Ошибка соединения с сервером');
    }
}

// Открытие модального окна для создания/редактирования карточки
function openCardModal(columnId, cardId = null) {
    currentColumnId = columnId;
    editingCardId = cardId;

    const modal = document.getElementById('card-modal');
    const modalTitle = document.getElementById('modal-title');
    const form = document.getElementById('card-form');

    // Очищаем форму
    form.reset();
    document.getElementById('modal-error').classList.add('hidden');

    if (cardId) {
        // Режим редактирования
        modalTitle.textContent = 'Редактировать карточку';
        loadCardData(cardId);
    } else {
        // Режим создания
        modalTitle.textContent = 'Новая карточка';
    }

    modal.style.display = 'block';
}

// Загрузка данных карточки для редактирования
async function loadCardData(cardId) {
    try {
        const response = await fetch(`${API_BASE}/board`, {
            credentials: 'include'
        });

        if (response.ok) {
            const board = await response.json();
            let foundCard = null;

            // Находим карточку
            for (const column of board.columns) {
                foundCard = column.cards.find(card => card.id === cardId);
                if (foundCard) break;
            }

            if (foundCard) {
                document.getElementById('card-title').value = foundCard.title;
                document.getElementById('card-description').value = foundCard.description || '';
                document.getElementById('card-assignee').value = foundCard.assignee || '';
            }
        }
    } catch (error) {
        console.error('Ошибка загрузки данных карточки:', error);
    }
}

// Закрытие модального окна
function closeCardModal() {
    document.getElementById('card-modal').style.display = 'none';
    currentColumnId = null;
    editingCardId = null;
}

// Сохранение карточки
async function saveCard(e) {
    e.preventDefault();

    const title = document.getElementById('card-title').value;
    const description = document.getElementById('card-description').value;
    const assignee = document.getElementById('card-assignee').value;

    const errorDiv = document.getElementById('modal-error');
    errorDiv.classList.add('hidden');

    try {
        let response;

        if (editingCardId) {
            // Обновление существующей карточки
            response = await fetch(`${API_BASE}/cards/${editingCardId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ title, description, assignee }),
            });
        } else {
            // Создание новой карточки
            response = await fetch(`${API_BASE}/cards`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    title,
                    description,
                    assignee,
                    column_id: currentColumnId
                }),
            });
        }

        if (response.ok) {
            closeCardModal();
            await refreshBoard();
        } else {
            const error = await response.json();
            showError('modal-error', error.error);
        }
    } catch (error) {
        showError('modal-error', 'Ошибка соединения с сервером');
    }
}

// Редактирование карточки
function editCard(cardId) {
    // Находим колонку карточки
    const cardElement = document.querySelector(`[data-card-id="${cardId}"]`);
    const columnElement = cardElement.closest('.column');
    const columnId = columnElement.dataset.columnId;

    openCardModal(columnId, cardId);
}

// Удаление карточки
async function deleteCard(cardId) {
    if (!confirm('Удалить карточку?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/cards/${cardId}`, {
            method: 'DELETE',
            credentials: 'include',
        });

        if (response.ok) {
            await refreshBoard();
        } else {
            const error = await response.json();
            alert(`Ошибка удаления: ${error.error}`);
        }
    } catch (error) {
        alert('Ошибка соединения с сервером');
    }
}

// Открытие модального окна для создания/редактирования колонки
function openColumnModal(columnId = null, currentName = '') {
    const modal = document.getElementById('column-modal');
    const modalTitle = document.getElementById('column-modal-title');
    const form = document.getElementById('column-form');
    const nameInput = document.getElementById('column-name');

    // Очищаем форму
    form.reset();
    document.getElementById('column-modal-error').classList.add('hidden');

    if (columnId) {
        // Режим редактирования
        modalTitle.textContent = 'Редактировать колонку';
        nameInput.value = currentName;
        form.onsubmit = (e) => saveColumn(e, columnId);
    } else {
        // Режим создания
        modalTitle.textContent = 'Новая колонка';
        form.onsubmit = (e) => saveColumn(e);
    }

    modal.style.display = 'block';
}

// Закрытие модального окна колонки
function closeColumnModal() {
    document.getElementById('column-modal').style.display = 'none';
}

// Сохранение колонки
async function saveColumn(e, columnId = null) {
    e.preventDefault();

    const name = document.getElementById('column-name').value;
    const errorDiv = document.getElementById('column-modal-error');
    errorDiv.classList.add('hidden');

    try {
        let response;

        if (columnId) {
            // Обновление существу��щей колонки
            response = await fetch(`${API_BASE}/columns/${columnId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ name }),
            });
        } else {
            // Создание новой колонки
            response = await fetch(`${API_BASE}/columns`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ name }),
            });
        }

        if (response.ok) {
            closeColumnModal();
            await refreshBoard();
        } else {
            const error = await response.json();
            showError('column-modal-error', error.error);
        }
    } catch (error) {
        showError('column-modal-error', 'Ошибка соединения с сервером');
    }
}

// Редактирование колонки
function editColumn(columnId, currentName) {
    openColumnModal(columnId, currentName);
}

// Удаление колонки
async function deleteColumn(columnId) {
    if (!confirm('Удалить колонку? Все карточки в ней будут также удалены.')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/columns/${columnId}`, {
            method: 'DELETE',
            credentials: 'include',
        });

        if (response.ok) {
            await refreshBoard();
        } else {
            const error = await response.json();
            alert(`Ошибка удаления: ${error.error}`);
        }
    } catch (error) {
        alert('Ошибка соединения с сервером');
    }
}

// Обновление доски
async function refreshBoard() {
    try {
        const response = await fetch(`${API_BASE}/board`, {
            credentials: 'include'
        });

        if (response.ok) {
            const board = await response.json();
            renderBoard(board);
        }
    } catch (error) {
        console.error('Ошибка обновления доски:', error);
    }
}

// Выход из системы
async function logout() {
    try {
        await fetch(`${API_BASE}/logout`, {
            method: 'POST',
            credentials: 'include',
        });
    } catch (error) {
        console.error('Ошибка выхода:', error);
    } finally {
        // Перезагружаем страницу для сброса состояния
        window.location.reload();
    }
}

// Показать форму входа
function showLoginForm() {
    document.getElementById('create-board-form').classList.add('hidden');
    document.getElementById('login-form').classList.remove('hidden');
}

// Показать форму создания
function showCreateForm() {
    document.getElementById('login-form').classList.add('hidden');
    document.getElementById('create-board-form').classList.remove('hidden');
}

// Показать ошибку
function showError(elementId, message) {
    const errorDiv = document.getElementById(elementId);
    errorDiv.textContent = message;
    errorDiv.classList.remove('hidden');
}

// Удаление эффекта drag-over при уходе мыши
document.addEventListener('dragleave', function(e) {
    if (e.target.classList.contains('drop-zone')) {
        e.target.classList.remove('drag-over');
    }
});
