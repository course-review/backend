INSERT INTO users (user_id, admin, moderator) VALUES
('u_001', TRUE, TRUE),
('u_002', FALSE, FALSE),
('u_003', FALSE, FALSE),
('u_004', FALSE, FALSE),
('u_005', FALSE, FALSE),
('u_006', FALSE, FALSE),
('u_007', FALSE, FALSE),
('u_008', FALSE, FALSE);

INSERT INTO courses (course_number, course_name) VALUES
('263-3010-00L', 'Introduction to Computer Science'),
('252-3900-00L', 'Calculus II'),
('252-0025-00L', 'Introduction to Programming'),
('263-2400-00L', 'Reliable and Interactive Systems'),
('401-0212-08L', 'Analysis III'),
('401-0663-00L', 'Numerical Methods for CSE'),
('402-0044-00L', 'Physics I'),
('151-0102-00L', 'Dynamics'),
('227-0102-00G', 'Systems and Signals'),
('363-1100-00L', 'Technology and Innovation Management'),
('851-0740-00S', 'Introduction to Philosophy of Science'),
('529-0191-00L', 'General Chemistry'),
('551-0007-00L', 'General Biology'),
('101-0128-00G', 'Structural Design I'),
('052-0500-00S', 'History of Urban Design'),
('636-0007-00L', 'Introduction to Systems Biology'),
('701-0023-00L', 'Introduction to Climate Systems'),
('252-0058-00U', 'Programming Exercises'),
('252-0331-00P', 'Software Engineering Lab'),
('888-1234-00L', 'Special Topics Colloquium');

INSERT INTO actions (name) VALUES
('evaluation_submitted'),
('review_published');

INSERT INTO current_semester (semester) VALUES
('23FS'),
('23HS');

INSERT INTO course_evaluation_map (id, user_id, course_number, semester) VALUES
(1, 'u_001', '263-3010-00L', '23FS'),
(2, 'u_002', '252-3900-00L', '23HS'),
-- Introduction to Programming (D-INFK, Lecture) — 6 ratings
(3, 'u_003', '252-0025-00L', '23HS'),
(4, 'u_004', '252-0025-00L', '23HS'),
(5, 'u_005', '252-0025-00L', '24FS'),
(6, 'u_006', '252-0025-00L', '24FS'),
(7, 'u_007', '252-0025-00L', '24FS'),
(8, 'u_008', '252-0025-00L', '24HS'),
-- Reliable and Interactive Systems (D-INFK, Lecture) — 5 ratings
(9, 'u_003', '263-2400-00L', '23HS'),
(10, 'u_004', '263-2400-00L', '23HS'),
(11, 'u_005', '263-2400-00L', '24FS'),
(12, 'u_006', '263-2400-00L', '24FS'),
(13, 'u_007', '263-2400-00L', '24HS'),
-- Analysis III (D-MATH, Lecture) — 4 ratings, mixed/negative
(14, 'u_003', '401-0212-08L', '23FS'),
(15, 'u_004', '401-0212-08L', '23FS'),
(16, 'u_005', '401-0212-08L', '23HS'),
(17, 'u_006', '401-0212-08L', '24FS'),
-- Numerical Methods for CSE (D-MATH, Lecture) — 2 ratings
(18, 'u_003', '401-0663-00L', '23HS'),
(19, 'u_004', '401-0663-00L', '24FS'),
-- Physics I (D-PHYS, Lecture) — 3 ratings, low scores
(20, 'u_005', '402-0044-00L', '23FS'),
(21, 'u_006', '402-0044-00L', '23FS'),
(22, 'u_007', '402-0044-00L', '23HS'),
-- Dynamics (D-MAVT, Lecture) — 1 rating
(23, 'u_008', '151-0102-00L', '24FS'),
-- Systems and Signals (D-ITET, Lecture+Exercises) — 4 ratings
(24, 'u_003', '227-0102-00G', '23FS'),
(25, 'u_004', '227-0102-00G', '23HS'),
(26, 'u_005', '227-0102-00G', '24FS'),
(27, 'u_006', '227-0102-00G', '24HS'),
-- Technology and Innovation Management (D-MTEC, Lecture) — 2 ratings
(28, 'u_007', '363-1100-00L', '23HS'),
(29, 'u_008', '363-1100-00L', '24FS'),
-- Introduction to Philosophy of Science (D-GESS, Seminar) — 3 ratings
(30, 'u_003', '851-0740-00S', '23FS'),
(31, 'u_004', '851-0740-00S', '23HS'),
(32, 'u_005', '851-0740-00S', '24FS'),
-- General Chemistry (D-CHAB, Lecture) — 1 rating
(33, 'u_006', '529-0191-00L', '24HS'),
-- General Biology (D-BIOL, Lecture) — 2 ratings
(34, 'u_007', '551-0007-00L', '23FS'),
(35, 'u_008', '551-0007-00L', '23HS'),
-- Structural Design I (D-BAUG, Lecture+Exercises) — 1 rating
(36, 'u_003', '101-0128-00G', '24FS'),
-- History of Urban Design (D-ARCH, Seminar) — 1 rating
(37, 'u_004', '052-0500-00S', '23HS'),
-- Introduction to Systems Biology (D-BSSE, Lecture) — 1 rating
(38, 'u_005', '636-0007-00L', '24FS'),
-- Introduction to Climate Systems (D-USYS, Lecture) — 3 ratings
(39, 'u_006', '701-0023-00L', '23FS'),
(40, 'u_007', '701-0023-00L', '23HS'),
(41, 'u_008', '701-0023-00L', '24FS'),
-- Programming Exercises (D-INFK, Exercises) — 2 ratings
(42, 'u_003', '252-0058-00U', '23HS'),
(43, 'u_004', '252-0058-00U', '24FS'),
-- Software Engineering Lab (D-INFK, Lab/Praktikum) — 1 rating
(44, 'u_005', '252-0331-00P', '24HS'),
-- Special Topics Colloquium (unrecognized prefix -> Other) — 1 rating
(45, 'u_006', '888-1234-00L', '23FS');

