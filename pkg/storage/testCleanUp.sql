-- удаляем все данные из таблиц
TRUNCATE TABLE tasks_labels, labels, tasks, users;

-- восстанавливаем последовательности
ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 5;
ALTER SEQUENCE IF EXISTS labels_id_seq RESTART WITH 6;
ALTER SEQUENCE IF EXISTS tasks_id_seq RESTART WITH 5;

-- восстанавливаем данные в таблцах
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