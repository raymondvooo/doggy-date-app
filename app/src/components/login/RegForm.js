import React, { Component } from "react";
import axios from "axios";
import { v1 as uuid } from 'uuid';
import  gql from 'graphql-tag';
import validator from 'validator';
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
    dogAge: null,
    dogBreed: "",
    dogs: [],
    submitState: false,
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
    }
    const query = CREATE_USER; //mutation query

    if (!this.state.validEmail) {
      alert("Error. Please enter a valid email!");
    }
    else if (!this.state.validAge) {
      alert("Error. Please enter a number for your dog's age!")
    }
    else {
    axios.post('https://doggy-date-go.herokuapp.com/graphql', {query, variables})   // send post request to graphql endpoint with mutation query and variables
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
          validAge: false,
        });
      });
    }
  };

  checkEmail = () => {
    this.setState({
    validEmail: validator.isEmail(this.state.email)
    })
  }
  checkAge = () => {
    this.setState({
      validAge: validator.isInt(this.state.dogAge)
    })
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <input
          type="text"
          value={this.state.firstName} 
          onChange={event => this.setState({ firstName: event.target.value })}
          placeholder="First Name"
          required
        />
        <input
          type="text"
          value={this.state.lastName}
          onChange={event => this.setState({ lastName: event.target.value })}
          placeholder="Last Name"
          required
        />
        <input
          type="text"
          value={this.state.email}
          onChange={event => this.setState({ email: event.target.value })}
          placeholder="Email"
          onBlur={this.checkEmail}
          required
        />
        <input
          type="text"
          value={this.state.dogName}
          onChange={event => this.setState({ dogName: event.target.value })}
          placeholder="Dog's Name"
          onBlue={this.checkAge}
          required
        />
        <input
          type="text"
          value={this.state.dogAge}
          onChange={event => this.setState({ dogAge: event.target.value })}
          placeholder="Dog's Age"
          required
        />
        <input
          type="text"
          value={this.state.dogBreed}
          onChange={event => this.setState({ dogBreed: event.target.value })}
          placeholder="Dog's Breed"
          required
        />
        {/* <input
          type="password"
          value={this.state.password}
          onChange={event => this.setState({ password: event.target.value })}
          placeholder="Password"
          required
        /> */}
        <button type="submit" className="btn btn-primary">
          Register
        </button>
      </form>
    );
  }
}

// name, password, email, dogId
//id, name, age, breed, ownerTag
