CREATE TABLE courses (
    course_number VARCHAR(12) PRIMARY KEY, -- Unique identifier for the course
    course_name TEXT NOT NULL -- Name of the course
);

CREATE TABLE users (
    user_id VARCHAR(16) PRIMARY KEY, -- Unique identifier for the user
    user_name TEXT, -- Name of the user
    admin BOOLEAN DEFAULT FALSE, -- Indicates if the user is an admin
    moderator BOOLEAN DEFAULT FALSE, -- Indicates if the user is a moderator
    deactivated BOOLEAN DEFAULT FALSE, -- Indicates if the user account is deactivated
    last_logged_out DATE DEFAULT NULL -- Timestamp of the last logout
);

CREATE TABLE course_evaluation_map (
    id SERIAL PRIMARY KEY, -- Unique ID for each mapping
    user_id VARCHAR(16) NOT NULL, -- User associated with the evaluation
    course_number VARCHAR(12) NOT NULL, -- Course associated with the evaluation
    semester VARCHAR(4) DEFAULT NULL, -- Semester in which the evaluation was performed
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (course_number) REFERENCES courses(course_number)
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY, -- Unique identifier for the review
    evaluation_id INTEGER NOT NULL, -- Reference to the evaluation
    date DATE DEFAULT NOW(), -- Date of the review
    published BOOLEAN DEFAULT FALSE, -- Indicates if the review is published
    review TEXT NOT NULL, -- Content of the review
    requested_changes TEXT DEFAULT NULL, -- Changes requested for the review
    updated_review TEXT DEFAULT NULL, -- Updated version of the review
    FOREIGN KEY (evaluation_id) REFERENCES course_evaluation_map(id)
);

CREATE TABLE ratings (
    id SERIAL PRIMARY KEY, -- Unique identifier for the rating
    evaluation_id INTEGER NOT NULL, -- Reference to the evaluation
    date DATE NOT NULL, -- Date of the rating
    recommended INTEGER DEFAULT NULL, -- Rating for recommendation
    engaging INTEGER DEFAULT NULL, -- Rating for engagement
    difficulty INTEGER DEFAULT NULL, -- Rating for difficulty
    effort INTEGER DEFAULT NULL, -- Rating for effort
    resources INTEGER DEFAULT NULL, -- Rating for resources
    FOREIGN KEY (evaluation_id) REFERENCES course_evaluation_map(id)
);

CREATE TABLE event_log (
    id SERIAL PRIMARY KEY, -- Unique identifier for the log entry
    evaluation_id INTEGER, -- Reference to the evaluation
    user_id VARCHAR(16), -- User associated with the log entry
    action_id INTEGER, -- Action performed
    info TEXT, -- Additional information
    date DATE DEFAULT NOW(), -- Date of the event
    FOREIGN KEY (evaluation_id) REFERENCES course_evaluation_map(id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (action_id) REFERENCES actions(id)
);

CREATE TABLE actions (
    id SERIAL PRIMARY KEY, -- Unique identifier for the action
    name TEXT NOT NULL -- Name of the action
);

CREATE TABLE course_number_alias (
    id SERIAL PRIMARY KEY, -- Unique identifier for the alias
    source VARCHAR(12) NOT NULL, -- Original course number
    target VARCHAR(12) NOT NULL, -- Alias course number
    FOREIGN KEY (source) REFERENCES courses(course_number),
    FOREIGN KEY (target) REFERENCES courses(course_number)
);

CREATE TABLE current_semester (
    semester VARCHAR(4) PRIMARY KEY -- Current semester
);
