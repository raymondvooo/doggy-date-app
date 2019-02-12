import React, { Component } from "react";
import RegForm from "./RegForm";
import LoginForm from "./LoginForm";
import "../../App.css";
import "./Login.css";

class Login extends Component {
  state = {
    showReg: false,
    showLog: true
  };

  toggleReg = () => {
    if (this.state.showReg === false) {
      this.setState({
        showReg: true,
        showLog: false
      });
    } else
      this.setState({
        showReg: false,
        showLog: true
      });
  };

  toggleLog = () => {
    if (this.state.showLog === false) {
      this.setState({
        showReg: false,
        showLog: true
      });
    } else
      this.setState({
        showReg: true,
        showLog: false
      });
  };

  render() {
    return (
      <div className="login-wrapper">
        <section className="title">
          <h3 className="welcome">Welcome to</h3> Doggy Date!
          <h5 className="hint">Sign in below</h5>
          <p>
            <br />
          </p>
        </section>

        <div className="login-group">
          {this.state.showReg ? (
            <RegForm showLog={this.toggleLog} />
          ) : (
            <LoginForm showReg={this.toggleReg} />
          )}
          {/* {this.state.showLog ? <LoginForm showReg={this.toggleReg} /> : <RegForm showLog={this.toggleLog}/>} */}
        </div>
      </div>
    );
  }
}

export default Login;
