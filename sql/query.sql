-- name: GetAllTheData :many
SELECT
    *
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN courses ON course_evaluation_map.course_number = courses.course_number
    JOIN ratings ON reviews.evaluation_id = ratings.evaluation_id
WHERE
    reviews.published = 'verified';

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
    reviews.published = 'verified'
GROUP BY
    courses.course_name,
    courses.course_number
ORDER BY
    DATE DESC;

-- name: GetStats :one
SELECT
    COUNT(DISTINCT course_number) AS total_courses,
    COUNT(*) AS total_reviews
FROM
    course_evaluation_map
    JOIN reviews ON reviews.evaluation_id = course_evaluation_map.id
WHERE
    published = 'verified';

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
    AND reviews.published = 'verified'
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
    course_evaluation_map.course_number,
    course_name,
    reviews.id as ReviewId,
    ratings.id as RatingId
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN ratings ON reviews.evaluation_id = ratings.evaluation_id
    JOIN courses ON course_evaluation_map.course_number = courses.course_number
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
INSERT INTO
    reviews (evaluation_id, review)
VALUES
    (@evaluation_id, @review) RETURNING *;

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
    review,
    course_evaluation_map.id
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
WHERE
    user_id = @user_id
    AND course_number = @course_number;

-- name: SetCourseEvaluationMap :one
INSERT INTO
    course_evaluation_map (user_id, course_number, semester)
VALUES
    (@user_id, @course_number, @semester) RETURNING id;

-- name: GetCourseEvaluationMap :one
SELECT
    id
FROM
    course_evaluation_map
WHERE
    user_id = @user_id
    AND course_number = @course_number;

-- name: UpdateSemester :one
UPDATE
    course_evaluation_map
SET
    semester = @semester
WHERE
    id = @evaluation_id
    AND user_id = @user_id RETURNING *;

-- name: UpdateReview :one
UPDATE
    reviews
SET
    review = @review
FROM
    course_evaluation_map
WHERE
    reviews.evaluation_id = course_evaluation_map.id
    AND course_evaluation_map.id = @evaluation_id
    AND course_evaluation_map.user_id = @user_id RETURNING *;

-- name: GetCurrentSemester :many
SELECT
    semester
FROM
    current_semester;

-- name: RemoveCurrentSemester :one
TRUNCATE TABLE current_semester;

-- name: SetCurrentSemester :many
INSERT INTO
    current_semester (semester)
VALUES
    (@semester) RETURNING *;

-- name: DeleteReview :one
DELETE FROM
    reviews
WHERE
    evaluation_id = @evaluation_id RETURNING *;

-- name: UpdateRating :one
UPDATE
    ratings
SET
    recommended = @recommended,
    engaging = @engaging,
    difficulty = @difficulty,
    effort = @effort,
    resources = @resources
WHERE
    evaluation_id = @evaluation_id RETURNING *;

-- name: DeleteRating :one
DELETE FROM
    ratings
WHERE
    evaluation_id = @evaluation_id RETURNING *;

-- name: SetModerator :one
UPDATE
    users
SET
    moderator = NOT moderator
WHERE
    user_id = @user_id RETURNING *;

-- name: GetModerators :many
SELECT
    user_id,
    user_name
FROM
    users
WHERE
    moderator = TRUE;

-- name: GetUnverifiedReviews :many
SELECT
    reviews.review,
    course_evaluation_map.course_number,
    courses.course_name,
    course_evaluation_map.user_id,
    course_evaluation_map.id
FROM
    reviews
    JOIN course_evaluation_map ON reviews.evaluation_id = course_evaluation_map.id
    JOIN courses ON course_evaluation_map.course_number = courses.course_number
WHERE
    reviews.published = 'pending';

-- name: VerifyReview :one
UPDATE
    reviews
SET
    published = 'verified'
WHERE
    evaluation_id = @evaluation_id RETURNING *;

-- name: RejectReview :one
UPDATE
    reviews
SET
    published = 'rejected'
WHERE
    evaluation_id = @evaluation_id RETURNING *;

-- name: AddCourse :one
INSERT INTO
    courses (course_number, course_name)
VALUES
    (@course_number, @course_name) RETURNING *;

-- name: GetCourseName :one
SELECT
    course_name
FROM
    courses
WHERE
    course_number = @course_number;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    user_id = @user_id;

-- name: CheckUserWithId :one
SELECT
    *
FROM
    course_evaluation_map
WHERE
    id = @evaluation_id
    AND user_id = @user_id;