CREATE TABLE IF NOT EXISTS colleges (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    token CHAR(43) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS campuses (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    college_id INT UNSIGNED NOT NULL,

    UNIQUE(name, college_id),
    FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS student_groups (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    campus_id INT UNSIGNED NOT NULL,

    UNIQUE(name, campus_id),
    FOREIGN KEY (campus_id) REFERENCES campuses(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS lessons (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

    title VARCHAR(255) NOT NULL,
    cabinet VARCHAR(20) NOT NULL,
    teacher VARCHAR(255) NOT NULL,

    date DATE NOT NULL,
    lesson_order TINYINT UNSIGNED NOT NULL,

    student_group_id INT UNSIGNED NOT NULL,

    UNIQUE(student_group_id, date, lesson_order),
    FOREIGN KEY (student_group_id) REFERENCES student_groups(id)
);

CREATE TABLE IF NOT EXISTS calls (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

    weekday TINYINT UNSIGNED NOT NULL,
    begins TIME NOT NULL,
    ends TIME NOT NULL,
    call_order TINYINT UNSIGNED NOT NULL,

    college_id INT UNSIGNED NOT NULL,

    UNIQUE (college_id, weekday, call_order),
    FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE,
    CHECK (weekday BETWEEN 0 AND 6),
    CHECK (begins < ends)
);
