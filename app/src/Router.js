import React from 'react';
import { BrowserRouter as Router, Route, Link } from "react-router-dom";
import Login from "./components/login/Login.js"

const AppRouter = () => (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/login/">Login</Link>
            </li>
            {/* <li>
              <Link to="/cards/">Cards</Link>
            </li>
            <li>
              <Link to="/game/">Game</Link>
            </li> */}
          </ul>
        </nav>
  
        {/* <Route path="/" exact component={Index} /> */}
        <Route path="/login/" component={Login} />
        {/* <Route path="/cards/" component={Cards} />
        <Route path="/game/" component={Game} /> */}
      </div>
    </Router>
  );

  export default AppRouter;