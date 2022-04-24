# Example application demoing DB migrations

For people who don't have shell access to production

## Setup

Create a database

    $ psql -c 'CREATE DATABASE gomigrationexample'
	$ export DATABASE_URL=postgres://localhost/gomigrationexample

Start the webserver

    $ go run -v ./cmd/go-migration-example

## Usage

    $ curl localhost:8000/users

This fails until the database tables have been created.
Run the migrations by sending a POST request.

    $ curl -X POST localhost:8000/admin/migrations/apply

Then try the database operations again.

## Notes

It's a great idea to rewrite SQL that you run from database shells
as code. It's an even better idea to version control that code
and put it behind a validating API.

Don't limit yourself to migrations either, debugging databases from
the application, and driven by Web APIs is a great way to share
knowledge for the code in the code itself!

In real life, make sure to put these admin/debug pages behind authentication.
