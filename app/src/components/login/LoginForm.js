import React, { Component, createContext} from "react";
import { Button, Form, Input } from "semantic-ui-react";
import validator from "validator";
import { MyContext } from "../../MyProvider";


export default class LoginForm extends Component {
  state = {
    email: "",
    password: "",
    id: null,
    name: "",
    dogs: [],
    validEmail: false
  };

  checkEmail = () => {
    this.setState({
      validEmail: validator.isEmail(this.state.email)
    });
  };


  render() {
    return (
      <MyContext.Consumer>
      {value =>  (       
        <React.Fragment>
          <Form onSubmit={(e) => {
            e.preventDefault();
            value(this.state.email, this.state.validEmail);
            }}>
            <Input
              type="text"
              value={this.state.email}
              onChange={event => this.setState({ email: event.target.value })}
              placeholder="Email"
              onBlur={this.checkEmail}
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
        </React.Fragment>
        )}
        </MyContext.Consumer>

    );
  }
}
