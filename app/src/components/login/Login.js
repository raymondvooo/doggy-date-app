import React, { Component } from 'react';
import '../../App.css';
import './Login.css';

class Login extends Component {
  render() {
    return (
      <div className="login-wrapper">
      <form className="loginForm">
      <section className="title">
      <h3 className="welcome">Welcome to</h3> Doggy Date!
      <h5 className="hint">Sign in below</h5>
      <p>
          <br />
          </p>
      </section>

      <div className="login-group">
         <input className="username" type="text" name="input" placeholder="Username"></input>
         <input className="password" type="password" name="input" placeholder="Password"></input>


         <button type="submit" className="btn btn-primary">Login</button>


         <button className="btn btn-link">Need an Account?</button>
    </div>
      
      </form>
</div>


    );
  }
}

export default Login;


// <div class="login-wrapper">
// <form class="login">
//     <section class="title">
//         <h3 class="welcome">Welcome to</h3>
//         Stock Search
//         <h5 class="hint" *ngIf="userService.uLogin === true && userService.id != null">You are already signed in!</h5>
//         <h5 class="hint" *ngIf="userService.uLogin === true && userService.id === null">Sign in below</h5>
//         <p>
//           <br>
//           <!-- if user is already signed in, create logout button -->
//         <p> <button *ngIf="userService.uLogin === true && userService.id != null" class="btn btn-link" (click)= "onLogout()">Logout</button>
//         <p>
            
//         <button *ngIf="userService.uLogin === true && userService.id != null" class="btn btn-link" (click)="userService.uRegister = true; userService.uLogin = false" >Need an Account?</button>
//         </p>
//     </section>
//     <!-- if user needs to login, display this -->
//     <div class="login-group" *ngIf="userService.uLogin === true && userService.id === null">
//         <input class="username" type="text" name="input" placeholder="Username" [(ngModel)]="user.email">
//         <input class="password" type="password" name="input" placeholder="Password" [(ngModel)]="user.password">

//         <div class="error active" *ngIf="loginS === false">
//             Invalid user name or password
//         </div>
//         <button type="submit" class="btn btn-primary" (click)="onLogin()">Login</button>


//         <button class="btn btn-link" (click)="userService.uRegister = true; userService.uLogin = false" >Need an Account?</button>
//     </div>
    
//     <!-- if user wants to register, display this -->
//     <h5 class="hint" *ngIf="userService.uRegister === true">Register below</h5>
//     <div class="login-group" *ngIf="userService.uRegister === true">
      
//         <input class="first name" type="text" name="input" placeholder="First Name" [(ngModel)]="user.firstName">
//         <input class="last name" type="text" name="input" placeholder="Last Name" [(ngModel)]="user.lastName">
//         <input class="email" type="text" name="input" placeholder="Email" [(ngModel)]="user.email">
//         <input class="password" type="password" name="input" placeholder="Password" [(ngModel)]="user.password">

       
//         <button type="submit" class="btn btn-primary" (click)="onRegister()">Register</button>

//         <button class="btn btn-link" (click)="userService.uRegister = false; userService.uLogin = true" >Go Back</button>
//     </div>
// </form>
// </div>