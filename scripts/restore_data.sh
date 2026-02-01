#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

echo "Copying images from uploads_backup to uploads..."
mkdir -p uploads
cp -r uploads_backup/* uploads/

echo "Resetting Product DB schema..."
docker exec -i product_db psql -U $PRODUCT_DB_USER -d $PRODUCT_DB_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

echo "Restoring Product DB..."
if [ -f backup_product_db.sql ]; then
    cat backup_product_db.sql | docker exec -i product_db psql -U $PRODUCT_DB_USER -d $PRODUCT_DB_NAME
else
    echo "backup_product_db.sql not found!"
fi

echo "Resetting User DB schema..."
docker exec -i user_db psql -U $USER_DB_USER -d $USER_DB_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

echo "Restoring User DB..."
if [ -f backup_user_db.sql ]; then
    cat backup_user_db.sql | docker exec -i user_db psql -U $USER_DB_USER -d $USER_DB_NAME
else
    echo "backup_user_db.sql not found!"
fi

echo "Restoration complete."
