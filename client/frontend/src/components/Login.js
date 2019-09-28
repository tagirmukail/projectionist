import React from 'react';
import { connect } from 'react-redux';
import {
    TextField,
    FormControl,
    IconButton,
} from '@material-ui/core';
import MeetingRoomIcon from '@material-ui/icons/MeetingRoom';
import withWidth from '@material-ui/core/withWidth';
import { fetchLogin } from '../actions/login';

class Login extends React.Component {
    state = {
        ...this.props.login.result.data,
        username: null,
        password: null,
        error: '',
    }

    clickLogin = () => {
        // TODO: fix fetchLogin
        fetchLogin({
            username: this.state.username,
            password: this.state.password,
        });
    };

    render() {
        return (
            <FormControl style={{
                display: 'flex',
                flexWrap: 'wrap',
            }}
            >
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
                    onClick={(e) => this.clickLogin()}
                >
                    <MeetingRoomIcon />
                </IconButton>
            </FormControl>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        login: state.login,
    }
}

export default withWidth()(connect(
    mapStateToProps,
    { fetchLogin }
)(Login));