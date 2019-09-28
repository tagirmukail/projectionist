import React from 'react';
import { connect } from 'react-redux';
import withWidth from '@material-ui/core/withWidth';
import { fetchUsers } from '../../actions/fetchUsers';

class Users extends React.Component {
    state = {
        error: '',
    }

    componentDidMount () {
        this.props.fetchUsers(0, 10);
    }

    render() {
        console.log(this.props.users.result);
        return (<div>Users</div>)
    }
}

const mapStateToProps = (state) => {
    return {
        users: state.users,
    }
}

export default withWidth()(connect(
    mapStateToProps,
    { fetchUsers }
)(Users));