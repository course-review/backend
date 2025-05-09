-- up migration: initial schema
CREATE TABLE IF NOT EXISTS courses (
    course_number VARCHAR(12) PRIMARY KEY, -- Unique identifier for the course
    course_name TEXT NOT NULL -- Name of the course
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(128) PRIMARY KEY, -- Unique identifier for the user
    admin BOOLEAN DEFAULT FALSE, -- Indicates if the user is an admin
    moderator BOOLEAN DEFAULT FALSE -- Indicates if the user is a moderator
);

CREATE TABLE IF NOT EXISTS course_evaluation_map (
    id SERIAL PRIMARY KEY, -- Unique ID for each mapping
    user_id VARCHAR(128) NOT NULL, -- User associated with the evaluation
    course_number VARCHAR(12) NOT NULL, -- Course associated with the evaluation
    semester VARCHAR(4) DEFAULT NULL, -- Semester in which the evaluation was performed
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (course_number) REFERENCES courses(course_number)
);

CREATE TYPE status AS ENUM ('pending', 'verified', 'rejected');

CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY, -- Unique identifier for the review
    evaluation_id INTEGER NOT NULL, -- Reference to the evaluation
    date DATE DEFAULT NOW(), -- Date of the review
    published status DEFAULT 'pending', -- Indicates if the review is published
    review TEXT NOT NULL, -- Content of the review
    requested_changes TEXT DEFAULT NULL, -- Changes requested for the review
    old_review TEXT DEFAULT NULL, -- old version of the review after edit
    FOREIGN KEY (evaluation_id) REFERENCES course_evaluation_map(id),
    UNIQUE (evaluation_id)
);

CREATE TABLE IF NOT EXISTS ratings (
    id SERIAL PRIMARY KEY, -- Unique identifier for the rating
    evaluation_id INTEGER NOT NULL, -- Reference to the evaluation
    date DATE DEFAULT NOW(), -- Date of the rating
    recommended FLOAT DEFAULT NULL CHECK (recommended BETWEEN 1.0 AND 5.0), -- Enforces 1.0 - 5.0 range
    engaging FLOAT DEFAULT NULL CHECK (engaging BETWEEN 1.0 AND 5.0),
    difficulty FLOAT DEFAULT NULL CHECK (difficulty BETWEEN 1.0 AND 5.0),
    effort FLOAT DEFAULT NULL CHECK (effort BETWEEN 1.0 AND 5.0),
    resources FLOAT DEFAULT NULL CHECK (resources BETWEEN 1.0 AND 5.0),
    FOREIGN KEY (evaluation_id) REFERENCES course_evaluation_map(id),
    UNIQUE (evaluation_id) -- Ensures one rating per evaluation
);

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY, -- Unique identifier for the action
    name TEXT NOT NULL -- Name of the action
);

CREATE TABLE IF NOT EXISTS event_log (
    id SERIAL PRIMARY KEY, -- Unique identifier for the log entry
    evaluation_id INTEGER, -- Reference to the evaluation
    user_id VARCHAR(128), -- User associated with the log entry
    action_id INTEGER, -- Action performed
    info TEXT, -- Additional information
    date DATE DEFAULT NOW(), -- Date of the event
    FOREIGN KEY (evaluation_id) REFERENCES course_evaluation_map(id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (action_id) REFERENCES actions(id)
);

CREATE TABLE IF NOT EXISTS course_number_alias (
    id SERIAL PRIMARY KEY, -- Unique identifier for the alias
    source VARCHAR(12) NOT NULL, -- Original course number
    target VARCHAR(12) NOT NULL, -- Alias course number
    FOREIGN KEY (source) REFERENCES courses(course_number),
    FOREIGN KEY (target) REFERENCES courses(course_number)
);

CREATE TABLE IF NOT EXISTS current_semester (
    semester VARCHAR(4) PRIMARY KEY -- Current semester
);
