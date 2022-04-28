DROP TABLE IF EXISTS users, labels, tasks, tasks_labels;

-- таблица пользователи
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

-- таблица метки
CREATE TABLE IF NOT EXISTS labels (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

-- таблица задачи
CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
	opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
	closed BIGINT NOT NULL DEFAULT 0,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	author_id INTEGER DEFAULT 0,
	assigned_id INTEGER DEFAULT 0,
	FOREIGN KEY(author_id) REFERENCES users(id),
	FOREIGN KEY(assigned_id) REFERENCES users(id)
);

-- таблица многие-ко-многим между задачами и метками
CREATE TABLE IF NOT EXISTS tasks_labels (
	task_id INTEGER,
	label_id INTEGER,
	PRIMARY KEY(task_id, label_id),
	FOREIGN KEY (task_id) REFERENCES tasks(id),
	FOREIGN KEY (label_id) REFERENCES labels(id)
);

-- заполняем таблицы
INSERT INTO users(id, name) 
VALUES (1, 'Walter White'), 
(2, 'Jesse Pinkman'), 
(3, 'Gus Fring'), 
(4, 'Jimmy "Saul Goodman" McGill');

INSERT INTO tasks(id, opened, closed, title, content, author_id, assigned_id) 
VALUES (1, 1648805555, 1649583155, 'Cook meth', 'Make 99.1% pure crystals', 3, 1),
(2, 1649669555, 1650879155, 'Law protection', 'Get Jesse Pinkman out of jail', 3, 4),
(3, 1650015155, 1650274355, 'Spreading', 'Drugs spreading', 1, 2),
(4, 1648848755, 1648935155, 'Meth contraband', 'Deliver meth over the mexican border', 3, 3);

INSERT INTO labels(id, name) 
VALUES (1, 'meth'), 
(2, 'law cover up'), 
(3, 'contraband'), 
(4, 'spreading'),
(5, 'production');

INSERT INTO tasks_labels(task_id, label_id) 
VALUES (1, 5), (1, 1), (2, 2), (3, 4), (3, 1), (4, 3), (4, 1);