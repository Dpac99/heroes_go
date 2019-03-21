# How to run

* Make sure you have **golang** installed.
* Clone and cd into what we'll call **$DIR**
* Run `mongod --dbpath="./db"`. This will start your database.
* The command above will ocupy your terminal unless you run it with `--fork`. So, to make it simpler, just open another terminal and cd onto $DIR
* cd into $DIR/src/heroes.
* Run `go build`. It will generate an executable file called **heroes**.
* And you guessed it, run it with `./heroes`.

# How to test

   Swagger is not yet implemented, so for now I recommend using an API tester such as Postman.

# Interaction
  
 Url | Method | Description | Body | Response
 --- | --- | --- | --- | ---
 `/hero` | **GET** | Gets all heroes | none | An array of heroes
 `/hero` | **POST** | Creates a new hero | {"id" : int, "name": string} | The new hero
 `/hero/{id}` | **GET** | Gets a hero indexed by id | none | An hero
 `/hero/{id}` | **DELETE** | Deletes an hero | none | The deleted hero
 `/hero/{id}` | **PUT** | Updates an hero's name | {"name": string} | The updated hero

---
### That's it, have fun!

Any bugs found please report them
