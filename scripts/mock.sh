#!/bin/bash
set -e

# clear tables
psql -c "TRUNCATE TABLE courses, users, course_evaluation_map, reviews, ratings, actions, event_log, course_number_alias, current_semester RESTART IDENTITY CASCADE;"

# add mock data
psql -f scripts/mock.sql
