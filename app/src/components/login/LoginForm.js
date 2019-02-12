import React, { Component } from "react";
import axios from "axios";
import { Button, Form, Input } from "semantic-ui-react";

const LOGIN_USER = `
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

export default class LoginForm extends Component {
  state = {
    email: "",
    password: "",
    id: null,
    name: "",
    dogId: null,
    dogName: "",
    dogAge: "",
    dogBreed: ""
  };

  handleSubmit = event => {
    event.preventDefault();
    axios
      .post("https://doggy-date-go.herokuapp.com/graphql", JSON.stringify())
      .then(resp => {
        console.log(resp.data);
        this.setState({ email: "", password: "" });
      });
  }


  render() {
    return (
      <div>
        <Form onSubmit={this.handleSubmit}>
          <Input
            type="text"
            value={this.state.email}
            onChange={event => this.setState({ email: event.target.value })}
            placeholder="Email"
            required
          />
          {/* <input
          type="password"
          value={this.state.password}
          onChange={event => this.setState({ password: event.target.value })}
          placeholder="Password"
          required
        /> */}
          <Button type="submit" primary>
            Login
          </Button>
        </Form>
        <Button secondary id="reg" onClick={this.props.showReg}>
          Need an Account?
        </Button>
      </div>
    );
  }
}
