DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS labels (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
    closed BIGINT DEFAULT 0,
    author_id INTEGER REFERENCES users(id) DEFAULT 0,
    assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks_labels (
    task_id INTEGER REFERENCES tasks(id),
    label_id INTEGER REFERENCES labels(id)
);

TRUNCATE TABLE users, labels CASCADE;
INSERT INTO users (id, name) VALUES (0, 'default');
INSERT INTO users (name) VALUES ('John Smit'), ('Sarah Brown'), ('Edward Peacock');
INSERT INTO labels (id, name) VALUES (0, 'default');
INSERT INTO labels (name) VALUES ('Task'), ('Bug');
INSERT INTO tasks (title, content) VALUES ('Write a report', 'Write a complete report on current situation in IT marketing'), ('Prepare presentation', 'Presentation about program new features for customers');
INSERT INTO tasks_labels (task_id, label_id) VALUES (1, 1), (2,1), (3,2);