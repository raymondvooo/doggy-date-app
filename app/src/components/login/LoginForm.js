import React, { Component } from "react";
import axios from "axios";

const url = 'https://doggy-date-go.herokuapp.com/';


export default class LoginForm extends Component {
  state = {
    email: "",
    password: "",
    id: null,
    name: "",
    dogId: null,
    dogName:"",
    dogAge: null,
    dogBreed:""
  };
  

  
  handleSubmit = event => {
    event.preventDefault();
    axios
      .post(url, JSON.stringify())
      .then(resp => {
        console.log(resp.data);
        this.setState({ email: "", password: ""});
      });
  };
  render() {
    return (
        <div>
      <form onSubmit={this.handleSubmit}>
        <input
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
        <button type="submit" className="btn btn-primary">Login</button>
      </form>
      <button className="btn btn-link">Need an Account?</button>
        </div>
    );
  }
}


