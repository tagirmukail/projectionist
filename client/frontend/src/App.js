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

        this.setState({
            isAuthenticated: authenticated,
            session: cookies.get("session") || ""
        });
    };

    userClearAuthenticated = () => {
        const {cookies} = this.props;

        cookies.remove('session');

      this.setState({
          isAuthenticated: false,
          session: ""
      })
    };

    render() {
        const childProps = {
            isAuthenticated: this.state.isAuthenticated,
            userHasAuthenticated: this.userHasAuthenticated,
            userClearAuthenticated: this.userClearAuthenticated
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
