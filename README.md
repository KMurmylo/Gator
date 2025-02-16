Require:
Go v1.20+
postgres 16.6+

setup
create database for "gator"  by running command 
CREATE DATABASE gator;
 in database.
Then, apply the provided schema by running:

psql -U <username> -d gator -f sql/schema.sql
Replace <username> with your Postgres username."

In the gatorconfig.json file, set the db_url field to point to your Postgres database. For example:

{
  "db_url": "postgres://<DB username>:<DB password>@localhost:5432/gator?sslmode=disable"
}
Replace <DB username> and <DB password> with your Postgres credentials."

gator commands:
register (name)
Registers a name in the database and switches to that user

login (name)
switches to previously registered user

reset
Deletes all users, posts and feeds 

users
lists users

agg (time)
starts listening to the feeds at given internal , leave blank for 2seconds

addfeed (name) (url)
Adds feed to the current user

feeds 
lists all the feeds

follow (url)
Adds current user to followers of the given feed

following
lists feeds that current user is following

unfollow (url)
removes user from the followers of the given feed 

browse (number)
display the latest number of posts that the current user is following, leave blank for 2 latest posts
