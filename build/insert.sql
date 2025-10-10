-- =======================
-- Таблица users
-- =======================
INSERT INTO users (id, login, hashed_password, is_moderator)
VALUES
(1, 'admin', 'hashed_password_admin', true),
(2, 'user1', 'hashed_password_user1', false);

-- =======================
-- Таблица texts
-- =======================
INSERT INTO texts (id, is_delete, image_url, title, description, price)
VALUES
(1, false, 'http://localhost:9000/img/img/sci.jpg',
 'Научная статья',
 'Мы анализируем научные тексты: проверяем структуру, терминологическую насыщенность и вычисляем индексы читабельности (Флеша, Фога, SMOG). Отчёт покажет, насколько материал удобен для чтения и как его упростить.',
 125),
(2, false, 'http://localhost:9000/img/img/his.jpg',
 'Исторический документ',
 'Подходит для архивных и исторических материалов. Оцениваем длину предложений, редкие конструкции и устаревшие формы, строим метрики читабельности и даём рекомендации по адаптации.',
 400),
(3, false, 'http://localhost:9000/img/img/lit.jpg',
 'Художественный текст',
 'Анализ художественных фрагментов: ритм и длины фраз, доля диалогов, простота восприятия. Сравниваем индексы с популярной литературой и подсказываем, где можно повысить «легкость» чтения.',
 500);

-- =======================
-- Таблица read_indxs
-- =======================
INSERT INTO read_indxs (id, status, date_create, creator_id, date_form, date_end, moderator_id, comments, contacts)
VALUES
(1, 'DRAFT', NOW(), 2, NOW(), NULL, 1,
 'Примерный комментарий к расчёту индекса читабельности',
 'test@example.com');

-- =======================
-- Таблица read_indxs_to_text
-- =======================
INSERT INTO read_indxs_to_texts (text_id, read_indxs_id, calculation, count_words, count_sentences, count_syllables)
VALUES
(1, 1, 0, 1000, 100, 120),
(2, 1, 0, 1000, 100, 120),
(3, 1, 0, 1000, 100, 120);
