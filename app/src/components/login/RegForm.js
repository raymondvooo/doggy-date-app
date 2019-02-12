import React, { Component } from "react";
import axios from "axios";
import { v1 as uuid } from "uuid";
import gql from "graphql-tag";
import validator from "validator";
import { Button, Form, Input } from "semantic-ui-react";

// CREATE_USER mutation query
const CREATE_USER = `
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

export default class RegForm extends Component {
  state = {
    firstName: "",
    lastName: "",
    email: "",
    dogName: "",
    dogAge: "",
    dogBreed: "",
    dogs: [],
    validEmail: false,
    validAge: false
  };

  handleSubmit = event => {
    event.preventDefault();
    // input variables for mutation query
    const variables = {
      id: uuid(),
      name: this.state.firstName + " " + this.state.lastName,
      email: this.state.email,
      dogId: uuid(),
      dogAge: parseInt(this.state.dogAge),
      dogName: this.state.dogName,
      dogBreed: this.state.dogBreed
    };
    const query = CREATE_USER; //mutation query

    if (!this.state.validEmail) {
      alert("Error. Please enter a valid email!");
    } else if (!this.state.validAge) {
      alert("Error. Please enter a number for your dog's age!");
    } else {
      axios
        .post("https://doggy-date-go.herokuapp.com/graphql", {
          query,
          variables
        }) // send post request to graphql endpoint with mutation query and variables
        .then(resp => {
          console.log("Graphql User response: ", resp.data);
          const userData = resp.data.createUser;
          this.setState({
            firstName: "",
            lastName: "",
            email: "",
            dogName: "",
            dogAge: "",
            dogBreed: "",
            dogs: [],
            validEmail: false,
            validAge: false
          });
        });
    }
  };

  checkEmail = () => {
    this.setState({
      validEmail: validator.isEmail(this.state.email)
    });
  };
  checkAge = () => {
    this.setState({
      validAge: validator.isInt(this.state.dogAge)
    });
  };

  render() {
    return (
      <div>
        <Form onSubmit={this.handleSubmit}>
          <Input
            type="text"
            value={this.state.firstName}
            onChange={event => this.setState({ firstName: event.target.value })}
            placeholder="First Name"
            required
          />
          <Input
            type="text"
            value={this.state.lastName}
            onChange={event => this.setState({ lastName: event.target.value })}
            placeholder="Last Name"
            required
          />
          <Input
            type="text"
            value={this.state.email}
            onChange={event => this.setState({ email: event.target.value })}
            placeholder="Email"
            onBlur={this.checkEmail}
            required
          />
          <Input
            type="text"
            value={this.state.dogName}
            onChange={event => this.setState({ dogName: event.target.value })}
            placeholder="Dog's Name"
            onBlur={this.checkAge}
            required
          />
          <Input
            type="text"
            value={this.state.dogAge}
            onChange={event => this.setState({ dogAge: event.target.value })}
            placeholder="Dog's Age"
            required
          />
          <Input
            type="text"
            value={this.state.dogBreed}
            onChange={event => this.setState({ dogBreed: event.target.value })}
            placeholder="Dog's Breed"
            required
          />

          <Button primary type="submit">
            Register
          </Button>
        </Form>
        <Button secondary id="reg" onClick={this.props.showLog}>
          Already Have an Account? Login Here.
        </Button>
      </div>
    );
  }
}

// name, password, email, dogId
//id, name, age, breed, ownerTag
