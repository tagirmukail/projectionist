import React, {Component} from 'react';
import {Link} from 'react-router-dom';

class Header extends Component {
    handleClick = event => {
        event.preventDefault();

        this.props.childProps.userHasAuthenticated(false);
    };

    render() {
        return (
                <nav className="navbar navbar-expand-lg navbar-light bg-light">
                        {!this.props.childProps.isAuthenticated
                            ? <ul className="navbar-nav mr-auto">
                                <li><Link to={'/login'} className="nav-link"> Login</Link></li>
                              </ul>
                            : <ul className="navbar-nav mr-auto">
                                <li><Link to={'/'} className="nav-link"> Services</Link></li>
                                <li><Link to={'/configs'} className="nav-link"> Configs</Link></li>
                                <li><Link to={'/logout'} onClick={this.handleClick} className="nav-link"> Logout</Link></li>
                            </ul>
                        }
                </nav>
        )
    }

}

export default Header;