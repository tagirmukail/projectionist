import React from 'react';
import { connect } from 'react-redux';
import withWidth from '@material-ui/core/withWidth';
import { fetchUsers } from '../../actions/fetchUsers';
import {Redirect} from "react-router-dom";

import {
    Container
} from '@material-ui/core';

class Users extends React.Component {
    state = {
        error: '',
    };

    componentDidMount () {
        this.props.fetchUsers(1, 10);
    }

    render() {
        if (!localStorage.getItem("token")) {
            return <Redirect to="/login"/>
        }

        return (<Container>Users</Container>)
    }
}

const mapStateToProps = (state) => {
    return {
        users: state.users,
    }
};

export default withWidth()(connect(
    mapStateToProps,
    { fetchUsers }
)(Users));