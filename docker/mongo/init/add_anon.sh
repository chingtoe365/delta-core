#!/bin/bash
set -e

mongosh "mongodb://localhost:27017" --username ${DB_USER} --password ${DB_PASS} <<EOF
use ${DB_NAME}
db.users.insertOne({"_id": ObjectId("${MONGO_ANON_USER_ID}"), "name": "${MONGO_ANON_USER_NAME}", "email": "${MONGO_ANON_EMAIL}"})
EOF
