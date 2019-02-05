import React, { Component } from "react";
import axios from "axios";

export default class RegForm extends Component {
  state = {
    email: "",
    password: "",
    firstName: "",
    lastName: ""
  };
  handleSubmit = event => {
    event.preventDefault();
    axios
      // .post(`https://api.github.com/users/${this.state.userName}`, JSON.stringify(this.state))
      .then(resp => {
        console.log(resp.data);
        this.setState({ email: "", password: "", firstName: "", lastName: "" });
      });
  };
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
          required
        />
        <input
          type="password"
          value={this.state.password}
          onChange={event => this.setState({ password: event.target.value })}
          placeholder="Password"
          required
        />
        <button type="submit">Register</button>
      </form>
    );
  }
}

// name, password, email, dogId
//id, name, age, breed, ownerTag
