import React from 'react';
import { connect } from 'react-redux';
import {
    TextField,
    FormControl,
    IconButton,
    LinearProgress,
    SnackbarContent,
    Icon,
    Container,
} from '@material-ui/core';
import {Redirect} from "react-router-dom";
import MeetingRoomIcon from '@material-ui/icons/MeetingRoom';
import withWidth from '@material-ui/core/withWidth';
import { fetchLogin } from '../actions/login';

function SnackbarContentWrapper(props) {
    const { message } = props;

    return (
        <SnackbarContent
            aria-describedby="login-snackbar"
            message={
                <span id="login-snackbar">
                    <Icon />
                    {message}
                </span>
            }
        />
    );
}

class Login extends React.Component {
    state = {
        ...this.props.login.result.data,
        username: null,
        password: null,
        error: '',
    };

    clickLogin = () =>{
       this.props.fetchLogin({
            username: this.state.username,
            password: this.state.password,
        });
    };


    render() {
        if (this.props.login.wait) {
            return (<LinearProgress/>)
        } else {
            if (localStorage.getItem("token")) {
                return <Redirect to='/users'/>
            }
        return (
            <Container maxWidth="sm">
                <FormControl
                    maxwidth='sm'
                    maxhight='sm'
                    style={{
                        marginTop: 150,
                        display: 'flex'
                       }}
                >
                    {localStorage.getItem("token") && <SnackbarContentWrapper
                        message={this.props.login.result.data.message}
                    />}
                    <TextField
                        inputProps={{
                            placeholder: "Username"
                        }}
                        id="username"
                        required
                        value={(this.state.username) ? this.state.username : ''}
                        onChange={(e) => this.setState({ username: e.target.value })}
                        margin="normal"
                    >
                        {this.state.username ? this.state.username : ''}
                    </TextField>
                    <TextField
                        inputProps={{
                            placeholder: "Password"
                        }}
                        value={this.state.password ? this.state.password : ''}
                        onChange={(e) => this.setState({ password: e.target.value })}
                        required
                        margin="normal"
                        id="password"
                        type="password"
                    >
                        {this.state.password ? this.state.password : ''}
                    </TextField>
                    <IconButton
                        onClick={this.clickLogin}
                    >
                        <MeetingRoomIcon />
                    </IconButton>
                </FormControl>
            </Container>
        )}
    }
}

const mapStateToProps = (state) => {
    return {
        login: state.login,
    }
};

export default withWidth(400)(connect(
    mapStateToProps,
    { fetchLogin }
)(Login));