import React, {Component} from 'react';
import {Switch} from 'react-router-dom';
import Login from "./components/Login";
import Service from "./components/Service";
import Config from "./components/Config";
import AppliedRoute from "./components/AppliedRoute"

class Main extends Component {
    render() {
        return (
            <Switch>
                {!this.props.childProps.isAuthenticated
                    ? <AppliedRoute exact path='/' component={Login} props={this.props.childProps}/>
                    : ""
                }
                {this.props.childProps.isAuthenticated
                ?<AppliedRoute path='/services' component={Service} props={this.props.childProps}/>
                :""
                }
                {this.props.childProps.isAuthenticated
                    ? <AppliedRoute path='/configs' component={Config} props={this.props.childProps}/>
                    : ""
                }
            </Switch>
        );
    }
}

export default Main;