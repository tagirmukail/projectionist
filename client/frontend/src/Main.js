import React, {Component} from 'react';
import {Switch} from 'react-router-dom';
import Login from "./components/Login";
import Service from "./components/Service";
import Config from "./components/Config";
import AppliedRoute from "./components/AppliedRoute"

class Main extends Component {
    render() {
        return (!this.props.childProps.isAuthenticated
                ? <Switch><AppliedRoute path='/login' component={Login} props={this.props.childProps}/></Switch>
                : <Switch>
                    <AppliedRoute path='/' exact component={Service} props={this.props.childProps}/>
                    <AppliedRoute path='/configs' component={Config} props={this.props.childProps}/>
                  </Switch>
        );
    }
}

export default Main;