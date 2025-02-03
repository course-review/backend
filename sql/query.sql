-- name: GetAllTheData :many
SELECT
    *
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN courses ON course_evaluation_map.course_number = courses.course_number
    JOIN ratings ON reviews.evaluation_id = ratings.evaluation_id
WHERE
    reviews.published = TRUE;

-- name: GetReviewedCourses :many
SELECT
    courses.course_name,
    courses.course_number,
    MAX(reviews.date) AS date
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN courses ON course_evaluation_map.course_number = courses.course_number
WHERE
    reviews.published = TRUE
GROUP BY
    courses.course_name,
    courses.course_number;

-- name: GetStats :one
SELECT
    COUNT(DISTINCT course_number) AS total_courses,
    COUNT(*) AS total_reviews
FROM
    course_evaluation_map
WHERE
    published = TRUE;

-- name: GetRatings :many
SELECT
    recommended,
    engaging,
    difficulty,
    effort,
    resources
FROM
    ratings
    JOIN course_evaluation_map ON ratings.evaluation_id = course_evaluation_map.id
WHERE
    course_number = @course_number;

-- name: GetReviews :many
SELECT
    review,
    semester
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
WHERE
    course_number = @course_number
    AND reviews.published = TRUE
ORDER BY
    reviews.date DESC;

-- name: GetUserData :many
SELECT
    reviews.review,
    ratings.recommended,
    ratings.engaging,
    ratings.difficulty,
    ratings.effort,
    ratings.resources,
    semester,
    course_number
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN ratings ON reviews.evaluation_id = ratings.evaluation_id
WHERE
    user_id = @user_id;

-- -- name: Review :one
-- WITH evaluation AS (
--     INSERT INTO course_evaluation_map (user_id, course_number, semester)
--     VALUES (@user_id, @course_number, @semester)
--     RETURNING id
-- )
-- INSERT INTO reviews (evaluation_id, review)
-- SELECT
--     id,
--     @review
-- FROM
--     evaluation
-- RETURNING *;

-- name: SetMapping :many
INSERT INTO
    course_evaluation_map (user_id, course_number, semester)
VALUES
    (@user_id, @course_number, @semester) RETURNING *;

-- name: SetReview :many
WITH evaluation AS (
    SELECT
        id
    FROM
        course_evaluation_map
    WHERE
        user_id = @user_id
        AND course_number = @course_number
)
INSERT INTO
    reviews (evaluation_id, review)
VALUES
    (evaluation, @review) RETURNING *;

-- name: SetRating :many
INSERT INTO
    ratings (
        evaluation_id,
        recommended,
        engaging,
        difficulty,
        effort,
        resources
    )
VALUES
    (
        @evaluation_id,
        @recommended,
        @engaging,
        @difficulty,
        @effort,
        @resources
    ) RETURNING *;

-- name: SetEventLog :many
INSERT INTO
    event_log (evaluation_id, user_id, action_id, info)
VALUES
    (@evaluation_id, @user_id, @action_id, @info) RETURNING *;

-- name: SetUser :many
INSERT INTO
    users (user_id, user_name)
VALUES
    (@user_id, @user_name) RETURNING *;

-- name: SetCourse :many
INSERT INTO
    courses (course_number, course_name)
VALUES
    (@course_number, @course_name) RETURNING *;

-- name: SetCourseAlias :many
INSERT INTO
    course_number_alias (source, target)
VALUES
    (@source, @target) RETURNING *;

-- name: SetAction :many
INSERT INTO
    actions (name)
VALUES
    (@name) RETURNING *;

-- name: GetCourses :many
SELECT
    course_number,
    course_name
FROM
    courses;

-- name: GetUserReviewsAndRatings :many
SELECT
    reviews.review,
    ratings.recommended,
    ratings.engaging,
    ratings.difficulty,
    ratings.effort,
    ratings.resources,
    course_evaluation_map.semester
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN ratings ON reviews.evaluation_id = ratings.evaluation_id
WHERE
    user_id = @user_id;

-- name: GetCourseReviews :many
SELECT
    reviews.review,
    course_evaluation_map.semester
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
WHERE
    course_number = @course_number;

-- name: GetCourseRatings :many
SELECT
    recommended,
    engaging,
    difficulty,
    effort,
    resources
FROM
    ratings
    JOIN course_evaluation_map ON ratings.evaluation_id = course_evaluation_map.id
WHERE
    course_number = @course_number;

-- name: GetLogs :many
SELECT
    courses.course_number,
    course_name,
    user_name,
    actions.name,
    info,
    date
FROM
    event_log
    JOIN actions ON event_log.action_id = actions.id
    JOIN course_evaluation_map ON event_log.evaluation_id = course_evaluation_map.id
    JOIN courses ON course_evaluation_map.course_number = courses.course_number
    JOIN users ON event_log.user_id = users.user_id
WHERE
    DATE(date) BETWEEN @start_date
    AND @end_date;

-- name: GetReview :one
SELECT
    review
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
WHERE
    user_id = @user_id
    AND course_number = @course_number;

-- name: GetCourseEvaluationMap :one
SELECT
    id
FROM
    course_evaluation_map
WHERE
    user_id = @user_id
    AND course_number = @course_number
    AND semester = @semester;