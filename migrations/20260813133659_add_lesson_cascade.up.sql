ALTER TABLE lessons
    DROP FOREIGN KEY lessons_ibfk_1;

ALTER TABLE lessons
    ADD CONSTRAINT fk_lessons_student_group
    FOREIGN KEY (student_group_id)
    REFERENCES student_groups(id)
    ON DELETE CASCADE;