import React, {Component} from 'react';
import {Link} from 'react-router-dom';

class Header extends Component {
    render() {
        return (
                <nav className="navbar navbar-expand-lg navbar-light bg-light">
                    <ul className="navbar-nav mr-auto">
                        {!this.props.childProps.isAuthenticated
                            ? <li><Link to={'/'} className="nav-link"> Login</Link></li>
                            : ""
                        }
                        {this.props.childProps.isAuthenticated
                        ? <li><Link to={'/services'} className="nav-link"> Services</Link></li>
                            : ""
                        }
                        {this.props.childProps.isAuthenticated
                            ?<li><Link to={'/configs'} className="nav-link"> Configs</Link></li>
                            : ""
                        }
                    </ul>
                </nav>
        )
    }

}

export default Header;