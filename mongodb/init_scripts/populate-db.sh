#!/bin/bash

# MongoDB connection details
MONGO_HOST="localhost"
MONGO_PORT="27017"
DATABASE_NAME="cinecrew"

# Connect to MongoDB
mongo_connect="mongosh --host ${MONGO_HOST} --port ${MONGO_PORT} ${DATABASE_NAME}"

# Insert data into 'people' collection
${mongo_connect} <<EOF
db.people.insertMany([
  { "name": "Tom Hanks", "birthday": ISODate("1956-07-09"), "profession": "actor" },
  { "name": "Steven Spielberg", "birthday": ISODate("1946-12-18"), "profession": "director" },
  { "name": "Meryl Streep", "birthday": ISODate("1949-06-22"), "profession": "actor" },
  { "name": "Christopher Nolan", "birthday": ISODate("1970-07-30"), "profession": "director" },
  { "name": "Brad Pitt", "birthday": ISODate("1963-12-18"), "profession": "actor" },
  { "name": "Leonardo DiCaprio", "birthday": ISODate("1974-11-11"), "profession": "actor"},
  { "name": "Morgan Freeman", "birthday": ISODate("1937-06-01"), "profession": "actor" },
  { "name": "Edward Norton", "birthday": ISODate("1969-08-18"), "profession": "actor" },
  { "name": "Christian Bale", "birthday": ISODate("1974-01-30"), "profession": "actor" },
  { "name": "Heath Ledger", "birthday": ISODate("1979-04-04"), "profession": "actor" }
])
EOF

# Insert data into 'genres' collection
${mongo_connect} <<EOF
db.genres.insertMany([
  { "name": "Drama" },
  { "name": "Comedy" },
  { "name": "Action" },
  { "name": "Sci-Fi" },
  { "name": "Thriller" },
  { "name": "Horror" },
  { "name": "Romance" },
  { "name": "Mystery" },
  { "name": "Crime" },
  { "name": "Adventure" }
])
EOF

# Insert data into 'movies' collection (referencing 'people' and 'genres')
${mongo_connect} <<EOF
db.movies.insertMany([
  {
    "name": "Forrest Gump",
    "genre": "Drama",
    "release_date": ISODate("1994-07-06"),
    "description": "The story of a man with a low IQ who accomplished great things.",
    "director": "Steven Spielberg",
    "cast": ["Tom Hanks", "Meryl Streep"],
    "rating": { "users_voted": 500, "score": 4.5 }
  },
  {
    "name": "Inception",
    "genre": "Sci-Fi",
    "release_date": ISODate("2010-07-16"),
    "description": "A thief who enters the dreams of others to steal their secrets.",
    "director": "Christopher Nolan",
    "cast": ["Leonardo DiCaprio"],
    "rating": { "users_voted": 400, "score": 4.7 }
  },
  {
    "name": "The Shawshank Redemption",
    "genre": "Drama",
    "release_date": ISODate("1994-10-14"),
    "description": "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
    "director": "Steven Spielberg",
    "cast": ["Tom Hanks", "Morgan Freeman"],
    "rating": { "users_voted": 450, "score": 4.8 }
  },
  {
    "name": "Fight Club",
    "genre": "Drama",
    "release_date": ISODate("1999-10-15"),
    "description": "An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.",
    "director": "Christopher Nolan",
    "cast": ["Brad Pitt", "Edward Norton"],
    "rating": { "users_voted": 550, "score": 4.6 }
  },
  {
    "name": "The Dark Knight",
    "genre": "Action",
    "release_date": ISODate("2008-07-18"),
    "description": "When the menace known as the Joker emerges from his mysterious past, he wreaks havoc and chaos on the people of Gotham.",
    "director": "Christopher Nolan",
    "cast": ["Christian Bale", "Heath Ledger"],
    "rating": { "users_voted": 600, "score": 4.9 }
  }
])
EOF

echo "Data insertion complete."