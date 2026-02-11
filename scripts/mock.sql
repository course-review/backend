INSERT INTO users (user_id, admin, moderator) VALUES
('u_001', TRUE, TRUE),
('u_002', FALSE, FALSE);

INSERT INTO courses (course_number, course_name) VALUES
('263-3010-00L', 'Introduction to Computer Science'),
('252-3900-00L', 'Calculus II');

INSERT INTO actions (name) VALUES
('evaluation_submitted'),
('review_published');

INSERT INTO current_semester (semester) VALUES
('23FS'),
('23HS');

INSERT INTO course_evaluation_map (id, user_id, course_number, semester) VALUES
(1, 'u_001', '263-3010-00L', '23FS'),
(2, 'u_002', '252-3900-00L', '23HS');

INSERT INTO reviews (evaluation_id, published, review, requested_changes) VALUES
(1, 'verified', 'Excellent course content.', NULL),
(2, 'pending', 'The pacing was very fast.', 'Please elaborate on the difficulty.');

INSERT INTO ratings (evaluation_id, recommended, engaging, difficulty, effort, resources) VALUES
(1, 5, 5, 3, 4, 5),
(2, null, 2, 5, 5, 3);

INSERT INTO event_log (evaluation_id, user_id, action_id, info) VALUES
(1, 'u_001', 1, 'User submitted initial evaluation'),
(2, 'u_001', 2, 'Admin published the review');

INSERT INTO course_number_alias (source, target) VALUES
('263-3010-00L', '252-3900-00L'),
('252-3900-00L', '263-3010-00L');
