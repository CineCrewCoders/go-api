## General Description
REST api for rating movies.

## Pre-requisites
-   [Golang](https://golang.org/dl/)
-   [Docker](https://docs.docker.com/engine/install/)

## How to run

```bash
# database
cd mongodb # go to mongodb directory
docker compose -f docker-compose.yml up # run the app in a container
docker compose -f docker-compose.yml down # shut down the container but keep the volumes
docker compose -f docker-compose.yml down --volumes # shut down the container and delete the volumes (the database will be empty)
```
```bash
# go api
cd src
go run .
```

## Endpoints (with postman examples)
- "/movies" (GET all movies  - **no need to be logged in**)
	```
	GET http://localhost:5678/movies
	```
- "/movies/:id" (GET a specific movie - **you have to be logged in**)
	```
	GET http://localhost:5678/movies/movieId
	```
- "/search" (GET a list of movies with any specified combination of this query params: title, min_average, genres - **no need to be logged in**)
	```
	GET http://localhost:5678/search?title=the martian&min_score=0.0&genres=Drama&genres=Comedy
	```
- "/signup" (POST - create a user with specified id (from firebase) and username -  **not logged in**)
	```
	POST http://localhost:5678/signup
	body:
	{
		"userID": "66217079769296d67c049493",
		"username": "darlena"
	}
	```
- "/user/:username" (GET the details for the specified user - **testing only for safety reasons**)
	```
	GET http://localhost:5678/user/darlena
	```
- "/movies/list" (POST - add movie to list (ids in body) **you have to be logged in**)
	```
	POST http://localhost:5678/movies/list
	body:
	{
		"list": "PlanToWatch",
		"movieId": "6622e243b51d6cd2b75ed6ff"
	}
	```
- "/movies/plan_to_watch" (GET plan to watch list **you have to be logged in**)
	```
	GET http://localhost:5678/movies/plan_to_watch
	```
- "/movies/watched" (GET watched list **you have to be logged in**)
	```
	GET http://localhost:5678/movies/watched
	```
- "/movies/plan_to_watch" (DELETE movie from plan to watch list (movieId as query param) **you have to be logged in**)
	```
	DELETE http://localhost:5678/movies/plan_to_watch?movieId=movieId
	```
 - "/movies/watched" (DELETE movie from watched list (movieId as query param) **you have to be logged in**)
	```
	DELETE http://localhost:5678/movies/watched?movieId=movieId
	```
- "/movies/rate" (POST - rate a movie (movieId and grade in body) **you have to be logged in**)
	```
	POST http://localhost:5678/movies/rate
	body:
	{
		"movieId": "6622e243b51d6cd2b75ed745",
		"score": 8
	}
	```
- "/movies/rate" (PUT - modify a grade **you have to be logged in**)
	```
	PUT http://localhost:5678/movies/rate
	body:
	{
		"movieId": "6622e243b51d6cd2b75ed745",
		"score": 8.5
	}
	```
!!! The UserId must be included in header for the requests that need a logged in user : Key: UserId, Value: 66217079769296d67c049493

#### A movie has the following details:
- title
- year
- runtime
- genres
- actors
- director
- plot
- poster
- average and number of votes