INSERT INTO reviews (evaluation_id, published, review, requested_changes) VALUES
(1, 'verified', 'Excellent course content.', NULL),
(2, 'pending', 'The pacing was very fast.', 'Please elaborate on the difficulty.'),
(3, 'verified', 'Fantastic introduction to CS, very well structured.', NULL),
(9, 'verified', 'Great systems course, projects were super instructive.', NULL),
(14, 'pending', 'Analysis III moves fast, could use more examples.', 'Please add more examples.'),
(18, 'verified', 'Solid overview of numerical methods.', NULL),
(20, 'verified', 'Physics I is tough, exams are brutal.', NULL),
(24, 'rejected', 'ok', 'Please write a longer, more constructive review.'),
(28, 'verified', 'Learned a lot about innovation management.', NULL),
(30, 'verified', 'Interesting philosophy seminar, light workload.', NULL),
(34, 'verified', 'Solid intro to biology, engaging lectures.', NULL),
(36, 'verified', 'Good structural design fundamentals.', NULL),
(37, 'verified', 'Loved this architecture history course!', NULL),
(39, 'verified', 'Great intro to climate systems, relevant today.', NULL),
(42, 'verified', 'Exercises reinforce the lecture material well.', NULL),
(44, 'verified', 'Hands-on lab, learned practical skills.', NULL);

INSERT INTO ratings (evaluation_id, recommended, engaging, difficulty, effort, resources) VALUES
(1, 5, 5, 3, 4, 5),
(2, null, 2, 5, 5, 3),
-- Introduction to Programming — consistently high
(3, 5, 5, 4, 4, 5),
(4, 5, 4, 3, 4, 5),
(5, 4, 5, 4, 3, 4),
(6, 5, 4, 4, 4, 5),
(7, 4, 5, 3, 4, 4),
(8, 5, 4, 4, 3, 5),
-- Reliable and Interactive Systems — good, slightly harder
(9, 4, 4, 3, 3, 4),
(10, 4, 5, 3, 3, 4),
(11, 5, 4, 2, 2, 3),
(12, 3, 4, 3, 3, 4),
(13, 4, 3, 3, 3, 4),
-- Analysis III — mixed, too difficult
(14, 3, 2, 2, 2, 3),
(15, 2, 2, 1, 1, 3),
(16, 3, 3, 2, 2, 4),
(17, 3, 2, 2, 2, 3),
-- Numerical Methods for CSE — good
(18, 4, 4, 3, 4, 4),
(19, 5, 4, 4, 4, 5),
-- Physics I — low, brutal
(20, 2, 2, 1, 1, 2),
(21, 2, 3, 2, 2, 3),
(22, 1, 2, 1, 1, 2),
-- Dynamics
(23, 4, 4, 3, 4, 4),
-- Systems and Signals
(24, 4, 3, 3, 3, 4),
(25, 3, 4, 3, 3, 3),
(26, 4, 4, 2, 3, 4),
(27, 5, 4, 3, 4, 4),
-- Technology and Innovation Management
(28, 5, 5, 4, 4, 4),
(29, 4, 4, 4, 4, 5),
-- Introduction to Philosophy of Science — liked, light workload
(30, 5, 5, 4, 3, 4),
(31, 5, 4, 4, 3, 5),
(32, 4, 5, 3, 2, 4),
-- General Chemistry
(33, 3, 3, 3, 3, 3),
-- General Biology
(34, 4, 4, 3, 4, 4),
(35, 4, 5, 4, 3, 4),
-- Structural Design I
(36, 4, 3, 4, 4, 3),
-- History of Urban Design
(37, 5, 5, 4, 5, 4),
-- Introduction to Systems Biology
(38, 4, 4, 3, 4, 4),
-- Introduction to Climate Systems
(39, 5, 5, 4, 4, 5),
(40, 4, 4, 3, 4, 4),
(41, 5, 5, 4, 4, 5),
-- Programming Exercises — one partially-filled rating (null resources)
(42, 4, 3, 4, 4, 4),
(43, 3, 3, 3, 3, null),
-- Software Engineering Lab — partially-filled rating (null resources)
(44, 5, 5, 4, 5, null),
-- Special Topics Colloquium (Other dept)
(45, 3, 3, 3, 3, 3);

INSERT INTO event_log (evaluation_id, user_id, action_id, info) VALUES
(1, 'u_001', 1, 'User submitted initial evaluation'),
(2, 'u_001', 2, 'Admin published the review');

INSERT INTO course_number_alias (source, target) VALUES
('263-3010-00L', '252-3900-00L'),
('252-3900-00L', '263-3010-00L');
