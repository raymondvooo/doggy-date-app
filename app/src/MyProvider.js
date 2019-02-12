import React, { Component, createContext } from 'react';
import { LOGIN_USER } from './schema';
import axios from "axios";
export const MyContext = createContext();

export default class MyProvider extends Component {
    state = {
    };

    getLoginState = (email, validEmail) => {
        const variables = {
        email: email
        };
        const query = LOGIN_USER;
        if (!validEmail) {
        alert("Error. Please enter a valid email!");
        } else {
        axios
            .post("https://doggy-date-go.herokuapp.com/graphql", {
            query,
            variables
            }) // send post request to graphql endpoint with login query and variables
            .then(resp => {
            console.log("Graphql User response: ", resp.data);
            const userData = resp.data.data.loginUser;
            this.setState({
                userData
            });
            });
        }
    }
    render() {
        return (
            <MyContext.Provider value={this.getLoginState}>
                {this.props.children}
            </MyContext.Provider>
        )
    }
}