This is a project from boot.dev that used PostgreSQL and Golang along with numerous tools like goose and sqlc.

## About:

Gator is a CLI tool that can be used to follow RSS feeds from multiple sources, save them as posts, and then the user can browse the RSS feeds that they followed. Gator also support
multiple users and only show users the feeds that they follow.

## Requirements:

Users are required to have PostgreSQL and Go installed.

### 1. Install PostgreSQL

**On macOS (using Homebrew):**

brew install postgresql

**on Ubuntu/linux**
sudo apt update
sudo apt install postgresql postgresql-contrib

### 2. Start PostgreSQL

**On macOS**

brew services start postgresql

**On Linux**

sudo systemctl start postgresql
sudo systemctl enable postgresql  # To start automatically on boot

### 3. Create the Database

**Connect to PostgreSQL**

psql -U postgres

**Create the database and user**

CREATE DATABASE gator;
CREATE USER gatoruser WITH PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE gator TO gatoruser;
\q

### 4. Test your connection

psql -U gatoruser -d gator -h localhost

If this connect successfully, you're ready to go!

## Installation:

Open terminal and run 'go install github.com/paul39-33/gator@latest'

## Setup:

### 1. Create the config file

Create a file called `.gatorconfig.json` in your home directory with the following content:

{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}

### 2. Update the database URL

Replace the db_url with your actual PostgreSQL connection string:

username: Your PostgreSQL username

password: Your PostgreSQL password

localhost:5432: Your PostgreSQL host and port

gator: Your database name

Example: {
  "db_url": "postgres://john:mypassword@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}

### 3. Run the program

Use "gator" + "command name" + "input"(if required) to run the program.
For example:

gator register yourusername

## Command lists

 - login {username} : login to the specified username if exists
 - register {username} : register new user with given username
 - reset : reset the program (THIS WILL DELETE ALL USER AND DATA)
 - users : list all registered users
 - agg {time_duration} : aggregate feeds that the user follows with the given time duration as its frequencies and save it as posts to be viewed.
 - addfeed {"feed_name"} {feed_url} : add new feed that will be automatically followed by the current user and can also be followed by other user using "follow {feed_url}"
 - feeds : list all feeds and their creators
 - follow {feed_url} : follow a feed. If the feed doesn't exist yet then it needs to be added using "addfeed {"feed_name"} {feed_url}"
 - following : list all feeds that's followed by the current user
 - unfollow {feed_url} : unfollow the given feed
 - browse {feed_limit}(optional) : shows all posts (saved feeds) for the current user. If no feed_limit is given then the default value is 2.
 


## **Additional notes:**

- **Default PostgreSQL user**: Usually `postgres` with no password on local installations
- **Port**: PostgreSQL typically runs on port 5432
- **Connection string example**: `postgres://gatoruser:yourpassword@localhost:5432/gator?sslmode=disable`

## **Troubleshooting section:**

### Troubleshooting PostgreSQL

**If you get "connection refused":**
- Make sure PostgreSQL is running: `brew services list | grep postgresql` (macOS)
- Check if it's listening on port 5432: `lsof -i :5432`

**If you get "authentication failed":**
- Check your username/password in the config file
- Try connecting with `psql` first to verify credentials

