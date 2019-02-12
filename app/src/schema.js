export const LOGIN_USER = `
query loginUser($email: String!) {
  loginUser(email: $email) {
    id
    name
    email
    profileImageURL
    dogs {
      id
      name
      age
      breed
      profileImageURL
    }
  }
}`;

// CREATE_USER mutation query
export const CREATE_USER = `
mutation CreateUser($id: ID!, $name: String!, $email: String!, $dogId: ID!, $dogName: String!, $dogAge: Int!, $dogBreed: String!) {
  createUser(
    id: $id,
    name: $name,
    email: $email,
    dogId: $dogId,
    dogName: $dogName,
    dogAge: $dogAge,
    dogBreed: $dogBreed
  ) {
    id
    name
    email
    dogs {
      id
      name
      age
      breed
    }
  }
}`;