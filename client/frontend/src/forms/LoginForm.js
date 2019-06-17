import React, {Component} from 'react';
import { Button, FormGroup, FormControl, FormLabel } from 'react-bootstrap';
import "./Login.css";

const LoginPath = "/v1/api/login";

class LoginForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            login: "",
            password: ""
        };
    }

    validateForm() {
        return this.state.login.length > 0 && this.state.password.length > 0;
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleSubmit = event => {
        event.preventDefault();


            let payload = JSON.stringify({
                "username": this.state.login,
                "password": this.state.password
            });

            fetch(LoginPath, {
                method: 'POST',
                credentials: "include",
                headers: {
                    'Access-Control-Allow-Origin': '*',
                    'Content-Type': 'application/json',
                },
                body: payload,
            })
                .then(results => results.json())
                .then(response => {
                    if (!response.status) {
                        console.log(response.message);
                        return;
                    }

                    this.props.userHasAuthenticated(true);
                })
    };

    render() {
        return (
            <div className="Login">
                <form onSubmit={this.handleSubmit}>
                    <FormGroup controlId='login'>
                        <FormLabel>Login</FormLabel>
                        <FormControl
                        autoFocus
                        type="text"
                        value={this.state.login}
                        onChange={this.handleChange}
                        />
                    </FormGroup>
                    <FormGroup controlId='password'>
                        <FormLabel>Password</FormLabel>
                        <FormControl
                            type="password"
                            value={this.state.password}
                            onChange={this.handleChange}
                        />
                    </FormGroup>
                    <Button
                    block
                    disabled={!this.validateForm()}
                    type="submit"
                    >Login</Button>
                </form>
            </div>
        );
    }
}

export default LoginForm;