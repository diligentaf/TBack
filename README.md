# Before you start this module, you need a DB running with a schema named "HaruDB"
# For those of you who don't how to set up DB, simply run this command

sudo docker run -d --name haruDB -p 3377:3306 -v /opt/datadir:/var/lib/mariadb -e MYSQL_ROOT_PASSWORD=BillionDollar12! mariadb:10.2

# The command above will download and run Maria DB's image with a port listening on 3377
# Once the Database starts running, you need to create a schema called "HaruDB"
# You may ask your fellow developer who sits next to you for help

# Once DB is completely set up, open TBack's module and run these commands:
dep init
dep ensure

# Press F5 to run the module
# You're all set 