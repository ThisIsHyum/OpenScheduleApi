ALTER TABLE lessons
    DROP FOREIGN KEY fk_lessons_student_group;

ALTER TABLE lessons
    ADD CONSTRAINT lessons_ibfk_1
    FOREIGN KEY (student_group_id)
    REFERENCES student_groups(id);