#!/bin/bash
sudo mysql -e "CREATE DATABASE IF NOT EXISTS tay_photos;"
sudo mysql -e "CREATE USER IF NOT EXISTS tayphoto_user@localhost IDENTIFIED BY 'tayphoto_pass';"
sudo mysql -e "GRANT ALL PRIVILEGES ON tay_photos.* TO tayphoto_user@localhost;"
sudo mysql -e "FLUSH PRIVILEGES;"
mysql -utayphoto_user -ptayphoto_pass tay_photos < ../server/sql/create_tables.sql
mysql -utayphoto_user -ptayphoto_pass tay_photos < ../server/sql/insert_test_data.sql