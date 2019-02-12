import React, { Component } from "react";
import "./App.css";
import AppRouter from "./Router";
import "semantic-ui-css/semantic.min.css";
import Login from "./components/login/Login";
import MyProvider from "./MyProvider";

// import { library } from '@fortawesome/fontawesome-svg-core'
// import { fab } from '@fortawesome/free-brands-svg-icons'
// import { faSync, faStar } from '@fortawesome/free-solid-svg-icons'

// library.add(fab, faStar, faSync)

const Index = () => <h2>Index</h2>;

class App extends Component {
  render() {
    return (
      <div>
        <MyProvider>
          <AppRouter />
        </MyProvider>
      </div>
    );
  }
}

export default App;
