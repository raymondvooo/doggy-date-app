# Doggy Date App
Connect with other dog owners and plan a doggy date!

## Server
The backend is written in go and uses graphql to query or create users and dogs. <br/>
Demo: https://doggy-date-go.herokuapp.com/<br/>
Copy and paste this code into the playground!
```
# Example user query
query {
  user(id: "2a3ce71c-2a8f-11e9-9fd2-22000b860eee") {
    name
    email
    dogs {
      name
      age
      breed
    }
  }
}
```
