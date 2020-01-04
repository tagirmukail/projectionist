import React, {Fragment} from 'react';
import ReactDOM from 'react-dom';
import { createStore, applyMiddleware, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';
import { Switch, Route, Link, BrowserRouter} from 'react-router-dom';
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import {AppBar} from "@material-ui/core";
import { routes, loginRoute, logoutRoute } from './router';
import './index.css';
import reducers from './reducers/reducers'
import {makeStyles} from "@material-ui/core/styles";

const store = createStore(combineReducers(reducers), applyMiddleware(thunk));

class App extends React.Component {
    useStyles () {
        return makeStyles(theme => ({
            root: {
                flexGrow: 1,
                backgroundColor: theme.palette.background.paper,
            },
            button: {
                
            },
        }));
    }

    render() {

        const useStyle = this.useStyles();

        return (<BrowserRouter>
            <div className="App">
                <Route
                    path="/"
                    render={({location}) => (
                        <Fragment>
                            <div className={useStyle.root}>
                                <AppBar position="static">
                                    <Tabs value={location.pathname}>
                                        {routes.map((route, i) =>(localStorage.getItem("token") &&
                                            <Tab key={i} label={route.name} value={route.url} component={Link} to={route.url}/>
                                            ))}
                                            {!localStorage.getItem("token") &&
                                            <Tab
                                                key={routes.length+1}
                                                label={loginRoute.name}
                                                value={loginRoute.url}
                                                component={Link}
                                                to={loginRoute.url}
                                            />}
                                        {localStorage.getItem("token") &&
                                        <Tab
                                            key={routes.length+2}
                                            label={logoutRoute.name}
                                            value={logoutRoute.url}
                                            component={Link}
                                            to={logoutRoute.url}
                                        />}
                                    </Tabs>
                                </AppBar>
                            </div>
                            <Switch>
                                {routes.map((route, i) =>(localStorage.getItem("token") &&
                                    <Route
                                        key={i}
                                        path={route.url}
                                        component={route.component}
                                    />)
                                )}
                                {!localStorage.getItem("token") && <Route
                                    key={routes.length + 1}
                                    path={loginRoute.url}
                                    component={loginRoute.component}
                                />}
                                {localStorage.getItem("token") && <Route
                                    key={routes.length + 1}
                                    path={logoutRoute.url}
                                    component={logoutRoute.component}
                                />}
                            </Switch>
                        </Fragment>
                    )}
                />
            </div>
        </BrowserRouter>);
    }
}

ReactDOM.render(
    <Provider store={store}>
        <App/>
    </Provider>, document.getElementById('root')
);