import React, {Component} from 'react';
import { BrowserRouter } from 'react-router-dom';
import Header from './Header';
import Main from './Main';
import { instanceOf } from 'prop-types';
import { withCookies, Cookies } from 'react-cookie';

class App extends Component {
    static propTypes = {
        cookies: instanceOf(Cookies).isRequired
    };

    constructor(props) {
        super(props);

        const {cookies} = props;

        const session = cookies.get("session");

        this.state = {
            isAuthenticated: !!session,
            session: session || ""
        };
    }

    userHasAuthenticated = authenticated => {
        const {cookies} = this.props;

        let session = (authenticated) ? cookies.get("session"): "";
        let location  = (authenticated) ? "/": "/login";

        if (!authenticated) {cookies.remove('session');}

        this.setState({
            isAuthenticated: authenticated,
            session: session || ""
        });

        document.location = location;
    };

    render() {
        const childProps = {
            isAuthenticated: this.state.isAuthenticated,
            userHasAuthenticated: this.userHasAuthenticated,
        };

        return (
            <BrowserRouter>
                <div>
                    <Header childProps={childProps}/>
                    <hr/>
                    <Main childProps={childProps}/>
                </div>
            </BrowserRouter>
        );
    }
}

export default withCookies(App);
